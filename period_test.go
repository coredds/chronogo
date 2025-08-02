package chronogo

import (
	"context"
	"testing"
	"time"
)

func TestNewPeriod(t *testing.T) {
	start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := Date(2023, time.January, 2, 0, 0, 0, 0, time.UTC)

	period := NewPeriod(start, end)

	if !period.Start.Equal(start) {
		t.Errorf("Period start should equal provided start time")
	}

	if !period.End.Equal(end) {
		t.Errorf("Period end should equal provided end time")
	}
}

func TestPeriodDuration(t *testing.T) {
	start := Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC)
	end := Date(2023, time.January, 1, 15, 0, 0, 0, time.UTC)

	period := NewPeriod(start, end)
	duration := period.Duration()

	expected := 3 * time.Hour
	if duration != expected {
		t.Errorf("Expected duration %v, got %v", expected, duration)
	}
}

func TestPeriodContains(t *testing.T) {
	start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := Date(2023, time.January, 3, 0, 0, 0, 0, time.UTC)
	middle := Date(2023, time.January, 2, 0, 0, 0, 0, time.UTC)
	outside := Date(2023, time.January, 4, 0, 0, 0, 0, time.UTC)

	period := NewPeriod(start, end)

	if !period.Contains(start) {
		t.Errorf("Period should contain start time")
	}

	if !period.Contains(end) {
		t.Errorf("Period should contain end time")
	}

	if !period.Contains(middle) {
		t.Errorf("Period should contain time in the middle")
	}

	if period.Contains(outside) {
		t.Errorf("Period should not contain time outside range")
	}
}

func TestPeriodIsNegative(t *testing.T) {
	start := Date(2023, time.January, 2, 0, 0, 0, 0, time.UTC)
	end := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)

	period := NewPeriod(start, end)

	if !period.IsNegative() {
		t.Errorf("Period with end before start should be negative")
	}

	// Test positive period
	positivePeriod := NewPeriod(end, start)
	if positivePeriod.IsNegative() {
		t.Errorf("Period with end after start should not be negative")
	}
}

func TestPeriodAbs(t *testing.T) {
	start := Date(2023, time.January, 2, 0, 0, 0, 0, time.UTC)
	end := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)

	period := NewPeriod(start, end)
	abs := period.Abs()

	if abs.IsNegative() {
		t.Errorf("Abs() should return positive period")
	}

	if !abs.Start.Equal(end) || !abs.End.Equal(start) {
		t.Errorf("Abs() should swap start and end for negative periods")
	}
}

func TestPeriodYears(t *testing.T) {
	tests := []struct {
		start    DateTime
		end      DateTime
		expected int
	}{
		{
			Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			3,
		},
		{
			Date(2020, time.January, 15, 0, 0, 0, 0, time.UTC),
			Date(2023, time.January, 10, 0, 0, 0, 0, time.UTC),
			2, // Not quite 3 full years
		},
		{
			Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			-3, // Negative period
		},
	}

	for _, test := range tests {
		period := NewPeriod(test.start, test.end)
		result := period.Years()

		if result != test.expected {
			t.Errorf("Period from %v to %v: expected %d years, got %d",
				test.start.ToDateString(), test.end.ToDateString(), test.expected, result)
		}
	}
}

func TestPeriodMonths(t *testing.T) {
	tests := []struct {
		start    DateTime
		end      DateTime
		expected int
	}{
		{
			Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			Date(2023, time.March, 1, 0, 0, 0, 0, time.UTC),
			2,
		},
		{
			Date(2023, time.January, 15, 0, 0, 0, 0, time.UTC),
			Date(2023, time.March, 10, 0, 0, 0, 0, time.UTC),
			1, // Not quite 2 full months
		},
	}

	for _, test := range tests {
		period := NewPeriod(test.start, test.end)
		result := period.Months()

		if result != test.expected {
			t.Errorf("Period from %v to %v: expected %d months, got %d",
				test.start.ToDateString(), test.end.ToDateString(), test.expected, result)
		}
	}
}

