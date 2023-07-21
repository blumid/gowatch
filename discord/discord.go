package discord

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func getEnv(key string) string {
	return os.Getenv(key)
}
func Annoy(a string) {
	value := getEnv("")
	// send message:
	fmt.Println(value)
}

func Connect() {
	dg, err1 := discordgo.New("Bot " + getEnv("Bot_Token"))
	if err1 != nil {
		log.Fatal("discord.go - Connect() :", err1)
	}
	dg.AddHandler(messageHandler)
	fmt.Println("handler added.")

	err2 := dg.Open()
	if err2 != nil {
		log.Fatal("discord.go - Connect() :", err2)
	}
}

func init() {

	err1 := godotenv.Load(".env")
	if err1 != nil {
		log.Fatal("discord.go - init() :", err1)
	}
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println(m.Content)
	fmt.Println(m.ChannelID)
	switch m.Content {
	case "fuck":
		s.ChannelMessageSend("1130071729067261972", "fuck yourself!")

	}
}
