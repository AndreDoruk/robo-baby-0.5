package starboard

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
	"github.com/trustig/robobaby0.5/internal/database"
	"github.com/trustig/robobaby0.5/internal/discord/logging"
)

var board_channel_id string = os.Getenv("STARBOARD_CHANNEL")
var server_id string = os.Getenv("SERVER_ID")

var image_file_formats []string = []string{
	"gif",
	"png",
	"jpg",
	"jpeg",
	"webm",
	"gifv",
	"webp",
}

const tenor_link string = "https://tenor.com/"

const emoji_name string = "quality5"
const min_reaction_num int = 3

const message_url_prefix string = "https://discordapp.com/channels/"

type StarredMessage struct {
	UserID             string
	ChannelID          string
	StarboardMessageID string
	StarNum            int
}

func OnReact(session *discordgo.Session, reactionAdd *discordgo.MessageReactionAdd) {
	if !(reactionAdd.Emoji.Name == emoji_name) {
		return
	}

	message, err := session.ChannelMessage(reactionAdd.ChannelID, reactionAdd.MessageID)

	if err != nil {
		logging.LogError(session, "getting message for quality5 board", err)
		return
	}

	q5Num := getQuality5Reactions(message)

	if q5Num < min_reaction_num {
		return
	}

	boardMessages := make(map[string]StarredMessage)
	database.LoadJson("db/boardmessages.json", &boardMessages)
	defer database.SaveJson("db/boardmessages.json", boardMessages)

	boardMessage, exists := boardMessages[message.ID]
	if !exists {
		messageId := SendBoardMessage(session, message, &discordgo.MessageEmbedFooter{
			IconURL: "https://cdn.discordapp.com/attachments/1146774215127744533/1149431894912544868/fbac9e72-de98-45d2-9089-4d471f6783be.png",
			Text:    strconv.Itoa(q5Num),
		})

		if messageId == "" {
			return
		}

		boardMessages[message.ID] = StarredMessage{message.Author.ID, message.ChannelID, messageId, q5Num}
	} else {
		editBoardMessage(session, boardMessage.StarboardMessageID, q5Num)

		boardMessage.StarNum = q5Num
		boardMessages[message.ID] = boardMessage
	}
}

func getQuality5Reactions(message *discordgo.Message) int {
	for _, emoji := range message.Reactions {
		if emoji.Emoji.Name == emoji_name {
			return emoji.Count
		}
	}
	return 0
}

func SendBoardMessage(session *discordgo.Session, message *discordgo.Message, footer *discordgo.MessageEmbedFooter) string {
	embed := discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: message.Author.AvatarURL("128"),
			Name:    message.Author.Username,
		},
		Fields:      make([]*discordgo.MessageEmbedField, 0, 4),
		Description: message.Content,
		Color:       rand.Intn(16777215),
		Footer:      footer,
		Timestamp:   message.Timestamp.Format(time.RFC3339),
	}

	embed = addFields(session, message, embed)
	embed = addFiles(session, message, embed)

	message, err := session.ChannelMessageSendComplex(board_channel_id, &discordgo.MessageSend{
		Embed: &embed,
	})

	if err != nil {
		logging.LogError(session, "sending <:quality5:1146794549210001511> board message", err)
		return ""
	}

	return message.ID
}

func addFields(session *discordgo.Session, message *discordgo.Message, embed discordgo.MessageEmbed) discordgo.MessageEmbed {

	if message.ReferencedMessage != nil {
		replyMessage, err := session.ChannelMessage(message.MessageReference.ChannelID, message.MessageReference.MessageID)

		if err != nil {
			logging.LogError(session, "getting <:quality5:1146794549210001511> message reply", err)
			return embed
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "Replying to " + replyMessage.Author.Username,
			Value: replyMessage.Content,
		})
	}

	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Value: "[Message](" + message_url_prefix + server_id + "/" + message.ChannelID + "/" + message.ID + ")",
	})

	return embed
}

