package commit

import (
	td "atcoder-streak/commit/test_data"
	tm "atcoder-streak/timer"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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

func NewMockClient(c *Client) {
	m := &MockClient{
		edit:     tm.NewEditTime(),
		testData: "TwoDay",
	}
	m.client = c
}

func (*MockClient) GetCommit(ctx context.Context, req *http.Request) []Commits {
	if req.URL.String() == baseURL {
		resp, err := toCommits(td.ResultAll)
		if err != nil {
			return nil
		}
		return resp
	} else if req.URL.String() == (baseURL + "?per_page=1") {
		resp, err := toCommits(td.ResultLast)
		if err != nil {
			return nil
		}
		return resp
	}
	return nil
}

func (m *MockClient) GetAllCommit(ctx context.Context) ([]Commits, bool) {
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, false
	}
	return m.GetCommit(ctx, req), true
}

func (m *MockClient) GetLastCommit(ctx context.Context) ([]Commits, bool) {
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, false
	}
	u := req.URL.Query()
	u.Set("per_page", "1")
	req.URL.RawQuery = u.Encode()
	return m.GetCommit(ctx, req), true
}

func (m *MockClient) InitStreak() {
	var commits []Commits
	var ok bool
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	switch m.testData {
	case "twoDay":
		commits, ok = m.GetAllCommit(ctx)
		if !ok {
			commits = nil
			return
		}
	case "oneDay":
		commits, err = toCommits(td.StreakOneDay)
		if err != nil {
			commits = nil
			return
		}
	}
	mp := make(map[string]bool)
	var days []string
	var target time.Time
	var formatT string
	pre := "-1"

	m.client.latestCommit = m.client.edit.ConvJST((commits[0].Commit.Author.Date)).Format(m.client.edit.Layout)
	for _, v := range commits {
		target = m.client.edit.ConvJST(v.Commit.Author.Date)
		formatT = target.Format(m.client.edit.Layout)
		if !mp[formatT] {
			mp[formatT] = true
			if !m.isStreak(target, pre) {
				days = nil
			}
			pre = formatT
			days = append(days, formatT)
		}
	}

	m.client.streak = len(days)
	m.client.ShowStreak()
}

func (m *MockClient) Update(ctx context.Context) {
	resp, ok := m.GetLastCommit(ctx)
	fmt.Println(ok)
	latest := resp[0]
	lastDate := m.edit.ConvJST(latest.Commit.Author.Date)
	DayAgo := (lastDate.AddDate(0, 0, -1)).Format(m.client.edit.Layout)
	if m.client.latestCommit == DayAgo {
		m.client.streak += 1
		m.client.latestCommit = lastDate.Format(m.client.edit.Layout)
	}
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
