// Package chronogo provides a Go implementation of powerful datetime handling
// for easy-to-use datetime and timezone operations.
package chronogo

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Unit represents a logical time unit used by various helpers.
type Unit int

const (
	UnitSecond Unit = iota
	UnitMinute
	UnitHour
	UnitDay
	UnitWeek
	UnitMonth
	UnitQuarter
	UnitYear
)

// DateTime wraps Go's time.Time to extend functionality while maintaining compatibility.
// It provides timezone-aware datetime operations with a fluent API.
type DateTime struct {
	time.Time
}

// cache for standard (non-DST) timezone offsets per (location, year)
var standardOffsetCache sync.Map // key string -> int (seconds east of UTC)

func standardOffsetKey(loc *time.Location, year int) string {
	// loc.String() returns IANA name (e.g., "America/New_York") or "Local"
	return loc.String() + "|" + strconv.Itoa(year)
}

func getStandardOffset(loc *time.Location, year int) int {
	key := standardOffsetKey(loc, year)
	if v, ok := standardOffsetCache.Load(key); ok {
		return v.(int)
	}
	// Compute minimum offset observed across the year as standard offset
	minOffset := int(1<<31 - 1)
	for month := time.January; month <= time.December; month++ {
		t := time.Date(year, month, 1, 0, 0, 0, 0, loc)
		_, off := t.Zone()
		if off < minOffset {
			minOffset = off
		}
	}
	standardOffsetCache.Store(key, minOffset)
	return minOffset
}

// Now returns the current datetime in the local timezone.
// When testing helpers are active (SetTestNow, FreezeTime, TravelTo),
// this will return the mocked time instead of the actual current time.
func Now() DateTime {
	return DateTime{getTestableNow()}
}

// NowUTC returns the current datetime in UTC timezone.
// When testing helpers are active, this will return the mocked time in UTC.
func NowUTC() DateTime {
	return DateTime{getTestableNow().UTC()}
}

