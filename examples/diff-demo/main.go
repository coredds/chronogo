package main

import (
	"fmt"
	"time"

	"github.com/coredds/chronogo"
)

func main() {
	fmt.Println("ChronoGo Diff Type Demo")
	fmt.Println("=======================")
	fmt.Println()

	// Example 1: Basic Diff creation
	fmt.Println("1. Basic Diff Creation")
	fmt.Println("----------------------")
	dt1 := chronogo.Date(2023, time.January, 15, 10, 0, 0, 0, time.UTC)
	dt2 := chronogo.Date(2024, time.March, 20, 14, 30, 0, 0, time.UTC)

	diff := dt2.Diff(dt1)
	fmt.Printf("Start: %s\n", dt1.Format("2006-01-02 15:04:05"))
	fmt.Printf("End:   %s\n", dt2.Format("2006-01-02 15:04:05"))
	fmt.Printf("Diff:  %s\n", diff.String())
	fmt.Println()

	// Example 2: Calendar-Aware Differences
	fmt.Println("2. Calendar-Aware Differences")
	fmt.Println("------------------------------")
	fmt.Printf("Years:  %d\n", diff.Years())
	fmt.Printf("Months: %d\n", diff.Months())
	fmt.Printf("Weeks:  %d\n", diff.Weeks())
	fmt.Printf("Days:   %d\n", diff.Days())
	fmt.Printf("Hours:  %d\n", diff.Hours())
	fmt.Println()

	// Example 3: Precise Differences (with fractional parts)
	fmt.Println("3. Precise Differences (Float)")
	fmt.Println("-------------------------------")
	fmt.Printf("In Years:  %.2f\n", diff.InYears())
	fmt.Printf("In Months: %.2f\n", diff.InMonths())
	fmt.Printf("In Weeks:  %.2f\n", diff.InWeeks())
	fmt.Printf("In Days:   %.2f\n", diff.InDays())
	fmt.Printf("In Hours:  %.2f\n", diff.InHours())
	fmt.Println()

	// Example 4: Human-Readable Strings
	fmt.Println("4. Human-Readable Strings")
	fmt.Println("-------------------------")
	fmt.Printf("For Humans:       %s\n", diff.ForHumans())
	fmt.Printf("For Comparison:   %s\n", diff.ForHumansComparison())
	fmt.Printf("Compact:          %s\n", diff.CompactString())
	fmt.Println()

	// Example 5: Absolute Differences
	fmt.Println("5. Absolute Differences")
	fmt.Println("-----------------------")
	pastDiff := dt1.Diff(dt2) // dt1 < dt2, so negative
	fmt.Printf("Negative diff:  %s\n", pastDiff.CompactString())
	fmt.Printf("Is negative:    %v\n", pastDiff.IsNegative())
	fmt.Printf("Absolute:       %s\n", pastDiff.Abs().CompactString())
	fmt.Printf("Using DiffAbs:  %s\n", dt1.DiffAbs(dt2).CompactString())
	fmt.Println()

	// Example 6: Invert Differences
	fmt.Println("6. Invert Differences")
	fmt.Println("---------------------")
	original := dt2.Diff(dt1)
	inverted := original.Invert()
	fmt.Printf("Original: %s (%s)\n", original.ForHumansComparison(), original.CompactString())
	fmt.Printf("Inverted: %s (%s)\n", inverted.ForHumansComparison(), inverted.CompactString())
	fmt.Println()

	// Example 7: Comparing Differences
	fmt.Println("7. Comparing Differences")
	fmt.Println("------------------------")
	dt3 := chronogo.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	dt4 := chronogo.Date(2023, time.January, 10, 0, 0, 0, 0, time.UTC)
	dt5 := chronogo.Date(2023, time.January, 20, 0, 0, 0, 0, time.UTC)

	diff1 := dt4.Diff(dt3) // 9 days
	diff2 := dt5.Diff(dt3) // 19 days

	fmt.Printf("Diff1: %s\n", diff1.CompactString())
	fmt.Printf("Diff2: %s\n", diff2.CompactString())
	fmt.Printf("Diff1 shorter than Diff2: %v\n", diff1.ShorterThan(diff2))
	fmt.Printf("Diff2 longer than Diff1:  %v\n", diff2.LongerThan(diff1))
	fmt.Printf("Diff1 equals Diff1:       %v\n", diff1.EqualTo(diff1))
	fmt.Println()

	// Example 8: Real-World Use Case - Age Calculator
	fmt.Println("8. Real-World Use Case: Age Calculator")
	fmt.Println("---------------------------------------")
	birthdate := chronogo.Date(1990, time.June, 15, 0, 0, 0, 0, time.UTC)
	now := chronogo.Now()
	age := now.Diff(birthdate)

	fmt.Printf("Birthdate: %s\n", birthdate.Format("January 2, 2006"))
	fmt.Printf("Today:     %s\n", now.Format("January 2, 2006"))
	fmt.Printf("Age:       %d years old\n", age.Years())
	fmt.Printf("Or:        %d months old\n", age.Months())
	fmt.Printf("Or:        %d days old\n", age.Days())
	fmt.Printf("Human:     %s\n", age.ForHumans())
	fmt.Println()

	// Example 9: Project Timeline
	fmt.Println("9. Real-World Use Case: Project Timeline")
	fmt.Println("-----------------------------------------")
	projectStart := chronogo.Date(2023, time.September, 1, 9, 0, 0, 0, time.UTC)
	projectEnd := chronogo.Date(2024, time.March, 15, 17, 0, 0, 0, time.UTC)
	projectDiff := projectEnd.Diff(projectStart)

	fmt.Printf("Project Start:    %s\n", projectStart.Format("January 2, 2006"))
	fmt.Printf("Project End:      %s\n", projectEnd.Format("January 2, 2006"))
	fmt.Printf("Duration:         %s\n", projectDiff.String())
	fmt.Printf("Compact:          %s\n", projectDiff.CompactString())
	fmt.Printf("Total Work Days:  %.0f days\n", projectDiff.InDays())
	fmt.Printf("Total Work Hours: %.0f hours\n", projectDiff.InHours())
	fmt.Println()

	// Example 10: Countdown Timer
	fmt.Println("10. Real-World Use Case: Countdown Timer")
	fmt.Println("-----------------------------------------")
	newYear := chronogo.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	countdown := newYear.Diff(chronogo.Now())

	if countdown.IsPositive() {
		fmt.Printf("New Year 2026:    %s\n", newYear.Format("January 2, 2006 15:04:05"))
		fmt.Printf("Countdown:        %s\n", countdown.String())
		fmt.Printf("Compact:          %s\n", countdown.CompactString())
		fmt.Printf("Days remaining:   %d\n", countdown.Days())
		fmt.Printf("Hours remaining:  %d\n", countdown.Hours())
		fmt.Printf("Human:            %s\n", countdown.ForHumans())
	} else {
		fmt.Println("Happy New Year! 🎉")
	}
	fmt.Println()

	// Example 11: Accessing Duration and Period
	fmt.Println("11. Accessing Duration and Period")
	fmt.Println("----------------------------------")
	diffExample := dt2.Diff(dt1)
	duration := diffExample.Duration()
	period := diffExample.Period()

	fmt.Printf("As time.Duration: %v\n", duration)
	fmt.Printf("As Period Start:  %s\n", period.Start.Format("2006-01-02"))
	fmt.Printf("As Period End:    %s\n", period.End.Format("2006-01-02"))
	fmt.Printf("Period Days():    %d\n", period.Days())
	fmt.Println()

	// Example 12: Zero Difference
	fmt.Println("12. Zero Difference")
	fmt.Println("-------------------")
	same := chronogo.Date(2023, time.June, 15, 12, 0, 0, 0, time.UTC)
	zeroDiff := same.Diff(same)

	fmt.Printf("Is Zero:   %v\n", zeroDiff.IsZero())
	fmt.Printf("String:    %s\n", zeroDiff.String())
	fmt.Printf("Compact:   %s\n", zeroDiff.CompactString())
	fmt.Println()

	// Example 13: Performance Comparison
	fmt.Println("13. Performance Note")
	fmt.Println("--------------------")
	fmt.Println("✓ Diff type provides unified API for both precise and calendar-aware differences")
	fmt.Println("✓ Zero allocations for most operations")
	fmt.Println("✓ Lazy evaluation - only computes what you need")
	fmt.Println("✓ Chainable with all DateTime methods")
	fmt.Println()

	fmt.Println("Done! The Diff type unifies time.Duration and Period into one convenient API.")
}
