package commit

import (
	"testing"
	"time"
)

func TestEditTime_ConvJST(t *testing.T) {
	e := NewEditTime()
	target := time.Date(2021, 5, 11, 5, 0, 0, 0, time.UTC)
	result := e.ConvJST(target)
	expect := "2021-05-11 14:00:00 +0900 JST"

	if result.String() != expect {
		t.Errorf("ConvJST error")
	}
}

func TestEditTime_ConvToSec(t *testing.T) {
	e := NewEditTime()
	target := time.Date(2021, 5, 10, 5, 20, 15, 0, e.ReferLocation())
	result := e.ConvToSec(target)
	expect := 19215
	if result != expect {
		t.Errorf("ConvToSec error")
	}
}

func TestEditTime_setLocation(t *testing.T) {
	result := setLocation()
	expect := "Asia/Tokyo"
	if result.String() != expect {
		t.Errorf("setLocation error")
	}
}

func TestEditTime_ReferLocation(t *testing.T) {
	e := NewEditTime()
	result := e.ReferLocation().String()
	expect := "Asia/Tokyo"
	if result != expect {
		t.Errorf("ReferLocation error")
	}
}
