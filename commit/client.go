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
	resetFlag     bool
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
		Path:   "repos/" + Name + "/" + Repository + "/commits",
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
			Timeout: time.Second * 10,
		},
		edit: tm.NewEditTime(),
	}
	c.createURL()
	return c
}

func (c *Client) ShowStreak() {
	fmt.Println(c.streak)
}

func (c *Client) SetStreak(s int) {
	c.streak = s
}

func (c *Client) SetLatest(l string) {
	c.latestCommit = l
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
	useCurStreak := false
	c.latestCommit = c.edit.ConvJST((commits[0].Commit.Author.Date)).Format(c.edit.Layout)
	for _, v := range commits {
		target := c.edit.ConvJST(v.Commit.Author.Date)
		formatT := target.Format(c.edit.Layout)
		if !mp[formatT] {
			mp[formatT] = true
			if !c.isStreak(target, pre) {
				break
			}
			if c.latestCommit == formatT && c.streak != 0 {
				useCurStreak = true
				break
			}
			pre = formatT
			days = append(days, formatT)
		}
	}
	if useCurStreak {
		c.streak = len(days) + c.streak
	} else {
		c.streak = len(days)
	}
	c.timeoutFlag = false
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
	ctx, cancel := context.WithTimeout(context.Background(), c.HTTPClient.Timeout)
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
		return err
	}
	latest := resp[0]
	lastDate := c.edit.ConvJST(latest.Commit.Author.Date)
	if lastDate.Format(c.edit.Layout) == c.latestCommit {
		return nil
	}
	if lastDate.AddDate(0, 0, -1).Format(c.edit.Layout) == c.latestCommit {
		c.streak += 1
	} else {
		c.streak = 1
	}
	c.latestCommit = lastDate.Format(c.edit.Layout)
	return nil
}

func (c *Client) UpdateStreak() error {
	if !c.updateFlag {
		ctx, cancel := context.WithTimeout(context.Background(), c.HTTPClient.Timeout)
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

		select {
		case <-ctx.Done():
			c.timeoutFlag = true
			return errors.New("UpdateStreak timeout")
		case <-endCh:
			return err
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

func (c *Client) ReferTimeoutFlag() bool {
	return c.timeoutFlag
}

func (c *Client) SetUpdateFlag(tf bool) {
	c.updateFlag = tf
}

func (c *Client) ReferUpdateFlag() bool {
	return c.updateFlag
}

func (c *Client) SetResetFlag(f bool) {
	c.resetFlag = f
}

func (c *Client) ReferResetFlag() bool {
	return c.resetFlag
}
