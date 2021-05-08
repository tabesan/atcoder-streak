package commit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const layout = "2006-01-02"

type Client struct {
	name          string
	repo          string
	latestCommit  string
	streak        int
	LongestStreak int
	URL           *url.URL
	HTTPClient    *http.Client
}

func NewClient(name, repo string) *Client {
	c := &Client{
		name: name,
		repo: repo,
		HTTPClient: &http.Client{
			Timeout: time.Second * 15,
		},
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
	var t string
	for _, v := range commits {
		t = (v.Commit.Author.Date).Format(layout)
		if !mp[t] {
			mp[t] = true
			days = append(days, t)
		}
	}

	c.latestCommit = days[0]
	c.streak = len(days)
}

func (c *Client) UpdateStreak() {
	latest := (c.GetLastCommit())[0]
	lastDate := latest.Commit.Author.Date
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
