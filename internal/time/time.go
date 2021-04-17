package time

import (
	"time"

	"github.com/nhatthm/timeparser"
)

// Period returns a time period from given input.
func Period(now time.Time, from, to string) (start time.Time, end time.Time, err error) {
	if from == "" && to == "" {
		return now.AddDate(0, 0, -1), now, nil
	}

	if from != "" {
		start, err = timeparser.Parse(from)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	if to != "" {
		end, err = timeparser.Parse(to)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	if from == "" {
		return end.AddDate(0, 0, -1), end, nil
	}

	if to == "" {
		return start, start.AddDate(0, 0, 1), nil
	}

	if start.After(end) {
		return time.Time{}, time.Time{}, timeparser.ErrInvalidTimePeriod
	}

	return start, end, nil
}
