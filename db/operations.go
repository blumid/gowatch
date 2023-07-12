package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/blumid/gowatch/structure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddProgram(program *structure.Program) error {
	program.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	// program.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	data, _ := bson.Marshal(program)

	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

/*
func UpdateProgram(program *structure.Program) (*mongo.UpdateResult, error) {
	program.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	//if the document not exist this item cause to add.
	options := options.Update().SetUpsert(true)

	filter := bson.M{"name": program.Name}

	update, _ := bson.Marshal(program)
	result, err := collection.UpdateOne(ctx, filter, update, options)

	if err != nil {
		return nil, err
	}
	return result, nil
}
*/

func FindProgram(filter interface{}) bool {

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		fmt.Println("FindProgram-err: ", err)
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

func CheckScope() {

}

func UpdateArray(name string, array []structure.InScope) interface{} {

	// Define the search criteria
	filter := bson.M{"name": name}

	// Define the update
	update := bson.M{"$addToSet": bson.M{"target.inscope": bson.M{"$each": array}}}

	// arrayFilter := bson.D{{"elem.assettype", bson.E{"$eq", "URL"}}}
	// opts := options.Update().SetArrayFilters(options.ArrayFilters{Filters: []interface{}{arrayFilter}})

	// Execute the update
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		panic(err)
	}

	return ""
}

func FandU(name string, array []structure.InScope) {

	/*
		Find program
			- if it already exsits find get list
			- make diff between that and db

			- if not add it to db.
	*/

	// Find
	// filter := bson.M{"name": name}

	// creating stages
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "name", Value: name}}}}
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
		Inscope []structure.InScope `bson:"in,omitempty" json:"in,omitempty"`
		Name    string              `bson:"name,omitempty" json:"name,omitempty"`
	}

	for cursor.Next(context.Background()) {
		var doc result
		if err := cursor.Decode(&doc); err != nil {
			log.Fatal("err2: ", err)
		}
		fmt.Println("doc :", doc.Inscope)
	}

}
