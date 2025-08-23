package chronogo

import (
	"context"
	"testing"
	"time"
)

// Additional comprehensive benchmarks for performance monitoring

func BenchmarkParse(b *testing.B) {
	inputs := []string{
		"2023-12-25T15:30:45Z",
		"2023/12/25",
		"2023-12-25",
		"15:30:45",
		"2023-12-25 15:30:45",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input := inputs[i%len(inputs)]
		if _, err := Parse(input); err != nil {
			// Skip invalid inputs in benchmark
			continue
		}
	}
}

func BenchmarkFromFormat(b *testing.B) {
	value := "25/12/2023 15:30"
	layout := "02/01/2006 15:04"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := FromFormat(value, layout); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTimezoneConversion(b *testing.B) {
	utc := Date(2023, time.December, 25, 15, 30, 45, 0, time.UTC)
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		b.Skip("NY timezone not available")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = utc.In(ny)
	}
}

func BenchmarkDateArithmetic(b *testing.B) {
	base := Date(2023, time.December, 25, 12, 0, 0, 0, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = base.AddDays(1).AddHours(2).AddMinutes(30)
	}
}

func BenchmarkHumanDiffCalculation(b *testing.B) {
	base := Date(2023, time.December, 25, 12, 0, 0, 0, time.UTC)
	future := base.AddDays(42).AddHours(3)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = future.DiffForHumans(base)
	}
}

func BenchmarkPeriodCreation(b *testing.B) {
	start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewPeriod(start, end)
	}
}

func BenchmarkPeriodIteration(b *testing.B) {
	start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := Date(2023, time.January, 31, 0, 0, 0, 0, time.UTC)
	p := NewPeriod(start, end)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		count := 0
		for range p.RangeDays() {
			count++
			if count > 31 { // Safety break
				break
			}
		}
	}
}

func BenchmarkBusinessDayCalculations(b *testing.B) {
	checker := NewUSHolidayChecker()
	start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = start.AddBusinessDays(10, checker)
	}
}

func BenchmarkHolidayChecking(b *testing.B) {
	checker := NewUSHolidayChecker()
	dates := []DateTime{
		Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),   // New Year
		Date(2023, time.July, 4, 0, 0, 0, 0, time.UTC),      // Independence Day
		Date(2023, time.December, 25, 0, 0, 0, 0, time.UTC), // Christmas
		Date(2023, time.March, 15, 0, 0, 0, 0, time.UTC),    // Regular day
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		date := dates[i%len(dates)]
		_ = date.IsHoliday(checker)
	}
}

func BenchmarkDateFormatting(b *testing.B) {
	dt := Date(2023, time.December, 25, 15, 30, 45, 123456789, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dt.Format("2006-01-02 15:04:05.000000000 -0700 MST")
	}
}

func BenchmarkJSONMarshal(b *testing.B) {
	dt := Date(2023, time.December, 25, 15, 30, 45, 0, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := dt.MarshalJSON(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONUnmarshal(b *testing.B) {
	data := []byte(`"2023-12-25T15:30:45Z"`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var dt DateTime
		if err := dt.UnmarshalJSON(data); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConcurrentOperations(b *testing.B) {
	dt := Date(2023, time.December, 25, 15, 30, 45, 0, time.UTC)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = dt.AddDays(1).Format("2006-01-02")
		}
	})
}

func BenchmarkMemoryAllocation(b *testing.B) {
	start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := Date(2023, time.January, 7, 0, 0, 0, 0, time.UTC)
	p := NewPeriod(start, end)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
		count := 0
		for range p.RangeWithContext(ctx, "days", 1) {
			count++
			if count > 10 { // Safety limit
				break
			}
		}
		cancel()
	}
}

// Benchmark comparison with standard library
func BenchmarkCompareStdLibTimeNow(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = time.Now()
	}
}

func BenchmarkCompareChronoGoNow(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Now()
	}
}

func BenchmarkCompareStdLibTimeParse(b *testing.B) {
	layout := "2006-01-02T15:04:05Z"
	value := "2023-12-25T15:30:45Z"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := time.Parse(layout, value); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareChronoGoParse(b *testing.B) {
	value := "2023-12-25T15:30:45Z"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := ParseISO8601(value); err != nil {
			b.Fatal(err)
		}
	}
}
