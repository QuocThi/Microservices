package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"random-service/cmd/database"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
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
	log log.Logger
}

func NewLogs() Logs {
	var buf bytes.Buffer
	logger := log.New(&buf, "logger: ", log.Lshortfile)
	return Logs{
		log: *logger,
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
	db, err := database.NewDB("random", mongoURL)
	if err != nil {
		panic(err)
	}

	// close db connection
	defer func() {
		if err = db.Close(); err != nil {
			panic(err)
		}
	}()

	log := NewLogs()
	app := NewConfig(db, "random-service", log)

	// Listen normal requests
	server := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = server.ListenAndServe()
	if err != nil {
		panic("failed to start random-service")
	}
}
