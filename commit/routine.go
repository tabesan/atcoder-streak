package commit

import (
	"time"
)

func (c *Client) convJST(t time.Time) time.Time {
	return t.In(c.location)
}

func (c *Client) convToSec(t time.Time) int {
	jst := c.convJST(t)
	toSec := jst.Hour()*3600 + jst.Minute()*60 + jst.Second()
	return toSec
}

func (c *Client) UpdateRoutine() {
	const interval = 24 * 60 * 60
	var (
		now  time.Time
		wait int
	)

	curTime := c.convToSec(time.Now())
	initSleep := interval - curTime

	time.Sleep(time.Duration(initSleep) * time.Second)
	for {
		now = c.convJST(time.Now())
		if now.Hour() == 0 {
			c.updateStreak()
		} else {
			wait = interval - c.convToSec(now)
			time.Sleep(time.Duration(wait) * time.Second)
			c.updateStreak()
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
