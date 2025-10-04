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

	// ISO 8601 ordinal date pattern (YYYY-DDD or YYYYDDD)
	ordinalDatePattern = regexp.MustCompile(`^(\d{4})-?(\d{3})$`)

	// ISO 8601 week date patterns (YYYY-Www-D, YYYY-Www, YYYYWwwD, YYYYWww)
	weekDatePattern = regexp.MustCompile(`^(\d{4})-?W(\d{2})-?([1-7])?$`)

	// ISO 8601 interval patterns
	intervalPattern = regexp.MustCompile(`^(.+)/(.+)$`)

	// Duration pattern for ISO 8601 intervals (P[n]Y[n]M[n]DT[n]H[n]M[n]S)
	durationPattern = regexp.MustCompile(`^P(?:(\d+)Y)?(?:(\d+)M)?(?:(\d+)D)?(?:T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+(?:\.\d+)?)S)?)?$`)

	// Common parsing errors
	ErrEmptyString      = errors.New("empty string")
	ErrNoMatchingFormat = errors.New("no matching format found")
)

// ParseOptions defines options for advanced parsing
type ParseOptions struct {
	Exact  bool // Return exact type (Date, Time, Interval) if true
	Strict bool // Use strict parsing (RFC3339/ISO8601 only) if true
}

// Parse is an intelligent datetime parser that handles:
// - Technical formats (ISO 8601, RFC 3339, Unix timestamps) via fast path
// - Natural language expressions ("tomorrow", "next Monday", "3 days ago")
// - Common datetime formats with lenient parsing
// - Multi-language support (English, Spanish, Portuguese, French, German, Chinese, Japanese)
//
// The parser automatically tries technical formats first for performance, then falls back
// to natural language parsing via godateparser for maximum flexibility.
//
// Examples:
//
//	Parse("2024-01-15T14:30:00Z")     // ISO 8601
//	Parse("tomorrow")                  // Natural language
//	Parse("next Monday")               // Relative date
//	Parse("3 days ago")                // Relative with quantity
//	Parse("demain")                    // French for "tomorrow"
//	Parse("明天")                       // Chinese for "tomorrow"
func Parse(value string, options ...ParseOptions) (DateTime, error) {
	return ParseInLocation(value, time.UTC, options...)
}

// ParseInLocation parses a datetime string in the specified location.
// Supports the same formats as Parse() but uses the provided location for timezone-naive inputs
// and relative date calculations.
func ParseInLocation(value string, loc *time.Location, options ...ParseOptions) (DateTime, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return DateTime{}, ParseError(value, ErrEmptyString)
	}

	var opts ParseOptions
	if len(options) > 0 {
		opts = options[0]
	}

	// Build ParseConfig from options
	config := ParseConfig{
		Strict:    opts.Strict,
		Languages: DefaultParseConfig.Languages,
		Location:  loc,
	}

	return ParseWith(value, config)
}

// ParseStrict parses using only technical formats (RFC3339, ISO8601, Unix timestamps).
// No natural language parsing is attempted.
func ParseStrict(value string) (DateTime, error) {
	return ParseStrictInLocation(value, time.UTC)
}

// ParseStrictInLocation is the location-aware strict parser.
// Only technical formats are accepted (RFC3339, ISO8601, Unix timestamps).
func ParseStrictInLocation(value string, loc *time.Location) (DateTime, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return DateTime{}, ParseError(value, ErrEmptyString)
	}

	config := ParseConfig{
		Strict:   true,
		Location: loc,
	}

	return ParseWith(value, config)
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
// This is similar to strptime but uses Go's time format syntax.
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

// FromFormatTokens parses a datetime string using chronogo-style format tokens.
// Converts token-based format to Go time format and then parses.
func FromFormatTokens(value, format string) (DateTime, error) {
	return FromFormatTokensInLocation(value, format, time.UTC)
}

// FromFormatTokensInLocation parses using token-style format in the specified location.
func FromFormatTokensInLocation(value, format string, loc *time.Location) (DateTime, error) {
	goLayout := convertTokenFormat(format)
	return FromFormatInLocation(value, goLayout, loc)
}

