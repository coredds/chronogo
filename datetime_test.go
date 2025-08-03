package chronogo

import (
	"testing"
	"time"
)

func TestNow(t *testing.T) {
	dt := Now()
	now := time.Now()

	// Should be within a few milliseconds
	diff := dt.Sub(Instance(now))
	if diff > time.Second || diff < -time.Second {
		t.Errorf("Now() returned time too far from actual now: %v", diff)
	}
}

func TestNowUTC(t *testing.T) {
	dt := NowUTC()
	if dt.Location() != time.UTC {
		t.Errorf("NowUTC() should return UTC time, got %v", dt.Location())
	}
}

func TestNowIn(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip("Could not load America/New_York timezone")
	}

	dt := NowIn(loc)
	if dt.Location() != loc {
		t.Errorf("NowIn() should return time in specified location")
	}
}

func TestToday(t *testing.T) {
	dt := Today()
	now := time.Now()

	if dt.Year() != now.Year() || dt.Month() != now.Month() || dt.Day() != now.Day() {
		t.Errorf("Today() should return today's date")
	}

	if dt.Hour() != 0 || dt.Minute() != 0 || dt.Second() != 0 || dt.Nanosecond() != 0 {
		t.Errorf("Today() should return time at midnight")
	}
}

func TestTomorrowYesterday(t *testing.T) {
	today := Today()
	tomorrow := Tomorrow()
	yesterday := Yesterday()

	if tomorrow.Sub(today) != 24*time.Hour {
		t.Errorf("Tomorrow should be 24 hours after today")
	}

	if today.Sub(yesterday) != 24*time.Hour {
		t.Errorf("Yesterday should be 24 hours before today")
	}
}

func TestDate(t *testing.T) {
	dt := Date(2023, time.December, 25, 15, 30, 45, 123456789, time.UTC)

	if dt.Year() != 2023 {
		t.Errorf("Expected year 2023, got %d", dt.Year())
	}
	if dt.Month() != time.December {
		t.Errorf("Expected December, got %v", dt.Month())
	}
	if dt.Day() != 25 {
		t.Errorf("Expected day 25, got %d", dt.Day())
	}
	if dt.Hour() != 15 {
		t.Errorf("Expected hour 15, got %d", dt.Hour())
	}
	if dt.Minute() != 30 {
		t.Errorf("Expected minute 30, got %d", dt.Minute())
	}
	if dt.Second() != 45 {
		t.Errorf("Expected second 45, got %d", dt.Second())
	}
}

func TestUTC(t *testing.T) {
	dt := UTC(2023, time.December, 25, 15, 30, 45, 0)

	if dt.Location() != time.UTC {
		t.Errorf("UTC() should create time in UTC location")
	}
}

func TestFromUnix(t *testing.T) {
	timestamp := int64(1640995200) // 2022-01-01 00:00:00 UTC
	dt := FromUnix(timestamp, 0, time.UTC)

	if dt.Year() != 2022 || dt.Month() != time.January || dt.Day() != 1 {
		t.Errorf("FromUnix() produced incorrect date: %v", dt)
	}
}

func TestTimezoneConversions(t *testing.T) {
	dt := UTC(2023, time.December, 25, 12, 0, 0, 0)

	// Convert to UTC (should be no change)
	utc := dt.UTC()
	if !utc.Equal(dt) {
		t.Errorf("UTC conversion should not change UTC time")
	}

	// Convert to local
	local := dt.Local()
	if local.Location() != time.Local {
		t.Errorf("Local() should convert to local timezone")
	}
}

