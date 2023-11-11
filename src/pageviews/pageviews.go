package pageviews

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/mpaktiti/wikimedia-pageviews-api/internal/utilities"
)

// TODO Move this to the config file
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

// curl http://localhost:3000/article/Albert_Einstein/weekly/2023/03
func GetPageviewsByWeek(article, year, week string) (int, error) {
	// Convert input year and week to integers
	// TODO this is duplicate code with articles.go, extract it to utilities
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		// TODO log error properly
		fmt.Println("Error during conversion")
	}
	weekInt, err := strconv.Atoi(week)
	if err != nil {
		// TODO log error properly
		fmt.Println("Error during conversion")
	}
	// Get week range
	startDate := utilities.WeekStart(yearInt, weekInt)

	// Build URL
	// TODO duplicate code
	month := fmt.Sprint(int(startDate.Month()))
	if len(month) == 1 {
		// Wikipedia API expects months and days as 2 digits each so add a zero at the beginning if needed
		month = "0" + month
	}
	firstDay := fmt.Sprint(startDate.Year()) + month + fmt.Sprint(startDate.Day()) + "00"
	lastDay := fmt.Sprint(startDate.Year()) + month + fmt.Sprint(startDate.AddDate(0, 0, 7).Day()) + "00"
	url := fmt.Sprintf("%s/%s/daily/%s/%s", baseURL, article, firstDay, lastDay)

	// Call the wikipedia API
	response, err := http.Get(url)
	if err != nil {
		// TODO log error properly
		fmt.Print(err.Error())
		return 0, err
	}

	// Parse response
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	// fmt.Println(string(responseData))

	// Aggregate view counts
	var items Items
	err = json.Unmarshal(responseData, &items)
	if err != nil {
		fmt.Println("error:", err)
		return 0, err
	}
	sum := 0
	for _, item := range items.Items {
		sum += item.Views
	}

	return sum, nil
}

// curl http://localhost:3000/article/Albert_Einstein/monthly/2023/04
func GetPageviewsByMonth(article, year, month string) (int, error) {
	lastOfMonth, err := utilities.LastDayOfMonth(year, month)
	if err != nil {
		// TODO log error properly
		fmt.Println("Error calculating month's end date")
		return 0, err
	}

	// Build URL
	// TODO duplicate code
	if len(month) == 1 {
		// Wikipedia API expects months and days as 2 digits each so add a zero at the beginning if needed
		month = "0" + month
	}
	firstDay := year + month + "0100"
	lastDay := year + month + fmt.Sprint(lastOfMonth.Day()) + "00"
	url := fmt.Sprintf("%s/%s/monthly/%s/%s", baseURL, article, firstDay, lastDay)
	fmt.Println("First day: ", firstDay)
	fmt.Println("Last day: ", lastDay)
	fmt.Println("URL: ", url)

	// Call the wikipedia API
	response, err := http.Get(url)
	if err != nil {
		// TODO log error properly
		fmt.Print(err.Error())
		return 0, err
	}

	// Parse response
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	fmt.Println(string(responseData))

	// Parse response and retrieve pageviews number
	var items Items
	err = json.Unmarshal(responseData, &items)
	if err != nil {
		fmt.Println("error:", err)
		return 0, err
	}

	return items.Items[0].Views, nil
}

// curl http://localhost:3000/article/Albert_Einstein/top/monthly/2023/04
func GetDayWithMostPageviews(article, year, month string) (string, int, error) {
	// Get month's last day
	lastOfMonth, err := utilities.LastDayOfMonth(year, month)
	if err != nil {
		// TODO log error properly
		fmt.Println("Error calculating month's end date")
		return "", 0, err
	}

	// Build URL
	// TODO duplicate code
	if len(month) == 1 {
		// Wikipedia API expects months and days as 2 digits each so add a zero at the beginning if needed
		month = "0" + month
	}
	firstDay := year + month + "0100"
	lastDay := year + month + fmt.Sprint(lastOfMonth.Day()) + "00"
	url := fmt.Sprintf("%s/%s/daily/%s/%s", baseURL, article, firstDay, lastDay)
	fmt.Println("First day: ", firstDay)
	fmt.Println("Last day: ", lastDay)
	fmt.Println("URL: ", url)

	// Call the wikipedia API
	response, err := http.Get(url)
	if err != nil {
		// TODO log error properly
		fmt.Print(err.Error())
		return "", 0, err
	}

	// Parse response
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return "", 0, err
	}

	// Loop through results and find the max pageviews
	var topDay string
	var items Items
	err = json.Unmarshal(responseData, &items)
	if err != nil {
		fmt.Println("error:", err)
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
