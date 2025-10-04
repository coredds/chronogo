// Demo program showcasing chronogo's comprehensive datetime functionality
package main

import (
	"fmt"
	"time"

	"github.com/coredds/chronogo"
)

func main() {
	printHeader("ChronoGo - Comprehensive DateTime Library Demo")

	// Core DateTime Operations
	demoCreationAndBasics()
	demoParsingAndFormatting()
	demoArithmeticOperations()
	demoConvenienceMethods()
	
	// Advanced Features
	demoDiffType()
	demoTimezoneOperations()
	demoPeriodOperations()
	demoFluentAPI()
	
	// Specialized Features
	demoWeekdayOperations()
	demoComparisonMethods()
	demoBusinessDateOperations()
	demoLocalization()
	demoTestingHelpers()

	printFooter()
}

func printHeader(title string) {
	fmt.Println()
	fmt.Println(title)
	fmt.Println(repeat("=", len(title)))
	fmt.Println()
	}

func printSection(title string) {
	fmt.Printf("\n%s\n%s\n", title, repeat("-", len(title)))
	}

func printFooter() {
	fmt.Println()
	fmt.Println("Demo Complete")
	fmt.Println("For more information: https://github.com/coredds/chronogo")
	fmt.Println()
}

func repeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

func demoCreationAndBasics() {
	printSection("1. DateTime Creation and Basics")

	now := chronogo.Now()
	fmt.Printf("Current time:        %s\n", now.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("Current time (UTC):  %s\n", chronogo.NowUTC().Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("Today:               %s\n", chronogo.Today().ToDateString())
	fmt.Printf("Tomorrow:            %s\n", chronogo.Tomorrow().ToDateString())
	fmt.Printf("Yesterday:           %s\n", chronogo.Yesterday().ToDateString())

	custom := chronogo.Date(2024, time.June, 15, 14, 30, 0, 0, time.UTC)
	fmt.Printf("Custom date:         %s\n", custom.Format("2006-01-02 15:04:05"))

	fromUnix := chronogo.FromUnix(1718461800, 0, time.UTC)
	fmt.Printf("From Unix timestamp: %s\n", fromUnix.Format("2006-01-02 15:04:05"))
}

func demoParsingAndFormatting() {
	printSection("2. Parsing and Formatting")

	// Natural language parsing (multi-language support)
	examples := []string{"tomorrow", "next Monday", "3 days ago", "in 2 weeks"}
	fmt.Println("Natural language parsing:")
	for _, example := range examples {
		if dt, err := chronogo.Parse(example); err == nil {
			fmt.Printf("  '%s' -> %s\n", example, dt.Format("2006-01-02"))
		}
	}

	// Technical format parsing
	iso, _ := chronogo.Parse("2024-06-15T14:30:00Z")
	fmt.Printf("\nISO 8601:            %s\n", iso.Format("2006-01-02 15:04:05"))

	custom, _ := chronogo.FromFormat("15/06/2024 14:30", "02/01/2006 15:04")
	fmt.Printf("Custom format:       %s\n", custom.Format("2006-01-02 15:04:05"))

	// Various output formats
	dt := chronogo.Date(2024, time.June, 15, 14, 30, 0, 0, time.UTC)
	fmt.Printf("\nFormatting options:\n")
	fmt.Printf("  ISO 8601:          %s\n", dt.ToISO8601String())
	fmt.Printf("  RFC 3339:          %s\n", dt.Format(time.RFC3339))
	fmt.Printf("  Date only:         %s\n", dt.ToDateString())
	fmt.Printf("  Time only:         %s\n", dt.ToTimeString())
	fmt.Printf("  Cookie format:     %s\n", dt.ToCookieString())
	fmt.Printf("  Human readable:    %s\n", dt.DiffForHumans())
}

func demoArithmeticOperations() {
	printSection("3. Arithmetic Operations")

	base := chronogo.Date(2024, time.January, 15, 10, 0, 0, 0, time.UTC)
	fmt.Printf("Base date:           %s\n", base.Format("2006-01-02 15:04:05"))

	// Addition
	fmt.Printf("Add 1 year:          %s\n", base.AddYears(1).Format("2006-01-02"))
	fmt.Printf("Add 3 months:        %s\n", base.AddMonths(3).Format("2006-01-02"))
	fmt.Printf("Add 7 days:          %s\n", base.AddDays(7).Format("2006-01-02"))
	fmt.Printf("Add 5 hours:         %s\n", base.AddHours(5).Format("2006-01-02 15:04:05"))

	// Subtraction
	fmt.Printf("Subtract 2 weeks:    %s\n", base.SubtractDays(14).Format("2006-01-02"))
	fmt.Printf("Subtract 30 minutes: %s\n", base.SubtractMinutes(30).Format("2006-01-02 15:04:05"))

	// Chaining
	result := base.AddYears(1).AddMonths(2).AddDays(10).AddHours(5)
	fmt.Printf("Chained operations:  %s\n", result.Format("2006-01-02 15:04:05"))
}

func demoConvenienceMethods() {
	printSection("4. Convenience Methods")

	base := chronogo.Date(2024, time.March, 15, 14, 30, 45, 0, time.UTC)
	fmt.Printf("Base:                %s\n", base.Format("2006-01-02 15:04:05"))

	// On() - Set date in one call (4.3x faster than Set chain)
	newDate := base.On(2024, time.December, 25)
	fmt.Printf("On(2024, Dec, 25):   %s (time preserved)\n", newDate.Format("2006-01-02 15:04:05"))

	// At() - Set time in one call (3.7x faster than Set chain)
	newTime := base.At(9, 0, 0)
	fmt.Printf("At(9, 0, 0):         %s (date preserved)\n", newTime.Format("2006-01-02 15:04:05"))

	// Chaining On() and At()
	meeting := chronogo.Now().On(2024, time.June, 15).At(14, 30, 0)
	fmt.Printf("Meeting scheduled:   %s\n", meeting.Format("2006-01-02 15:04:05"))

	// Boundary methods
	fmt.Printf("\nBoundary operations:\n")
	fmt.Printf("  Start of day:      %s\n", base.StartOfDay().Format("2006-01-02 15:04:05"))
	fmt.Printf("  End of day:        %s\n", base.EndOfDay().Format("2006-01-02 15:04:05"))
	fmt.Printf("  Start of month:    %s\n", base.StartOfMonth().Format("2006-01-02"))
	fmt.Printf("  End of month:      %s\n", base.EndOfMonth().Format("2006-01-02"))
	fmt.Printf("  Start of quarter:  %s\n", base.StartOfQuarter().Format("2006-01-02"))
	fmt.Printf("  End of quarter:    %s\n", base.EndOfQuarter().Format("2006-01-02"))

	// ISO 8601 features
	fmt.Printf("\nISO 8601 features:\n")
	fmt.Printf("  Is long year:      %v\n", chronogo.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).IsLongYear())
	fmt.Printf("  ISO week:          %d\n", base.ISOWeekNumber())
	fmt.Printf("  Day of year:       %d\n", base.DayOfYear())
	fmt.Printf("  Quarter:           %d\n", base.Quarter())
}

func demoDiffType() {
	printSection("5. Diff Type - Rich DateTime Differences")

	start := chronogo.Date(2023, time.January, 15, 10, 0, 0, 0, time.UTC)
	end := chronogo.Date(2024, time.March, 20, 14, 30, 0, 0, time.UTC)
	diff := end.Diff(start)

	fmt.Printf("Start: %s\n", start.Format("2006-01-02 15:04:05"))
	fmt.Printf("End:   %s\n", end.Format("2006-01-02 15:04:05"))

	// Calendar-aware differences
	fmt.Printf("\nCalendar-aware (integer):\n")
	fmt.Printf("  Years:             %d\n", diff.Years())
	fmt.Printf("  Months:            %d\n", diff.Months())
	fmt.Printf("  Weeks:             %d\n", diff.Weeks())
	fmt.Printf("  Days:              %d\n", diff.Days())
	fmt.Printf("  Hours:             %d\n", diff.Hours())

	// Precise differences
	fmt.Printf("\nPrecise (float):\n")
	fmt.Printf("  In years:          %.2f\n", diff.InYears())
	fmt.Printf("  In months:         %.2f\n", diff.InMonths())
	fmt.Printf("  In weeks:          %.2f\n", diff.InWeeks())
	fmt.Printf("  In days:           %.2f\n", diff.InDays())
	fmt.Printf("  In hours:          %.2f\n", diff.InHours())

	// Human-readable formats
	fmt.Printf("\nHuman-readable:\n")
	fmt.Printf("  For humans:        %s\n", diff.ForHumans())
	fmt.Printf("  Comparison:        %s\n", diff.ForHumansComparison())
	fmt.Printf("  Detailed:          %s\n", diff.String())
	fmt.Printf("  Compact:           %s\n", diff.CompactString())

	// Absolute and inverted
	pastDiff := start.Diff(end)
	fmt.Printf("\nTransformations:\n")
	fmt.Printf("  Negative:          %s\n", pastDiff.CompactString())
	fmt.Printf("  Absolute:          %s\n", pastDiff.Abs().CompactString())
	fmt.Printf("  Inverted:          %s\n", diff.Invert().CompactString())
}

func demoTimezoneOperations() {
	printSection("6. Timezone Operations")

	utc := chronogo.Date(2024, time.June, 15, 14, 30, 0, 0, time.UTC)
	fmt.Printf("UTC time:            %s\n", utc.Format("2006-01-02 15:04:05 MST"))

	ny, _ := chronogo.LoadLocation("America/New_York")
	nyTime := utc.In(ny)
	fmt.Printf("New York:            %s (DST: %v)\n", nyTime.Format("2006-01-02 15:04:05 MST"), nyTime.IsDST())

	tokyo, _ := chronogo.LoadLocation("Asia/Tokyo")
	tokyoTime := utc.In(tokyo)
	fmt.Printf("Tokyo:               %s (DST: %v)\n", tokyoTime.Format("2006-01-02 15:04:05 MST"), tokyoTime.IsDST())

	london, _ := chronogo.LoadLocation("Europe/London")
	londonTime := utc.In(london)
	fmt.Printf("London:              %s (DST: %v)\n", londonTime.Format("2006-01-02 15:04:05 MST"), londonTime.IsDST())
}

func demoPeriodOperations() {
	printSection("7. Period Operations")

	start := chronogo.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := chronogo.Date(2024, time.January, 10, 0, 0, 0, 0, time.UTC)
	period := chronogo.NewPeriod(start, end)

	fmt.Printf("Period:              %s to %s\n", start.ToDateString(), end.ToDateString())
	fmt.Printf("Duration:            %s\n", period.String())
	fmt.Printf("Days:                %d\n", period.Days())
	fmt.Printf("Hours:               %d\n", period.Hours())

	checkDate := chronogo.Date(2024, time.January, 5, 0, 0, 0, 0, time.UTC)
	fmt.Printf("Contains %s:      %v\n", checkDate.ToDateString(), period.Contains(checkDate))

	// Period iteration
	fmt.Print("Dates in period:     ")
	count := 0
	for date := range period.RangeDays() {
		if count > 0 {
			fmt.Print(", ")
		}
		fmt.Print(date.Format("Jan 2"))
		count++
		if count >= 5 {
			fmt.Print("...")
			break
		}
	}
	fmt.Println()

	// Period operations
	p1 := chronogo.NewPeriod(
		chronogo.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		chronogo.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
	)
	p2 := chronogo.NewPeriod(
		chronogo.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		chronogo.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	)
	fmt.Printf("Periods overlap:     %v\n", p1.Overlaps(p2))
	fmt.Printf("Period encompasses:  %v\n", p1.Encompasses(p2))
}

func demoFluentAPI() {
	printSection("8. Fluent API")

	base := chronogo.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	fmt.Printf("Base:                %s\n", base.Format("2006-01-02 15:04:05"))

	// Fluent duration addition
	result := base.AddFluent().Years(1).Months(2).Days(10).Hours(5).Minutes(30).To(base)
	fmt.Printf("Fluent addition:     %s\n", result.Format("2006-01-02 15:04:05"))

	// Fluent setting
	configured := base.Set().Year(2025).Month(time.December).Day(25).Hour(15).Minute(30).Build()
	fmt.Printf("Fluent setting:      %s\n", configured.Format("2006-01-02 15:04:05"))

	// ChronoDuration
	duration := chronogo.NewDuration(25*time.Hour + 30*time.Minute)
	fmt.Printf("\nChronoDuration:      %s\n", duration.String())
	fmt.Printf("Human readable:      %s\n", duration.HumanString())
	fmt.Printf("Days:                %.2f\n", duration.Days())
	fmt.Printf("Weeks:               %.2f\n", duration.Weeks())
}

func demoWeekdayOperations() {
	printSection("9. Weekday Operations")

	base := chronogo.Date(2024, time.June, 15, 0, 0, 0, 0, time.UTC)
	fmt.Printf("Base (%s):        %s\n", base.Weekday(), base.Format("2006-01-02"))

	fmt.Printf("Next Monday:         %s\n", base.NextWeekday(time.Monday).Format("2006-01-02"))
	fmt.Printf("Previous Friday:     %s\n", base.PreviousWeekday(time.Friday).Format("2006-01-02"))
	fmt.Printf("Closest Wednesday:   %s\n", base.ClosestWeekday(time.Wednesday).Format("2006-01-02"))

	// Nth weekday operations
	fmt.Printf("\nNth weekday of month:\n")
	fmt.Printf("  1st Monday:        %s\n", base.FirstWeekdayOf(time.Monday).Format("2006-01-02"))
	fmt.Printf("  Last Friday:       %s\n", base.LastWeekdayOf(time.Friday).Format("2006-01-02"))
	fmt.Printf("  2nd Tuesday:       %s\n", base.NthWeekdayOfMonth(2, time.Tuesday).Format("2006-01-02"))

	fmt.Printf("\nWeekday checks:\n")
	fmt.Printf("  Is weekend:        %v\n", base.IsWeekend())
	fmt.Printf("  Is weekday:        %v\n", base.IsWeekday())
}

func demoComparisonMethods() {
	printSection("10. Comparison Methods")

	dt1 := chronogo.Date(2024, time.June, 15, 0, 0, 0, 0, time.UTC)
	dt2 := chronogo.Date(2024, time.June, 20, 0, 0, 0, 0, time.UTC)

	fmt.Printf("Date 1:              %s\n", dt1.ToDateString())
	fmt.Printf("Date 2:              %s\n", dt2.ToDateString())

	fmt.Printf("\nBasic comparisons:\n")
	fmt.Printf("  Before:            %v\n", dt1.Before(dt2))
	fmt.Printf("  After:             %v\n", dt1.After(dt2))
	fmt.Printf("  Equal:             %v\n", dt1.Equal(dt2))
	fmt.Printf("  Between:           %v\n", dt1.Between(dt1, dt2, true))

	fmt.Printf("\nCalendar comparisons:\n")
	fmt.Printf("  Same day:          %v\n", dt1.IsSameDay(dt2))
	fmt.Printf("  Same month:        %v\n", dt1.IsSameMonth(dt2))
	fmt.Printf("  Same year:         %v\n", dt1.IsSameYear(dt2))
	fmt.Printf("  Same quarter:      %v\n", dt1.IsSameQuarter(dt2))

	birthday := chronogo.Date(1990, time.June, 15, 0, 0, 0, 0, time.UTC)
	fmt.Printf("\nSpecial comparisons:\n")
	fmt.Printf("  Is birthday:       %v\n", dt1.IsBirthday(birthday))
	fmt.Printf("  Is anniversary:    %v\n", dt1.IsAnniversary(birthday))
	fmt.Printf("  Is leap year:      %v\n", dt1.IsLeapYear())
	fmt.Printf("  Is past:           %v\n", dt1.IsPast())
	fmt.Printf("  Is future:         %v\n", dt1.IsFuture())

	// Closest/Farthest
	dates := []chronogo.DateTime{
		chronogo.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC),
		chronogo.Date(2024, 6, 25, 0, 0, 0, 0, time.UTC),
		chronogo.Date(2024, 6, 18, 0, 0, 0, 0, time.UTC),
	}
	closest := dt1.Closest(dates...)
	farthest := dt1.Farthest(dates...)
	fmt.Printf("  Closest date:      %s\n", closest.ToDateString())
	fmt.Printf("  Farthest date:     %s\n", farthest.ToDateString())

	// Average
	average := dt1.Average(dt2)
	fmt.Printf("  Average:           %s\n", average.ToDateString())
}

