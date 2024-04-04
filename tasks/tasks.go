package tasks

import (
	"encoding/json"
	"os"

	"github.com/blumid/gowatch/db"
	"github.com/blumid/gowatch/discord"
	"github.com/blumid/gowatch/structure"
	logrus "github.com/sirupsen/logrus"
)

func Start() {

	/*  --------- initial ----------- */
	// connect discord bot:
	discord.Open()

	//download json file
	file := download()
	if db.DBExists {
		task_update_db(&file)
	} else {
		task_init(&file)
	}

}
func task_init(file *[]byte) bool {
	var temp []structure.Program
	err := json.Unmarshal(*file, &temp)
	if err != nil {
		logrus.Error("task_init(): ", err)
		return false
	}
	for _, v := range temp {
		err2 := db.AddProgram(&v)
		if err2 != nil {
			logrus.Fatal("task_init(): ", err2)
		}
	}
	return true
}

func download() []byte {

	//temp:

	file, err := os.ReadFile("temp.json")

	if err != nil {
		logrus.Error("tasks.download(): ", err)
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

func task_update_db(file *[]byte) {
	var temp []structure.Program
	err := json.Unmarshal(*file, &temp)
	if err != nil {
		logrus.Error("task_update_db(): ", err)
		return
	}

	for _, v := range temp {

		scopes, err := db.GetInScopes(v.Name)
		if err != nil {
			logrus.Fatal("task_update_db(): ", err)
			continue
		}
		if scopes != nil {
			diff := db.ScopeDiff(v.Target.InScope, scopes.Target.InScope)
			if diff != nil {
				logrus.Info(v.Name, ", diff is: ", diff)
				// db.UpdateInScopes(scopes.ID, diff)
				// discord.NotifyNewAsset(&v, diff)
			} else {
				logrus.Info(v.Name, ", no diff")
				continue
			}
		} else {
			// err := db.AddProgram(&v)
			// if err != nil {
			// 	logrus.Fatal("task_update_db(): ", err)
			// 	continue
			// }
			// discord.NotifyNewProgram(&v)
		}

	}
}
