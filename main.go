package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/blumid/gowatch/structure"
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

	file, err := os.ReadFile("temp.json")

	if err != nil {
		fmt.Println("err in opening file: ", err)
	}

	var temp []structure.Program
	json.Unmarshal(file, &temp)

	for k, v := range temp {
		fmt.Println("program ", k+1, " :", v.Name)
		fmt.Println("Targets-InScope ", k+1, " :", v.Target.InScope)
		fmt.Println("Targets-OutScope", k+1, " :", v.Target.OutScope)
		fmt.Println("Bounty? ", k+1, " :", v.Bounty)
	}
}