// convertTokenFormat converts token-style format to Go time layout
func convertTokenFormat(format string) string {
	// Use a state machine approach to replace tokens without conflicts
	result := format

	// Define token mappings in order of specificity (longest first)
	tokens := []struct {
		token       string
		replacement string
	}{
		{"YYYY", "2006"},
		{"MMMM", "January"},
		{"MMM", "Jan"},
		{"dddd", "Monday"},
		{"ddd", "Mon"},
		{"MM", "01"},
		{"DD", "02"},
		{"Do", "2nd"}, // Ordinal day - Go doesn't have native support, but we'll handle this specially
		{"HH", "15"},
		{"hh", "03"},
		{"mm", "04"},
		{"ss", "05"},
		{"ZZ", "Z0700"},
		{"YY", "06"},
		{"Y", "2006"},
		{"M", "1"},
		{"D", "2"},
		{"H", "15"},
		{"h", "3"},
		{"m", "4"},
		{"s", "5"},
		{"A", "PM"},
		{"a", "pm"},
		{"Z", "Z07:00"},
	}

	// Process each position in the string
	i := 0
	for i < len(result) {
		matched := false

		// Try to match any token at current position
		for _, token := range tokens {
			if i+len(token.token) <= len(result) && result[i:i+len(token.token)] == token.token {
				// Check if this is a complete token (not part of a larger identifier)
				validStart := i == 0 || !isTokenChar(result[i-1])
				validEnd := i+len(token.token) == len(result) || !isTokenChar(result[i+len(token.token)])

				if validStart && validEnd {
					// Replace the token
					result = result[:i] + token.replacement + result[i+len(token.token):]
					i += len(token.replacement)
					matched = true
					break
				}
			}
		}

		if !matched {
			i++
		}
	}

	return result
}

// isTokenChar checks if a character can be part of a format token
func isTokenChar(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
}

// parseOrdinalDate parses ISO 8601 ordinal date format (YYYY-DDD or YYYYDDD)
func parseOrdinalDate(value string, loc *time.Location) (DateTime, error) {
	matches := ordinalDatePattern.FindStringSubmatch(value)
	if len(matches) != 3 {
		return DateTime{}, ParseError(value, errors.New("not an ordinal date"))
	}

	year, err := strconv.Atoi(matches[1])
	if err != nil {
		return DateTime{}, ParseError(value, errors.New("invalid year in ordinal date"))
	}

	dayOfYear, err := strconv.Atoi(matches[2])
	if err != nil {
		return DateTime{}, ParseError(value, errors.New("invalid day of year in ordinal date"))
	}

	// Validate day of year
	isLeap := year%4 == 0 && (year%100 != 0 || year%400 == 0)
	maxDays := 365
	if isLeap {
		maxDays = 366
	}

	if dayOfYear < 1 || dayOfYear > maxDays {
		return DateTime{}, ParseError(value, fmt.Errorf("invalid day of year: %d", dayOfYear))
	}

	// Create date from ordinal day
	t := time.Date(year, 1, 1, 0, 0, 0, 0, loc).AddDate(0, 0, dayOfYear-1)
	return DateTime{t}, nil
}

// parseWeekDate parses ISO 8601 week date format (YYYY-Www-D, YYYY-Www, YYYYWwwD, YYYYWww)
func parseWeekDate(value string, loc *time.Location) (DateTime, error) {
	matches := weekDatePattern.FindStringSubmatch(value)
	if len(matches) != 4 {
		return DateTime{}, ParseError(value, errors.New("not a week date"))
	}

	year, err := strconv.Atoi(matches[1])
	if err != nil {
		return DateTime{}, ParseError(value, errors.New("invalid year in week date"))
	}

	week, err := strconv.Atoi(matches[2])
	if err != nil {
		return DateTime{}, ParseError(value, errors.New("invalid week in week date"))
	}

	// Default to Monday (1) if day is not specified
	dayOfWeek := 1
	if matches[3] != "" {
		dayOfWeek, err = strconv.Atoi(matches[3])
		if err != nil || dayOfWeek < 1 || dayOfWeek > 7 {
			return DateTime{}, ParseError(value, errors.New("invalid day of week in week date"))
		}
	}

	// Validate week number
	if week < 1 || week > 53 {
		return DateTime{}, ParseError(value, fmt.Errorf("invalid week number: %d", week))
	}

	// Calculate the date from ISO week
	t := isoWeekToDate(year, week, dayOfWeek, loc)
	return DateTime{t}, nil
}

// isoWeekToDate converts ISO week year, week number, and day of week to a time.Time
func isoWeekToDate(year, week, dayOfWeek int, loc *time.Location) time.Time {
	// January 4th is always in week 1 of the ISO week-numbering year
	jan4 := time.Date(year, 1, 4, 0, 0, 0, 0, loc)

	// Find Monday of week 1
	jan4Weekday := int(jan4.Weekday())
	if jan4Weekday == 0 { // Sunday
		jan4Weekday = 7
	}

	mondayWeek1 := jan4.AddDate(0, 0, -(jan4Weekday - 1))

	// Add weeks and days to get to the target date
	targetDate := mondayWeek1.AddDate(0, 0, (week-1)*7+(dayOfWeek-1))

	return targetDate
}

