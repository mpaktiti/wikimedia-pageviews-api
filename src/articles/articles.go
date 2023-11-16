package articles

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"

	"github.com/mpaktiti/wikimedia-pageviews-api/src/utilities"
)

// Move this to the config file
const baseURL = "https://wikimedia.org/api/rest_v1/metrics/pageviews/top/en.wikipedia/all-access"

type Items struct {
	Items []Item
}

type Item struct {
	Project  string
	Access   string
	Year     string
	Month    string
	Day      string
	Articles []Article
}

type Article struct {
	Article string
	Views   int
	Rank    int
}

func sortMap(input map[string]int) []string {
	// Create slice of key-value pairs
	pairs := make([][2]interface{}, 0, len(input))
	for k, v := range input {
		pairs = append(pairs, [2]interface{}{k, v})
	}

	// Sort slice based on values
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i][1].(int) > pairs[j][1].(int)
	})

	// Extract sorted keys
	keys := make([]string, len(pairs))
	for i, p := range pairs {
		keys[i] = p[0].(string)
	}

	// Print sorted map
	// for _, k := range keys {
	// 	fmt.Printf("%s: %d\n", k, input[k])
	// }

	return keys
}

// curl http://localhost:8080/articles/top/weekly/2023/03
// Returns a list of the most viewed articles for a week
// If an article is not listed on a given day, we assume it has 0 views
// It is assumed that the week starts on Monday
func GetTopArticlesByWeek(year, week string) (string, error) {
	// Convert input year and week to integers
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return "", err
	}
	weekInt, err := strconv.Atoi(week)
	if err != nil {
		return "", err
	}

	// Validate that input year is not out of bounds
	err = utilities.ValidateInputYear(yearInt)
	if err != nil {
		return "", err
	}

	// Validate that input week is not out of bounds
	err = utilities.ValidateInputWeek(year, weekInt)
	if err != nil {
		return "", err
	}

	// Get the first day of the week
	startDate := utilities.WeekStart(yearInt, weekInt)

	var urls [7]string
	var errorStatus, errorDetails string
	articlesMap := map[string]int{}
	for i := 0; i < 7; i++ {
		// Build the URL
		month := utilities.PadString(fmt.Sprint(int(startDate.Month())))
		day := utilities.PadString(fmt.Sprint(startDate.AddDate(0, 0, i).Day()))
		urls[i] = fmt.Sprintf("%s/%s/%s/%s", baseURL, fmt.Sprint(startDate.Year()), month, day)

		// Call the wikipedia API
		response, err := http.Get(urls[i])
		if err != nil {
			return "", err
		}

		// Parse response
		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return "", err
		}

		// If an error happens during any of the API calls stop processing, exit the loop, and return the error details
		if response.StatusCode != http.StatusOK {
			errorStatus = response.Status
			errorDetails, err = utilities.ParseErrorDetails(responseData)
			if err != nil {
				errorDetails = "Failed to process error details"
			}
			break
		}

		// Go through articles and add them to a map
		var items Items
		err = json.Unmarshal(responseData, &items)
		if err != nil {
			return "", err
		}

		//If any items were found extract them from the response and add them to the map
		if len(items.Items) > 0 {
			for _, article := range items.Items[0].Articles {
				if val, ok := articlesMap[article.Article]; ok {
					articlesMap[article.Article] = article.Views + val
				} else {
					articlesMap[article.Article] = article.Views
				}
			}
		}
	}

	if errorDetails != "" {
		return "", fmt.Errorf(errorStatus + ": " + errorDetails)
	}

	// if there are no results return empty result set
	if len(articlesMap) == 0 {
		return "", nil
	}

	// Sort results and retrieve 10 most viewed
	sortKeysDesc := sortMap(articlesMap)
	var top10Articles []Article
	for i := 0; i < 10; i++ {
		top10Articles = append(top10Articles, Article{
			Article: sortKeysDesc[i],
			Views:   articlesMap[sortKeysDesc[i]],
			Rank:    i + 1,
		})
	}

	// Convert to JSON and return string
	jsonResult, err := json.Marshal(top10Articles)
	if err != nil {
		return "", err
	}

	return string(jsonResult), nil
}

// curl http://localhost:8080/articles/top/monthly/2023/03
func GetTopArticlesByMonth(year, month string) (string, error) {
	// Convert input year to integer
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return "", err
	}

	// TODO add test case for this
	// Validate that input year is not out of bounds
	err = utilities.ValidateInputYear(yearInt)
	if err != nil {
		return "", err
	}

	// Build URL
	url := fmt.Sprintf("%s/%s/%s/all-days", baseURL, year, utilities.PadString(month))

	// Call the wikipedia API
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}

	// Parse response
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// Check if the response in anything else than HTTP 200 and return error
	if response.StatusCode != http.StatusOK {
		errorDetails, err := utilities.ParseErrorDetails(responseData)
		if err != nil {
			errorDetails = "Failed to process error details"
		}
		return "", fmt.Errorf(response.Status + ": " + errorDetails)
	}

	// Get top 10 articles
	var items Items
	err = json.Unmarshal(responseData, &items)
	if err != nil {
		return "", err
	}

	// I haven't found a test case where there are no articles returned but I'm adding a check just in case, to avoid nil pointer exceptions
	var top10Articles []Article
	numOfArticles := len(items.Items[0].Articles)
	if numOfArticles > 0 {
		if numOfArticles > 10 {
			numOfArticles = 10
		}
		top10Articles = items.Items[0].Articles[0:numOfArticles]
	}

	// Convert to JSON and return string
	jsonResult, err := json.Marshal(top10Articles)
	if err != nil {
		return "", err
	}

	return string(jsonResult), nil
}
