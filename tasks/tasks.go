package tasks

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/blumid/gowatch/db"
	"github.com/blumid/gowatch/discord"
	"github.com/blumid/gowatch/structure"
)

func Start() {

	/*  --------- initial ----------- */
	// connect discord bot:
	discord.Connect()

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
	err2 := json.Unmarshal(*file, &temp)
	if err2 != nil {
		fmt.Println("err2 is: ", err2)
	}

	for _, v := range temp {

		res := db.FandU(v.Name, v.Target.InScope)
		if !res {
			if err := db.AddProgram(&v); err != nil {
				fmt.Println("new one add: ", v.Name)
			}

		}

	}
}
