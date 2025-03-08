package balancecmd

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/AndreDoruk/robo-baby-0.5/internal/database"
)

var COMMAND *discordgo.ApplicationCommand = &discordgo.ApplicationCommand{
	Name:        "balance",
	Type:        discordgo.ChatApplicationCommand,
	Description: "Show your balance",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "userid",
			Description: "Id of the user to check",
			Required:    false,
		},
	},
}

func Command(session *discordgo.Session, commandData discordgo.ApplicationCommandInteractionData, interaction *discordgo.InteractionCreate) string {
	balances := make(map[string]int)
	database.LoadJson("db/balance.json", &balances)

	if len(commandData.Options) == 0 {
		return "Your balance: " + strconv.Itoa(balances[interaction.Member.User.ID]) + " 🍅"
	} else {
		userId := commandData.Options[0].StringValue()

		_, err := session.User(userId)
		if err != nil {
			return "```go\n" + err.Error() + "```"
		}

		return "<@" + userId + ">'s balance: " + strconv.Itoa(balances[userId]) + " 🍅"
	}
}
