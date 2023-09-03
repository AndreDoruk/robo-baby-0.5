package voting

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/timshannon/bolthold"
	"github.com/trustig/robobaby0.5/internal/discord/upload"
	"github.com/trustig/robobaby0.5/internal/discord/whitelist"
	"github.com/trustig/robobaby0.5/internal/images"
)

var channel_id string = os.Getenv("VOTING_CHANNEL_ID")

const vote_duration time.Duration = 12 * time.Hour
const for_each_compare float64 = -10000000

const min_ratio_go_through float32 = 70 / 30
const overwhelming_difference_ratio float32 = 90 / 10
const min_overwhelming_difference_voters int = 9
const min_votes int = 5

type Vote struct {
	MessageId   string
	UserId      string
	TimeStarted time.Time
	LastHour    float64
}

func CreateVote(session *discordgo.Session, userId string) error {
	//TODO: Stop from creating vote for whitelisted user
	user, err := session.User(userId)

	if err != nil {
		fmt.Println(err)
		return err
	}

	image := images.CreateVoteImage(session, user)
	path := upload.UploadFileAndReturnUrl(session, "vote.png", images.ImageToBytesReader(image))

	message, err := session.ChannelMessageSend(channel_id, path)

	session.MessageReactionAdd(channel_id, message.ID, "🍏")
	session.MessageReactionAdd(channel_id, message.ID, "🍅")

	if err != nil {
		log.Fatalln(err)
	}

	vote := Vote{message.ID, userId, time.Now(), -1}

	store, err := bolthold.Open("db/votes.db", 0666, nil)

	if err != nil {
		log.Fatalln(err)
	}

	store.Insert(userId, vote)
	store.Close()

	return nil
}

func UpdateVoting(session *discordgo.Session) {
	store, err := bolthold.Open("db/votes.db", 0666, nil)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(store.Count(Vote{}, bolthold.Where("LastHour").Gt(for_each_compare)))

	store.ForEach(bolthold.Where("LastHour").Gt(for_each_compare), func(vote *Vote) error {
		go updateVote(session, vote, *store)
		return nil
	})
}

func updateVote(session *discordgo.Session, vote *Vote, store bolthold.Store) {
	message, err := session.ChannelMessage(channel_id, vote.MessageId)

	if err != nil {
		log.Fatalln(err)
	}

	remainingTime := time.Until(vote.TimeStarted.Add(vote_duration))
	overwhelmingDifference := overwhelmingDifferenceInVotes(message)

	fmt.Println("C!")

	if remainingTime.Minutes() > 0 && !overwhelmingDifference {
		hours := math.Floor(remainingTime.Hours())

		if !(hours > 1 && vote.LastHour == hours) {
			go editImageTimestamp(session, message, math.Floor(remainingTime.Minutes()))
			vote.LastHour = hours
		}

		store.Update(vote.UserId, vote)
	} else {
		go finishVote(session, message, vote)

		err = store.Delete(vote.UserId, Vote{})

		if err != nil {
			log.Fatalln(err)
		}
	}
}

func overwhelmingDifferenceInVotes(message *discordgo.Message) bool {
	votesFavor, votesAgainst := getFavorAndAgainstVotes(message)

	if votesFavor+votesAgainst < min_overwhelming_difference_voters {
		return false
	}

	ratio := float32(votesFavor) / float32(votesAgainst)

	return ratio <= 1/overwhelming_difference_ratio || ratio >= overwhelming_difference_ratio
}

func editImageTimestamp(session *discordgo.Session, message *discordgo.Message, minutes float64) {
	image := images.GetImageFromUrl(message.Content)
	newImage := images.UpdateVoteTimestamp(image, minutes)

	reader := images.ImageToBytesReader(newImage)

	url := upload.UploadFileAndReturnUrl(session, "b.png", reader)
	session.ChannelMessageEdit(message.ChannelID, message.ID, url)
}

func finishVote(session *discordgo.Session, message *discordgo.Message, vote *Vote) {
	store, err := bolthold.Open("db/whitelist.db", 0666, nil)

	if err != nil {
		log.Fatalln(err)
	}

	defer store.Close()

	votesFavor, votesAgainst := getFavorAndAgainstVotes(message)
	ratio := float32(votesFavor) / float32(votesAgainst)

	wentThrough := (votesFavor+votesAgainst) > min_votes && ratio >= min_ratio_go_through

	if wentThrough {
		whitelist := whitelist.Whitelist{
			FavorVotes:   votesFavor,
			AgainstVotes: votesAgainst,
			UserId:       vote.UserId,
		}

		err = store.Insert(vote.UserId, whitelist)

		if err != nil {
			log.Fatalln(err)
		}
	}

	go updateVoteVictory(session, message, wentThrough)
}

func updateVoteVictory(session *discordgo.Session, message *discordgo.Message, win bool) {
	image := images.GetImageFromUrl(message.Content)
	newImage := images.UpdateVoteVictoryText(image, win)

	reader := images.ImageToBytesReader(newImage)

	url := upload.UploadFileAndReturnUrl(session, "vote.png", reader)
	_, err := session.ChannelMessageEdit(message.ChannelID, message.ID, url)

	if err != nil {
		log.Fatalln(err)
	}
}

func getFavorAndAgainstVotes(message *discordgo.Message) (int, int) {
	votesFavor, votesAgainst := 0, 0

	for _, reaction := range message.Reactions {
		if reaction.Emoji.Name == "🍏" {
			votesFavor = reaction.Count
		} else if reaction.Emoji.Name == "🍅" {
			votesAgainst = reaction.Count
		}
	}

	return votesFavor, votesAgainst
}
