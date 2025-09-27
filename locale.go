package chronogo

import (
	"fmt"
	"strings"
	"sync"
)

// Locale represents a specific locale configuration for formatting dates and times.
type Locale struct {
	Code         string                   // Locale code (e.g., "en-US", "es-ES")
	Name         string                   // Human-readable name
	MonthNames   []string                 // Full month names (January, February, ...)
	MonthAbbr    []string                 // Abbreviated month names (Jan, Feb, ...)
	WeekdayNames []string                 // Full weekday names (Monday, Tuesday, ...)
	WeekdayAbbr  []string                 // Abbreviated weekday names (Mon, Tue, ...)
	AMPMNames    []string                 // AM/PM indicators
	Ordinals     map[int]string           // Ordinal suffixes (1st, 2nd, 3rd, ...)
	TimeUnits    map[string]TimeUnitNames // Time unit names for human differences
	DateFormats  map[string]string        // Common date format patterns
}

// TimeUnitNames contains singular and plural forms for time units
type TimeUnitNames struct {
	Singular string
	Plural   string
}

// LocaleRegistry manages available locales
type LocaleRegistry struct {
	locales map[string]*Locale
	mutex   sync.RWMutex
}

var (
	// Global locale registry
	localeRegistry = &LocaleRegistry{
		locales: make(map[string]*Locale),
	}

	// Default locale
	defaultLocale = "en-US"
)

// RegisterLocale registers a new locale in the global registry
func RegisterLocale(locale *Locale) {
	localeRegistry.mutex.Lock()
	defer localeRegistry.mutex.Unlock()
	localeRegistry.locales[locale.Code] = locale
}

// GetLocale retrieves a locale by code
func GetLocale(code string) (*Locale, error) {
	localeRegistry.mutex.RLock()
	defer localeRegistry.mutex.RUnlock()

	locale, exists := localeRegistry.locales[code]
	if !exists {
		return nil, fmt.Errorf("locale %q not found", code)
	}
	return locale, nil
}

// GetAvailableLocales returns all registered locale codes
func GetAvailableLocales() []string {
	localeRegistry.mutex.RLock()
	defer localeRegistry.mutex.RUnlock()

	codes := make([]string, 0, len(localeRegistry.locales))
	for code := range localeRegistry.locales {
		codes = append(codes, code)
	}
	return codes
}

// SetDefaultLocale sets the default locale for operations
func SetDefaultLocale(code string) error {
	if _, err := GetLocale(code); err != nil {
		return err
	}
	defaultLocale = code
	return nil
}

// GetDefaultLocale returns the current default locale code
func GetDefaultLocale() string {
	return defaultLocale
}

// FormatLocalized formats the datetime using locale-specific patterns
func (dt DateTime) FormatLocalized(pattern, localeCode string) (string, error) {
	locale, err := GetLocale(localeCode)
	if err != nil {
		return "", err
	}

	return dt.formatWithLocale(pattern, locale), nil
}

// FormatLocalizedDefault formats using the default locale
func (dt DateTime) FormatLocalizedDefault(pattern string) string {
	locale, err := GetLocale(defaultLocale)
	if err != nil {
		// Fallback to English if default locale fails
		locale, _ = GetLocale("en-US")
	}
	return dt.formatWithLocale(pattern, locale)
}

// formatWithLocale performs the actual formatting with locale data
func (dt DateTime) formatWithLocale(pattern string, locale *Locale) string {
	// First, convert all standard tokens to Go format
	goLayout := convertTokenFormat(pattern)

	// Format with Go's time package
	result := dt.Format(goLayout)

	// Now replace the English parts with localized versions
	// We need to be careful about the order and use the actual formatted values

	// Get the English month/weekday names that Go would have produced
	englishMonth := dt.Format("January")
	englishMonthAbbr := dt.Format("Jan")
	englishWeekday := dt.Format("Monday")
	englishWeekdayAbbr := dt.Format("Mon")

	// Replace with localized versions
	localizedMonth := locale.MonthNames[dt.Month()-1]
	localizedMonthAbbr := locale.MonthAbbr[dt.Month()-1]
	localizedWeekday := locale.WeekdayNames[dt.Weekday()]
	localizedWeekdayAbbr := locale.WeekdayAbbr[dt.Weekday()]

	// Replace English names with localized names
	result = strings.ReplaceAll(result, englishMonth, localizedMonth)
	result = strings.ReplaceAll(result, englishMonthAbbr, localizedMonthAbbr)
	result = strings.ReplaceAll(result, englishWeekday, localizedWeekday)
	result = strings.ReplaceAll(result, englishWeekdayAbbr, localizedWeekdayAbbr)

	// Handle ordinals - Go's "2nd" format becomes the actual day with suffix
	if strings.Contains(pattern, "Do") {
		day := dt.Day()
		// Go will format "Do" token as "2nd" but we want actual day with localized suffix
		goOrdinalPattern := dt.Format("2nd") // This gives us the Go-formatted ordinal
		localizedOrdinal := fmt.Sprintf("%d%s", day, locale.getOrdinalSuffix(day))
		result = strings.ReplaceAll(result, goOrdinalPattern, localizedOrdinal)
	}

	// Handle AM/PM
	if strings.Contains(pattern, "A") || strings.Contains(pattern, "a") {
		englishAM := "AM"
		englishPM := "PM"
		localizedAM := locale.AMPMNames[0]
		localizedPM := locale.AMPMNames[1]

		result = strings.ReplaceAll(result, englishAM, localizedAM)
		result = strings.ReplaceAll(result, englishPM, localizedPM)
		result = strings.ReplaceAll(result, strings.ToLower(englishAM), strings.ToLower(localizedAM))
		result = strings.ReplaceAll(result, strings.ToLower(englishPM), strings.ToLower(localizedPM))
	}

	return result
}

