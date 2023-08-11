package main

import (
	"log"

	"random-service/cmd/database"
	"random-service/cmd/database/models"
)

type (
	RPCServer struct {
		db   database.DB
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

func NewRPCServer(name string, db database.DB) *RPCServer {
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
