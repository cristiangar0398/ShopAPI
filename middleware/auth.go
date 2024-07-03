package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/cristiangar0398/ShopAPI/models"
	"github.com/cristiangar0398/ShopAPI/server"
	"github.com/golang-jwt/jwt"
)

var (
	NO_AUTH_NEEDED = []string{
		"login",
		"signup",
		"/",
	}
)

func shoulCheckYocken(route string) bool {

	for _, p := range NO_AUTH_NEEDED {
		if strings.Contains(route, p) {
			return false
		}
	}

	return true
}

func CheckAuthMiddleware(s server.Server) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !shoulCheckYocken(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			_, err := TokenParseString(w, s, r)
			if err != nil {
				log.Fatal(err)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func TokenParseString(w http.ResponseWriter, s server.Server, r *http.Request) (*jwt.Token, error) {
	tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
	jwtSecret := s.Config().JWTSecret

	if jwtSecret == "" {
		return nil, errors.New("failed to retrieve JWT secret from configuration")
	}
	token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(s.Config().JWTSecret), nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token is not valid")
	}

	return token, nil
}
