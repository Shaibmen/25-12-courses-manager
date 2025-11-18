package utils

import "time"

func TimeParse(date string) (*time.Time, error) {
	result, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
