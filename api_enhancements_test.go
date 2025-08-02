package chronogo

import (
	"testing"
	"time"
)

// Test utility methods
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

// Test fluent API
func TestFluentDuration(t *testing.T) {
	// Test with a fixed base date to get predictable results
	base := Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC)

	// Test chaining different units - now uses accurate calendar arithmetic
	result := base.AddFluent().Years(1).Months(2).Days(3).Hours(4).Minutes(5).Seconds(6).To(base)

	// Expected: 2023-01-01 + 1 year = 2024-01-01
	// + 2 months = 2024-03-01 + 3 days = 2024-03-04
	// + 4 hours 5 minutes 6 seconds = 2024-03-04 16:05:06
	expected := Date(2024, time.March, 4, 16, 5, 6, 0, time.UTC)

	if !result.Equal(expected) {
		t.Errorf("Fluent duration calculation: got %v, expected %v", result, expected)
	}
}

// TestFluentDurationAccuracy tests accurate calendar arithmetic
func TestFluentDurationAccuracy(t *testing.T) {
	t.Run("Year arithmetic with leap year", func(t *testing.T) {
		base := Date(2020, time.February, 29, 12, 0, 0, 0, time.UTC) // Leap year
		
		// Add 1 year - Go's time package handles this by moving to March 1st in non-leap years
		result := base.AddFluent().Years(1).To(base)
		expected := Date(2021, time.March, 1, 12, 0, 0, 0, time.UTC) // Go's behavior for Feb 29 + 1 year
		
		if !result.Equal(expected) {
			t.Errorf("Year addition from leap year: got %v, expected %v", result, expected)
		}
	})

	t.Run("Month arithmetic with overflow", func(t *testing.T) {
		base := Date(2023, time.January, 31, 12, 0, 0, 0, time.UTC)
		
		// Add 1 month - Go's time package handles this by adding overflow days to next month
		result := base.AddFluent().Months(1).To(base)
		expected := Date(2023, time.March, 3, 12, 0, 0, 0, time.UTC) // Jan 31 + 1 month = Mar 3 (Feb has 28 days)
		
		if !result.Equal(expected) {
			t.Errorf("Month addition overflow: got %v, expected %v", result, expected)
		}
	})

	t.Run("Combined calendar and time arithmetic", func(t *testing.T) {
		base := Date(2020, time.December, 31, 15, 30, 45, 0, time.UTC)
		
		// Add 1 year, 2 months, 5 days, 3 hours
		result := base.AddFluent().
			Years(1).
			Months(2).
			Days(5).
			Hours(3).
			To(base)
		
		// Expected: 2020-12-31 + 1 year = 2021-12-31
		// + 2 months = 2022-03-03 (Dec 31 + 2 months overflow to March)
		// + 5 days = 2022-03-08
		// + 3 hours = 2022-03-08 18:30:45
		expected := Date(2022, time.March, 8, 18, 30, 45, 0, time.UTC)
		
		if !result.Equal(expected) {
			t.Errorf("Combined arithmetic: got %v, expected %v", result, expected)
		}
	})

	t.Run("Zero values", func(t *testing.T) {
		base := Date(2023, time.June, 15, 12, 30, 45, 0, time.UTC)
		result := base.AddFluent().Years(0).Months(0).Days(0).To(base)
		
		if !result.Equal(base) {
			t.Errorf("Zero additions should return unchanged date: got %v, expected %v", result, base)
		}
	})
}

func TestFluentDateTime(t *testing.T) {
	dt := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)

	result := dt.Set().
		Year(2024).
		Month(time.December).
		Day(25).
		Hour(15).
		Minute(30).
		Second(45).
		Build()

	expected := Date(2024, time.December, 25, 15, 30, 45, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("Fluent datetime setting = %v, want %v", result, expected)
	}
}

func TestFluentDateTimeTimezone(t *testing.T) {
	ny, err := LoadLocation("America/New_York")
	if err != nil {
		t.Skip("Could not load America/New_York timezone")
	}

	dt := Date(2023, time.December, 25, 15, 30, 45, 0, time.UTC)

	result := dt.Set().
		Timezone(ny).
		Hour(12).
		Build()

	if result.Location() != ny {
		t.Errorf("Fluent timezone setting failed: got %v, want %v", result.Location(), ny)
	}

	if result.Hour() != 12 {
		t.Errorf("Fluent hour setting failed: got %d, want 12", result.Hour())
	}
}

// Test ChronoDuration
func TestChronoDuration(t *testing.T) {
	d := NewDuration(25*time.Hour + 30*time.Minute)

	if d.Days() < 1.0 || d.Days() > 1.1 {
		t.Errorf("Days() = %f, want approximately 1.0", d.Days())
	}

	if d.Weeks() < 0.14 || d.Weeks() > 0.16 {
		t.Errorf("Weeks() = %f, want approximately 0.15", d.Weeks())
	}
}

func TestChronoDurationFromComponents(t *testing.T) {
	d := NewDurationFromComponents(2, 30, 45)
	expected := 2*time.Hour + 30*time.Minute + 45*time.Second

	if d.Duration != expected {
		t.Errorf("NewDurationFromComponents() = %v, want %v", d.Duration, expected)
	}
}

func TestChronoDurationOperations(t *testing.T) {
	d1 := NewDuration(2 * time.Hour)
	d2 := NewDuration(30 * time.Minute)

	// Test Add
	sum := d1.Add(d2)
	expected := 2*time.Hour + 30*time.Minute
	if sum.Duration != expected {
		t.Errorf("Add() = %v, want %v", sum.Duration, expected)
	}

	// Test Subtract
	diff := d1.Subtract(d2)
	expected = 1*time.Hour + 30*time.Minute
	if diff.Duration != expected {
		t.Errorf("Subtract() = %v, want %v", diff.Duration, expected)
	}

	// Test Multiply
	product := d1.Multiply(2.5)
	expected = 5 * time.Hour
	if product.Duration != expected {
		t.Errorf("Multiply() = %v, want %v", product.Duration, expected)
	}

	// Test Divide
	quotient := d1.Divide(2)
	expected = 1 * time.Hour
	if quotient.Duration != expected {
		t.Errorf("Divide() = %v, want %v", quotient.Duration, expected)
	}
}

func TestChronoDurationBooleans(t *testing.T) {
	positive := NewDuration(1 * time.Hour)
	negative := NewDuration(-1 * time.Hour)
	zero := NewDuration(0)

	if !positive.IsPositive() || positive.IsNegative() || positive.IsZero() {
		t.Error("Positive duration boolean checks failed")
	}

	if negative.IsPositive() || !negative.IsNegative() || negative.IsZero() {
		t.Error("Negative duration boolean checks failed")
	}

	if zero.IsPositive() || zero.IsNegative() || !zero.IsZero() {
		t.Error("Zero duration boolean checks failed")
	}
}

func TestChronoDurationAbs(t *testing.T) {
	negative := NewDuration(-2 * time.Hour)
	abs := negative.Abs()

	expected := 2 * time.Hour
	if abs.Duration != expected {
		t.Errorf("Abs() = %v, want %v", abs.Duration, expected)
	}

	positive := NewDuration(2 * time.Hour)
	abs = positive.Abs()
	if abs.Duration != expected {
		t.Errorf("Abs() of positive = %v, want %v", abs.Duration, expected)
	}
}
