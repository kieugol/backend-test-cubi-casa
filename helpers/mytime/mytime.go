package mytime

import (
	"time"
)

var loc *time.Location

func SetTimezone(timezone string) error {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return err
	}
	loc = location
	return nil
}

func Now() time.Time {
	return time.Now().In(loc)
}

func NowUTC() time.Time {
	loc, _ = time.LoadLocation("UTC")
	return time.Now().In(loc)
}
