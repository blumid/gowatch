package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection_program, collection_domain *mongo.Collection
var ctx = context.TODO()

func init() {
	//connection string
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection_program = client.Database("gowatch").Collection("programs")
	// collection_domain = client.Database("gowatch").Collection("domains")
}
