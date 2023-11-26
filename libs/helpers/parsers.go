package helpers

import "time"

func ParseTime(dateStr string) (time.Time, error) {
	layout := "2006-01-02" // LAYOUT - datepicker string
	parsedTime, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, err // Zero time and err
	}
	return parsedTime, nil
}