func TestIsDST(t *testing.T) {
	// Test cases for different timezones and seasons
	testCases := []struct {
		name     string
		location string
		date     time.Time
		expected bool
	}{
		{
			name:     "New York summer (DST)",
			location: "America/New_York",
			date:     time.Date(2023, time.July, 15, 12, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "New York winter (no DST)",
			location: "America/New_York",
			date:     time.Date(2023, time.January, 15, 12, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "London summer (BST)",
			location: "Europe/London",
			date:     time.Date(2023, time.July, 15, 12, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "London winter (GMT)",
			location: "Europe/London",
			date:     time.Date(2023, time.January, 15, 12, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "UTC (never DST)",
			location: "UTC",
			date:     time.Date(2023, time.July, 15, 12, 0, 0, 0, time.UTC),
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			loc, err := time.LoadLocation(tc.location)
			if err != nil {
				t.Skipf("Skipping test: timezone %s not available", tc.location)
			}

			dt := DateTime{tc.date.In(loc)}
			result := dt.IsDST()

			if result != tc.expected {
				t.Errorf("IsDST() for %s = %v, expected %v", tc.name, result, tc.expected)
			}
		})
	}
}

func TestIsLeapYear(t *testing.T) {
	tests := []struct {
		year     int
		expected bool
	}{
		{2000, true},  // Divisible by 400
		{1900, false}, // Divisible by 100 but not 400
		{2004, true},  // Divisible by 4
		{2001, false}, // Not divisible by 4
	}

	for _, test := range tests {
		dt := Date(test.year, time.January, 1, 0, 0, 0, 0, time.UTC)
		if dt.IsLeapYear() != test.expected {
			t.Errorf("Year %d: expected IsLeapYear()=%v, got %v",
				test.year, test.expected, dt.IsLeapYear())
		}
	}
}

func TestIsPastFuture(t *testing.T) {
	now := Now()
	past := now.AddDays(-1)
	future := now.AddDays(1)

	if !past.IsPast() {
		t.Errorf("Yesterday should be in the past")
	}

	if !future.IsFuture() {
		t.Errorf("Tomorrow should be in the future")
	}

	if past.IsFuture() {
		t.Errorf("Yesterday should not be in the future")
	}

	if future.IsPast() {
		t.Errorf("Tomorrow should not be in the past")
	}
}

func TestArithmetic(t *testing.T) {
	dt := Date(2023, time.January, 15, 12, 30, 45, 0, time.UTC)

	// Test years
	future := dt.AddYears(1)
	if future.Year() != 2024 {
		t.Errorf("AddYears(1) should add one year")
	}

	past := dt.SubtractYears(1)
	if past.Year() != 2022 {
		t.Errorf("SubtractYears(1) should subtract one year")
	}

	// Test months
	future = dt.AddMonths(1)
	if future.Month() != time.February {
		t.Errorf("AddMonths(1) should add one month")
	}

	// Test days
	future = dt.AddDays(1)
	if future.Day() != 16 {
		t.Errorf("AddDays(1) should add one day")
	}

	// Test hours
	future = dt.AddHours(1)
	if future.Hour() != 13 {
		t.Errorf("AddHours(1) should add one hour")
	}

	// Test minutes
	future = dt.AddMinutes(30)
	if future.Minute() != 0 {
		t.Errorf("AddMinutes(30) should add 30 minutes")
	}

	// Test seconds
	future = dt.AddSeconds(15)
	if future.Second() != 0 {
		t.Errorf("AddSeconds(15) should add 15 seconds")
	}
}

func TestSetters(t *testing.T) {
	dt := Date(2023, time.January, 15, 12, 30, 45, 0, time.UTC)

	// Test SetYear
	newDt := dt.SetYear(2024)
	if newDt.Year() != 2024 {
		t.Errorf("SetYear(2024) should set year to 2024")
	}
	if dt.Year() == 2024 {
		t.Errorf("SetYear should return new instance, not modify original")
	}

	// Test SetMonth
	newDt = dt.SetMonth(time.February)
	if newDt.Month() != time.February {
		t.Errorf("SetMonth should set month to February")
	}

	// Test SetDay
	newDt = dt.SetDay(20)
	if newDt.Day() != 20 {
		t.Errorf("SetDay should set day to 20")
	}

	// Test SetHour
	newDt = dt.SetHour(18)
	if newDt.Hour() != 18 {
		t.Errorf("SetHour should set hour to 18")
	}

	// Test SetMinute
	newDt = dt.SetMinute(45)
	if newDt.Minute() != 45 {
		t.Errorf("SetMinute should set minute to 45")
	}

	// Test SetSecond
	newDt = dt.SetSecond(30)
	if newDt.Second() != 30 {
		t.Errorf("SetSecond should set second to 30")
	}
}

func TestComparisons(t *testing.T) {
	dt1 := Date(2023, time.January, 15, 12, 30, 45, 0, time.UTC)
	dt2 := Date(2023, time.January, 15, 12, 30, 45, 0, time.UTC)
	dt3 := Date(2023, time.January, 16, 12, 30, 45, 0, time.UTC)

	// Test Equal
	if !dt1.Equal(dt2) {
		t.Errorf("Equal DateTimes should be equal")
	}

	if dt1.Equal(dt3) {
		t.Errorf("Different DateTimes should not be equal")
	}

	// Test Before/After
	if !dt1.Before(dt3) {
		t.Errorf("dt1 should be before dt3")
	}

	if !dt3.After(dt1) {
		t.Errorf("dt3 should be after dt1")
	}

	if dt1.After(dt3) {
		t.Errorf("dt1 should not be after dt3")
	}

	if dt3.Before(dt1) {
		t.Errorf("dt3 should not be before dt1")
	}
}

func TestSub(t *testing.T) {
	dt1 := Date(2023, time.January, 15, 12, 0, 0, 0, time.UTC)
	dt2 := Date(2023, time.January, 15, 13, 0, 0, 0, time.UTC)

	diff := dt2.Sub(dt1)
	if diff != time.Hour {
		t.Errorf("Expected 1 hour difference, got %v", diff)
	}
}

func TestStringFormats(t *testing.T) {
	dt := Date(2023, time.December, 25, 15, 30, 45, 0, time.UTC)

	// Test ToDateString
	dateStr := dt.ToDateString()
	expected := "2023-12-25"
	if dateStr != expected {
		t.Errorf("ToDateString(): expected %s, got %s", expected, dateStr)
	}

	// Test ToTimeString
	timeStr := dt.ToTimeString()
	expected = "15:30:45"
	if timeStr != expected {
		t.Errorf("ToTimeString(): expected %s, got %s", expected, timeStr)
	}

	// Test ToDateTimeString
	dateTimeStr := dt.ToDateTimeString()
	expected = "2023-12-25 15:30:45"
	if dateTimeStr != expected {
		t.Errorf("ToDateTimeString(): expected %s, got %s", expected, dateTimeStr)
	}

	// Test ToISO8601String
	iso8601Str := dt.ToISO8601String()
	expected = "2023-12-25T15:30:45Z"
	if iso8601Str != expected {
		t.Errorf("ToISO8601String(): expected %s, got %s", expected, iso8601Str)
	}

	// Test String (should be same as ISO8601)
	str := dt.String()
	if str != iso8601Str {
		t.Errorf("String() should equal ToISO8601String()")
	}
}

func TestFormat(t *testing.T) {
	dt := Date(2023, time.December, 25, 15, 30, 45, 0, time.UTC)

	formatted := dt.Format("2006-01-02 15:04:05")
	expected := "2023-12-25 15:30:45"
	if formatted != expected {
		t.Errorf("Format(): expected %s, got %s", expected, formatted)
	}
}

func TestIsFirstDayOfMonth(t *testing.T) {
	tests := []struct {
		name     string
		date     DateTime
		expected bool
	}{
		{"First day of January", Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC), true},
		{"First day of February", Date(2023, time.February, 1, 0, 0, 0, 0, time.UTC), true},
		{"Second day of January", Date(2023, time.January, 2, 12, 0, 0, 0, time.UTC), false},
		{"Last day of January", Date(2023, time.January, 31, 12, 0, 0, 0, time.UTC), false},
		{"Mid-month", Date(2023, time.June, 15, 12, 0, 0, 0, time.UTC), false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.date.IsFirstDayOfMonth()
			if result != test.expected {
				t.Errorf("IsFirstDayOfMonth() for %s: expected %t, got %t", test.date.ToDateString(), test.expected, result)
			}
		})
	}
}

func TestIsLastDayOfMonth(t *testing.T) {
	tests := []struct {
		name     string
		date     DateTime
		expected bool
	}{
		{"Last day of January", Date(2023, time.January, 31, 12, 0, 0, 0, time.UTC), true},
		{"Last day of February (non-leap)", Date(2023, time.February, 28, 12, 0, 0, 0, time.UTC), true},
		{"Last day of February (leap)", Date(2024, time.February, 29, 12, 0, 0, 0, time.UTC), true},
		{"Last day of April", Date(2023, time.April, 30, 12, 0, 0, 0, time.UTC), true},
		{"First day of month", Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC), false},
		{"Mid-month", Date(2023, time.June, 15, 12, 0, 0, 0, time.UTC), false},
		{"February 28 in leap year", Date(2024, time.February, 28, 12, 0, 0, 0, time.UTC), false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.date.IsLastDayOfMonth()
			if result != test.expected {
				t.Errorf("IsLastDayOfMonth() for %s: expected %t, got %t", test.date.ToDateString(), test.expected, result)
			}
		})
	}
}

