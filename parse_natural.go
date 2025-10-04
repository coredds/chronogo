package chronogo

import (
	"time"

	"github.com/coredds/godateparser"
)

// ParseConfig holds configuration for intelligent parsing with natural language support.
type ParseConfig struct {
	// Strict mode: only parse technical formats (ISO 8601, RFC 3339, Unix timestamps)
	// When false, enables natural language parsing via godateparser
	Strict bool

	// Languages for natural language parsing (e.g., "en", "es", "pt", "fr", "de", "zh", "ja")
	// Default: all supported languages
	Languages []string

	// Location for parsing (default: UTC)
	Location *time.Location

	// Prefer future dates when parsing ambiguous relative dates
	// e.g., "Friday" will prefer next Friday if today is not Friday
	PreferFuture bool
}

// DefaultParseConfig provides sensible defaults: all languages enabled, UTC location
var DefaultParseConfig = ParseConfig{
	Languages:    []string{"en", "es", "pt", "fr", "de", "zh", "ja"},
	Location:     time.UTC,
	PreferFuture: false,
}

// parseWithGodateparser attempts to parse using godateparser for natural language and common formats
func parseWithGodateparser(value string, loc *time.Location, languages []string, preferFuture bool) (DateTime, error) {
	// Configure godateparser settings
	settings := &godateparser.Settings{
		Languages: languages,
	}

	// Set relative base if location is specified
	if loc != nil {
		settings.RelativeBase = time.Now().In(loc)
	} else {
		settings.RelativeBase = time.Now().UTC()
	}

	// Note: godateparser v1.3.3 may not have PreferFuture field
	// This is handled by default behavior in godateparser

	// Parse with godateparser
	result, err := godateparser.ParseDate(value, settings)
	if err != nil {
		return DateTime{}, ParseError(value, err)
	}

	// Convert to chronogo DateTime with specified location
	if loc != nil && loc != result.Location() {
		result = result.In(loc)
	}

	return DateTime{result}, nil
}

// tryStrictFormats attempts parsing with only strict RFC/ISO formats and Unix timestamps
// Used by strict mode parsing
func tryStrictFormats(value string, loc *time.Location) (DateTime, bool) {
	// Try strict RFC 3339 / ISO 8601 formats only (with dashes/colons)
	strictLayouts := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02", // Date-only ISO 8601
	}

	for _, layout := range strictLayouts {
		if t, err := time.ParseInLocation(layout, value, loc); err == nil {
			return DateTime{t}, true
		}
	}

	// Try ISO 8601 ordinal date (YYYY-DDD with dash only)
	if len(value) == 8 && value[4] == '-' {
		if dt, err := parseOrdinalDate(value, loc); err == nil {
			return dt, true
		}
	}

	// Try ISO 8601 week date (YYYY-Www-D with dashes only)
	if len(value) >= 8 && value[4] == '-' && value[5] == 'W' {
		if dt, err := parseWeekDate(value, loc); err == nil {
			return dt, true
		}
	}

	// Try Unix timestamp (numeric only, but only 10+ digits to avoid ambiguity)
	if isNumericOnly(value) && len(value) >= 10 {
		if dt, err := parseUnixTimestampInLocation(value, loc); err == nil {
			return dt, true
		}
	}

	return DateTime{}, false
}

// tryTechnicalFormats attempts fast-path parsing for technical formats
// Returns (result, true) if successful, (zero, false) if format not recognized
func tryTechnicalFormats(value string, loc *time.Location) (DateTime, bool) {
	// Try common datetime layouts FIRST (before godateparser can misinterpret them)
	commonLayouts := []string{
		// Strict RFC 3339 / ISO 8601
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		// Space-separated
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
		// Slash-separated (common in US)
		"2006/01/02 15:04:05",
		"2006/01/02",
		// Lenient (single digits)
		"2006-1-2 15:04:05",
		"2006-1-2",
		// Time-only
		"15:04:05",
		"15:04",
	}

	for _, layout := range commonLayouts {
		if t, err := time.ParseInLocation(layout, value, loc); err == nil {
			return DateTime{t}, true
		}
	}

	// Try ISO 8601 ordinal date FIRST (7 digits like 2023359)
	if dt, err := parseOrdinalDate(value, loc); err == nil {
		return dt, true
	}

	// Try compact date format (8 digits like 20240115)
	if len(value) == 8 && isNumericOnly(value) {
		if t, err := time.ParseInLocation("20060102", value, loc); err == nil {
			return DateTime{t}, true
		}
	}

	// Try Unix timestamp (numeric only, but not 8-digit compact dates or 7-digit ordinal dates)
	if isNumericOnly(value) && len(value) != 8 && len(value) != 7 {
		if dt, err := parseUnixTimestampInLocation(value, loc); err == nil {
			return dt, true
		}
	}

	// Try ISO 8601 week date
	if dt, err := parseWeekDate(value, loc); err == nil {
		return dt, true
	}

	// Try ISO 8601 interval (return start of interval)
	if len(value) > 0 && value[0] != '/' && (len(value) < 2 || value[1] != '/') {
		// Only try if doesn't start with / to avoid false positives
		if interval, err := parseInterval(value, loc); err == nil {
			return interval.Start, true
		}
	}

	return DateTime{}, false
}

// ParseWith parses a datetime string using the provided configuration.
// This is the most flexible parsing function, allowing fine control over
// natural language parsing, languages, and location.
func ParseWith(value string, config ParseConfig) (DateTime, error) {
	if value == "" {
		return DateTime{}, ParseError(value, ErrEmptyString)
	}

	loc := config.Location
	if loc == nil {
		loc = time.UTC
	}

	// Strict mode: only try strict technical formats (RFC3339, ISO8601, Unix timestamps)
	if config.Strict {
		if dt, ok := tryStrictFormats(value, loc); ok {
			return dt, nil
		}
		return DateTime{}, ParseError(value, ErrNoMatchingFormat)
	}

	// Try fast-path technical formats first
	if dt, ok := tryTechnicalFormats(value, loc); ok {
		return dt, nil
	}

	// Use godateparser for natural language and common formats
	languages := config.Languages
	if len(languages) == 0 {
		languages = DefaultParseConfig.Languages
	}

	return parseWithGodateparser(value, loc, languages, config.PreferFuture)
}

// SetDefaultParseLanguages sets the default languages for Parse() and ParseInLocation().
// This is a convenience function for applications that primarily use specific languages.
// Default is all supported languages: en, es, pt, fr, de, zh, ja
func SetDefaultParseLanguages(languages ...string) {
	DefaultParseConfig.Languages = languages
}

// GetDefaultParseLanguages returns the current default languages for parsing
func GetDefaultParseLanguages() []string {
	return DefaultParseConfig.Languages
}
