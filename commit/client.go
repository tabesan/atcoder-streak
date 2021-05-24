package commit

import (
	tm "atcoder-streak/timer"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
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
		c.InitStreak()
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

func (c *Client) SetStreak(commits []Commits) error {
	mp := make(map[string]bool)
	var days []string
	pre := "-1"
	c.latestCommit = c.edit.ConvJST((commits[0].Commit.Author.Date)).Format(c.edit.Layout)
	for _, v := range commits {
		target := c.edit.ConvJST(v.Commit.Author.Date)
		formatT := target.Format(c.edit.Layout)
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

	return nil
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

func (c *Client) InitStreak() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var commits []Commits
	errCount := 1
	for {
		if errCount == 5 {
			return errors.New("errorCount is at 5")
		}
		var err error
		commits, err = c.Getter.GetAllCommit(ctx)
		if err == nil {
			break
		}

		select {
		case <-ctx.Done():
			c.timeoutFlag = true
			time.Sleep(1 * time.Hour)
			errCount += 1
			continue
		}
	}

	return c.SetStreak(commits)
}
