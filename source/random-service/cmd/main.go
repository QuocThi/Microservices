package main

import (
	"fmt"
	"net/http"
	"net/rpc"

	"random-service/api"
	"random-service/internal/database"
	mylog "random-service/pkg/log"
)

const (
	webPort  = "80"
	rpcPort  = "5002"
	gRPCPort = "50001"
	mongoURL = "mongodb://mongo:27017"
)

func main() {
	l := mylog.NewCustomLogger()
	l.Info("connecting to mongodb ...")
	db, err := db.NewDB("random", mongoURL)
	if err != nil {
		panic(err)
	}
	l.Info("connected to mongodb")

	// close db connection
	defer func() {
		if err = db.Close(); err != nil {
			panic(err)
		}
	}()

	app := api.NewServer(db, "random-service", l)

	rpcServer := api.NewRPCServer("randomRPC", *app.DB)
	// register to RPC
	err = rpc.Register(rpcServer)
	if err != nil {
		panic("failed to register RPC for RPCServer reveiver")
	}
	app.Log.Info("register RPC successed")

	go app.ListenRPC(rpcPort)
	app.Log.Info("started RPC server", "port", rpcPort)

	go app.RegisterGPRC(gRPCPort, "RandomGRPC")
	app.Log.Info("started GRPC server successfully")

	// Listen normal requests
	server := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.Routes(),
	}

	err = server.ListenAndServe()
	if err != nil {
		panic("failed to start random-service")
	}
	app.Log.Info("start http server listen on %s \n", webPort)
}
