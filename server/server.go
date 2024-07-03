package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/cristiangar0398/ShopAPI/database"
	"github.com/cristiangar0398/ShopAPI/repository"
	"github.com/cristiangar0398/ShopAPI/websocket"
	"github.com/gorilla/mux"
)

type Config struct {
	Port        string
	JWTSecret   string
	BatabaseUrl string
}

type Server interface {
	Config() *Config
	Hub() *websocket.Hub
}

type Broker struct {
	config *Config
	router *mux.Router
	hub    *websocket.Hub
}

func (b *Broker) Config() *Config {
	return b.config
}

func (b *Broker) Hub() *websocket.Hub {
	return b.hub
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("port is required")
	}

	if config.JWTSecret == "" {
		return nil, errors.New("Secret is required")
	}

	if config.BatabaseUrl == "" {
		return nil, errors.New("Data Base is required")
	}

	broker := &Broker{
		config: config,
		router: mux.NewRouter(),
		hub:    websocket.NewHub(),
	}

	return broker, nil
}

func (b *Broker) Start(bainder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter()
	bainder(b, b.router)

	repo, err := database.NewPostgresRepository(b.config.BatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	go b.hub.Run()
	repository.SetRepository(repo)
	port := b.Config().Port

	log.Println(">>> >>> >>> 🚀 El servidor está despegando en el puerto", port, ">>> >>> >>>")
	if err := http.ListenAndServe(b.config.Port, b.router); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
