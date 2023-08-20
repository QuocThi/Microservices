package api

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-logr/logr"
	"google.golang.org/grpc"

	db "log-service/internal/database"
	logs "log-service/proto"
)

type Config struct {
	DB  db.DB
	Log logr.Logger
}

func NewConfig(db db.DB, l logr.Logger) Config {
	return Config{
		DB:  db,
		Log: l,
	}
}

func (app *Config) Routes() http.Handler {
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

	mux.Post("/log", app.WriteLog)

	return mux
}

func (app *Config) GRPCListen(grpcPort string) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	s := grpc.NewServer()

	logs.RegisterLogServiceServer(s, NewGRPCLogServer(app.DB, app.Log))

	app.Log.Info("gRPC Server", "port", grpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}

func (app *Config) RPCListen(rpcPort string) error {
	app.Log.Info("Starting RPC", "port", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}
