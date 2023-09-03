package vote

import (
	"github.com/bwmarrin/discordgo"
	"github.com/trustig/robobaby0.5/internal/discord/voting"
)

var COMMAND *discordgo.ApplicationCommand = &discordgo.ApplicationCommand{
	Name:        "vote",
	Type:        discordgo.ChatApplicationCommand,
	Description: "Initiate a vote",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "userid",
			Description: "Id of the user to be voted on",
		},
	},
}

func Command(session *discordgo.Session, commandData discordgo.ApplicationCommandInteractionData, interaction *discordgo.InteractionCreate) string {
	if len(commandData.Options) == 0 {
		return "Please specify a user"
	}

	userId := commandData.Options[0].Value
	err := voting.CreateVote(session, userId.(string))

	if err == nil {
		return "Successfuly made vote"
	} else {
		return "Error while making vote: ```go\n" + err.Error() + "```"
	}
}
