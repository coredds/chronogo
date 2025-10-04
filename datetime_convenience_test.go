package chronogo

import (
	"testing"
	"time"
)

// TestIsLongYear tests the IsLongYear() method for ISO 8601 long year detection
func TestIsLongYear(t *testing.T) {
	tests := []struct {
		year     int
		expected bool
		reason   string
	}{
		{2004, true, "2004 starts on Thursday"},
		{2009, true, "2009 starts on Thursday"},
		{2015, true, "2015 starts on Thursday"},
		{2020, true, "2020 is leap year starting on Wednesday"},
		{2026, true, "2026 starts on Thursday"},
		{2032, true, "2032 is leap year starting on Thursday"},
		{2037, true, "2037 starts on Thursday"},
		{2043, true, "2043 starts on Thursday"},
		{2048, true, "2048 is leap year starting on Thursday"},
		{2054, true, "2054 starts on Thursday"},
		{2060, true, "2060 is leap year starting on Thursday"},
		{2065, true, "2065 starts on Thursday"},
		{2071, true, "2071 starts on Thursday"},
		{2076, true, "2076 is leap year starting on Thursday"},

		// Non-long years
		{2019, false, "2019 starts on Tuesday"},
		{2021, false, "2021 starts on Friday"},
		{2022, false, "2022 starts on Saturday"},
		{2023, false, "2023 starts on Sunday"},
		{2024, false, "2024 is leap year but starts on Monday"},
		{2025, false, "2025 starts on Wednesday but not leap"},
		{2027, false, "2027 starts on Friday"},
		{2028, false, "2028 is leap year but starts on Saturday"},
		{2029, false, "2029 starts on Monday"},
		{2030, false, "2030 starts on Tuesday"},
	}

	for _, tt := range tests {
		t.Run(tt.reason, func(t *testing.T) {
			dt := Date(tt.year, time.January, 1, 0, 0, 0, 0, time.UTC)
			result := dt.IsLongYear()

			if result != tt.expected {
				t.Errorf("IsLongYear() for %d = %v, want %v (%s)",
					tt.year, result, tt.expected, tt.reason)
			}
		})
	}
}

// TestIsLongYearEdgeCases tests edge cases for long year detection
func TestIsLongYearEdgeCases(t *testing.T) {
	t.Run("Different dates in same year give same result", func(t *testing.T) {
		year := 2020
		dates := []time.Time{
			time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(year, 6, 15, 0, 0, 0, 0, time.UTC),
			time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC),
		}

		expected := true // 2020 is a long year
		for _, date := range dates {
			dt := DateTime{date}
			if result := dt.IsLongYear(); result != expected {
				t.Errorf("IsLongYear() for %v = %v, want %v", date, result, expected)
			}
		}
	})

	t.Run("Different timezones give same result", func(t *testing.T) {
		year := 2015
		ny, _ := time.LoadLocation("America/New_York")
		tokyo, _ := time.LoadLocation("Asia/Tokyo")

		dt1 := Date(year, 6, 1, 12, 0, 0, 0, time.UTC)
		dt2 := Date(year, 6, 1, 12, 0, 0, 0, ny)
		dt3 := Date(year, 6, 1, 12, 0, 0, 0, tokyo)

		expected := true // 2015 is a long year
		if dt1.IsLongYear() != expected || dt2.IsLongYear() != expected || dt3.IsLongYear() != expected {
			t.Errorf("IsLongYear() should be consistent across timezones")
		}
	})
}

