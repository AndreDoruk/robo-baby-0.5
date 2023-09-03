package logging

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

var logging_channel = os.Getenv("LOGGING_CHANNEL")

func Log(session *discordgo.Session, text string) {
	session.ChannelMessageSendEmbed(logging_channel, &discordgo.MessageEmbed{
		Title:       "Log",
		Description: text,
	})
}
