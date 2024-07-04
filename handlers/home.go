package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cristiangar0398/ShopAPI/server"
)

type HomeResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func HomeHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(HomeResponse{
			Message: "Welcome ShopAPI",
			Status:  http.StatusOK,
		})

	}
}
