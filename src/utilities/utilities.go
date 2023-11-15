package utilities

import (
	"encoding/json"
	"fmt"
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

// TODO write tests
func ParseErrorDetails(response []byte) (string, error) {
	var errorResponse ErrorResponse
	err := json.Unmarshal(response, &errorResponse)
	if err != nil {
		if err.Error() == "json: cannot unmarshal array into Go struct field ErrorResponse.Detail of type string" {
			// the response contains an array of strings as details, use different object
			var errorResponse ErrorResponseWithMultipleDetails
			err := json.Unmarshal(response, &errorResponse)
			if err != nil {
				fmt.Println("error:", err)
				return "", err
			}
			errors := ""
			for _, error := range errorResponse.Detail {
				errors += error + ". "
			}
			return errors, nil
		}
		fmt.Println("error:", err)
		return "", err
	}

	return errorResponse.Detail, nil

	// Try approach using map[string]interface{} instead of structs
	// var result map[string]interface{}
	// err := json.Unmarshal([]byte(string(response)), &result)
	// if err != nil {
	// 	// print out if error is not nil
	// 	fmt.Println(err)
	// }
	// for key, value := range result {
	// 	fmt.Println(key, ":", value)
	// }

	// return result["detail"].(string), nil
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

// Wikipedia API expects months and days as 2 digits
// This function adds a zero at the beginning if needed
func PadString(input string) string {
	if len(input) == 1 {
		return "0" + input
	}
	return input
}
