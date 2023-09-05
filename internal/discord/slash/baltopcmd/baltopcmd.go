package baltopcmd

import (
	"sort"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/trustig/robobaby0.5/internal/database"
)

type userTomatoes struct {
	id  string
	num int
}

var COMMAND *discordgo.ApplicationCommand = &discordgo.ApplicationCommand{
	Name:        "baltop",
	Type:        discordgo.ChatApplicationCommand,
	Description: "See the wealthiest people on the server",
}

const text_length int = 25

func Command(session *discordgo.Session, commandData discordgo.ApplicationCommandInteractionData, interaction *discordgo.InteractionCreate) string {
	userId := interaction.Member.User.ID

	tomatoes := make(map[string]int)
	database.LoadJson("db/balance.json", &tomatoes)

	sortedTomatoes := sortByValue(tomatoes)
	leaderboard := ""

	for i := 0; i < 10; i++ {
		if len(sortedTomatoes) <= i {
			continue
		}

		userTomato := sortedTomatoes[i]

		text := "**#" + strconv.Itoa(i+1) + "**: <@" + userTomato.id + ">\n"
		leaderboard += text
	}

	leaderboard += `‗‗‗‗‗‗‗‗‗‗‗‗‗‗‗‗‗‗‗
	Your placement: **` + getPlacementFromUserId(userId, sortedTomatoes) + "** (" + strconv.Itoa(tomatoes[userId]) + " 🍅)"

	return leaderboard
}

func getPlacementFromUserId(userId string, userTomatoes []userTomatoes) string {
	index := 0
	for _, value := range userTomatoes {
		if value.id == userId {
			return "#" + strconv.Itoa(index+1)
		}
		index += 1
	}
	return "Unranked"
}

func sortByValue(tomatoes map[string]int) []userTomatoes {
	var sortedArray []userTomatoes

	for key, value := range tomatoes {
		sortedArray = append(sortedArray, userTomatoes{key, value})
	}

	sort.Slice(sortedArray, func(index1 int, index2 int) bool {
		return sortedArray[index1].num > sortedArray[index2].num
	})

	return sortedArray
}
