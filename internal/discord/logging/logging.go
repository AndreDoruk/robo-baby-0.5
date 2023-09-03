package logging

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

var logging_channel = os.Getenv("LOGGING_CHANNEL_ID")

func Log(session *discordgo.Session, text string) {
	_, err := session.ChannelMessageSendEmbed(logging_channel, &discordgo.MessageEmbed{
		Title:       "Log",
		Description: text,
	})

	if err != nil {
		log.Fatalln(err)
	}
}