func TestIsFirstDayOfYear(t *testing.T) {
	tests := []struct {
		name     string
		date     DateTime
		expected bool
	}{
		{"January 1st, 2023", Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC), true},
		{"January 1st, 2024", Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC), true},
		{"January 2nd", Date(2023, time.January, 2, 12, 0, 0, 0, time.UTC), false},
		{"December 31st", Date(2023, time.December, 31, 12, 0, 0, 0, time.UTC), false},
		{"Mid-year", Date(2023, time.June, 15, 12, 0, 0, 0, time.UTC), false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.date.IsFirstDayOfYear()
			if result != test.expected {
				t.Errorf("IsFirstDayOfYear() for %s: expected %t, got %t", test.date.ToDateString(), test.expected, result)
			}
		})
	}
}

func TestIsLastDayOfYear(t *testing.T) {
	tests := []struct {
		name     string
		date     DateTime
		expected bool
	}{
		{"December 31st, 2023", Date(2023, time.December, 31, 12, 0, 0, 0, time.UTC), true},
		{"December 31st, 2024", Date(2024, time.December, 31, 23, 59, 59, 0, time.UTC), true},
		{"December 30th", Date(2023, time.December, 30, 12, 0, 0, 0, time.UTC), false},
		{"January 1st", Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC), false},
		{"Mid-year", Date(2023, time.June, 15, 12, 0, 0, 0, time.UTC), false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.date.IsLastDayOfYear()
			if result != test.expected {
				t.Errorf("IsLastDayOfYear() for %s: expected %t, got %t", test.date.ToDateString(), test.expected, result)
			}
		})
	}
}

