package utils

import "time"

// FormatTimestamp converts an ISO 8601 timestamp string to a human-readable date
// If the timestamp cannot be parsed, it returns the original string
func FormatTimestamp(timestamp string) string {
	if timestamp == "" {
		return ""
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
