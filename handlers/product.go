package handlers

import (
	"encoding/json"
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
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ImageUrl    string  `json"image_url"`
	Price       float64 `json:"price"`
}

type PostResponse struct {
	Id          string  `json:"id`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ImageUrl    string  `json"image_url"`
	Price       float64 `json:"price"`
}

type PostUpdateResponse struct {
	Message string `json:"message"`
}

func InsertProducttHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := middleware.TokenParseString(w, s, r)
		if err != nil {
			log.Fatal(err)
		}
		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			var productRequest = UpsertPostRequest{}
			if err := json.NewDecoder(r.Body).Decode(&productRequest); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			id, err := ksuid.NewRandom()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			Product := models.Products{
				Id:          id.String(),
				Title:       productRequest.Title,
				Description: productRequest.Description,
				ImageUrl:    productRequest.ImageUrl,
				Price:       productRequest.Price,
				UserId:      claims.UserId,
			}

			err = repository.InsertProduct(r.Context(), &Product)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("content-type", "aaplication/json")
			json.NewEncoder(w).Encode(PostResponse{
				Id:          Product.Id,
				Title:       productRequest.Title,
				Description: productRequest.Description,
				ImageUrl:    productRequest.ImageUrl,
				Price:       productRequest.Price,
			})
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetProductByIdHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		product, err := repository.GetProductById(r.Context(), params["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "aaplication/json")
		json.NewEncoder(w).Encode(product)
	}
}

func UpdateProducttHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := middleware.TokenParseString(w, s, r)
		if err != nil {
			log.Fatal(err)
		}
		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			var productRequest = UpsertPostRequest{}
			if err := json.NewDecoder(r.Body).Decode(&productRequest); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			params := mux.Vars(r)
			NewProduct := models.Products{
				Id:          params["id"],
				Title:       productRequest.Title,
				Description: productRequest.Description,
				ImageUrl:    productRequest.ImageUrl,
				Price:       productRequest.Price,
				UserId:      claims.UserId,
			}

			err = repository.UpdateProduct(r.Context(), &NewProduct)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("content-type", "aaplication/json")
			json.NewEncoder(w).Encode(PostUpdateResponse{
				Message: "Product Update",
			})
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func ListProductHandler(s server.Server) http.HandlerFunc {
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
		products, err := repository.ListProducts(r.Context(), page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}

func DeleteProductHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := middleware.TokenParseString(w, s, r)
		if err != nil {
			log.Fatal(err)
		}
		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			params := mux.Vars(r)
			err = repository.DeleteProduct(r.Context(), params["id"], claims.UserId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("content-type", "aaplication/json")
			json.NewEncoder(w).Encode(PostUpdateResponse{
				Message: "Delete product",
			})
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
