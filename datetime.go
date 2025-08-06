// Package chronogo provides a Go implementation of Python's Pendulum library
// for powerful and easy-to-use datetime and timezone handling.
package chronogo

import (
	"time"
)

// DateTime wraps Go's time.Time to extend functionality while maintaining compatibility.
// It provides timezone-aware datetime operations with a fluent API.
type DateTime struct {
	time.Time
}

// Now returns the current datetime in the local timezone.
func Now() DateTime {
	return DateTime{time.Now()}
}

// NowUTC returns the current datetime in UTC timezone.
func NowUTC() DateTime {
	return DateTime{time.Now().UTC()}
}

// NowIn returns the current datetime in the specified timezone.
func NowIn(loc *time.Location) DateTime {
	return DateTime{time.Now().In(loc)}
}

// Today returns today's date at midnight in the local timezone.
func Today() DateTime {
	return Now().StartOfDay()
}

// TodayIn returns today's date at midnight in the specified timezone.
func TodayIn(loc *time.Location) DateTime {
	now := time.Now().In(loc)
	return DateTime{time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)}
}

// Tomorrow returns tomorrow's date at midnight in the local timezone.
func Tomorrow() DateTime {
	return Now().AddDays(1).StartOfDay()
}

// Yesterday returns yesterday's date at midnight in the local timezone.
func Yesterday() DateTime {
	return Now().AddDays(-1).StartOfDay()
}

// Date creates a DateTime similar to time.Date() but returns our DateTime type.
func Date(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) DateTime {
	return DateTime{time.Date(year, month, day, hour, min, sec, nsec, loc)}
}

// UTC creates a DateTime in UTC timezone.
func UTC(year int, month time.Month, day, hour, min, sec, nsec int) DateTime {
	return DateTime{time.Date(year, month, day, hour, min, sec, nsec, time.UTC)}
}

// FromUnix creates a DateTime from Unix timestamp.
func FromUnix(sec int64, nsec int64, loc *time.Location) DateTime {
	return DateTime{time.Unix(sec, nsec).In(loc)}
}

// In converts the datetime to the specified timezone.
func (dt DateTime) In(loc *time.Location) DateTime {
	return DateTime{dt.Time.In(loc)}
}

// UTC converts the datetime to UTC timezone.
func (dt DateTime) UTC() DateTime {
	return DateTime{dt.Time.UTC()}
}

// Local converts the datetime to the local timezone.
func (dt DateTime) Local() DateTime {
	return DateTime{dt.Time.Local()}
}

// Location returns the current timezone location.
func (dt DateTime) Location() *time.Location {
	return dt.Time.Location()
}

// IsDST returns whether the datetime is in daylight saving time.
func (dt DateTime) IsDST() bool {
	// A time is in DST if its zone offset is different from the standard offset.
	// We find the standard offset by checking a time in deep winter (January 1st).
	_, regularOffset := time.Date(dt.Year(), time.January, 1, 0, 0, 0, 0, dt.Location()).Zone()
	_, dstOffset := dt.Zone()
	return regularOffset != dstOffset
}

// IsUTC returns whether the datetime is in UTC timezone.
func (dt DateTime) IsUTC() bool {
	return dt.Time.Location() == time.UTC
}

// IsLocal returns whether the datetime is in the local timezone.
func (dt DateTime) IsLocal() bool {
	return dt.Time.Location() == time.Local
}

