package whitelist

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/timshannon/bolthold"
)

type Whitelist struct {
	FavorVotes   int
	AgainstVotes int
}

func OnJoin(session *discordgo.Session, event *discordgo.GuildMemberAdd) {
	store, err := bolthold.Open("db/whitelist.db", 0666, nil)

	if err != nil {
		fmt.Println(err)
		session.GuildBanCreate(event.Member.GuildID, event.Member.User.ID, 0)
	}

	defer store.Close()

	var whitelist Whitelist
	err = store.Get(event.User.ID, &whitelist)

	if err != nil {
		fmt.Println(err)
		session.GuildBanCreate(event.Member.GuildID, event.Member.User.ID, 0)
	}
}
