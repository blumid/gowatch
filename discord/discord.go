package discord

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func getEnv(key string) string {
	fmt.Println("key is: ", key)
	return os.Getenv(key)
}
func Annoy(a string) {
	value := getEnv("")
	// send message:
	fmt.Println(value)
}

func Connect() {
	dg, err1 := discordgo.New("Bot " + "<my bot token>")
	if err1 != nil {
		log.Fatal("discord.go - Connect() :", err1)
	}
	dg.AddHandler(func() {

	})

	err2 := dg.Open()
	if err2 != nil {
		log.Fatal("discord.go - Connect() :", err2)
	}
}

func init() {

	err1 := godotenv.Load("../.env")
	if err1 != nil {
		log.Fatal("discord.go - init() :", err1)
	}
}
