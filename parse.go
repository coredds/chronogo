package chronogo

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	// Common datetime patterns
	iso8601Pattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})$`)

	// Common layouts for parsing
	commonLayouts = []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"15:04:05",
		"15:04",
	}
)

// ParseError represents an error that occurred during parsing.
type ParseError struct {
	Input  string
	Reason string
}

func (e ParseError) Error() string {
	return fmt.Sprintf("failed to parse '%s': %s", e.Input, e.Reason)
}

// Parse parses a datetime string using common formats.
// It supports ISO 8601, RFC 3339, and other common datetime formats.
func Parse(value string) (DateTime, error) {
	return ParseInLocation(value, time.UTC)
}

// ParseInLocation parses a datetime string in the specified location.
func ParseInLocation(value string, loc *time.Location) (DateTime, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return DateTime{}, ParseError{Input: value, Reason: "empty string"}
	}

	// Try each common layout
	for _, layout := range commonLayouts {
		if t, err := time.ParseInLocation(layout, value, loc); err == nil {
			return DateTime{t}, nil
		}
	}

	// Try without timezone info, defaulting to provided location
	if t, err := time.ParseInLocation("2006-01-02T15:04:05", value, loc); err == nil {
		return DateTime{t}, nil
	}

	return DateTime{}, ParseError{
		Input:  value,
		Reason: "no matching format found",
	}
}

// ParseISO8601 parses an ISO 8601 formatted datetime string.
func ParseISO8601(value string) (DateTime, error) {
	if !iso8601Pattern.MatchString(value) {
		return DateTime{}, ParseError{Input: value, Reason: "invalid ISO 8601 format"}
	}

	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		// Try with nanoseconds
		t, err = time.Parse(time.RFC3339Nano, value)
		if err != nil {
			return DateTime{}, ParseError{Input: value, Reason: err.Error()}
		}
	}

	return DateTime{t}, nil
}

// ParseRFC3339 parses an RFC 3339 formatted datetime string.
func ParseRFC3339(value string) (DateTime, error) {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return DateTime{}, ParseError{Input: value, Reason: err.Error()}
	}
	return DateTime{t}, nil
}

// FromFormat parses a datetime string using a custom format layout.
// This is similar to Pendulum's from_format method but uses Go's time format syntax.
func FromFormat(value, layout string) (DateTime, error) {
	return FromFormatInLocation(value, layout, time.UTC)
}

// FromFormatInLocation parses a datetime string using a custom format layout in the specified location.
func FromFormatInLocation(value, layout string, loc *time.Location) (DateTime, error) {
	t, err := time.ParseInLocation(layout, value, loc)
	if err != nil {
		return DateTime{}, ParseError{Input: value, Reason: err.Error()}
	}
	return DateTime{t}, nil
}

// LoadLocation loads a timezone by name.
// This is a convenience wrapper around time.LoadLocation.
func LoadLocation(name string) (*time.Location, error) {
	if name == "local" {
		return time.Local, nil
	}
	return time.LoadLocation(name)
}

// MustLoadLocation loads a timezone by name and panics if it fails.
func MustLoadLocation(name string) *time.Location {
	loc, err := LoadLocation(name)
	if err != nil {
		panic(fmt.Sprintf("failed to load location '%s': %v", name, err))
	}
	return loc
}

// Instance creates a DateTime from a standard time.Time.
func Instance(t time.Time) DateTime {
	return DateTime{t}
}

// parseUnixTimestamp attempts to parse a Unix timestamp string.
func parseUnixTimestamp(value string) (DateTime, error) {
	s := strings.TrimSpace(value)
	if s == "" {
		return DateTime{}, errors.New("invalid Unix timestamp")
	}

	// Detect length ignoring leading sign
	signless := s
	if s[0] == '-' || s[0] == '+' {
		signless = s[1:]
	}

	// Parse as int64
	ts, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return DateTime{}, errors.New("invalid Unix timestamp")
	}

	switch l := len(signless); l {
	case 1, 2, 3, 4, 5, 6, 7, 8, 9, 10: // seconds (support small absolute values too)
		return DateTime{time.Unix(ts, 0).UTC()}, nil
	case 13: // milliseconds
		return DateTime{time.UnixMilli(ts).UTC()}, nil
	case 16: // microseconds
		return DateTime{time.UnixMicro(ts).UTC()}, nil
	case 19: // nanoseconds
		return DateTime{time.Unix(0, ts).UTC()}, nil
	default:
		return DateTime{}, errors.New("invalid Unix timestamp length")
	}
}

// TryParseUnix attempts to parse a string as a Unix timestamp.
func TryParseUnix(value string) (DateTime, error) {
	return parseUnixTimestamp(value)
}
