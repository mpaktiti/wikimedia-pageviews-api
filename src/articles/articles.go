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
	for _, k := range keys {
		fmt.Printf("%s: %d\n", k, input[k])
	}

	return keys
}

// curl http://localhost:3000/articles/top/weekly/2023/03
// Returns a list of the most viewed articles for a week
// If an article is not listed on a given day, we assume it has 0 views
func GetTopArticlesByWeek(year, week string) {
	// Convert input year and week to integers
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
	fmt.Println("Start Date: ", startDate.Day())
	fmt.Println("End Date + 6: ", startDate.AddDate(0, 0, 6).Day())
	fmt.Println("Month as a number: ", int(startDate.Month()))

	var urls [7]string
	var articles []Article
	articlesMap := map[string]int{}
	for i := 0; i < 7; i++ {
		// 1. Build URLs
		// Wikipedia API expects months and days as 2 digits each so add a zero at the beginning if needed
		month := fmt.Sprint(int(startDate.Month()))
		if len(month) == 1 {
			month = "0" + month
		}
		day := fmt.Sprint(startDate.AddDate(0, 0, i).Day())
		if len(day) == 1 {
			day = "0" + day
		}
		urls[i] = fmt.Sprintf("%s/%s/%s/%s", baseURL, fmt.Sprint(startDate.Year()), month, day)
		// log.Println("URL: ", urls[i])

		// 2. Call the wikipedia API
		response, err := http.Get(urls[i])
		if err != nil {
			// TODO log error properly
			fmt.Print(err.Error())
		}

		// 3. Parse response
		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		// 4. Go through articles and add them to a struct?
		var items Items
		err = json.Unmarshal(responseData, &items)
		if err != nil {
			fmt.Println("error:", err)
		}
		// fmt.Printf("Call %+v: %+v results\n", i, len(items.Items[0].Articles))
		articles = append(articles, items.Items[0].Articles...)

		for _, article := range items.Items[0].Articles {
			if val, ok := articlesMap[article.Article]; ok {
				articlesMap[article.Article] = article.Views + val
			} else {
				articlesMap[article.Article] = article.Views
			}
		}
	}

	fmt.Printf("Total results: %+v\n", len(articles))
	fmt.Printf("map: %+v\n", articlesMap)
	fmt.Printf("map len: %+v\n", len(articlesMap))

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
	// fmt.Printf("TOP 10 ARTICLES: %+v", top10Articles)

	b, err := json.Marshal(top10Articles)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	fmt.Println("TOP 10 ARTICLES: ", string(b))
}

// curl http://localhost:3000/articles/top/monthly/2023/03
func GetTopArticlesByMonth(year, month string) {
	// Build URL
	url := fmt.Sprintf("%s/%s/%s/all-days", baseURL, year, month)

	// Call the wikipedia API
	response, err := http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
	}

	// Parse response
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(string(responseData))

	// Get top 10 articles
	var items Items
	err = json.Unmarshal(responseData, &items)
	if err != nil {
		fmt.Println("error:", err)
	}
	top10Articles := items.Items[0].Articles[0:9]
	fmt.Printf("%+v", top10Articles)

	// b, err := json.Marshal(top10Articles)
	// if err != nil {
	// 	fmt.Printf("Error: %s", err)
	// 	return
	// }
	// fmt.Println(string(b))
}