// NowIn returns the current datetime in the specified timezone.
// When testing helpers are active, this will return the mocked time in the specified timezone.
func NowIn(loc *time.Location) DateTime {
	return DateTime{getTestableNow().In(loc)}
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
	// Determine standard (non-DST) offset via cached minimum offset across the year.
	loc := dt.Location()
	year := dt.Year()
	minOffset := getStandardOffset(loc, year)
	_, currentOffset := dt.Zone()
	return currentOffset != minOffset
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

// IsLongYear returns whether the datetime's year is an ISO 8601 long year.
// A long year has 53 ISO weeks instead of the normal 52.
// This occurs when the year starts on a Thursday or is a leap year starting on Wednesday.
//
// Examples:
//
//	chronogo.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).IsLongYear() // true - 2020 has 53 ISO weeks
//	chronogo.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).IsLongYear() // false - 2021 has 52 ISO weeks
func (dt DateTime) IsLongYear() bool {
	year := dt.Year()

	// A year is a long year if it has 53 ISO weeks.
	// This happens if:
	// 1. January 1st is a Thursday, OR
	// 2. January 1st is a Wednesday AND it's a leap year
	//
	// Equivalently, we can check if December 28 (which is always in the last week)
	// is in week 53 of that year.
	dec28 := time.Date(year, 12, 28, 0, 0, 0, 0, dt.Location())
	_, week := dec28.ISOWeek()

	return week == 53
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

// UnixMilli returns t as a Unix time, the number of milliseconds elapsed
// since January 1, 1970 UTC.
func (dt DateTime) UnixMilli() int64 {
	return dt.Time.UnixMilli()
}

// UnixMicro returns t as a Unix time, the number of microseconds elapsed
// since January 1, 1970 UTC.
func (dt DateTime) UnixMicro() int64 {
	return dt.Time.UnixMicro()
}

// UnixNano returns t as a Unix time, the number of nanoseconds elapsed
// since January 1, 1970 UTC.
func (dt DateTime) UnixNano() int64 {
	return dt.Time.UnixNano()
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

// On is a convenience method that sets the date components (year, month, day) in one call.
// This is a simpler alternative to using Set().Year().Month().Day() for date modifications.
// The time components (hour, minute, second, nanosecond) remain unchanged.
//
// Examples:
//
//	dt := chronogo.Now()
//	dt2 := dt.On(2024, time.January, 15)  // Sets date to 2024-01-15, keeps current time
//	dt3 := dt.On(2025, time.December, 31) // Sets date to 2025-12-31, keeps current time
func (dt DateTime) On(year int, month time.Month, day int) DateTime {
	return DateTime{time.Date(year, month, day, dt.Hour(), dt.Minute(), dt.Second(), dt.Nanosecond(), dt.Location())}
}

// At is a convenience method that sets the time components (hour, minute, second) in one call.
// This is a simpler alternative to using Set().Hour().Minute().Second() for time modifications.
// The date components (year, month, day) and nanoseconds remain unchanged.
//
// Examples:
//
//	dt := chronogo.Today()
//	dt2 := dt.At(14, 30, 0)    // Sets time to 14:30:00, keeps current date
//	dt3 := dt.At(9, 0, 0)      // Sets time to 09:00:00, keeps current date
//	dt4 := dt.At(23, 59, 59)   // Sets time to 23:59:59, keeps current date
func (dt DateTime) At(hour, minute, second int) DateTime {
	return DateTime{time.Date(dt.Year(), dt.Month(), dt.Day(), hour, minute, second, dt.Nanosecond(), dt.Location())}
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

// ToCookieString returns the datetime in HTTP cookie format (RFC 1123).
// Example: "Mon, 15 Jan 2024 12:00:00 GMT"
func (dt DateTime) ToCookieString() string {
	return dt.UTC().Format(time.RFC1123)
}

// ToRSSString returns the datetime in RSS feed format (RFC 1123Z).
// Example: "Mon, 15 Jan 2024 12:00:00 +0000"
func (dt DateTime) ToRSSString() string {
	return dt.Format(time.RFC1123Z)
}

// ToW3CString returns the datetime in W3C format (ISO 8601 / RFC 3339).
// Example: "2024-01-15T12:00:00Z" or "2024-01-15T12:00:00+00:00"
func (dt DateTime) ToW3CString() string {
	return dt.Format(time.RFC3339)
}

// ToAtomString returns the datetime in Atom feed format (RFC 3339).
// This is an alias for ToW3CString().
func (dt DateTime) ToAtomString() string {
	return dt.ToW3CString()
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

// IsZero reports whether the time instant is January 1, year 1, 00:00:00 UTC.
func (dt DateTime) IsZero() bool {
	return dt.Time.IsZero()
}

// Unwrap returns the underlying time.Time value.
func (dt DateTime) Unwrap() time.Time {
	return dt.Time
}

// Truncate returns dt truncated to the start of the given unit.
// For calendar units (day/week/month/quarter/year) this aligns to the logical
// start boundary in the current location (e.g., StartOfDay, Monday StartOfWeek).
func (dt DateTime) Truncate(unit Unit) DateTime {
	switch unit {
	case UnitSecond:
		return DateTime{time.Date(dt.Year(), dt.Month(), dt.Day(), dt.Hour(), dt.Minute(), dt.Second(), 0, dt.Location())}
	case UnitMinute:
		return DateTime{time.Date(dt.Year(), dt.Month(), dt.Day(), dt.Hour(), dt.Minute(), 0, 0, dt.Location())}
	case UnitHour:
		return DateTime{time.Date(dt.Year(), dt.Month(), dt.Day(), dt.Hour(), 0, 0, 0, dt.Location())}
	case UnitDay:
		return dt.StartOfDay()
	case UnitWeek:
		return dt.StartOfWeek()
	case UnitMonth:
		return dt.StartOfMonth()
	case UnitQuarter:
		return dt.StartOfQuarter()
	case UnitYear:
		return dt.StartOfYear()
	default:
		return dt
	}
}

// Round returns dt rounded to the nearest boundary of the given unit.
// Ties are rounded up to the next boundary.
// Calendar-aware for day/week/month/quarter/year using local timezone boundaries.
func (dt DateTime) Round(unit Unit) DateTime {
	start := dt.Truncate(unit)

	var next DateTime
	switch unit {
	case UnitSecond:
		next = start.AddSeconds(1)
	case UnitMinute:
		next = start.AddMinutes(1)
	case UnitHour:
		next = start.AddHours(1)
	case UnitDay:
		next = start.AddDays(1)
	case UnitWeek:
		next = start.AddDays(7)
	case UnitMonth:
		next = start.AddMonths(1)
	case UnitQuarter:
		next = start.AddMonths(3)
	case UnitYear:
		next = start.AddYears(1)
	default:
		return dt
	}

	// Use duration between boundaries to decide rounding
	toStart := dt.Sub(start)
	boundary := next.Sub(start)
	if toStart*2 < boundary {
		return start
	}
	return next
}

// Clamp returns dt clamped to the [min, max] range (order-agnostic).
func (dt DateTime) Clamp(a, b DateTime) DateTime {
	min := a
	max := b
	if b.Before(a) {
		min, max = b, a
	}
	if dt.Before(min) {
		return min
	}
	if dt.After(max) {
		return max
	}
	return dt
}

// Between reports whether dt is within the range (a, b) or [a, b] depending on inclusive.
// The order of a and b does not matter.
func (dt DateTime) Between(a, b DateTime, inclusive bool) bool {
	min := a
	max := b
	if b.Before(a) {
		min, max = b, a
	}
	if inclusive {
		return !dt.Before(min) && !dt.After(max)
	}
	return dt.After(min) && dt.Before(max)
}

// MarshalText implements encoding.TextMarshaler.
func (dt DateTime) MarshalText() ([]byte, error) {
	return []byte(dt.ToISO8601String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (dt *DateTime) UnmarshalText(data []byte) error {
	s := strings.TrimSpace(string(data))
	if s == "" {
		*dt = DateTime{}
		return nil
	}
	parsed, err := Parse(s)
	if err != nil {
		return err
	}
	*dt = parsed
	return nil
}

// MarshalJSON implements json.Marshaler.
func (dt DateTime) MarshalJSON() ([]byte, error) {
	// Quote the ISO 8601 string
	return []byte(fmt.Sprintf("\"%s\"", dt.ToISO8601String())), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (dt *DateTime) UnmarshalJSON(data []byte) error {
	s := strings.TrimSpace(string(data))
	if s == "null" || s == "" {
		*dt = DateTime{}
		return nil
	}
	// Trim surrounding quotes if present
	if len(s) >= 2 && ((s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'')) {
		s = s[1 : len(s)-1]
	}
	parsed, err := Parse(s)
	if err != nil {
		return err
	}
	*dt = parsed
	return nil
}

// Value implements the driver.Valuer interface for database serialization.
func (dt DateTime) Value() (driver.Value, error) {
	return dt.Time, nil
}

// Scan implements the sql.Scanner interface for database deserialization.
func (dt *DateTime) Scan(value any) error {
	switch v := value.(type) {
	case time.Time:
		*dt = DateTime{v}
		return nil
	case string:
		parsed, err := Parse(v)
		if err != nil {
			return err
		}
		*dt = parsed
		return nil
	case []byte:
		parsed, err := Parse(string(v))
		if err != nil {
			return err
		}
		*dt = parsed
		return nil
	case nil:
		*dt = DateTime{}
		return nil
	default:
		return fmt.Errorf("unsupported Scan type %T", value)
	}
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

// IsFirstDayOfMonth returns whether the datetime is the first day of the month.
func (dt DateTime) IsFirstDayOfMonth() bool {
	return dt.Day() == 1
}

// IsLastDayOfMonth returns whether the datetime is the last day of the month.
func (dt DateTime) IsLastDayOfMonth() bool {
	return dt.Day() == dt.DaysInMonth()
}

// IsFirstDayOfYear returns whether the datetime is the first day of the year (January 1st).
func (dt DateTime) IsFirstDayOfYear() bool {
	return dt.Month() == time.January && dt.Day() == 1
}

// IsLastDayOfYear returns whether the datetime is the last day of the year (December 31st).
func (dt DateTime) IsLastDayOfYear() bool {
	return dt.Month() == time.December && dt.Day() == 31
}

// WeekOfMonth returns the week number within the month (1-6).
// The first week of the month is the week containing the first day of the month.
func (dt DateTime) WeekOfMonth() int {
	// Simple calculation: (day - 1) / 7 + 1
	// This ensures that days 1-7 are in week 1, days 8-14 are in week 2, etc.
	return ((dt.Day() - 1) / 7) + 1
}

// WeekOfMonthISO returns the ISO-style week of month using Monday as the first day of week
// and accounting for the weekday of the month's first day.
func (dt DateTime) WeekOfMonthISO() int {
	firstOfMonth := Date(dt.Year(), dt.Month(), 1, 0, 0, 0, 0, dt.Location())
	// Convert Go's weekday (Sun=0..Sat=6) to ISO (Mon=0..Sun=6)
	offset := (int(firstOfMonth.Weekday()) + 6) % 7
	return ((offset + dt.Day() - 1) / 7) + 1
}

// WeekOfMonthWithStart returns the week of the month using a custom week start day.
// For example, start = time.Sunday yields Sunday-start weeks.
func (dt DateTime) WeekOfMonthWithStart(start time.Weekday) int {
	firstOfMonth := Date(dt.Year(), dt.Month(), 1, 0, 0, 0, 0, dt.Location())
	// Compute offset from desired start day
	offset := (int(firstOfMonth.Weekday()) - int(start) + 7) % 7
	return ((offset + dt.Day() - 1) / 7) + 1
}

// DaysInMonth returns the number of days in the datetime's month.
func (dt DateTime) DaysInMonth() int {
	year, month, _ := dt.Date()
	// Create the first day of the next month and subtract one day
	firstOfNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, dt.Location())
	lastOfCurrentMonth := firstOfNextMonth.AddDate(0, 0, -1)
	return lastOfCurrentMonth.Day()
}

// DaysInYear returns the number of days in the datetime's year (365 or 366 for leap years).
func (dt DateTime) DaysInYear() int {
	if dt.IsLeapYear() {
		return 366
	}
	return 365
}

// FromUnixMilli creates a DateTime from a Unix timestamp in milliseconds in the specified location.
func FromUnixMilli(ms int64, loc *time.Location) DateTime {
	return DateTime{time.UnixMilli(ms).In(loc)}
}

// FromUnixMicro creates a DateTime from a Unix timestamp in microseconds in the specified location.
func FromUnixMicro(us int64, loc *time.Location) DateTime {
	return DateTime{time.UnixMicro(us).In(loc)}
}

// FromUnixNano creates a DateTime from a Unix timestamp in nanoseconds in the specified location.
func FromUnixNano(ns int64, loc *time.Location) DateTime {
	return DateTime{time.Unix(0, ns).In(loc)}
}

// DST optimization cache entry
type dstCacheEntry struct {
	standardOffset int
	lastYear       int
}

var dstCache sync.Map // map[*time.Location]*dstCacheEntry

// IsDSTOptimized returns whether the datetime is in daylight saving time using optimized caching
func (dt DateTime) IsDSTOptimized() bool {
	loc := dt.Location()
	year := dt.Year()

	// Fast path for UTC - never DST
	if loc == time.UTC {
		return false
	}

	// Check cache
	if entry, ok := dstCache.Load(loc); ok {
		cacheEntry := entry.(*dstCacheEntry)
		if cacheEntry.lastYear == year {
			_, currentOffset := dt.Zone()
			return currentOffset != cacheEntry.standardOffset
		}
	}

	// Cache miss or stale year - calculate and cache
	standardOffset := getStandardOffsetOptimized(loc, year)
	dstCache.Store(loc, &dstCacheEntry{
		standardOffset: standardOffset,
		lastYear:       year,
	})

	_, currentOffset := dt.Zone()
	return currentOffset != standardOffset
}

// getStandardOffsetOptimized calculates standard offset with optimizations
func getStandardOffsetOptimized(loc *time.Location, year int) int {
	// For most locations, the standard time is the non-DST time, which typically occurs in winter
	// We'll use January (definitely winter) to get the standard offset
	winterTime := time.Date(year, 1, 15, 12, 0, 0, 0, loc)
	_, winterOffset := winterTime.Zone()

	// Also check December to be sure
	decemberTime := time.Date(year, 12, 15, 12, 0, 0, 0, loc)
	_, decemberOffset := decemberTime.Zone()

	// In the northern hemisphere, winter time is standard time
	// In the southern hemisphere, it's more complex, but January/December should be consistent
	if winterOffset == decemberOffset {
		return winterOffset
	}

	// If they're different, take the smaller offset (standard time is usually less than DST)
	if winterOffset < decemberOffset {
		return winterOffset
	}
	return decemberOffset
}

// ClearDSTCache clears the DST cache (useful for testing or memory management)
func ClearDSTCache() {
	dstCache = sync.Map{}
}
