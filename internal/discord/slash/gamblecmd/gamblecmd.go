package gamblecmd

import (
	"math/rand"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/AndreDoruk/robo-baby-0.5/internal/database"
)

var COMMAND *discordgo.ApplicationCommand = &discordgo.ApplicationCommand{
	Name:        "gamble",
	Type:        discordgo.ChatApplicationCommand,
	Description: "Gamble",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "tomatoes",
			Description: "Amount of tomatoes to gamble",
		},
	},
}

func Command(session *discordgo.Session, commandData discordgo.ApplicationCommandInteractionData, interaction *discordgo.InteractionCreate) string {
	if len(commandData.Options) == 0 {
		return "Please specify the amount to gamble"
	}

	userId := interaction.Member.User.ID

	tomatoes := make(map[string]int)
	database.LoadJson("db/balance.json", &tomatoes)
	defer database.SaveJson("db/balance.json", tomatoes)

	gambleNum, err := strconv.Atoi(commandData.Options[0].Value.(string))

	if err != nil {
		return "Bro can you input like an actual number: ```yaml\n" + err.Error() + "```"
	}

	if gambleNum < 0 {
		return "🐎"
	}

	if gambleNum > tomatoes[userId] {
		return "You lack the sufficient 🍅 to gamble" //BROKE N
	}

	var verb string

	if rand.Intn(2) == 0 {
		tomatoes[userId] += gambleNum
		verb = "win"
	} else {
		tomatoes[userId] -= gambleNum
		verb = "lost"
	}

	return "You " + verb + " " + strconv.Itoa(gambleNum) + " 🍅"
}
