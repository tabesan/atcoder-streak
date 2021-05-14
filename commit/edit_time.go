package commit

import "time"

func setLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		loc = time.FixedZone("Asia/Tokyo", 9*60*60)
	}

	return loc
}

type editTime struct {
	location *time.Location
}

func NewEditTime() *editTime {
	e := &editTime{
		location: setLocation(),
	}

	return e
}

func (e *editTime) convJST(t time.Time) time.Time {
	return t.In(e.location)
}

func (e *editTime) convToSec(t time.Time) int {
	jst := e.convJST(t)
	toSec := jst.Hour()*3600 + jst.Minute()*60 + jst.Second()
	return toSec
}
