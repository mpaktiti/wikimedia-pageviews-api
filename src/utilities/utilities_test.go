package utilities

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWeekStart(t *testing.T) {
	testCases := []struct {
		name        string
		year        int
		week        int
		expectedDay int
	}{
		{
			name:        "week 01",
			year:        2023,
			week:        01,
			expectedDay: 2,
		},
		{
			name:        "week 31",
			year:        2023,
			week:        31,
			expectedDay: 31,
		},
		{
			name:        "week 52",
			year:        2023,
			week:        52,
			expectedDay: 25,
		},
	}
	for tcNum, tc := range testCases {
		got := WeekStart(tc.year, tc.week)
		if got.Day() != tc.expectedDay {
			t.Errorf("test %d failed: got %v want %v", tcNum+1, got.Day(), tc.expectedDay)
		}
	}
}

func TestLastDayOfMonth(t *testing.T) {
	testCases := []struct {
		name        string
		year        string
		month       string
		expectedDay time.Time
	}{
		{
			name:        "January 2020",
			year:        "2020",
			month:       "01",
			expectedDay: time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:        "February 2020",
			year:        "2020",
			month:       "02",
			expectedDay: time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			name:        "February 2022",
			year:        "2022",
			month:       "02",
			expectedDay: time.Date(2022, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name:        "December 2023",
			year:        "2023",
			month:       "12",
			expectedDay: time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
		},
	}
	for tcNum, tc := range testCases {
		got, _ := LastDayOfMonth(tc.year, tc.month)
		if got != tc.expectedDay {
			t.Errorf("test %d failed: got %v want %v", tcNum+1, got, tc.expectedDay)
		}
	}
}

func TestPadString(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedOutput string
	}{
		{
			name:           "string with 1 character",
			input:          "1",
			expectedOutput: "01",
		},
		{
			name:           "string with 2 characters",
			input:          "01",
			expectedOutput: "01",
		},
	}
	for tcNum, tc := range testCases {
		got := PadString(tc.input)
		if got != tc.expectedOutput {
			t.Errorf("test %d failed:  got %v want %v", tcNum+1, got, tc.expectedOutput)
		}
	}
}

func TestParseErrorDetails(t *testing.T) {
	testCases := []struct {
		name           string
		input          []byte
		expectedOutput string
	}{
		{
			name:           "error contains more than 1 error details (array of strings)",
			input:          []byte(`{"type":"https://mediawiki.org/wiki/HyperSwitch/errors/invalid_request","method":"get","detail":["end timestamp is invalid, must be a valid date in YYYYMMDD format"],"uri":"/analytics.wikimedia.org/v1/pageviews/per-article/en.wikipedia/all-access/all-agents/Albert_Einstein/daily/2031063000/203106700"}`),
			expectedOutput: "end timestamp is invalid, must be a valid date in YYYYMMDD format. ",
		},
		{
			name:           "error contains a single entry for error details (string)",
			input:          []byte(`{"type":"https://mediawiki.org/wiki/HyperSwitch/errors/not_found","title":"Not found.","method":"get","detail":"The date(s) you used are valid, but we either do not have data for those date(s), or the project you asked for is not loaded yet.  Please check https://wikimedia.org/api/rest_v1/?doc for more information.","uri":"/analytics.wikimedia.org/v1/pageviews/per-article/en.wikipedia/all-access/all-agents/Albert_Einstein/daily/2025011300/2025012000"}`),
			expectedOutput: "The date(s) you used are valid, but we either do not have data for those date(s), or the project you asked for is not loaded yet.  Please check https://wikimedia.org/api/rest_v1/?doc for more information.",
		},
	}
	for tcNum, tc := range testCases {
		got, err := ParseErrorDetails(tc.input)
		require.NoError(t, err)
		if got != tc.expectedOutput {
			t.Errorf("test %d failed:  got %v want %v", tcNum+1, got, tc.expectedOutput)
		}
	}
}
