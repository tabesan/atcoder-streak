package commit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Getter interface {
	GetCommit(req *http.Request) []Commits
	LastCommitReq() []Commits
	AllCommitReq() []Commits
	InitStreak()
}

func (c *Client) GetLastCommit() []Commits {
	req, err := http.NewRequest("GET", c.URL.String(), nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	q := req.URL.Query()
	q.Set("per_page", "1")
	req.URL.RawQuery = q.Encode()
	resp := c.GetCommit(req)
	return resp
}

func (c *Client) GetAllCommit() []Commits {
	req, err := http.NewRequest("GET", c.URL.String(), nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	resp := c.GetCommit(req)
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
		t = c.edit.ConvJST(v.Commit.Author.Date)
		formatT = t.Format(c.edit.Layout)
		if !mp[formatT] {
			mp[formatT] = true
			days = append(days, formatT)
		}
	}

	c.latestCommit = days[0]
	c.streak = len(days)
	c.ShowStreak()
}
