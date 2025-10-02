package chronogo

import (
	"testing"
	"time"
)

func TestPeriodOverlaps(t *testing.T) {
	tests := []struct {
		name     string
		p1Start  DateTime
		p1End    DateTime
		p2Start  DateTime
		p2End    DateTime
		expected bool
	}{
		{
			name:     "Complete overlap",
			p1Start:  Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			p1End:    Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			p2Start:  Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			p2End:    Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "No overlap - p1 before p2",
			p1Start:  Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			p1End:    Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			p2Start:  Date(2024, 1, 11, 0, 0, 0, 0, time.UTC),
			p2End:    Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "No overlap - p2 before p1",
			p1Start:  Date(2024, 1, 11, 0, 0, 0, 0, time.UTC),
			p1End:    Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			p2Start:  Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			p2End:    Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "Touching at boundary",
			p1Start:  Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			p1End:    Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			p2Start:  Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			p2End:    Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "p1 encompasses p2",
			p1Start:  Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			p1End:    Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
			p2Start:  Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			p2End:    Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p1 := NewPeriod(tt.p1Start, tt.p1End)
			p2 := NewPeriod(tt.p2Start, tt.p2End)

			result := p1.Overlaps(p2)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}

			// Test symmetry
			result2 := p2.Overlaps(p1)
			if result2 != tt.expected {
				t.Errorf("Overlaps is not symmetric: p1.Overlaps(p2)=%v, p2.Overlaps(p1)=%v", result, result2)
			}
		})
	}
}

func TestPeriodGap(t *testing.T) {
	// Two periods with a gap
	p1 := NewPeriod(
		Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
	)
	p2 := NewPeriod(
		Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
	)

	gap := p1.Gap(p2)

	if gap.Start.Day() != 10 || gap.End.Day() != 15 {
		t.Errorf("Expected gap from Jan 10 to Jan 15, got %v to %v",
			gap.Start.Format("2006-01-02"), gap.End.Format("2006-01-02"))
	}

	// Test symmetry
	gap2 := p2.Gap(p1)
	if gap2.Start.Day() != 10 || gap2.End.Day() != 15 {
		t.Error("Gap is not symmetric")
	}
}

func TestPeriodGapOverlapping(t *testing.T) {
	// Overlapping periods should have zero gap
	p1 := NewPeriod(
		Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	)
	p2 := NewPeriod(
		Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
	)

	gap := p1.Gap(p2)

	if !gap.Start.IsZero() || !gap.End.IsZero() {
		t.Error("Expected zero gap for overlapping periods")
	}
}

func TestPeriodEncompasses(t *testing.T) {
	outer := NewPeriod(
		Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
	)
	inner := NewPeriod(
		Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
		Date(2024, 6, 30, 0, 0, 0, 0, time.UTC),
	)

	if !outer.Encompasses(inner) {
		t.Error("Expected outer to encompass inner")
	}

	if inner.Encompasses(outer) {
		t.Error("Expected inner NOT to encompass outer")
	}

	// Period encompasses itself
	if !outer.Encompasses(outer) {
		t.Error("Expected period to encompass itself")
	}
}

func TestPeriodMerge(t *testing.T) {
	p1 := NewPeriod(
		Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	)
	p2 := NewPeriod(
		Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
	)

	merged := p1.Merge(p2)

	expectedStart := Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedEnd := Date(2024, 1, 25, 0, 0, 0, 0, time.UTC)

	if !merged.Start.Equal(expectedStart) {
		t.Errorf("Expected merged start %v, got %v", expectedStart.Format("2006-01-02"), merged.Start.Format("2006-01-02"))
	}

	if !merged.End.Equal(expectedEnd) {
		t.Errorf("Expected merged end %v, got %v", expectedEnd.Format("2006-01-02"), merged.End.Format("2006-01-02"))
	}
}

func TestPeriodMergeNonOverlapping(t *testing.T) {
	p1 := NewPeriod(
		Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
	)
	p2 := NewPeriod(
		Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
		Date(2024, 1, 30, 0, 0, 0, 0, time.UTC),
	)

	merged := p1.Merge(p2)

	// Should span from earliest to latest
	expectedStart := Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedEnd := Date(2024, 1, 30, 0, 0, 0, 0, time.UTC)

	if !merged.Start.Equal(expectedStart) || !merged.End.Equal(expectedEnd) {
		t.Errorf("Expected merged from %v to %v, got %v to %v",
			expectedStart.Format("2006-01-02"), expectedEnd.Format("2006-01-02"),
			merged.Start.Format("2006-01-02"), merged.End.Format("2006-01-02"))
	}
}

