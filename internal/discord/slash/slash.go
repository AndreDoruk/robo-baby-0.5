package slash

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/timshannon/bolthold"
	"github.com/trustig/robobaby0.5/internal/discord/logging"
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

var whitelist_all_command *discordgo.ApplicationCommand = &discordgo.ApplicationCommand{
	Name:        "whitelist-all",
	Type:        discordgo.ChatApplicationCommand,
	Description: "Whitelist everyone currently on the server",
}

func CreateCommands(session *discordgo.Session) {
	appId := session.State.Application.ID

	session.ApplicationCommandCreate(appId, server_id, vote_command)
	session.ApplicationCommandCreate(appId, server_id, whitelist_command)
	session.ApplicationCommandCreate(appId, server_id, whitelist_all_command)
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

			respondInteraction(session, interaction, "Successfuly started vote")
			logging.Log(session, commandData.Name+" used by "+interaction.User.Username)
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
			err = store.Insert(userId, whitelist.Whitelist{FavorVotes: 69, AgainstVotes: 420})

			if err == nil {
				respondInteraction(session, interaction, "Sucessfuly whitelisted user")
			} else {
				respondInteraction(session, interaction, "Error when whitelisting user: "+err.Error())
			}
			logging.Log(session, commandData.Name+" used by "+interaction.User.Username)
		}

	case whitelist_all_command.Name:
		{
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

			respondInteraction(session, interaction, "Sucessfuly whitelisted everyone in the server")
			logging.Log(session, commandData.Name+" used by "+interaction.User.Username)
		}
	}

}

func respondInteraction(session *discordgo.Session, interaction *discordgo.InteractionCreate, text string) {
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:         text,
			AllowedMentions: &discordgo.MessageAllowedMentions{},
		},
	})
}
