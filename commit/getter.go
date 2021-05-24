package commit

import (
	tm "atcoder-streak/timer"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type Getter interface {
	GetCommit(ctx context.Context, req *http.Request) ([]Commits, error)
	GetLastCommit(ctx context.Context) ([]Commits, error)
	GetAllCommit(ctx context.Context) ([]Commits, error)
	Update(ctx context.Context) error
	RefTestData() interface{}
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

func (g *getter) RefTestData() interface{} {
	return nil
}
