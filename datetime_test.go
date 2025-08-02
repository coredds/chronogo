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
