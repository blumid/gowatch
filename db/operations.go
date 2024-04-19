package db

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/blumid/gowatch/structure"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddProgram(program *structure.Program) error {
	program.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	program.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	data, _ := bson.Marshal(program)

	_, err := collection_program.InsertOne(ctx, data)
	if err != nil {
		return err
	}
	return nil
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
	for _, item := range new {
		// for _, el := range []string{"url", "wildcard", "cidr", "api"} {
		// 	if el == strings.ToLower(item.Type) {
		// 		diff = append(diff, item)
		// 		break
		// 	}
		// }
		if (strings.ToLower(item.Type) == "url" || strings.ToLower(item.Type) == "wildcard" || strings.ToLower(item.Type) == "cidr" || strings.ToLower(item.Type) == "api") && !m[item] {
			diff = append(diff, item)
		}
	}
	return diff
}

func AddSub(domain *structure.Domain) error {
	data, _ := bson.Marshal(domain)

	_, err := collection_domain.InsertOne(ctx, data)

	if err != nil {
		return err
	}
	return nil
}
