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

// ParseOptions defines options for advanced parsing
type ParseOptions struct {
	Exact  bool // Return exact type (Date, Time, Interval) if true
	Strict bool // Use strict parsing (RFC3339/ISO8601 only) if true
}

// Parse parses a datetime string using common formats (lenient by default) in UTC.
func Parse(value string, options ...ParseOptions) (DateTime, error) {
	return ParseInLocation(value, time.UTC, options...)
}

// ParseInLocation parses a datetime string in the specified location using a lenient set of formats.
func ParseInLocation(value string, loc *time.Location, options ...ParseOptions) (DateTime, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return DateTime{}, ParseError(value, errors.New("empty string"))
	}

	var opts ParseOptions
	if len(options) > 0 {
		opts = options[0]
	}

	// Try ISO 8601 ordinal date format
	if dt, err := parseOrdinalDate(value, loc); err == nil {
		return dt, nil
	}

	// Try ISO 8601 week date format
	if dt, err := parseWeekDate(value, loc); err == nil {
		return dt, nil
	}

	// Try ISO 8601 interval format
	if strings.Contains(value, "/") {
		if interval, err := parseInterval(value, loc); err == nil {
			// For now, return the start of the interval
			// TODO: Consider returning a separate Interval type when exact=true
			return interval.Start, nil
		}
	}

	// Use strict or lenient layouts based on options
	layouts := commonLayouts
	if opts.Strict {
		layouts = []string{
			time.RFC3339,
			time.RFC3339Nano,
			"2006-01-02T15:04:05Z07:00",
			"2006-01-02T15:04:05",
		}
	}

	// Try each layout
	for _, layout := range layouts {
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

// FromFormatTokens parses a datetime string using ChronoGo-style format tokens.
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
		{"MM", "01"},
		{"DD", "02"},
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

// ParseWithFallback attempts to parse using the main parser, then falls back to a more lenient parser
func ParseWithFallback(value string) (DateTime, error) {
	return ParseWithFallbackInLocation(value, time.UTC)
}

// ParseWithFallbackInLocation attempts parsing with fallback support
func ParseWithFallbackInLocation(value string, loc *time.Location) (DateTime, error) {
	// First try the optimized parser
	if dt, err := ParseOptimizedInLocation(value, loc); err == nil {
		return dt, nil
	}

	// Try Unix timestamp parsing
	if dt, err := parseUnixTimestampInLocation(value, loc); err == nil {
		return dt, nil
	}

	// Try some additional lenient formats
	lenientFormats := []string{
		"1/2/2006 15:04:05",
		"1/2/2006 3:04:05 PM",
		"1/2/2006",
		"1-2-2006",
		"Jan 2, 2006",
		"January 2, 2006",
		"Jan 2, 2006 15:04:05",
		"January 2, 2006 3:04:05 PM",
		"2 Jan 2006",
		"2 January 2006",
		"Mon, 2 Jan 2006 15:04:05 MST",
		"Monday, 2 January 2006 15:04:05 MST",
	}

	for _, format := range lenientFormats {
		if t, err := time.ParseInLocation(format, value, loc); err == nil {
			return DateTime{t}, nil
		}
	}

	return DateTime{}, ParseError(value, errors.New("no matching format found even with fallback"))
}

// ParseMultiple attempts to parse a string that might contain multiple possible formats
func ParseMultiple(value string, formats []string) (DateTime, error) {
	return ParseMultipleInLocation(value, formats, time.UTC)
}

// ParseMultipleInLocation parses using multiple possible formats
func ParseMultipleInLocation(value string, formats []string, loc *time.Location) (DateTime, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return DateTime{}, ParseError(value, errors.New("empty string"))
	}

	// Convert any token-style formats to Go formats
	goFormats := make([]string, len(formats))
	for i, format := range formats {
		goFormats[i] = convertTokenFormat(format)
	}

	// Try each format
	for _, format := range goFormats {
		if t, err := time.ParseInLocation(format, value, loc); err == nil {
			return DateTime{t}, nil
		}
	}

	return DateTime{}, ParseError(value, errors.New("no matching format found"))
}

