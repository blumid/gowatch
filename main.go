package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/blumid/gowatch/db"
	"github.com/blumid/gowatch/structure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {

	// I just commented this part to test with local file:

	// res, err := http.Get("https://github.com/arkadiyt/bounty-targets-data/blob/main/data/hackerone_data.json")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer res.Body.Close()
	// _, err2 := io.ReadAll(res.Body)
	// if err != nil {
	// 	log.Fatal(err2)
	// }
	process()

}

func process() {
	file, err := os.ReadFile("temp.json")

	if err != nil {
		fmt.Println("err in opening file: ", err)
	}

	var temp []structure.Program
	err2 := json.Unmarshal(file, &temp)
	if err2 != nil {
		fmt.Println("err2 is: ", err2)
	}

	for k, v := range temp {
		fmt.Println("program ", k+1, " :", v.Name)
		filter := bson.M{"name": v.Name}
		if db.FindProgram(filter) {
			continue
		} else {
			v.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
			v.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
			data, _ := bson.Marshal(v)
			err := db.AddProgram(data)
			if err != nil {
				log.Fatal("process - adding to db: ", err)
			}
		}
	}

}
