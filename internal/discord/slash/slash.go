package slash

import (
	"math/rand"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/trustig/robobaby0.5/internal/discord/logging"
	"github.com/trustig/robobaby0.5/internal/discord/slash/vote"
	"github.com/trustig/robobaby0.5/internal/discord/slash/whitelist"
	"github.com/trustig/robobaby0.5/internal/discord/slash/whitelistall"
)

var server_id string = os.Getenv("SERVER_ID")

type command_func func(*discordgo.Session, discordgo.ApplicationCommandInteractionData, *discordgo.InteractionCreate) string

type Command struct {
	Command  *discordgo.ApplicationCommand
	Function command_func
}

var commands = map[string]Command{
	"vote":          Command{vote.COMMAND, vote.Command},
	"whitelist":     Command{whitelist.COMMAND, whitelist.Command},
	"whitelist-all": Command{whitelistall.COMMAND, whitelistall.Command},
}

func CreateCommands(session *discordgo.Session) {
	appId := session.State.Application.ID

	for _, command := range commands {
		session.ApplicationCommandCreate(appId, server_id, command.Command)
	}
}

func OnInteract(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	commandData := interaction.ApplicationCommandData()
	command, ok := commands[commandData.Name]

	if ok {
		response := command.Function(session, commandData, interaction)
		respondInteraction(session, interaction, commandData, response)

		logging.LogCommand(session, interaction, commandData)
	}
}

func respondInteraction(session *discordgo.Session, interaction *discordgo.InteractionCreate, commandData discordgo.ApplicationCommandInteractionData, text string) {
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Author:      &discordgo.MessageEmbedAuthor{Name: interaction.Member.Nick + " used '/" + commandData.Name + "'", IconURL: interaction.Member.User.AvatarURL("128")},
					Color:       rand.Intn(16777215),
					Description: text,
				},
			},
		},
	})
}
