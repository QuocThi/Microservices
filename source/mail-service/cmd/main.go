package main

import (
	"fmt"
	"net/http"

	"mailer-service/api"
	"mailer-service/internal/mailer"
	mylog "mailer-service/pkg/log"
)

const webPort = "80"

func main() {
	l := mylog.NewCustomLogger()
	app := api.NewConfig(mailer.NewMailer(), l)

	app.Log.Info("starting mail service... ", "port", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.Routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		app.Log.Error(err, "server listen failed")
		panic(err)
	}
}