func TestPeriodDays(t *testing.T) {
	start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := Date(2023, time.January, 4, 0, 0, 0, 0, time.UTC)

	period := NewPeriod(start, end)
	days := period.Days()

	if days != 3 {
		t.Errorf("Expected 3 days, got %d", days)
	}
}

func TestPeriodHours(t *testing.T) {
	start := Date(2023, time.January, 1, 10, 0, 0, 0, time.UTC)
	end := Date(2023, time.January, 1, 15, 0, 0, 0, time.UTC)

	period := NewPeriod(start, end)
	hours := period.Hours()

	if hours != 5 {
		t.Errorf("Expected 5 hours, got %d", hours)
	}
}

func TestPeriodInUnits(t *testing.T) {
	start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := Date(2023, time.January, 2, 12, 0, 0, 0, time.UTC) // 1.5 days

	period := NewPeriod(start, end)

	days := period.InDays()
	if days != 1.5 {
		t.Errorf("Expected 1.5 days, got %f", days)
	}

	hours := period.InHours()
	if hours != 36.0 {
		t.Errorf("Expected 36.0 hours, got %f", hours)
	}
}

func TestPeriodString(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{0, "0 seconds"},
		{1 * time.Second, "1 second"},
		{5 * time.Second, "5 seconds"},
		{1 * time.Minute, "1 minute"},
		{2 * time.Minute, "2 minutes"},
		{1 * time.Hour, "1 hour"},
		{3 * time.Hour, "3 hours"},
		{25 * time.Hour, "1 day and 1 hour"},
		{48 * time.Hour, "2 days"},
	}

	for _, test := range tests {
		start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
		end := start.Add(test.duration)
		period := NewPeriod(start, end)

		result := period.String()
		if result != test.expected {
			t.Errorf("Duration %v: expected '%s', got '%s'", test.duration, test.expected, result)
		}
	}
}

func TestPeriodRange(t *testing.T) {
	start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := Date(2023, time.January, 4, 0, 0, 0, 0, time.UTC)

	period := NewPeriod(start, end)

	// Test ranging by days
	var dates []DateTime
	for dt := range period.Range("days", 1) {
		dates = append(dates, dt)
	}

	expected := 4 // Jan 1, 2, 3, 4
	if len(dates) != expected {
		t.Errorf("Expected %d dates, got %d", expected, len(dates))
	}

	// Test step size
	dates = nil
	for dt := range period.Range("days", 2) {
		dates = append(dates, dt)
	}

	expected = 2 // Jan 1, 3
	if len(dates) != expected {
		t.Errorf("With step 2, expected %d dates, got %d", expected, len(dates))
	}
}

func TestPeriodRangeDays(t *testing.T) {
	start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := Date(2023, time.January, 3, 0, 0, 0, 0, time.UTC)

	period := NewPeriod(start, end)

	var count int
	for range period.RangeDays() {
		count++
	}

	if count != 3 {
		t.Errorf("Expected 3 days in range, got %d", count)
	}
}

func TestPeriodRangeHours(t *testing.T) {
	start := Date(2023, time.January, 1, 10, 0, 0, 0, time.UTC)
	end := Date(2023, time.January, 1, 13, 0, 0, 0, time.UTC)

	period := NewPeriod(start, end)

	var count int
	for range period.RangeHours() {
		count++
	}

	if count != 4 { // 10:00, 11:00, 12:00, 13:00
		t.Errorf("Expected 4 hours in range, got %d", count)
	}
}

func TestPeriodForEach(t *testing.T) {
	start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := Date(2023, time.January, 3, 0, 0, 0, 0, time.UTC)

	period := NewPeriod(start, end)

	var count int
	period.ForEach("days", 1, func(dt DateTime) {
		count++
	})

	if count != 3 {
		t.Errorf("ForEach should iterate 3 times, got %d", count)
	}
}

// TestPeriodRangeWithContext tests the context cancellation in Period.Range
func TestPeriodRangeWithContext(t *testing.T) {
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

		// Test that the original Range method still works and uses RangeWithContext internally
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
