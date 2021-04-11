package time

import (
	"errors"
	"time"

	"github.com/nhatthm/n26cli/internal/parser"
)

// ErrInvalidTimePeriod indicates that the time input is invalid.
var ErrInvalidTimePeriod = errors.New("invalid time period")

// Period returns a time period from given input.
func Period(now time.Time, from, to string) (start time.Time, end time.Time, err error) {
	if from == "" && to == "" {
		return now.AddDate(0, 0, -1), now, nil
	}

	if from != "" {
		start, err = parser.DateTime(from)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	if to != "" {
		end, err = parser.DateTime(to)
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
		return time.Time{}, time.Time{}, ErrInvalidTimePeriod
	}

	return start, end, nil
}
