package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-logr/logr"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"

	"authentication/data"
	mylog "authentication/pkg/log"
)

const webPort = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
	Log    logr.Logger
}

func main() {
	logger := mylog.NewCustomLogger()

	logger.Info("Connecting to Postgres...")
	conn := connectToDB(logger)
	if conn == nil {
		panic("Can't connect to Postgres")
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
		Log:    logger,
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	logger.Info("Starting authentication service...")
	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB(logger logr.Logger) *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			logger.Info("Postgres not yet ready ...")
			counts++
		} else {
			logger.Info("Connected to Postgres!")
			return connection
		}

		if counts > 10 {
			logger.Error(err, "connectToDB failed!", "total retries: ", counts)
			return nil
		}

		logger.Info("Backing off for two seconds...", "Retry: ", counts)
		time.Sleep(2 * time.Second)
		continue
	}
}
