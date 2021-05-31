package commit

import (
	"testing"
	"time"
)

func TestChTimer_FlagTimer(t *testing.T) {
	c := NewTimer()
	c.flagInter = 3
	e := NewEditTime()
	go c.FlagTimer()
	endPoint := time.Now().In(e.ReferLocation()).Add(3 * time.Second)
	time.Sleep(time.Second)
	select {
	case <-c.ChFlag:
		if endPoint.Second() != time.Now().Second() {
			t.Errorf("FlagTimer error")
		}
	}
}

func TestChTimer_UpdateTimer(t *testing.T) {
	c := NewTimer()
	e := NewEditTime()
	start := 3597
	go c.UpdateTimer(start)
	endPoint := time.Now().In(e.ReferLocation()).Add(3 * time.Second)
	select {
	case <-c.ChUpdate:
		if endPoint.Second() != time.Now().Second() {
			t.Errorf("UpdateTimer error")
		}
	}
}
