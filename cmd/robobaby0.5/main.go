package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/trustig/robobaby0.5/internal/discord/items"
	"github.com/trustig/robobaby0.5/internal/discord/name"
	"github.com/trustig/robobaby0.5/internal/discord/slash"
	"github.com/trustig/robobaby0.5/internal/discord/voting"
	"github.com/trustig/robobaby0.5/internal/discord/whitelist"
	"github.com/trustig/robobaby0.5/internal/schedule"
	"github.com/trustig/robobaby0.5/internal/workshop"
)

const update_items_frequency time.Duration = 12 * time.Hour
const server_name_frequency time.Duration = 24 * time.Hour

func main() {
	session, err := discordgo.New("Bot " + os.Getenv("TOKEN"))

	if err != nil {
		log.Fatalln(err)
	}

	session.Identify.Intents = discordgo.IntentsAll

	session.AddHandler(whitelist.OnJoin)
	session.AddHandler(slash.OnInteract)

	session.Open()
	defer session.Close()

	if err != nil {
		log.Fatalln(err)
	}

	slash.CreateCommands(session)
	go voting.Loop(session)

	schedule.Loop("workshopItems", update_items_frequency, func() {
		items.SendWorkshopItems(session, workshop.GetMostPopularItems())
	})

	schedule.Loop("servername", server_name_frequency, func() {
		name.ChangeServerName(session)
	})

	err = session.UpdateGameStatus(10, "The Binding of Isaac: Antibirth")

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Bot online")
	for {
	}
}
