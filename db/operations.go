package db

import (
	"context"
	"fmt"
	"time"

	"github.com/blumid/gowatch/structure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	// Execute the update
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		panic(err)
	}
	fmt.Println("modifiedcount is :", result.ModifiedCount)
	// Retrieve the added items
	projection := bson.M{"name": 1, "target.inscope": bson.M{"$slice": -result.ModifiedCount}}
	cursor, err := collection.Find(context.Background(), filter, options.Find().SetProjection(projection))
	if err != nil {
		panic(err)
	}

	// Iterate over the results and print out the added items
	for cursor.Next(context.Background()) {
		var doc bson.M

		err := cursor.Decode(&doc)

		fmt.Println("document is: ", doc)
		if err != nil {
			panic(err)
		}

	}

	return ""
}
