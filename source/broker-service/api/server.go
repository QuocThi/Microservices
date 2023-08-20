package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-logr/logr"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Server struct {
	Rabbit *amqp.Connection
	Log    logr.Logger
}

func NewServer(rabbit *amqp.Connection, l logr.Logger) Server {
	return Server{
		Rabbit: rabbit,
		Log:    l,
	}
}

func (app *Server) Routes() http.Handler {
	mux := chi.NewRouter()

	// specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/", app.Broker)

	mux.Post("/log-grpc", app.LogViaGRPC)
	mux.Post("/random-grpc", app.CallRandomGRPC)

	mux.Post("/handle", app.HandleSubmission)

	return mux
}
