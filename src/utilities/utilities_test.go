package utilities

import (
	"testing"
	"time"
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
	for _, tc := range testCases {
		got := WeekStart(tc.year, tc.week)
		if got.Day() != tc.expectedDay {
			t.Errorf("code returned unexpected day: got %v want %v", got.Day(), tc.expectedDay)
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
	for _, tc := range testCases {
		got, _ := LastDayOfMonth(tc.year, tc.month)
		if got != tc.expectedDay {
			t.Errorf("code returned unexpected day: got %v want %v", got, tc.expectedDay)
		}
	}
}
