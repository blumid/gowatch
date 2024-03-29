package db

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/blumid/gowatch/structure"
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

func UpdateInScope(id primitive.ObjectID, diff []structure.InScope) bool {

	update := bson.M{"$addToSet": bson.M{"target.inscope": bson.M{"$each": diff}}, "$set": bson.M{"updatedat": primitive.NewDateTimeFromTime(time.Now())}}
	filter := bson.M{"_id": id}
	_, err := collection_program.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
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
		if (item.AssetType == "CIDR" || item.AssetType == "URL") && !m[item] {
			diff = append(diff, item)
		}
	}
	return diff
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

	cursor, err := collection_program.Aggregate(context.Background(), pipeline)
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
		diff := ScopeDiff(array, res.Inscopes)

		// update using this diff thing:
		update := bson.M{"$addToSet": bson.M{"target.inscope": bson.M{"$each": diff}}, "$set": bson.M{"updatedat": primitive.NewDateTimeFromTime(time.Now())}}
		filter := bson.M{"name": name}
		collection_program.UpdateOne(context.Background(), filter, update)
		// fmt.Println("diff is: ", diff)

		// fmt.Println("doc :", res.Inscopes)
	}
	return has
}

func AddSub(domain *structure.Domain) error {
	data, _ := bson.Marshal(domain)

	_, err := collection_domain.InsertOne(ctx, data)

	if err != nil {
		return err
	}
	return nil
}
