package chronogo

import (
	"errors"
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
		// strict RFC/ISO
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z07:00",
		// space and naive
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
		// lenient separators and compact
		"2006/01/02 15:04:05",
		"2006/01/02",
		"2006-1-2 15:04:05",
		"2006-1-2",
		"20060102",
		// time-only
		"15:04:05",
		"15:04",
	}
)

// Parse parses a datetime string using common formats (lenient by default) in UTC.
func Parse(value string) (DateTime, error) {
	return ParseInLocation(value, time.UTC)
}

// ParseInLocation parses a datetime string in the specified location.
// ParseInLocation parses a datetime string in the specified location using a lenient set of formats.
func ParseInLocation(value string, loc *time.Location) (DateTime, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return DateTime{}, ParseError(value, errors.New("empty string"))
	}

	// Try each common layout (lenient set)
	for _, layout := range commonLayouts {
		if t, err := time.ParseInLocation(layout, value, loc); err == nil {
			return DateTime{t}, nil
		}
	}

	return DateTime{}, ParseError(value, errors.New("no matching format found"))
}

// ParseStrict parses using a stricter set of layouts (RFC3339 and ISO8601 variants only).
func ParseStrict(value string) (DateTime, error) {
	return ParseStrictInLocation(value, time.UTC)
}

// ParseStrictInLocation is the location-aware strict parser.
func ParseStrictInLocation(value string, loc *time.Location) (DateTime, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return DateTime{}, ParseError(value, errors.New("empty string"))
	}
	strictLayouts := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
	}
	for _, layout := range strictLayouts {
		if t, err := time.ParseInLocation(layout, value, loc); err == nil {
			return DateTime{t}, nil
		}
	}
	return DateTime{}, ParseError(value, errors.New("no matching strict format found"))
}

// ParseISO8601 parses an ISO 8601 formatted datetime string.
func ParseISO8601(value string) (DateTime, error) {
	if !iso8601Pattern.MatchString(value) {
		return DateTime{}, ParseError(value, errors.New("invalid ISO 8601 format"))
	}

	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		// Try with nanoseconds
		t, err = time.Parse(time.RFC3339Nano, value)
		if err != nil {
			return DateTime{}, ParseError(value, err)
		}
	}

	return DateTime{t}, nil
}

// ParseRFC3339 parses an RFC 3339 formatted datetime string.
func ParseRFC3339(value string) (DateTime, error) {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return DateTime{}, ParseError(value, err)
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
		return DateTime{}, FormatError(layout, err)
	}
	return DateTime{t}, nil
}

// LoadLocation loads a timezone by name.
// This is a convenience wrapper around time.LoadLocation.
func LoadLocation(name string) (*time.Location, error) {
	if name == "local" {
		return time.Local, nil
	}
	loc, err := time.LoadLocation(name)
	if err != nil {
		return nil, TimezoneError(name, err)
	}
	return loc, nil
}

// Instance creates a DateTime from a standard time.Time.
func Instance(t time.Time) DateTime {
	return DateTime{t}
}

// parseUnixTimestamp attempts to parse a Unix timestamp string.
func parseUnixTimestamp(value string) (DateTime, error) {
	s := strings.TrimSpace(value)
	if s == "" {
		return DateTime{}, ParseError(value, errors.New("invalid Unix timestamp"))
	}

	// Detect length ignoring leading sign
	signless := s
	if s[0] == '-' || s[0] == '+' {
		signless = s[1:]
	}

	// Parse as int64
	ts, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return DateTime{}, ParseError(value, errors.New("invalid Unix timestamp"))
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
		return DateTime{}, ParseError(value, errors.New("invalid Unix timestamp length"))
	}
}

// TryParseUnix attempts to parse a string as a Unix timestamp.
func TryParseUnix(value string) (DateTime, error) {
	return parseUnixTimestamp(value)
}

// AvailableTimezones returns a list of commonly used timezone names.
// This is helpful for error suggestions and validation.
func AvailableTimezones() []string {
	return []string{
		"UTC",
		"Local",
		// Americas
		"America/New_York",
		"America/Chicago",
		"America/Denver",
		"America/Los_Angeles",
		"America/Toronto",
		"America/Vancouver",
		"America/Mexico_City",
		"America/Sao_Paulo",
		"America/Argentina/Buenos_Aires",
		// Europe
		"Europe/London",
		"Europe/Paris",
		"Europe/Berlin",
		"Europe/Rome",
		"Europe/Madrid",
		"Europe/Amsterdam",
		"Europe/Stockholm",
		"Europe/Moscow",
		// Asia
		"Asia/Tokyo",
		"Asia/Shanghai",
		"Asia/Hong_Kong",
		"Asia/Singapore",
		"Asia/Bangkok",
		"Asia/Jakarta",
		"Asia/Manila",
		"Asia/Seoul",
		"Asia/Kolkata",
		"Asia/Dubai",
		// Australia/Pacific
		"Australia/Sydney",
		"Australia/Melbourne",
		"Australia/Perth",
		"Pacific/Auckland",
		"Pacific/Honolulu",
		// Africa
		"Africa/Cairo",
		"Africa/Johannesburg",
		"Africa/Lagos",
	}
}

// IsValidTimezone checks if a timezone name is valid.
func IsValidTimezone(name string) bool {
	_, err := time.LoadLocation(name)
	return err == nil
}
