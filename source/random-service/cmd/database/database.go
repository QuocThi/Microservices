package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	db *mongo.Database
}

func NewDB(name, url string) (*DB, error) {
	opts := options.Client().
		ApplyURI(url).SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(name)
	return &DB{
		db: db,
	}, err
}

func (db *DB) Close() error {
	return db.db.Client().Disconnect(context.Background())
}
