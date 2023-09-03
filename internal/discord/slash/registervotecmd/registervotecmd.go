package registervotecmd

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/timshannon/bolthold"
	"github.com/trustig/robobaby0.5/internal/discord/voting"
)

var channel_id string = os.Getenv("VOTING_CHANNEL_ID")

var COMMAND *discordgo.ApplicationCommand = &discordgo.ApplicationCommand{
	Name:        "register-vote",
	Type:        discordgo.ChatApplicationCommand,
	Description: "Register a message as a vote",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "userid",
			Description: "Id of the user to be voted on",
		},

		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "messageid",
			Description: "Id of the message to be turned into a vote",
		},
	},
}

func Command(session *discordgo.Session, commandData discordgo.ApplicationCommandInteractionData, interaction *discordgo.InteractionCreate) string {
	if len(commandData.Options) > 2 {
		return "Please specify both a message and a user"
	}

	userId := commandData.Options[0].Value.(string)
	messageId := commandData.Options[1].Value.(string)

	message, err := session.ChannelMessage(channel_id, messageId)

	if err != nil {
		return "Error while trying to create vote: ```yaml\n" + err.Error() + "```"
	}

	store, err := bolthold.Open("db/votes.db", 0666, nil)
	defer store.Close()

	if err != nil {
		return "Error while trying to create vote: ```yaml\n" + err.Error() + "```"
	}

	vote := voting.Vote{
		MessageId:   messageId,
		UserId:      userId,
		TimeStarted: message.Timestamp,
		LastHour:    -1,
	}

	err = store.Insert(userId, vote)

	if err != nil {
		return "Error while trying to create vote: ```yaml\n" + err.Error() + "```"
	}

	return "Successfuly registered message as vote for <@" + userId + ">"
}
