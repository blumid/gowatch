package db

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/blumid/gowatch/structure"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddProgram(program *structure.Program) (primitive.ObjectID, error) {
	program.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	program.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	data, _ := bson.Marshal(program)

	res, err := collection_program.InsertOne(ctx, data)
	if err != nil {
		return primitive.NilObjectID, err
	}
	insertedID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("inserted ID is not a valid ObjectID")
	}

	return insertedID, nil
}

func GetInScopes(name string) (*structure.Result_1, error) {
	filter := bson.D{{Key: "name", Value: name}}
	opts := options.FindOne().SetProjection(bson.D{{Key: "target.inscope", Value: 1}})
	var res structure.Result_1

	err := collection_program.FindOne(context.Background(), filter, opts).Decode(&res)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &res, nil
}

func UpdateInScopes(id primitive.ObjectID, diff []structure.InScope) bool {

	update := bson.M{"$addToSet": bson.M{"target.inscope": bson.M{"$each": diff}}, "$set": bson.M{"updatedat": primitive.NewDateTimeFromTime(time.Now())}}
	filter := bson.M{"_id": id}
	_, err := collection_program.UpdateOne(context.Background(), filter, update)
	if err != nil {
		logrus.Fatal("UpdateInScopes(): ", err)
		return false

	}
	return true
}

func ScopeDiff(new, old []structure.InScope) []structure.InScope {

	m := make(map[structure.InScope]bool)
	for _, item := range old {
		m[item] = true
	}

	var diff []structure.InScope
	var s_type string
	for _, item := range new {
		s_type = strings.ToLower(item.Type)
		if (s_type == "url" || s_type == "wildcard" || s_type == "cidr" || s_type == "api") && !m[item] {
			diff = append(diff, item)
		}
	}
	return diff
}

func AddSub(domain *structure.Subdomain) error {
	data, _ := bson.Marshal(domain)

	_, err := collection_sub.InsertOne(ctx, data)

	if err != nil {
		return err
	}
	return nil
}
