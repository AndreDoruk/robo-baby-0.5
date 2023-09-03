package upload

import (
	"bytes"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

var image_channel_id string = os.Getenv("IMAGE_CHANNEL_ID")

func UploadFileAndReturnUrl(session *discordgo.Session, filename string, reader *bytes.Reader) string {
	message, err := session.ChannelFileSend(image_channel_id, filename, reader)

	if err != nil {
		log.Fatalln(err)
	}

	return message.Attachments[0].URL
}
