package parser

import "time"

const (
	layoutRFC3339 = time.RFC3339
	layoutYMD     = "2006-01-02"
)

// DateTime parses string into time.Time.
func DateTime(s string) (time.Time, error) {
	if len(s) > 10 {
		return time.Parse(layoutRFC3339, s)
	}

	return time.Parse(layoutYMD, s)
}
