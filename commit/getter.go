package commit

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Getter interface {
	GetCommit(ctx context.Context, req *http.Request) []Commits
	GetLastCommit(ctx context.Context) ([]Commits, bool)
	GetAllCommit(ctx context.Context) ([]Commits, bool)
	InitStreak()
}

func (c *Client) GetLastCommit(ctx context.Context) ([]Commits, bool) {
	req, err := http.NewRequest("GET", c.URL.String(), nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	q := req.URL.Query()
	q.Set("per_page", "1")
	req.URL.RawQuery = q.Encode()
	resp := c.GetCommit(ctx, req)
	return resp, true
}

func (c *Client) GetAllCommit(ctx context.Context) ([]Commits, bool) {
	req, err := http.NewRequest("GET", c.URL.String(), nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	resp := c.GetCommit(ctx, req)
	return resp, true
}

func (c *Client) GetCommit(ctx context.Context, req *http.Request) []Commits {
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

func (c *Client) isStreak(target time.Time, pre string) bool {
	if pre == "-1" {
		return true
	}

	dayLater := target.AddDate(0, 0, 1).Format(c.edit.Layout)
	if dayLater == pre {
		return true
	} else {
		return false
	}
}

func (c *Client) InitStreak() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var commits []Commits
	ok := false
	for {
		commits, ok = c.GetAllCommit(ctx)
		if ok {
			break
		}

		select {
		case <-ctx.Done():
			c.timeoutFlag = true
			time.Sleep(10 * time.Minute)
			continue
		}
	}

	mp := make(map[string]bool)
	var days []string
	var target time.Time
	var formatT string
	pre := "-1"
	c.latestCommit = c.edit.ConvJST((commits[0].Commit.Author.Date)).Format(c.edit.Layout)
	for _, v := range commits {
		target = c.edit.ConvJST(v.Commit.Author.Date)
		formatT = target.Format(c.edit.Layout)
		if !mp[formatT] {
			mp[formatT] = true
			if !c.isStreak(target, pre) {
				days = nil
			}
			pre = formatT
			days = append(days, formatT)
		}
	}

	c.streak = len(days)
	c.ShowStreak()
}
