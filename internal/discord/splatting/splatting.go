package splatting

import (
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/trustig/robobaby0.5/internal/database"
	"github.com/trustig/robobaby0.5/internal/discord/starboard"
	"github.com/trustig/robobaby0.5/internal/discord/logging"
)

var server_id string = os.Getenv("SERVER_ID")
var splatted_role string = os.Getenv("SPLATTED_ROLE")

const required_reaction_count int = 5

const splat_window time.Duration = 15 * time.Minute
const splat_cooldown time.Duration = 15 * time.Minute

const timeout_time time.Duration = 30 * time.Second

func OnReact(session *discordgo.Session, reactionAdd *discordgo.MessageReactionAdd) {
	message, err := session.ChannelMessage(reactionAdd.ChannelID, reactionAdd.MessageID)

	if err != nil {
		logging.LogError(session, "splatting message (get message) "+reactionAdd.MessageID, err)
		return
	}

	if !shouldSplatMessage(session, message) {
		return
	}

	targetMember, err := session.GuildMember(server_id, message.Author.ID)

	if err != nil {
		logging.LogError(session, "splatting message (get member) "+reactionAdd.MessageID, err)
		return
	}

	if hasSplatRole(targetMember) {
		return
	}

	splattedMessages := make(map[string]bool)
	database.LoadJson("db/messagesplats.json", &splattedMessages)

	if splattedMessages[reactionAdd.MessageID] {
		return
	}

	splattedMessages[reactionAdd.MessageID] = true
	database.SaveJson("db/messagesplats.json", &splattedMessages)

	timeout_time := time.Now().Add(timeout_time)
	err = session.GuildMemberTimeout(server_id, targetMember.User.ID, &timeout_time)

	if err != nil {
		logging.LogError(session, "splatting message (timeout user) "+reactionAdd.MessageID, err)
		return
	}

	reactors, err := session.MessageReactions(message.ChannelID, message.ID, "🍅", required_reaction_count, "", "")

	if err != nil {
		logging.LogError(session, "splatting message (get reactors)"+reactionAdd.MessageID, err)
		return
	}

	_, err = session.ChannelMessageSendEmbedReply(message.ChannelID, &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{IconURL: "https://media.tenor.com/Ix7ojVpGhRUAAAAC/tomato-explode.gif", Name: "Splatted!"},
		Description: "<@" + targetMember.User.ID + "> has been splatted! \n Timeout ends in " + "<t:" + strconv.Itoa(int(timeout_time.Unix())) + ":R>",
		Fields: []*discordgo.MessageEmbedField{{
			Name:  "Splatters",
			Value: getSplatters(reactors),
		}},
		Color: rand.Intn(16777215),
	}, message.Reference())

	starboard.SendBoardMessage(session, message, &discordgo.MessageEmbedFooter{
			IconURL: "https://media.tenor.com/Ix7ojVpGhRUAAAAC/tomato-explode.gif",
			Text: "Splatted!",
	})

	if err != nil {
		logging.LogError(session, "splatting message (sending embed) "+reactionAdd.MessageID, err)
		return
	}

	err = session.GuildMemberRoleAdd(server_id, targetMember.User.ID, splatted_role)

	if err != nil {
		logging.LogError(session, "splatting message (adding role) "+reactionAdd.MessageID, err)
		return
	}

	splattedUsers := make(map[string]time.Time)
	database.LoadJson("db/usersplats.json", &splattedUsers)
	defer database.SaveJson("db/usersplats.json", &splattedUsers)

	splattedUsers[targetMember.User.ID] = time.Now().Add(splat_cooldown)
}

func shouldSplatMessage(session *discordgo.Session, message *discordgo.Message) bool {
	if message.Author.ID == session.State.User.ID {
		return false
	}

	if time.Now().After(message.Timestamp.Add(splat_window)) {
		return false
	}

	for _, reaction := range message.Reactions {
		if reaction.Emoji.Name == "🍅" {
			return reaction.Count == required_reaction_count
		}
	}

	return false
}

func hasSplatRole(member *discordgo.Member) bool {
	for _, role := range member.Roles {
		if role == splatted_role {
			return true
		}
	}
	return false
}

func getSplatters(reactors []*discordgo.User) string {
	text := ""
	for _, user := range reactors {
		text += "<@" + user.ID + ">\n"
	}
	return text
}

func UpdateSplattedRole(session *discordgo.Session) {
	splattedUsers := make(map[string]time.Time)
	database.LoadJson("db/usersplats.json", &splattedUsers)
	defer database.SaveJson("db/usersplats.json", splattedUsers)

	now := time.Now()

	for userId, endTime := range splattedUsers {
		if !now.After(endTime) {
			continue
		}

		err := session.GuildMemberRoleRemove(server_id, userId, splatted_role)

		if err != nil {
			logging.LogError(session, " removing splat role from <@"+userId+">", err)
			continue
		}

		delete(splattedUsers, userId)
	}
}
