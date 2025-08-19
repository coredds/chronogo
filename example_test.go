package chronogo_test

import (
	"fmt"
	"time"

	chronogo "github.com/coredds/ChronoGo"
)

// Example demonstrates basic usage of ChronoGo similar to the PRD sample.
func Example() {
	// Load timezone
	loc, _ := chronogo.LoadLocation("America/New_York")

	// Get current time
	now := chronogo.Now()
	fmt.Println("Current time:", now.ToDateTimeString())

	// Create custom datetime
	dt := chronogo.Date(2025, time.July, 30, 11, 0, 0, 0, loc)
	fmt.Println("Custom datetime:", dt.String())

	// Add time
	dt2 := dt.AddDays(5).AddHours(3)
	fmt.Println("After 5 days 3 hours:", dt2.String())

	// Convert to UTC
	fmt.Println("In UTC:", dt2.UTC().String())

	// Check if DST is in effect
	if dt2.IsDST() {
		fmt.Println("DST is in effect")
	}

	// Get duration difference
	diff := dt2.Sub(now)
	fmt.Println("Duration from now:", diff)

	// Human-friendly difference
	fmt.Println("Diff for humans:", dt2.DiffForHumans(now))

	// Custom formatting
	fmt.Println("Formatted:", dt2.Format("2006-01-02 15:04:05 MST"))
}

// Example_parsing demonstrates various parsing capabilities.
func Example_parsing() {
	// Parse ISO 8601
	dt1, _ := chronogo.ParseISO8601("2023-12-25T15:30:45Z")
	fmt.Println("Parsed ISO 8601:", dt1.String())

	// Parse common formats
	dt2, _ := chronogo.Parse("2023-12-25 15:30:45")
	fmt.Println("Parsed datetime:", dt2.String())

	// Parse with custom format
	dt3, _ := chronogo.FromFormat("25/12/2023 15:30", "02/01/2006 15:04")
	fmt.Println("Parsed custom format:", dt3.String())

	// Parse Unix timestamp
	dt4, _ := chronogo.TryParseUnix("1703516445")
	fmt.Println("From Unix timestamp:", dt4.String())
}

// Example_period demonstrates Period usage for time intervals.
func Example_period() {
	start := chronogo.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := chronogo.Date(2023, time.January, 10, 0, 0, 0, 0, time.UTC)

	period := chronogo.NewPeriod(start, end)
	fmt.Println("Period duration:", period.String())
	fmt.Println("Days in period:", period.Days())
	fmt.Println("Hours in period:", period.Hours())

	// Check if a date is in the period
	middle := chronogo.Date(2023, time.January, 5, 0, 0, 0, 0, time.UTC)
	fmt.Println("Contains middle date:", period.Contains(middle))

	// Iterate over period
	fmt.Println("Dates in period (every 2 days):")
	count := 0
	for dt := range period.Range("days", 2) {
		if count < 3 { // Limit output for example
			fmt.Printf("  %s\n", dt.ToDateString())
		}
		count++
	}
}

// Example_timezones demonstrates timezone handling.
func Example_timezones() {
	// Create datetime in different timezones
	utc := chronogo.UTC(2023, time.December, 25, 12, 0, 0, 0)
	fmt.Println("UTC time:", utc.String())

	// Convert to different timezones
	if ny, err := chronogo.LoadLocation("America/New_York"); err == nil {
		nyTime := utc.In(ny)
		fmt.Println("New York time:", nyTime.String())
		fmt.Println("Is DST in NY:", nyTime.IsDST())
	}

	if tokyo, err := chronogo.LoadLocation("Asia/Tokyo"); err == nil {
		tokyoTime := utc.In(tokyo)
		fmt.Println("Tokyo time:", tokyoTime.String())
	}

	if london, err := chronogo.LoadLocation("Europe/London"); err == nil {
		londonTime := utc.In(london)
		fmt.Println("London time:", londonTime.String())
	}
}

