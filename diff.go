package chronogo

import (
	"fmt"
	"math"
	"time"
)

// Diff represents the difference between two DateTime instances.
// It provides both precise (Duration-based) and calendar-aware (Period-based) difference methods.
// This type unifies the functionality of time.Duration and Period into a single, convenient API.
type Diff struct {
	start    DateTime
	end      DateTime
	duration time.Duration
	period   Period
}

// Diff returns a Diff object representing the difference between two DateTimes.
// The returned Diff provides both precise duration-based and calendar-aware period-based methods.
//
// Examples:
//
//	dt1 := chronogo.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC)
//	dt2 := chronogo.Date(2024, 3, 20, 14, 30, 0, 0, time.UTC)
//	diff := dt2.Diff(dt1)
//
//	// Precise differences
//	fmt.Println(diff.InHours())  // Total hours as float
//	fmt.Println(diff.InDays())   // Total days as float
//
//	// Calendar-aware differences
//	fmt.Println(diff.Years())    // Full years
//	fmt.Println(diff.Months())   // Total full months
//
//	// Human-readable
//	fmt.Println(diff.ForHumans())
func (dt DateTime) Diff(other DateTime) Diff {
	return Diff{
		start:    other,
		end:      dt,
		duration: dt.Sub(other),
		period:   NewPeriod(other, dt),
	}
}

// DiffAbs returns the absolute difference between two DateTimes.
// The returned Diff always represents a positive duration.
func (dt DateTime) DiffAbs(other DateTime) Diff {
	diff := dt.Diff(other)
	if diff.IsNegative() {
		return diff.Abs()
	}
	return diff
}

// Duration returns the precise time.Duration between the two DateTimes.
func (d Diff) Duration() time.Duration {
	return d.duration
}

// Period returns the Period representation for calendar-aware operations.
func (d Diff) Period() Period {
	return d.period
}

// Start returns the earlier DateTime in the comparison.
func (d Diff) Start() DateTime {
	return d.start
}

// End returns the later DateTime in the comparison.
func (d Diff) End() DateTime {
	return d.end
}

// IsNegative returns true if the end DateTime is before the start DateTime.
func (d Diff) IsNegative() bool {
	return d.duration < 0
}

// IsPositive returns true if the end DateTime is after the start DateTime.
func (d Diff) IsPositive() bool {
	return d.duration > 0
}

// IsZero returns true if the two DateTimes are identical.
func (d Diff) IsZero() bool {
	return d.duration == 0
}

// Abs returns a new Diff with positive duration.
func (d Diff) Abs() Diff {
	if d.IsNegative() {
		return Diff{
			start:    d.end,
			end:      d.start,
			duration: -d.duration,
			period:   d.period.Abs(),
		}
	}
	return d
}

// Invert returns a new Diff with start and end swapped.
func (d Diff) Invert() Diff {
	return Diff{
		start:    d.end,
		end:      d.start,
		duration: -d.duration,
		period:   Period{Start: d.end, End: d.start},
	}
}

// Years returns the number of full calendar years in the difference.
// This is calendar-aware and accounts for leap years and varying month lengths.
func (d Diff) Years() int {
	return d.period.Years()
}

// Months returns the total number of full calendar months in the difference.
// This is calendar-aware and accounts for varying month lengths.
func (d Diff) Months() int {
	return d.period.Months()
}

// Weeks returns the number of full weeks in the difference.
func (d Diff) Weeks() int {
	return d.Days() / 7
}

// Days returns the number of full 24-hour days in the difference.
func (d Diff) Days() int {
	return d.period.Days()
}

// Hours returns the number of full hours in the difference.
func (d Diff) Hours() int {
	return d.period.Hours()
}

// Minutes returns the number of full minutes in the difference.
func (d Diff) Minutes() int {
	return d.period.Minutes()
}

// Seconds returns the number of full seconds in the difference.
func (d Diff) Seconds() int {
	return d.period.Seconds()
}

// InYears returns the total difference expressed as years (with fractional part).
// This uses an approximate calculation (365.25 days per year).
func (d Diff) InYears() float64 {
	return d.InDays() / 365.25
}

// InMonths returns the total difference expressed as months (with fractional part).
// This uses an approximate calculation (30.44 days per month).
func (d Diff) InMonths() float64 {
	return d.InDays() / 30.44
}

// InWeeks returns the total difference expressed as weeks (with fractional part).
func (d Diff) InWeeks() float64 {
	return d.InDays() / 7.0
}

// InDays returns the total difference expressed as days (with fractional part).
func (d Diff) InDays() float64 {
	return d.period.InDays()
}

// InHours returns the total difference expressed as hours (with fractional part).
func (d Diff) InHours() float64 {
	return d.period.InHours()
}

// InMinutes returns the total difference expressed as minutes (with fractional part).
func (d Diff) InMinutes() float64 {
	return d.period.InMinutes()
}

// InSeconds returns the total difference expressed as seconds (with fractional part).
func (d Diff) InSeconds() float64 {
	return d.period.InSeconds()
}

