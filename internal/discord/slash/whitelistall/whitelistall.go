package whitelistall

import (
	"fmt"
	"log"
	"os"

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
	guild, err := session.Guild(server_id)

	if err != nil {
		log.Fatalln(err)
	}

	store, err := bolthold.Open("db/whitelist.db", 0666, nil)

	if err != nil {
		log.Fatalln(err)
	}

	for _, member := range guild.Members {
		store.Insert(member.User.ID, whitelist.Whitelist{FavorVotes: -1, AgainstVotes: -1})
		fmt.Println("Whitelisted " + member.User.Username)
	}

	return "Successfuly whitelisted all users"
}
