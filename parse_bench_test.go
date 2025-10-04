package chronogo

import (
	"testing"
)

// Parse performance benchmarks

func BenchmarkParseComparison(b *testing.B) {
	inputs := []string{
		"2023-12-25T15:30:45Z",
		"2023-12-25T15:30:45+00:00",
		"2023-12-25 15:30:45",
		"2023-12-25",
		"2023/12/25",
		"15:30:45",
		"20231225",
	}

	b.Run("Original", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			input := inputs[i%len(inputs)]
			if _, err := Parse(input); err != nil {
				// Skip errors in benchmark
				continue
			}
		}
	})

	b.Run("Optimized", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			input := inputs[i%len(inputs)]
			if _, err := Parse(input); err != nil {
				// Skip errors in benchmark
				continue
			}
		}
	})
}

func BenchmarkParseISO8601Specific(b *testing.B) {
	value := "2023-12-25T15:30:45Z"

	b.Run("Original", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = Parse(value)
		}
	})

	b.Run("Optimized", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = Parse(value)
		}
	})

	b.Run("ParseISO8601", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = ParseISO8601(value)
		}
	})
}

func BenchmarkParseHeavyMix(b *testing.B) {
	inputs := []string{
		"2023-12-25T15:30:45Z",     // ISO 8601
		"2023-12-25 15:30:45",      // Space separator
		"2023-12-25",               // Date only
		"2023/12/25 15:30:45",      // Slash separator
		"15:30:45",                 // Time only
		"2023-1-2 3:04:05",         // Lenient format
		"20231225",                 // Compact
		"2023-12-25T15:30:45.123Z", // With milliseconds
	}

	b.Run("Original", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			input := inputs[i%len(inputs)]
			_, _ = Parse(input)
		}
	})

	b.Run("Optimized", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			input := inputs[i%len(inputs)]
			_, _ = Parse(input)
		}
	})
}

func BenchmarkBatchParsing(b *testing.B) {
	values := make([]string, 100)
	for i := 0; i < 100; i++ {
		values[i] = "2023-12-25T15:30:45Z"
	}

	b.Run("Individual", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, value := range values {
				_, _ = Parse(value)
			}
		}
	})
}
