package discord

import (
	"fmt"
	"os"
	"time"

	"github.com/blumid/gowatch/structure"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/kz/discordrus"
	"github.com/sirupsen/logrus"
)

var dg *discordgo.Session

func init() {

	err1 := godotenv.Load()
	if err1 != nil {
		logrus.Fatal("discord.init():", err1)
	}

	var err error
	dg, err = discordgo.New("Bot " + getEnv("Bot_Token"))
	if err != nil {
		logrus.Fatal("discord.init():", err)
	}

	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stderr)
	logrus.SetLevel(logrus.InfoLevel)

	logrus.AddHook(discordrus.NewHook(
		// Use environment variable for security reasons
		getEnv("WebHook_URL"),
		// Set minimum level to DebugLevel to receive all log entries
		logrus.InfoLevel,
		&discordrus.Opts{
			Username:           "captain hook",
			Author:             "",                   // Setting this to a non-empty string adds the author text to the message header
			DisableTimestamp:   false,                // Setting this to true will disable timestamps from appearing in the footer
			TimestampFormat:    "Jan 2 15:04:05 MST", // The timestamp takes this format; if it is unset, it will take logrus' default format
			TimestampLocale:    nil,                  // The timestamp uses this locale; if it is unset, it will use time.Local
			EnableCustomColors: true,                 // If set to true, the below CustomLevelColors will apply
			CustomLevelColors: &discordrus.LevelColors{
				Trace: 3092790,
				Debug: 10170623,
				Info:  3581519,
				Warn:  14327864,
				Error: 13631488,
				Panic: 13631488,
				Fatal: 13631488,
			},
			DisableInlineFields: false,
		},
	))

}

func getEnv(key string) string {
	return os.Getenv(key)
}

func Open() {

	// dg.AddHandler(messageHandler)

	// dg.AddHandler(replyHandler)

	err2 := dg.Open()
	if err2 != nil {
		logrus.Fatal("discord.Open(): ", err2)
	}
}

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

	// thumbnail
	thumb := discordgo.MessageEmbedThumbnail{
		URL:    "https://raw.githubusercontent.com/blumid/gowatch/main/static/hackerone.jpg",
		Width:  30,
		Height: 30,
	}
	switch p.Owner {
	case "hackerone":
		thumb.URL = "https://raw.githubusercontent.com/blumid/gowatch/main/static/hackerone.jpg"
	case "bugcrowd":
		thumb.URL = "https://raw.githubusercontent.com/blumid/gowatch/main/static/bugcrowd.jpg"
	case "intigriti":
		thumb.URL = "https://raw.githubusercontent.com/blumid/gowatch/main/static/intigriti.jpg"
	}

	cID := getEnv("ChannelId_general")
	embed := &discordgo.MessageEmbed{
		Title:       p.Name,
		URL:         p.Url,
		Description: "*" + p.Owner + "*",
		Color:       0xff6666,
		Thumbnail:   &thumb,
	}
	dg.ChannelMessageSendEmbed(cID, embed)
	return true
}

func NotifyNewAsset(p *structure.Program, s []structure.InScope) bool {
	cID := getEnv("ChannelId_general")

	// thumbnail
	thumb := discordgo.MessageEmbedThumbnail{
		URL:    "https://raw.githubusercontent.com/blumid/gowatch/main/static/hackerone.jpg",
		Width:  30,
		Height: 30,
	}
	switch p.Owner {
	case "hackerone":
		thumb.URL = "https://raw.githubusercontent.com/blumid/gowatch/main/static/hackerone.jpg"
	case "bugcrowd":
		thumb.URL = "https://raw.githubusercontent.com/blumid/gowatch/main/static/bugcrowd.jpg"
	case "intigriti":
		thumb.URL = "https://raw.githubusercontent.com/blumid/gowatch/main/static/intigriti.jpg"
	}

	// fields
	fields := []*discordgo.MessageEmbedField{}

	for _, item := range s {
		temp := discordgo.MessageEmbedField{
			Name:   item.Asset,
			Value:  item.Type,
			Inline: true}
		fields = append(fields, &temp)
	}
	embed := &discordgo.MessageEmbed{
		Title:       p.Name,
		URL:         p.Url,
		Description: "*" + p.Owner + "*",
		Color:       0x0080ff,
		Fields:      fields,
		Thumbnail:   &thumb,
	}
	dg.ChannelMessageSendEmbed(cID, embed)

	return true
}
