package api

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	db "log-service/internal/database"
)

// RPCServer is the type for our RPC Server. Methods that take this as a receiver are available
// over RPC, as long as they are exported.
type RPCServer struct {
	client *mongo.Client
}

// RPCPayload is the type for data we receive from RPC
type RPCPayload struct {
	Name string
	Data string
}

func NewRPCServer(client *mongo.Client) *RPCServer {
	return &RPCServer{
		client: client,
	}
}

// LogInfo writes our payload to mongo
func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	collection := r.client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), db.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error writing to mongo", err)
		return err
	}

	// resp is the message sent back to the RPC caller
	*resp = "Processed payload via RPC:" + payload.Name
	return nil
}
