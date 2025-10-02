package chronogo

import (
	"time"
)

// NextWeekday returns the next occurrence of the specified weekday.
// If the current day is the specified weekday, it returns the same weekday next week.
//
// Example:
//
//	dt := chronogo.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC) // Monday
//	next := dt.NextWeekday(time.Wednesday) // Returns next Wednesday
func (dt DateTime) NextWeekday(weekday time.Weekday) DateTime {
	current := dt.Weekday()
	daysToAdd := int(weekday - current)

	if daysToAdd <= 0 {
		daysToAdd += 7
	}

	return dt.AddDays(daysToAdd)
}

// PreviousWeekday returns the previous occurrence of the specified weekday.
// If the current day is the specified weekday, it returns the same weekday last week.
//
// Example:
//
//	dt := chronogo.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC) // Monday
//	prev := dt.PreviousWeekday(time.Friday) // Returns previous Friday
func (dt DateTime) PreviousWeekday(weekday time.Weekday) DateTime {
	current := dt.Weekday()
	daysToSubtract := int(current - weekday)

	if daysToSubtract <= 0 {
		daysToSubtract += 7
	}

	return dt.AddDays(-daysToSubtract)
}

// ClosestWeekday returns the closest occurrence of the specified weekday.
// If the current day is the specified weekday, it returns the current DateTime.
// If two occurrences are equidistant, it returns the future one.
//
// Example:
//
//	dt := chronogo.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC) // Monday
//	closest := dt.ClosestWeekday(time.Wednesday) // Returns nearest Wednesday
func (dt DateTime) ClosestWeekday(weekday time.Weekday) DateTime {
	current := dt.Weekday()

	// If it's already the target weekday, return as is
	if current == weekday {
		return dt
	}

	// Calculate days to next and previous occurrence
	daysToNext := int(weekday - current)
	if daysToNext < 0 {
		daysToNext += 7
	}

	daysToPrevious := int(current - weekday)
	if daysToPrevious < 0 {
		daysToPrevious += 7
	}

	// Return the closest one (prefer future if equidistant)
	if daysToNext <= daysToPrevious {
		return dt.AddDays(daysToNext)
	}
	return dt.AddDays(-daysToPrevious)
}

// FarthestWeekday returns the farthest occurrence of the specified weekday within the week.
// This returns the occurrence that is furthest from the current day (3-4 days away).
//
// Example:
//
//	dt := chronogo.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC) // Monday
//	farthest := dt.FarthestWeekday(time.Thursday) // Returns the Thursday that's ~3-4 days away
func (dt DateTime) FarthestWeekday(weekday time.Weekday) DateTime {
	current := dt.Weekday()

	// If it's already the target weekday, return next week's occurrence
	if current == weekday {
		return dt.AddDays(7)
	}

	// Calculate days to next occurrence (always positive, 1-6)
	daysToNext := int(weekday - current)
	if daysToNext <= 0 {
		daysToNext += 7
	}

	// Calculate days to previous occurrence (always positive, 1-6)
	daysToPrevious := int(current - weekday)
	if daysToPrevious <= 0 {
		daysToPrevious += 7
	}

	// Return the farthest one (prefer forward if equal)
	if daysToNext >= daysToPrevious {
		return dt.AddDays(daysToNext)
	}
	return dt.AddDays(-daysToPrevious)
}

// NextOrSameWeekday returns the next occurrence of the specified weekday,
// or the current DateTime if it's already that weekday.
//
// Example:
//
//	dt := chronogo.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC) // Monday
//	next := dt.NextOrSameWeekday(time.Monday) // Returns the same Monday
//	next := dt.NextOrSameWeekday(time.Tuesday) // Returns next Tuesday
func (dt DateTime) NextOrSameWeekday(weekday time.Weekday) DateTime {
	if dt.Weekday() == weekday {
		return dt
	}
	return dt.NextWeekday(weekday)
}

// PreviousOrSameWeekday returns the previous occurrence of the specified weekday,
// or the current DateTime if it's already that weekday.
//
// Example:
//
//	dt := chronogo.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC) // Monday
//	prev := dt.PreviousOrSameWeekday(time.Monday) // Returns the same Monday
//	prev := dt.PreviousOrSameWeekday(time.Sunday) // Returns previous Sunday
func (dt DateTime) PreviousOrSameWeekday(weekday time.Weekday) DateTime {
	if dt.Weekday() == weekday {
		return dt
	}
	return dt.PreviousWeekday(weekday)
}

