package utilities

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWeekStart(t *testing.T) {
	testCases := []struct {
		name           string
		year           int
		week           int
		expectedOutput int
	}{
		{
			name:           "week 01",
			year:           2023,
			week:           01,
			expectedOutput: 2,
		},
		{
			name:           "week 31",
			year:           2023,
			week:           31,
			expectedOutput: 31,
		},
		{
			name:           "week 52",
			year:           2023,
			week:           52,
			expectedOutput: 25,
		},
	}
	for tcNum, tc := range testCases {
		got := WeekStart(tc.year, tc.week)
		assertExpectedOutput(t, tcNum, got.Day(), tc.expectedOutput)
	}
}

func TestLastDayOfMonth(t *testing.T) {
	testCases := []struct {
		name           string
		year           string
		month          string
		expectedOutput time.Time
	}{
		{
			name:           "January 2020",
			year:           "2020",
			month:          "01",
			expectedOutput: time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:           "February 2020",
			year:           "2020",
			month:          "02",
			expectedOutput: time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			name:           "February 2022",
			year:           "2022",
			month:          "02",
			expectedOutput: time.Date(2022, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name:           "December 2023",
			year:           "2023",
			month:          "12",
			expectedOutput: time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
		},
	}
	for tcNum, tc := range testCases {
		got, _ := LastDayOfMonth(tc.year, tc.month)
		assertExpectedOutput(t, tcNum, got, tc.expectedOutput)
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
		assertExpectedOutput(t, tcNum, got, tc.expectedOutput)
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
		assertExpectedOutput(t, tcNum, got, tc.expectedOutput)
	}
}

func TestValidateInputYear(t *testing.T) {
	testCases := []struct {
		name           string
		input          int
		expectedOutput error
	}{
		{
			name:           "valid year",
			input:          2023,
			expectedOutput: nil,
		},
		{
			name:           "invalid year",
			input:          2030,
			expectedOutput: fmt.Errorf("400 Bad Request: input year cannot be greater than current year "),
		},
	}
	for _, tc := range testCases {
		got := ValidateInputYear(tc.input)
		if tc.expectedOutput != nil {
			require.Error(t, got)
		} else {
			require.NoError(t, got)
		}
	}
}

func TestValidateInputWeek(t *testing.T) {
	testCases := []struct {
		name           string
		year           string
		week           int
		expectedOutput error
	}{
		{
			name:           "valid week",
			year:           "2023",
			week:           3,
			expectedOutput: nil,
		},
		{
			name:           "invalid week",
			year:           "2020",
			week:           54,
			expectedOutput: fmt.Errorf("400 Bad Request: input week cannot be greater than 53"),
		},
	}
	for _, tc := range testCases {
		got := ValidateInputWeek(tc.year, tc.week)
		if tc.expectedOutput != nil {
			require.Error(t, got)
		} else {
			require.NoError(t, got)
		}
	}
}

func assertExpectedOutput(t testing.TB, testNum int, got, want interface{}) {
	t.Helper()
	if got != want {
		t.Errorf("test %d failed: got %v want %v", testNum+1, got, want)
	}
}
