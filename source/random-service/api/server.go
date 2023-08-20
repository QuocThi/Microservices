package api

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/go-chi/chi/middleware"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-logr/logr"
	"google.golang.org/grpc"

	db "random-service/internal/database"
	api_proto "random-service/proto"
)

type Server struct {
	Log  logr.Logger
	DB   *db.DB
	Name string
}

func NewServer(db *db.DB, name string, l logr.Logger) Server {
	return Server{
		Log:  l,
		DB:   db,
		Name: name,
	}
}

func (app *Server) Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/random", app.SaveRandom)

	return mux
}

func (app *Server) RegisterGPRC(port, name string) {
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	randomServer := NewRandomGRPCServer(app.DB, name, app.Log)
	grpcServer := grpc.NewServer()
	api_proto.RegisterRandomServiceServer(grpcServer, randomServer)

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("failed to serve GRPC: %v", err)
	}
}

func (app *Server) ListenRPC(rpcPort string) {
	listenURL := fmt.Sprintf("0.0.0.0:%s", rpcPort)
	listen, err := net.Listen("tcp", listenURL)
	if err != nil {
		app.Log.Error(err, "rpc listen failed")
		panic("failed to listen rpc...")
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			continue
		}

		go rpc.ServeConn(conn)
	}
}
