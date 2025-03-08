package voting

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/AndreDoruk/robo-baby-0.5/internal/database"
	"github.com/AndreDoruk/robo-baby-0.5/internal/discord/upload"
	"github.com/AndreDoruk/robo-baby-0.5/internal/discord/whitelist"
	"github.com/AndreDoruk/robo-baby-0.5/internal/images"
)

var channel_id string = os.Getenv("VOTING_CHANNEL_ID")

const sleep_time time.Duration = 1 * time.Minute

const vote_duration time.Duration = 12 * time.Hour

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

	votes := make(map[string]Vote)
	database.LoadJson("db/votes.json", &votes)

	defer database.SaveJson("db/votes.json", votes)

	if err != nil {
		log.Fatalln(err)
	}

	votes[userId] = vote

	return nil
}

func UpdateVoting(session *discordgo.Session) {
	votes := make(map[string]Vote)

	database.LoadJson("db/votes.json", &votes)
	defer database.SaveJson("db/votes.json", votes)

	for _, vote := range votes {
		votes = updateVote(session, vote, votes)
	}
}

func updateVote(session *discordgo.Session, vote Vote, votes map[string]Vote) map[string]Vote {
	message, err := session.ChannelMessage(channel_id, vote.MessageId)

	if err != nil {
		log.Fatalln(err)
	}

	remainingTime := time.Until(vote.TimeStarted.Add(vote_duration))
	overwhelmingDifference := overwhelmingDifferenceInVotes(message)

	if remainingTime.Minutes() > 0 && !overwhelmingDifference {
		hours := math.Floor(remainingTime.Hours())

		if !(hours > 1 && vote.LastHour == hours) {
			editImageTimestamp(session, message, math.Floor(remainingTime.Minutes()))
			vote.LastHour = hours
		}

		votes[vote.UserId] = vote
	} else {
		finishVote(session, message, vote)

		delete(votes, vote.UserId)
	}

	return votes
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

func finishVote(session *discordgo.Session, message *discordgo.Message, vote Vote) {
	whitelists := make(map[string]whitelist.Whitelist)
	database.LoadJson("db/whitelist.json", &whitelists)
	defer database.SaveJson("db/whitelist.json", whitelists)

	votesFavor, votesAgainst := getFavorAndAgainstVotes(message)
	ratio := float32(votesFavor) / float32(votesAgainst)

	wentThrough := (votesFavor+votesAgainst) > min_votes && ratio >= min_ratio_go_through

	if wentThrough {
		whitelist := whitelist.Whitelist{
			FavorVotes:   votesFavor,
			AgainstVotes: votesAgainst,
			UserId:       vote.UserId,
		}
		whitelists[vote.UserId] = whitelist
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