func demoBusinessDateOperations() {
	printSection("11. Business Date Operations")

	checker := chronogo.NewGoHolidayChecker("US")

	testDate := chronogo.Date(2024, time.July, 4, 0, 0, 0, 0, time.UTC)
	fmt.Printf("Date:                %s (%s)\n", testDate.Format("2006-01-02"), testDate.Weekday())
	fmt.Printf("Is holiday:          %v\n", checker.IsHoliday(testDate))
	fmt.Printf("Holiday name:        %s\n", checker.GetHolidayName(testDate))
	fmt.Printf("Is business day:     %v\n", testDate.IsBusinessDay(checker))

	start := chronogo.Date(2024, time.July, 1, 0, 0, 0, 0, time.UTC)
	fmt.Printf("\nBusiness day operations from %s:\n", start.Format("2006-01-02"))
	fmt.Printf("  Next business day: %s\n", start.NextBusinessDay(checker).Format("2006-01-02"))
	fmt.Printf("  Add 5 biz days:    %s\n", start.AddBusinessDays(5, checker).Format("2006-01-02"))

	end := chronogo.Date(2024, time.July, 31, 0, 0, 0, 0, time.UTC)
	fmt.Printf("  Biz days in range: %d\n", start.BusinessDaysBetween(end, checker))
	fmt.Printf("  Biz days in month: %d\n", start.BusinessDaysInMonth(checker))
	fmt.Printf("  Biz days in year:  %d\n", start.BusinessDaysInYear(checker))

	// Multi-country support
	fmt.Printf("\nMulti-country support (34 countries):\n")
	countries := []string{"US", "GB", "FR", "DE", "JP", "BR"}
	newYear := chronogo.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	for _, country := range countries {
		checker := chronogo.NewGoHolidayChecker(country)
		fmt.Printf("  %s: %v\n", country, checker.IsHoliday(newYear))
	}
}

