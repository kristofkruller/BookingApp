package helpers

import "time"

// GENERAL HELPERS
func ParseTime(dateStr string) (time.Time, error) {
	layout := "2006-01-02" // LAYOUT - datepicker string
	parsedTime, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, err // Zero time and err
	}
	return parsedTime, nil
}

func IsValidDateRange(start, end time.Time) bool {
	return start.Before(end) || start.Equal(end)
}

func IsValidPriceLogic(min, max float64) bool {
	return min <= max
}
