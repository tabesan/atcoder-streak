package commit

import (
	"fmt"
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
		fmt.Println("cache")
		if endPoint.Second() != time.Now().Second() {
			t.Errorf("FlagTimer error")
		}
	}
}

func TestChTimer_UpdateTimer(t *testing.T) {
	c := NewTimer()
	c.updateInter = 3
	e := NewEditTime()
	go c.UpdateTimer(time.Second)
	endPoint := time.Now().In(e.ReferLocation()).Add(3 * time.Second)
	select {
	case <-c.ChUpdate:
		if endPoint.Second() != time.Now().Second() {
			t.Errorf("UpdateTimer error")
		}
	}
}
