package commit

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	name = Name
	repo = Repository
)

func TestClient_createURL(t *testing.T) {
	fmt.Println("TestClient_createURL")
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

func TestClient_InitStreak(t *testing.T) {
	t.Run("StreakTwoDays", func(t *testing.T) {
		fmt.Println("Test StreakTwoDays")
		latest := "2021-05-11"
		streak := 2
		client := NewClient(name, repo)
		NewMockClient(client)
		err := client.InitStreak()
		if err != nil {
			t.Errorf("InitStreak error")
		}
		assert.Equal(t, latest, client.latestCommit)
		assert.Equal(t, streak, client.streak)
	})

	t.Run("StreakOneDay", func(t *testing.T) {
		fmt.Println("Test StreakOneDay")
		latest := "2021-05-11"
		streak := 1
		client := NewClient(name, repo)
		NewMockClient(client, "oneDay")
		client.InitStreak()
		assert.Equal(t, latest, client.latestCommit)
		assert.Equal(t, streak, client.streak)
	})
}

func TestClient_isStreak(t *testing.T) {
	c := NewClient(name, repo)
	t.Run("later == -1", func(t *testing.T) {
		later := "-1"
		target := time.Now()
		result := c.isStreak(target, later)
		expect := true
		if result != expect {
			t.Errorf("pre == -1 error")
		}
	})

	t.Run("Is streak", func(t *testing.T) {
		target := c.edit.ConvJST(time.Date(2001, 05, 20, 23, 0, 0, 0, c.edit.ReferLocation()))
		later := "2001-05-21"
		result := c.isStreak(target, later)
		expect := true
		if result != expect {
			t.Errorf("Is streak error")
		}
	})

	t.Run("Is not streak", func(t *testing.T) {
		target := c.edit.ConvJST(time.Date(2001, 05, 20, 23, 0, 0, 0, c.edit.ReferLocation()))
		later := "2001-05-22"
		result := c.isStreak(target, later)
		expect := false
		if result != expect {
			t.Errorf("Is not streak error")
		}
	})
}

func TestClient_Update(t *testing.T) {
	t.Run("Is streak", func(t *testing.T) {
		fmt.Println("Test Update Is streak")
		c := NewClient(name, repo)
		NewMockClient(c)
		latest := time.Date(2021, 5, 10, 5, 0, 0, 0, c.edit.ReferLocation())
		c.latestCommit = c.edit.ConvJST(latest).Format(c.edit.Layout)
		c.streak = 1
		c.UpdateStreak()
		expectStreak := 2
		expectLatest := "2021-05-11"

		if expectStreak != c.ReferStreak() {
			t.Errorf("streak update missed")
		}
		if expectLatest != c.ReferLatestCommit() {
			t.Errorf("latestCommit is wrong")
		}
	})

	t.Run("Is not streak", func(t *testing.T) {
		fmt.Println("Teset Update Is not streak")
		c := NewClient(name, repo)
		NewMockClient(c)
		latest := time.Date(2021, 5, 9, 5, 0, 0, 0, time.UTC)
		c.latestCommit = c.edit.ConvJST(latest).Format(c.edit.Layout)
		c.streak = 1
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		c.Update(ctx)
		expectStreak := 1
		expectLatest := "2021-05-11"
		if expectStreak != c.streak {
			t.Errorf("streak update missed")
		}
		if expectLatest != c.latestCommit {
			t.Errorf("latestCommit is wrong")
		}
	})
}

func TestClient_UpdateStreak(t *testing.T) {
	t.Run("update", func(t *testing.T) {
		fmt.Println("Test UpdateStreak update")
		c := NewClient(name, repo)
		NewMockClient(c)
		latest := time.Date(2021, 5, 10, 5, 0, 0, 0, c.edit.ReferLocation())
		c.latestCommit = c.edit.ConvJST(latest).Format(c.edit.Layout)
		c.streak = 1
		c.UpdateStreak()

		expectStreak := 2
		expectLatest := "2021-05-11"
		if expectStreak != c.streak {
			t.Errorf("Updated streak is wrong")
		}
		if expectLatest != c.latestCommit {
			t.Errorf("Updated latestCommit is wrong")
		}
	})
	t.Run("!updateFlag", func(t *testing.T) {
		fmt.Println("Test UpdateStreak !updateFlag")
		c := NewClient(name, repo)
		NewMockClient(c)
		latest := time.Date(2021, 5, 10, 5, 0, 0, 0, time.UTC)
		c.latestCommit = c.edit.ConvJST(latest).Format(c.edit.Layout)
		c.updateFlag = false
		c.timeoutFlag = true
		c.UpdateStreak()
		if !c.updateFlag {
			t.Errorf("update missed")
		}
		if c.updateFlag && c.timeoutFlag {
			t.Errorf("UpdateStreak error")
		}
	})

	t.Run("Timeout", func(t *testing.T) {
		fmt.Println("Test UpdateStreak Timeout")
		c := NewClient(name, repo)
		NewMockClient(c)
		rtd := c.Getter.RefTestPara().(*RefTest)
		rtd.SetIsTimeOut()
		c.UpdateStreak()
	})
}
