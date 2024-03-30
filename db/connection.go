package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection_program, collection_domain *mongo.Collection
var ctx = context.TODO()
var DBExists = false

// connection string
const url = "mongodb://localhost:27017/"

func init() {

	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err1 := client.Ping(ctx, nil)
	if err1 != nil {
		log.Fatal(err1)
		return
	}

	//check db existence?
	DBExists = existDB(client, "gowatch")

	collection_program = client.Database("gowatch").Collection("programs")

	// collection_domain = client.Database("gowatch").Collection("domains")
}

func existDB(client *mongo.Client, name string) bool {
	databases, _ := client.ListDatabaseNames(context.Background(), bson.M{"name": name})

	for _, db := range databases {
		if db == name {
			return true
		}
	}
	return false
}
