package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"bufio"

	"github.com/bwmarrin/discordgo"

	"github.com/AndreDoruk/robo-baby-0.5/internal/discord/commentgame"
	"github.com/AndreDoruk/robo-baby-0.5/internal/discord/items"
	"github.com/AndreDoruk/robo-baby-0.5/internal/discord/name"
	"github.com/AndreDoruk/robo-baby-0.5/internal/discord/slash"
	"github.com/AndreDoruk/robo-baby-0.5/internal/discord/splatting"
	"github.com/AndreDoruk/robo-baby-0.5/internal/discord/starboard"
	"github.com/AndreDoruk/robo-baby-0.5/internal/discord/voting"
	"github.com/AndreDoruk/robo-baby-0.5/internal/discord/whitelist"
	"github.com/AndreDoruk/robo-baby-0.5/internal/schedule"
	"github.com/AndreDoruk/robo-baby-0.5/internal/workshop"
)

var isTesting bool = os.Getenv("TESTING") != ""

const update_items_frequency time.Duration = 12 * time.Hour
const server_name_frequency time.Duration = 24 * time.Hour
const vote_update_frequency time.Duration = 2 * time.Minute
const splatting_role_frequency time.Duration = 2 * time.Minute

func main() {
	session, err := discordgo.New("Bot " + os.Getenv("TOKEN"))

	if err != nil {
		log.Fatalln(err)
	}

	session.Identify.Intents = discordgo.IntentsAll

	if false {
		session.AddHandler(whitelist.OnJoin)
		session.AddHandler(commentgame.OnInteract)
		session.AddHandler(splatting.OnReact)
		session.AddHandler(starboard.OnUnreact)
		session.AddHandler(starboard.OnReact)
	}

	session.AddHandler(slash.OnInteract)

	session.Open()
	defer session.Close()

	if err != nil {
		log.Fatalln(err)
	}

	slash.CreateCommands(session)

	if !isTesting {
		schedule.Loop("workshopItems", update_items_frequency, func() {
			items.SendWorkshopItems(session, workshop.GetMostPopularItems())
		})

		schedule.Loop("serverName", server_name_frequency, func() {
			name.ChangeServerName(session)
		})

		schedule.Loop("voteUpdate", vote_update_frequency, func() {
			voting.UpdateVoting(session)
		})

		schedule.Loop("splattingRole", splatting_role_frequency, func() {
			splatting.UpdateSplattedRole(session)
		})
	}

	err = session.UpdateGameStatus(10, "The Binding of Isaac: Antibirth")

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Bot online")

	input := bufio.NewScanner(os.Stdin)
    input.Scan()
	session.Close()
}