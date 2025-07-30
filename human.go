package chronogo

import (
	"fmt"
	"math"
	"time"
)

// DiffForHumans returns a human-readable string describing the difference
// between this DateTime and another DateTime or the current time.
func (dt DateTime) DiffForHumans(other ...DateTime) string {
	var reference DateTime
	if len(other) > 0 {
		reference = other[0]
	} else {
		reference = Now()
	}

	return diffForHumans(dt, reference, false) // Always use ago/in format for DiffForHumans
}

// DiffForHumansNow returns a human-readable string describing the difference
// between this DateTime and the current time.
func (dt DateTime) DiffForHumansNow() string {
	return dt.DiffForHumans()
}

// diffForHumans implements the core logic for human-readable time differences.
func diffForHumans(dt, reference DateTime, isComparison bool) string {
	duration := dt.Sub(reference)
	absDuration := time.Duration(math.Abs(float64(duration)))

	// Determine if it's past or future
	isPast := duration < 0

	// Get the appropriate unit and value
	unit, value := getHumanTimeUnit(absDuration)

	// Handle special case for very small durations
	if value == 0 {
		if isComparison {
			if isPast {
				return "a few seconds before"
			}
			return "a few seconds after"
		} else {
			if isPast {
				return "a few seconds ago"
			}
			return "in a few seconds"
		}
	}

	// Format the string based on context
	if isComparison {
		if isPast {
			return fmt.Sprintf("%d %s before", value, unit)
		}
		return fmt.Sprintf("%d %s after", value, unit)
	} else {
		if isPast {
			return fmt.Sprintf("%d %s ago", value, unit)
		}
		return fmt.Sprintf("in %d %s", value, unit)
	}
}

// getHumanTimeUnit determines the most appropriate unit for human display.
func getHumanTimeUnit(duration time.Duration) (string, int) {
	seconds := int(duration.Seconds())
	minutes := int(duration.Minutes())
	hours := int(duration.Hours())
	days := int(duration.Hours() / 24)
	weeks := days / 7
	months := days / 30 // Approximate
	years := days / 365 // Approximate

	switch {
	case years > 0:
		unit := "year"
		if years != 1 {
			unit = "years"
		}
		return unit, years

	case months > 0:
		unit := "month"
		if months != 1 {
			unit = "months"
		}
		return unit, months

	case weeks > 0:
		unit := "week"
		if weeks != 1 {
			unit = "weeks"
		}
		return unit, weeks

	case days > 0:
		unit := "day"
		if days != 1 {
			unit = "days"
		}
		return unit, days

	case hours > 0:
		unit := "hour"
		if hours != 1 {
			unit = "hours"
		}
		return unit, hours

	case minutes > 0:
		unit := "minute"
		if minutes != 1 {
			unit = "minutes"
		}
		return unit, minutes

	case seconds >= 10:
		return "seconds", seconds

	default:
		return "seconds", 0 // "a few seconds"
	}
}

// Humanize returns a human-readable representation of a duration.
func Humanize(duration time.Duration) string {
	if duration == 0 {
		return "0 seconds"
	}

	absDuration := time.Duration(math.Abs(float64(duration)))
	unit, value := getHumanTimeUnit(absDuration)

	if value == 0 {
		if duration < 0 {
			return "-a few seconds"
		}
		return "a few seconds"
	}

	if duration < 0 {
		return fmt.Sprintf("-%d %s", value, unit)
	}

	return fmt.Sprintf("%d %s", value, unit)
}

// Age returns the age of the DateTime compared to now.
func (dt DateTime) Age() string {
	now := Now()
	if dt.After(now) {
		return "not yet born"
	}

	duration := now.Sub(dt)
	years := int(duration.Hours() / 24 / 365.25) // More accurate calculation

	if years == 0 {
		months := int(duration.Hours() / 24 / 30.44) // Average month length
		if months == 0 {
			days := int(duration.Hours() / 24)
			if days == 1 {
				return "1 day old"
			}
			return fmt.Sprintf("%d days old", days)
		}
		if months == 1 {
			return "1 month old"
		}
		return fmt.Sprintf("%d months old", months)
	}

	if years == 1 {
		return "1 year old"
	}
	return fmt.Sprintf("%d years old", years)
}

// TimeFromNow returns a human-readable string representing when this DateTime
// will occur relative to now.
func (dt DateTime) TimeFromNow() string {
	now := Now()
	if dt.Before(now) {
		return dt.DiffForHumans()
	}
	return dt.DiffForHumans()
}

// TimeAgo returns a human-readable string representing how long ago this DateTime occurred.
func (dt DateTime) TimeAgo() string {
	return dt.DiffForHumans()
}

// DiffForHumansComparison returns a human-readable string describing the difference
// between this DateTime and another using "before/after" format for explicit comparisons.
func (dt DateTime) DiffForHumansComparison(other DateTime) string {
	return diffForHumans(dt, other, true) // Use before/after format for explicit comparisons
}
