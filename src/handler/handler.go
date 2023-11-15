package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mpaktiti/wikimedia-pageviews-api/src/articles"
	"github.com/mpaktiti/wikimedia-pageviews-api/src/converters"
	"github.com/mpaktiti/wikimedia-pageviews-api/src/pageviews"
)

func parseStatusCode(input string) int {
	// Enhancement: don't count only on converting the first 3 characters to integer, verify the result against an enum
	statusCode, conversionErr := strconv.Atoi(input)
	if conversionErr != nil {
		// if you cannot parse the status returned from the API, it must be an internal error, return 500
		statusCode = http.StatusInternalServerError
	}
	return statusCode
}

func convertError(input string) ([]byte, error) {
	// Enhancement: the err.Error() contains the HTTP status at the beginning, parse the string so it returns only the error description
	res, err := converters.ConvertErrorToJson(input)
	if err != nil {
		// If the conversion fails create the JSON here
		// Enhancement: as is this will return the conversion error but now the API error has been lost. This could aggregate errors.
		error := fmt.Sprintf(`{"Error": "%s"}`, err.Error())
		return []byte(error), err
	}
	return res, nil
}

func TopArticlesWeeklyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	res, err := articles.GetTopArticlesByWeek(vars["year"], vars["week"])
	if err != nil {
		// TODO enhancement: properly log errors
		fmt.Println("ERROR: ", err)
		statusCode, conversionErr := strconv.Atoi(err.Error()[:3])
		if conversionErr != nil {
			fmt.Println("ERROR: ", err)
			// if you cannot parse the status returned from the API, return 500
			statusCode = http.StatusInternalServerError
		}
		w.WriteHeader(statusCode)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

func TopArticlesMonthlyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	res, err := articles.GetTopArticlesByMonth(vars["year"], vars["month"])
	if err != nil {
		// TODO enhancement: properly log errors
		fmt.Println("ERROR: ", err)
		statusCode, conversionErr := strconv.Atoi(err.Error()[:3])
		if conversionErr != nil {
			fmt.Println("ERROR: ", err)
			// if you cannot parse the status returned from the API, return 500
			statusCode = http.StatusInternalServerError
		}
		w.WriteHeader(statusCode)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// TODO maybe just get directly []bytes output from GetTopArticlesByMonth?
	w.Write([]byte(res))
}

func ViewsPerArticleWeeklyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageviews, err := pageviews.GetPageviewsByWeek(vars["article"], vars["year"], vars["week"])
	if err != nil {
		fmt.Println("ERROR: ", err)
		// Parse returned error to see if there is an HTTP status there
		statusCode := parseStatusCode(err.Error()[:3])
		w.WriteHeader(statusCode)

		// Convert returned error to JSON
		res, err := convertError(err.Error())
		if err != nil {
			fmt.Println("ERROR: ", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write(res)
		return
	}

	// Convert pageviews result to JSON
	res, err := converters.ConvertPageviewsToJson(pageviews)
	if err != nil {
		fmt.Println("ERROR: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		res, _ := convertError(err.Error())
		w.Write([]byte(res))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func ViewsPerArticleMonthlyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageviews, err := pageviews.GetPageviewsByMonth(vars["article"], vars["year"], vars["month"])
	if err != nil {
		fmt.Println("ERROR: ", err)
		// Parse returned error to see if there is an HTTP status there
		statusCode := parseStatusCode(err.Error()[:3])
		w.WriteHeader(statusCode)

		// Convert returned error to JSON
		res, err := convertError(err.Error())
		if err != nil {
			fmt.Println("ERROR: ", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write(res)
		return
	}

	// Convert pageviews result to JSON
	res, err := converters.ConvertPageviewsToJson(pageviews)
	if err != nil {
		fmt.Println("ERROR: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		res, _ := convertError(err.Error())
		w.Write([]byte(res))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func TopViewsPerArticleMonthlyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	timestamp, pageviews, err := pageviews.GetDayWithMostPageviews(vars["article"], vars["year"], vars["month"])
	if err != nil {
		fmt.Println("ERROR: ", err)
		// Parse returned error to see if there is an HTTP status there
		statusCode := parseStatusCode(err.Error()[:3])
		w.WriteHeader(statusCode)

		// Convert returned error to JSON
		res, err := convertError(err.Error())
		if err != nil {
			fmt.Println("ERROR: ", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write(res)
		return
	}

	// Convert pageviews result to JSON
	res, err := converters.ConvertTopDayPageviewsToJson(timestamp, pageviews)
	if err != nil {
		fmt.Println("ERROR: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		res, _ := convertError(err.Error())
		w.Write([]byte(res))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
