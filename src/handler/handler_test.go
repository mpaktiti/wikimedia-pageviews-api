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
		assertResponseField(t, "wrong status code", rr.Code, http.StatusOK)

		// Check the response body is what we expect.
		expected := `[{"Article":"Main_Page","Views":35124815,"Rank":1},{"Article":"Index_(statistics)","Views":11321482,"Rank":2},{"Article":"Special:Search","Views":9513645,"Rank":3},{"Article":"The_Last_of_Us_(TV_series)","Views":2502335,"Rank":4},{"Article":"XXX:_Return_of_Xander_Cage","Views":2458723,"Rank":5},{"Article":"Index_(economics)","Views":1577466,"Rank":6},{"Article":"The_Last_of_Us","Views":1540964,"Rank":7},{"Article":"Index,_Washington","Views":1438865,"Rank":8},{"Article":"Wikipedia:Featured_pictures","Views":1415908,"Rank":9},{"Article":"ChatGPT","Views":1329459,"Rank":10}]`
		assertResponseField(t, "unexpected body", rr.Body.String(), expected)
	})
}

func TestGETTopArticlesMonthly(t *testing.T) {
	t.Run("returns top articles monthly", func(t *testing.T) {
		// Build the request URL.
		year, month := "2023", "03"
		path := fmt.Sprintf("/articles/top/monthly/%s/%s", year, month)

		// Create a request to pass to the handler.
		req, err := http.NewRequest(http.MethodGet, path, nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()

		// Create a router through which we can pass the request vars.
		router := mux.NewRouter()
		router.HandleFunc("/articles/top/monthly/{year}/{month}", TopArticlesMonthlyHandler)
		router.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		assertResponseField(t, "wrong status code", rr.Code, http.StatusOK)

		// Check the response body is what we expect.
		expected := `[{"Article":"Main_Page","Views":145431456,"Rank":1},{"Article":"Special:Search","Views":42163260,"Rank":2},{"Article":"YouTube","Views":7716744,"Rank":3},{"Article":"Wikipedia:Featured_pictures","Views":7460936,"Rank":4},{"Article":"ChatGPT","Views":6916888,"Rank":5},{"Article":"Cleopatra","Views":5063272,"Rank":6},{"Article":"Everything_Everywhere_All_at_Once","Views":5061529,"Rank":7},{"Article":"The_Last_of_Us_(TV_series)","Views":4811343,"Rank":8},{"Article":"Deaths_in_2023","Views":4124371,"Rank":9},{"Article":"Lance_Reddick","Views":3937033,"Rank":10}]`
		assertResponseField(t, "unexpected body", rr.Body.String(), expected)
	})
}

func TestGETViewsPerArticleWeekly(t *testing.T) {
	t.Run("returns total pageviews for an article for a specific week", func(t *testing.T) {
		// Build the request URL.
		article, year, week := "Albert_Einstein", "2023", "03"
		path := fmt.Sprintf("/article/%s/weekly/%s/%s", article, year, week)

		// Create a request to pass to the handler.
		req, err := http.NewRequest(http.MethodGet, path, nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()

		// Create a router through which we can pass the request vars.
		router := mux.NewRouter()
		router.HandleFunc("/article/{article}/weekly/{year}/{week}", ViewsPerArticleWeeklyHandler)
		router.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		assertResponseField(t, "wrong status code", rr.Code, http.StatusOK)

		// Check the response body is what we expect.
		expected := `{"Pageviews":"157023"}`
		assertResponseField(t, "unexpected body", rr.Body.String(), expected)
	})
}

func TestGETViewsPerArticleMonthly(t *testing.T) {
	t.Run("returns total pageviews for an article for a specific month", func(t *testing.T) {
		// Build the request URL.
		article, year, month := "Albert_Einstein", "2023", "04"
		path := fmt.Sprintf("/article/%s/monthly/%s/%s", article, year, month)

		// Create a request to pass to the handler.
		req, err := http.NewRequest(http.MethodGet, path, nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()

		// Create a router through which we can pass the request vars.
		router := mux.NewRouter()
		router.HandleFunc("/article/{article}/monthly/{year}/{month}", ViewsPerArticleMonthlyHandler)
		router.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		assertResponseField(t, "wrong status code", rr.Code, http.StatusOK)

		// Check the response body is what we expect.
		expected := `{"Pageviews":"485684"}`
		assertResponseField(t, "unexpected body", rr.Body.String(), expected)
	})
}

func TestGETDayWithMostPageviews(t *testing.T) {
	t.Run("returns the day with the most pageviews for an article for a specific month", func(t *testing.T) {
		// Build the request URL.
		article, year, month := "Albert_Einstein", "2023", "04"
		path := fmt.Sprintf("/article/%s/top/monthly/%s/%s", article, year, month)

		// Create a request to pass to the handler.
		req, err := http.NewRequest(http.MethodGet, path, nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()

		// Create a router through which we can pass the request vars.
		router := mux.NewRouter()
		router.HandleFunc("/article/{article}/top/monthly/{year}/{month}", TopViewsPerArticleMonthlyHandler)
		router.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		assertResponseField(t, "wrong status code", rr.Code, http.StatusOK)

		// Check the response body is what we expect.
		expected := `{"Pageviews":"30724","Timestamp":"2023042200"}`
		assertResponseField(t, "unexpected body", rr.Body.String(), expected)
	})
}

// use interface{} for input so it can be either string or int
func assertResponseField(t testing.TB, fieldAsserted string, got, want interface{}) {
	t.Helper()
	if got != want {
		t.Errorf("code returned %v: got %v want %v", fieldAsserted, got, want)
	}
}
