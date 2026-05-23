package timeparser

import (
	"fmt"
	"time"
)

// Common log timestamp formats to attempt when parsing.
var knownFormats = []string{
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02T15:04:05",
	"2006-01-02T15:04:05.999999999",
	"2006-01-02 15:04:05",
	"2006-01-02 15:04:05.999999999",
	"2006/01/02 15:04:05",
	"02/Jan/2006:15:04:05 -0700",
	"Jan 02 15:04:05",
}

// Parse attempts to parse a timestamp string using a list of known formats.
// It returns the parsed time and the matched format string, or an error if none match.
func Parse(value string) (time.Time, string, error) {
	for _, layout := range knownFormats {
		t, err := time.Parse(layout, value)
		if err == nil {
			return t, layout, nil
		}
	}
	return time.Time{}, "", fmt.Errorf("timeparser: unrecognized timestamp format: %q", value)
}

// ParseWithLocation parses a timestamp using known formats and applies the given location.
func ParseWithLocation(value string, loc *time.Location) (time.Time, string, error) {
	for _, layout := range knownFormats {
		t, err := time.ParseInLocation(layout, value, loc)
		if err == nil {
			return t, layout, nil
		}
	}
	return time.Time{}, "", fmt.Errorf("timeparser: unrecognized timestamp format %q in location %s", value, loc)
}

// KnownFormats returns a copy of the list of supported timestamp formats.
func KnownFormats() []string {
	copy := make([]string, len(knownFormats))
	for i, f := range knownFormats {
		copy[i] = f
	}
	return copy
}