// TestOn tests the On() convenience method
func TestOn(t *testing.T) {
	base := Date(2023, time.March, 15, 14, 30, 45, 123456789, time.UTC)

	tests := []struct {
		name     string
		year     int
		month    time.Month
		day      int
		wantDate string
		wantTime string
	}{
		{
			name:     "Change to New Year's Day",
			year:     2024,
			month:    time.January,
			day:      1,
			wantDate: "2024-01-01",
			wantTime: "14:30:45",
		},
		{
			name:     "Change to leap day",
			year:     2024,
			month:    time.February,
			day:      29,
			wantDate: "2024-02-29",
			wantTime: "14:30:45",
		},
		{
			name:     "Change to end of year",
			year:     2025,
			month:    time.December,
			day:      31,
			wantDate: "2025-12-31",
			wantTime: "14:30:45",
		},
		{
			name:     "Change only year",
			year:     2030,
			month:    time.March,
			day:      15,
			wantDate: "2030-03-15",
			wantTime: "14:30:45",
		},
		{
			name:     "Change to same date",
			year:     2023,
			month:    time.March,
			day:      15,
			wantDate: "2023-03-15",
			wantTime: "14:30:45",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := base.On(tt.year, tt.month, tt.day)

			// Check date changed
			gotDate := result.Format("2006-01-02")
			if gotDate != tt.wantDate {
				t.Errorf("On() date = %v, want %v", gotDate, tt.wantDate)
			}

			// Check time unchanged
			gotTime := result.Format("15:04:05")
			if gotTime != tt.wantTime {
				t.Errorf("On() time = %v, want %v", gotTime, tt.wantTime)
			}

			// Check nanoseconds unchanged
			if result.Nanosecond() != base.Nanosecond() {
				t.Errorf("On() nanosecond = %v, want %v", result.Nanosecond(), base.Nanosecond())
			}

			// Check timezone unchanged
			if result.Location() != base.Location() {
				t.Errorf("On() location = %v, want %v", result.Location(), base.Location())
			}
		})
	}
}

// TestAt tests the At() convenience method
func TestAt(t *testing.T) {
	base := Date(2023, time.March, 15, 14, 30, 45, 123456789, time.UTC)

	tests := []struct {
		name     string
		hour     int
		minute   int
		second   int
		wantDate string
		wantTime string
	}{
		{
			name:     "Set to midnight",
			hour:     0,
			minute:   0,
			second:   0,
			wantDate: "2023-03-15",
			wantTime: "00:00:00",
		},
		{
			name:     "Set to noon",
			hour:     12,
			minute:   0,
			second:   0,
			wantDate: "2023-03-15",
			wantTime: "12:00:00",
		},
		{
			name:     "Set to afternoon",
			hour:     15,
			minute:   45,
			second:   30,
			wantDate: "2023-03-15",
			wantTime: "15:45:30",
		},
		{
			name:     "Set to end of day",
			hour:     23,
			minute:   59,
			second:   59,
			wantDate: "2023-03-15",
			wantTime: "23:59:59",
		},
		{
			name:     "Set to same time",
			hour:     14,
			minute:   30,
			second:   45,
			wantDate: "2023-03-15",
			wantTime: "14:30:45",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := base.At(tt.hour, tt.minute, tt.second)

			// Check date unchanged
			gotDate := result.Format("2006-01-02")
			if gotDate != tt.wantDate {
				t.Errorf("At() date = %v, want %v", gotDate, tt.wantDate)
			}

			// Check time changed
			gotTime := result.Format("15:04:05")
			if gotTime != tt.wantTime {
				t.Errorf("At() time = %v, want %v", gotTime, tt.wantTime)
			}

			// Check nanoseconds unchanged
			if result.Nanosecond() != base.Nanosecond() {
				t.Errorf("At() nanosecond = %v, want %v", result.Nanosecond(), base.Nanosecond())
			}

			// Check timezone unchanged
			if result.Location() != base.Location() {
				t.Errorf("At() location = %v, want %v", result.Location(), base.Location())
			}
		})
	}
}

// TestOnAndAtChaining tests that On() and At() can be chained
func TestOnAndAtChaining(t *testing.T) {
	base := Now()

	// Chain On() and At()
	result := base.On(2024, time.June, 15).At(9, 30, 0)

	expected := "2024-06-15 09:30:00"
	got := result.Format("2006-01-02 15:04:05")

	if got != expected {
		t.Errorf("On().At() chaining = %v, want %v", got, expected)
	}

	// Reverse order
	result2 := base.At(18, 45, 30).On(2025, time.December, 25)
	expected2 := "2025-12-25 18:45:30"
	got2 := result2.Format("2006-01-02 15:04:05")

	if got2 != expected2 {
		t.Errorf("At().On() chaining = %v, want %v", got2, expected2)
	}
}

