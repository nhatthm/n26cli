package time_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	. "github.com/nhatthm/n26cli/internal/time" //nolint: revive
)

func TestPeriod(t *testing.T) {
	t.Parallel()

	tsStr := "2020-01-02T03:04:05Z"
	ts := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)

	testCases := []struct {
		scenario      string
		now           time.Time
		from          string
		to            string
		expectedStart time.Time
		expectedEnd   time.Time
		expectedError string
	}{
		{
			scenario:      "invalid from",
			from:          "foobar",
			expectedError: `parsing time "foobar" as "2006-01-02": cannot parse "foobar" as "2006"`,
		},
		{
			scenario:      "invalid to",
			to:            "foobar",
			expectedError: `parsing time "foobar" as "2006-01-02": cannot parse "foobar" as "2006"`,
		},
		{
			scenario:      "from and to are empty",
			now:           ts,
			expectedStart: ts.AddDate(0, 0, -1),
			expectedEnd:   ts,
		},
		{
			scenario:      "from is empty",
			now:           ts,
			to:            tsStr,
			expectedStart: ts.AddDate(0, 0, -1),
			expectedEnd:   ts,
		},
		{
			scenario:      "to is empty",
			now:           ts,
			from:          tsStr,
			expectedStart: ts,
			expectedEnd:   ts.AddDate(0, 0, 1),
		},
		{
			scenario:      "from and to are not empty and valid",
			now:           ts,
			from:          "2020-01-02T03:04:05Z",
			to:            "2020-02-02T03:04:05Z",
			expectedStart: time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC),
			expectedEnd:   time.Date(2020, 2, 2, 3, 4, 5, 0, time.UTC),
		},
		{
			scenario:      "from and to are not empty and invalid",
			now:           ts,
			from:          "2020-02-02T03:04:05Z",
			to:            "2020-01-02T03:04:05Z",
			expectedError: `invalid time period`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			start, end, err := Period(tc.now, tc.from, tc.to)

			assert.Equal(t, tc.expectedStart, start)
			assert.Equal(t, tc.expectedEnd, end)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
