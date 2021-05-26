package commit

import "time"

func setLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		loc = time.FixedZone("Asia/Tokyo", 9*60*60)
	}

	return loc
}

type EditTime struct {
	location *time.Location
	Layout   string
}

func NewEditTime() *EditTime {
	e := &EditTime{
		location: setLocation(),
		Layout:   Layout,
	}

	return e
}

func (e *EditTime) ConvJST(t time.Time) time.Time {
	return t.In(e.location)
}

func (e *EditTime) ConvToSec(t time.Time) int {
	jst := e.ConvJST(t)
	toSec := jst.Hour()*3600 + jst.Minute()*60 + jst.Second()
	return toSec
}

func (e *EditTime) ReferLocation() *time.Location {
	return e.location
}
