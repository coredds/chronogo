package chronogo

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestDiff(t *testing.T) {
	dt1 := Date(2023, time.January, 15, 10, 0, 0, 0, time.UTC)
	dt2 := Date(2024, time.March, 20, 14, 30, 0, 0, time.UTC)

	diff := dt2.Diff(dt1)

	// Test basic properties
	if diff.IsNegative() {
		t.Error("Diff should be positive")
	}

	if !diff.IsPositive() {
		t.Error("Diff should be positive")
	}

	if diff.IsZero() {
		t.Error("Diff should not be zero")
	}

	// Test Start/End
	if !diff.Start().Equal(dt1) {
		t.Errorf("Start should be %v, got %v", dt1, diff.Start())
	}

	if !diff.End().Equal(dt2) {
		t.Errorf("End should be %v, got %v", dt2, diff.End())
	}
}

func TestDiffAbs(t *testing.T) {
	dt1 := Date(2024, time.March, 20, 14, 30, 0, 0, time.UTC)
	dt2 := Date(2023, time.January, 15, 10, 0, 0, 0, time.UTC)

	// dt1 > dt2, so diff will be negative
	diff := dt2.Diff(dt1)

	if !diff.IsNegative() {
		t.Error("Diff should be negative")
	}

	// Test DiffAbs always returns positive
	diffAbs := dt2.DiffAbs(dt1)
	if !diffAbs.IsPositive() {
		t.Error("DiffAbs should be positive")
	}

	// Test Abs() method
	absManual := diff.Abs()
	if !absManual.IsPositive() {
		t.Error("Abs() should return positive diff")
	}
}

func TestDiffCalendarMethods(t *testing.T) {
	tests := []struct {
		name     string
		dt1      DateTime
		dt2      DateTime
		wantYear int
		wantMon  int
		wantDays int
	}{
		{
			name:     "One year exact",
			dt1:      Date(2023, time.January, 15, 10, 0, 0, 0, time.UTC),
			dt2:      Date(2024, time.January, 15, 10, 0, 0, 0, time.UTC),
			wantYear: 1,
			wantMon:  12,
			wantDays: 365, // Exact 1 year duration
		},
		{
			name:     "Two years and 3 months",
			dt1:      Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
			dt2:      Date(2023, time.April, 1, 0, 0, 0, 0, time.UTC),
			wantYear: 2,
			wantMon:  27, // 2*12 + 3
			wantDays: 820,
		},
		{
			name:     "Less than a year",
			dt1:      Date(2023, time.January, 15, 0, 0, 0, 0, time.UTC),
			dt2:      Date(2023, time.June, 20, 0, 0, 0, 0, time.UTC),
			wantYear: 0,
			wantMon:  5,
			wantDays: 156,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff := tt.dt2.Diff(tt.dt1)

			if got := diff.Years(); got != tt.wantYear {
				t.Errorf("Years() = %d, want %d", got, tt.wantYear)
			}

			if got := diff.Months(); got != tt.wantMon {
				t.Errorf("Months() = %d, want %d", got, tt.wantMon)
			}

			if got := diff.Days(); got != tt.wantDays {
				t.Errorf("Days() = %d, want %d", got, tt.wantDays)
			}
		})
	}
}

func TestDiffPreciseMethods(t *testing.T) {
	dt1 := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	dt2 := Date(2023, time.January, 3, 5, 30, 45, 0, time.UTC)

	diff := dt2.Diff(dt1)

	// Test Hours
	wantHours := 2*24 + 5 // 2 days + 5 hours
	if got := diff.Hours(); got != wantHours {
		t.Errorf("Hours() = %d, want %d", got, wantHours)
	}

	// Test Minutes
	wantMinutes := wantHours*60 + 30
	if got := diff.Minutes(); got != wantMinutes {
		t.Errorf("Minutes() = %d, want %d", got, wantMinutes)
	}

	// Test Seconds
	wantSeconds := wantMinutes*60 + 45
	if got := diff.Seconds(); got != wantSeconds {
		t.Errorf("Seconds() = %d, want %d", got, wantSeconds)
	}

	// Test InHours (float)
	wantInHours := float64(wantHours) + 30.0/60.0 + 45.0/3600.0
	if got := diff.InHours(); got < wantInHours-0.01 || got > wantInHours+0.01 {
		t.Errorf("InHours() = %f, want ~%f", got, wantInHours)
	}
}

