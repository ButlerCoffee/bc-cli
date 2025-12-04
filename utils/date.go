package utils

import (
	"strconv"
	"time"
)

// FormatTimestamp converts an ISO 8601 timestamp string or Unix milliseconds to a human-readable date
// If the timestamp cannot be parsed, it returns the original string
func FormatTimestamp(timestamp string) string {
	if timestamp == "" {
		return ""
	}

	// Try parsing as Unix timestamp in milliseconds first (e.g., "1764427190000")
	if unixMs, err := strconv.ParseInt(timestamp, 10, 64); err == nil {
		// Check if it's a reasonable Unix timestamp (after year 2000 and before year 2100)
		// Unix ms for 2000-01-01: 946684800000
		// Unix ms for 2100-01-01: 4102444800000
		if unixMs > 946684800000 && unixMs < 4102444800000 {
			t := time.Unix(unixMs/1000, (unixMs%1000)*1000000)
			return t.Format("January 2, 2006 at 3:04 PM")
		}
	}

	// Try common timestamp formats
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05.999999999Z07:00",
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timestamp); err == nil {
			// Format as a human-readable date
			return t.Format("January 2, 2006 at 3:04 PM")
		}
	}

	// If we can't parse it, return the original
	return timestamp
}

// ParseTimestamp converts an ISO 8601 timestamp string or Unix milliseconds to a time.Time object
func ParseTimestamp(timestamp string) (time.Time, error) {
	if timestamp == "" {
		return time.Time{}, strconv.ErrSyntax
	}

	// Try parsing as Unix timestamp in milliseconds first
	if unixMs, err := strconv.ParseInt(timestamp, 10, 64); err == nil {
		if unixMs > 946684800000 && unixMs < 4102444800000 {
			return time.Unix(unixMs/1000, (unixMs%1000)*1000000), nil
		}
	}

	// Try common timestamp formats
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05.999999999Z07:00",
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timestamp); err == nil {
			return t, nil
		}
	}

	return time.Time{}, strconv.ErrSyntax
}

// FormatDate formats a time.Time as a date-only string (no time component)
func FormatDate(t time.Time) string {
	return t.Format("January 2, 2006")
}
