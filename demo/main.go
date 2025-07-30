// Demo program showing ChronoGo functionality
package main

import (
	"fmt"
	"time"

	"github.com/coredds/ChronoGo"
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

	fmt.Println("=== Demo Complete ===")
}
