package commit

import (
	tm "atcoder-streak/timer"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type Getter interface {
	GetCommit(ctx context.Context, req *http.Request) ([]Commits, error)
	GetLastCommit(ctx context.Context) ([]Commits, error)
	GetAllCommit(ctx context.Context) ([]Commits, error)
	InitStreak() error
	Update(ctx context.Context) error
}

type getter struct {
	edit   *tm.EditTime
	client *Client
}

func NewGetter(c *Client) {
	g := &getter{
		edit:   tm.NewEditTime(),
		client: c,
	}
	c.Getter = g
}

func (g *getter) GetLastCommit(ctx context.Context) ([]Commits, error) {
	req, err := http.NewRequest("GET", g.client.URL.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "NewRequest missed in GetLastCommit")
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	q := req.URL.Query()
	q.Set("per_page", "1")
	req.URL.RawQuery = q.Encode()
	resp, err := g.GetCommit(ctx, req)
	return resp, nil
}

func (g *getter) GetAllCommit(ctx context.Context) ([]Commits, error) {
	req, err := http.NewRequest("GET", g.client.URL.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "NewRequest missed in GetAllCommit")
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	resp, err := g.GetCommit(ctx, req)
	return resp, nil
}

func (g *getter) GetCommit(ctx context.Context, req *http.Request) ([]Commits, error) {
	resp, err := g.client.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "HTTPClient.Do missed in GetCommit")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "ReadAll missed in GetCommit")
	}

	bytes := []byte(body)
	var data []Commits
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal missed in GetCommit")
	}
	return data, nil
}

func (g *getter) isStreak(target time.Time, pre string) bool {
	if pre == "-1" {
		return true
	}

	dayLater := target.AddDate(0, 0, 1).Format(g.edit.Layout)
	if dayLater == pre {
		return true
	} else {
		return false
	}
}

func (g *getter) InitStreak() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var commits []Commits
	errorCount := 1
	for {
		if errorCount == 5 {
			return errors.New("errorCount is at 5")
		}
		var err error
		commits, err = g.GetAllCommit(ctx)
		if err == nil {
			break
		}

		select {
		case <-ctx.Done():
			g.client.timeoutFlag = true
			time.Sleep(1 * time.Hour)
			errorCount += 1
			continue
		}
	}

	mp := make(map[string]bool)
	var days []string
	pre := "-1"
	g.client.latestCommit = g.edit.ConvJST((commits[0].Commit.Author.Date)).Format(g.edit.Layout)
	for _, v := range commits {
		target := g.edit.ConvJST(v.Commit.Author.Date)
		formatT := target.Format(g.edit.Layout)
		if !mp[formatT] {
			mp[formatT] = true
			if !g.isStreak(target, pre) {
				days = nil
			}
			pre = formatT
			days = append(days, formatT)
		}
	}
	g.client.streak = len(days)

	return nil
}

func (g *getter) Update(ctx context.Context) error {
	resp, err := g.GetLastCommit(ctx)
	if err != nil {
		return errors.Wrap(err, "GetLastCommit missed in Update")
	}

	latest := resp[0]
	lastDate := g.edit.ConvJST(latest.Commit.Author.Date)
	DayAgo := (lastDate.AddDate(0, 0, -1)).Format(g.client.edit.Layout)
	if g.client.latestCommit == DayAgo {
		g.client.streak += 1
		g.client.latestCommit = lastDate.Format(g.client.edit.Layout)
	}
	return nil
}
