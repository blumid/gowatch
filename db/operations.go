package db

import (
	"context"
	"log"
	"time"

	"github.com/blumid/gowatch/structure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddProgram(program *structure.Program) error {
	program.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	program.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	data, _ := bson.Marshal(program)

	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func FandU(name string, array []structure.InScope) bool {

	/*
		Find program
			- if it already exsits find get list
			- make diff between that and db
			- if not add it to db.
	*/

	// creating stages
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "name", Value: name}}}}

	// just comment the projection temparially:
	projection := bson.D{{Key: "$project", Value: bson.M{
		"in": bson.M{
			"$filter": bson.M{
				"input": "$target.inscope",
				"as":    "elem",
				"cond":  bson.M{"$in": []interface{}{"$$elem.assettype", bson.A{"URL", "CIDR"}}},
			},
		},
		"name": 1,
		"_id":  0},
	},
	}

	pipeline := mongo.Pipeline{matchStage, projection}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		log.Fatal("err1: ", err)
	}

	type result struct {
		Inscopes []structure.InScope `bson:"in,omitempty" json:"in,omitempty"`
		Name     string              `bson:"name,omitempty" json:"name,omitempty"`
	}
	defer cursor.Close(context.TODO())

	var has bool = false

	for cursor.Next(context.Background()) {

		has = true
		var res result
		if err := cursor.Decode(&res); err != nil {
			log.Fatal("err2: ", err)
		}
		diff := scopeDifference(array, res.Inscopes)

		// update using this diff thing:
		update := bson.M{"$addToSet": bson.M{"target.inscope": bson.M{"$each": diff}}, "$set": bson.M{"updatedat": primitive.NewDateTimeFromTime(time.Now())}}
		filter := bson.M{"name": name}
		collection.UpdateOne(context.Background(), filter, update)
		// fmt.Println("diff is: ", diff)

		// fmt.Println("doc :", res.Inscopes)
	}
	return has
}

func scopeDifference(a, b []structure.InScope) []structure.InScope {

	m := make(map[structure.InScope]bool)
	for _, item := range b {
		m[item] = true
	}

	var diff []structure.InScope
	for _, item := range a {
		if (item.AssetType == "CIDR" || item.AssetType == "URL") && !m[item] {
			diff = append(diff, item)
		}
	}
	return diff
}
