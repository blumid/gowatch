package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/blumid/gowatch/db"
	"github.com/blumid/gowatch/discord"
	"github.com/blumid/gowatch/structure"
	logrus "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/projectdiscovery/goflags"
	httpx "github.com/projectdiscovery/httpx/runner"
	"github.com/projectdiscovery/subfinder/v2/pkg/resolve"
	subfinder "github.com/projectdiscovery/subfinder/v2/pkg/runner"
)

// #phase 1:

func Start() {

	/*  --------- initial ----------- */
	// connect discord bot:
	discord.Open()

	//download json file
	dls := getDls()

	for k, v := range dls {
		var file = readFile(v)
		if db.DBExists {
			task_update_db(file, k)
			break
			//
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
		id, err2 := db.AddProgram(&v)
		if err2 != nil {
			logrus.Fatal("task_init() - addProgram : ", err2)
			continue
		}
		doScopes(id, v.Target.InScope)
	}
	return true
}

func task_update_db(file *[]byte, owner string) {
	var temp []structure.Program
	err := json.Unmarshal(*file, &temp)
	if err != nil {
		// logrus.Error("task_update_db(): ", err)
		fmt.Println("task_update_db(): ", err)
		return
	}

	for _, v := range temp {

		scopes, err := db.GetInScopes(v.Name)
		if err != nil {
			// logrus.Fatal("task_update_db(): ", err)
			fmt.Println("task_update_db(): ", err)
			continue
		}
		if scopes != nil {
			diff := db.ScopeDiff(v.Target.InScope, scopes.Target.InScope)
			if diff != nil {
				v.Owner = owner

				// test
				// logrus.Info(v.Name, ", diff is: ", diff)
				// discord.NotifyNewAsset(&v, diff)
				//

				db.UpdateInScopes(scopes.ID, diff)

				doScopes(scopes.ID, diff)
			}
		} else {
			v.Owner = owner
			logrus.Info(v.Name, ", is a new program!")
			discord.NotifyNewProgram(&v)
			id, err := db.AddProgram(&v)

			if err != nil {
				logrus.Fatal("task_update_db(): ", err)
				continue
			}
			doScopes(id, v.Target.InScope)

		}

	}
	// test
	// logrus.Info(owner + " updated")
	//
}

func getDls() map[string]string {

	commands := map[int]string{
		// hackerOne
		0: "wget -O HackerOne.json -A json https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/master/data/hackerone_data.json",
		1: "jq 'map(.bounty = (.offers_bounties | tostring) | del(.offers_bounties) | .targets |= with_entries(if .key == \"in_scope\" or .key == \"out_of_scope\" then .value |= map(if .asset_identifier then .asset = .asset_identifier | del(.asset_identifier) else . end | if .asset_type then .type = .asset_type | del(.asset_type) else . end) else . end))' HackerOne.json > temp.json && mv temp.json HackerOne.json",

		// intigriti
		2: "wget -O Intigriti.json -A json https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/main/data/intigriti_data.json",
		3: "jq 'map(.bounty = (.min_bounty.value | tostring) + \"-\" + (.max_bounty.value | tostring) + \" \" + .min_bounty.currency | .targets |= with_entries(if .key == \"in_scope\" or .key == \"out_of_scope\" then .value |= map(if .endpoint then .asset = .endpoint | del(.endpoint) else . end) else . end))' Intigriti.json > temp.json && mv temp.json Intigriti.json",

		// // bugCrowd
		4: "wget -O BugCrowd.json -A json https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/main/data/bugcrowd_data.json",
		5: "jq 'map(.bounty = \"max: \" + (.max_payout | tostring) | del(.max_payout) | .targets |= with_entries(if .key == \"in_scope\" or .key == \"out_of_scope\" then .value |= map(if .target then .asset = .target | del(.target) else . end | if .type == \"website\" then .type = \"url\" else . end) else . end))' BugCrowd.json > temp.json && mv temp.json BugCrowd.json",

		// wildcards
		// 6: "wget https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/main/data/wildcards.txt",
	}

	for i := 0; i < len(commands); i++ {

		cmd := fmt.Sprintf(commands[i])
		runCommand(cmd)
		time.Sleep(time.Millisecond * 500)
	}

	return map[string]string{"hackerone": "HackerOne.json", "intigriti": "Intigriti.json", "bugcrowd": "BugCrowd.json"}
}

func runCommand(command string) {
	com := exec.Command("bash", "-c", command)
	if err := com.Run(); err != nil {
		fmt.Println("tasks.runCommand():", err)
		// logrus.Error("tasks.runCommand():", err)
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

// #phase 2:

func enumerateSubs(domain string) string {

	subfinderOpts := &subfinder.Options{
		Threads:            10,
		Silent:             true,
		Timeout:            30,
		MaxEnumerationTime: 10,
		ResultCallback: func(s *resolve.HostEntry) {
			// I commented this for get all output at same time and save in a file.
			// enumerateTech(prog_id, s.Host)
		},

		ProviderConfig: "~/.config/subfinder/provider-config.yaml",
	}

	subfinder, err := subfinder.NewRunner(subfinderOpts)
	if err != nil {
		// log.Fatalf("failed to create subfinder runner: %v", err)
		// logrus.Error("EnumerateSub-newRunner:",err)
		fmt.Println(err)
	}

	// I commented this:
	// output := &bytes.Buffer{}

	// Open a file for writing
	fileName := fmt.Sprintf("subdomains_%04x.txt", rand.Intn(0x10000))
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return ""
	}
	defer file.Close()

	// To run subdomain enumeration on a single domain
	if err = subfinder.EnumerateSingleDomainWithCtx(context.Background(), domain, []io.Writer{file}); err != nil {
		// log.Fatalf("failed to enumerate single domain: %v", err)
		fmt.Println(err)
	}

	return fileName
}

func enumerateTechSingle(prog_id primitive.ObjectID, domain string) {

	temp := structure.Subdomain{ProgramID: prog_id, Sub: domain}

	options := httpx.Options{
		Methods:                 "GET",
		Silent:                  true,
		Threads:                 1,
		InputTargetHost:         goflags.StringSlice{domain},
		ResponseHeadersInStdout: true,
		OnResult: func(r httpx.Result) {
			// handle error
			if r.Err != nil {
				// fmt.Printf("[Err] %s: %s\n", r.Input, r.Err)
				return
			}
			// fill a temp var && calling AddSub()
			temp.SC = r.StatusCode

			if r.Location != "" {
				temp.Locatoin = r.Location
			}

			temp.Icon = r.FavIconMMH3
			temp.CDN = r.CDN

			temp.Detail.A = r.A
			temp.Detail.Cname = r.CNAMEs
			temp.Detail.Tech = r.Technologies
			temp.Detail.Headers = r.ResponseHeaders

		},
	}
	if err := options.ValidateOptions(); err != nil {
		// logrus.Fatal(err)
		fmt.Println("tasks.EnumerateTech(): ", err)
	}
	httpxRunner, err := httpx.New(&options)
	if err != nil {
		// logrus.Fatal(err)
		fmt.Println("tasks.EnumerateTech(): ", err)
	}
	defer httpxRunner.Close()
	httpxRunner.RunEnumeration()

	if temp.SC != 0 {
		db.AddSub(&temp)
	}
}

func enumerateTechMulti(prog_id primitive.ObjectID, file_name string) {

	temp := structure.Subdomain{ProgramID: prog_id}

	options := httpx.Options{
		Methods:                 "GET",
		Silent:                  true,
		Threads:                 10,
		InputFile:               file_name,
		ResponseHeadersInStdout: true,
		OnResult: func(r httpx.Result) {
			// handle error
			if r.Err != nil {
				// fmt.Printf("[Err] %s: %s\n", r.Input, r.Err)
				return
			}
			// fill a temp var && calling AddSub()
			temp.Sub = r.Input
			temp.SC = r.StatusCode

			if r.Location != "" {
				temp.Locatoin = r.Location
			}

			temp.Icon = r.FavIconMMH3
			temp.CDN = r.CDN

			temp.Detail.A = r.A
			temp.Detail.Cname = r.CNAMEs
			temp.Detail.Tech = r.Technologies
			temp.Detail.Headers = r.ResponseHeaders

		},
	}
	if err := options.ValidateOptions(); err != nil {
		// logrus.Fatal(err)
		fmt.Println("tasks.EnumerateTech(): ", err)
	}
	httpxRunner, err := httpx.New(&options)
	if err != nil {
		// logrus.Fatal(err)
		fmt.Println("tasks.EnumerateTech(): ", err)
	}
	defer httpxRunner.Close()
	httpxRunner.RunEnumeration()

	if temp.SC != 0 {
		db.AddSub(&temp)
	}
}

func isWild(domain string) bool {
	wildcardPattern := `^\*\.[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(wildcardPattern, domain)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	return matched
}

func doScopes(id primitive.ObjectID, assets []structure.InScope) {

	var s_type, d string

	for _, v := range assets {
		s_type = strings.ToLower(v.Type)
		if s_type == "url" || s_type == "wildcard" || s_type == "api" {

			if isWild(v.Asset) {

				//get rid of Asterisk
				d = strings.TrimLeft(v.Asset, "*.")

				file_name := enumerateSubs(d)
				enumerateTechMulti(id, file_name)

				// delete file
				if err := os.Remove(file_name); err != nil {
					fmt.Printf("Error deleting file: %v\n", err)
					continue
				}
			} else {
				enumerateTechSingle(id, v.Asset)
			}
		}
		continue
	}

}
