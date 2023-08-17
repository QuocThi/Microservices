package event

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/go-logr/logr"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
	log       logr.Logger
}

func NewConsumer(conn *amqp.Connection, name string, log logr.Logger) (Consumer, error) {
	consumer := Consumer{
		conn:      conn,
		queueName: name,
		log:       log,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) GetLog() logr.Logger {
	return consumer.log
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		ch.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload, consumer.log)
		}
	}()

	consumer.log.Info("Waiting for message [Exchange, Queue]", "queueInfo", q)
	<-forever

	return nil
}

func handlePayload(payload Payload, l logr.Logger) {
	switch payload.Name {
	case "log", "event":
		// log whatever we get
		err := logEvent(payload)
		if err != nil {
			l.Error(err, "handle payload failed")
		}

	case "auth":
		// authenticate

	// you can have as many cases as you want, as long as you write the logic

	default:
		err := logEvent(payload)
		if err != nil {
			l.Error(err, "handle default payload failed")
		}
	}
}

func logEvent(entry Payload) error {
	jsonData, err := json.MarshalIndent(entry, "", "\t")
	if err != nil {
		return err
	}

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}
