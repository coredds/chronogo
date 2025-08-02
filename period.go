package chronogo

import (
	"context"
	"fmt"
	"time"
)

// Period represents a time interval between two DateTime instances.
// It provides iteration capabilities and human-friendly representations.
type Period struct {
	Start DateTime
	End   DateTime
}

// NewPeriod creates a new Period between two DateTime instances.
func NewPeriod(start, end DateTime) Period {
	return Period{Start: start, End: end}
}

// Duration returns the time.Duration of the period.
func (p Period) Duration() time.Duration {
	return p.End.Sub(p.Start)
}

// Contains checks if a DateTime falls within the period.
func (p Period) Contains(dt DateTime) bool {
	return !dt.Before(p.Start) && !dt.After(p.End)
}

// IsNegative returns true if the period represents a negative duration (end before start).
func (p Period) IsNegative() bool {
	return p.End.Before(p.Start)
}

// Abs returns a new Period with positive duration.
func (p Period) Abs() Period {
	if p.IsNegative() {
		return Period{Start: p.End, End: p.Start}
	}
	return p
}

// Years returns the number of full years in the period.
func (p Period) Years() int {
	if p.IsNegative() {
		return -p.Abs().Years()
	}

	years := p.End.Year() - p.Start.Year()

	// Adjust if we haven't reached the anniversary
	endMonth := p.End.Month()
	startMonth := p.Start.Month()
	endDay := p.End.Day()
	startDay := p.Start.Day()

	if endMonth < startMonth || (endMonth == startMonth && endDay < startDay) {
		years--
	}

	return years
}

// Months returns the total number of full months in the period.
func (p Period) Months() int {
	if p.IsNegative() {
		return -p.Abs().Months()
	}

	years := p.Years()
	months := years * 12

	// Add remaining months
	endMonth := int(p.End.Month())
	startMonth := int(p.Start.Month())
	months += endMonth - startMonth

	// Adjust if we haven't reached the day anniversary
	if p.End.Day() < p.Start.Day() {
		months--
	}

	return months
}

// Days returns the number of full days in the period.
func (p Period) Days() int {
	duration := p.Duration()
	if duration < 0 {
		return -int((-duration).Hours() / 24)
	}
	return int(duration.Hours() / 24)
}

// Hours returns the number of full hours in the period.
func (p Period) Hours() int {
	duration := p.Duration()
	if duration < 0 {
		return -int((-duration).Hours())
	}
	return int(duration.Hours())
}

// Minutes returns the number of full minutes in the period.
func (p Period) Minutes() int {
	duration := p.Duration()
	if duration < 0 {
		return -int((-duration).Minutes())
	}
	return int(duration.Minutes())
}

// Seconds returns the number of full seconds in the period.
func (p Period) Seconds() int {
	duration := p.Duration()
	if duration < 0 {
		return -int((-duration).Seconds())
	}
	return int(duration.Seconds())
}

// InDays returns the total period expressed in days as a float.
func (p Period) InDays() float64 {
	return p.Duration().Hours() / 24
}

// InHours returns the total period expressed in hours as a float.
func (p Period) InHours() float64 {
	return p.Duration().Hours()
}

// InMinutes returns the total period expressed in minutes as a float.
func (p Period) InMinutes() float64 {
	return p.Duration().Minutes()
}

// InSeconds returns the total period expressed in seconds as a float.
func (p Period) InSeconds() float64 {
	return p.Duration().Seconds()
}

// String returns a string representation of the period.
func (p Period) String() string {
	duration := p.Duration()
	if duration == 0 {
		return "0 seconds"
	}

	if p.IsNegative() {
		return fmt.Sprintf("-%s", p.Abs().String())
	}

	parts := []string{}

	days := int(duration.Hours() / 24)
	if days > 0 {
		if days == 1 {
			parts = append(parts, "1 day")
		} else {
			parts = append(parts, fmt.Sprintf("%d days", days))
		}
		duration -= time.Duration(days) * 24 * time.Hour
	}

	hours := int(duration.Hours())
	if hours > 0 {
		if hours == 1 {
			parts = append(parts, "1 hour")
		} else {
			parts = append(parts, fmt.Sprintf("%d hours", hours))
		}
		duration -= time.Duration(hours) * time.Hour
	}

	minutes := int(duration.Minutes())
	if minutes > 0 {
		if minutes == 1 {
			parts = append(parts, "1 minute")
		} else {
			parts = append(parts, fmt.Sprintf("%d minutes", minutes))
		}
		duration -= time.Duration(minutes) * time.Minute
	}

	seconds := int(duration.Seconds())
	if seconds > 0 || len(parts) == 0 {
		if seconds == 1 {
			parts = append(parts, "1 second")
		} else {
			parts = append(parts, fmt.Sprintf("%d seconds", seconds))
		}
	}

	result := ""
	for i, part := range parts {
		if i == 0 {
			result = part
		} else if i == len(parts)-1 {
			result += " and " + part
		} else {
			result += ", " + part
		}
	}

	return result
}

// Range returns a channel that yields DateTime instances within the period.
// The step parameter determines the unit: "days", "hours", "minutes", "seconds".
func (p Period) Range(unit string, step ...int) <-chan DateTime {
	return p.RangeWithContext(context.Background(), unit, step...)
}

// RangeWithContext returns a channel that yields DateTime instances within the period with context cancellation.
// The step parameter determines the unit: "years", "months", "days", "hours", "minutes", "seconds".
// This method provides memory-safe iteration by respecting context cancellation and preventing goroutine leaks.
func (p Period) RangeWithContext(ctx context.Context, unit string, step ...int) <-chan DateTime {
	stepSize := 1
	if len(step) > 0 {
		stepSize = step[0]
	}

	ch := make(chan DateTime)

	go func() {
		defer close(ch)

		current := p.Start

		for !current.After(p.End) {
			select {
			case <-ctx.Done():
				return // Context cancelled, stop iteration
			case ch <- current:
				// Successfully sent, continue
			}

			switch unit {
			case "years":
				current = current.AddYears(stepSize)
			case "months":
				current = current.AddMonths(stepSize)
			case "days":
				current = current.AddDays(stepSize)
			case "hours":
				current = current.AddHours(stepSize)
			case "minutes":
				current = current.AddMinutes(stepSize)
			case "seconds":
				current = current.AddSeconds(stepSize)
			default:
				return // Invalid unit
			}
		}
	}()

	return ch
}

// RangeDays is a convenience method for ranging by days.
func (p Period) RangeDays(step ...int) <-chan DateTime {
	return p.Range("days", step...)
}

// RangeHours is a convenience method for ranging by hours.
func (p Period) RangeHours(step ...int) <-chan DateTime {
	return p.Range("hours", step...)
}

// ForEach iterates over the period with the given unit and step, calling fn for each DateTime.
func (p Period) ForEach(unit string, step int, fn func(DateTime)) {
	for dt := range p.Range(unit, step) {
		fn(dt)
	}
}
