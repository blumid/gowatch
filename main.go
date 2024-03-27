package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"os"

	"github.com/blumid/gowatch/db"
	"github.com/blumid/gowatch/discord"
	"github.com/blumid/gowatch/structure"
)

func main() {
	// connect discord bot:
	discord.Connect()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigchan
	fmt.Println("Received signal:", sig)

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

		res := db.FandU(v.Name, v.Target.InScope)
		if !res {
			if err := db.AddProgram(&v); err != nil {
				fmt.Println("new one add: ", v.Name)
			}

		}

	}

}

func download() {
	res, err := http.Get("https://github.com/arkadiyt/bounty-targets-data/blob/main/data/hackerone_data.json")
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	_, err2 := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err2)
	}
}
