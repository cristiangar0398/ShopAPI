package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/cristiangar0398/ShopAPI/middleware"
	"github.com/cristiangar0398/ShopAPI/models"
	"github.com/cristiangar0398/ShopAPI/repository"
	"github.com/cristiangar0398/ShopAPI/server"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
)

type UpsertPostRequest struct {
	PostContent string `json:"post_content"`
}

type PostResponse struct {
	Id          string `json:"id`
	PostContent string `json:"post_content"`
}

type PostUpdateResponse struct {
	Message string `json:"message"`
}

func InsertPostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := middleware.TokenParseString(w, s, r)
		if err != nil {
			log.Fatal(err)
		}
		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			var postRequest = UpsertPostRequest{}
			if err := json.NewDecoder(r.Body).Decode(&postRequest); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			id, err := ksuid.NewRandom()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			post := models.Post{
				Id:           id.String(),
				Post_content: postRequest.PostContent,
				UserId:       claims.UserId,
			}

			err = repository.InsertPost(r.Context(), &post)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var postMessage = models.WebsocketMessage{
				Type:    "Post_Created",
				Payload: post,
			}

			s.Hub().Broadcast(postMessage, nil)
			w.Header().Set("content-type", "aaplication/json")
			json.NewEncoder(w).Encode(PostResponse{
				Id:          post.Id,
				PostContent: post.Post_content,
			})
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetPostByIdHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		post, err := repository.GetPostById(r.Context(), params["id"])
		fmt.Println("post_content :", post.Post_content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "aaplication/json")
		json.NewEncoder(w).Encode(post)
	}
}

func UpdatePostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := middleware.TokenParseString(w, s, r)
		if err != nil {
			log.Fatal(err)
		}
		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			var postRequest = UpsertPostRequest{}
			if err := json.NewDecoder(r.Body).Decode(&postRequest); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			params := mux.Vars(r)
			Newpost := models.Post{
				Id:           params["id"],
				Post_content: postRequest.PostContent,
				UserId:       claims.UserId,
			}

			err = repository.UpdatePost(r.Context(), &Newpost)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("content-type", "aaplication/json")
			json.NewEncoder(w).Encode(PostUpdateResponse{
				Message: "Post Update",
			})
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func ListPostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		pageStr := r.URL.Query().Get("page")
		var page = uint64(0)
		if pageStr != "" {
			page, err = strconv.ParseUint(pageStr, 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		posts, err := repository.ListPost(r.Context(), page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}

func DeletePostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := middleware.TokenParseString(w, s, r)
		if err != nil {
			log.Fatal(err)
		}
		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			params := mux.Vars(r)
			err = repository.DeletePost(r.Context(), params["postId"], claims.UserId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("content-type", "aaplication/json")
			json.NewEncoder(w).Encode(PostUpdateResponse{
				Message: "Delete Post",
			})
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
