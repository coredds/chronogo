package chronogo

import (
	"testing"
	"time"
)

// Optimized benchmark tests to measure performance improvements

func BenchmarkPeriodRangeByUnitDaysOptimized(b *testing.B) {
	start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC)
	p := NewPeriod(start, end)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		count := 0
		for range p.RangeByUnitSlice(UnitDay, 7) { // weekly steps
			count++
		}
		if count == 0 {
			b.Fatal("unexpected zero count")
		}
	}
}

func BenchmarkPeriodRangeByUnitSlice(b *testing.B) {
	start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC)
	p := NewPeriod(start, end)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		items := p.RangeByUnitSlice(UnitDay, 7) // weekly steps
		if len(items) == 0 {
			b.Fatal("unexpected zero count")
		}
	}
}

func BenchmarkPeriodFastRangeDays(b *testing.B) {
	start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC)
	p := NewPeriod(start, end)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		items := p.FastRangeDays(7) // weekly steps
		if len(items) == 0 {
			b.Fatal("unexpected zero count")
		}
	}
}

// Benchmark memory allocation patterns
func BenchmarkPeriodMemoryOptimization(b *testing.B) {
	start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := Date(2023, time.January, 31, 0, 0, 0, 0, time.UTC) // Smaller range
	p := NewPeriod(start, end)

	b.Run("Original", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			count := 0
			for range p.RangeByUnit(UnitDay, 1) {
				count++
				if count > 31 { // Safety limit
					break
				}
			}
		}
	})

	b.Run("Optimized", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			count := 0
			for range p.RangeByUnitSlice(UnitDay, 1) {
				count++
				if count > 31 { // Safety limit
					break
				}
			}
		}
	})

	b.Run("Slice", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			items := p.RangeByUnitSlice(UnitDay, 1)
			_ = len(items) // Use the result
		}
	})
}

// Test concurrent period operations
func BenchmarkPeriodConcurrent(b *testing.B) {
	start := Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := Date(2023, time.March, 31, 0, 0, 0, 0, time.UTC)
	p := NewPeriod(start, end)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			items := p.FastRangeDays(7)
			_ = len(items)
		}
	})
}
