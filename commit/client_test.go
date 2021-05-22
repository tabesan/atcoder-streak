package commit

import (
	td "atcoder-streak/commit/test_data"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_createURL(t *testing.T) {
	const path = "repos/" + name + "/" + repo + "/commits"
	var err error
	getter := NewMockClient()
	c := NewClient(name, repo, getter)
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
