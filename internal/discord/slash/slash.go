package slash

import (
	"math/rand"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/trustig/robobaby0.5/internal/discord/logging"
	"github.com/trustig/robobaby0.5/internal/discord/slash/balancecmd"
	"github.com/trustig/robobaby0.5/internal/discord/slash/baltopcmd"
	"github.com/trustig/robobaby0.5/internal/discord/slash/commentgamecmd"
	"github.com/trustig/robobaby0.5/internal/discord/slash/gamblecmd"
	"github.com/trustig/robobaby0.5/internal/discord/slash/registervotecmd"
	"github.com/trustig/robobaby0.5/internal/discord/slash/startopcmd"
	"github.com/trustig/robobaby0.5/internal/discord/slash/votecmd"
	"github.com/trustig/robobaby0.5/internal/discord/slash/whitelistallcmd"
	"github.com/trustig/robobaby0.5/internal/discord/slash/whitelistcmd"
	"github.com/trustig/robobaby0.5/internal/discord/slash/whitelistedcmd"
)

var server_id string = os.Getenv("SERVER_ID")

type command_func func(*discordgo.Session, discordgo.ApplicationCommandInteractionData, *discordgo.InteractionCreate) string

type Command struct {
	Command  *discordgo.ApplicationCommand
	Function command_func
}

var commands = map[string]Command{
	"vote":          Command{votecmd.COMMAND, votecmd.Command},
	"whitelist":     Command{whitelistcmd.COMMAND, whitelistcmd.Command},
	"whitelist-all": Command{whitelistallcmd.COMMAND, whitelistallcmd.Command},
	"whitelisted":   Command{whitelistedcmd.COMMAND, whitelistedcmd.Command},
	"register-vote": Command{registervotecmd.COMMAND, registervotecmd.Command},
	"balance":       Command{balancecmd.COMMAND, balancecmd.Command},
	"comment-game":  Command{commentgamecmd.COMMAND, commentgamecmd.Command},
	"gamble":        Command{gamblecmd.COMMAND, gamblecmd.Command},
	"baltop":        Command{baltopcmd.COMMAND, baltopcmd.Command},
	"startop":       Command{startopcmd.COMMAND, startopcmd.Command},
}

func CreateCommands(session *discordgo.Session) {
	appId := session.State.Application.ID

	for _, command := range commands {
		session.ApplicationCommandCreate(appId, server_id, command.Command)
	}
}

func OnInteract(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Type != discordgo.InteractionApplicationCommand {
		return
	}

	commandData := interaction.ApplicationCommandData()
	command, ok := commands[commandData.Name]

	if ok {
		go func() {
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponsePong,
			})
		}()

		response := command.Function(session, commandData, interaction)
		if response != "" {
			respondInteraction(session, interaction, commandData, response)
		}

		logging.LogCommand(session, interaction, commandData)
	}
}

func respondInteraction(session *discordgo.Session, interaction *discordgo.InteractionCreate, commandData discordgo.ApplicationCommandInteractionData, text string) {
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Author:      &discordgo.MessageEmbedAuthor{Name: interaction.Member.User.Username + " used '/" + commandData.Name + "'", IconURL: interaction.Member.User.AvatarURL("128")},
					Color:       rand.Intn(16777215),
					Description: text,
				},
			},
		},
	})
}
