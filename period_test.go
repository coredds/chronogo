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

	// Test Abs() on already positive period
	positivePeriod := NewPeriod(end, start) // end before start makes it positive
	absPositive := positivePeriod.Abs()
	if !absPositive.Start.Equal(positivePeriod.Start) || !absPositive.End.Equal(positivePeriod.End) {
		t.Errorf("Abs() should return same period when already positive")
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

	// Test negative period months
	start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := Date(2023, time.March, 1, 0, 0, 0, 0, time.UTC)
	negativePeriod := NewPeriod(end, start)
	negativeMonths := negativePeriod.Months()
	if negativeMonths != -2 {
		t.Errorf("Expected -2 months for negative period, got %d", negativeMonths)
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

	// Test negative period days
	negativePeriod := NewPeriod(end, start)
	negativeDays := negativePeriod.Days()
	if negativeDays != -3 {
		t.Errorf("Expected -3 days for negative period, got %d", negativeDays)
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

	// Test negative period hours
	negativePeriod := NewPeriod(end, start)
	negativeHours := negativePeriod.Hours()
	if negativeHours != -5 {
		t.Errorf("Expected -5 hours for negative period, got %d", negativeHours)
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

	// Test negative period string
	start := Date(2023, time.January, 2, 0, 0, 0, 0, time.UTC)
	end := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	negativePeriod := NewPeriod(start, end)
	negativeString := negativePeriod.String()
	if negativeString != "-1 day" {
		t.Errorf("Expected '-1 day' for negative period, got '%s'", negativeString)
	}

	// Test complex duration string with multiple parts
	start = Date(2023, time.January, 1, 10, 30, 45, 0, time.UTC)
	end = Date(2023, time.January, 3, 12, 32, 50, 0, time.UTC) // 2 days, 2 hours, 2 minutes, 5 seconds
	complexPeriod := NewPeriod(start, end)
	complexString := complexPeriod.String()
	expectedComplex := "2 days, 2 hours, 2 minutes and 5 seconds"
	if complexString != expectedComplex {
		t.Errorf("Expected '%s' for complex period, got '%s'", expectedComplex, complexString)
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

	// Test invalid unit
	dates = nil
	for dt := range period.Range("invalid", 1) {
		dates = append(dates, dt)
	}
	if len(dates) != 1 {
		t.Errorf("Expected 1 date for invalid unit (start date only), got %d", len(dates))
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
			defer func() { done <- true }() // Ensure done is always sent
			for range ch {
				received++
				if received == 3 {
					cancel() // Cancel after receiving 3 items
					return   // Exit immediately after cancellation
				}
			}
		}()

		// Wait for completion or timeout
		select {
		case <-done:
			// Test completed
		case <-time.After(1 * time.Second):
			cancel() // Ensure cleanup
			t.Fatal("Test timed out waiting for context cancellation")
		}

		// Should have received exactly 3 items before cancellation
		if received < 3 {
			t.Errorf("Expected to receive at least 3 items before cancellation, got %d", received)
		}
		if received > 5 {
			t.Errorf("Expected cancellation to stop iteration quickly, but got %d items", received)
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

// Quick wins tests: typed iteration helper
func TestPeriodRangeByUnit(t *testing.T) {
	loc := time.UTC
	start := Date(2023, time.January, 1, 0, 0, 0, 0, loc)
	end := Date(2023, time.January, 10, 0, 0, 0, 0, loc)
	p := NewPeriod(start, end)

	// Every 3 days
	cnt := 0
	prev := start
	for d := range p.RangeByUnit(UnitDay, 3) {
		if cnt == 0 && !d.Equal(start) {
			t.Fatalf("First day should be start, got %v", d)
		}
		if cnt > 0 {
			if diff := d.Sub(prev); diff != 72*time.Hour { // 3 days
				t.Fatalf("Step mismatch: got %v", diff)
			}
		}
		prev = d
		cnt++
	}
	if cnt == 0 {
		t.Fatalf("Expected at least one iteration")
	}
}

func TestPeriodMinutesAndSeconds(t *testing.T) {
	start := Date(2023, time.January, 1, 10, 30, 45, 0, time.UTC)
	end := Date(2023, time.January, 1, 11, 32, 50, 0, time.UTC) // 1h 2m 5s later

	period := NewPeriod(start, end)

	// Test Minutes method
	minutes := period.Minutes()
	if minutes < 62 || minutes > 63 {
		t.Errorf("Minutes(): expected around %d, got %d", 62, minutes)
	}

	// Test Seconds method
	seconds := period.Seconds()
	expectedSeconds := 3725 // 1h 2m 5s = 3725 seconds
	if seconds != expectedSeconds {
		t.Errorf("Seconds(): expected %d, got %d", expectedSeconds, seconds)
	}

	// Test InMinutes method
	inMinutes := period.InMinutes()
	expectedMinutesFloat := 62.0 + 5.0/60.0 // 1h 2m 5s = 62.083... minutes
	if inMinutes < 62.0 || inMinutes > 63.0 {
		t.Errorf("InMinutes(): expected around %f, got %f", expectedMinutesFloat, inMinutes)
	}

	// Test InSeconds method
	inSeconds := period.InSeconds()
	expectedInSeconds := 3725.0 // 1h 2m 5s = 3725 seconds
	if inSeconds != expectedInSeconds {
		t.Errorf("InSeconds(): expected %f, got %f", expectedInSeconds, inSeconds)
	}

	// Test negative period minutes and seconds
	negativePeriod := NewPeriod(end, start)

	negativeMinutes := negativePeriod.Minutes()
	if negativeMinutes != -minutes {
		t.Errorf("Expected -%d minutes for negative period, got %d", minutes, negativeMinutes)
	}

	negativeSeconds := negativePeriod.Seconds()
	if negativeSeconds != -seconds {
		t.Errorf("Expected -%d seconds for negative period, got %d", seconds, negativeSeconds)
	}
}

// TestPeriodRangeByUnitWithContext tests the typed unit range method with context cancellation
func TestPeriodRangeByUnitWithContext(t *testing.T) {
	t.Run("Context cancellation stops iteration", func(t *testing.T) {
		start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
		end := Date(2023, time.January, 10, 0, 0, 0, 0, time.UTC)
		period := NewPeriod(start, end)

		ctx, cancel := context.WithCancel(context.Background())

		// Start iteration with typed unit
		ch := period.RangeByUnitWithContext(ctx, UnitDay, 1)

		received := 0
		done := make(chan bool)

		go func() {
			defer func() { done <- true }()
			for range ch {
				received++
				if received == 3 {
					cancel() // Cancel after receiving 3 items
					return
				}
			}
		}()

		<-done

		if received != 3 {
			t.Errorf("Expected to receive 3 items before cancellation, got %d", received)
		}
	})

	t.Run("All typed units work correctly", func(t *testing.T) {
		start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
		end := Date(2023, time.January, 1, 2, 0, 0, 0, time.UTC) // 2 hours
		period := NewPeriod(start, end)

		// Test hour unit
		ctx := context.Background()
		ch := period.RangeByUnitWithContext(ctx, UnitHour, 1)
		count := 0
		for range ch {
			count++
		}
		if count != 3 { // 0, 1, 2 hours = 3 items
			t.Errorf("Expected 3 hours, got %d", count)
		}

		// Test minute unit
		ch = period.RangeByUnitWithContext(ctx, UnitMinute, 30)
		count = 0
		for range ch {
			count++
		}
		if count != 5 { // 0, 30, 60, 90, 120 minutes = 5 items
			t.Errorf("Expected 5 30-minute intervals, got %d", count)
		}
	})

	t.Run("Invalid unit stops iteration", func(t *testing.T) {
		start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
		end := Date(2023, time.January, 2, 0, 0, 0, 0, time.UTC)
		period := NewPeriod(start, end)

		ctx := context.Background()
		// Use an invalid unit value
		ch := period.RangeByUnitWithContext(ctx, Unit(999), 1)

		count := 0
		for range ch {
			count++
		}
		if count != 1 {
			t.Errorf("Expected 1 item for invalid unit (start date only), got %d", count)
		}
	})

	t.Run("Week unit works correctly", func(t *testing.T) {
		start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
		end := Date(2023, time.January, 22, 0, 0, 0, 0, time.UTC) // 3 weeks
		period := NewPeriod(start, end)

		ctx := context.Background()
		ch := period.RangeByUnitWithContext(ctx, UnitWeek, 1)
		count := 0
		for range ch {
			count++
		}
		if count != 4 { // 0, 7, 14, 21 days = 4 items
			t.Errorf("Expected 4 weeks, got %d", count)
		}
	})
}

func TestRangeByUnitSlice(t *testing.T) {
	start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := Date(2023, time.January, 5, 0, 0, 0, 0, time.UTC)
	period := NewPeriod(start, end)

	tests := []struct {
		name     string
		unit     Unit
		step     []int
		expected int
	}{
		{"Daily step 1", UnitDay, []int{1}, 5}, // Jan 1, 2, 3, 4, 5
		{"Daily step 2", UnitDay, []int{2}, 3}, // Jan 1, 3, 5
		{"Default step", UnitDay, nil, 5},      // Default step should be 1
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result []DateTime
			if test.step == nil {
				result = period.RangeByUnitSlice(test.unit)
			} else {
				result = period.RangeByUnitSlice(test.unit, test.step[0])
			}

			if len(result) != test.expected {
				t.Errorf("RangeByUnitSlice %s: expected %d items, got %d", test.name, test.expected, len(result))
			}

			// Verify first and last elements
			if len(result) > 0 {
				if !result[0].Equal(period.Start) {
					t.Errorf("First element should equal period start")
				}
				// Check if the last element is within the period
				if result[len(result)-1].After(period.End) {
					t.Errorf("Last element should not exceed period end")
				}
			}

			// Verify slice is properly ordered
			for i := 1; i < len(result); i++ {
				if !result[i].After(result[i-1]) {
					t.Errorf("RangeByUnitSlice should return chronologically ordered slice")
				}
			}
		})
	}
}

// Test with smaller time ranges to avoid memory issues
func TestRangeByUnitSliceHourly(t *testing.T) {
	start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := Date(2023, time.January, 1, 6, 0, 0, 0, time.UTC) // Just 6 hours
	period := NewPeriod(start, end)

	result := period.RangeByUnitSlice(UnitHour, 1)
	expected := 7 // 0, 1, 2, 3, 4, 5, 6 hours

	if len(result) != expected {
		t.Errorf("Hourly RangeByUnitSlice: expected %d items, got %d", expected, len(result))
	}
}

func TestRangeByUnitSliceComparison(t *testing.T) {
	// Compare RangeByUnitSlice with RangeByUnit to ensure same results
	start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := Date(2023, time.January, 8, 0, 0, 0, 0, time.UTC)
	period := NewPeriod(start, end)

	// Get results from both methods
	slice := period.RangeByUnitSlice(UnitDay, 1)

	var iterator []DateTime
	for dt := range period.RangeByUnit(UnitDay, 1) {
		iterator = append(iterator, dt)
	}

	if len(slice) != len(iterator) {
		t.Errorf("RangeByUnitSlice and RangeByUnit should return same number of items: slice=%d, iterator=%d",
			len(slice), len(iterator))
	}

	// Compare each element
	for i := 0; i < len(slice) && i < len(iterator); i++ {
		if !slice[i].Equal(iterator[i]) {
			t.Errorf("Element %d differs: slice=%v, iterator=%v", i, slice[i], iterator[i])
		}
	}
}

func TestFastRangeDays(t *testing.T) {
	start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := Date(2023, time.January, 5, 0, 0, 0, 0, time.UTC)
	period := NewPeriod(start, end)

	tests := []struct {
		name     string
		step     []int
		expected int
	}{
		{"Step 1", []int{1}, 5},  // Jan 1, 2, 3, 4, 5
		{"Step 2", []int{2}, 3},  // Jan 1, 3, 5
		{"Step 3", []int{3}, 2},  // Jan 1, 4
		{"Default step", nil, 5}, // Default step should be 1
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result []DateTime
			if test.step == nil {
				result = period.FastRangeDays()
			} else {
				result = period.FastRangeDays(test.step[0])
			}

			if len(result) != test.expected {
				t.Errorf("FastRangeDays %s: expected %d items, got %d", test.name, test.expected, len(result))
			}

			// Verify first and last elements
			if len(result) > 0 {
				if !result[0].Equal(period.Start) {
					t.Errorf("First element should equal period start")
				}
				// Check that the last element doesn't exceed the period end
				if result[len(result)-1].After(period.End) {
					t.Errorf("Last element should not exceed period end")
				}
			}

			// Verify all elements are at midnight (start of day)
			for i, dt := range result {
				if dt.Hour() != 0 || dt.Minute() != 0 || dt.Second() != 0 || dt.Nanosecond() != 0 {
					t.Errorf("Element %d should be at midnight: %v", i, dt)
				}
			}

			// Verify step size
			if len(result) > 1 && test.step != nil {
				step := test.step[0]
				for i := 1; i < len(result); i++ {
					diff := result[i].Sub(result[i-1]).Hours() / 24 // days
					if int(diff) != step {
						t.Errorf("Step size should be %d days, got %f", step, diff)
					}
				}
			}
		})
	}
}

func TestFastRangeDaysComparison(t *testing.T) {
	// Compare FastRangeDays with RangeDays to ensure same results
	start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := Date(2023, time.January, 10, 0, 0, 0, 0, time.UTC)
	period := NewPeriod(start, end)

	// Get results from both methods
	fast := period.FastRangeDays(1)

	var standard []DateTime
	for dt := range period.RangeDays() {
		standard = append(standard, dt)
	}

	if len(fast) != len(standard) {
		t.Errorf("FastRangeDays and RangeDays should return same number of items: fast=%d, standard=%d",
			len(fast), len(standard))
	}

	// Compare each element
	for i := 0; i < len(fast) && i < len(standard); i++ {
		if !fast[i].Equal(standard[i]) {
			t.Errorf("Element %d differs: fast=%v, standard=%v", i, fast[i], standard[i])
		}
	}
}

func TestRangeByUnitSliceInvalidUnit(t *testing.T) {
	start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := Date(2023, time.January, 5, 0, 0, 0, 0, time.UTC)
	period := NewPeriod(start, end)

	// Test with invalid unit (use a Unit value that doesn't exist)
	invalidUnit := Unit(999) // Not defined in the enum
	result := period.RangeByUnitSlice(invalidUnit)
	if len(result) != 0 {
		t.Errorf("Invalid unit should return empty slice, got %d items", len(result))
	}
}

func TestRangeByUnitSliceZeroStep(t *testing.T) {
	start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := Date(2023, time.January, 5, 0, 0, 0, 0, time.UTC)
	period := NewPeriod(start, end)

	// Test with zero step (should default to 1)
	result := period.RangeByUnitSlice(UnitDay, 0)
	expected := 5 // Should behave like step=1

	if len(result) != expected {
		t.Errorf("Zero step should default to 1: expected %d items, got %d", expected, len(result))
	}
}
