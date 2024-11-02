package config

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Ctx = context.TODO()

func ConnectDB() (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(Ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	return client.Database("medicalDB"), nil
}
