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
		fmt.Println("there is! ")
		return true
	}
	return false
}

func CheckScope() {

}

func Updated() interface{} {
	// Define the search criteria
	filter := bson.M{"_id": "document_id"}

	// Define the update
	update := bson.M{"$addToSet": bson.M{"myArrayField": bson.M{"$each": []string{"new_item1", "new_item2"}}}}

	// Execute the update
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		panic(err)
	}
	// Retrieve the added items
	projection := bson.M{"myArrayField": bson.M{"$slice": bson.A{-result.ModifiedCount, result.ModifiedCount}}}
	cursor, err := collection.Find(context.Background(), filter, options.Find().SetProjection(projection))
	if err != nil {
		panic(err)
	}

	// Iterate over the results and print out the added items
	for cursor.Next(context.Background()) {
		var document bson.M
		err := cursor.Decode(&document)
		if err != nil {
			panic(err)
		}
		addedItems := document["myArrayField"].([]interface{})[:result.ModifiedCount]
		fmt.Println("Added items:", addedItems)
	}

	return ""
}
