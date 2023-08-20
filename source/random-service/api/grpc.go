package api

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"

	db "random-service/internal/database"
	"random-service/internal/models"
	api_proto "random-service/proto"
)

type RandomServer struct {
	api_proto.UnimplementedRandomServiceServer
	DB   *db.DB
	Name string
	Log  logr.Logger
}

func NewRandomGRPCServer(DB *db.DB, name string, l logr.Logger) *RandomServer {
	return &RandomServer{
		DB:   DB,
		Name: name,
		Log:  l,
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
		r.Log.Error(err, "failed to insert db")
		return &api_proto.RandomResponse{Result: false}, fmt.Errorf("failed to store random data: %v", err)
	}

	res := api_proto.RandomResponse{
		Method: "GRPC",
		Data:   "invoked GRPC successfully",
		Result: false,
	}
	return &res, nil
}