// IsLeapYear returns whether the datetime's year is a leap year.
func (dt DateTime) IsLeapYear() bool {
	year := dt.Year()
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// IsPast returns whether the datetime is in the past compared to now.
func (dt DateTime) IsPast() bool {
	return dt.Time.Before(time.Now())
}

// IsFuture returns whether the datetime is in the future compared to now.
func (dt DateTime) IsFuture() bool {
	return dt.Time.After(time.Now())
}

// AddYears adds the specified number of years.
func (dt DateTime) AddYears(years int) DateTime {
	return DateTime{dt.Time.AddDate(years, 0, 0)}
}

// AddMonths adds the specified number of months.
func (dt DateTime) AddMonths(months int) DateTime {
	return DateTime{dt.Time.AddDate(0, months, 0)}
}

// AddDays adds the specified number of days.
func (dt DateTime) AddDays(days int) DateTime {
	return DateTime{dt.Time.AddDate(0, 0, days)}
}

// AddHours adds the specified number of hours.
func (dt DateTime) AddHours(hours int) DateTime {
	return DateTime{dt.Time.Add(time.Duration(hours) * time.Hour)}
}

// AddMinutes adds the specified number of minutes.
func (dt DateTime) AddMinutes(minutes int) DateTime {
	return DateTime{dt.Time.Add(time.Duration(minutes) * time.Minute)}
}

// AddSeconds adds the specified number of seconds.
func (dt DateTime) AddSeconds(seconds int) DateTime {
	return DateTime{dt.Time.Add(time.Duration(seconds) * time.Second)}
}

// Add adds a time.Duration to the datetime.
func (dt DateTime) Add(duration time.Duration) DateTime {
	return DateTime{dt.Time.Add(duration)}
}

// SubtractYears subtracts the specified number of years.
func (dt DateTime) SubtractYears(years int) DateTime {
	return dt.AddYears(-years)
}

// SubtractMonths subtracts the specified number of months.
func (dt DateTime) SubtractMonths(months int) DateTime {
	return dt.AddMonths(-months)
}

// SubtractDays subtracts the specified number of days.
func (dt DateTime) SubtractDays(days int) DateTime {
	return dt.AddDays(-days)
}

// SubtractHours subtracts the specified number of hours.
func (dt DateTime) SubtractHours(hours int) DateTime {
	return dt.AddHours(-hours)
}

// SubtractMinutes subtracts the specified number of minutes.
func (dt DateTime) SubtractMinutes(minutes int) DateTime {
	return dt.AddMinutes(-minutes)
}

// SubtractSeconds subtracts the specified number of seconds.
func (dt DateTime) SubtractSeconds(seconds int) DateTime {
	return dt.AddSeconds(-seconds)
}

// Subtract subtracts a time.Duration from the datetime.
func (dt DateTime) Subtract(duration time.Duration) DateTime {
	return DateTime{dt.Time.Add(-duration)}
}

// Sub returns the time.Duration between two DateTime instances.
func (dt DateTime) Sub(other DateTime) time.Duration {
	return dt.Time.Sub(other.Time)
}

// SetYear returns a new DateTime with the year set to the specified value.
func (dt DateTime) SetYear(year int) DateTime {
	return DateTime{time.Date(year, dt.Month(), dt.Day(), dt.Hour(), dt.Minute(), dt.Second(), dt.Nanosecond(), dt.Location())}
}

// SetMonth returns a new DateTime with the month set to the specified value.
func (dt DateTime) SetMonth(month time.Month) DateTime {
	return DateTime{time.Date(dt.Year(), month, dt.Day(), dt.Hour(), dt.Minute(), dt.Second(), dt.Nanosecond(), dt.Location())}
}

// SetDay returns a new DateTime with the day set to the specified value.
func (dt DateTime) SetDay(day int) DateTime {
	return DateTime{time.Date(dt.Year(), dt.Month(), day, dt.Hour(), dt.Minute(), dt.Second(), dt.Nanosecond(), dt.Location())}
}

// SetHour returns a new DateTime with the hour set to the specified value.
func (dt DateTime) SetHour(hour int) DateTime {
	return DateTime{time.Date(dt.Year(), dt.Month(), dt.Day(), hour, dt.Minute(), dt.Second(), dt.Nanosecond(), dt.Location())}
}

// SetMinute returns a new DateTime with the minute set to the specified value.
func (dt DateTime) SetMinute(minute int) DateTime {
	return DateTime{time.Date(dt.Year(), dt.Month(), dt.Day(), dt.Hour(), minute, dt.Second(), dt.Nanosecond(), dt.Location())}
}

// SetSecond returns a new DateTime with the second set to the specified value.
func (dt DateTime) SetSecond(second int) DateTime {
	return DateTime{time.Date(dt.Year(), dt.Month(), dt.Day(), dt.Hour(), dt.Minute(), second, dt.Nanosecond(), dt.Location())}
}

// Before reports whether the datetime is before other.
func (dt DateTime) Before(other DateTime) bool {
	return dt.Time.Before(other.Time)
}

// After reports whether the datetime is after other.
func (dt DateTime) After(other DateTime) bool {
	return dt.Time.After(other.Time)
}

// Equal reports whether the datetime is equal to other.
func (dt DateTime) Equal(other DateTime) bool {
	return dt.Time.Equal(other.Time)
}

// ToDateString returns the date portion as a string (YYYY-MM-DD).
func (dt DateTime) ToDateString() string {
	return dt.Time.Format("2006-01-02")
}

// ToTimeString returns the time portion as a string (HH:MM:SS).
func (dt DateTime) ToTimeString() string {
	return dt.Time.Format("15:04:05")
}

// ToDateTimeString returns the datetime as a string (YYYY-MM-DD HH:MM:SS).
func (dt DateTime) ToDateTimeString() string {
	return dt.Time.Format("2006-01-02 15:04:05")
}

// ToISO8601String returns the datetime in ISO 8601 format.
func (dt DateTime) ToISO8601String() string {
	return dt.Time.Format("2006-01-02T15:04:05Z07:00")
}

// String returns the default string representation (ISO 8601 format).
func (dt DateTime) String() string {
	return dt.ToISO8601String()
}

// Format formats the datetime using Go's time format layout.
func (dt DateTime) Format(layout string) string {
	return dt.Time.Format(layout)
}

// StartOfDay returns a new DateTime set to the beginning of the day (00:00:00).
func (dt DateTime) StartOfDay() DateTime {
	return DateTime{time.Date(dt.Year(), dt.Month(), dt.Day(), 0, 0, 0, 0, dt.Location())}
}

// EndOfDay returns a new DateTime set to the end of the day (23:59:59.999999999).
func (dt DateTime) EndOfDay() DateTime {
	return DateTime{time.Date(dt.Year(), dt.Month(), dt.Day(), 23, 59, 59, 999999999, dt.Location())}
}

// StartOfMonth returns a new DateTime set to the beginning of the month (first day at 00:00:00).
func (dt DateTime) StartOfMonth() DateTime {
	return DateTime{time.Date(dt.Year(), dt.Month(), 1, 0, 0, 0, 0, dt.Location())}
}

// EndOfMonth returns a new DateTime set to the end of the month (last day at 23:59:59.999999999).
func (dt DateTime) EndOfMonth() DateTime {
	return dt.StartOfMonth().AddMonths(1).AddDays(-1).EndOfDay()
}

// StartOfWeek returns a new DateTime set to the beginning of the week (Monday at 00:00:00).
func (dt DateTime) StartOfWeek() DateTime {
	weekday := dt.Weekday()
	// In Go, Sunday = 0, Monday = 1, etc. We want Monday = 0 for ISO 8601
	daysFromMonday := (int(weekday) + 6) % 7
	startOfWeek := dt.AddDays(-daysFromMonday).StartOfDay()
	return startOfWeek
}

// EndOfWeek returns a new DateTime set to the end of the week (Sunday at 23:59:59.999999999).
func (dt DateTime) EndOfWeek() DateTime {
	return dt.StartOfWeek().AddDays(6).EndOfDay()
}

// StartOfYear returns a new DateTime set to the beginning of the year (January 1st at 00:00:00).
func (dt DateTime) StartOfYear() DateTime {
	return DateTime{time.Date(dt.Year(), time.January, 1, 0, 0, 0, 0, dt.Location())}
}

// EndOfYear returns a new DateTime set to the end of the year (December 31st at 23:59:59.999999999).
func (dt DateTime) EndOfYear() DateTime {
	return DateTime{time.Date(dt.Year(), time.December, 31, 23, 59, 59, 999999999, dt.Location())}
}

// IsWeekend returns whether the datetime falls on a weekend (Saturday or Sunday).
func (dt DateTime) IsWeekend() bool {
	weekday := dt.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// IsWeekday returns whether the datetime falls on a weekday (Monday through Friday).
func (dt DateTime) IsWeekday() bool {
	return !dt.IsWeekend()
}

// Quarter returns the quarter of the year (1-4).
func (dt DateTime) Quarter() int {
	month := int(dt.Month())
	return (month-1)/3 + 1
}

// StartOfQuarter returns a new DateTime set to the beginning of the quarter.
func (dt DateTime) StartOfQuarter() DateTime {
	quarter := dt.Quarter()
	month := time.Month((quarter-1)*3 + 1)
	return DateTime{time.Date(dt.Year(), month, 1, 0, 0, 0, 0, dt.Location())}
}

// EndOfQuarter returns a new DateTime set to the end of the quarter.
func (dt DateTime) EndOfQuarter() DateTime {
	return dt.StartOfQuarter().AddMonths(3).AddDays(-1).EndOfDay()
}

// ISOWeek returns the ISO 8601 year and week number.
// Week 1 is the first week with at least 4 days in the new year.
func (dt DateTime) ISOWeek() (year, week int) {
	return dt.Time.ISOWeek()
}

// ISOWeekYear returns the ISO 8601 year for the week containing the datetime.
func (dt DateTime) ISOWeekYear() int {
	year, _ := dt.Time.ISOWeek()
	return year
}

// ISOWeekNumber returns the ISO 8601 week number (1-53).
func (dt DateTime) ISOWeekNumber() int {
	_, week := dt.Time.ISOWeek()
	return week
}

// DayOfYear returns the day of the year (1-366).
func (dt DateTime) DayOfYear() int {
	return dt.Time.YearDay()
}
