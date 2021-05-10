package commit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		edit: NewEditTime(),
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

func (c *Client) GetLastCommit() []Commits {
	req, err := http.NewRequest("GET", c.URL.String(), nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.URL.Query().Set("per_page", "1")

	var resp []Commits
	resp = c.GetCommit(req)
	return resp
}

func (c *Client) GetAllCommit() []Commits {
	req, err := http.NewRequest("GET", c.URL.String(), nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	var resp []Commits
	resp = c.GetCommit(req)
	return resp
}

func (c *Client) GetCommit(req *http.Request) []Commits {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	bytes := []byte(body)
	var data []Commits
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		fmt.Println(err)
	}
	return data
}

func (c *Client) InitStreak() {
	commits := c.GetAllCommit()
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
}

func (c *Client) update() {
	latest := (c.GetLastCommit())[0]
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
