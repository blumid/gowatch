package tasks

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	// "github.com/blumid/gowatch/discord"
	"github.com/blumid/gowatch/db"
	"github.com/blumid/gowatch/structure"
)

func Start() {

	/*  --------- initial ----------- */
	// connect discord bot:
	// discord.Connect()

	//download json file
	file := download()
	parseJson(&file)

}

func download() []byte {

	//temp:

	file, err := os.ReadFile("temp.json")

	if err != nil {
		fmt.Println("err in opening file: ", err)
		return nil
	}

	// res, err := http.Get("https://github.com/arkadiyt/bounty-targets-data/blob/main/data/hackerone_data.json")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer res.Body.Close()
	// file, err2 := io.ReadAll(res.Body)
	// if err != nil {
	// 	log.Fatal(err2)
	// }
	return file
}

func parseJson(file *[]byte) {
	var temp []structure.Program
	err := json.Unmarshal(*file, &temp)
	if err != nil {
		fmt.Println("json Unmarshal- err is: ", err)
		return
	}

	for _, v := range temp {

		// fmt.Println("\n", k, v)

		// db.GetInScopes(v.Name, v.Target.InScope)
		scopes, err := db.GetInScopes(v.Name)
		if err != nil {
			log.Fatal(err)
			continue
		}
		if scopes != nil {
			diff := db.ScopeDiff(v.Target.InScope, scopes.Target.InScope)
			if diff != nil {
				fmt.Println("diff is: ", v.Name, ": ", diff)
				//update db
			} else {
				fmt.Println(v.Name, ": ", "no diff")
				//return nothing
			}
		} else {
			fmt.Println("need to add program: ", v.Name)
		}

		// res := db.FandU(v.Name, v.Target.InScope)
		// if !res {
		// 	if err := db.AddProgram(&v); err != nil {
		// 		fmt.Println("new one add: ", v.Name)
		// 	}

		// } else {
		// 	continue
		// }

	}
}