func TestDiffWeeks(t *testing.T) {
	dt1 := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	dt2 := Date(2023, time.January, 22, 0, 0, 0, 0, time.UTC) // 21 days = 3 weeks

	diff := dt2.Diff(dt1)

	if got := diff.Weeks(); got != 3 {
		t.Errorf("Weeks() = %d, want 3", got)
	}

	if got := diff.InWeeks(); got < 2.99 || got > 3.01 {
		t.Errorf("InWeeks() = %f, want ~3.0", got)
	}
}

func TestDiffInvert(t *testing.T) {
	dt1 := Date(2023, time.January, 15, 10, 0, 0, 0, time.UTC)
	dt2 := Date(2024, time.March, 20, 14, 30, 0, 0, time.UTC)

	diff := dt2.Diff(dt1)
	inverted := diff.Invert()

	// Check sign flip
	if diff.IsNegative() == inverted.IsNegative() {
		t.Error("Invert() should flip the sign")
	}

	// Check start/end swap
	if !inverted.Start().Equal(dt2) {
		t.Error("Inverted start should be original end")
	}

	if !inverted.End().Equal(dt1) {
		t.Error("Inverted end should be original start")
	}

	// Check duration magnitude is same
	if diff.Abs().InSeconds() != inverted.Abs().InSeconds() {
		t.Error("Inverted diff should have same magnitude")
	}
}

func TestDiffString(t *testing.T) {
	tests := []struct {
		name     string
		dt1      DateTime
		dt2      DateTime
		contains []string
	}{
		{
			name:     "Zero diff",
			dt1:      Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			dt2:      Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			contains: []string{"0 seconds"},
		},
		{
			name:     "Years and months",
			dt1:      Date(2020, time.January, 15, 0, 0, 0, 0, time.UTC),
			dt2:      Date(2022, time.March, 20, 0, 0, 0, 0, time.UTC),
			contains: []string{"year", "month"},
		},
		{
			name:     "Days and hours",
			dt1:      Date(2023, time.January, 1, 10, 0, 0, 0, time.UTC),
			dt2:      Date(2023, time.January, 3, 14, 30, 0, 0, time.UTC),
			contains: []string{"day", "hour"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff := tt.dt2.Diff(tt.dt1)
			str := diff.String()

			for _, want := range tt.contains {
				if !strings.Contains(str, want) {
					t.Errorf("String() = %q, should contain %q", str, want)
				}
			}
		})
	}
}

func TestDiffCompactString(t *testing.T) {
	tests := []struct {
		name     string
		dt1      DateTime
		dt2      DateTime
		contains []string
	}{
		{
			name:     "Zero diff",
			dt1:      Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			dt2:      Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			contains: []string{"0s"},
		},
		{
			name:     "Years and months",
			dt1:      Date(2020, time.January, 15, 0, 0, 0, 0, time.UTC),
			dt2:      Date(2022, time.March, 20, 0, 0, 0, 0, time.UTC),
			contains: []string{"y", "mo"},
		},
		{
			name:     "Hours and minutes",
			dt1:      Date(2023, time.January, 1, 10, 0, 0, 0, time.UTC),
			dt2:      Date(2023, time.January, 1, 14, 30, 0, 0, time.UTC),
			contains: []string{"h", "m"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff := tt.dt2.Diff(tt.dt1)
			str := diff.CompactString()

			for _, want := range tt.contains {
				if !strings.Contains(str, want) {
					t.Errorf("CompactString() = %q, should contain %q", str, want)
				}
			}
		})
	}
}

func TestDiffTypeForHumans(t *testing.T) {
	now := Date(2023, time.June, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		dt       DateTime
		contains string
	}{
		{
			name:     "Future",
			dt:       now.AddDays(5),
			contains: "in",
		},
		{
			name:     "Past",
			dt:       now.AddDays(-5),
			contains: "ago",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff := tt.dt.Diff(now)
			str := diff.ForHumans()

			if !strings.Contains(str, tt.contains) {
				t.Errorf("ForHumans() = %q, should contain %q", str, tt.contains)
			}
		})
	}
}

func TestDiffTypeForHumansComparison(t *testing.T) {
	// Set to English for consistent testing
	SetDefaultLocale("en-US")
	defer SetDefaultLocale("en-US")

	dt1 := Date(2023, time.June, 15, 12, 0, 0, 0, time.UTC)
	dt2 := Date(2023, time.June, 20, 12, 0, 0, 0, time.UTC)

	diff := dt2.Diff(dt1)
	str := diff.ForHumansComparison()

	// Should contain "days" since there's a 5-day difference
	// Note: In the refactored version, ForHumansComparison uses locale-aware
	// patterns which may not distinguish "ago/in" from "before/after"
	if !strings.Contains(str, "day") {
		t.Errorf("ForHumansComparison() = %q, should contain 'day'", str)
	}

	// Inverted should also contain "days"
	inverted := dt1.Diff(dt2)
	strInv := inverted.ForHumansComparison()

	if !strings.Contains(strInv, "day") {
		t.Errorf("ForHumansComparison() = %q, should contain 'day'", strInv)
	}
}

