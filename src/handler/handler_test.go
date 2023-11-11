package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestGETTopArticlesWeekly(t *testing.T) {
	t.Run("returns top articles weekly", func(t *testing.T) {
		// Build the request URL.
		year, week := "2023", "03"
		path := fmt.Sprintf("/articles/top/weekly/%s/%s", year, week)

		// Create a request to pass to the handler.
		req, err := http.NewRequest(http.MethodGet, path, nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()

		// Create a router through which we can pass the request vars.
		router := mux.NewRouter()
		router.HandleFunc("/articles/top/weekly/{year}/{week}", TopArticlesWeeklyHandler)
		router.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		assertStatus(t, rr.Code, http.StatusOK)

		// Check the response body is what we expect.
		expected := fmt.Sprintf(`{"year": %s, "week": %s}`, year, week)
		assertResponseBody(t, rr.Body.String(), expected)
	})
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("handler returned unexpected body: got %v want %v", got, want)
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("handler returned wrong status code: got %v, want %v", got, want)
	}
}
