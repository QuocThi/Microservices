package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"

	"github.com/go-logr/logr"

	"random-service/cmd/database"
	mylog "random-service/pkg/log"
)

const (
	webPort  = "80"
	rpcPort  = "5002"
	gRPCPort = "50001"
	mongoURL = "mongodb://mongo:27017"
)

type Config struct {
	log  logr.Logger
	db   *database.DB
	name string
}

func NewConfig(db *database.DB, name string, l logr.Logger) Config {
	return Config{
		db:   db,
		name: name,
		log:  l,
	}
}

func main() {
	l := mylog.NewCustomLogger()
	l.Info("connecting to mongodb ...")
	db, err := database.NewDB("random", mongoURL)
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

	app := NewConfig(db, "random-service", l)

	rpcServer := NewRPCServer("randomRPC", *app.db)
	// register to RPC
	err = rpc.Register(rpcServer)
	if err != nil {
		panic("failed to register RPC for RPCServer reveiver")
	}
	app.log.Info("register RPC successed")

	go app.listenRPC(rpcPort)
	app.log.Info("started RPC server", "port", rpcPort)

	go app.RegisterGPRC(gRPCPort, "RandomGRPC")
	app.log.Info("started GRPC server successfully")

	// Listen normal requests
	server := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = server.ListenAndServe()
	if err != nil {
		panic("failed to start random-service")
	}
	app.log.Info("start http server listen on %s \n", webPort)
}

func (app *Config) listenRPC(port string) {
	listenURL := fmt.Sprintf("0.0.0.0:%s", rpcPort)
	listen, err := net.Listen("tcp", listenURL)
	if err != nil {
		app.log.Error(err, "rpc listen failed")
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
