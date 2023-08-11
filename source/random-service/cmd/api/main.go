package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"random-service/cmd/database"
)

const (
	webPort  = "80"
	rpcPort  = "5002"
	gRPC     = "50001"
	mongoURL = "mongodb://mongo:27017"
)

// TODO: Add logger to this Config, to log at any internal server error for debuging. Ex:
//
//	type Log struct {
//		// enabled/disabled
//		Env       string `json:"env" mapstructure:"env"`
//		Timestamp bool   `json:"timestamp" mapstructure:"timestamp"`
//		// empty mean StdOut
//		FileName string `json:"file_name" mapstructure:"file_name"`
//	}
type Config struct {
	db   *database.DB
	name string
	log  Logs
}

type Logs struct {
	buf *bytes.Buffer
	log *log.Logger
}

func NewLogs() Logs {
	var buf bytes.Buffer
	logger := log.New(&buf, "logger: ", log.Lshortfile)
	return Logs{
		log: logger,
		buf: &buf,
	}
}

func NewConfig(db *database.DB, name string, log Logs) Config {
	return Config{
		db:   db,
		name: name,
		log:  log,
	}
}

func main() {
	log := NewLogs()
	log.PrintLogs("connecting to mongodb ...")
	db, err := database.NewDB("random", mongoURL)
	if err != nil {
		panic(err)
	}
	log.PrintLogs("connected to mongodb")

	// close db connection
	defer func() {
		if err = db.Close(); err != nil {
			panic(err)
		}
	}()

	app := NewConfig(db, "random-service", log)

	rpcServer := NewRPCServer("randomRPC", *app.db)
	// register to RPC
	err = rpc.Register(rpcServer)
	if err != nil {
		panic("failed to register RPC for RPCServer reveiver")
	}
	log.PrintLogs("register RPC successed")

	go app.listenRPC(rpcPort)
	log.log.Printf("start RPC server listen on %s \n", rpcPort)

	// Listen normal requests
	server := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = server.ListenAndServe()
	if err != nil {
		panic("failed to start random-service")
	}
	log.log.Printf("start http server listen on %s \n", webPort)
}

func (app *Config) listenRPC(port string) {
	listenURL := fmt.Sprintf("0.0.0.0:%s", rpcPort)
	listen, err := net.Listen("tcp", listenURL)
	if err != nil {
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
