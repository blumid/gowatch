package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/blumid/gowatch/db"
	"github.com/blumid/gowatch/structure"
	"go.mongodb.org/mongo-driver/bson"
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

	for _, v := range temp {
		filter := bson.M{"name": v.Name}
		if db.FindProgram(filter) {
			db.UpdateArray(v.Name, v.Target.InScope)
			// we call UpdateArray() ? function to get new things:

			continue
		} else {

			err := db.AddProgram(&v)
			if err != nil {
				log.Fatal("process - adding to db: ", err)
			}
		}
	}

}
