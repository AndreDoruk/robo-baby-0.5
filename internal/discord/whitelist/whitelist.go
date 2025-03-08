package whitelist

import (
	"github.com/bwmarrin/discordgo"
	"github.com/AndreDoruk/robo-baby-0.5/internal/database"
)

type Whitelist struct {
	FavorVotes   int
	AgainstVotes int
	UserId       string
}

func OnJoin(session *discordgo.Session, event *discordgo.GuildMemberAdd) {
	whitelist := make(map[string]Whitelist)
	database.LoadJson("db/whitelist.json", &whitelist)

	if _, exists := whitelist[event.User.ID]; !exists {
		session.GuildBanCreate(event.Member.GuildID, event.Member.User.ID, 0)
	}
}
