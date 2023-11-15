package articles

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
		// TODO log error properly
		fmt.Println("Error during conversion")
		return "", err
	}
	weekInt, err := strconv.Atoi(week)
	if err != nil {
		// TODO log error properly
		fmt.Println("Error during conversion")
		// TODO return HTTP code for bad input
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
			// TODO log error properly
			fmt.Print(err.Error())
			return "", err
		}

		// Parse response
		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			//TODO why do I log.Fatal() here and fmt.Print() earlier?
			log.Fatal(err)
			return "", err
		}

		// If an error happens during any of the API calls stop processing, exit the loop, and return the error details
		if response.StatusCode != http.StatusOK {
			fmt.Println("Response status: ", response.StatusCode)
			fmt.Println(string(responseData))
			errorStatus = response.Status
			errorDetails, err = utilities.ParseErrorDetails(responseData)
			if err != nil {
				fmt.Print(err.Error())
				errorDetails = "Failed to process error details"
			}
			break
		}

		// Go through articles and add them to a map
		var items Items
		err = json.Unmarshal(responseData, &items)
		if err != nil {
			fmt.Println("error:", err)
			return "", err
		}

		//If any items were found extract them from the response and add them to the map
		// TODO is this if needed? is there any case where I get HTTP 200 but no articles?
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
		fmt.Println(err)
		return "", err
	}

	return string(jsonResult), nil
}

// curl http://localhost:8080/articles/top/monthly/2023/03
func GetTopArticlesByMonth(year, month string) (string, error) {
	// Build URL
	url := fmt.Sprintf("%s/%s/%s/all-days", baseURL, year, month)

	// Call the wikipedia API
	response, err := http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
		return "", err
	}

	// Parse response
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	// fmt.Println(string(responseData))

	// Check if the response in anything else than HTTP 200 and return error
	if response.StatusCode != http.StatusOK {
		fmt.Println("Response status: ", response.StatusCode)
		fmt.Println(string(responseData))
		errorDetails, err := utilities.ParseErrorDetails(responseData)
		if err != nil {
			fmt.Print(err.Error())
			errorDetails = "Failed to process error details"
		}
		return "", fmt.Errorf(response.Status + ": " + errorDetails)
	}

	// Get top 10 articles
	var items Items
	err = json.Unmarshal(responseData, &items)
	if err != nil {
		fmt.Println("error:", err)
		return "", err
	}

	// TODO should I check here if len(items.Items) > 0 ??
	top10Articles := items.Items[0].Articles[0:10]

	// fmt.Printf("%+v", top10Articles)

	// Convert to JSON and return string
	jsonResult, err := json.Marshal(top10Articles)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return string(jsonResult), nil
}
