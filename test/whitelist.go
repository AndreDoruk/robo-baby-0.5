package main

import (
	"fmt"
	"log"

	"github.com/timshannon/bolthold"
	"github.com/trustig/robobaby0.5/internal/discord/whitelist"
)

func main() {
	store, err := bolthold.Open("db/whitelist.db", 0666, nil)

	if err != nil {
		log.Fatalln(err)
	}

	store.ForEach(bolthold.Where("Exists").Eq(true), func(whitelist *whitelist.Whitelist) error {
		fmt.Print(whitelist)
		return nil
	})
}
