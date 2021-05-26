package commit

import (
	td "atcoder-streak/commit/test_data"
	tm "atcoder-streak/timer"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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

type RefTest struct {
	testData string
}

func NewRefTest() *RefTest {
	r := &RefTest{
		testData: "twoDay",
	}
	return r
}

func (r *RefTest) SetTestData(str string) {
	r.testData = str
}

func (r *RefTest) RefTestData() string {
	return r.testData
}

type MockClient struct {
	edit    *tm.EditTime
	client  *Client
	refTest *RefTest
}

func NewMockClient(c *Client, testData ...string) {
	m := &MockClient{
		edit:    tm.NewEditTime(),
		refTest: NewRefTest(),
		client:  c,
	}
	for _, t := range testData {
		m.refTest.testData = t
	}
	c.Getter = m
}

func (m *MockClient) RefTestData() interface{} {
	return m.refTest
}

func (m *MockClient) GetCommit(ctx context.Context, req *http.Request) ([]Commits, error) {
	var data []Commits
	var err error

	if req.URL.String() == baseURL {
		rtd := m.RefTestData().(*RefTest)
		oneOrTwo := rtd.RefTestData()
		switch oneOrTwo {
		case "oneDay":
			data, err = toCommits(td.StreakOneDay)
			if err != nil {
				return data, errors.Wrap(err, "toCommit error in GetCommit")
			}
		default:
			data, err = toCommits(td.ResultAll)
			if err != nil {
				return data, errors.Wrap(err, "toCommits error in GetCommit")
			}
		}
	} else if req.URL.String() == (baseURL + "?per_page=1") {
		data, err = toCommits(td.ResultLast)
		if err != nil {
			return nil, errors.Wrap(err, "toCommits error in GetCommit")
		}
	}

	return data, nil
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

//func (m *MockClient) Update(ctx context.Context) error {
//	resp, ok := m.GetLastCommit(ctx)
//	fmt.Println(ok)
//	latest := resp[0]
//	lastDate := m.edit.ConvJST(latest.Commit.Author.Date)
//	DayAgo := (lastDate.AddDate(0, 0, -1)).Format(m.client.edit.Layout)
//	if m.client.latestCommit == DayAgo {
//		m.client.streak += 1
//		m.client.latestCommit = lastDate.Format(m.client.edit.Layout)
//	}
//
//	return nil
//}
