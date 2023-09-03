package whitelistedcmd

import (
	"github.com/bwmarrin/discordgo"
	"github.com/trustig/robobaby0.5/internal/database"
	"github.com/trustig/robobaby0.5/internal/discord/whitelist"
)

var COMMAND *discordgo.ApplicationCommand = &discordgo.ApplicationCommand{
	Name:        "whitelisted",
	Type:        discordgo.ChatApplicationCommand,
	Description: "See all whitelisted users",
}

func Command(session *discordgo.Session, commandData discordgo.ApplicationCommandInteractionData, interaction *discordgo.InteractionCreate) string {
	whitelists := make(map[string]whitelist.Whitelist)
	database.LoadJson("db/whitelist.json", &whitelists)
	defer database.SaveJson("db/whitelist.json", whitelists)

	text := "Whitelisted users: \n"

	for _, whitelist := range whitelists {
		text += "<@" + whitelist.UserId + "> \n"
	}

	return text
}
