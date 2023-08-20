package api

import (
	"log"

	db "random-service/internal/database"
	"random-service/internal/models"
)

type (
	RPCServer struct {
		db   db.DB
		name string
	}

	RPCPayload struct {
		Data string
	}

	RPCResponse struct {
		Method string
		Data   string
	}
)

func NewRPCServer(name string, db db.DB) *RPCServer {
	return &RPCServer{
		name: name,
		db:   db,
	}
}

func (rpc *RPCServer) RandomRPC(payload RPCPayload, res *RPCResponse) error {
	entry := models.RandomData{
		Data: payload.Data,
	}
	err := rpc.db.Insert(entry)
	if err != nil {
		log.Print("RPC Random service failed to store data")
		return err
	}

	res.Data = "RPC response from Random Service"
	res.Method = "RPC"
	return nil
}
