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

var dg *discordgo.Session

func init() {

	err1 := godotenv.Load()
	if err1 != nil {
		log.Fatal("Error loading .env file:", err1)
	}

	var err error
	dg, err = discordgo.New("Bot " + getEnv("Bot_Token"))
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func getEnv(key string) string {
	return os.Getenv(key)
}

func Connect() {

	// dg.AddHandler(messageHandler)

	// dg.AddHandler(replyHandler)

	err2 := dg.Open()
	if err2 != nil {
		log.Fatal("discord.go - Connect:", err2)
	}
}

// func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
// 	switch m.Content {
// 	case "fuck":
// 		// s.ChannelMessageSendReply(m.ChannelID, "fucking reply", m.MessageReference)
// 		s.ChannelMessageSend(m.ChannelID, "fuck yourself!")

// 	}

// }

// func replyHandler(s *discordgo.Session, m *discordgo.MessageReactionAdd) {

// 	switch m.Emoji.Name {
// 	case "💩":
// 		s.ChannelMessageSend(m.ChannelID, "poop yourself 2 ! ")

// 	case "🔍":
// 		ref := discordgo.MessageReference{
// 			MessageID: m.MessageID,
// 			ChannelID: m.ChannelID,
// 			GuildID:   m.GuildID,
// 		}
// 		s.ChannelMessageSendReply(m.ChannelID, "starting for enum...", &ref)
// 		// s.ChannelMessageSend(m.ChannelID, "")
// 		// run gosub get final.json file and give to db/operations.go

// 	case "📃":

// 	case "🕷":
// 	}

// }

func Timer() {
	/*
		set timer()
	*/
	now := time.Now()

	year, month, day := now.Year(), now.Month(), now.Day()

	desiredTime := time.Date(year, month, day, 12, 0, 0, 0, time.Local)

	duration := desiredTime.Sub(now)

	fmt.Println(duration)
}

func NotifyNewProgram(p *structure.Program) bool {
	cID := getEnv("ChannelId_general")
	embed := &discordgo.MessageEmbed{
		Title:       p.Name,
		URL:         p.Url,
		Description: "*newProgram*",
		Timestamp:   time.Now().Format("2006-1-2 15:4:5"),
		Color:       0xff6666,
	}
	dg.ChannelMessageSendEmbed(cID, embed)
	return true
}

func NotifyNewAsset(p *structure.Program, s []structure.InScope) bool {
	cID := getEnv("ChannelId_general")
	fields := []*discordgo.MessageEmbedField{}
	for _, item := range s {
		temp := discordgo.MessageEmbedField{
			Name:   item.AssetIdentifier,
			Value:  item.AssetType,
			Inline: true}
		fields = append(fields, &temp)
	}
	embed := &discordgo.MessageEmbed{
		Title:       p.Name,
		URL:         p.Url,
		Description: "*newAsset*",
		Timestamp:   time.Now().Format("2006-1-2 15:4:5"),
		Color:       0x0080ff,
		Fields:      fields,
	}
	dg.ChannelMessageSendEmbed(cID, embed)

	return true
}