func TestWeekOfMonth(t *testing.T) {
	tests := []struct {
		name     string
		date     DateTime
		expected int
	}{
		// Days 1-7 are in week 1, days 8-14 are in week 2, etc.
		{"Day 1", Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC), 1},
		{"Day 2", Date(2023, time.January, 2, 12, 0, 0, 0, time.UTC), 1},
		{"Day 7", Date(2023, time.January, 7, 12, 0, 0, 0, time.UTC), 1},
		{"Day 8", Date(2023, time.January, 8, 12, 0, 0, 0, time.UTC), 2},
		{"Day 14", Date(2023, time.January, 14, 12, 0, 0, 0, time.UTC), 2},
		{"Day 15", Date(2023, time.January, 15, 12, 0, 0, 0, time.UTC), 3},
		{"Day 21", Date(2023, time.January, 21, 12, 0, 0, 0, time.UTC), 3},
		{"Day 22", Date(2023, time.January, 22, 12, 0, 0, 0, time.UTC), 4},
		{"Day 28", Date(2023, time.January, 28, 12, 0, 0, 0, time.UTC), 4},
		{"Day 29", Date(2023, time.January, 29, 12, 0, 0, 0, time.UTC), 5},
		{"Day 31", Date(2023, time.January, 31, 12, 0, 0, 0, time.UTC), 5},

		// Test different months
		{"June 1st", Date(2023, time.June, 1, 12, 0, 0, 0, time.UTC), 1},
		{"June 8th", Date(2023, time.June, 8, 12, 0, 0, 0, time.UTC), 2},
		{"June 15th", Date(2023, time.June, 15, 12, 0, 0, 0, time.UTC), 3},
		{"June 30th", Date(2023, time.June, 30, 12, 0, 0, 0, time.UTC), 5},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.date.WeekOfMonth()
			if result != test.expected {
				t.Errorf("WeekOfMonth() for %s: expected %d, got %d", test.date.ToDateString(), test.expected, result)
			}
		})
	}
}

