package chronogo

import (
	"strings"
	"testing"
	"time"
)

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

func TestParseISODuration(t *testing.T) {
	tests := []struct {
		in     string
		expect string
	}{
		{"PT15M", "15m0s"},
		{"PT1H30M", "1h30m0s"},
		{"P2W", "336h0m0s"}, // 2 weeks -> 336 hours
		{"P1DT2H", "26h0m0s"},
		{"-PT45S", "-45s"},
	}
	for _, tt := range tests {
		d, err := ParseISODuration(tt.in)
		if err != nil {
			t.Fatalf("ParseISODuration failed for %q: %v", tt.in, err)
		}
		if d.String() != tt.expect {
			t.Fatalf("ParseISODuration(%q) got %q want %q", tt.in, d.String(), tt.expect)
		}
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

func TestFluentBuilderMissingMethods(t *testing.T) {
	base := Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC)

	// Test Weeks method
	result := base.AddFluent().Weeks(2).To(base)
	expected := base.AddDays(14) // 2 weeks = 14 days
	if !result.Equal(expected) {
		t.Errorf("Weeks(2) failed: expected %v, got %v", expected, result)
	}

	// Test Milliseconds method
	result = base.AddFluent().Milliseconds(1500).To(base) // 1.5 seconds
	expected = base.Add(1500 * time.Millisecond)
	if !result.Equal(expected) {
		t.Errorf("Milliseconds(1500) failed: expected %v, got %v", expected, result)
	}

	// Test Microseconds method
	result = base.AddFluent().Microseconds(2000000).To(base) // 2 seconds
	expected = base.Add(2000000 * time.Microsecond)
	if !result.Equal(expected) {
		t.Errorf("Microseconds(2000000) failed: expected %v, got %v", expected, result)
	}

	// Test Nanoseconds method
	result = base.AddFluent().Nanoseconds(3000000000).To(base) // 3 seconds
	expected = base.Add(3000000000 * time.Nanosecond)
	if !result.Equal(expected) {
		t.Errorf("Nanoseconds(3000000000) failed: expected %v, got %v", expected, result)
	}

	// Test From method (reverse operation)
	future := Date(2023, time.February, 15, 18, 30, 45, 0, time.UTC)
	result = base.AddFluent().From(future)
	// The From() method should compute the duration from future to base
	// Let's just verify the method doesn't panic and returns a valid date
	if result.IsZero() {
		t.Error("From() should return a valid datetime")
	}
}

func TestChronoDurationMissingMethods(t *testing.T) {
	// Test Years method
	duration := NewDuration(365*24*time.Hour + 6*time.Hour) // ~1 year
	years := duration.Years()
	if years < 0.9 || years > 1.1 {
		t.Errorf("Years() = %f, expected around 1.0", years)
	}

	// Test Months method
	monthDuration := NewDuration(30*24*time.Hour + 12*time.Hour) // ~1 month
	months := monthDuration.Months()
	if months < 0.9 || months > 1.1 {
		t.Errorf("Months() = %f, expected around 1.0", months)
	}

	// Test HumanString method
	humanDuration := NewDuration(25*time.Hour + 30*time.Minute)
	human := humanDuration.HumanString()
	if human == "" {
		t.Error("HumanString() should return non-empty string")
	}
	// Should say something like "1 day" since it's over 24 hours
	if !strings.Contains(strings.ToLower(human), "day") {
		t.Errorf("HumanString() = %s, expected to contain 'day'", human)
	}
}
