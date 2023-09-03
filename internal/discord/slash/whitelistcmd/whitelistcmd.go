package whitelistcmd

import (
	"github.com/bwmarrin/discordgo"
	"github.com/trustig/robobaby0.5/internal/database"
	"github.com/trustig/robobaby0.5/internal/discord/whitelist"
)

var COMMAND *discordgo.ApplicationCommand = &discordgo.ApplicationCommand{
	Name:        "whitelist",
	Type:        discordgo.ChatApplicationCommand,
	Description: "Whitelist someone",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "userid",
			Description: "Id of the user",
		},
	},
}

func Command(session *discordgo.Session, commandData discordgo.ApplicationCommandInteractionData, interaction *discordgo.InteractionCreate) string {
	if len(commandData.Options) == 0 {
		return "Please specify user"
	}

	whitelists := make(map[string]whitelist.Whitelist)
	database.LoadJson("db/whitelist.json", &whitelists)
	defer database.SaveJson("db/whitelist.json", whitelists)

	userId := commandData.Options[0].Value.(string)
	whitelists[userId] = whitelist.Whitelist{FavorVotes: -1, AgainstVotes: -1, UserId: userId}

	/*if err != nil {
		return "Error while whitelisting user: ```yaml\n" + err.Error() + "```"
	} else { */

	return "Succesfuly whitelisted user <@" + userId + ">"
	//}
}