func TestDaysInMonth(t *testing.T) {
	tests := []struct {
		name     string
		date     DateTime
		expected int
	}{
		{"January 2023", Date(2023, time.January, 15, 12, 0, 0, 0, time.UTC), 31},
		{"February 2023 (non-leap)", Date(2023, time.February, 15, 12, 0, 0, 0, time.UTC), 28},
		{"February 2024 (leap)", Date(2024, time.February, 15, 12, 0, 0, 0, time.UTC), 29},
		{"March 2023", Date(2023, time.March, 15, 12, 0, 0, 0, time.UTC), 31},
		{"April 2023", Date(2023, time.April, 15, 12, 0, 0, 0, time.UTC), 30},
		{"May 2023", Date(2023, time.May, 15, 12, 0, 0, 0, time.UTC), 31},
		{"June 2023", Date(2023, time.June, 15, 12, 0, 0, 0, time.UTC), 30},
		{"July 2023", Date(2023, time.July, 15, 12, 0, 0, 0, time.UTC), 31},
		{"August 2023", Date(2023, time.August, 15, 12, 0, 0, 0, time.UTC), 31},
		{"September 2023", Date(2023, time.September, 15, 12, 0, 0, 0, time.UTC), 30},
		{"October 2023", Date(2023, time.October, 15, 12, 0, 0, 0, time.UTC), 31},
		{"November 2023", Date(2023, time.November, 15, 12, 0, 0, 0, time.UTC), 30},
		{"December 2023", Date(2023, time.December, 15, 12, 0, 0, 0, time.UTC), 31},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.date.DaysInMonth()
			if result != test.expected {
				t.Errorf("DaysInMonth() for %s: expected %d, got %d", test.date.ToDateString(), test.expected, result)
			}
		})
	}
}

func TestDaysInYear(t *testing.T) {
	tests := []struct {
		name     string
		date     DateTime
		expected int
	}{
		{"2023 (non-leap year)", Date(2023, time.June, 15, 12, 0, 0, 0, time.UTC), 365},
		{"2024 (leap year)", Date(2024, time.June, 15, 12, 0, 0, 0, time.UTC), 366},
		{"2000 (leap year)", Date(2000, time.June, 15, 12, 0, 0, 0, time.UTC), 366},
		{"1900 (non-leap year)", Date(1900, time.June, 15, 12, 0, 0, 0, time.UTC), 365},
		{"2100 (non-leap year)", Date(2100, time.June, 15, 12, 0, 0, 0, time.UTC), 365},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.date.DaysInYear()
			if result != test.expected {
				t.Errorf("DaysInYear() for %s: expected %d, got %d", test.date.ToDateString(), test.expected, result)
			}
		})
	}
}

// Start/End operations tests
func TestStartOfDay(t *testing.T) {
	dt := Date(2023, time.December, 25, 15, 30, 45, 123456789, time.UTC)
	startOfDay := dt.StartOfDay()

	expected := Date(2023, time.December, 25, 0, 0, 0, 0, time.UTC)
	if !startOfDay.Equal(expected) {
		t.Errorf("StartOfDay() = %v, want %v", startOfDay, expected)
	}
}

func TestEndOfDay(t *testing.T) {
	dt := Date(2023, time.December, 25, 15, 30, 45, 123456789, time.UTC)
	endOfDay := dt.EndOfDay()

	expected := Date(2023, time.December, 25, 23, 59, 59, 999999999, time.UTC)
	if !endOfDay.Equal(expected) {
		t.Errorf("EndOfDay() = %v, want %v", endOfDay, expected)
	}
}

func TestStartOfMonth(t *testing.T) {
	dt := Date(2023, time.December, 25, 15, 30, 45, 123456789, time.UTC)
	startOfMonth := dt.StartOfMonth()

	expected := Date(2023, time.December, 1, 0, 0, 0, 0, time.UTC)
	if !startOfMonth.Equal(expected) {
		t.Errorf("StartOfMonth() = %v, want %v", startOfMonth, expected)
	}
}

