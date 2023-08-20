package main

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/go-logr/logr"
	amqp "github.com/rabbitmq/amqp091-go"

	"broker/api"
	mylog "broker/pkg/log"
)

const webPort = "80"

func main() {
	l := mylog.NewCustomLogger()
	// try to connect to rabbitmq
	rabbitConn, err := connect(l)
	if err != nil {
		l.Error(err, "connect rabbitmq failed")
		os.Exit(1)
	}
	defer rabbitConn.Close()

	app := api.NewServer(rabbitConn, l)

	app.Log.Info("Starting broker service on port", webPort)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.Routes(),
	}

	// start the server
	err = srv.ListenAndServe()
	if err != nil {
		app.Log.Error(err, "server listen failed")
		panic(err)
	}
}

func connect(l logr.Logger) (*amqp.Connection, error) {
	var counts int64
	backOff := 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			l.Info("rabbitmq not yet ready...")
			counts++
		} else {
			l.Info("connected to rabbitmq!")
			connection = c
			break
		}

		if counts > 5 {
			l.Error(err, "connect to rabbitmq failed", "retried", counts)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		l.Info("backing off...", "retrying", counts)
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
