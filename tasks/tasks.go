package tasks

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

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
	files := download()

	for k, v := range files {
		var file = readFile(v)
		if db.DBExists {
			task_update_db(file, k)
		} else {
			res := task_init(file, k)
			if res {
				logrus.Info("DB created succesfully!")
			}
		}
	}
}

func task_init(file *[]byte, owner string) bool {
	var temp []structure.Program
	err := json.Unmarshal(*file, &temp)
	if err != nil {
		logrus.Error("task_init() - unmarshal : ", err)
		return false
	}
	for _, v := range temp {
		v.Owner = owner
		err2 := db.AddProgram(&v)
		if err2 != nil {
			logrus.Fatal("task_init() - addProgram : ", err2)
			continue
		}
	}
	return true
}

func task_update_db(file *[]byte, owner string) {
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
				v.Owner = owner
				logrus.Info(v.Name, ", diff is: ", diff)
				discord.NotifyNewAsset(&v, diff)

				db.UpdateInScopes(scopes.ID, diff)
			}
		} else {
			v.Owner = owner
			logrus.Info(v.Name, ", is a new program!")
			discord.NotifyNewProgram(&v)
			err := db.AddProgram(&v)
			if err != nil {
				logrus.Fatal("task_update_db(): ", err)
				continue
			}

		}

	}
	logrus.Info(owner + " updated")
}

func download() map[string]string {

	commands := map[int]string{
		// hackerOne
		0: "wget -O HackerOne.json -A json https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/master/data/hackerone_data.json",
		1: "jq 'map(.targets |= with_entries(if .key == \"in_scope\" or .key == \"out_of_scope\" then .value |= map(if .asset_identifier then .asset = .asset_identifier | del(.asset_identifier) else . end | if .asset_type then .type = .asset_type | del(.asset_type) else . end) else . end))' HackerOne.json > temp.json && mv temp.json HackerOne.json",

		// intigriti
		2: "wget -O Intigriti.json -A json https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/main/data/intigriti_data.json",
		3: "jq 'map(.targets |= with_entries(if .key == \"in_scope\" or .key == \"out_of_scope\" then .value |= map(if .endpoint then .asset = .endpoint | del(.endpoint) else . end) else . end))' Intigriti.json > temp.json && mv temp.json Intigriti.json",

		// bugCrowd
		4: "wget -O BugCrowd.json -A json https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/main/data/bugcrowd_data.json",
		5: "jq 'map(.targets |= with_entries(if .key == \"in_scope\" or .key == \"out_of_scope\" then .value |= map(if .target then .asset = .target | del(.target) else . end | if .type == \"website\" then .type = \"url\" else . end) else . end))' BugCrowd.json > temp.json && mv temp.json BugCrowd.json",
	}

	for i := 0; i < len(commands); i++ {

		cmd := fmt.Sprintf(commands[i])
		runCommand(cmd)
		time.Sleep(time.Millisecond * 150)
	}
	return map[string]string{"hackerone": "HackerOne.json", "intigriti": "Intigriti.json", "bugcrowd": "BugCrowd.json"}
}

func runCommand(command string) {
	com := exec.Command("bash", "-c", command)
	if err := com.Run(); err != nil {
		fmt.Println("error is:", err)
		logrus.Error("tasks.runCommand():", err)
		os.Exit(1)
	}
}

func readFile(name string) *[]byte {
	file, err := os.ReadFile(name)

	if err != nil {
		logrus.Error("tasks.download(): ", err)
		return nil
	}
	return &file
}
