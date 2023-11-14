package converters

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertPageviewsToJson(t *testing.T) {
	t.Run("convert input to JSON", func(t *testing.T) {
		pageviews := 485684
		want := []byte(`{"Pageviews":"485684"}`)
		got, err := ConvertPageviewsToJson(pageviews)
		require.NoError(t, err)
		assertJSON(t, got, want)
	})
}

func TestConvertTopDayPageviewsToJson(t *testing.T) {
	t.Run("convert input to JSON", func(t *testing.T) {
		pageviews := 30724
		timestamp := "2023042200"
		want := []byte(`{"Pageviews":"30724","Timestamp":"2023042200"}`)
		got, err := ConvertTopDayPageviewsToJson(timestamp, pageviews)
		require.NoError(t, err)
		assertJSON(t, got, want)
	})
}

func TestErrorToJson(t *testing.T) {
	t.Run("convert error to JSON", func(t *testing.T) {
		error := "400 Bad Request: start timestamp is invalid, must be a valid date in YYYYMMDD format"
		want := []byte(`{"Error":"400 Bad Request: start timestamp is invalid, must be a valid date in YYYYMMDD format"}`)
		got, err := ConvertErrorToJson(error)
		require.NoError(t, err)
		assertJSON(t, got, want)
	})
}

func assertJSON(t testing.TB, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