func TestDiffCompare(t *testing.T) {
	dt1 := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	dt2 := Date(2023, time.January, 10, 0, 0, 0, 0, time.UTC)
	dt3 := Date(2023, time.January, 20, 0, 0, 0, 0, time.UTC)

	diff1 := dt2.Diff(dt1) // 9 days
	diff2 := dt3.Diff(dt1) // 19 days

	// diff1 should be shorter than diff2
	if diff1.Compare(diff2) >= 0 {
		t.Error("diff1 should be shorter than diff2")
	}

	if !diff1.ShorterThan(diff2) {
		t.Error("ShorterThan() should return true")
	}

	if !diff2.LongerThan(diff1) {
		t.Error("LongerThan() should return true")
	}

	// Same diff should be equal
	if !diff1.EqualTo(diff1) {
		t.Error("EqualTo() should return true for same diff")
	}
}

func TestDiffInYearsAndMonths(t *testing.T) {
	dt1 := Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	dt2 := Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)

	diff := dt2.Diff(dt1)

	// Should be approximately 1 year
	inYears := diff.InYears()
	if inYears < 0.99 || inYears > 1.01 {
		t.Errorf("InYears() = %f, want ~1.0", inYears)
	}

	// Should be approximately 12 months
	inMonths := diff.InMonths()
	if inMonths < 11.9 || inMonths > 12.1 {
		t.Errorf("InMonths() = %f, want ~12.0", inMonths)
	}
}

func TestDiffNegative(t *testing.T) {
	dt1 := Date(2023, time.June, 20, 0, 0, 0, 0, time.UTC)
	dt2 := Date(2023, time.June, 10, 0, 0, 0, 0, time.UTC)

	diff := dt2.Diff(dt1) // dt2 < dt1, so negative

	if !diff.IsNegative() {
		t.Error("Diff should be negative")
	}

	if diff.IsPositive() {
		t.Error("Diff should not be positive")
	}

	// String should include negative sign
	str := diff.String()
	if !strings.HasPrefix(str, "-") {
		t.Errorf("String() = %q, should start with '-'", str)
	}

	// CompactString should include negative sign
	compact := diff.CompactString()
	if !strings.HasPrefix(compact, "-") {
		t.Errorf("CompactString() = %q, should start with '-'", compact)
	}
}

func TestDiffMarshalJSON(t *testing.T) {
	dt1 := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	dt2 := Date(2023, time.January, 10, 12, 30, 0, 0, time.UTC)

	diff := dt2.Diff(dt1)

	data, err := json.Marshal(diff)
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}

	// Should be valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Unmarshal error = %v", err)
	}

	// Should have required fields
	if _, ok := result["duration"]; !ok {
		t.Error("JSON should contain 'duration' field")
	}

	if _, ok := result["start"]; !ok {
		t.Error("JSON should contain 'start' field")
	}

	if _, ok := result["end"]; !ok {
		t.Error("JSON should contain 'end' field")
	}
}

func TestDiffDurationAndPeriod(t *testing.T) {
	dt1 := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	dt2 := Date(2023, time.January, 10, 12, 30, 0, 0, time.UTC)

	diff := dt2.Diff(dt1)

	// Test Duration() method
	duration := diff.Duration()
	expected := dt2.Sub(dt1)
	if duration != expected {
		t.Errorf("Duration() = %v, want %v", duration, expected)
	}

	// Test Period() method
	period := diff.Period()
	if !period.Start.Equal(dt1) || !period.End.Equal(dt2) {
		t.Error("Period() should return correct start and end")
	}
}

// BenchmarkDiff benchmarks the Diff creation
func BenchmarkDiff(b *testing.B) {
	dt1 := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	dt2 := Date(2024, time.March, 20, 14, 30, 0, 0, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dt2.Diff(dt1)
	}
}

// BenchmarkDiffMethods benchmarks various Diff methods
func BenchmarkDiffMethods(b *testing.B) {
	dt1 := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	dt2 := Date(2024, time.March, 20, 14, 30, 0, 0, time.UTC)
	diff := dt2.Diff(dt1)

	b.Run("Years", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = diff.Years()
		}
	})

	b.Run("Months", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = diff.Months()
		}
	})

	b.Run("Days", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = diff.Days()
		}
	})

	b.Run("String", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = diff.String()
		}
	})

	b.Run("ForHumans", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = diff.ForHumans()
		}
	})
}