// getOrdinalSuffix returns the ordinal suffix for a number in the locale
func (locale *Locale) getOrdinalSuffix(n int) string {
	if suffix, exists := locale.Ordinals[n]; exists {
		return suffix
	}

	// Default English-style ordinals as fallback
	switch n % 10 {
	case 1:
		if n%100 != 11 {
			return "st"
		}
	case 2:
		if n%100 != 12 {
			return "nd"
		}
	case 3:
		if n%100 != 13 {
			return "rd"
		}
	}
	return "th"
}

// HumanStringLocalized returns a human-readable difference in the specified locale
func (dt DateTime) HumanStringLocalized(localeCode string, other ...DateTime) (string, error) {
	locale, err := GetLocale(localeCode)
	if err != nil {
		return "", err
	}

	var reference DateTime
	if len(other) > 0 {
		reference = other[0]
	} else {
		reference = Now()
	}

	return dt.humanStringWithLocale(reference, locale), nil
}

// HumanStringLocalizedDefault returns a human-readable difference using the default locale
func (dt DateTime) HumanStringLocalizedDefault(other ...DateTime) string {
	locale, err := GetLocale(defaultLocale)
	if err != nil {
		// Fallback to English
		locale, _ = GetLocale("en-US")
	}

	var reference DateTime
	if len(other) > 0 {
		reference = other[0]
	} else {
		reference = Now()
	}

	return dt.humanStringWithLocale(reference, locale)
}

// humanStringWithLocale generates human-readable time differences using locale data
func (dt DateTime) humanStringWithLocale(reference DateTime, locale *Locale) string {
	duration := dt.Sub(reference)
	isPast := duration < 0
	if isPast {
		duration = -duration
	}

	// Determine the appropriate unit and value
	var unit string
	var value int

	seconds := int(duration.Seconds())
	minutes := int(duration.Minutes())
	hours := int(duration.Hours())
	days := int(duration.Hours() / 24)
	weeks := days / 7
	months := days / 30 // Approximate
	years := days / 365 // Approximate

	switch {
	case years > 0:
		unit = "year"
		value = years
	case months > 0:
		unit = "month"
		value = months
	case weeks > 0:
		unit = "week"
		value = weeks
	case days > 0:
		unit = "day"
		value = days
	case hours > 0:
		unit = "hour"
		value = hours
	case minutes > 0:
		unit = "minute"
		value = minutes
	default:
		unit = "second"
		value = seconds
		if value < 10 {
			return locale.formatFewMoments(isPast)
		}
	}

	return locale.formatTimeUnit(unit, value, isPast)
}

// formatTimeUnit formats a time unit with proper singular/plural and tense
func (locale *Locale) formatTimeUnit(unit string, value int, isPast bool) string {
	unitNames, exists := locale.TimeUnits[unit]
	if !exists {
		// Fallback to English
		return fmt.Sprintf("%d %s", value, unit)
	}

	unitName := unitNames.Singular
	if value != 1 {
		unitName = unitNames.Plural
	}

	// Check for locale-specific patterns
	if patterns, exists := locale.TimeUnits["patterns"]; exists {
		if isPast {
			return fmt.Sprintf(patterns.Singular, value, unitName) // past pattern
		} else {
			return fmt.Sprintf(patterns.Plural, value, unitName) // future pattern
		}
	}

	// Default English-style formatting
	if isPast {
		return fmt.Sprintf("%d %s ago", value, unitName)
	} else {
		return fmt.Sprintf("in %d %s", value, unitName)
	}
}

// formatFewMoments formats "a few moments" type messages
func (locale *Locale) formatFewMoments(isPast bool) string {
	if moments, exists := locale.TimeUnits["moments"]; exists {
		if isPast {
			return moments.Singular // "hace unos momentos"
		} else {
			return moments.Plural // "en unos momentos"
		}
	}

	// Fallback to English
	if isPast {
		return "a few seconds ago"
	}
	return "in a few seconds"
}

// GetMonthName returns the localized month name
func (dt DateTime) GetMonthName(localeCode string) (string, error) {
	locale, err := GetLocale(localeCode)
	if err != nil {
		return "", err
	}
	return locale.MonthNames[dt.Month()-1], nil
}

// GetMonthNameDefault returns the localized month name using default locale
func (dt DateTime) GetMonthNameDefault() string {
	name, _ := dt.GetMonthName(defaultLocale)
	return name
}

// GetWeekdayName returns the localized weekday name
func (dt DateTime) GetWeekdayName(localeCode string) (string, error) {
	locale, err := GetLocale(localeCode)
	if err != nil {
		return "", err
	}
	return locale.WeekdayNames[dt.Weekday()], nil
}

// GetWeekdayNameDefault returns the localized weekday name using default locale
func (dt DateTime) GetWeekdayNameDefault() string {
	name, _ := dt.GetWeekdayName(defaultLocale)
	return name
}

// init registers default locales
func init() {
	registerDefaultLocales()
}
