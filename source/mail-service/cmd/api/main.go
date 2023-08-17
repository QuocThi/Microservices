package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-logr/logr"

	mylog "mailer-service/pkg/log"
)

type Config struct {
	Mailer Mail
	Log    logr.Logger
}

func NewConfig(mail Mail, l logr.Logger) Config {
	return Config{
		Mailer: mail,
		Log:    l,
	}
}

const webPort = "80"

func main() {
	l := mylog.NewCustomLogger()
	app := NewConfig(createMail(), l)

	app.Log.Info("starting mail service... ", "port", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		app.Log.Error(err, "server listen failed")
		panic(err)
	}
}

func createMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	m := Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromName:    os.Getenv("FROM_NAME"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
	}

	return m
}
