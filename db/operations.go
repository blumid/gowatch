package db

import (
	"fmt"
	"time"

	"github.com/blumid/gowatch/structure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