// ForHumans returns a human-readable string describing the difference.
// Uses the default locale (set via SetDefaultLocale). Defaults to English.
//
// Examples:
//   - English: "2 years ago", "in 3 months"
//   - Spanish: "hace 2 años", "en 3 meses"
//   - Japanese: "2年前", "3ヶ月後"
func (d Diff) ForHumans() string {
	return d.end.DiffForHumans(d.start)
}

// ForHumansLocalized returns a human-readable string in the specified locale.
func (d Diff) ForHumansLocalized(localeCode string) (string, error) {
	return d.end.HumanStringLocalized(localeCode, d.start)
}

// ForHumansComparison returns a human-readable comparison string.
// Uses the default locale.
//
// Examples:
//   - English: "2 years before", "3 months after"
//   - Spanish: "hace 2 años", "en 3 meses"
//   - Japanese: "2年前", "3ヶ月後"
func (d Diff) ForHumansComparison() string {
	return d.end.DiffForHumansComparison(d.start)
}

// String returns a detailed string representation of the difference.
func (d Diff) String() string {
	if d.IsZero() {
		return "0 seconds"
	}

	abs := d.Abs()
	sign := ""
	if d.IsNegative() {
		sign = "-"
	}

	years := abs.Years()
	months := abs.Months() % 12
	days := abs.Days() % 30 // Approximate
	hours := abs.Hours() % 24
	minutes := abs.Minutes() % 60
	seconds := abs.Seconds() % 60

	parts := []string{}

	if years > 0 {
		if years == 1 {
			parts = append(parts, "1 year")
		} else {
			parts = append(parts, fmt.Sprintf("%d years", years))
		}
	}

	if months > 0 {
		if months == 1 {
			parts = append(parts, "1 month")
		} else {
			parts = append(parts, fmt.Sprintf("%d months", months))
		}
	}

	if days > 0 {
		if days == 1 {
			parts = append(parts, "1 day")
		} else {
			parts = append(parts, fmt.Sprintf("%d days", days))
		}
	}

	if hours > 0 {
		if hours == 1 {
			parts = append(parts, "1 hour")
		} else {
			parts = append(parts, fmt.Sprintf("%d hours", hours))
		}
	}

	if minutes > 0 {
		if minutes == 1 {
			parts = append(parts, "1 minute")
		} else {
			parts = append(parts, fmt.Sprintf("%d minutes", minutes))
		}
	}

	if seconds > 0 || len(parts) == 0 {
		if seconds == 1 {
			parts = append(parts, "1 second")
		} else {
			parts = append(parts, fmt.Sprintf("%d seconds", seconds))
		}
	}

	result := ""
	if len(parts) == 1 {
		result = parts[0]
	} else if len(parts) == 2 {
		result = parts[0] + " and " + parts[1]
	} else if len(parts) > 2 {
		result = ""
		for i, part := range parts[:len(parts)-1] {
			if i > 0 {
				result += ", "
			}
			result += part
		}
		result += " and " + parts[len(parts)-1]
	}

	return sign + result
}

// CompactString returns a compact string representation (e.g., "2y 3m 5d").
func (d Diff) CompactString() string {
	if d.IsZero() {
		return "0s"
	}

	abs := d.Abs()
	sign := ""
	if d.IsNegative() {
		sign = "-"
	}

	years := abs.Years()
	months := abs.Months() % 12
	days := abs.Days() % 30
	hours := abs.Hours() % 24
	minutes := abs.Minutes() % 60
	seconds := abs.Seconds() % 60

	parts := []string{}

	if years > 0 {
		parts = append(parts, fmt.Sprintf("%dy", years))
	}
	if months > 0 {
		parts = append(parts, fmt.Sprintf("%dmo", months))
	}
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}
	if seconds > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%ds", seconds))
	}

	result := ""
	for i, part := range parts {
		if i > 0 {
			result += " "
		}
		result += part
	}

	return sign + result
}

// Compare compares this Diff with another Diff and returns:
//   - -1 if this Diff is shorter
//   - 0 if they are equal
//   - 1 if this Diff is longer
func (d Diff) Compare(other Diff) int {
	d1 := math.Abs(float64(d.duration))
	d2 := math.Abs(float64(other.duration))

	if d1 < d2 {
		return -1
	} else if d1 > d2 {
		return 1
	}
	return 0
}

// LongerThan returns true if this Diff is longer than the other Diff.
func (d Diff) LongerThan(other Diff) bool {
	return d.Compare(other) > 0
}

// ShorterThan returns true if this Diff is shorter than the other Diff.
func (d Diff) ShorterThan(other Diff) bool {
	return d.Compare(other) < 0
}

// EqualTo returns true if this Diff has the same duration as the other Diff.
func (d Diff) EqualTo(other Diff) bool {
	return d.Compare(other) == 0
}

// MarshalJSON implements json.Marshaler.
func (d Diff) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"duration":%d,"start":%q,"end":%q}`,
		d.duration.Milliseconds(),
		d.start.Format(time.RFC3339Nano),
		d.end.Format(time.RFC3339Nano))), nil
}
