package chronogo

import (
	"testing"
	"time"
)

func BenchmarkParseISO8601(b *testing.B) {
	s := "2023-12-25T15:30:45Z"
	for i := 0; i < b.N; i++ {
		if _, err := ParseISO8601(s); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDiffForHumans(b *testing.B) {
	base := Date(2023, time.December, 25, 12, 0, 0, 0, time.UTC)
	other := base.AddDays(42).AddHours(3)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = other.DiffForHumans(base)
	}
}

func BenchmarkIsDST(b *testing.B) {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		b.Skip("NY tz not available")
	}
	summer := Date(2023, time.July, 15, 12, 0, 0, 0, ny)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = summer.IsDST()
	}
}

func BenchmarkPeriodRangeByUnitDays(b *testing.B) {
	start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC)
	p := NewPeriod(start, end)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		count := 0
		for range p.RangeByUnit(UnitDay, 7) { // weekly steps
			count++
		}
		if count == 0 {
			b.Fatal("unexpected zero count")
		}
	}
}
