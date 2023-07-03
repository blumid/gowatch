package db

import (
	"fmt"
)

func AddProgram(data []byte) error {
	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

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
