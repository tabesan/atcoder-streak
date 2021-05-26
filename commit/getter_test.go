package commit

import (
	td "atcoder-streak/commit/test_data"
	"context"
	"reflect"
	"testing"
	"time"
)

func TestGetter_GetCommit(t *testing.T) {
	t.Run("GetAll", func(t *testing.T) {
		c := NewClient(name, repo)
		NewMockClient(c)
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
		c := NewClient(name, repo)
		NewMockClient(c)
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
