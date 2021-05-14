package commit

import (
	"time"
)

const layout = "2006-01-02"

type ChTimer struct {
	ChFlag   chan string
	ChUpdate chan string
	location *time.Location
	edit     *editTime
}

func NewTimer() *ChTimer {
	ct := &ChTimer{
		ChFlag:   make(chan string),
		ChUpdate: make(chan string),
		location: setLocation(),
		edit:     NewEditTime(),
	}

	return ct
}

func (c *ChTimer) FlagTimer() {
	const interval = 24 * 60 * 60
	var (
		now  time.Time
		wait int
	)

	curTime := c.edit.convToSec(time.Now())
	initSleep := interval - curTime

	time.Sleep(time.Duration(initSleep) * time.Second)
	for range time.Tick(time.Duration(interval) * time.Second) {
		now = c.edit.convJST(time.Now())
		if now.Hour() != 0 {
			wait = interval - c.edit.convToSec(now)
			time.Sleep(time.Duration(wait) * time.Second)
		}
		c.ChFlag <- "UpdateFlag"
	}
}

func (c *ChTimer) UpdateTimer() {
	const interval = 1

	curTime := c.edit.convToSec(time.Now())
	initSleep := interval - curTime

	time.Sleep(time.Duration(initSleep) * time.Second)
	for range time.Tick(interval * time.Hour) {
		c.ChUpdate <- "UpdateStreak"
	}
}
