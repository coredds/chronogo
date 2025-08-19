// Demo program showing ChronoGo functionality
package main

import (
	"errors"
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

	// Business Date Operations (NEW!)
	fmt.Println("12. Business Date Operations:")

	// Set up holiday checker
	holidayChecker := chronogo.NewUSHolidayChecker()

	// Add custom company holiday
	companyHoliday := chronogo.Holiday{
		Name:  "Company Founding Day",
		Month: time.March,
		Day:   15,
	}
	holidayChecker.AddHoliday(companyHoliday)

	// Test various dates
	testDates := []chronogo.DateTime{
		chronogo.Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC), // MLK Day (Monday)
		chronogo.Date(2024, time.January, 16, 0, 0, 0, 0, time.UTC), // Tuesday
		chronogo.Date(2024, time.January, 13, 0, 0, 0, 0, time.UTC), // Saturday
		chronogo.Date(2024, time.March, 15, 0, 0, 0, 0, time.UTC),   // Custom holiday
	}

	for _, date := range testDates {
		fmt.Printf("   %s (%s): Business day? %t\n",
			date.Format("2006-01-02"), date.Weekday(), date.IsBusinessDay(holidayChecker))
	}

	// Business day arithmetic
	startBiz := chronogo.Date(2024, time.January, 8, 0, 0, 0, 0, time.UTC) // Monday
	fmt.Printf("   Start: %s (%s)\n", startBiz.Format("2006-01-02"), startBiz.Weekday())

	nextBiz := startBiz.NextBusinessDay(holidayChecker)
	fmt.Printf("   Next business day: %s (%s)\n", nextBiz.Format("2006-01-02"), nextBiz.Weekday())

	add5Biz := startBiz.AddBusinessDays(5, holidayChecker)
	fmt.Printf("   Add 5 business days: %s (%s)\n", add5Biz.Format("2006-01-02"), add5Biz.Weekday())

	endBiz := chronogo.Date(2024, time.January, 22, 0, 0, 0, 0, time.UTC)
	bizCount := startBiz.BusinessDaysBetween(endBiz, holidayChecker)
	fmt.Printf("   Business days between %s and %s: %d\n",
		startBiz.Format("01-02"), endBiz.Format("01-02"), bizCount)

	// Month/year business day counts
	janDate := chronogo.Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC)
	fmt.Printf("   Business days in January 2024: %d\n", janDate.BusinessDaysInMonth(holidayChecker))
	fmt.Printf("   Business days in 2024: %d\n", janDate.BusinessDaysInYear(holidayChecker))

	// List holidays for the year
	fmt.Printf("   US Holidays in 2024:\n")
	holidays2024 := holidayChecker.GetHolidays(2024)
	for _, holiday := range holidays2024[:5] { // Show first 5
		fmt.Printf("     %s (%s)\n", holiday.Format("2006-01-02"), holiday.Weekday())
	}
	fmt.Printf("     ... and %d more\n", len(holidays2024)-5)
	fmt.Println()

	// Enhanced Error Handling (NEW!)
	fmt.Println("13. Enhanced Error Handling:")

	// Parse error with helpful suggestion
	_, err := chronogo.Parse("25/12/2023") // Wrong format
	if err != nil {
		var chronoErr *chronogo.ChronoError
		if errors.As(err, &chronoErr) {
			fmt.Printf("   Parse error with suggestion:\n   %s\n", chronoErr.Error())
		}
	}

	// Timezone error with suggestion
	_, err = chronogo.LoadLocation("EST") // Ambiguous
	if err != nil {
		var chronoErr *chronogo.ChronoError
		if errors.As(err, &chronoErr) {
			fmt.Printf("   Timezone error with suggestion:\n   %s\n", chronoErr.Error())
		}
	}

	// Validation error
	zeroDate := chronogo.DateTime{}
	if err := zeroDate.Validate(); err != nil {
		fmt.Printf("   Validation error:\n   %s\n", err.Error())
	}

	// Show available timezones (sample)
	fmt.Printf("   Sample available timezones:\n")
	timezones := chronogo.AvailableTimezones()
	for _, tz := range timezones[:5] {
		fmt.Printf("     %s\n", tz)
	}
	fmt.Printf("     ... and %d more\n", len(timezones)-5)
	fmt.Println()

	// Must functions for constants
	fmt.Println("14. Must Functions (for constants):")
	// These would panic if the input was invalid, but are safe for known-good values
	appLaunchDate := chronogo.MustParse("2024-01-01T00:00:00Z")
	fmt.Printf("   App launch date: %s\n", appLaunchDate.Format("2006-01-02"))

	eastCoast := chronogo.MustLoadLocation("America/New_York")
	fmt.Printf("   East coast timezone: %s\n", eastCoast.String())
	fmt.Println()

	fmt.Println("=== Demo Complete ===")
}
