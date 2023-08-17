package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	api_proto "random-service/api-proto"
	"random-service/cmd/database"
	"random-service/cmd/database/models"
)

type RandomServer struct {
	api_proto.UnimplementedRandomServiceServer
	DB   *database.DB
	name string
}

func NewRandomGRPCServer(DB *database.DB, name string) *RandomServer {
	return &RandomServer{
		DB:   DB,
		name: name,
	}
}

func (r *RandomServer) RandomGPRC(ctx context.Context, request *api_proto.RandomRequest) (*api_proto.RandomResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("empty request")
	}

	entry := models.RandomData{
		Data: request.Data,
	}

	err := r.DB.Insert(entry)
	if err != nil {
		return &api_proto.RandomResponse{Result: false}, fmt.Errorf("failed to store random data: %v", err)
	}

	res := api_proto.RandomResponse{
		Method: "GRPC",
		Data:   "invoked GRPC successfully",
		Result: false,
	}
	return &res, nil
}

func (app *Config) RegisterGPRC(port, name string) {
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	randomServer := NewRandomGRPCServer(app.db, name)
	grpcServer := grpc.NewServer()
	api_proto.RegisterRandomServiceServer(grpcServer, randomServer)

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("failed to serve GRPC: %v", err)
	}
}