// parseInterval parses ISO 8601 interval format (start/end, start/duration, duration/end)
func parseInterval(value string, loc *time.Location) (Period, error) {
	matches := intervalPattern.FindStringSubmatch(value)
	if len(matches) != 3 {
		return Period{}, ParseError(value, errors.New("not an interval"))
	}

	part1 := strings.TrimSpace(matches[1])
	part2 := strings.TrimSpace(matches[2])

	// Case 1: start/end (both are datetimes)
	if !strings.HasPrefix(part1, "P") && !strings.HasPrefix(part2, "P") {
		start, err := ParseInLocation(part1, loc)
		if err != nil {
			return Period{}, ParseError(value, fmt.Errorf("invalid start datetime: %v", err))
		}

		end, err := ParseInLocation(part2, loc)
		if err != nil {
			return Period{}, ParseError(value, fmt.Errorf("invalid end datetime: %v", err))
		}

		return NewPeriod(start, end), nil
	}

	// Case 2: start/duration
	if !strings.HasPrefix(part1, "P") && strings.HasPrefix(part2, "P") {
		start, err := ParseInLocation(part1, loc)
		if err != nil {
			return Period{}, ParseError(value, fmt.Errorf("invalid start datetime: %v", err))
		}

		duration, err := parseDuration(part2)
		if err != nil {
			return Period{}, ParseError(value, fmt.Errorf("invalid duration: %v", err))
		}

		end := addDurationToDateTime(start, duration)
		return NewPeriod(start, end), nil
	}

	// Case 3: duration/end
	if strings.HasPrefix(part1, "P") && !strings.HasPrefix(part2, "P") {
		duration, err := parseDuration(part1)
		if err != nil {
			return Period{}, ParseError(value, fmt.Errorf("invalid duration: %v", err))
		}

		end, err := ParseInLocation(part2, loc)
		if err != nil {
			return Period{}, ParseError(value, fmt.Errorf("invalid end datetime: %v", err))
		}

		start := subtractDurationFromDateTime(end, duration)
		return NewPeriod(start, end), nil
	}

	return Period{}, ParseError(value, errors.New("invalid interval format"))
}

// ISO8601Duration represents a parsed ISO 8601 duration
type ISO8601Duration struct {
	Years   int
	Months  int
	Days    int
	Hours   int
	Minutes int
	Seconds float64
}

// parseDuration parses ISO 8601 duration format (P[n]Y[n]M[n]DT[n]H[n]M[n]S)
func parseDuration(value string) (ISO8601Duration, error) {
	matches := durationPattern.FindStringSubmatch(value)
	if len(matches) != 7 {
		return ISO8601Duration{}, fmt.Errorf("invalid duration format: %s", value)
	}

	var duration ISO8601Duration
	var err error

	if matches[1] != "" { // Years
		duration.Years, err = strconv.Atoi(matches[1])
		if err != nil {
			return ISO8601Duration{}, fmt.Errorf("invalid years: %s", matches[1])
		}
	}

	if matches[2] != "" { // Months
		duration.Months, err = strconv.Atoi(matches[2])
		if err != nil {
			return ISO8601Duration{}, fmt.Errorf("invalid months: %s", matches[2])
		}
	}

	if matches[3] != "" { // Days
		duration.Days, err = strconv.Atoi(matches[3])
		if err != nil {
			return ISO8601Duration{}, fmt.Errorf("invalid days: %s", matches[3])
		}
	}

	if matches[4] != "" { // Hours
		duration.Hours, err = strconv.Atoi(matches[4])
		if err != nil {
			return ISO8601Duration{}, fmt.Errorf("invalid hours: %s", matches[4])
		}
	}

	if matches[5] != "" { // Minutes
		duration.Minutes, err = strconv.Atoi(matches[5])
		if err != nil {
			return ISO8601Duration{}, fmt.Errorf("invalid minutes: %s", matches[5])
		}
	}

	if matches[6] != "" { // Seconds
		duration.Seconds, err = strconv.ParseFloat(matches[6], 64)
		if err != nil {
			return ISO8601Duration{}, fmt.Errorf("invalid seconds: %s", matches[6])
		}
	}

	return duration, nil
}

// addDurationToDateTime adds an ISO8601Duration to a DateTime
func addDurationToDateTime(dt DateTime, duration ISO8601Duration) DateTime {
	result := dt.AddYears(duration.Years).AddMonths(duration.Months).AddDays(duration.Days)
	result = result.AddHours(duration.Hours).AddMinutes(duration.Minutes)

	// Handle fractional seconds
	seconds := int(duration.Seconds)
	nanoseconds := int((duration.Seconds - float64(seconds)) * 1e9)
	result = result.AddSeconds(seconds)

	if nanoseconds > 0 {
		result = result.Add(time.Duration(nanoseconds) * time.Nanosecond)
	}

	return result
}

// subtractDurationFromDateTime subtracts an ISO8601Duration from a DateTime
func subtractDurationFromDateTime(dt DateTime, duration ISO8601Duration) DateTime {
	result := dt.SubtractYears(duration.Years).SubtractMonths(duration.Months).SubtractDays(duration.Days)
	result = result.SubtractHours(duration.Hours).SubtractMinutes(duration.Minutes)

	// Handle fractional seconds
	seconds := int(duration.Seconds)
	nanoseconds := int((duration.Seconds - float64(seconds)) * 1e9)
	result = result.SubtractSeconds(seconds)

	if nanoseconds > 0 {
		result = result.Subtract(time.Duration(nanoseconds) * time.Nanosecond)
	}

	return result
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
