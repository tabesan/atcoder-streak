package commit

import (
	td "atcoder-streak/commit/test_data"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	client   *Client
	URL      *url.URL
	testData string
}

func NewMockClient() (*MockClient, error) {
	u, err := url.Parse(baseURL)
	m := &MockClient{
		client:   NewClient(name, repo),
		URL:      u,
		testData: "twoDay",
	}
	return m, err
}

func (*MockClient) GetCommit(req *http.Request) []Commits {
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

func (m *MockClient) GetAllCommit() []Commits {
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil
	}
	return m.GetCommit(req)
}

func (m *MockClient) GetLastCommit() []Commits {
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil
	}
	u := req.URL.Query()
	u.Set("per_page", "1")
	req.URL.RawQuery = u.Encode()
	return m.GetCommit(req)
}

func (m *MockClient) InitStreak() {
	c := NewClient(name, repo)
	var commits []Commits
	var err error
	switch m.testData {
	case "twoDay":
		commits = m.GetAllCommit()
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

	m.client.latestCommit = c.edit.ConvJST((commits[0].Commit.Author.Date)).Format(c.edit.Layout)
	for _, v := range commits {
		target = c.edit.ConvJST(v.Commit.Author.Date)
		formatT = target.Format(c.edit.Layout)
		if !mp[formatT] {
			mp[formatT] = true
			if !m.client.isStreak(target, pre) {
				days = nil
			}
			pre = formatT
			days = append(days, formatT)
		}
	}

	m.client.streak = len(days)
	c.ShowStreak()
}

func TestClient_createURL(t *testing.T) {
	const path = "repos/" + name + "/" + repo + "/commits"
	var err error
	c := NewClient(name, repo)
	if c == nil {
		t.Errorf("failed NewClient()")
	}

	if err != nil {
		t.Errorf("url.Parse missed at NewMockClient()")
	}

	assert.Equal(t, "https", c.URL.Scheme)
	assert.Equal(t, "api.github.com", c.URL.Host)
	assert.Equal(t, path, c.URL.Path)

	t.Logf("client: %p", c)
	t.Logf("Scheme: %s", c.URL.Scheme)
	t.Logf("Host: %s", c.URL.Host)
	t.Logf("Path: %s", c.URL.Path)
}

func TestClient_GetCommit(t *testing.T) {
	m, err := NewMockClient()
	if err != nil {
		t.Errorf("NewMockClient() missed at GetCommit")
	}

	t.Run("GetAll", func(t *testing.T) {
		result := m.GetAllCommit()
		expect, err := toCommits(td.ResultAll)
		if err != nil {
			t.Errorf("toCommits missed")
		}
		if !reflect.DeepEqual(result, expect) {
			t.Errorf("result is not equals to expect")
		}
	})

	t.Run("GetLast", func(t *testing.T) {
		result := m.GetLastCommit()
		expect, err := toCommits(td.ResultLast)
		if err != nil {
			t.Errorf("toCommits missed")
		}
		if !reflect.DeepEqual(result, expect) {
			t.Errorf("result is not equals to expect")
		}
	})
}

func TestClient_InitStreak(t *testing.T) {
	t.Run("StreakTwoDays", func(t *testing.T) {
		latest := "2021-05-11"
		streak := 2
		m, err := NewMockClient()
		if err != nil {
			t.Errorf("NewMockClient() missed at InitStreak()")
		}
		m.client = NewClient(name, repo)
		m.InitStreak()
		assert.Equal(t, latest, m.client.latestCommit)
		assert.Equal(t, streak, m.client.streak)
	})

	t.Run("StreakOneDay", func(t *testing.T) {
		latest := "2021-05-11"
		streak := 1
		m, err := NewMockClient()
		if err != nil {
			t.Errorf("NewMockClient() missed at InitStreak()")
		}
		m.client = NewClient(name, repo)
		m.testData = "oneDay"
		m.InitStreak()
		assert.Equal(t, latest, m.client.latestCommit)
		assert.Equal(t, streak, m.client.streak)
	})
}
