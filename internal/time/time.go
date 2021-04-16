package time

import (
	"time"

	"github.com/nhatthm/timeparser"
)

// Period returns a time period from given input.
func Period(now time.Time, from, to string) (time.Time, time.Time, error) {
	s, e, err := timeparser.ParsePeriod(from, to)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	if s == nil && e == nil {
		return now.AddDate(0, 0, -1), now, nil
	}

	if s == nil {
		return e.AddDate(0, 0, -1), *e, nil
	}

	if e == nil {
		return *s, s.AddDate(0, 0, 1), nil
	}

	return *s, *e, nil
}
