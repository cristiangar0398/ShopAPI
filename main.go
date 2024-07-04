package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/cristiangar0398/ShopAPI/handlers"
	"github.com/cristiangar0398/ShopAPI/middleware"
	"github.com/cristiangar0398/ShopAPI/server"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var (
	nombre string
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")
	JMT_SECRET := os.Getenv("JMT_SECRET")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	s, err := server.NewServer(context.Background(), &server.Config{
		Port:        PORT,
		JWTSecret:   JMT_SECRET,
		BatabaseUrl: DATABASE_URL,
	})

	if err != nil {
		log.Fatal(err)
	}

	s.Start(BindRoutes)

}

func BindRoutes(s server.Server, r *mux.Router) {

	r.Use(middleware.CheckAuthMiddleware(s))

	r.HandleFunc("/", handlers.HomeHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/signup", handlers.SignUpHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/login", handlers.LoginHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/me", handlers.MeHandler(s)).Methods(http.MethodGet)

	r.HandleFunc("/product", handlers.InsertProducttHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/product/{id}", handlers.GetProductByIdHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/product/{id}", handlers.UpdateProducttHandler(s)).Methods(http.MethodPut)
	r.HandleFunc("/product/{id}", handlers.DeleteProductHandler(s)).Methods(http.MethodDelete)
	r.HandleFunc("/product", handlers.ListProductHandler(s)).Methods(http.MethodGet)

	r.HandleFunc("/ws", s.Hub().HandleWebSockey)
}
