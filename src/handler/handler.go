package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mpaktiti/wikimedia-pageviews-api/internal/articles"
	"github.com/mpaktiti/wikimedia-pageviews-api/internal/pageviews"
)

func TopArticlesWeeklyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// w.WriteHeader(http.StatusNotFound)
	// TODO expect response (or error) from articles function
	articles.GetTopArticlesByWeek(vars["year"], vars["week"])
	fmt.Fprintf(w, `{"year": %s, "week": %s}`, vars["year"], vars["week"])
}

func TopArticlesMonthlyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Year: %v\n", vars["year"])
	fmt.Fprintf(w, "Month: %v\n", vars["month"])
	// TODO expect response (or error) from articles function
	articles.GetTopArticlesByMonth(vars["year"], vars["month"])
}

func ViewsPerArticleWeeklyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	res, err := pageviews.GetPageviewsByWeek(vars["article"], vars["year"], vars["week"])
	if err != nil {
		// TODO properly log error and return failure http code (which one?)
		fmt.Println("error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Total pageviews: %v\n", res)
}

func ViewsPerArticleMonthlyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	res, err := pageviews.GetPageviewsByMonth(vars["article"], vars["year"], vars["month"])
	if err != nil {
		// TODO properly log error and return failure http code (which one?)
		fmt.Println("error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Total pageviews: %v\n", res)
}

func TopViewsPerArticleMonthlyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resDay, resViews, err := pageviews.GetDayWithMostPageviews(vars["article"], vars["year"], vars["month"])
	if err != nil {
		// TODO properly log error and return failure http code (which one?)
		fmt.Println("error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Top day is %v with %v pageviews\n", resDay, resViews)
}
