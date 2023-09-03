package whitelistedcmd

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/timshannon/bolthold"
	"github.com/trustig/robobaby0.5/internal/discord/whitelist"
)

var COMMAND *discordgo.ApplicationCommand = &discordgo.ApplicationCommand{
	Name:        "whitelisted",
	Type:        discordgo.ChatApplicationCommand,
	Description: "See all whitelisted users",
}

func Command(session *discordgo.Session, commandData discordgo.ApplicationCommandInteractionData, interaction *discordgo.InteractionCreate) string {
	store, err := bolthold.Open("db/whitelist.db", 0666, nil)

	if err != nil {
		log.Fatalln(err)
	}

	defer store.Close()

	text := "Whitelisted users: \n"

	store.ForEach(bolthold.Where("FavorVotes").Gt(-100), func(whitelist *whitelist.Whitelist) error {
		text += "<@" + whitelist.UserId + "> \n"
		return nil
	})

	return text
}
