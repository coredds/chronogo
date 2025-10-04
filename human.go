package chronogo

import (
	"fmt"
	"math"
	"time"
)

// DiffForHumans returns a human-readable string describing the difference
// between this DateTime and another DateTime or the current time.
// Uses the default locale (set via SetDefaultLocale). Defaults to English.
//
// Examples:
//   - English: "2 hours ago", "in 3 days"
//   - Spanish: "hace 2 horas", "en 3 días"
//   - Japanese: "2時間前", "3日後"
func (dt DateTime) DiffForHumans(other ...DateTime) string {
	var reference DateTime
	if len(other) > 0 {
		reference = other[0]
	} else {
		reference = Now()
	}

	// Use default locale
	locale, err := GetLocale(defaultLocale)
	if err != nil {
		// Fallback to English if default locale fails
		locale, _ = GetLocale("en-US")
	}

	return dt.humanStringWithLocale(reference, locale)
}

// DiffForHumansNow returns a human-readable string describing the difference
// between this DateTime and the current time.
// Uses the default locale.
func (dt DateTime) DiffForHumansNow() string {
	return dt.DiffForHumans()
}

// Humanize returns a human-readable representation of a duration.
// Uses the default locale for time unit names.
//
// Note: This provides a simple duration representation. For relative time
// differences (e.g., "2 hours ago"), use DiffForHumans() instead.
func Humanize(duration time.Duration) string {
	if duration == 0 {
		locale, _ := GetLocale(defaultLocale)
		if locale == nil {
			locale, _ = GetLocale("en-US")
		}
		// Return "0 seconds" in the appropriate language
		if units, ok := locale.TimeUnits["second"]; ok {
			return fmt.Sprintf("0 %s", units.Plural)
		}
		return "0 seconds"
	}

	// Get locale for unit names
	locale, err := GetLocale(defaultLocale)
	if err != nil {
		locale, _ = GetLocale("en-US")
	}

	absDuration := time.Duration(math.Abs(float64(duration)))
	
	// Calculate the appropriate unit
	seconds := int(absDuration.Seconds())
	minutes := int(absDuration.Minutes())
	hours := int(absDuration.Hours())
	days := int(absDuration.Hours() / 24)
	weeks := days / 7
	months := days / 30
	years := days / 365

	var unit string
	var value int

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
	}

	// Get localized unit name
	unitNames, exists := locale.TimeUnits[unit]
	if !exists {
		// Fallback to English
		if value == 1 {
			return fmt.Sprintf("%d %s", value, unit)
		}
		return fmt.Sprintf("%d %ss", value, unit)
	}

	unitName := unitNames.Singular
	if value != 1 {
		unitName = unitNames.Plural
	}

	if duration < 0 {
		return fmt.Sprintf("-%d %s", value, unitName)
	}

	return fmt.Sprintf("%d %s", value, unitName)
}

// Age returns the age of the DateTime compared to now.
// Uses the default locale for output.
//
// Note: This is a simplified age calculation. For more precise
// age calculations with years, months, and days, use Diff().
func (dt DateTime) Age() string {
	now := Now()
	if dt.After(now) {
		// TODO: Localize "not yet born"
		return "not yet born"
	}

	// Use DiffForHumans which is already locale-aware
	// But format it as "X old" instead of "X ago"
	duration := now.Sub(dt)
	years := int(duration.Hours() / 24 / 365.25)

	locale, err := GetLocale(defaultLocale)
	if err != nil {
		locale, _ = GetLocale("en-US")
	}

	if years == 0 {
		months := int(duration.Hours() / 24 / 30.44)
		if months == 0 {
			days := int(duration.Hours() / 24)
			if units, ok := locale.TimeUnits["day"]; ok {
				unitName := units.Singular
				if days != 1 {
					unitName = units.Plural
				}
				// TODO: Localize "old"
				return fmt.Sprintf("%d %s old", days, unitName)
			}
			return fmt.Sprintf("%d days old", days)
		}
		if units, ok := locale.TimeUnits["month"]; ok {
			unitName := units.Singular
			if months != 1 {
				unitName = units.Plural
			}
			return fmt.Sprintf("%d %s old", months, unitName)
		}
		return fmt.Sprintf("%d months old", months)
	}

	if units, ok := locale.TimeUnits["year"]; ok {
		unitName := units.Singular
		if years != 1 {
			unitName = units.Plural
		}
		return fmt.Sprintf("%d %s old", years, unitName)
	}
	return fmt.Sprintf("%d years old", years)
}

// TimeFromNow returns a human-readable string representing when this DateTime
// will occur relative to now. Uses the default locale.
func (dt DateTime) TimeFromNow() string {
	return dt.DiffForHumans()
}

// TimeAgo returns a human-readable string representing how long ago this DateTime occurred.
// Uses the default locale.
func (dt DateTime) TimeAgo() string {
	return dt.DiffForHumans()
}

// DiffForHumansComparison returns a human-readable string describing the difference
// between this DateTime and another using "before/after" format for explicit comparisons.
// Uses the default locale.
//
// Examples:
//   - English: "2 hours before", "3 days after"
//   - Spanish: "hace 2 horas", "en 3 días" (same as DiffForHumans in Spanish)
//   - Japanese: "2時間前", "3日後" (same as DiffForHumans in Japanese)
func (dt DateTime) DiffForHumansComparison(other DateTime) string {
	// Note: Some languages don't distinguish between "ago/in" and "before/after"
	// They use the same patterns, so we just use the standard human string
	locale, err := GetLocale(defaultLocale)
	if err != nil {
		locale, _ = GetLocale("en-US")
	}

	return dt.humanStringWithLocale(other, locale)
}
