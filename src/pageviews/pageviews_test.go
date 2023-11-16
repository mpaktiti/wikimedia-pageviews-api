package pageviews

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPageviewsByWeek(t *testing.T) {
	testCases := []struct {
		name              string
		article           string
		year              string
		week              string
		expectedPageviews int
		expectedError     string
	}{
		{
			name:              "total pageviews for Albert Einstein article on the 3rd week of 2023",
			article:           "Albert_Einstein",
			year:              "2023",
			week:              "03",
			expectedPageviews: 157023,
			expectedError:     "",
		},
		{
			name:              "total pageviews for Albert Einstein article on the 1st week of 2020 (which starts in 2019)",
			article:           "Albert_Einstein",
			year:              "2023",
			week:              "01",
			expectedPageviews: 170511,
			expectedError:     "",
		},
		{
			name:              "total pageviews for Albert Einstein article on the last week of 2020 (which ends in 2021)",
			article:           "Albert_Einstein",
			year:              "2020",
			week:              "53",
			expectedPageviews: 109207,
			expectedError:     "",
		},
		{
			name:              "error case: HTTP 400 for invalid input (week)",
			article:           "Albert_Einstein",
			year:              "2020",
			week:              "100",
			expectedPageviews: 0,
			expectedError:     "400 Bad Request: input week cannot be greater than 53",
		},
		{
			name:              "error case: HTTP 400 for invalid input (year)",
			article:           "Albert_Einstein",
			year:              "2030",
			week:              "01",
			expectedPageviews: 0,
			expectedError:     "400 Bad Request: input year cannot be greater than current year",
		},
		{
			name:              "error case: HTTP 404 when the article does not exist",
			article:           "JHKJHK123",
			year:              "2023",
			week:              "03",
			expectedPageviews: 0,
			expectedError:     "404 Not Found: The date(s) you used are valid, but we either do not have data for those date(s), or the project you asked for is not loaded yet. Please check documentation for more information.",
		},
	}
	for i, tc := range testCases {
		gotPageviews, gotError := GetPageviewsByWeek(tc.article, tc.year, tc.week)
		if tc.expectedError != "" {
			require.Error(t, gotError)
			assertResponseField(t, i, gotError.Error(), tc.expectedError)
		} else {
			require.NoError(t, gotError)
		}
		assertResponseField(t, i, gotPageviews, tc.expectedPageviews)
	}
}

func TestGetPageviewsByMonth(t *testing.T) {
	testCases := []struct {
		name              string
		article           string
		year              string
		month             string
		expectedPageviews int
		expectedError     string
	}{
		{
			name:              "total pageviews for Albert Einstein article on the 3rd week of 2023",
			article:           "Albert_Einstein",
			year:              "2023",
			month:             "04",
			expectedPageviews: 485684,
			expectedError:     "",
		},
		{
			name:              "error case: HTTP 400 for invalid input",
			article:           "Albert_Einstein",
			year:              "2023",
			month:             "14",
			expectedPageviews: 0,
			expectedError:     "400 Bad Request: start timestamp is invalid, must be a valid date in YYYYMMDD format",
		},
		{
			name:              "error case: HTTP 404 when the article does not exist",
			article:           "JHKJHK123",
			year:              "2023",
			month:             "04",
			expectedPageviews: 0,
			expectedError:     "404 Not Found: The date(s) you used are valid, but we either do not have data for those date(s), or the project you asked for is not loaded yet. Please check documentation for more information.",
		},
	}
	for i, tc := range testCases {
		gotPageviews, gotError := GetPageviewsByMonth(tc.article, tc.year, tc.month)
		if tc.expectedError != "" {
			require.Error(t, gotError)
			assertResponseField(t, i, gotError.Error(), tc.expectedError)
		} else {
			require.NoError(t, gotError)
		}
		assertResponseField(t, i, gotPageviews, tc.expectedPageviews)
	}
}

func TestGetDayWithMostPageviews(t *testing.T) {
	testCases := []struct {
		name              string
		article           string
		year              string
		month             string
		expectedDay       string
		expectedPageviews int
		expectedError     string
	}{
		{
			name:              "total pageviews for Albert Einstein article on the 3rd week of 2023",
			article:           "Albert_Einstein",
			year:              "2023",
			month:             "04",
			expectedDay:       "2023042200",
			expectedPageviews: 30724,
			expectedError:     "",
		},
		{
			name:              "error case: HTTP 400 for invalid input",
			article:           "Albert_Einstein",
			year:              "2023",
			month:             "14",
			expectedDay:       "",
			expectedPageviews: 0,
			expectedError:     "400 Bad Request: start timestamp is invalid, must be a valid date in YYYYMMDD format",
		},
		{
			name:              "error case: HTTP 404 when the article does not exist",
			article:           "JHKJHK123",
			year:              "2023",
			month:             "04",
			expectedDay:       "",
			expectedPageviews: 0,
			expectedError:     "404 Not Found: The date(s) you used are valid, but we either do not have data for those date(s), or the project you asked for is not loaded yet. Please check documentation for more information.",
		},
	}
	for i, tc := range testCases {
		gotDay, gotPageviews, gotError := GetDayWithMostPageviews(tc.article, tc.year, tc.month)
		if tc.expectedError != "" {
			require.Error(t, gotError)
			assertResponseField(t, i, gotError.Error(), tc.expectedError)
		} else {
			require.NoError(t, gotError)
		}
		assertResponseField(t, i, gotDay, tc.expectedDay)
		assertResponseField(t, i, gotPageviews, tc.expectedPageviews)
	}
}

func assertResponseField(t testing.TB, testNum int, got, want interface{}) {
	t.Helper()
	if got != want {
		t.Errorf("test %d failed: got %v want %v", testNum+1, got, want)
	}
}
