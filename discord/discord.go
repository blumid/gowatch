package discord

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/blumid/gowatch/structure"
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

	//run download github file at a specific time

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

func magnifier() {
	var domain *structure.Domain

	// run gosub

	// get directory and run addsub for each domain

	sub := &structure.Sub{}
	domain.Subs = append(domain.Subs, *sub)

	// call db.AddSub()

	//
}

func timer() {
	/*
		set timer()
	*/
	now := time.Now()

	year, month, day := now.Year(), now.Month(), now.Day()

	desiredTime := time.Date(year, month, day, 12, 0, 0, 0, time.Local)

	duration := desiredTime.Sub(now)

	fmt.Println(duration)
}
