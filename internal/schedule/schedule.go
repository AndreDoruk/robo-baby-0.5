package schedule

import (
	"fmt"
	"log"
	"time"

	"github.com/timshannon/bolthold"
)

type validFunc func()

func Loop(name string, duration time.Duration, function validFunc) {
	sleepTime := getSleepTime(name, duration)

	time.AfterFunc(sleepTime, func() {
		function()
		updateSleepTime(name)
		Loop(name, duration, function)
	})
}

func getSleepTime(name string, duration time.Duration) time.Duration {
	store, err := bolthold.Open("db/timed.db", 0666, nil)

	if err != nil {
		log.Fatalln(err)
	}

	var lastTime time.Time
	err = store.Get(name, &lastTime)

	if err != nil {
		fmt.Println(err)

		store.Insert(name, time.Now())
		store.Close()

		return 0 * time.Hour
	}

	store.Close()
	return time.Until(lastTime.Add(duration))
}

func updateSleepTime(name string) {
	store, err := bolthold.Open("db/timed.db", 0666, nil)

	if err != nil {
		log.Fatalln(err)
	}

	store.Update(name, time.Now())
	store.Close()
}
