package commit

import (
	tm "atcoder-streak/timer"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	name          string
	repo          string
	latestCommit  string
	streak        int
	LongestStreak int
	updateFlag    bool
	URL           *url.URL
	HTTPClient    *http.Client
	edit          *tm.EditTime
	timeoutFlag   bool
}

func (c *Client) createURL() {
	c.URL = &url.URL{
		Scheme: "https",
		Host:   "api.github.com",
		Path:   "repos/" + c.name + "/" + c.repo + "/commits",
	}

}

func NewClient(name, repo string) *Client {
	c := &Client{
		name:          name,
		repo:          repo,
		latestCommit:  "",
		streak:        0,
		LongestStreak: 0,
		updateFlag:    false,
		HTTPClient: &http.Client{
			Timeout: time.Second * 15,
		},
		edit:        tm.NewEditTime(),
		timeoutFlag: false,
	}
	c.createURL()
	return c
}

func (c *Client) Timeouted() {
	if c.updateFlag == false {
		c.timeoutFlag = true
	}
}

func (c *Client) update(ctx context.Context) {
	resp, ok := c.GetLastCommit(ctx)
	fmt.Println(ok)
	latest := resp[0]
	lastDate := c.edit.ConvJST(latest.Commit.Author.Date)
	DayAgo := (lastDate.AddDate(0, 0, -1)).Format(c.edit.Layout)
	if c.latestCommit == DayAgo {
		c.streak += 1
		c.latestCommit = lastDate.Format(c.edit.Layout)
	} else if c.latestCommit != lastDate.Format(c.edit.Layout) {
		c.streak = 0
	}
}

func (c *Client) ShowStreak() {
	fmt.Println(c.streak)
}

func (c *Client) FlagReset() {
	if c.updateFlag == false && c.timeoutFlag == true {
		c.InitStreak()
		c.timeoutFlag = false
	}
	c.updateFlag = false
}

func (c *Client) UpdateStreak(ctx context.Context) {
	if !c.updateFlag {
		c.update(ctx)
		time.Sleep(20 * time.Second)
		c.updateFlag = true
		c.timeoutFlag = false
	}
}
