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

	httpx "github.com/projectdiscovery/httpx/runner"
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
				// EnumerateSubs(diff)
				// DoScopes(diff)
			}
		} else {
			v.Owner = owner
			logrus.Info(v.Name, ", is a new program!")
			discord.NotifyNewProgram(&v)
			err := db.AddProgram(&v)
			// EnumerateSubs(v.in_scopes)
			if err != nil {
				logrus.Fatal("task_update_db(): ", err)
				continue
			}

		}

	}
	logrus.Info(owner + " updated")
}

func getDls() map[string]string {

	commands := map[int]string{
		// hackerOne
		0: "wget -O HackerOne.json -A json https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/master/data/hackerone_data.json",
		1: "jq 'map(.bounty = (.offers_bounties | tostring) | del(.offers_bounties) | .targets |= with_entries(if .key == \"in_scope\" or .key == \"out_of_scope\" then .value |= map(if .asset_identifier then .asset = .asset_identifier | del(.asset_identifier) else . end | if .asset_type then .type = .asset_type | del(.asset_type) else . end) else . end))' HackerOne.json > temp.json && mv temp.json HackerOne.json",

		// intigriti
		2: "wget -O Intigriti.json -A json https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/main/data/intigriti_data.json",
		3: "jq 'map(.bounty = (.min_bounty.value | tostring) + \"-\" + (.max_bounty.value | tostring) + \" \" + .min_bounty.currency | .targets |= with_entries(if .key == \"in_scope\" or .key == \"out_of_scope\" then .value |= map(if .endpoint then .asset = .endpoint | del(.endpoint) else . end) else . end))' Intigriti.json > temp.json && mv temp.json Intigriti.json",

		// bugCrowd
		4: "wget -O BugCrowd.json -A json https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/main/data/bugcrowd_data.json",
		5: "jq 'map(.bounty = \"max: \" + (.max_payout | tostring) | del(.max_payout) | .targets |= with_entries(if .key == \"in_scope\" or .key == \"out_of_scope\" then .value |= map(if .target then .asset = .target | del(.target) else . end | if .type == \"website\" then .type = \"url\" else . end) else . end))' BugCrowd.json > temp.json && mv temp.json BugCrowd.json",

		// wildcards
		6: "wget https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/main/data/wildcards.txt",
		// 7: "",
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

// #phase 2:

func EnumerateSubs(domains []string) {
	subfinderOpts := &subfinder.Options{
		Threads:            10, // Thread controls the number of threads to use for active enumerations
		Timeout:            30, // Timeout is the seconds to wait for sources to respond
		MaxEnumerationTime: 10, // MaxEnumerationTime is the maximum amount of time in mins to wait for enumeration
		// ResultCallback: func(s *resolve.HostEntry) {
		// callback function executed after each unique subdomain is found
		// },
		ProviderConfig: "~/.config/subfinder/provider-config.yaml",
		// and other config related options
	}

	subfinder, err := subfinder.NewRunner(subfinderOpts)
	if err != nil {
		// log.Fatalf("failed to create subfinder runner: %v", err)
		// logrus.Error("EnumerateSub-newRunner:",err)
		fmt.Println(err)
	}

	output := &bytes.Buffer{}
	// To run subdomain enumeration on a single domain
	if err = subfinder.EnumerateSingleDomainWithCtx(context.Background(), "hackerone.com", []io.Writer{output}); err != nil {
		// log.Fatalf("failed to enumerate single domain: %v", err)
		fmt.Println(err)
	}

	// To run subdomain enumeration on a list of domains from file/reader
	file, err := os.Open("domains.txt")
	if err != nil {
		// log.Fatalf("failed to open domains file: %v", err)
		fmt.Println(err)
	}
	defer file.Close()

	if err = subfinder.EnumerateMultipleDomainsWithCtx(context.Background(), file, []io.Writer{output}); err != nil {
		// log.Fatalf("failed to enumerate subdomains from file: %v", err)
		fmt.Println(err)
	}
}

func EnumerateTech(domain string) {
	options := httpx.Options{
		Methods: "GET",
		// InputTargetHost: goflags.StringSlice{"scanme.sh", "projectdiscovery.io", "localhost"},
		//InputFile: "./targetDomains.txt", // path to file containing the target domains list
		OnResult: func(r httpx.Result) {
			// handle error
			if r.Err != nil {
				fmt.Printf("[Err] %s: %s\n", r.Input, r.Err)
				return
			}
			fmt.Printf("%s %s %d\n", r.Input, r.Host, r.StatusCode)
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
}

func isWild(domain string) bool {
	wildcardPattern := `^\*\.[^.]+\.[^.]+$`
	matched, err := regexp.MatchString(wildcardPattern, domain)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	return matched
}

func DoScopes(assets []structure.InScope) {

	var s_type string

	for _, v := range assets {
		s_type = strings.ToLower(v.Type)
		if s_type == "url" || s_type == "wildcard" || s_type == "api" {
			if isWild(v.Asset) {
				// EnumerateSubs()
			} else {
				EnumerateTech(v.Asset)
			}
		} else {
			continue
		}
	}

}
