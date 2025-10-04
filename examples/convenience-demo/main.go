package main

import (
	"fmt"
	"time"

	"github.com/coredds/chronogo"
)

func main() {
	fmt.Println("ChronoGo Convenience Methods Demo")
	fmt.Println("==================================")
	fmt.Println()

	// IsLongYear() - Check for ISO 8601 long years (53 weeks)
	fmt.Println("1. IsLongYear() - ISO 8601 Long Year Detection")
	fmt.Println("-----------------------------------------------")

	longYears := []int{2004, 2009, 2015, 2020, 2026}
	regularYears := []int{2019, 2021, 2022, 2023, 2024}

	fmt.Println("Long years (53 ISO weeks):")
	for _, year := range longYears {
		dt := chronogo.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
		if dt.IsLongYear() {
			fmt.Printf("  ✓ %d is a long year\n", year)
		}
	}

	fmt.Println("\nRegular years (52 ISO weeks):")
	for _, year := range regularYears {
		dt := chronogo.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
		if !dt.IsLongYear() {
			fmt.Printf("  ✓ %d is a regular year\n", year)
		}
	}

	fmt.Println()

	// On() - Quick date setting
	fmt.Println("2. On() - Convenience Method for Setting Dates")
	fmt.Println("----------------------------------------------")

	base := chronogo.Now()
	fmt.Printf("Current datetime: %s\n", base.Format("2006-01-02 15:04:05"))

	// Set date to New Year's Day (keeping current time)
	newYear := base.On(2024, time.January, 1)
	fmt.Printf("On New Year's Day: %s (time preserved)\n", newYear.Format("2006-01-02 15:04:05"))

	// Set date to Valentine's Day
	valentine := base.On(2024, time.February, 14)
	fmt.Printf("On Valentine's Day: %s (time preserved)\n", valentine.Format("2006-01-02 15:04:05"))

	// Set date to Christmas
	christmas := base.On(2024, time.December, 25)
	fmt.Printf("On Christmas: %s (time preserved)\n", christmas.Format("2006-01-02 15:04:05"))

	fmt.Println()

	// At() - Quick time setting
	fmt.Println("3. At() - Convenience Method for Setting Times")
	fmt.Println("----------------------------------------------")

	today := chronogo.Today()
	fmt.Printf("Today at midnight: %s\n", today.Format("2006-01-02 15:04:05"))

	// Morning meeting
	morning := today.At(9, 30, 0)
	fmt.Printf("Morning meeting: %s (date preserved)\n", morning.Format("2006-01-02 15:04:05"))

	// Lunch time
	lunch := today.At(12, 0, 0)
	fmt.Printf("Lunch time: %s (date preserved)\n", lunch.Format("2006-01-02 15:04:05"))

	// End of work
	endOfDay := today.At(17, 30, 0)
	fmt.Printf("End of work: %s (date preserved)\n", endOfDay.Format("2006-01-02 15:04:05"))

	// Last second of the day
	lastSecond := today.At(23, 59, 59)
	fmt.Printf("Last second: %s (date preserved)\n", lastSecond.Format("2006-01-02 15:04:05"))

	fmt.Println()

	// Combining On() and At()
	fmt.Println("4. Chaining On() and At() Together")
	fmt.Println("----------------------------------")

	// Schedule a meeting on a specific date and time
	meeting := chronogo.Now().
		On(2024, time.June, 15).
		At(14, 30, 0)
	fmt.Printf("Meeting scheduled: %s\n", meeting.Format("Monday, January 2, 2006 at 3:04 PM"))

	// Create a deadline (end of year)
	deadline := chronogo.Now().
		On(2024, time.December, 31).
		At(23, 59, 59)
	fmt.Printf("Project deadline: %s\n", deadline.Format("Monday, January 2, 2006 at 3:04 PM"))

	// Start of fiscal year Q2
	fiscalQ2 := chronogo.Now().
		On(2024, time.April, 1).
		At(0, 0, 0)
	fmt.Printf("Fiscal Q2 starts: %s\n", fiscalQ2.Format("Monday, January 2, 2006 at 3:04 PM"))

	fmt.Println()

	// Practical use cases
	fmt.Println("5. Practical Use Cases")
	fmt.Println("---------------------")

	// Schedule recurring weekly meetings for the next 4 weeks
	fmt.Println("Weekly team meetings (Mondays at 10:00 AM):")
	currentWeek := chronogo.Now().NextWeekday(time.Monday).At(10, 0, 0)
	for i := 0; i < 4; i++ {
		meeting := currentWeek.AddDays(7 * i)
		fmt.Printf("  Week %d: %s\n", i+1, meeting.Format("Mon Jan 2, 2006 at 3:04 PM"))
	}

	fmt.Println()

	// Set daily reminders
	fmt.Println("Daily reminders:")
	reminder1 := chronogo.Today().At(8, 0, 0)  // Morning standup
	reminder2 := chronogo.Today().At(14, 0, 0) // Afternoon check-in
	reminder3 := chronogo.Today().At(18, 0, 0) // End of day report

	fmt.Printf("  Morning standup: %s\n", reminder1.Format("3:04 PM"))
	fmt.Printf("  Afternoon check-in: %s\n", reminder2.Format("3:04 PM"))
	fmt.Printf("  End of day report: %s\n", reminder3.Format("3:04 PM"))

	fmt.Println()

	// Performance comparison
	fmt.Println("6. Performance Benefits")
	fmt.Println("----------------------")
	fmt.Println("On() is ~4.3x faster than Set().Year().Month().Day().Build()")
	fmt.Println("At() is ~3.7x faster than Set().Hour().Minute().Second().Build()")
	fmt.Println()
	fmt.Println("Use On() and At() for cleaner, faster code!")
	fmt.Println()

	// ISO long year practical example
	fmt.Println("7. ISO Long Year Practical Example")
	fmt.Println("----------------------------------")
	currentYear := time.Now().Year()
	dt := chronogo.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC)

	if dt.IsLongYear() {
		fmt.Printf("FYI: %d is a long year with 53 ISO weeks!\n", currentYear)
		fmt.Println("This affects week-based scheduling and reporting systems.")
	} else {
		fmt.Printf("%d is a regular year with 52 ISO weeks.\n", currentYear)
	}

	// Find next long year
	fmt.Println("\nNext long years:")
	count := 0
	for year := currentYear; count < 5; year++ {
		dt := chronogo.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		if dt.IsLongYear() {
			fmt.Printf("  - %d\n", year)
			count++
		}
	}
}
