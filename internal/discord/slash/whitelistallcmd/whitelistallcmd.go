package whitelistallcmd

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/timshannon/bolthold"
	"github.com/trustig/robobaby0.5/internal/discord/whitelist"
)

var server_id string = os.Getenv("SERVER_ID")

var COMMAND *discordgo.ApplicationCommand = &discordgo.ApplicationCommand{
	Name:        "whitelist-all",
	Type:        discordgo.ChatApplicationCommand,
	Description: "Whitelist everyone currently on the server",
}

func Command(session *discordgo.Session, commandData discordgo.ApplicationCommandInteractionData, interaction *discordgo.InteractionCreate) string {
	store, err := bolthold.Open("db/whitelist.db", 0666, nil)

	if err != nil {
		log.Fatalln(err)
	}

	defer store.Close()

	members, err := session.GuildMembers(server_id, "2006", 1000)

	if err != nil {
		log.Fatalln(err)
	}

	for _, member := range members {
		if member.User.Bot {
			continue
		}

		err = store.Insert(member.User.ID, whitelist.Whitelist{FavorVotes: -1, AgainstVotes: -1, UserId: member.User.ID})

		if err != nil {
			fmt.Println(err)
		}
	}

	return "Successfuly whitelisted all " + strconv.Itoa(len(members)) + " users"
}
