package pageviews

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/mpaktiti/wikimedia-pageviews-api/src/utilities"
)

const baseURL = "https://wikimedia.org/api/rest_v1/metrics/pageviews/per-article/en.wikipedia/all-access/all-agents"

type Items struct {
	Items []Item
}

type Item struct {
	Project     string
	Article     string
	Granularity string
	Timestamp   string
	Access      string
	Agent       string
	Views       int
}

// curl http://localhost:8080/article/Albert_Einstein/weekly/2023/03
func GetPageviewsByWeek(article, year, week string) (int, error) {
	// Convert input year and week to integers
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return 0, err
	}
	weekInt, err := strconv.Atoi(week)
	if err != nil {
		return 0, err
	}

	// Validate that input year is not out of bounds
	err = utilities.ValidateInputYear(yearInt)
	if err != nil {
		return 0, err
	}

	// Validate that input week is not out of bounds
	err = utilities.ValidateInputWeek(year, weekInt)
	if err != nil {
		return 0, err
	}

	// Get the first and last days of the week
	startDate := utilities.WeekStart(yearInt, weekInt)
	endDate := startDate.AddDate(0, 0, 6)

	// Build URL
	// Wikipedia API expects months and days as 2 digits each so add a zero at the beginning if needed (done by PadString())
	startDateMonth := utilities.PadString(fmt.Sprint(int(startDate.Month())))
	startDateDay := utilities.PadString(fmt.Sprint(startDate.Day()))
	endDateMonth := utilities.PadString(fmt.Sprint(int(endDate.Month())))
	endDateDay := utilities.PadString(fmt.Sprint(endDate.Day()))
	firstDay := fmt.Sprint(startDate.Year()) + startDateMonth + startDateDay + "00"
	lastDay := fmt.Sprint(endDate.Year()) + endDateMonth + endDateDay + "00"
	url := fmt.Sprintf("%s/%s/daily/%s/%s", baseURL, article, firstDay, lastDay)

	// Call the wikipedia API
	response, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	// Parse response
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	// If the request was not successful parse the response for the error and return it
	if response.StatusCode != http.StatusOK {
		// return error
		errorDetails, err := utilities.ParseErrorDetails(responseData)
		if err != nil {
			fmt.Print(err.Error())
			return 0, fmt.Errorf(response.Status + ": Failed to process error details")
		}
		return 0, fmt.Errorf(response.Status + ": " + errorDetails)
	}

	// Aggregate view counts
	var items Items
	err = json.Unmarshal(responseData, &items)
	if err != nil {
		return 0, err
	}
	sum := 0
	for _, item := range items.Items {
		sum += item.Views
	}

	return sum, nil
}

// curl http://localhost:8080/article/Albert_Einstein/monthly/2023/04
func GetPageviewsByMonth(article, year, month string) (int, error) {
	// Convert input year to integer
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return 0, err
	}

	// TODO add test case for this
	// Validate that input year is not out of bounds
	err = utilities.ValidateInputYear(yearInt)
	if err != nil {
		return 0, err
	}

	lastOfMonth, err := utilities.LastDayOfMonth(year, month)
	if err != nil {
		return 0, err
	}

	// Build URL
	month = utilities.PadString(month)
	firstDay := year + month + "0100"
	lastDay := year + month + fmt.Sprint(lastOfMonth.Day()) + "00"
	url := fmt.Sprintf("%s/%s/monthly/%s/%s", baseURL, article, firstDay, lastDay)

	// Call the wikipedia API
	response, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	// Parse response
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	// If the request was not successful parse the response for the error and return it
	if response.StatusCode != http.StatusOK {
		// return error
		errorDetails, err := utilities.ParseErrorDetails(responseData)
		if err != nil {
			fmt.Print(err.Error())
			return 0, fmt.Errorf(response.Status + ": Failed to process error details")
		}
		return 0, fmt.Errorf(response.Status + ": " + errorDetails)
	}

	// Parse response and retrieve pageviews number
	var items Items
	err = json.Unmarshal(responseData, &items)
	if err != nil {
		return 0, err
	}

	return items.Items[0].Views, nil
}

// curl http://localhost:8080/article/Albert_Einstein/top/monthly/2023/04
func GetDayWithMostPageviews(article, year, month string) (string, int, error) {
	// Convert input year to integer
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return "", 0, err
	}

	// TODO add test case for this
	// Validate that input year is not out of bounds
	err = utilities.ValidateInputYear(yearInt)
	if err != nil {
		return "", 0, err
	}

	// Get month's last day
	lastOfMonth, err := utilities.LastDayOfMonth(year, month)
	if err != nil {
		return "", 0, err
	}

	// Build URL
	month = utilities.PadString(month)
	firstDay := year + month + "0100"
	lastDay := year + month + fmt.Sprint(lastOfMonth.Day()) + "00"
	url := fmt.Sprintf("%s/%s/daily/%s/%s", baseURL, article, firstDay, lastDay)

	// Call the wikipedia API
	response, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}

	// Parse response
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return "", 0, err
	}

	// If the request was not successful parse the response for the error and return it
	if response.StatusCode != http.StatusOK {
		// return error
		errorDetails, err := utilities.ParseErrorDetails(responseData)
		if err != nil {
			fmt.Print(err.Error())
			return "", 0, fmt.Errorf(response.Status + ": Failed to process error details")
		}
		return "", 0, fmt.Errorf(response.Status + ": " + errorDetails)
	}

	// Loop through results and find the max pageviews
	var topDay string
	var items Items
	err = json.Unmarshal(responseData, &items)
	if err != nil {
		return "", 0, err
	}
	topPageviews := 0
	for _, item := range items.Items {
		if item.Views > topPageviews {
			topPageviews = item.Views
			topDay = item.Timestamp
		}
	}

	return topDay, topPageviews, nil
}
