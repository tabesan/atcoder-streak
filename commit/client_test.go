package commit

import (
	td "atcoder-streak/commit/test_data"
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient_createURL(t *testing.T) {
	const path = "repos/" + name + "/" + repo + "/commits"
	var err error
	c := NewClient(name, repo)
	NewMockClient(c)
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
	c := NewClient(name, repo)
	NewMockClient(c)

	t.Run("GetAll", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()

		result, err := c.Getter.GetAllCommit(ctx)
		if err != nil {
			t.Errorf("GetAllCommit error")
		}
		expect, err := toCommits(td.ResultAll)
		if err != nil {
			t.Errorf("toCommits missed")
		}
		if !reflect.DeepEqual(result, expect) {
			t.Errorf("result is not equals to expect")
		}
	})

	t.Run("GetLast", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()

		result, err := c.Getter.GetLastCommit(ctx)
		if err != nil {
			t.Errorf("GetLastCommit error")
		}
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
		client := NewClient(name, repo)
		NewMockClient(client)
		err := client.Getter.InitStreak()
		if err != nil {
			t.Errorf("InitStreak error")
		}
		assert.Equal(t, latest, client.latestCommit)
		assert.Equal(t, streak, client.streak)
	})

	t.Run("StreakOneDay", func(t *testing.T) {
		latest := "2021-05-11"
		streak := 1
		client := NewClient(name, repo)
		NewMockClient(client, "oneDay")
		client.Getter.InitStreak()
		assert.Equal(t, latest, client.latestCommit)
		assert.Equal(t, streak, client.streak)
	})
}
