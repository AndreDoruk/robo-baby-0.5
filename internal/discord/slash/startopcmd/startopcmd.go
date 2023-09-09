package startopcmd

import (
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/trustig/robobaby0.5/internal/database"
	"github.com/trustig/robobaby0.5/internal/discord/starboard"
)

var COMMAND *discordgo.ApplicationCommand = &discordgo.ApplicationCommand{
	Name:        "startop",
	Type:        discordgo.ChatApplicationCommand,
	Description: "Shows starboard stats",
}

var server_id string = os.Getenv("SERVER_ID")

const message_url_prefix string = "https://discordapp.com/channels/"

const leaderboard_entry_num int = 3

type StarredUser struct {
	UserID  string
	StarNum int
}

type ArrayMessage struct {
	MessageID      string
	StarredMessage starboard.StarredMessage
}

func Command(session *discordgo.Session, commandData discordgo.ApplicationCommandInteractionData, interaction *discordgo.InteractionCreate) string {
	boardMessages := make(map[string]starboard.StarredMessage)
	database.LoadJson("db/boardmessages.json", &boardMessages)

	boardArray := toArray(boardMessages)

	totalStars := getTotalStars(boardArray)
	topMessages := getTopMessages(session, boardArray)
	topRecievers := getTopRecievers(boardArray)

	embed := discordgo.MessageEmbed{
		Title:       "Starboard Stats",
		Description: strconv.Itoa(len(boardArray)-1) + " starred messages with a total of " + strconv.Itoa(totalStars) + " <:quality5:1146794549210001511>",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Top Messages",
				Value: topMessages,
			},

			{
				Name:  "Top Members",
				Value: topRecievers,
			},
		},
		Color: rand.Intn(16777215),
	}

	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{&embed},
		},
	})

	return ""
}

func toArray(mapToConvert map[string]starboard.StarredMessage) []ArrayMessage {
	array := make([]ArrayMessage, len(mapToConvert)-1)

	for messageId, value := range mapToConvert {
		array = append(array, ArrayMessage{messageId, value})
	}

	return array
}

func getTotalStars(boardArray []ArrayMessage) int {
	totalStars := 0
	for _, starredMessage := range boardArray {
		totalStars += starredMessage.StarredMessage.StarNum
	}
	return totalStars
}

func getTopMessages(session *discordgo.Session, boardArray []ArrayMessage) string {
	sort.Slice(boardArray, func(index1 int, index2 int) bool {
		return boardArray[index1].StarredMessage.StarNum > boardArray[index2].StarredMessage.StarNum
	})

	returnString := ""

	for index := 0.0; index < math.Min(float64(leaderboard_entry_num), float64(len(boardArray)-1)); index++ {
		message := boardArray[int(index)]
		returnString += "[Message by](" + message_url_prefix + server_id + "/" + message.StarredMessage.ChannelID + "/" + message.MessageID + ") <@" + message.StarredMessage.UserID + "> - " + strconv.Itoa(message.StarredMessage.StarNum) + " <:quality5:1146794549210001511>\n"
	}

	return returnString
}

func getTopRecievers(boardArray []ArrayMessage) string {
	starredUsers := make(map[string]int)

	for _, messageData := range boardArray {
		starredUsers[messageData.StarredMessage.UserID] += messageData.StarredMessage.StarNum
	}

	sortedUsers := make([]StarredUser, len(starredUsers))

	for userId, starNum := range starredUsers {
		sortedUsers = append(sortedUsers, StarredUser{userId, starNum})
	}

	sort.Slice(sortedUsers, func(index1 int, index2 int) bool {
		return sortedUsers[index1].StarNum > sortedUsers[index2].StarNum
	})

	returnString := ""

	for index := 0.0; index < math.Min(float64(leaderboard_entry_num), float64(len(boardArray)-1)); index++ {
		user := sortedUsers[int(index)]
		returnString += "<@" + user.UserID + "> - " + strconv.Itoa(user.StarNum) + " <:quality5:1146794549210001511> \n"
	}

	return returnString
}
