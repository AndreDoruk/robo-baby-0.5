package whitelistallcmd

import (
	"log"
	"os"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/trustig/robobaby0.5/internal/database"
	"github.com/trustig/robobaby0.5/internal/discord/whitelist"
)

var server_id string = os.Getenv("SERVER_ID")

var COMMAND *discordgo.ApplicationCommand = &discordgo.ApplicationCommand{
	Name:        "whitelist-all",
	Type:        discordgo.ChatApplicationCommand,
	Description: "Whitelist everyone currently on the server",
}

func Command(session *discordgo.Session, commandData discordgo.ApplicationCommandInteractionData, interaction *discordgo.InteractionCreate) string {
	whitelists := make(map[string]whitelist.Whitelist)
	database.LoadJson("db/whitelist.json", &whitelists)
	defer database.SaveJson("db/whitelist.json", whitelists)

	members, err := session.GuildMembers(server_id, "2006", 1000)

	if err != nil {
		log.Fatalln(err)
	}

	for _, member := range members {
		if member.User.Bot {
			continue
		}

		whitelists[member.User.ID] = whitelist.Whitelist{FavorVotes: -1, AgainstVotes: -1, UserId: member.User.ID}
	}

	return "Successfuly whitelisted all " + strconv.Itoa(len(members)) + " users"
}
