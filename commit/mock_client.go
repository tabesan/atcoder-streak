package commit

import (
	td "atcoder-streak/commit/test_data"
	tm "atcoder-streak/timer"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const (
	name    = "tabesan"
	repo    = "Atcoder"
	baseURL = "https://api.github.com/repos/tabesan/Atcoder/commits"
)

func toCommits(str string) ([]Commits, error) {
	bytes := []byte(str)
	var data []Commits
	err := json.Unmarshal(bytes, &data)
	if err != nil {
		fmt.Println(err)
	}
	return data, err
}

type MockClient struct {
	edit     *tm.EditTime
	client   *Client
	testData string
}

func NewMockClient(c *Client, testData ...string) {
	m := &MockClient{
		edit:     tm.NewEditTime(),
		testData: "twoDay",
		client:   c,
	}
	for _, t := range testData {
		m.testData = t
	}
	fmt.Println(m.testData)
	c.Getter = m
}

func (*MockClient) GetCommit(ctx context.Context, req *http.Request) ([]Commits, error) {
	if req.URL.String() == baseURL {
		resp, err := toCommits(td.ResultAll)
		if err != nil {
			return nil, errors.Wrap(err, "toCommits error in GetCommit")
		}
		return resp, nil
	} else if req.URL.String() == (baseURL + "?per_page=1") {
		resp, err := toCommits(td.ResultLast)
		if err != nil {
			return nil, errors.Wrap(err, "toCommits error in GetCommit")
		}
		return resp, nil
	}

	return nil, errors.New("GetCommit error")
}

func (m *MockClient) GetAllCommit(ctx context.Context) ([]Commits, error) {
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "NewRequest error in GetAllCommit")
	}
	resp, err := m.GetCommit(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "GetCommit error in GetAllCommit")
	}
	return resp, nil
}

func (m *MockClient) GetLastCommit(ctx context.Context) ([]Commits, error) {
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "NewRequest error in GetLastCommit")
	}
	u := req.URL.Query()
	u.Set("per_page", "1")
	req.URL.RawQuery = u.Encode()
	resp, err := m.GetCommit(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "GetCommit error in GetLastCommit")
	}
	return resp, nil
}

func (m *MockClient) InitStreak() error {
	var commits []Commits
	var ok bool
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client := NewClient(name, repo)
	switch m.testData {
	case "twoDay":
		fmt.Println("initStreak")
		NewMockClient(client)
		commits, err = m.GetAllCommit(ctx)
		if !ok {
			commits = nil
			return errors.Wrap(err, "GetAllCommit error in InitStreak")
		}
	case "oneDay":
		NewMockClient(client, "oneDay")
		commits, err = toCommits(td.StreakOneDay)
		if err != nil {
			commits = nil
			return errors.Wrap(err, "toCommits error in InitStreak")
		}
	}

	mp := make(map[string]bool)
	var days []string
	pre := "-1"

	client.latestCommit = client.edit.ConvJST((commits[0].Commit.Author.Date)).Format(m.client.edit.Layout)
	for _, v := range commits {
		target := client.edit.ConvJST(v.Commit.Author.Date)
		formatT := target.Format(client.edit.Layout)
		if !mp[formatT] {
			mp[formatT] = true
			if !m.isStreak(target, pre) {
				days = nil
			}
			pre = formatT
			days = append(days, formatT)
		}
	}

	client.streak = len(days)
	client.ShowStreak()
	return nil
}

func (m *MockClient) Update(ctx context.Context) error {
	resp, ok := m.GetLastCommit(ctx)
	fmt.Println(ok)
	latest := resp[0]
	lastDate := m.edit.ConvJST(latest.Commit.Author.Date)
	DayAgo := (lastDate.AddDate(0, 0, -1)).Format(m.client.edit.Layout)
	if m.client.latestCommit == DayAgo {
		m.client.streak += 1
		m.client.latestCommit = lastDate.Format(m.client.edit.Layout)
	}

	return nil
}

func (m *MockClient) isStreak(target time.Time, pre string) bool {
	if pre == "-1" {
		return true
	}

	dayLater := target.AddDate(0, 0, 1).Format(m.edit.Layout)
	if dayLater == pre {
		return true
	} else {
		return false
	}
}
