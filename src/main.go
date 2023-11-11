package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mpaktiti/wikimedia-pageviews-api/internal/handler"
)

func main() {
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	// })
	// log.Println("Listening on localhost:3000")
	// log.Fatal(http.ListenAndServe(":3000", nil))

	// handler := http.HandlerFunc(TopArticlesWeeklyHandler)
	// log.Fatal(http.ListenAndServe(":3000", handler))

	// Create a router and register the available routes.
	// Use regular expressions to validate the variable part of the path.
	r := mux.NewRouter()
	r.HandleFunc("/articles/top/weekly/{year:[0-9]+}/{week:[0-9]+}", handler.TopArticlesWeeklyHandler)
	r.HandleFunc("/articles/top/monthly/{year:[0-9]+}/{month:[0-9]+}", handler.TopArticlesMonthlyHandler)
	r.HandleFunc("/article/{article:[\\w%]+}/weekly/{year:[0-9]+}/{week:[0-9]+}", handler.ViewsPerArticleWeeklyHandler)
	r.HandleFunc("/article/{article}/monthly/{year}/{month}", handler.ViewsPerArticleMonthlyHandler)
	r.HandleFunc("/article/{article}/top/monthly/{year}/{month}", handler.TopViewsPerArticleMonthlyHandler)
	http.Handle("/", r)

	log.Println("Listening on localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
