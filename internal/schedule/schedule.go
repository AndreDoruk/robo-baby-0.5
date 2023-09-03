package schedule

import (
	"time"

	"github.com/trustig/robobaby0.5/internal/database"
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
	times := make(map[string]time.Time)

	database.LoadJson("db/timed.json", &times)
	defer database.SaveJson("db/timed.json", times)

	lastTime, exists := times[name]

	if !exists {
		times[name] = time.Now()
		return 0 * time.Hour
	}

	return time.Until(lastTime.Add(duration))
}

func updateSleepTime(name string) {
	times := make(map[string]time.Time)

	database.LoadJson("db/timed.json", &times)
	defer database.SaveJson("db/timed.json", times)

	times[name] = time.Now()
}
