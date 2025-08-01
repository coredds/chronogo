package chronogo

import (
	"time"
)

// FluentDateTime provides a fluent API for building and manipulating DateTime instances.
// It allows method chaining for more readable and expressive date/time operations.
type FluentDateTime struct {
	base DateTime
}

// FluentDuration provides a fluent API for building durations with human-readable methods.
type FluentDuration struct {
	duration time.Duration
}

// AddFluent returns a FluentDuration for adding time units to the DateTime.
func (dt DateTime) AddFluent() *FluentDuration {
	return &FluentDuration{duration: 0}
}

// Set returns a FluentDateTime for setting specific components of the DateTime.
func (dt DateTime) Set() *FluentDateTime {
	return &FluentDateTime{base: dt}
}

// Years adds the specified number of years to the duration.
func (fd *FluentDuration) Years(years int) *FluentDuration {
	// For fluent duration, we need to approximate years as days
	// This is approximate since years have different numbers of days
	fd.duration += time.Duration(years) * 365 * 24 * time.Hour
	return fd
}

// Months adds the specified number of months to the duration.
func (fd *FluentDuration) Months(months int) *FluentDuration {
	// Approximate months as 30 days
	fd.duration += time.Duration(months) * 30 * 24 * time.Hour
	return fd
}

// Weeks adds the specified number of weeks to the duration.
func (fd *FluentDuration) Weeks(weeks int) *FluentDuration {
	fd.duration += time.Duration(weeks) * 7 * 24 * time.Hour
	return fd
}

// Days adds the specified number of days to the duration.
func (fd *FluentDuration) Days(days int) *FluentDuration {
	fd.duration += time.Duration(days) * 24 * time.Hour
	return fd
}

// Hours adds the specified number of hours to the duration.
func (fd *FluentDuration) Hours(hours int) *FluentDuration {
	fd.duration += time.Duration(hours) * time.Hour
	return fd
}

// Minutes adds the specified number of minutes to the duration.
func (fd *FluentDuration) Minutes(minutes int) *FluentDuration {
	fd.duration += time.Duration(minutes) * time.Minute
	return fd
}

// Seconds adds the specified number of seconds to the duration.
func (fd *FluentDuration) Seconds(seconds int) *FluentDuration {
	fd.duration += time.Duration(seconds) * time.Second
	return fd
}

// Milliseconds adds the specified number of milliseconds to the duration.
func (fd *FluentDuration) Milliseconds(milliseconds int) *FluentDuration {
	fd.duration += time.Duration(milliseconds) * time.Millisecond
	return fd
}

// Microseconds adds the specified number of microseconds to the duration.
func (fd *FluentDuration) Microseconds(microseconds int) *FluentDuration {
	fd.duration += time.Duration(microseconds) * time.Microsecond
	return fd
}

// Nanoseconds adds the specified number of nanoseconds to the duration.
func (fd *FluentDuration) Nanoseconds(nanoseconds int) *FluentDuration {
	fd.duration += time.Duration(nanoseconds) * time.Nanosecond
	return fd
}

// To applies the accumulated duration to a DateTime and returns the result.
func (fd *FluentDuration) To(dt DateTime) DateTime {
	return dt.Add(fd.duration)
}

// From subtracts the accumulated duration from a DateTime and returns the result.
func (fd *FluentDuration) From(dt DateTime) DateTime {
	return dt.Subtract(fd.duration)
}

// Year sets the year component.
func (fdt *FluentDateTime) Year(year int) *FluentDateTime {
	fdt.base = fdt.base.SetYear(year)
	return fdt
}

// Month sets the month component.
func (fdt *FluentDateTime) Month(month time.Month) *FluentDateTime {
	fdt.base = fdt.base.SetMonth(month)
	return fdt
}

// Day sets the day component.
func (fdt *FluentDateTime) Day(day int) *FluentDateTime {
	fdt.base = fdt.base.SetDay(day)
	return fdt
}

// Hour sets the hour component.
func (fdt *FluentDateTime) Hour(hour int) *FluentDateTime {
	fdt.base = fdt.base.SetHour(hour)
	return fdt
}

// Minute sets the minute component.
func (fdt *FluentDateTime) Minute(minute int) *FluentDateTime {
	fdt.base = fdt.base.SetMinute(minute)
	return fdt
}

// Second sets the second component.
func (fdt *FluentDateTime) Second(second int) *FluentDateTime {
	fdt.base = fdt.base.SetSecond(second)
	return fdt
}

// Timezone sets the timezone.
func (fdt *FluentDateTime) Timezone(loc *time.Location) *FluentDateTime {
	fdt.base = fdt.base.In(loc)
	return fdt
}

// Build returns the final DateTime with all modifications applied.
func (fdt *FluentDateTime) Build() DateTime {
	return fdt.base
}

// ChronoDuration extends Go's time.Duration with human-readable formatting and additional operations.
type ChronoDuration struct {
	time.Duration
}

// NewDuration creates a new ChronoDuration from a time.Duration.
func NewDuration(d time.Duration) ChronoDuration {
	return ChronoDuration{d}
}

// NewDurationFromComponents creates a ChronoDuration from individual components.
func NewDurationFromComponents(hours, minutes, seconds int) ChronoDuration {
	d := time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second
	return ChronoDuration{d}
}

// Years returns the approximate number of years in the duration.
func (cd ChronoDuration) Years() float64 {
	return cd.Hours() / (365.25 * 24) // Account for leap years
}

// Months returns the approximate number of months in the duration.
func (cd ChronoDuration) Months() float64 {
	return cd.Hours() / (30.44 * 24) // Average month length
}

// Weeks returns the number of weeks in the duration.
func (cd ChronoDuration) Weeks() float64 {
	return cd.Hours() / (7 * 24)
}

// Days returns the number of days in the duration.
func (cd ChronoDuration) Days() float64 {
	return cd.Hours() / 24
}

// HumanString returns a human-readable representation of the duration.
func (cd ChronoDuration) HumanString() string {
	return Humanize(cd.Duration)
}

// String returns a string representation of the duration.
func (cd ChronoDuration) String() string {
	return cd.Duration.String()
}

// Add adds another duration to this one.
func (cd ChronoDuration) Add(other ChronoDuration) ChronoDuration {
	return ChronoDuration{cd.Duration + other.Duration}
}

// Subtract subtracts another duration from this one.
func (cd ChronoDuration) Subtract(other ChronoDuration) ChronoDuration {
	return ChronoDuration{cd.Duration - other.Duration}
}

// Multiply multiplies the duration by a factor.
func (cd ChronoDuration) Multiply(factor float64) ChronoDuration {
	return ChronoDuration{time.Duration(float64(cd.Duration) * factor)}
}

// Divide divides the duration by a factor.
func (cd ChronoDuration) Divide(factor float64) ChronoDuration {
	return ChronoDuration{time.Duration(float64(cd.Duration) / factor)}
}

// IsPositive returns true if the duration is positive.
func (cd ChronoDuration) IsPositive() bool {
	return cd.Duration > 0
}

// IsNegative returns true if the duration is negative.
func (cd ChronoDuration) IsNegative() bool {
	return cd.Duration < 0
}

// IsZero returns true if the duration is zero.
func (cd ChronoDuration) IsZero() bool {
	return cd.Duration == 0
}

// Abs returns the absolute value of the duration.
func (cd ChronoDuration) Abs() ChronoDuration {
	if cd.Duration < 0 {
		return ChronoDuration{-cd.Duration}
	}
	return cd
}
