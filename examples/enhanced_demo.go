package main

import (
	"fmt"
	"strings"
	"time"

	chronogo "github.com/coredds/ChronoGo"
)

func main() {
	fmt.Println("ChronoGo v0.6.5 - Enhanced Business Operations Demo")
	fmt.Println("GoHoliday v0.6.3+ - 34 Countries Support")
	fmt.Println(strings.Repeat("=", 60))

	// Enhanced Business Day Calculator
	fmt.Println("\nEnhanced Business Day Calculator:")
	calc := chronogo.NewEnhancedBusinessDayCalculator("US")

	today := chronogo.Today()
	nextBiz := calc.NextBusinessDay(today)
	bizDays := calc.BusinessDaysBetween(today, nextBiz.AddDays(30))

	fmt.Printf("   Today: %s\n", today.Format("2006-01-02"))
	fmt.Printf("   Next business day: %s\n", nextBiz.Format("2006-01-02"))
	fmt.Printf("   Business days in next 30 days: %d\n", bizDays)

	// Multi-Country Holiday Checking (Enhanced in GoHoliday v0.6.3)
	fmt.Println("\nMulti-Country Holiday Checking (33 Countries):")
	countries := []string{"US", "BR", "IN", "KR", "IT", "ES", "NL", "PT", "PL", "CN", "TH", "SG", "ZA", "EG"}
	newYear := chronogo.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)

	for _, country := range countries {
		checker := chronogo.NewGoHolidayChecker(country)
		if checker.IsHoliday(newYear) {
			name := checker.GetHolidayName(newYear)
			fmt.Printf("   %s: %s\n", country, name)
		}
	}

	// Custom weekends for international business
	fmt.Println("\nInternational Business (Custom Weekends):")
	calc.SetCustomWeekends([]time.Weekday{time.Friday, time.Saturday})
	fmt.Printf("   With Fri-Sat weekends, Sunday is business day: %t\n",
		calc.IsBusinessDay(today.AddDays(1))) // Assuming today is Saturday

	// Holiday-Aware Scheduler
	fmt.Println("\nHoliday-Aware Scheduler:")
	scheduler := chronogo.NewHolidayAwareScheduler("US")

	// Schedule weekly team meetings
	weeklyMeetings := scheduler.ScheduleRecurring(
		chronogo.Date(2025, time.September, 1, 9, 0, 0, 0, time.UTC),
		7*24*time.Hour, 4)

	fmt.Printf("   Next 4 weekly meetings (avoiding holidays):\n")
	for i, meeting := range weeklyMeetings {
		fmt.Printf("     %d. %s (%s)\n", i+1,
			meeting.Format("2006-01-02"), meeting.Weekday())
	}

	// End-of-month reports
	monthlyReports := scheduler.ScheduleMonthlyEndOfMonth(
		chronogo.Date(2025, time.September, 1, 17, 0, 0, 0, time.UTC), 3)

	fmt.Printf("   Next 3 month-end reports (business days only):\n")
	for i, report := range monthlyReports {
		fmt.Printf("     %d. %s\n", i+1, report.Format("2006-01-02"))
	}

	// Holiday Calendar
	fmt.Println("\nHoliday Calendar:")
	calendar := chronogo.NewHolidayCalendar("US")

	// Upcoming holidays
	upcoming := calendar.GetUpcomingHolidays(chronogo.Now(), 3)
	fmt.Printf("   Next 3 holidays:\n")
	for i, holiday := range upcoming {
		fmt.Printf("     %d. %s\n", i+1, holiday.String())
	}

	// December holidays
	decHolidays := calendar.GetMonthlyHolidays(2025, time.December)
	fmt.Printf("   December 2025 holidays: %d found\n", len(decHolidays))
	for _, holiday := range decHolidays {
		fmt.Printf("     - %s\n", holiday.String())
	}

	// Convenience methods
	fmt.Println("\nConvenient DateTime Methods:")
	dt := chronogo.Date(2025, time.January, 15, 0, 0, 0, 0, time.UTC)

	// Get enhanced calculator
	enhancedCalc := dt.WithEnhancedBusinessDays("US")
	fmt.Printf("   Enhanced calc for %s: %t\n",
		dt.Format("2006-01-02"), enhancedCalc.IsBusinessDay(dt))

	// Get upcoming holidays
	upcomingFromDate := dt.GetUpcomingHolidays("US", 2)
	fmt.Printf("   Upcoming holidays from %s: %d found\n",
		dt.Format("2006-01-02"), len(upcomingFromDate))

	fmt.Println("\nAll features work seamlessly with existing ChronoGo functionality!")
}
