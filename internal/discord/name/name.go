package name

import (
	"log"
	"math/rand"
	"os"

	"github.com/bwmarrin/discordgo"
)

var server_id string = os.Getenv("SERVER_ID")

var server_names = []string{
	"Basement '95",
	"Rotten Tomato",
	"Banned 𝘇𝗮𝗺𝗶𝗲𝗹",
	"Brontulous Orange",
	"im so green j could green a hors",
	"typical toxic isaac community",
	"Poop National Park",
	"Toxicing of Isaac",
	"Better than IsaacScript",
	"I REMEMBERED THE FABLES",
	"the entirety of america",
	"piss",
	"green isaac commits modded murder",
	"fiend and green isaac making out in the basement",
	"sheriff executes your family (real)",
	"Botting of Isaac",
	"THINGS A CREEP MIGHT DO WHEN NEAR CHILDREN 2",
	"oiled up",
}

func ChangeServerName(session *discordgo.Session) {
	selected_name := server_names[rand.Intn(len(server_names))]

	_, err := session.GuildEdit(server_id, &discordgo.GuildParams{
		Name: selected_name,
	})

	if err != nil {
		log.Fatalln(err)
	}
}
