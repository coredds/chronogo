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
	return parseUnixTimestampInLocation(value, time.UTC)
}

// parseUnixTimestampInLocation attempts to parse a Unix timestamp string in a specific location.
func parseUnixTimestampInLocation(value string, loc *time.Location) (DateTime, error) {
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

	var t time.Time
	switch l := len(signless); l {
	case 1, 2, 3, 4, 5, 6, 7, 8, 9, 10: // seconds (support small absolute values too)
		t = time.Unix(ts, 0)
	case 13: // milliseconds
		t = time.UnixMilli(ts)
	case 16: // microseconds
		t = time.UnixMicro(ts)
	case 19: // nanoseconds
		t = time.Unix(0, ts)
	default:
		return DateTime{}, ParseError(value, errors.New("invalid Unix timestamp length"))
	}

	// Convert to the specified location
	return DateTime{t.In(loc)}, nil
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

// ParseOptimized provides faster parsing by using heuristics to detect the format
func ParseOptimized(value string) (DateTime, error) {
	return ParseOptimizedInLocation(value, time.UTC)
}

// ParseOptimizedInLocation parses with optimized format detection
func ParseOptimizedInLocation(value string, loc *time.Location) (DateTime, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return DateTime{}, ParseError(value, errors.New("empty string"))
	}

	// Check for compact date format first (8 digits like 20231225)
	if len(value) == 8 && isAllDigits(value) {
		if t, err := time.ParseInLocation("20060102", value, loc); err == nil {
			return DateTime{t}, nil
		}
	}

	// Fast path for numeric-only strings (Unix timestamps) - but not 8-digit dates
	if isNumericOnly(value) && len(value) != 8 {
		return parseUnixTimestampInLocation(value, loc)
	}

	// Detect format based on string characteristics
	layout := detectLayout(value)
	if layout != "" {
		if t, err := time.ParseInLocation(layout, value, loc); err == nil {
			return DateTime{t}, nil
		}
	}

	// Fallback to trying common layouts in optimized order
	return parseWithOptimizedOrder(value, loc)
}

// detectLayout uses heuristics to detect the likely format
func detectLayout(value string) string {
	length := len(value)

	// ISO 8601 / RFC 3339 patterns and space-separated datetime patterns
	if length >= 16 && value[4] == '-' && value[7] == '-' {
		if value[10] == 'T' && length >= 19 {
			if strings.HasSuffix(value, "Z") {
				if length == 20 {
					return "2006-01-02T15:04:05Z"
				}
				return time.RFC3339Nano
			}
			if length >= 25 && (value[19] == '+' || value[19] == '-') {
				return time.RFC3339
			}
			if length == 19 {
				return "2006-01-02T15:04:05"
			}
		} else if length >= 16 && len(value) > 10 && value[10] == ' ' {
			if length == 19 {
				return "2006-01-02 15:04:05"
			}
			if length == 16 {
				return "2006-01-02 15:04"
			}
		}
	}

	// Date-only patterns
	if length == 10 && value[4] == '-' && value[7] == '-' {
		return "2006-01-02"
	}

	// Compact date pattern
	if length == 8 && isAllDigits(value) {
		return "20060102"
	}

	// Slash-separated patterns
	if strings.Contains(value, "/") {
		if length >= 19 && value[4] == '/' && value[7] == '/' {
			return "2006/01/02 15:04:05"
		}
		if length == 10 && value[4] == '/' && value[7] == '/' {
			return "2006/01/02"
		}
	}

	// Time-only patterns
	if length >= 5 && strings.Contains(value, ":") && !strings.Contains(value, "-") && !strings.Contains(value, "/") {
		if length == 8 {
			return "15:04:05"
		}
		if length == 5 {
			return "15:04"
		}
	}

	return ""
}

// parseWithOptimizedOrder tries layouts in order of likelihood
func parseWithOptimizedOrder(value string, loc *time.Location) (DateTime, error) {
	// Reorder layouts based on frequency of use in real applications
	optimizedLayouts := []string{
		// Most common API formats first
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",

		// Less common but still used
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04",
		"2006/01/02 15:04:05",
		"2006/01/02",

		// Rare formats last
		"2006-1-2 15:04:05",
		"2006-1-2",
		"20060102",
		"15:04:05",
		"15:04",
	}

	for _, layout := range optimizedLayouts {
		if t, err := time.ParseInLocation(layout, value, loc); err == nil {
			return DateTime{t}, nil
		}
	}

	return DateTime{}, ParseError(value, errors.New("no matching format found"))
}

// isAllDigits checks if string contains only digits (faster than regex)
func isAllDigits(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// ParseBatch parses multiple values efficiently (useful for CSV imports, etc.)
func ParseBatch(values []string, loc *time.Location) ([]DateTime, []error) {
	results := make([]DateTime, len(values))
	parseErrors := make([]error, len(values))

	// Try to detect common format from first few values
	var detectedLayout string
	for i := 0; i < min(len(values), 3); i++ {
		if values[i] != "" {
			detectedLayout = detectLayout(strings.TrimSpace(values[i]))
			if detectedLayout != "" {
				break
			}
		}
	}

	// Parse all values, using detected layout first if available
	for i, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			parseErrors[i] = ParseError(value, errors.New("empty string"))
			continue
		}

		// Try detected layout first
		if detectedLayout != "" {
			if t, err := time.ParseInLocation(detectedLayout, value, loc); err == nil {
				results[i] = DateTime{t}
				continue
			}
		}

		// Fallback to optimized parsing
		if dt, err := ParseOptimizedInLocation(value, loc); err == nil {
			results[i] = dt
		} else {
			parseErrors[i] = err
		}
	}

	return results, parseErrors
}

// Helper function (Go 1.21+ has this in stdlib)
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
