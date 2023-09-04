package commentgamecmd

import (
	"github.com/bwmarrin/discordgo"
	"github.com/trustig/robobaby0.5/internal/discord/commentgame"
)

var COMMAND *discordgo.ApplicationCommand = &discordgo.ApplicationCommand{
	Name:        "comment-game",
	Type:        discordgo.ChatApplicationCommand,
	Description: "Begin the comment game",
}

func Command(session *discordgo.Session, commandData discordgo.ApplicationCommandInteractionData, interaction *discordgo.InteractionCreate) string {
	commentgame.Begin(session, interaction)
	return ""
}
