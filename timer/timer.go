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
		//ChFlag:      make(chan string),
		ChUpdate:    make(chan string),
		location:    setLocation(),
		edit:        NewEditTime(),
		updateInter: 30,
		flagInter:   24 * 60 * 60,
	}

	return ct
}

func (c *ChTimer) FlagTimer() {
	var (
		now  time.Time
		wait int
	)

	c.ChFlag <- "ResetFlag"
	if c.edit.ConvJST(time.Now()).Hour() != 0 {
		curTime := c.edit.ConvToSec(time.Now())
		initSleep := c.flagInter - curTime
		time.Sleep(time.Duration(initSleep) * time.Second)
	}
	c.ChFlag <- "ResetFlag"
	for range time.Tick(time.Duration(c.flagInter) * time.Second) {
		now = c.edit.ConvJST(time.Now())
		if now.Hour() != 0 {
			wait = c.flagInter - c.edit.ConvToSec(now)
			time.Sleep(time.Duration(wait) * time.Second)
		}
		c.ChFlag <- "ResetFlag"
	}
}

func (c *ChTimer) UpdateTimer(duration ...time.Duration) {
	now := c.edit.ConvJST(time.Now())
	curTime := (now.Minute()%30)*60 + now.Second()
	interval := time.Minute
	for _, d := range duration {
		interval = d
		curTime = 3
	}
	initSleep := c.updateInter*60 - curTime
	if initSleep < 0 {
		initSleep *= -1
	}
	c.ChUpdate <- "UpdateStreak"
	time.Sleep(time.Duration(initSleep) * time.Second)
	c.ChUpdate <- "UpdateStreak"
	for range time.Tick(time.Duration(c.updateInter) * interval) {
		c.ChUpdate <- "UpdateStreak"
	}
}

func (c *ChTimer) Timer() {
	const interval = 12 * 60 * 60
	now := c.edit.ConvToSec(time.Now())
	time.Sleep(time.Duration(interval-now) * time.Second)
	c.ChUpdate <- "Update"
	for range time.Tick(time.Hour * 24) {
		c.ChUpdate <- "Update"
	}
}
