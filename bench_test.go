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
