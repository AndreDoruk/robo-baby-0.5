package balancecmd

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/trustig/robobaby0.5/internal/database"
)

var COMMAND *discordgo.ApplicationCommand = &discordgo.ApplicationCommand{
	Name:        "balance",
	Type:        discordgo.ChatApplicationCommand,
	Description: "Show your balance",
}

func Command(session *discordgo.Session, commandData discordgo.ApplicationCommandInteractionData, interaction *discordgo.InteractionCreate) string {
	balances := make(map[string]int)
	database.LoadJson("db/balance.json", &balances)

	return strconv.Itoa(balances[interaction.Member.User.ID]) + " 🍅"
}
