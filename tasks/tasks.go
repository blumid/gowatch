package tasks

import (
	"encoding/json"
	"io"
	"net/http"

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
		task_update_db(file)
	} else {
		res := task_init(file)
		if res {
			logrus.Info("DB created succesfully!")
		}
	}

}
func task_init(file *[]byte) bool {
	var temp []structure.Program
	err := json.Unmarshal(*file, &temp)
	if err != nil {
		logrus.Error("task_init() - unmarshal : ", err)
		return false
	}
	for _, v := range temp {
		err2 := db.AddProgram(&v)
		if err2 != nil {
			logrus.Fatal("task_init() - addProgram : ", err2)
			continue
		}
	}
	return true
}

func download() *[]byte {

	//temp:

	// file, err := os.ReadFile("hackerone_data.json")

	// if err != nil {
	// 	logrus.Error("tasks.download(): ", err)
	// 	return nil
	// }
	url := "https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/main/data/hackerone_data.json"

	res, err := http.Get(url)
	if err != nil {
		logrus.Error("tasks.download() - connecting to url: ", err)
	}

	defer res.Body.Close()
	file, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		logrus.Error("tasks.download() - reading response body: ", err2)
	}
	return &file
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
				db.UpdateInScopes(scopes.ID, diff)
				discord.NotifyNewAsset(&v, diff)
			}
		} else {
			err := db.AddProgram(&v)
			if err != nil {
				logrus.Fatal("task_update_db(): ", err)
				continue
			}
			discord.NotifyNewProgram(&v)
		}

	}
	logrus.Info("task_update: done!")
}
