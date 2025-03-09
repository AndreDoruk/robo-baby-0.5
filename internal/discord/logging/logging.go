package logging

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/bwmarrin/discordgo"
)

var isTesting bool = os.Getenv("TESTING") != ""
var logging_channel = os.Getenv("LOGGING_CHANNEL_ID")

/// Logs a string to the logging channel
/// In testing mode also logs it to console
func LogString(session *discordgo.Session, toLog string) {
	if isTesting {
		fmt.Println("Logged string: " + toLog)
	}
	
	_, err := session.ChannelMessageSend(logging_channel, toLog)

	if err != nil {
		log.Fatalln(err)
	}
}
/// Logs a command's name, user and option values to the logging channel
/// In testing mode also logs them to the console
func LogCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate, commandData discordgo.ApplicationCommandInteractionData) {
	hasOptions := len(commandData.Options) > 0
	fieldNum := 1

	if isTesting {
		fmt.Println("Logged command: " + "'/" + commandData.Name + "' used by <@" + interaction.Member.User.ID + ">")
		if hasOptions {
			fmt.Println("     Command options:")
			for _, option := range commandData.Options {
				fmt.Println("                     " + option.Name + ": " + option.StringValue())
			}
		}
	}

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
/// Logs an error in the logging channel
/// In testing mode also logs it in the console
func LogError(session *discordgo.Session, activity string, err error) {
	if isTesting {
		fmt.Println("Error while " + activity)
	}

	_, err = session.ChannelMessageSendEmbed(logging_channel, &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{Name: "Error while " + activity, IconURL: session.State.User.AvatarURL("128")},
		Color:       rand.Intn(16777215),
		Description: "```go\n" + err.Error() + "```",
	})

	if err != nil {
		log.Fatalln(err)
	}
}
