package commit

import (
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
	edit          *editTime
	getter        Getter
}

func NewClient(name, repo string) *Client {
	c := &Client{
		name:          name,
		repo:          repo,
		latestCommit:  "",
		streak:        0,
		LongestStreak: 0,
		updateFlag:    false,
		edit:          NewEditTime(),
		getter:        NewComGetter(),
	}
	c.createURL()
	return c
}

func (c *Client) createURL() {
	c.URL = &url.URL{
		Scheme: "https",
		Host:   "api.github.com",
		Path:   "repos/" + c.name + "/" + c.repo + "/commits",
	}
}

func (c *Client) LastCommitReq() *http.Request {
	req, err := http.NewRequest("GET", c.URL.String(), nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.URL.Query().Set("per_page", "1")

	return req
}

func (c *Client) AllCommitReq() *http.Request {
	req, err := http.NewRequest("GET", c.URL.String(), nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	return req
}

func (c *Client) InitStreak() {
	commits := c.getter.GetCommit(c.AllCommitReq())
	mp := make(map[string]bool)
	var days []string
	var t time.Time
	var formatT string

	for _, v := range commits {
		t = c.edit.convJST(v.Commit.Author.Date)
		formatT = t.Format(layout)
		if !mp[formatT] {
			mp[formatT] = true
			days = append(days, formatT)
		}
	}

	c.latestCommit = days[0]
	c.streak = len(days)
	c.ShowStreak()
}

func (c *Client) update() {
	latest := (c.getter.GetCommit(c.LastCommitReq()))[0]
	lastDate := c.edit.convJST(latest.Commit.Author.Date)
	DayAgo := (lastDate.AddDate(0, 0, -1)).Format(layout)
	if c.latestCommit == DayAgo {
		c.streak += 1
		c.latestCommit = lastDate.Format(layout)
	} else if c.latestCommit != lastDate.Format(layout) {
		c.streak = 0
	}
}

func (c *Client) ShowStreak() {
	fmt.Println(c.streak)
}

func (c *Client) FlagReset() {
	c.updateFlag = false
}

func (c *Client) UpdateStreak() {
	if !c.updateFlag {
		c.update()
		c.updateFlag = true
	}
}