// Example_arithmetic demonstrates date/time arithmetic.
func Example_arithmetic() {
	dt := chronogo.Date(2023, time.January, 15, 12, 30, 45, 0, time.UTC)
	fmt.Println("Original:", dt.String())

	// Add various units
	fmt.Println("Add 1 year:", dt.AddYears(1).String())
	fmt.Println("Add 3 months:", dt.AddMonths(3).String())
	fmt.Println("Add 10 days:", dt.AddDays(10).String())
	fmt.Println("Add 5 hours:", dt.AddHours(5).String())

	// Subtract units
	fmt.Println("Subtract 2 weeks:", dt.AddDays(-14).String())

	// Method chaining
	future := dt.AddYears(1).AddMonths(6).AddDays(15)
	fmt.Println("Chained operations:", future.String())

	// Set specific components
	newYear := dt.SetYear(2025)
	fmt.Println("Set year to 2025:", newYear.String())

	newTime := dt.SetHour(18).SetMinute(0).SetSecond(0)
	fmt.Println("Set time to 18:00:00:", newTime.String())
}

// Example_comparison demonstrates date/time comparisons.
func Example_comparison() {
	dt1 := chronogo.Date(2023, time.January, 15, 12, 0, 0, 0, time.UTC)
	dt2 := chronogo.Date(2023, time.January, 16, 12, 0, 0, 0, time.UTC)
	dt3 := chronogo.Date(2023, time.January, 15, 12, 0, 0, 0, time.UTC)

	fmt.Println("dt1:", dt1.ToDateString())
	fmt.Println("dt2:", dt2.ToDateString())
	fmt.Println("dt3:", dt3.ToDateString())

	fmt.Println("dt1.Before(dt2):", dt1.Before(dt2))
	fmt.Println("dt1.After(dt2):", dt1.After(dt2))
	fmt.Println("dt1.Equal(dt3):", dt1.Equal(dt3))

	// Convenience methods
	now := chronogo.Now()
	past := now.AddDays(-1)
	future := now.AddDays(1)

	fmt.Println("Past is past:", past.IsPast())
	fmt.Println("Future is future:", future.IsFuture())

	// Leap year check
	fmt.Println("2024 is leap year:", chronogo.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC).IsLeapYear())
	fmt.Println("2023 is leap year:", chronogo.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC).IsLeapYear())
}

// Example_formatting demonstrates various formatting options.
func Example_formatting() {
	dt := chronogo.Date(2023, time.December, 25, 15, 30, 45, 0, time.UTC)

	// Built-in formats
	fmt.Println("Date string:", dt.ToDateString())
	fmt.Println("Time string:", dt.ToTimeString())
	fmt.Println("DateTime string:", dt.ToDateTimeString())
	fmt.Println("ISO 8601:", dt.ToISO8601String())

	// Custom formats using Go's layout
	fmt.Println("Custom format 1:", dt.Format("Monday, January 2, 2006"))
	fmt.Println("Custom format 2:", dt.Format("02/01/2006 15:04"))
	fmt.Println("Custom format 3:", dt.Format("Jan 2, 2006 at 3:04 PM"))

	// Human-friendly formats
	fmt.Println("Age:", dt.Age())
	fmt.Println("Time ago:", dt.DiffForHumans())
}

// Example_humanReadable demonstrates human-friendly time representations.
func Example_humanReadable() {
	now := chronogo.Now()

	// Various time differences
	times := []chronogo.DateTime{
		now.AddSeconds(-30),
		now.AddMinutes(-5),
		now.AddHours(-2),
		now.AddDays(-1),
		now.AddDays(-7),
		now.AddDays(-30),
		now.AddYears(-1),
		now.AddSeconds(30),
		now.AddMinutes(5),
		now.AddHours(2),
		now.AddDays(1),
		now.AddDays(7),
		now.AddDays(30),
		now.AddYears(1),
	}

	fmt.Println("Human-readable time differences:")
	for _, t := range times {
		fmt.Printf("  %s\n", t.DiffForHumans(now))
	}

	// Duration humanization
	durations := []time.Duration{
		30 * time.Second,
		5 * time.Minute,
		2 * time.Hour,
		25 * time.Hour,
		7 * 24 * time.Hour,
	}

	fmt.Println("\nHumanized durations:")
	for _, d := range durations {
		fmt.Printf("  %v = %s\n", d, chronogo.Humanize(d))
	}
}