// ParseAny attempts to parse any reasonable datetime string with maximum leniency
func ParseAny(value string) (DateTime, error) {
	return ParseAnyInLocation(value, time.UTC)
}

// ParseAnyInLocation parses with maximum leniency in the specified location
func ParseAnyInLocation(value string, loc *time.Location) (DateTime, error) {
	// Try optimized parsing first
	if dt, err := ParseOptimizedInLocation(value, loc); err == nil {
		return dt, nil
	}

	// Try fallback parsing
	if dt, err := ParseWithFallbackInLocation(value, loc); err == nil {
		return dt, nil
	}

	// Last resort: try to clean up the string and parse again
	cleaned := cleanDateTimeString(value)
	if cleaned != value {
		if dt, err := ParseWithFallbackInLocation(cleaned, loc); err == nil {
			return dt, nil
		}
	}

	return DateTime{}, ParseError(value, errors.New("unable to parse datetime string"))
}

// cleanDateTimeString attempts to clean up a datetime string for parsing
func cleanDateTimeString(value string) string {
	// Remove common prefixes/suffixes
	value = strings.TrimSpace(value)

	// Remove quotes
	if len(value) >= 2 {
		if (value[0] == '"' && value[len(value)-1] == '"') ||
			(value[0] == '\'' && value[len(value)-1] == '\'') {
			value = value[1 : len(value)-1]
		}
	}

	// Normalize whitespace
	value = regexp.MustCompile(`\s+`).ReplaceAllString(value, " ")

	// Remove ordinal suffixes (1st, 2nd, 3rd, 4th, etc.)
	value = regexp.MustCompile(`(\d+)(st|nd|rd|th)`).ReplaceAllString(value, "$1")

	// Convert some common separators
	value = strings.ReplaceAll(value, " at ", " ")
	value = strings.ReplaceAll(value, " on ", " ")

	return strings.TrimSpace(value)
}

// IsValidDateTimeString quickly checks if a string might be a valid datetime
func IsValidDateTimeString(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}

	// Basic length check
	if len(value) < 4 || len(value) > 50 {
		return false
	}

	// Must contain at least one digit
	hasDigit := false
	for _, r := range value {
		if r >= '0' && r <= '9' {
			hasDigit = true
			break
		}
	}

	return hasDigit
}

// GetSupportedFormats returns a list of all supported datetime formats
func GetSupportedFormats() []string {
	return []string{
		// ISO 8601 / RFC 3339
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05.999999999Z",
		"2006-01-02T15:04:05.999999999Z07:00",
		"2006-01-02T15:04:05",

		// Common datetime formats
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
		"2006/01/02 15:04:05",
		"2006/01/02",
		"20060102",

		// Time only
		"15:04:05",
		"15:04",

		// US formats
		"1/2/2006 15:04:05",
		"1/2/2006 3:04:05 PM",
		"1/2/2006",
		"1-2-2006",

		// Human readable
		"Jan 2, 2006",
		"January 2, 2006",
		"Jan 2, 2006 15:04:05",
		"January 2, 2006 3:04:05 PM",
		"2 Jan 2006",
		"2 January 2006",

		// Email/HTTP dates
		"Mon, 2 Jan 2006 15:04:05 MST",
		"Monday, 2 January 2006 15:04:05 MST",

		// ISO 8601 ordinal and week dates
		"2006-002",   // Ordinal date
		"2006002",    // Compact ordinal date
		"2006-W01-1", // Week date
		"2006W011",   // Compact week date
		"2006-W01",   // Week only
		"2006W01",    // Compact week only

		// ISO 8601 intervals (examples)
		"2006-01-02T15:04:05Z/2007-01-02T15:04:05Z",
		"2006-01-02T15:04:05Z/P1Y",
		"P1Y/2007-01-02T15:04:05Z",
	}
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
