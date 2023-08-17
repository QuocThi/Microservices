package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/go-logr/logr"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"log-service/data"
	mylog "log-service/pkg/log"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	DB  data.DB
	Log logr.Logger
}

func NewConfig(db data.DB, l logr.Logger) Config {
	return Config{
		DB:  db,
		Log: l,
	}
}

func main() {
	l := mylog.NewCustomLogger()
	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		l.Error(err, "connect to mongo failed")
		panic(err)
	}

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := NewConfig(data.New(mongoClient), l)
	// Register RPC Server
	err = rpc.Register(new(RPCServer))
	if err != nil {
		panic(err)
	}
	go app.rpcListen()

	go app.gRPCListen()

	// start web server
	app.Log.Info("Starting", "port", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic()
	}
}

func (app *Config) rpcListen() error {
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
