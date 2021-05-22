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
	timeoutFlag   bool
	HTTPClient    *http.Client
	URL           *url.URL
	edit          *tm.EditTime
	Getter        Getter
	rockFlag      bool
}

func (c *Client) createURL() {
	c.URL = &url.URL{
		Scheme: "https",
		Host:   "api.github.com",
		Path:   "repos/" + name + "/" + repo + "/commits",
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
		timeoutFlag:   false,
		HTTPClient: &http.Client{
			Timeout: time.Second * 15,
		},
		edit: tm.NewEditTime(),
	}
	c.createURL()
	return c
}

func (c *Client) InitStreak() {
	c.Getter.InitStreak()
}

func (c *Client) Timeouted() {
	if c.updateFlag == false {
		c.timeoutFlag = true
	}
}

func (c *Client) ShowStreak() {
	fmt.Println(c.streak)
}

func (c *Client) ResetFlag() {
	if c.updateFlag == false && c.timeoutFlag == true {
		c.rockFlag = true
		c.Getter.InitStreak()
		c.timeoutFlag = false
	}
	c.rockFlag = false
	c.updateFlag = false
}

func (c *Client) UpdateStreak(ctx context.Context) {
	if !c.updateFlag {
		c.Getter.Update(ctx)
		c.updateFlag = true
		c.timeoutFlag = false
	}
}
