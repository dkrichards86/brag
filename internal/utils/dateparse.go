package utils

import (
	"time"
)

// SupportedDateFormats are the date formats accepted by brag
var SupportedDateFormats = []string{
	"2006/01/02",
	"2006-01-02",
	"01/02/2006",
	"01-02-2006",
}

// ParseDate attempts to parse a date string using supported formats
// Returns the parsed time and true if successful, otherwise zero time and false
func ParseDate(dateStr string) (time.Time, bool) {
	for _, format := range SupportedDateFormats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}
