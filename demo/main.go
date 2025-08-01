// Demo program showing ChronoGo functionality
package main

import (
	"fmt"
	"time"

	chronogo "github.com/coredds/ChronoGo"
)

func main() {
	fmt.Println("=== ChronoGo Demo ===")
	fmt.Println()

	// Current time
	fmt.Println("1. Current Time:")
	now := chronogo.Now()
	fmt.Printf("   Now: %s\n", now)
	fmt.Printf("   Now UTC: %s\n", chronogo.NowUTC())
	fmt.Println()

	// Date creation
	fmt.Println("2. Date Creation:")
	christmas := chronogo.Date(2023, time.December, 25, 15, 30, 0, 0, time.UTC)
	fmt.Printf("   Christmas 2023: %s\n", christmas)
	fmt.Printf("   Today: %s\n", chronogo.Today().ToDateString())
	fmt.Printf("   Tomorrow: %s\n", chronogo.Tomorrow().ToDateString())
	fmt.Println()

	// Arithmetic
	fmt.Println("3. Date Arithmetic:")
	future := christmas.AddYears(1).AddMonths(6).AddDays(10)
	fmt.Printf("   Christmas + 1 year, 6 months, 10 days: %s\n", future)
	past := now.SubtractDays(30).SubtractHours(5)
	fmt.Printf("   30 days and 5 hours ago: %s\n", past)
	fmt.Println()

	// Parsing
	fmt.Println("4. Parsing:")
	if parsed, err := chronogo.Parse("2023-12-25T15:30:45Z"); err == nil {
		fmt.Printf("   Parsed ISO: %s\n", parsed)
	}
	if parsed, err := chronogo.FromFormat("25/12/2023 15:30", "02/01/2006 15:04"); err == nil {
		fmt.Printf("   Parsed custom: %s\n", parsed)
	}
	fmt.Println()

	// Human-friendly
	fmt.Println("5. Human-Friendly:")
	fmt.Printf("   Christmas was: %s\n", christmas.DiffForHumans())
	fmt.Printf("   Christmas age: %s\n", christmas.Age())
	duration := 2*time.Hour + 30*time.Minute
	fmt.Printf("   Duration: %s\n", chronogo.Humanize(duration))
	fmt.Println()

	// Timezones
	fmt.Println("6. Timezones:")
	if ny, err := chronogo.LoadLocation("America/New_York"); err == nil {
		nyTime := christmas.In(ny)
		fmt.Printf("   Christmas in NY: %s\n", nyTime)
		fmt.Printf("   DST in NY: %t\n", nyTime.IsDST())
	}
	fmt.Println()

	// Periods
	fmt.Println("7. Periods:")
	start := chronogo.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := chronogo.Date(2023, time.January, 5, 0, 0, 0, 0, time.UTC)
	period := chronogo.NewPeriod(start, end)
	fmt.Printf("   Period: %s\n", period.String())
	fmt.Printf("   Duration: %d days\n", period.Days())
	fmt.Printf("   Contains Jan 3rd: %t\n", period.Contains(chronogo.Date(2023, time.January, 3, 0, 0, 0, 0, time.UTC)))

	fmt.Print("   Dates in period: ")
	count := 0
	for date := range period.RangeDays() {
		if count > 0 {
			fmt.Print(", ")
		}
		fmt.Print(date.ToDateString())
		count++
	}
	fmt.Println()
	fmt.Println()

	// Comparisons
	fmt.Println("8. Comparisons:")
	dt1 := chronogo.Date(2023, time.January, 15, 12, 0, 0, 0, time.UTC)
	dt2 := chronogo.Date(2023, time.January, 16, 12, 0, 0, 0, time.UTC)
	fmt.Printf("   %s before %s: %t\n", dt1.ToDateString(), dt2.ToDateString(), dt1.Before(dt2))
	fmt.Printf("   %s is past: %t\n", dt1.ToDateString(), dt1.IsPast())
	fmt.Printf("   2024 is leap year: %t\n", chronogo.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC).IsLeapYear())
	fmt.Println()

	// New API Features
	fmt.Println("9. New API Features:")

	// Utility methods
	testDate := chronogo.Date(2023, time.October, 15, 14, 30, 45, 0, time.UTC) // Sunday
	fmt.Printf("   Date: %s\n", testDate.ToDateString())
	fmt.Printf("   Start of day: %s\n", testDate.StartOfDay().ToDateTimeString())
	fmt.Printf("   End of day: %s\n", testDate.EndOfDay().ToDateTimeString())
	fmt.Printf("   Start of month: %s\n", testDate.StartOfMonth().ToDateString())
	fmt.Printf("   End of month: %s\n", testDate.EndOfMonth().ToDateString())
	fmt.Printf("   Start of quarter: %s\n", testDate.StartOfQuarter().ToDateString())
	fmt.Printf("   End of quarter: %s\n", testDate.EndOfQuarter().ToDateString())
	fmt.Printf("   Quarter: %d\n", testDate.Quarter())
	fmt.Printf("   Is weekend: %t\n", testDate.IsWeekend())
	fmt.Printf("   Is weekday: %t\n", testDate.IsWeekday())
	fmt.Printf("   Day of year: %d\n", testDate.DayOfYear())
	year, week := testDate.ISOWeek()
	fmt.Printf("   ISO week: %d-%d\n", year, week)
	fmt.Println()

	// Fluent API
	fmt.Println("10. Fluent API:")
	base := chronogo.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)

	// Fluent duration addition
	fluentResult := base.AddFluent().Years(1).Months(2).Days(10).Hours(5).Minutes(30).To(base)
	fmt.Printf("   Base date: %s\n", base.ToDateString())
	fmt.Printf("   After fluent addition: %s\n", fluentResult.ToDateTimeString())

	// Fluent setting
	fluentSet := base.Set().Year(2024).Month(time.December).Day(25).Hour(15).Minute(30).Build()
	fmt.Printf("   Fluent set result: %s\n", fluentSet.ToDateTimeString())
	fmt.Println()

	// Enhanced Duration
	fmt.Println("11. Enhanced Duration:")
	chronoDuration := chronogo.NewDuration(25*time.Hour + 30*time.Minute + 45*time.Second)
	fmt.Printf("   Duration: %s\n", chronoDuration.String())
	fmt.Printf("   Human readable: %s\n", chronoDuration.HumanString())
	fmt.Printf("   Days: %.2f\n", chronoDuration.Days())
	fmt.Printf("   Weeks: %.2f\n", chronoDuration.Weeks())
	fmt.Printf("   Is positive: %t\n", chronoDuration.IsPositive())

	duration2 := chronogo.NewDurationFromComponents(2, 15, 30)
	sum := chronoDuration.Add(duration2)
	fmt.Printf("   Sum with 2h15m30s: %s\n", sum.String())
	fmt.Println()

	fmt.Println("=== Demo Complete ===")
}
