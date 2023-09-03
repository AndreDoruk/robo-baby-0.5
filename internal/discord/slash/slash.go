package slash

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/timshannon/bolthold"
	"github.com/trustig/robobaby0.5/internal/discord/voting"
	"github.com/trustig/robobaby0.5/internal/discord/whitelist"
)

var server_id string = os.Getenv("SERVER_ID")

var vote_command *discordgo.ApplicationCommand = &discordgo.ApplicationCommand{
	Name:        "vote",
	Type:        discordgo.ChatApplicationCommand,
	Description: "Initiate a vote",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "userid",
			Description: "Id of the user to be voted on",
		},
	},
}

var whitelist_command *discordgo.ApplicationCommand = &discordgo.ApplicationCommand{
	Name:        "whitelist",
	Type:        discordgo.ChatApplicationCommand,
	Description: "Whitelist someone",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "userid",
			Description: "Id of the user",
		},
	},
}

func CreateCommands(session *discordgo.Session) {
	appId := session.State.Application.ID

	session.ApplicationCommandCreate(appId, server_id, vote_command)
	session.ApplicationCommandCreate(appId, server_id, whitelist_command)
}

func OnInteract(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	commandData := interaction.ApplicationCommandData()

	switch commandData.Name {
	case vote_command.Name:
		{
			if len(commandData.Options) == 0 {
				return
			}

			userId := commandData.Options[0].Value
			voting.CreateVote(session, userId.(string))
		}

	case whitelist_command.Name:
		{
			if len(commandData.Options) == 0 {
				return
			}

			store, err := bolthold.Open("db/whitelist.db", 0666, nil)

			if err != nil {
				log.Fatalln(err)
			}

			userId := commandData.Options[0].Value
			store.Insert(userId, whitelist.Whitelist{FavorVotes: 69, AgainstVotes: 420})
		}
	}

}
