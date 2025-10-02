package chronogo

import (
	"testing"
	"time"
)

func TestIsBirthday(t *testing.T) {
	birthday := Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		date     DateTime
		expected bool
	}{
		{
			name:     "Same date different year",
			date:     Date(2024, 5, 15, 12, 30, 45, 0, time.UTC),
			expected: true,
		},
		{
			name:     "Different month",
			date:     Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "Different day",
			date:     Date(2024, 5, 16, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "Exact same date",
			date:     Date(1990, 5, 15, 0, 0, 0, 0, time.UTC),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.date.IsBirthday(birthday)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsAnniversary(t *testing.T) {
	anniversary := Date(2010, 6, 20, 14, 30, 0, 0, time.UTC)
	today := Date(2024, 6, 20, 10, 0, 0, 0, time.UTC)

	if !today.IsAnniversary(anniversary) {
		t.Error("Expected IsAnniversary to return true for matching month/day")
	}

	notAnniversary := Date(2024, 6, 21, 10, 0, 0, 0, time.UTC)
	if notAnniversary.IsAnniversary(anniversary) {
		t.Error("Expected IsAnniversary to return false for different day")
	}
}

func TestIsSameDay(t *testing.T) {
	dt1 := Date(2024, 5, 15, 8, 0, 0, 0, time.UTC)
	dt2 := Date(2024, 5, 15, 20, 0, 0, 0, time.UTC)
	dt3 := Date(2024, 5, 16, 8, 0, 0, 0, time.UTC)

	if !dt1.IsSameDay(dt2) {
		t.Error("Expected IsSameDay to return true for same calendar day")
	}

	if dt1.IsSameDay(dt3) {
		t.Error("Expected IsSameDay to return false for different day")
	}
}

func TestIsSameMonth(t *testing.T) {
	dt1 := Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	dt2 := Date(2024, 5, 31, 0, 0, 0, 0, time.UTC)
	dt3 := Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

	if !dt1.IsSameMonth(dt2) {
		t.Error("Expected IsSameMonth to return true for same month")
	}

	if dt1.IsSameMonth(dt3) {
		t.Error("Expected IsSameMonth to return false for different month")
	}
}

func TestIsSameYear(t *testing.T) {
	dt1 := Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	dt2 := Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	dt3 := Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	if !dt1.IsSameYear(dt2) {
		t.Error("Expected IsSameYear to return true for same year")
	}

	if dt1.IsSameYear(dt3) {
		t.Error("Expected IsSameYear to return false for different year")
	}
}

func TestIsSameQuarter(t *testing.T) {
	q1_jan := Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	q1_mar := Date(2024, 3, 20, 0, 0, 0, 0, time.UTC)
	q2_apr := Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)

	if !q1_jan.IsSameQuarter(q1_mar) {
		t.Error("Expected IsSameQuarter to return true for Q1")
	}

	if q1_jan.IsSameQuarter(q2_apr) {
		t.Error("Expected IsSameQuarter to return false for different quarter")
	}
}

func TestIsSameWeek(t *testing.T) {
	// Dates in the same ISO week
	dt1 := Date(2024, 1, 15, 0, 0, 0, 0, time.UTC) // Monday
	dt2 := Date(2024, 1, 17, 0, 0, 0, 0, time.UTC) // Wednesday (same week)
	dt3 := Date(2024, 1, 22, 0, 0, 0, 0, time.UTC) // Next Monday

	if !dt1.IsSameWeek(dt2) {
		t.Error("Expected IsSameWeek to return true for same ISO week")
	}

	if dt1.IsSameWeek(dt3) {
		t.Error("Expected IsSameWeek to return false for different week")
	}
}

func TestAverage(t *testing.T) {
	start := Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	midpoint := start.Average(end)
	expected := Date(2024, 1, 16, 0, 0, 0, 0, time.UTC)

	// Allow for small time difference due to calculation
	diff := midpoint.Time.Sub(expected.Time)
	if diff < 0 {
		diff = -diff
	}

	if diff > time.Hour {
		t.Errorf("Expected midpoint around %v, got %v", expected.Format("2006-01-02"), midpoint.Format("2006-01-02"))
	}
}

func TestClosest(t *testing.T) {
	dt := Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)

	dates := []DateTime{
		Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Date(2024, 6, 10, 0, 0, 0, 0, time.UTC), // 5 days before
		Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
	}

	closest := dt.Closest(dates...)
	expected := Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)

	if closest.Year() != expected.Year() || closest.Month() != expected.Month() || closest.Day() != expected.Day() {
		t.Errorf("Expected closest to be %v, got %v", expected.Format("2006-01-02"), closest.Format("2006-01-02"))
	}
}

func TestFarthest(t *testing.T) {
	dt := Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)

	dates := []DateTime{
		Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Date(2024, 6, 10, 0, 0, 0, 0, time.UTC),
		Date(2024, 12, 31, 0, 0, 0, 0, time.UTC), // Farthest
	}

	farthest := dt.Farthest(dates...)
	expected := Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	if farthest.Year() != expected.Year() || farthest.Month() != expected.Month() || farthest.Day() != expected.Day() {
		t.Errorf("Expected farthest to be %v, got %v", expected.Format("2006-01-02"), farthest.Format("2006-01-02"))
	}
}

func TestClosestEmptyList(t *testing.T) {
	dt := Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	closest := dt.Closest()

	if !closest.IsZero() {
		t.Error("Expected Closest to return zero DateTime for empty list")
	}
}

func TestFarthestEmptyList(t *testing.T) {
	dt := Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	farthest := dt.Farthest()

	if !farthest.IsZero() {
		t.Error("Expected Farthest to return zero DateTime for empty list")
	}
}

func TestToCookieString(t *testing.T) {
	dt := Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	result := dt.ToCookieString()

	// Should be in RFC1123 format
	if result == "" {
		t.Error("Expected non-empty cookie string")
	}

	// Should contain key components
	if len(result) < 20 {
		t.Errorf("Cookie string seems too short: %s", result)
	}
}

func TestToRSSString(t *testing.T) {
	dt := Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	result := dt.ToRSSString()

	// Should be in RFC1123Z format
	if result == "" {
		t.Error("Expected non-empty RSS string")
	}

	// Should contain timezone offset
	if len(result) < 20 {
		t.Errorf("RSS string seems too short: %s", result)
	}
}

func TestToW3CString(t *testing.T) {
	dt := Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	result := dt.ToW3CString()

	// Should be in RFC3339 format
	if result == "" {
		t.Error("Expected non-empty W3C string")
	}

	// Should look like ISO 8601
	if len(result) < 19 {
		t.Errorf("W3C string seems too short: %s", result)
	}
}

func TestToAtomString(t *testing.T) {
	dt := Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	atom := dt.ToAtomString()
	w3c := dt.ToW3CString()

	// Atom should be same as W3C (both RFC3339)
	if atom != w3c {
		t.Errorf("Expected Atom and W3C strings to match: %s vs %s", atom, w3c)
	}
}
