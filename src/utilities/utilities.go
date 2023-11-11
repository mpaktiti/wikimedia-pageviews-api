package utilities

import (
	"fmt"
	"strconv"
	"time"
)

func WeekStart(year, week int) time.Time {
	// Start from the middle of the year:
	t := time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC)

	// Roll back to Monday:
	if wd := t.Weekday(); wd == time.Sunday {
		t = t.AddDate(0, 0, -6)
	} else {
		t = t.AddDate(0, 0, -int(wd)+1)
	}

	// Difference in weeks:
	_, w := t.ISOWeek()
	t = t.AddDate(0, 0, (week-w)*7)

	return t
}

// TODO remove this, not needed
// TODO do I wanna return arguments like this in all functions?
func WeekRange(year, week int) (start, end time.Time) {
	start = WeekStart(year, week)
	end = start.AddDate(0, 0, 6)
	return
}

// Returns the last day of the input month
func LastDayOfMonth(year, month string) (time.Time, error) {
	// Convert input year and week to integers
	// TODO this is duplicate code with articles.go, extract it to utilities
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		// TODO log error properly
		fmt.Println("Error during conversion")
		return time.Time{}, err
	}
	monthInt, err := strconv.Atoi(month)
	if err != nil {
		// TODO log error properly
		fmt.Println("Error during conversion")
		return time.Time{}, err
	}
	firstOfMonth := time.Date(yearInt, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	// lastOfMonthAlternative := time.Date(yearInt, time.Month(monthInt+1), 0, 0, 0, 0, 0, time.UTC)

	return lastOfMonth, nil
}
