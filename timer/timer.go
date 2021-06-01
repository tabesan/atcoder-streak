package commit

import (
	"time"
)

const Layout = "2006-01-02"

type ChTimer struct {
	ChFlag      chan string
	ChUpdate    chan string
	location    *time.Location
	edit        *EditTime
	updateInter int
	flagInter   int
}

func NewTimer() *ChTimer {
	ct := &ChTimer{
		ChFlag:      make(chan string),
		ChUpdate:    make(chan string),
		location:    setLocation(),
		edit:        NewEditTime(),
		updateInter: 60 * 60,
		flagInter:   24 * 60 * 60,
	}

	return ct
}

func (c *ChTimer) FlagTimer() {
	var (
		now  time.Time
		wait int
	)

	curTime := c.edit.ConvToSec(time.Now())
	initSleep := c.flagInter - curTime

	time.Sleep(time.Duration(initSleep) * time.Second)
	for range time.Tick(time.Duration(c.flagInter) * time.Second) {
		now = c.edit.ConvJST(time.Now())
		if now.Hour() != 0 {
			wait = c.flagInter - c.edit.ConvToSec(now)
			time.Sleep(time.Duration(wait) * time.Second)
		}
		c.ChFlag <- "UpdateFlag"
	}
}

func (c *ChTimer) UpdateTimer(duration ...time.Duration) {
	now := time.Now()
	curTime := now.Minute()*60 + now.Second()
	interval := time.Hour
	for _, d := range duration {
		interval = d
		curTime = 3
	}
	initSleep := c.updateInter - curTime
	time.Sleep(time.Duration(initSleep) * interval)
	for range time.Tick(time.Duration(c.updateInter) * interval) {
		c.ChUpdate <- "UpdateStreak"
	}
}
