package chronogo

import (
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