func addFiles(session *discordgo.Session, message *discordgo.Message, embed discordgo.MessageEmbed) discordgo.MessageEmbed {
	if len(message.Attachments) > 0 {
		for _, attachment := range message.Attachments {
			isImage, _ := isImageFile(attachment.URL)

			if !isImage {
				continue
			}

			embed.Image = &discordgo.MessageEmbedImage{
				URL: attachment.URL,
			}

			return embed
		}
	}

	//https://stackoverflow.com/questions/1500260/detect-urls-in-text-with-javascript
	regexpCompiler := regexp.MustCompile(`((?:(http|https|Http|Https|rtsp|Rtsp):\/\/(?:(?:[a-zA-Z0-9\$\-\_\.\+\!\*\'\(\)\,\;\?\&\=]|(?:\%[a-fA-F0-9]{2})){1,64}(?:\:(?:[a-zA-Z0-9\$\-\_\.\+\!\*\'\(\)\,\;\?\&\=]|(?:\%[a-fA-F0-9]{2})){1,25})?\@)?)?((?:(?:[a-zA-Z0-9][a-zA-Z0-9\-]{0,64}\.)+(?:(?:aero|arpa|asia|a[cdefgilmnoqrstuwxz])|(?:biz|b[abdefghijmnorstvwyz])|(?:cat|com|coop|c[acdfghiklmnoruvxyz])|d[ejkmoz]|(?:edu|e[cegrstu])|f[ijkmor]|(?:gov|g[abdefghilmnpqrstuwy])|h[kmnrtu]|(?:info|int|i[delmnoqrst])|(?:jobs|j[emop])|k[eghimnrwyz]|l[abcikrstuvy]|(?:mil|mobi|museum|m[acdghklmnopqrstuvwxyz])|(?:name|net|n[acefgilopruz])|(?:org|om)|(?:pro|p[aefghklmnrstwy])|qa|r[eouw]|s[abcdeghijklmnortuvyz]|(?:tel|travel|t[cdfghjklmnoprtvwz])|u[agkmsyz]|v[aceginu]|w[fs]|y[etu]|z[amw]))|(?:(?:25[0-5]|2[0-4][0-9]|[0-1][0-9]{2}|[1-9][0-9]|[1-9])\.(?:25[0-5]|2[0-4][0-9]|[0-1][0-9]{2}|[1-9][0-9]|[1-9]|0)\.(?:25[0-5]|2[0-4][0-9]|[0-1][0-9]{2}|[1-9][0-9]|[1-9]|0)\.(?:25[0-5]|2[0-4][0-9]|[0-1][0-9]{2}|[1-9][0-9]|[0-9])))(?:\:\d{1,5})?)(\/(?:(?:[a-zA-Z0-9\;\/\?\:\@\&\=\#\~\-\.\+\!\*\'\(\)\,\_])|(?:\%[a-fA-F0-9]{2}))*)?(?:\b|$)`)

	links := regexpCompiler.FindAll([]byte(message.Content), -1)

	for _, linkBytes := range links {
		link := string(linkBytes)

		isImage, isTenor := isImageFile(link)
		if !isImage {
			continue
		}

		if isTenor {
			link = getLinkFromTenor(link)
		}

		embed.Image = &discordgo.MessageEmbedImage{
			URL: link,
		}

		return embed
	}

	return embed
}

func isImageFile(link string) (bool, bool) {
	loweredString := strings.ToLower(link)

	if strings.Contains(loweredString, tenor_link) {
		return true, true
	}

	for _, file_format := range image_file_formats {
		if !strings.HasSuffix(loweredString, file_format) {
			continue
		}
		return true, false
	}

	return false, false
}

func getLinkFromTenor(link string) string {
	response, err := http.Get(link)

	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		log.Fatalln(err)
	}

	selection := document.Find(".Gif")
	gifSelection := selection.Find("img")

	return gifSelection.AttrOr("src", "https://cdn.discordapp.com/attachments/1146773305693585420/1149444323952308305/F0W9ZKzaAAMykvf.jpeg")
}

func editBoardMessage(session *discordgo.Session, messageId string, q5Num int) {
	message, err := session.ChannelMessage(board_channel_id, messageId)

	if err != nil {
		logging.LogError(session, "getting :quality5: board message", err)
		return
	}

	embed := message.Embeds[0]
	embed.Footer.Text = strconv.Itoa(q5Num)

	session.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Embed:   embed,
		Channel: board_channel_id,
		ID:      message.ID,
	})
}

func OnUnreact(session *discordgo.Session, reactionRemove *discordgo.MessageReactionRemove) {
	if !(reactionRemove.Emoji.Name == emoji_name) {
		return
	}

	message, err := session.ChannelMessage(reactionRemove.ChannelID, reactionRemove.MessageID)

	if err != nil {
		logging.LogError(session, "getting message for removing <:quality5:1146794549210001511> board", err)
		return
	}

	boardMessages := make(map[string]StarredMessage)
	database.LoadJson("db/boardmessages.json", &boardMessages)
	defer database.SaveJson("db/boardmessages.json", boardMessages)

	if getQuality5Reactions(message) >= min_reaction_num {
		starredMessage := boardMessages[message.ID]
		q5Num := getQuality5Reactions(message)

		editBoardMessage(session, starredMessage.StarboardMessageID, q5Num)

		starredMessage.StarNum = q5Num
		boardMessages[message.ID] = starredMessage
	} else {
		err = session.ChannelMessageDelete(board_channel_id, boardMessages[message.ID].StarboardMessageID)

		if err != nil {
			logging.LogError(session, "deleting <:quality5:1146794549210001511> board message", err)
		}

		delete(boardMessages, message.ID)
	}
}
