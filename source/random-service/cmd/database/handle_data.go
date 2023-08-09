package database

import (
	"context"
	"log"
	"random-service/cmd/database/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func (db *DB) Insert(entry models.RandomData) error {
	collection := db.db.Collection("random")

	_, err := collection.InsertOne(context.TODO(), entry)
	if err != nil {
		log.Println("Error inserting into logs:", err)
		return err
	}

	return nil
}

func (db *DB) All() ([]*models.RandomData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := db.db.Collection("random")

	opts := options.Find()
	opts.SetSort(bson.D{{Name: "created_at", Value: -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Finding all docs error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*models.RandomData

	for cursor.Next(ctx) {
		var item models.RandomData

		err := cursor.Decode(&item)
		if err != nil {
			log.Print("Error decoding log into slice:", err)
			return nil, err
		} else {
			logs = append(logs, &item)
		}
	}

	return logs, nil
}

func (db *DB) GetOne(id string) (*models.RandomData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := db.db.Collection("logs")

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var entry models.RandomData
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (db *DB) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := db.db.Collection("logs")

	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil
}

func (db *DB) Update(entry models.RandomData) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := db.db.Collection("logs")

	docID, err := primitive.ObjectIDFromHex(entry.ID)
	if err != nil {
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{Name: "$set", Value: bson.D{
				{Name: "data", Value: entry.Data},
				{Name: "updated_at", Value: time.Now()},
			}},
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}
