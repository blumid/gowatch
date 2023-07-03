package db

import "github.com/blumid/gowatch/structure"

func CreateProgram(programs *structure.Program) error {
	_, err := collection.InsertOne(ctx, programs)
	return err
}