func demoLocalization() {
	printSection("12. Localization")

	dt := chronogo.Date(2024, time.June, 15, 14, 30, 0, 0, time.UTC)

	// Supported locales
	locales := []string{"en-US", "es-ES", "pt-BR", "fr-FR", "de-DE", "ja-JP"}
	fmt.Println("Human-readable differences in multiple languages:")
	for _, locale := range locales {
		humanStr, _ := dt.HumanStringLocalized(locale)
		fmt.Printf("  %s: %s\n", locale, humanStr)
	}

	// Month and weekday names
	fmt.Println("\nLocalized month names:")
	for _, locale := range []string{"en-US", "es-ES", "fr-FR", "de-DE"} {
		monthName, _ := dt.GetMonthName(locale)
		fmt.Printf("  %s: %s\n", locale, monthName)
	}

	fmt.Println("\nLocalized weekday names:")
	for _, locale := range []string{"en-US", "es-ES", "fr-FR", "de-DE"} {
		weekdayName, _ := dt.GetWeekdayName(locale)
		fmt.Printf("  %s: %s\n", locale, weekdayName)
	}
}

func demoTestingHelpers() {
	printSection("13. Testing Helpers")

	fmt.Println("Time manipulation for testing:")

	// Save current state
	original := chronogo.Now()
	fmt.Printf("Original time:       %s\n", original.Format("2006-01-02 15:04:05"))

	// Set test time
	testTime := chronogo.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)
	chronogo.SetTestNow(testTime)
	fmt.Printf("After SetTestNow:    %s\n", chronogo.Now().Format("2006-01-02 15:04:05"))

	// Freeze time
	chronogo.FreezeTime()
	fmt.Printf("Frozen time:         %s\n", chronogo.Now().Format("2006-01-02 15:04:05"))
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Still frozen:        %s\n", chronogo.Now().Format("2006-01-02 15:04:05"))

	// Travel in time
	chronogo.UnfreezeTime()
	chronogo.TravelForward(24 * time.Hour)
	fmt.Printf("After 24h forward:   %s\n", chronogo.Now().Format("2006-01-02 15:04:05"))

	chronogo.TravelBack(48 * time.Hour)
	fmt.Printf("After 48h back:      %s\n", chronogo.Now().Format("2006-01-02 15:04:05"))

	// Clean up
	chronogo.ClearTestNow()
	fmt.Printf("Back to real time:   %s\n", chronogo.Now().Format("2006-01-02 15:04:05"))

	fmt.Println("\nTest mode status:")
	chronogo.SetTestNow(testTime)
	fmt.Printf("  Is test mode:      %v\n", chronogo.IsTestMode())
	fmt.Printf("  Test time:         %s\n", chronogo.GetTestNow().Format("2006-01-02 15:04:05"))
	chronogo.ClearTestNow()
}