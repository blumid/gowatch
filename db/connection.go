package db

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection_program, collection_asset *mongo.Collection
var ctx = context.TODO()
var DBExists, AssetExist = false, false

// connection string
const url = "mongodb://localhost:27017/"

func init() {

	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logrus.Fatal("db.init(): ", err)
	}

	err1 := client.Ping(ctx, nil)
	if err1 != nil {
		logrus.Fatal("db.init(): ", err1)
		return
	}

	// check existence of db
	DBExists = existDB(client, "gowatch")
	// check  existence of assets
	AssetExist = existColl(client, "assets")
	// fmt.Println("assets existstence :", AssetExist)

	collection_program = client.Database("gowatch").Collection("programs")
	collection_asset = client.Database("gowatch").Collection("assets")
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

func existColl(client *mongo.Client, name string) bool {
	colllections, _ := client.Database("gowatch").ListCollectionNames(context.Background(), bson.M{"name": name})

	for _, db := range colllections {
		if db == name {
			return true
		}
	}
	return false
}
