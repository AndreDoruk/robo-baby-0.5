package logging

import (
	"log"
	"math/rand"
	"os"

	"github.com/bwmarrin/discordgo"
)

var logging_channel = os.Getenv("LOGGING_CHANNEL_ID")

func LogCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate, commandData discordgo.ApplicationCommandInteractionData) {
	hasOptions := len(commandData.Options) > 0
	fieldNum := 1

	if hasOptions {
		fieldNum += 1
	}

	messageFields := make([]*discordgo.MessageEmbedField, fieldNum)

	messageFields[0] = &discordgo.MessageEmbedField{Name: "Command", Value: "'/" + commandData.Name + "' used by <@" + interaction.Member.User.ID + ">"}

	if hasOptions {
		fieldText := "```yaml\n"

		for _, option := range commandData.Options {
			fieldText += option.Name + ": " + option.StringValue() + "\n"
		}

		fieldText += "```"
		messageFields[1] = &discordgo.MessageEmbedField{Name: "Values", Value: fieldText}
	}

	_, err := session.ChannelMessageSendEmbed(logging_channel, &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{Name: "<@" + interaction.Member.User.ID + ">", IconURL: interaction.Member.User.AvatarURL("128")},
		Color:  rand.Intn(16777215),
		Fields: messageFields,
	})

	if err != nil {
		log.Fatalln(err)
	}
}
func LogError(session *discordgo.Session, activity string, err error) {
	_, err = session.ChannelMessageSendEmbed(logging_channel, &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{Name: "Error while " + activity, IconURL: session.State.User.AvatarURL("128")},
		Color:       rand.Intn(16777215),
		Description: "```go\n" + err.Error() + "```",
	})

	if err != nil {
		log.Fatalln(err)
	}
}