// NthWeekdayOf returns the nth occurrence of the specified weekday in the given month/year.
// n can be 1-5 for first through fifth occurrence, or -1 for the last occurrence.
//
// Example:
//
//	// Get the 2nd Monday of March 2024
//	dt := chronogo.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
//	secondMonday := dt.NthWeekdayOf(2, time.Monday, "month")
//
//	// Get the last Friday of the year
//	lastFriday := dt.NthWeekdayOf(-1, time.Friday, "year")
func (dt DateTime) NthWeekdayOf(n int, weekday time.Weekday, unit string) DateTime {
	if n == 0 || n < -1 {
		return DateTime{} // Invalid n
	}

	// For month, limit to 5 occurrences; for year/quarter, allow more
	maxN := 53 // Maximum possible occurrences in a year
	if unit == "month" {
		maxN = 5
	}

	if n > maxN {
		return DateTime{} // Invalid n for unit
	}

	var start, end DateTime

	switch unit {
	case "month":
		start = dt.StartOfMonth()
		end = dt.EndOfMonth()
	case "year":
		start = dt.StartOfYear()
		end = dt.EndOfYear()
	case "quarter":
		start = dt.StartOfQuarter()
		end = dt.EndOfQuarter()
	default:
		return DateTime{} // Invalid unit
	}

	if n == -1 {
		// Find last occurrence - start from end and work backwards
		current := end
		for !current.Before(start) {
			if current.Weekday() == weekday {
				return current
			}
			current = current.AddDays(-1)
		}
		return DateTime{} // Not found
	}

	// Find nth occurrence from start
	current := start
	count := 0

	for !current.After(end) {
		if current.Weekday() == weekday {
			count++
			if count == n {
				return current
			}
		}
		current = current.AddDays(1)
	}

	return DateTime{} // Not found (month doesn't have n occurrences)
}

// FirstWeekdayOf returns the first occurrence of the specified weekday in the current month.
//
// Example:
//
//	dt := chronogo.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
//	firstMonday := dt.FirstWeekdayOf(time.Monday) // First Monday of March 2024
func (dt DateTime) FirstWeekdayOf(weekday time.Weekday) DateTime {
	return dt.NthWeekdayOf(1, weekday, "month")
}

// LastWeekdayOf returns the last occurrence of the specified weekday in the current month.
//
// Example:
//
//	dt := chronogo.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
//	lastFriday := dt.LastWeekdayOf(time.Friday) // Last Friday of March 2024
func (dt DateTime) LastWeekdayOf(weekday time.Weekday) DateTime {
	return dt.NthWeekdayOf(-1, weekday, "month")
}

// NthWeekdayOfMonth returns the nth occurrence of the specified weekday in the current month.
// This is a convenience wrapper for NthWeekdayOf with "month" unit.
//
// Example:
//
//	dt := chronogo.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
//	thirdTuesday := dt.NthWeekdayOfMonth(3, time.Tuesday) // 3rd Tuesday of March 2024
func (dt DateTime) NthWeekdayOfMonth(n int, weekday time.Weekday) DateTime {
	return dt.NthWeekdayOf(n, weekday, "month")
}

// NthWeekdayOfYear returns the nth occurrence of the specified weekday in the current year.
//
// Example:
//
//	dt := chronogo.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
//	tenthMonday := dt.NthWeekdayOfYear(10, time.Monday) // 10th Monday of 2024
func (dt DateTime) NthWeekdayOfYear(n int, weekday time.Weekday) DateTime {
	return dt.NthWeekdayOf(n, weekday, "year")
}

// IsNthWeekdayOf checks if the current DateTime is the nth occurrence of its weekday in the specified unit.
//
// Example:
//
//	dt := chronogo.Date(2024, 3, 11, 0, 0, 0, 0, time.UTC) // 2nd Monday of March
//	isSecond := dt.IsNthWeekdayOf(2, "month") // Returns true
func (dt DateTime) IsNthWeekdayOf(n int, unit string) bool {
	weekday := dt.Weekday()
	nthOccurrence := dt.NthWeekdayOf(n, weekday, unit)

	if nthOccurrence.IsZero() {
		return false
	}

	return dt.Year() == nthOccurrence.Year() &&
		dt.Month() == nthOccurrence.Month() &&
		dt.Day() == nthOccurrence.Day()
}

// WeekdayOccurrenceInMonth returns which occurrence of the weekday this is in the month (1-5).
// Returns 0 if it cannot be determined.
//
// Example:
//
//	dt := chronogo.Date(2024, 3, 11, 0, 0, 0, 0, time.UTC) // 2nd Monday of March
//	occurrence := dt.WeekdayOccurrenceInMonth() // Returns 2
func (dt DateTime) WeekdayOccurrenceInMonth() int {
	weekday := dt.Weekday()
	start := dt.StartOfMonth()

	count := 0
	current := start

	for current.Month() == dt.Month() {
		if current.Weekday() == weekday {
			count++
			if current.Day() == dt.Day() {
				return count
			}
		}
		current = current.AddDays(1)
	}

	return 0
}
