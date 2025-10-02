package chronogo

import (
	"testing"
	"time"
)

func TestNextWeekday(t *testing.T) {
	tests := []struct {
		name     string
		start    DateTime
		weekday  time.Weekday
		expected DateTime
	}{
		{
			name:     "Monday to Wednesday",
			start:    Date(2024, 1, 15, 12, 0, 0, 0, time.UTC), // Monday
			weekday:  time.Wednesday,
			expected: Date(2024, 1, 17, 12, 0, 0, 0, time.UTC), // Wednesday
		},
		{
			name:     "Friday to Monday (next week)",
			start:    Date(2024, 1, 19, 12, 0, 0, 0, time.UTC), // Friday
			weekday:  time.Monday,
			expected: Date(2024, 1, 22, 12, 0, 0, 0, time.UTC), // Next Monday
		},
		{
			name:     "Monday to Monday (next week)",
			start:    Date(2024, 1, 15, 12, 0, 0, 0, time.UTC), // Monday
			weekday:  time.Monday,
			expected: Date(2024, 1, 22, 12, 0, 0, 0, time.UTC), // Next Monday
		},
		{
			name:     "Sunday to Monday",
			start:    Date(2024, 1, 21, 12, 0, 0, 0, time.UTC), // Sunday
			weekday:  time.Monday,
			expected: Date(2024, 1, 22, 12, 0, 0, 0, time.UTC), // Monday
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.start.NextWeekday(tt.weekday)
			if !result.Equal(tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestPreviousWeekday(t *testing.T) {
	tests := []struct {
		name     string
		start    DateTime
		weekday  time.Weekday
		expected DateTime
	}{
		{
			name:     "Friday to Wednesday",
			start:    Date(2024, 1, 19, 12, 0, 0, 0, time.UTC), // Friday
			weekday:  time.Wednesday,
			expected: Date(2024, 1, 17, 12, 0, 0, 0, time.UTC), // Wednesday
		},
		{
			name:     "Monday to Friday (last week)",
			start:    Date(2024, 1, 15, 12, 0, 0, 0, time.UTC), // Monday
			weekday:  time.Friday,
			expected: Date(2024, 1, 12, 12, 0, 0, 0, time.UTC), // Previous Friday
		},
		{
			name:     "Monday to Monday (last week)",
			start:    Date(2024, 1, 15, 12, 0, 0, 0, time.UTC), // Monday
			weekday:  time.Monday,
			expected: Date(2024, 1, 8, 12, 0, 0, 0, time.UTC), // Previous Monday
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.start.PreviousWeekday(tt.weekday)
			if !result.Equal(tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestClosestWeekday(t *testing.T) {
	tests := []struct {
		name     string
		start    DateTime
		weekday  time.Weekday
		expected DateTime
	}{
		{
			name:     "Wednesday - closest is same day",
			start:    Date(2024, 1, 17, 12, 0, 0, 0, time.UTC), // Wednesday
			weekday:  time.Wednesday,
			expected: Date(2024, 1, 17, 12, 0, 0, 0, time.UTC), // Same Wednesday
		},
		{
			name:     "Wednesday to Friday (2 days forward)",
			start:    Date(2024, 1, 17, 12, 0, 0, 0, time.UTC), // Wednesday
			weekday:  time.Friday,
			expected: Date(2024, 1, 19, 12, 0, 0, 0, time.UTC), // Friday (2 days forward)
		},
		{
			name:     "Wednesday to Monday (3 days back vs 5 forward)",
			start:    Date(2024, 1, 17, 12, 0, 0, 0, time.UTC), // Wednesday
			weekday:  time.Monday,
			expected: Date(2024, 1, 15, 12, 0, 0, 0, time.UTC), // Previous Monday (3 days back)
		},
		{
			name:     "Wednesday to Saturday (prefer future if equidistant would be 3/4)",
			start:    Date(2024, 1, 17, 12, 0, 0, 0, time.UTC), // Wednesday
			weekday:  time.Saturday,
			expected: Date(2024, 1, 20, 12, 0, 0, 0, time.UTC), // Next Saturday (3 days)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.start.ClosestWeekday(tt.weekday)
			if !result.Equal(tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestNextOrSameWeekday(t *testing.T) {
	// Monday
	dt := Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	// Same weekday should return same day
	same := dt.NextOrSameWeekday(time.Monday)
	if !same.Equal(dt) {
		t.Errorf("Expected NextOrSameWeekday to return same day for same weekday")
	}

	// Different weekday should return next occurrence
	next := dt.NextOrSameWeekday(time.Wednesday)
	expected := Date(2024, 1, 17, 12, 0, 0, 0, time.UTC)
	if !next.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, next)
	}
}

func TestPreviousOrSameWeekday(t *testing.T) {
	// Monday
	dt := Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	// Same weekday should return same day
	same := dt.PreviousOrSameWeekday(time.Monday)
	if !same.Equal(dt) {
		t.Errorf("Expected PreviousOrSameWeekday to return same day for same weekday")
	}

	// Different weekday should return previous occurrence
	prev := dt.PreviousOrSameWeekday(time.Friday)
	expected := Date(2024, 1, 12, 12, 0, 0, 0, time.UTC)
	if !prev.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, prev)
	}
}

func TestNthWeekdayOfMonth(t *testing.T) {
	tests := []struct {
		name     string
		start    DateTime
		n        int
		weekday  time.Weekday
		expected DateTime
	}{
		{
			name:     "1st Monday of January 2024",
			start:    Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			n:        1,
			weekday:  time.Monday,
			expected: Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "2nd Monday of January 2024",
			start:    Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			n:        2,
			weekday:  time.Monday,
			expected: Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "3rd Friday of March 2024",
			start:    Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			n:        3,
			weekday:  time.Friday,
			expected: Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Last Friday of March 2024",
			start:    Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			n:        -1,
			weekday:  time.Friday,
			expected: Date(2024, 3, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Last Sunday of February 2024",
			start:    Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			n:        -1,
			weekday:  time.Sunday,
			expected: Date(2024, 2, 25, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.start.NthWeekdayOfMonth(tt.n, tt.weekday)
			if result.IsZero() {
				t.Fatalf("NthWeekdayOfMonth returned zero value")
			}
			// Compare dates only (year, month, day)
			if result.Year() != tt.expected.Year() || result.Month() != tt.expected.Month() || result.Day() != tt.expected.Day() {
				t.Errorf("Expected %v, got %v", tt.expected.Format("2006-01-02"), result.Format("2006-01-02"))
			}
		})
	}
}

func TestFirstWeekdayOf(t *testing.T) {
	dt := Date(2024, 3, 15, 0, 0, 0, 0, time.UTC) // March 2024

	firstMonday := dt.FirstWeekdayOf(time.Monday)
	expected := Date(2024, 3, 4, 0, 0, 0, 0, time.UTC)

	if firstMonday.Year() != expected.Year() || firstMonday.Month() != expected.Month() || firstMonday.Day() != expected.Day() {
		t.Errorf("Expected first Monday to be %v, got %v", expected.Format("2006-01-02"), firstMonday.Format("2006-01-02"))
	}
}

func TestLastWeekdayOf(t *testing.T) {
	dt := Date(2024, 3, 15, 0, 0, 0, 0, time.UTC) // March 2024

	lastFriday := dt.LastWeekdayOf(time.Friday)
	expected := Date(2024, 3, 29, 0, 0, 0, 0, time.UTC)

	if lastFriday.Year() != expected.Year() || lastFriday.Month() != expected.Month() || lastFriday.Day() != expected.Day() {
		t.Errorf("Expected last Friday to be %v, got %v", expected.Format("2006-01-02"), lastFriday.Format("2006-01-02"))
	}
}

func TestNthWeekdayOfYear(t *testing.T) {
	dt := Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// 10th Monday of 2024
	tenthMonday := dt.NthWeekdayOfYear(10, time.Monday)

	if tenthMonday.IsZero() {
		t.Fatal("Expected non-zero result for 10th Monday of 2024")
	}

	// Verify it's actually a Monday
	if tenthMonday.Weekday() != time.Monday {
		t.Errorf("Expected result to be a Monday, got %v", tenthMonday.Weekday())
	}

	// Verify it's the 10th one by counting
	count := 0
	current := Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	for current.Year() == 2024 && count < 15 { // Add safety limit
		if current.Weekday() == time.Monday {
			count++
			if count == 10 {
				if tenthMonday.Year() != current.Year() || tenthMonday.Month() != current.Month() || tenthMonday.Day() != current.Day() {
					t.Errorf("Expected 10th Monday to be %v, got %v", current.Format("2006-01-02"), tenthMonday.Format("2006-01-02"))
				}
				return
			}
		}
		current = current.AddDays(1)
	}

	if count < 10 {
		t.Errorf("Only found %d Mondays while counting, expected at least 10", count)
	}
}

func TestIsNthWeekdayOf(t *testing.T) {
	// March 11, 2024 is the 2nd Monday of March
	dt := Date(2024, 3, 11, 0, 0, 0, 0, time.UTC)

	if !dt.IsNthWeekdayOf(2, "month") {
		t.Error("Expected March 11, 2024 to be the 2nd Monday of the month")
	}

	if dt.IsNthWeekdayOf(1, "month") {
		t.Error("Expected March 11, 2024 NOT to be the 1st Monday of the month")
	}

	if dt.IsNthWeekdayOf(3, "month") {
		t.Error("Expected March 11, 2024 NOT to be the 3rd Monday of the month")
	}
}

func TestWeekdayOccurrenceInMonth(t *testing.T) {
	tests := []struct {
		name     string
		dt       DateTime
		expected int
	}{
		{
			name:     "March 4, 2024 - 1st Monday",
			dt:       Date(2024, 3, 4, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		{
			name:     "March 11, 2024 - 2nd Monday",
			dt:       Date(2024, 3, 11, 0, 0, 0, 0, time.UTC),
			expected: 2,
		},
		{
			name:     "March 29, 2024 - 5th Friday",
			dt:       Date(2024, 3, 29, 0, 0, 0, 0, time.UTC),
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.dt.WeekdayOccurrenceInMonth()
			if result != tt.expected {
				t.Errorf("Expected occurrence %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestNthWeekdayOfInvalidInputs(t *testing.T) {
	dt := Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)

	// Invalid n
	result := dt.NthWeekdayOf(0, time.Monday, "month")
	if !result.IsZero() {
		t.Error("Expected zero DateTime for invalid n=0")
	}

	result = dt.NthWeekdayOf(6, time.Monday, "month")
	if !result.IsZero() {
		t.Error("Expected zero DateTime for invalid n=6")
	}

	// Invalid unit
	result = dt.NthWeekdayOf(1, time.Monday, "invalid")
	if !result.IsZero() {
		t.Error("Expected zero DateTime for invalid unit")
	}

	// 6th occurrence doesn't exist
	result = dt.NthWeekdayOf(6, time.Monday, "month")
	if !result.IsZero() {
		t.Error("Expected zero DateTime when occurrence doesn't exist")
	}
}

func TestNthWeekdayOfQuarter(t *testing.T) {
	dt := Date(2024, 2, 1, 0, 0, 0, 0, time.UTC) // Q1 2024

	// First Monday of Q1 2024
	firstMonday := dt.NthWeekdayOf(1, time.Monday, "quarter")
	expected := Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	if !firstMonday.Equal(expected) {
		t.Errorf("Expected first Monday of Q1 to be %v, got %v", expected.Format("2006-01-02"), firstMonday.Format("2006-01-02"))
	}
}

func TestFarthestWeekday(t *testing.T) {
	// Wednesday
	dt := Date(2024, 1, 17, 12, 0, 0, 0, time.UTC)

	// Farthest Saturday from Wednesday: 3 days forward vs 4 days back - should go back
	farthest := dt.FarthestWeekday(time.Saturday)
	expected := Date(2024, 1, 13, 12, 0, 0, 0, time.UTC) // Previous Saturday (4 days back)

	if !farthest.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, farthest)
	}

	// When target is same as current, should return next week
	farthestSame := dt.FarthestWeekday(time.Wednesday)
	expectedSame := Date(2024, 1, 24, 12, 0, 0, 0, time.UTC)

	if !farthestSame.Equal(expectedSame) {
		t.Errorf("Expected %v, got %v", expectedSame, farthestSame)
	}

	// Farthest Sunday from Wednesday: 4 days forward vs 3 days back - should go forward
	farthestSunday := dt.FarthestWeekday(time.Sunday)
	expectedSunday := Date(2024, 1, 21, 12, 0, 0, 0, time.UTC) // Next Sunday (4 days forward)

	if !farthestSunday.Equal(expectedSunday) {
		t.Errorf("Expected farthest Sunday %v, got %v", expectedSunday, farthestSunday)
	}
}
