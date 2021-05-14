package commit

import (
	td "atcoder-streak/commit/test_data"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
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

type mockClient struct{}

func (*mockClient) GetCommit(req *http.Request) []Commits {
	if req.URL.String() == "All" {
		resp, err := toCommits(td.ResultAll)
		if err != nil {
			return nil
		}
		return resp
	} else if req.URL.String() == "Last" {
		resp, err := toCommits(td.ResultLast)
		if err != nil {
			return nil
		}
		return resp
	}
	return nil
}

func AllCommitReq() *http.Request {
	reqAll, err := http.NewRequest("GET", "All", nil)
	if err != nil {
		return nil
	}

	return reqAll
}

func LastCommitReq() *http.Request {
	reqLast, err := http.NewRequest("GET", "Last", nil)
	if err != nil {
		return nil
	}

	return reqLast
}

func TestClient_NewcreateURL(t *testing.T) {
	name := "name"
	repo := "repo"
	c := NewClient(name, repo)

	if c == nil {
		t.Errorf("failed NewClient()")
	}

	path := "repos/name/repo/commits"
	assert.Equal(t, "https", c.URL.Scheme)
	assert.Equal(t, "api.github.com", c.URL.Host)
	assert.Equal(t, path, c.URL.Path)

	t.Logf("client: %p", c)
	t.Logf("Scheme: %s", c.URL.Scheme)
	t.Logf("Host: %s", c.URL.Host)
	t.Logf("Path: %s", c.URL.Path)
}

func TestClient_GetCommit(t *testing.T) {
	const name = "name"
	const repo = "repo"

	c := NewClient(name, repo)
	c.getter = &mockClient{}

	t.Run("GetAll", func(t *testing.T) {
		result := c.getter.GetCommit(AllCommitReq())
		expect, err := toCommits(td.ResultAll)
		if err != nil {
			t.Errorf("toCommits missed")
		}
		if !reflect.DeepEqual(result, expect) {
			t.Errorf("result is not equals to expect")
		}
	})

	t.Run("GetLast", func(t *testing.T) {
		result := c.getter.GetCommit(LastCommitReq())
		expect, err := toCommits(td.ResultLast)
		if err != nil {
			t.Errorf("toCommits missed")
		}
		if !reflect.DeepEqual(result, expect) {
			t.Errorf("result is not equals to expect")
		}
	})
}