// TestOnWithDifferentTimezones tests On() preserves timezone
func TestOnWithDifferentTimezones(t *testing.T) {
	ny, _ := time.LoadLocation("America/New_York")
	tokyo, _ := time.LoadLocation("Asia/Tokyo")

	dtNY := Date(2023, time.June, 1, 14, 0, 0, 0, ny)
	dtTokyo := Date(2023, time.June, 1, 14, 0, 0, 0, tokyo)

	resultNY := dtNY.On(2024, time.January, 15)
	resultTokyo := dtTokyo.On(2024, time.January, 15)

	if resultNY.Location().String() != "America/New_York" {
		t.Errorf("On() should preserve New York timezone")
	}

	if resultTokyo.Location().String() != "Asia/Tokyo" {
		t.Errorf("On() should preserve Tokyo timezone")
	}
}

// TestAtWithDifferentTimezones tests At() preserves timezone
func TestAtWithDifferentTimezones(t *testing.T) {
	ny, _ := time.LoadLocation("America/New_York")
	tokyo, _ := time.LoadLocation("Asia/Tokyo")

	dtNY := Date(2023, time.June, 1, 14, 0, 0, 0, ny)
	dtTokyo := Date(2023, time.June, 1, 14, 0, 0, 0, tokyo)

	resultNY := dtNY.At(9, 30, 0)
	resultTokyo := dtTokyo.At(9, 30, 0)

	if resultNY.Location().String() != "America/New_York" {
		t.Errorf("At() should preserve New York timezone")
	}

	if resultTokyo.Location().String() != "Asia/Tokyo" {
		t.Errorf("At() should preserve Tokyo timezone")
	}
}

// TestOnVsSetComparison compares On() convenience with Set() chain
func TestOnVsSetComparison(t *testing.T) {
	base := Date(2023, time.March, 15, 14, 30, 45, 0, time.UTC)

	// Using On()
	result1 := base.On(2024, time.June, 20)

	// Using Set() chain
	result2 := base.Set().Year(2024).Month(time.June).Day(20).Build()

	if !result1.Equal(result2) {
		t.Errorf("On() and Set().Year().Month().Day() should produce same result")
	}
}

// TestAtVsSetComparison compares At() convenience with Set() chain
func TestAtVsSetComparison(t *testing.T) {
	base := Date(2023, time.March, 15, 14, 30, 45, 0, time.UTC)

	// Using At()
	result1 := base.At(9, 15, 30)

	// Using Set() chain
	result2 := base.Set().Hour(9).Minute(15).Second(30).Build()

	if !result1.Equal(result2) {
		t.Errorf("At() and Set().Hour().Minute().Second() should produce same result")
	}
}

// BenchmarkOn benchmarks the On() method
func BenchmarkOn(b *testing.B) {
	dt := Date(2023, time.March, 15, 14, 30, 45, 0, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dt.On(2024, time.June, 20)
	}
}

// BenchmarkOnVsSet benchmarks On() vs Set() chain
func BenchmarkOnVsSet(b *testing.B) {
	dt := Date(2023, time.March, 15, 14, 30, 45, 0, time.UTC)

	b.Run("On", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = dt.On(2024, time.June, 20)
		}
	})

	b.Run("Set chain", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = dt.Set().Year(2024).Month(time.June).Day(20).Build()
		}
	})
}

// BenchmarkAt benchmarks the At() method
func BenchmarkAt(b *testing.B) {
	dt := Date(2023, time.March, 15, 14, 30, 45, 0, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dt.At(9, 15, 30)
	}
}

// BenchmarkAtVsSet benchmarks At() vs Set() chain
func BenchmarkAtVsSet(b *testing.B) {
	dt := Date(2023, time.March, 15, 14, 30, 45, 0, time.UTC)

	b.Run("At", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = dt.At(9, 15, 30)
		}
	})

	b.Run("Set chain", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = dt.Set().Hour(9).Minute(15).Second(30).Build()
		}
	})
}
