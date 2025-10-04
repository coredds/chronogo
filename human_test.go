package chronogo

import (
	"strings"
	"testing"
	"time"
)

func TestDiffForHumans(t *testing.T) {
	now := Date(2023, time.January, 15, 12, 0, 0, 0, time.UTC)

	// Test relative mode (compared to now)
	relativeTests := []struct {
		dt       DateTime
		expected string
	}{
		{now.AddSeconds(-30), "30 seconds ago"},
		{now.AddMinutes(-1), "1 minute ago"},
		{now.AddMinutes(-5), "5 minutes ago"},
		{now.AddHours(-1), "1 hour ago"},
		{now.AddHours(-3), "3 hours ago"},
		{now.AddDays(-1), "1 day ago"},
		{now.AddDays(-7), "1 week ago"},
		{now.AddDays(-30), "1 month ago"},
		{now.AddYears(-1), "1 year ago"},
		{now.AddSeconds(30), "in 30 seconds"},
		{now.AddMinutes(1), "in 1 minute"},
		{now.AddMinutes(5), "in 5 minutes"},
		{now.AddHours(1), "in 1 hour"},
		{now.AddHours(3), "in 3 hours"},
		{now.AddDays(1), "in 1 day"},
		{now.AddDays(7), "in 1 week"},
		{now.AddDays(30), "in 1 month"},
		{now.AddYears(1), "in 1 year"},
	}

	for _, test := range relativeTests {
		result := test.dt.DiffForHumans(now)
		if result != test.expected {
			t.Errorf("DiffForHumans: expected '%s', got '%s'", test.expected, result)
		}
	}
}

func TestDiffForHumansComparison(t *testing.T) {
	// Set to English for consistent testing
	_ = SetDefaultLocale("en-US")
	defer func() { _ = SetDefaultLocale("en-US") }()

	dt1 := Date(2023, time.January, 15, 12, 0, 0, 0, time.UTC)
	dt2 := Date(2023, time.January, 15, 13, 0, 0, 0, time.UTC)
	dt3 := Date(2023, time.January, 15, 11, 0, 0, 0, time.UTC)

	// Test comparison mode - should return human-readable differences
	// Note: In the refactored version, DiffForHumansComparison uses the same
	// locale-aware patterns as DiffForHumans since many languages don't
	// distinguish between "ago/in" and "before/after"
	result1 := dt2.DiffForHumansComparison(dt1)
	if !strings.Contains(result1, "hour") {
		t.Errorf("Expected result to contain 'hour', got '%s'", result1)
	}

	result2 := dt3.DiffForHumansComparison(dt1)
	if !strings.Contains(result2, "hour") {
		t.Errorf("Expected result to contain 'hour', got '%s'", result2)
	}
}

func TestDiffForHumansNow(t *testing.T) {
	// Test that DiffForHumansNow returns a reasonable result
	dt := Now().AddHours(-1)
	result := dt.DiffForHumansNow()

	// Should contain "ago" since it's in the past
	if !strings.Contains(result, "ago") {
		t.Errorf("DiffForHumansNow should indicate past time, got '%s'", result)
	}
}

func TestHumanize(t *testing.T) {
	// Set to English for consistent testing
	_ = SetDefaultLocale("en-US")
	defer func() { _ = SetDefaultLocale("en-US") }()

	tests := []struct {
		duration time.Duration
		expected string
	}{
		{0, "0 seconds"},
		{30 * time.Second, "30 seconds"},
		{1 * time.Minute, "1 minute"},
		{5 * time.Minute, "5 minutes"},
		{1 * time.Hour, "1 hour"},
		{3 * time.Hour, "3 hours"},
		{25 * time.Hour, "1 day"},
		{8 * 24 * time.Hour, "1 week"},
		{-1 * time.Hour, "-1 hour"},
	}

	for _, test := range tests {
		result := Humanize(test.duration)
		if result != test.expected {
			t.Errorf("Humanize(%v): expected '%s', got '%s'", test.duration, test.expected, result)
		}
	}
}

func TestAge(t *testing.T) {
	// Test that Age() returns a non-empty string for past dates
	pastDate := Now().AddYears(-1)
	result := pastDate.Age()
	if len(result) == 0 {
		t.Errorf("Age() should return non-empty string")
	}

	// Test that future dates return "not yet born"
	futureDate := Now().AddDays(1)
	result = futureDate.Age()
	if result != "not yet born" {
		t.Errorf("Future dates should return 'not yet born', got '%s'", result)
	}

	// Test that very old dates return years
	veryOldDate := Now().AddYears(-10)
	result = veryOldDate.Age()
	if !strings.Contains(result, "year") {
		t.Errorf("Very old dates should mention years, got '%s'", result)
	}
}

func TestTimeFromNow(t *testing.T) {
	dt := Now().AddHours(1)
	result := dt.TimeFromNow()

	// Should return a non-empty string
	if len(result) == 0 {
		t.Errorf("TimeFromNow should return non-empty string")
	}
}

func TestTimeAgo(t *testing.T) {
	dt := Now().AddHours(-1)
	result := dt.TimeAgo()

	// Should return a non-empty string
	if len(result) == 0 {
		t.Errorf("TimeAgo should return non-empty string")
	}
}

// Test edge cases and boundary conditions
func TestDiffForHumansEdgeCases(t *testing.T) {
	dt1 := Date(2023, time.January, 15, 12, 0, 0, 0, time.UTC)
	dt2 := Date(2023, time.January, 15, 12, 0, 0, 0, time.UTC)

	// Same time should show minimal difference
	result := dt1.DiffForHumans(dt2)
	if !strings.Contains(result, "second") {
		t.Errorf("Same time comparison should show seconds, got '%s'", result)
	}
}

func TestPluralHandling(t *testing.T) {
	now := Date(2023, time.January, 15, 12, 0, 0, 0, time.UTC)

	// Test singular forms
	single := now.AddMinutes(-1)
	result := single.DiffForHumans(now)
	if !strings.Contains(result, "1 minute") {
		t.Errorf("Single minute should use singular form, got '%s'", result)
	}

	// Test plural forms
	multiple := now.AddMinutes(-5)
	result = multiple.DiffForHumans(now)
	if !strings.Contains(result, "5 minutes") {
		t.Errorf("Multiple minutes should use plural form, got '%s'", result)
	}
}
