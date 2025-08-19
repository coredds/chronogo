package chronogo

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// ChronoError represents errors that occur in ChronoGo operations.
type ChronoError struct {
	Op      string // Operation that caused the error
	Path    string // Path or context where error occurred
	Err     error  // Underlying error
	Input   string // Input that caused the error (for parsing errors)
	Suggestion string // Helpful suggestion for fixing the error
}

// Error implements the error interface.
func (e *ChronoError) Error() string {
	var parts []string
	
	if e.Op != "" {
		parts = append(parts, fmt.Sprintf("chronogo.%s", e.Op))
	}
	
	if e.Path != "" {
		parts = append(parts, fmt.Sprintf("(%s)", e.Path))
	}
	
	if e.Input != "" {
		parts = append(parts, fmt.Sprintf("input: %q", e.Input))
	}
	
	if e.Err != nil {
		parts = append(parts, e.Err.Error())
	}
	
	result := strings.Join(parts, ": ")
	
	if e.Suggestion != "" {
		result += fmt.Sprintf("\nSuggestion: %s", e.Suggestion)
	}
	
	return result
}

// Unwrap returns the underlying error for error wrapping support.
func (e *ChronoError) Unwrap() error {
	return e.Err
}

// Is implements error comparison for errors.Is().
func (e *ChronoError) Is(target error) bool {
	if target == nil {
		return false
	}
	
	if te, ok := target.(*ChronoError); ok {
		return e.Op == te.Op && e.Path == te.Path
	}
	
	return errors.Is(e.Err, target)
}

// Common error variables for easier error checking
var (
	ErrInvalidFormat    = errors.New("invalid datetime format")
	ErrInvalidTimezone  = errors.New("invalid timezone")
	ErrInvalidDuration  = errors.New("invalid duration")
	ErrInvalidRange     = errors.New("invalid range")
	ErrInvalidOperation = errors.New("invalid operation")
)

// ParseError creates a ChronoError for parsing operations.
func ParseError(input string, err error) *ChronoError {
	suggestion := suggestParseFormat(input)
	return &ChronoError{
		Op:         "Parse",
		Err:        err,
		Input:      input,
		Suggestion: suggestion,
	}
}

// TimezoneError creates a ChronoError for timezone operations.
func TimezoneError(timezone string, err error) *ChronoError {
	suggestion := suggestTimezone(timezone)
	return &ChronoError{
		Op:         "LoadLocation",
		Path:       timezone,
		Err:        err,
		Suggestion: suggestion,
	}
}

// FormatError creates a ChronoError for formatting operations.
func FormatError(format string, err error) *ChronoError {
	suggestion := suggestFormat(format)
	return &ChronoError{
		Op:         "Format",
		Input:      format,
		Err:        err,
		Suggestion: suggestion,
	}
}

// RangeError creates a ChronoError for range operations.
func RangeError(start, end DateTime, err error) *ChronoError {
	return &ChronoError{
		Op:         "Range",
		Path:       fmt.Sprintf("%s to %s", start.Format("2006-01-02"), end.Format("2006-01-02")),
		Err:        err,
		Suggestion: "Ensure start date is before end date",
	}
}

// suggestParseFormat provides helpful suggestions for parse errors.
func suggestParseFormat(input string) string {
	if input == "" {
		return "Provide a non-empty datetime string"
	}
	
	inputLower := strings.ToLower(input)
	
	// Detect likely patterns
	if strings.Contains(input, "T") && (strings.Contains(input, "Z") || strings.Contains(input, "+")) {
		return "For ISO 8601 format, try: chronogo.ParseISO8601() or chronogo.ParseRFC3339()"
	}
	
	if strings.Count(input, "/") == 2 {
		return "For date with slashes, try: chronogo.FromFormat(input, \"01/02/2006\") or \"02/01/2006\""
	}
	
	if strings.Count(input, "-") == 2 && len(input) >= 8 {
		return "For date with dashes, try: chronogo.Parse() which supports ISO format, or chronogo.FromFormat()"
	}
	
	if isNumericOnly(input) {
		if len(input) == 10 {
			return "For Unix timestamp (seconds), try: chronogo.FromUnix(timestamp, time.UTC)"
		} else if len(input) == 13 {
			return "For Unix timestamp (milliseconds), try: chronogo.FromUnixMilli(timestamp, time.UTC)"
		}
		return "For numeric timestamps, try: chronogo.FromUnix(), FromUnixMilli(), FromUnixMicro(), or FromUnixNano()"
	}
	
	if strings.Contains(inputLower, "ago") || strings.Contains(inputLower, "from now") {
		return "Human-readable relative times are not supported for parsing. Use absolute dates/times."
	}
	
	return "Try chronogo.Parse() for common formats, or chronogo.FromFormat() with a custom layout. See Go time package documentation for layout syntax."
}

