package utilities

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type ErrorResponse struct {
	Type   string
	Method string
	Detail string
	URI    string
}

type ErrorResponseWithMultipleDetails struct {
	Type   string
	Method string
	Detail []string
	URI    string
}

func ParseErrorDetails(response []byte) (string, error) {
	var errorResponse ErrorResponse
	err := json.Unmarshal(response, &errorResponse)
	if err != nil {
		if err.Error() == "json: cannot unmarshal array into Go struct field ErrorResponse.Detail of type string" {
			// the response contains an array of strings as details, use different object
			var errorResponse ErrorResponseWithMultipleDetails
			err := json.Unmarshal(response, &errorResponse)
			if err != nil {
				return "", err
			}
			errors := ""
			for _, error := range errorResponse.Detail {
				errors += error + ". "
			}
			return errors, nil
		}
		return "", err
	}

	return errorResponse.Detail, nil
}

func WeekStart(year, week int) time.Time {
	// Start from the middle of the year
	t := time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC)

	// Roll back to Monday
	if wd := t.Weekday(); wd == time.Sunday {
		t = t.AddDate(0, 0, -6)
	} else {
		t = t.AddDate(0, 0, -int(wd)+1)
	}

	// Difference in weeks
	_, w := t.ISOWeek()
	t = t.AddDate(0, 0, (week-w)*7)

	return t
}

// Returns the last day of the input month
func LastDayOfMonth(year, month string) (time.Time, error) {
	// Convert input year and week to integers
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return time.Time{}, err
	}
	monthInt, err := strconv.Atoi(month)
	if err != nil {
		return time.Time{}, err
	}
	firstOfMonth := time.Date(yearInt, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	// lastOfMonthAlternative := time.Date(yearInt, time.Month(monthInt+1), 0, 0, 0, 0, 0, time.UTC)

	return lastOfMonth, nil
}

// Returns the last week of the input year
func lastWeekOfYear(year string) (int, error) {
	yearToInt, err := strconv.Atoi(year)
	if err != nil {
		return 0, err
	}
	yearToDate := time.Date(yearToInt, time.January, 1, 0, 0, 0, 0, time.UTC)
	_, lastWeek := yearToDate.AddDate(1, 0, 0).Add(-time.Nanosecond).ISOWeek()

	return lastWeek, nil
}

// Validate that the input week is not greater than the last year of the input year
func ValidateInputWeek(year string, week int) error {
	validLastWeek, err := lastWeekOfYear(year)
	if err != nil {
		status500 := fmt.Sprint(http.StatusInternalServerError) + " " + http.StatusText(http.StatusInternalServerError)
		return fmt.Errorf(status500+": %s", err.Error())
	}
	if week > validLastWeek {
		status400 := fmt.Sprint(http.StatusBadRequest) + " " + http.StatusText(http.StatusBadRequest)
		return fmt.Errorf(status400+": input week cannot be greater than %d", validLastWeek)
	}
	return nil
}

// Validate that the input year is not greater than the current year
func ValidateInputYear(inputYear int) error {
	currentYear := time.Now().Year()
	if inputYear > currentYear {
		status400 := fmt.Sprint(http.StatusBadRequest) + " " + http.StatusText(http.StatusBadRequest)
		return fmt.Errorf(status400 + ": input year cannot be greater than current year")
	}
	return nil
}

// Wikipedia API expects months and days as 2 digits
// This function adds a zero at the beginning if needed
func PadString(input string) string {
	if len(input) == 1 {
		return "0" + input
	}
	return input
}
