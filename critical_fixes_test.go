package chronogo

import (
	"context"
	"testing"
	"time"
)

// TestIsDST_FixedImplementation tests the corrected IsDST method
func TestIsDST_FixedImplementation(t *testing.T) {
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

// TestFluentDuration_AccurateArithmetic tests the fixed FluentDuration implementation
func TestFluentDuration_AccurateArithmetic(t *testing.T) {
	t.Run("Year arithmetic", func(t *testing.T) {
		base := Date(2020, time.February, 29, 12, 0, 0, 0, time.UTC) // Leap year

		// Add 1 year - Go's time package handles this by moving to March 1st in non-leap years
		result := base.AddFluent().Years(1).To(base)
		expected := Date(2021, time.March, 1, 12, 0, 0, 0, time.UTC) // Go's behavior for Feb 29 + 1 year

		if !result.Equal(expected) {
			t.Errorf("Year addition from leap year: got %v, expected %v", result, expected)
		}
	})

	t.Run("Month arithmetic", func(t *testing.T) {
		base := Date(2023, time.January, 31, 12, 0, 0, 0, time.UTC)

		// Add 1 month - Go's time package handles this by adding overflow days to next month
		result := base.AddFluent().Months(1).To(base)
		expected := Date(2023, time.March, 3, 12, 0, 0, 0, time.UTC) // Jan 31 + 1 month = Mar 3 (Feb has 28 days)

		if !result.Equal(expected) {
			t.Errorf("Month addition overflow: got %v, expected %v", result, expected)
		}
	})

	t.Run("Combined year and month arithmetic", func(t *testing.T) {
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

	t.Run("Subtraction arithmetic", func(t *testing.T) {
		base := Date(2023, time.March, 31, 12, 0, 0, 0, time.UTC)

		// Subtract 1 month - Go handles this as March 31 - 1 month = March 3 (overflow)
		result := base.AddFluent().Months(-1).To(base)
		expected := Date(2023, time.March, 3, 12, 0, 0, 0, time.UTC) // Mar 31 - 1 month = Mar 3

		if !result.Equal(expected) {
			t.Errorf("Month subtraction: got %v, expected %v", result, expected)
		}
	})
}

// TestPeriodRange_ContextCancellation tests the context cancellation in Period.Range
func TestPeriodRange_ContextCancellation(t *testing.T) {
	t.Run("Context cancellation stops iteration", func(t *testing.T) {
		start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
		end := Date(2023, time.January, 10, 0, 0, 0, 0, time.UTC)
		period := NewPeriod(start, end)

		ctx, cancel := context.WithCancel(context.Background())

		// Start iteration
		ch := period.RangeWithContext(ctx, "days", 1)

		received := 0
		done := make(chan bool)

		go func() {
			for range ch {
				received++
				if received == 3 {
					cancel() // Cancel after receiving 3 items
				}
			}
			done <- true
		}()

		<-done

		// Should have received exactly 3 items before cancellation
		if received != 3 {
			t.Errorf("Expected to receive 3 items before cancellation, got %d", received)
		}
	})

	t.Run("Context timeout stops iteration", func(t *testing.T) {
		start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
		end := Date(2023, time.January, 10, 0, 0, 0, 0, time.UTC)
		period := NewPeriod(start, end)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		ch := period.RangeWithContext(ctx, "days", 1)

		received := 0
		for range ch {
			received++
			time.Sleep(20 * time.Millisecond) // Slow processing to trigger timeout
		}

		// Should not receive all 10 items due to timeout
		if received >= 10 {
			t.Errorf("Expected timeout to stop iteration early, but received all %d items", received)
		}
	})

	t.Run("Backward compatibility with Range method", func(t *testing.T) {
		start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
		end := Date(2023, time.January, 3, 0, 0, 0, 0, time.UTC)
		period := NewPeriod(start, end)

		// Test that the original Range method still works
		ch := period.Range("days", 1)
		
		received := 0
		for range ch {
			received++
		}

		// Should receive 3 items: Jan 1, 2, 3 (inclusive of start and end)
		expected := 3
		if received != expected {
			t.Errorf("Range method: expected %d items, got %d", expected, received)
		}
	})
}

// TestFluentDuration_EdgeCases tests edge cases for the FluentDuration fixes
func TestFluentDuration_EdgeCases(t *testing.T) {
	t.Run("Leap year edge case", func(t *testing.T) {
		// February 29, 2020 (leap year) + 4 years = February 28, 2024 (next leap year, but should be Feb 29)
		base := Date(2020, time.February, 29, 12, 0, 0, 0, time.UTC)
		result := base.AddFluent().Years(4).To(base)
		expected := Date(2024, time.February, 29, 12, 0, 0, 0, time.UTC) // 2024 is also a leap year

		if !result.Equal(expected) {
			t.Errorf("Leap year addition: got %v, expected %v", result, expected)
		}
	})

	t.Run("Month overflow edge case", func(t *testing.T) {
		// January 31 + 1 month in Go's time package results in March 3
		base := Date(2023, time.January, 31, 12, 0, 0, 0, time.UTC)
		result := base.AddFluent().Months(1).To(base)
		expected := Date(2023, time.March, 3, 12, 0, 0, 0, time.UTC) // Jan 31 + 1 month = Mar 3

		if !result.Equal(expected) {
			t.Errorf("Month overflow: got %v, expected %v", result, expected)
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
