package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"log-service/api"
	"log-service/internal/database"
	mylog "log-service/pkg/log"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

func main() {
	l := mylog.NewCustomLogger()
	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		l.Error(err, "connect to mongo failed")
		panic(err)
	}

	app := api.NewConfig(db.New(mongoClient), l)
	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = app.DB.GetClient().Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	rpcServer := api.NewRPCServer(mongoClient)
	// Register RPC Server
	err = rpc.Register(rpcServer)
	if err != nil {
		panic(err)
	}
	go app.RPCListen(rpcPort)

	go app.GRPCListen(gRpcPort)

	// start web server
	app.Log.Info("Starting", "port", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.Routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic()
	}
}

func connectToMongo() (*mongo.Client, error) {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	return c, nil
}
