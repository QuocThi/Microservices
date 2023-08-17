package main

import (
	"math"
	"os"
	"time"

	"github.com/go-logr/logr"
	amqp "github.com/rabbitmq/amqp091-go"

	"listener/event"
	mylog "listener/pkg/log"
)

func main() {
	l := mylog.NewCustomLogger()
	// try to connect to rabbitmq
	rabbitConn, err := connect(l)
	if err != nil {
		l.Error(err, "connect to rabbitmq failed")
		os.Exit(1)
	}
	defer rabbitConn.Close()

	// start listening for messages
	l.Info("listening for and consuming RabbitMQ messages...")

	// create consumer
	consumer, err := event.NewConsumer(rabbitConn, "listener", l)
	if err != nil {
		panic(err)
	}

	// watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		consumer.GetLog().Error(err, "comsumer listen failed")
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
			l.Info("rabbitmw not yet ready...")
			counts++
		} else {
			l.Info("connected to rabbitmq!")
			connection = c
			break
		}

		if counts > 5 {
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		l.Info("backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
