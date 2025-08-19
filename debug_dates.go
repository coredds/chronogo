package main

import (
	"fmt"
	"time"
	
	chronogo "github.com/coredds/ChronoGo"
)

func main() {
	// Check MLK Day 2024 (third Monday in January)
	firstOfJan := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	fmt.Printf("Jan 1, 2024: %s\n", firstOfJan.Weekday())
	
	// Find first Monday
	daysToMonday := (7 - int(firstOfJan.Weekday()) + int(time.Monday)) % 7
	firstMonday := firstOfJan.AddDate(0, 0, daysToMonday)
	thirdMonday := firstMonday.AddDate(0, 0, 14) // Add 2 weeks
	
	fmt.Printf("MLK Day 2024: %s (%s)\n", thirdMonday.Format("2006-01-02"), thirdMonday.Weekday())
	
	// Check Jan 16, 2024
	jan16 := time.Date(2024, time.January, 16, 0, 0, 0, 0, time.UTC)
	fmt.Printf("Jan 16, 2024: %s (%s)\n", jan16.Format("2006-01-02"), jan16.Weekday())
	
	// Count business days in January 2024
	checker := chronogo.NewUSHolidayChecker()
	jan := chronogo.Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC)
	businessDays := jan.BusinessDaysInMonth(checker)
	fmt.Printf("Business days in January 2024: %d\n", businessDays)
	
	// Count manually
	count := 0
	start := chronogo.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := chronogo.Date(2024, time.January, 31, 0, 0, 0, 0, time.UTC)
	current := start
	
	fmt.Println("All days in January 2024:")
	for !current.After(end) {
		isHoliday := checker.IsHoliday(current)
		isBizDay := current.IsBusinessDay(checker)
		if isBizDay {
			count++
		}
		fmt.Printf("  %s (%s): Weekend=%t, Holiday=%t, Business=%t\n", 
			current.Format("01-02"), current.Weekday(), current.IsWeekend(), isHoliday, isBizDay)
		current = current.AddDays(1)
	}
	
	fmt.Printf("Manual count: %d\n", count)
}