func TestEndOfMonth(t *testing.T) {
	tests := []struct {
		input    DateTime
		expected DateTime
	}{
		{
			Date(2023, time.December, 15, 15, 30, 45, 123456789, time.UTC),
			Date(2023, time.December, 31, 23, 59, 59, 999999999, time.UTC),
		},
		{
			Date(2024, time.February, 15, 12, 0, 0, 0, time.UTC), // Leap year
			Date(2024, time.February, 29, 23, 59, 59, 999999999, time.UTC),
		},
		{
			Date(2023, time.February, 15, 12, 0, 0, 0, time.UTC), // Non-leap year
			Date(2023, time.February, 28, 23, 59, 59, 999999999, time.UTC),
		},
	}

	for _, test := range tests {
		result := test.input.EndOfMonth()
		if !result.Equal(test.expected) {
			t.Errorf("EndOfMonth() for %v = %v, want %v", test.input, result, test.expected)
		}
	}
}

func TestStartOfWeek(t *testing.T) {
	// Test with a Thursday (2023-12-21)
	dt := Date(2023, time.December, 21, 15, 30, 45, 123456789, time.UTC)
	startOfWeek := dt.StartOfWeek()

	// Should be Monday 2023-12-18
	expected := Date(2023, time.December, 18, 0, 0, 0, 0, time.UTC)
	if !startOfWeek.Equal(expected) {
		t.Errorf("StartOfWeek() = %v, want %v", startOfWeek, expected)
	}
}

func TestEndOfWeek(t *testing.T) {
	// Test with a Thursday (2023-12-21)
	dt := Date(2023, time.December, 21, 15, 30, 45, 123456789, time.UTC)
	endOfWeek := dt.EndOfWeek()

	// Should be Sunday 2023-12-24
	expected := Date(2023, time.December, 24, 23, 59, 59, 999999999, time.UTC)
	if !endOfWeek.Equal(expected) {
		t.Errorf("EndOfWeek() = %v, want %v", endOfWeek, expected)
	}
}

func TestStartOfYear(t *testing.T) {
	dt := Date(2023, time.December, 25, 15, 30, 45, 123456789, time.UTC)
	startOfYear := dt.StartOfYear()

	expected := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	if !startOfYear.Equal(expected) {
		t.Errorf("StartOfYear() = %v, want %v", startOfYear, expected)
	}
}

func TestEndOfYear(t *testing.T) {
	dt := Date(2023, time.June, 15, 15, 30, 45, 123456789, time.UTC)
	endOfYear := dt.EndOfYear()

	expected := Date(2023, time.December, 31, 23, 59, 59, 999999999, time.UTC)
	if !endOfYear.Equal(expected) {
		t.Errorf("EndOfYear() = %v, want %v", endOfYear, expected)
	}
}

func TestIsWeekend(t *testing.T) {
	tests := []struct {
		dt       DateTime
		expected bool
	}{
		{Date(2023, time.December, 23, 12, 0, 0, 0, time.UTC), true},  // Saturday
		{Date(2023, time.December, 24, 12, 0, 0, 0, time.UTC), true},  // Sunday
		{Date(2023, time.December, 25, 12, 0, 0, 0, time.UTC), false}, // Monday
		{Date(2023, time.December, 22, 12, 0, 0, 0, time.UTC), false}, // Friday
	}

	for _, test := range tests {
		result := test.dt.IsWeekend()
		if result != test.expected {
			t.Errorf("IsWeekend() for %v = %v, want %v", test.dt.Weekday(), result, test.expected)
		}
	}
}

func TestIsWeekday(t *testing.T) {
	tests := []struct {
		dt       DateTime
		expected bool
	}{
		{Date(2023, time.December, 23, 12, 0, 0, 0, time.UTC), false}, // Saturday
		{Date(2023, time.December, 24, 12, 0, 0, 0, time.UTC), false}, // Sunday
		{Date(2023, time.December, 25, 12, 0, 0, 0, time.UTC), true},  // Monday
		{Date(2023, time.December, 22, 12, 0, 0, 0, time.UTC), true},  // Friday
	}

	for _, test := range tests {
		result := test.dt.IsWeekday()
		if result != test.expected {
			t.Errorf("IsWeekday() for %v = %v, want %v", test.dt.Weekday(), result, test.expected)
		}
	}
}

