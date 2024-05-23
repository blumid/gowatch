package tasks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/blumid/gowatch/db"
	"github.com/blumid/gowatch/discord"
	"github.com/blumid/gowatch/structure"
	logrus "github.com/sirupsen/logrus"

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

	// test

	for k, v := range dls {
		var file = readFile(v)
		if db.DBExists {
			task_update_db(file, k)
			// test
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
		err2 := db.AddProgram(&v)
		// EnumerateSubs(v.in_scopes)
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

				DoScopes(diff)
			}
		} else {
			v.Owner = owner
			// logrus.Info(v.Name, ", is a new program!")
			// discord.NotifyNewProgram(&v)
			// err := db.AddProgram(&v)
			// // EnumerateSubs(v.in_scopes)
			// if err != nil {
			// 	logrus.Fatal("task_update_db(): ", err)
			// 	continue
			// }
			DoScopes(v.Target.InScope)

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
		// 2: "wget -O Intigriti.json -A json https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/main/data/intigriti_data.json",
		// 3: "jq 'map(.bounty = (.min_bounty.value | tostring) + \"-\" + (.max_bounty.value | tostring) + \" \" + .min_bounty.currency | .targets |= with_entries(if .key == \"in_scope\" or .key == \"out_of_scope\" then .value |= map(if .endpoint then .asset = .endpoint | del(.endpoint) else . end) else . end))' Intigriti.json > temp.json && mv temp.json Intigriti.json",

		// // bugCrowd
		// 4: "wget -O BugCrowd.json -A json https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/main/data/bugcrowd_data.json",
		// 5: "jq 'map(.bounty = \"max: \" + (.max_payout | tostring) | del(.max_payout) | .targets |= with_entries(if .key == \"in_scope\" or .key == \"out_of_scope\" then .value |= map(if .target then .asset = .target | del(.target) else . end | if .type == \"website\" then .type = \"url\" else . end) else . end))' BugCrowd.json > temp.json && mv temp.json BugCrowd.json",

		// wildcards
		// 6: "wget https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/main/data/wildcards.txt",
	}

	for i := 0; i < len(commands); i++ {

		cmd := fmt.Sprintf(commands[i])
		runCommand(cmd)
		time.Sleep(time.Millisecond * 500)
	}

	// test
	return map[string]string{"hackerone": "HackerOne.json"}
	//
	// return map[string]string{"hackerone": "HackerOne.json", "intigriti": "Intigriti.json", "bugcrowd": "BugCrowd.json"}
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

func enumerateSubs(domain string) {
	// var subs []string
	subfinderOpts := &subfinder.Options{
		Threads:            10,
		Timeout:            30,
		MaxEnumerationTime: 10,
		ResultCallback: func(s *resolve.HostEntry) {
			enumerateTech(s.Host)
			// fmt.Println(s)
		},

		ProviderConfig: "~/.config/subfinder/provider-config.yaml",
	}

	subfinder, err := subfinder.NewRunner(subfinderOpts)
	if err != nil {
		// log.Fatalf("failed to create subfinder runner: %v", err)
		// logrus.Error("EnumerateSub-newRunner:",err)
		fmt.Println(err)
	}

	output := &bytes.Buffer{}
	// To run subdomain enumeration on a single domain
	if err = subfinder.EnumerateSingleDomainWithCtx(context.Background(), domain, []io.Writer{output}); err != nil {
		// log.Fatalf("failed to enumerate single domain: %v", err)
		fmt.Println(err)
	}

	// result := output.String()
	// subs = strings.Split(result, "\n")

	// return subs
}

func enumerateTech(domain string) {

	temp := structure.Subdomain{Sub: domain}

	options := httpx.Options{
		Methods:                 "GET",
		Silent:                  true,
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
			temp.CL = r.ContentLength

			if r.Location != "" {
				temp.Locatoin = r.Location
			}

			temp.Detail.Tech = r.Technologies
			temp.Detail.Icon = r.FavIconMMH3
			temp.Detail.Headers = r.ResponseHeaders
			temp.Detail.A = r.A
			temp.Detail.Cname = r.CNAMEs
			temp.Detail.CDN = r.CDN

			fmt.Println("raw headers: ", r)

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
		fmt.Println("let's add to db ", temp.Sub)
		// db.AddSub(&temp)
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

func DoScopes(assets []structure.InScope) {

	var s_type, d string

	for _, v := range assets {
		s_type = strings.ToLower(v.Type)
		if s_type == "url" || s_type == "wildcard" || s_type == "api" {

			if isWild(v.Asset) {

				//get rid of Asterisk
				d = strings.TrimLeft(v.Asset, "*.")
				enumerateSubs(d)
			} else {
				enumerateTech(v.Asset)
			}
		}
		continue
	}

	// fmt.Println("subs are: ", subs)
	// add to db here:

}
