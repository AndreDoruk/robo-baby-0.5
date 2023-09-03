package items

import (
	"bytes"
	"log"
	"os"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/trustig/robobaby0.5/internal/discord/upload"
	"github.com/trustig/robobaby0.5/internal/images"
	"github.com/trustig/robobaby0.5/internal/workshop"
)

var channel_id string = os.Getenv("WORKSHOP_CHANNEL_ID")

func SendWorkshopItems(session *discordgo.Session, items []workshop.WorkshopItem) {
	readers := make([]*bytes.Reader, len(items))

	var waitGroup sync.WaitGroup
	i := 0

	for _, item := range items {
		waitGroup.Add(1)

		go func(index int, item workshop.WorkshopItem) {
			readers[index] = images.ImageToBytesReader(images.CreateWorkshopImage(item))
			waitGroup.Done()
		}(i, item)

		i += 1
	}

	waitGroup.Wait()

	for _, reader := range readers {
		url := upload.UploadFileAndReturnUrl(session, "workshopItem.png", reader)
		_, err := session.ChannelMessageSend(channel_id, url)

		if err != nil {
			log.Fatalln(err)
		}
	}
}

func SendWorkshopItem(session *discordgo.Session, channelId string, item workshop.WorkshopItem) {
	reader := images.ImageToBytesReader(images.CreateWorkshopImage(item))
	url := upload.UploadFileAndReturnUrl(session, "workshopItem.png", reader)

	_, err := session.ChannelMessageSend(channelId, url)

	if err != nil {
		log.Fatalln(err)
	}
}
