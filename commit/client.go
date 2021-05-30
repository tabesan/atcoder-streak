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

func (c *Client) SetInitStreak(commits []Commits) error {
	mp := make(map[string]bool)
	var days []string
	pre := "-1"
	cnt := 0
	c.latestCommit = c.edit.ConvJST((commits[0].Commit.Author.Date)).Format(c.edit.Layout)
	for _, v := range commits {
		cnt += 1
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

func (c *Client) isStreak(target time.Time, later string) bool {
	if later == "-1" {
		return true
	}

	dayLater := target.AddDate(0, 0, 1).Format(c.edit.Layout)
	if dayLater == later {
		return true
	} else {
		return false
	}
}

func (c *Client) InitStreak() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	return c.SetInitStreak(commits)
}

func (c *Client) Update(ctx context.Context) error {
	resp, err := c.Getter.GetLastCommit(ctx)
	if err != nil {
		err = errors.Wrap(err, "GetLastCommit missed in Update")
	}
	latest := resp[0]
	lastDate := c.edit.ConvJST(latest.Commit.Author.Date)
	if lastDate.AddDate(0, 0, -1).Format(c.edit.Layout) == c.latestCommit {
		c.streak += 1
	} else {
		c.streak = 1
	}
	c.latestCommit = lastDate.Format(c.edit.Layout)
	return err
}

func (c *Client) UpdateStreak() error {
	if !c.updateFlag {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		var err error
		endCh := make(chan string)
		go func() {
			err = c.Update(ctx)
			if err != nil {
				c.timeoutFlag = false
				err = errors.Wrap(err, "Update error at UpdateStreak()")
				return
			}
			c.updateFlag = true
			c.timeoutFlag = false
			endCh <- "End"
		}()

		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			c.updateFlag = true
			return errors.New("UpdateStreak timeout")
		case <-endCh:
			return nil
		}
	}

	return nil
}

func (c *Client) ReferStreak() int {
	return c.streak
}

func (c *Client) ReferLatestCommit() string {
	return c.latestCommit
}
