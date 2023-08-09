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
	// fmt.Println("handler added 1.")

	dg.AddHandler(replyHandler)
	// fmt.Println("handler added 2.")

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
	switch m.Content {
	case "fuck":
		// s.ChannelMessageSendReply(m.ChannelID, "fucking reply", m.MessageReference)
		s.ChannelMessageSend(m.ChannelID, "fuck yourself!")

	}

}

func replyHandler(s *discordgo.Session, m *discordgo.MessageReactionAdd) {

	switch m.Emoji.Name {
	case "💩":
		s.ChannelMessageSend(m.ChannelID, "poop yourself 2 ! ")

	case "🔍":
		ref := discordgo.MessageReference{
			MessageID: m.MessageID,
			ChannelID: m.ChannelID,
			GuildID:   m.GuildID,
		}
		s.ChannelMessageSendReply(m.ChannelID, "starting for enum...", &ref)
		// s.ChannelMessageSend(m.ChannelID, "")
		// run gosub get final.json file and give to db/operations.go

	case "📃":

	case "🕷":
	}

}
