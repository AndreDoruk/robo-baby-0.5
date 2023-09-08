package commentgame

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/trustig/robobaby0.5/internal/database"
	"github.com/trustig/robobaby0.5/internal/workshop"
)

const item_num int = 4

type CommentGame struct {
	PlayerID      string
	WorkshopItems []workshop.WorkshopItem
	CorrectItem   workshop.WorkshopItem
	MessageID     string
}

const seconds int = 15

const comment_game_play_again_id = "PLAY_AGAIN_COMMENT_GAME"

var currentGames map[string]CommentGame = make(map[string]CommentGame)

var loss_text = []string{
	"do you just click a random button",
	"how are you this dumb bro",
	"they should kill you with hammers",
	"i've seen workshop users smarter than you",
	"Does He Know",
}

var win_text = []string{
	"intelligence unrivaled in the modding community",
	"Genius!",
	"congrat!",
	"so good",
	"How does he do it",
}

func OnInteract(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Type != discordgo.InteractionMessageComponent {
		return
	}

	game, isCommentGame := currentGames[interaction.Message.ID]

	selectedOption := interaction.Interaction.Data.(discordgo.MessageComponentInteractionData).CustomID

	if selectedOption == comment_game_play_again_id {
		Begin(session, interaction)

		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponsePong,
		})
	}

	if !isCommentGame {
		return
	}

	if interaction.Member.User.ID != game.PlayerID {
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content:         "this isn't your battle to fight....",
				AllowedMentions: &discordgo.MessageAllowedMentions{Users: []string{interaction.Member.User.ID}, RepliedUser: true},
				Flags:           discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})

	endGame(session, interaction.Message, selectedOption == game.CorrectItem.Name, selectedOption)
}

func Begin(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	workshopItems, comment, correctItem := getGameItemsCommentAndCorrect()

	var components []discordgo.MessageComponent

	for _, item := range workshopItems {
		components = append(components, discordgo.Button{
			Label:    item.Name,
			Style:    discordgo.SecondaryButton,
			Disabled: false,
			CustomID: item.Name,
		})
	}

	message, err := session.ChannelMessageSendComplex(interaction.ChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Author:      &discordgo.MessageEmbedAuthor{IconURL: comment.IconURL, Name: comment.Creator + " says: "},
			Description: comment.Comment,
			Color:       rand.Intn(16777215),
			Footer:      &discordgo.MessageEmbedFooter{IconURL: interaction.Member.AvatarURL("128"), Text: interaction.Member.User.Username + " is playing [" + strconv.Itoa(seconds) + " second(s) remaining]"},
		},

		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{Components: components},
		},
	})

	if err != nil {
		session.ChannelMessageSend(interaction.ChannelID, "```go\n"+err.Error()+"```")
		return
	}

	currentGame := CommentGame{
		PlayerID:      interaction.Member.User.ID,
		WorkshopItems: workshopItems,
		CorrectItem:   correctItem,
		MessageID:     message.ID,
	}
	currentGames[message.ID] = currentGame

	go loopTimer(session, message, seconds)
}

func loopTimer(session *discordgo.Session, message *discordgo.Message, seconds int) {
	for seconds >= 0 {
		_, exists := currentGames[message.ID]
		if !exists {
			return
		}

		embed := message.Embeds[0]

		prefix := strings.Split(embed.Footer.Text, "[")[0]
		footerText := prefix + "[" + strconv.Itoa(seconds) + " second(s) remaining]"

		embed.Footer.Text = footerText

		message, _ = session.ChannelMessageEditEmbed(message.ChannelID, message.ID, embed)

		time.Sleep(time.Second)
		seconds -= 1
	}

	endGame(session, message, false, "")
}

func endGame(session *discordgo.Session, message *discordgo.Message, victory bool, selectedButton string) {
	game, exists := currentGames[message.ID]

	if !exists {
		return
	}

	originalEmbed := message.Embeds[0]

	prefix := strings.Split(originalEmbed.Footer.Text, "[")[0]
	originalEmbed.Footer.Text = prefix + "[Over!]"

	actionsRow := message.Components[0].(*discordgo.ActionsRow)

	for _, component := range actionsRow.Components {
		button := component.(*discordgo.Button)
		button.Disabled = true

		if button.Label == selectedButton {
			button.Style = discordgo.PrimaryButton
		}
	}

	go func() {
		time.Sleep(time.Second)

		_, err := session.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Embed:      originalEmbed,
			Components: message.Components,
			ID:         message.ID,
			Channel:    message.ChannelID,
		})

		if err != nil {
			log.Fatalln(err)
		}
	}()

	delete(currentGames, message.ID)

	correctItem := game.CorrectItem

	embed := &discordgo.MessageEmbed{
		Fields: []*discordgo.MessageEmbedField{{
			Name:  "Answer",
			Value: "[" + correctItem.Name + "](" + correctItem.URL + ")",
		}},
		Color:  originalEmbed.Color,
		Footer: &discordgo.MessageEmbedFooter{IconURL: originalEmbed.Footer.IconURL, Text: "wow"},
	}

	if victory {
		embed.Author = &discordgo.MessageEmbedAuthor{IconURL: "https://cdn.discordapp.com/attachments/1133685184240287806/1148231464824090774/epico-mandela-catalog.gif", Name: win_text[rand.Intn(len(win_text))]}
		embed.Description = `you win, good job
		🍅
		🫴`
		embed.Footer.Text = "+1 tomato"

		addTomato(game.PlayerID)
	} else {
		embed.Author = &discordgo.MessageEmbedAuthor{IconURL: "https://media.discordapp.net/attachments/954489163820978196/1002423394236637204/edge.gif", Name: loss_text[rand.Intn(len(loss_text))]}
		embed.Description = "you lose lol, loser"
		embed.Footer.Text = "+0 tomatoes"
	}

	session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
		Embed: embed,
		Components: []discordgo.MessageComponent{discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{discordgo.Button{
				Label:    "Play Again",
				Style:    discordgo.PrimaryButton,
				Disabled: false,
				CustomID: comment_game_play_again_id,
			}},
		}},
		Reference: message.Reference(),
	})
}

func getGameItemsCommentAndCorrect() ([]workshop.WorkshopItem, workshop.WorkshopComment, workshop.WorkshopItem) {
	items := make([]workshop.WorkshopItem, 0, item_num)

	var waitGroup sync.WaitGroup

	for i := 0; i < item_num-1; i++ {
		waitGroup.Add(1)

		go func(i int) {
			items = append(items, workshop.GetRandomItem())
			waitGroup.Done()
		}(i)
	}

	var comment workshop.WorkshopComment
	var correctItem workshop.WorkshopItem

	waitGroup.Add(1)
	go func() {
		comment, correctItem = workshop.GetRandomCommentAndItem()
		waitGroup.Done()
	}()

	waitGroup.Wait()

	correctIndex := rand.Intn(item_num)

	items = append(items[:correctIndex+1], items[correctIndex:]...)
	items[correctIndex] = correctItem

	return items, comment, correctItem
}

func addTomato(userId string) {
	tomatoes := make(map[string]int)

	database.LoadJson("db/balance.json", &tomatoes)
	defer database.SaveJson("db/balance.json", tomatoes)

	tomatoes[userId] += 1
}
