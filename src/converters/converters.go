package converters

import (
	"encoding/json"
	"fmt"
)

type Pageviews struct {
	Pageviews string
}

type TopDayPageviews struct {
	Pageviews string
	Timestamp string
}

type Error struct {
	Error string
}

// TODO write tests
func ConvertPageviewsToJson(input int) ([]byte, error) {
	pageviews := &Pageviews{Pageviews: fmt.Sprint(input)}
	res, err := json.Marshal(pageviews)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// TODO write tests
func ConvertTopDayPageviewsToJson(timestamp string, pageviews int) ([]byte, error) {
	topDayPageviews := &TopDayPageviews{
		Pageviews: fmt.Sprint(pageviews),
		Timestamp: timestamp,
	}
	res, err := json.Marshal(topDayPageviews)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// TODO write tests
func ConvertErrorToJson(input string) ([]byte, error) {
	error := &Error{Error: input}
	res, err := json.Marshal(error)
	if err != nil {
		return nil, err
	}
	return res, nil
}