func TestQuarter(t *testing.T) {
	tests := []struct {
		dt       DateTime
		expected int
	}{
		{Date(2023, time.January, 15, 12, 0, 0, 0, time.UTC), 1},
		{Date(2023, time.March, 31, 12, 0, 0, 0, time.UTC), 1},
		{Date(2023, time.April, 1, 12, 0, 0, 0, time.UTC), 2},
		{Date(2023, time.June, 30, 12, 0, 0, 0, time.UTC), 2},
		{Date(2023, time.July, 1, 12, 0, 0, 0, time.UTC), 3},
		{Date(2023, time.September, 30, 12, 0, 0, 0, time.UTC), 3},
		{Date(2023, time.October, 1, 12, 0, 0, 0, time.UTC), 4},
		{Date(2023, time.December, 31, 12, 0, 0, 0, time.UTC), 4},
	}

	for _, test := range tests {
		result := test.dt.Quarter()
		if result != test.expected {
			t.Errorf("Quarter() for %v = %v, want %v", test.dt.Month(), result, test.expected)
		}
	}
}

func TestStartOfQuarter(t *testing.T) {
	tests := []struct {
		input    DateTime
		expected DateTime
	}{
		{
			Date(2023, time.February, 15, 12, 0, 0, 0, time.UTC),
			Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Date(2023, time.May, 15, 12, 0, 0, 0, time.UTC),
			Date(2023, time.April, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Date(2023, time.August, 15, 12, 0, 0, 0, time.UTC),
			Date(2023, time.July, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Date(2023, time.November, 15, 12, 0, 0, 0, time.UTC),
			Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, test := range tests {
		result := test.input.StartOfQuarter()
		if !result.Equal(test.expected) {
			t.Errorf("StartOfQuarter() for %v = %v, want %v", test.input, result, test.expected)
		}
	}
}

func TestEndOfQuarter(t *testing.T) {
	tests := []struct {
		input    DateTime
		expected DateTime
	}{
		{
			Date(2023, time.February, 15, 12, 0, 0, 0, time.UTC),
			Date(2023, time.March, 31, 23, 59, 59, 999999999, time.UTC),
		},
		{
			Date(2023, time.May, 15, 12, 0, 0, 0, time.UTC),
			Date(2023, time.June, 30, 23, 59, 59, 999999999, time.UTC),
		},
		{
			Date(2023, time.August, 15, 12, 0, 0, 0, time.UTC),
			Date(2023, time.September, 30, 23, 59, 59, 999999999, time.UTC),
		},
		{
			Date(2023, time.November, 15, 12, 0, 0, 0, time.UTC),
			Date(2023, time.December, 31, 23, 59, 59, 999999999, time.UTC),
		},
	}

	for _, test := range tests {
		result := test.input.EndOfQuarter()
		if !result.Equal(test.expected) {
			t.Errorf("EndOfQuarter() for %v = %v, want %v", test.input, result, test.expected)
		}
	}
}

func TestISOWeek(t *testing.T) {
	dt := Date(2023, time.December, 25, 12, 0, 0, 0, time.UTC)
	year, week := dt.ISOWeek()

	// December 25, 2023 should be in week 52 of 2023
	if year != 2023 || week != 52 {
		t.Errorf("ISOWeek() = (%d, %d), want (2023, 52)", year, week)
	}
}

func TestISOWeekYear(t *testing.T) {
	dt := Date(2023, time.December, 25, 12, 0, 0, 0, time.UTC)
	year := dt.ISOWeekYear()

	if year != 2023 {
		t.Errorf("ISOWeekYear() = %d, want 2023", year)
	}
}

func TestISOWeekNumber(t *testing.T) {
	dt := Date(2023, time.December, 25, 12, 0, 0, 0, time.UTC)
	week := dt.ISOWeekNumber()

	if week != 52 {
		t.Errorf("ISOWeekNumber() = %d, want 52", week)
	}
}

func TestDayOfYear(t *testing.T) {
	tests := []struct {
		dt       DateTime
		expected int
	}{
		{Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC), 1},
		{Date(2023, time.December, 31, 12, 0, 0, 0, time.UTC), 365},
		{Date(2024, time.December, 31, 12, 0, 0, 0, time.UTC), 366}, // Leap year
	}

	for _, test := range tests {
		result := test.dt.DayOfYear()
		if result != test.expected {
			t.Errorf("DayOfYear() for %v = %d, want %d", test.dt, result, test.expected)
		}
	}
}
