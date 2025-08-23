package chronogo

import (
	"sync"
	"testing"
	"time"
)

// DST performance benchmarks

func BenchmarkIsDSTOriginal(b *testing.B) {
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

func BenchmarkIsDSTOptimized(b *testing.B) {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		b.Skip("NY tz not available")
	}
	summer := Date(2023, time.July, 15, 12, 0, 0, 0, ny)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = summer.IsDSTOptimized()
	}
}

func BenchmarkIsDSTComparison(b *testing.B) {
	locations := []*time.Location{time.UTC}

	// Add more locations if available
	if ny, err := time.LoadLocation("America/New_York"); err == nil {
		locations = append(locations, ny)
	}
	if london, err := time.LoadLocation("Europe/London"); err == nil {
		locations = append(locations, london)
	}
	if tokyo, err := time.LoadLocation("Asia/Tokyo"); err == nil {
		locations = append(locations, tokyo)
	}

	dates := []DateTime{}
	for _, loc := range locations {
		dates = append(dates,
			Date(2023, time.January, 15, 12, 0, 0, 0, loc), // Winter
			Date(2023, time.July, 15, 12, 0, 0, 0, loc),    // Summer
		)
	}

	b.Run("Original", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			dt := dates[i%len(dates)]
			_ = dt.IsDST()
		}
	})

	b.Run("Optimized", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			dt := dates[i%len(dates)]
			_ = dt.IsDSTOptimized()
		}
	})
}

func BenchmarkIsDSTColdCache(b *testing.B) {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		b.Skip("NY tz not available")
	}

	b.Run("Original", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// Clear caches to simulate cold start
			standardOffsetCache = sync.Map{}
			summer := Date(2023, time.July, 15, 12, 0, 0, 0, ny)
			_ = summer.IsDST()
		}
	})

	b.Run("Optimized", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// Clear caches to simulate cold start
			ClearDSTCache()
			summer := Date(2023, time.July, 15, 12, 0, 0, 0, ny)
			_ = summer.IsDSTOptimized()
		}
	})
}

func BenchmarkIsDSTConcurrent(b *testing.B) {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		b.Skip("NY tz not available")
	}

	dates := []DateTime{
		Date(2023, time.January, 15, 12, 0, 0, 0, ny),
		Date(2023, time.July, 15, 12, 0, 0, 0, ny),
		Date(2023, time.April, 15, 12, 0, 0, 0, ny),
		Date(2023, time.October, 15, 12, 0, 0, 0, ny),
	}

	b.Run("Original", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				dt := dates[i%len(dates)]
				_ = dt.IsDST()
				i++
			}
		})
	})

	b.Run("Optimized", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				dt := dates[i%len(dates)]
				_ = dt.IsDSTOptimized()
				i++
			}
		})
	})
}
