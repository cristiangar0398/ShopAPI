package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cristiangar0398/ShopAPI/middleware"
	"github.com/cristiangar0398/ShopAPI/models"
	"github.com/cristiangar0398/ShopAPI/repository"
	"github.com/cristiangar0398/ShopAPI/server"
	"github.com/golang-jwt/jwt"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	request  SignUpRequest
	response SignUpResponse
)

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func SignUpHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		request, err := decodeSignUpRequest(r)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		userID, err := generateUserID()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, err := createUser(r.Context(), request, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-Type", "application/json")
		json.NewEncoder(w).Encode(SignUpResponse{
			Id:    user.Id,
			Email: user.Email,
		})
	}
}

func decodeSignUpRequest(r *http.Request) (SignUpRequest, error) {
	var request SignUpRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	return request, err
}

func generateUserID() (string, error) {
	id, err := ksuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func createUser(ctx context.Context, request SignUpRequest, userID string) (*models.User, error) {
	//traemos la variable del archivo .env para hacer el hash del pass
	hashCostStr := os.Getenv("HASH_COST")
	hashCost, err := strconv.Atoi(hashCostStr)
	if err != nil {
		return nil, fmt.Errorf("Error al convertir HASH_COST a entero:", err)
	}

	//logica de creacion de usuario
	isRegistered, err := isEmailRegistered(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	if isRegistered {
		return nil, fmt.Errorf("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), hashCost)

	user := &models.User{
		Email:    request.Email,
		Password: string(hashedPassword),
		Id:       userID,
	}
	err = repository.InsertUser(ctx, user)
	return user, err
}

func isEmailRegistered(ctx context.Context, email string) (bool, error) {
	user, err := repository.GetUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}

func LoginHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := repository.GetUserByEmail(r.Context(), request.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if user == nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
		if err != nil {
			http.Error(w, "invalid credential ", http.StatusUnauthorized)
			return
		}

		claims := models.AppClaims{
			UserId: user.Id,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(2 * time.Hour * 24).Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(s.Config().JWTSecret))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", "aaplication/json")
		json.NewEncoder(w).Encode(LoginResponse{
			Token: tokenString,
		})
	}
}

func MeHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := middleware.TokenParseString(w, s, r)
		if err != nil {
			log.Fatal(err)
		}

		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			user, err := repository.GetUserById(r.Context(), claims.UserId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("content-type", "aaplication/json")
			json.NewEncoder(w).Encode(user)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
