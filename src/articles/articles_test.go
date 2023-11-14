package articles

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSortMap(t *testing.T) {
	testCases := []struct {
		name        string
		unsortedMap map[string]int
		expectedMap []string
	}{
		{
			name: "sort map with 5 items",
			unsortedMap: map[string]int{
				"article1": 711,
				"article2": 2138,
				"article3": 1908,
				"article4": 912,
			},
			expectedMap: []string{"article2", "article3", "article4", "article1"},
		},
		{
			name: "sort map with 10 items",
			unsortedMap: map[string]int{
				"article1":  1,
				"article2":  2,
				"article3":  3,
				"article4":  4,
				"article5":  5,
				"article6":  6,
				"article7":  7,
				"article8":  8,
				"article9":  9,
				"article10": 10,
			},
			expectedMap: []string{"article10", "article9", "article8", "article7", "article6", "article5", "article4", "article3", "article2", "article1"},
		},
		{
			name:        "sort empty map",
			unsortedMap: map[string]int{},
			expectedMap: []string{},
		},
	}
	for tcNum, tc := range testCases {
		gotMap := sortMap(tc.unsortedMap)
		if !reflect.DeepEqual(gotMap, tc.expectedMap) {
			t.Errorf("test %d failed: got %v want %v", tcNum, gotMap, tc.expectedMap)
		}
	}
}

func TestGetTopArticlesByMonth(t *testing.T) {
	testCases := []struct {
		name             string
		year             string
		month            string
		expectedArticles string
		expectedError    string
	}{
		{
			name:             "top 10 most viewed articles on the 3rd month of 2023",
			year:             "2023",
			month:            "03",
			expectedArticles: `[{"Article":"Main_Page","Views":145431456,"Rank":1},{"Article":"Special:Search","Views":42163260,"Rank":2},{"Article":"YouTube","Views":7716744,"Rank":3},{"Article":"Wikipedia:Featured_pictures","Views":7460936,"Rank":4},{"Article":"ChatGPT","Views":6916888,"Rank":5},{"Article":"Cleopatra","Views":5063272,"Rank":6},{"Article":"Everything_Everywhere_All_at_Once","Views":5061529,"Rank":7},{"Article":"The_Last_of_Us_(TV_series)","Views":4811343,"Rank":8},{"Article":"Deaths_in_2023","Views":4124371,"Rank":9},{"Article":"Lance_Reddick","Views":3937033,"Rank":10}]`,
			expectedError:    "",
		},
		{
			name:             "error case: HTTP 400 for invalid input (month > 12)",
			year:             "2023",
			month:            "13",
			expectedArticles: "",
			expectedError:    "400 Bad Request: Given year/month/day is invalid date",
		},
		{
			name:             "error case: HTTP 404 for invalid input (future date)",
			year:             "2030",
			month:            "03",
			expectedArticles: "",
			expectedError:    "404 Not Found: The date(s) you used are valid, but we either do not have data for those date(s), or the project you asked for is not loaded yet. Please check documentation for more information.",
		},
	}
	for i, tc := range testCases {
		gotArticles, gotError := GetTopArticlesByMonth(tc.year, tc.month)
		if tc.expectedError != "" {
			require.Error(t, gotError)
			assertResponseField(t, i, gotError.Error(), tc.expectedError)
		} else {
			require.NoError(t, gotError)
		}
		assertResponseField(t, i, gotArticles, tc.expectedArticles)
	}
}

func TestGetTopArticlesByWeek(t *testing.T) {
	testCases := []struct {
		name             string
		year             string
		week             string
		expectedArticles string
		expectedError    string
	}{
		{
			name:             "top 10 most viewed articles on the 3rd week of 2023",
			year:             "2023",
			week:             "03",
			expectedArticles: `[{"Article":"Main_Page","Views":35124815,"Rank":1},{"Article":"Index_(statistics)","Views":11321482,"Rank":2},{"Article":"Special:Search","Views":9513645,"Rank":3},{"Article":"The_Last_of_Us_(TV_series)","Views":2502335,"Rank":4},{"Article":"XXX:_Return_of_Xander_Cage","Views":2458723,"Rank":5},{"Article":"Index_(economics)","Views":1577466,"Rank":6},{"Article":"The_Last_of_Us","Views":1540964,"Rank":7},{"Article":"Index,_Washington","Views":1438865,"Rank":8},{"Article":"Wikipedia:Featured_pictures","Views":1415908,"Rank":9},{"Article":"ChatGPT","Views":1329459,"Rank":10}]`,
			expectedError:    "",
		},
		{
			name:             "error case: HTTP 404 for invalid input (week > 52)",
			year:             "2023",
			week:             "55",
			expectedArticles: "",
			expectedError:    "404 Not Found: The date(s) you used are valid, but we either do not have data for those date(s), or the project you asked for is not loaded yet. Please check documentation for more information.",
		},
		{
			name:             "error case: HTTP 404 for invalid input (future date)",
			year:             "2030",
			week:             "03",
			expectedArticles: "",
			expectedError:    "404 Not Found: The date(s) you used are valid, but we either do not have data for those date(s), or the project you asked for is not loaded yet. Please check documentation for more information.",
		},
	}
	for i, tc := range testCases {
		gotArticles, gotError := GetTopArticlesByWeek(tc.year, tc.week)
		if tc.expectedError != "" {
			require.Error(t, gotError)
			assertResponseField(t, i, gotError.Error(), tc.expectedError)
		} else {
			require.NoError(t, gotError)
		}
		assertResponseField(t, i, gotArticles, tc.expectedArticles)
	}
}

func assertResponseField(t testing.TB, testNum int, got, want interface{}) {
	t.Helper()
	if got != want {
		t.Errorf("test %d failed: got %v want %v", testNum+1, got, want)
	}
}