// suggestTimezone provides helpful suggestions for timezone errors.
func suggestTimezone(timezone string) string {
	if timezone == "" {
		return "Provide a valid IANA timezone name like 'America/New_York' or 'Europe/London'"
	}
	
	common := []string{
		"UTC",
		"America/New_York",
		"America/Los_Angeles", 
		"America/Chicago",
		"Europe/London",
		"Europe/Paris",
		"Asia/Tokyo",
		"Asia/Shanghai",
		"Australia/Sydney",
	}
	
	// Simple similarity check
	timezoneLower := strings.ToLower(timezone)
	for _, tz := range common {
		if strings.Contains(strings.ToLower(tz), timezoneLower) || 
		   strings.Contains(timezoneLower, strings.ToLower(tz)) {
			return fmt.Sprintf("Did you mean '%s'? Use chronogo.LoadLocation(\"%s\")", tz, tz)
		}
	}
	
	if strings.Contains(timezoneLower, "est") || strings.Contains(timezoneLower, "eastern") {
		return "Try 'America/New_York' for Eastern Time"
	}
	if strings.Contains(timezoneLower, "pst") || strings.Contains(timezoneLower, "pacific") {
		return "Try 'America/Los_Angeles' for Pacific Time"
	}
	if strings.Contains(timezoneLower, "cst") || strings.Contains(timezoneLower, "central") {
		return "Try 'America/Chicago' for Central Time"
	}
	if strings.Contains(timezoneLower, "mst") || strings.Contains(timezoneLower, "mountain") {
		return "Try 'America/Denver' for Mountain Time"
	}
	
	return "Use IANA timezone names like 'America/New_York'. List available zones with: chronogo.AvailableTimezones()"
}

// suggestFormat provides helpful suggestions for format errors.
func suggestFormat(format string) string {
	if format == "" {
		return "Provide a format string using Go's reference time: 'Mon Jan 2 15:04:05 MST 2006'"
	}
	
	// Common format mistakes and corrections
	corrections := map[string]string{
		"YYYY": "2006",
		"yyyy": "2006", 
		"YY":   "06",
		"yy":   "06",
		"MM":   "01",
		"mm":   "04", // minutes
		"DD":   "02",
		"dd":   "02",
		"HH":   "15",
		"hh":   "03",
		"SS":   "05",
		"ss":   "05",
	}
	
	suggestion := "Use Go's reference time format. Common patterns:\n"
	suggestion += "  Date: '2006-01-02'\n"
	suggestion += "  Time: '15:04:05'\n"
	suggestion += "  DateTime: '2006-01-02 15:04:05'\n"
	suggestion += "  RFC3339: '2006-01-02T15:04:05Z07:00'"
	
	// Check for common mistakes
	for wrong, right := range corrections {
		if strings.Contains(format, wrong) {
			suggestion = fmt.Sprintf("Replace '%s' with '%s'. %s", wrong, right, suggestion)
			break
		}
	}
	
	return suggestion
}

// isNumericOnly checks if a string contains only digits.
func isNumericOnly(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// Validation helpers for better developer experience

// Validate checks if a DateTime is valid and returns a helpful error if not.
func (dt DateTime) Validate() error {
	if dt.IsZero() {
		return &ChronoError{
			Op:         "Validate",
			Err:        errors.New("zero DateTime"),
			Suggestion: "Initialize DateTime using chronogo.Now(), chronogo.Date(), or chronogo.Parse()",
		}
	}
	
	// Check for reasonable year range
	year := dt.Year()
	if year < 1 || year > 9999 {
		return &ChronoError{
			Op:         "Validate", 
			Path:       fmt.Sprintf("year=%d", year),
			Err:        errors.New("year out of reasonable range"),
			Suggestion: "Use years between 1 and 9999",
		}
	}
	
	return nil
}

// ValidateRange checks if a date range is valid.
func ValidateRange(start, end DateTime) error {
	if err := start.Validate(); err != nil {
		return fmt.Errorf("invalid start date: %w", err)
	}
	
	if err := end.Validate(); err != nil {
		return fmt.Errorf("invalid end date: %w", err)
	}
	
	if start.After(end) {
		return RangeError(start, end, ErrInvalidRange)
	}
	
	return nil
}

// MustParse is like Parse but panics on error. Useful for constants.
func MustParse(input string) DateTime {
	dt, err := Parse(input)
	if err != nil {
		panic(fmt.Sprintf("chronogo.MustParse: %v", err))
	}
	return dt
}

// MustParseInLocation is like ParseInLocation but panics on error.
func MustParseInLocation(input string, loc *time.Location) DateTime {
	dt, err := ParseInLocation(input, loc)
	if err != nil {
		panic(fmt.Sprintf("chronogo.MustParseInLocation: %v", err))
	}
	return dt
}

// MustFromFormat is like FromFormat but panics on error.
func MustFromFormat(input, layout string) DateTime {
	dt, err := FromFormat(input, layout)
	if err != nil {
		panic(fmt.Sprintf("chronogo.MustFromFormat: %v", err))
	}
	return dt
}

// MustLoadLocation is like LoadLocation but panics on error.
func MustLoadLocation(name string) *time.Location {
	loc, err := LoadLocation(name)
	if err != nil {
		panic(fmt.Sprintf("chronogo.MustLoadLocation: %v", err))
	}
	return loc
}
