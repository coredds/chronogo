package main

import (
	"fmt"
	"time"

	chronogo "github.com/coredds/chronogo"
)

func main() {
	fmt.Println("=== Testing Helpers Demo ===")
	testingHelpersDemo()

	fmt.Println("\n=== Weekday Navigation Demo ===")
	weekdayNavigationDemo()

	fmt.Println("\n=== Nth Weekday Occurrence Demo ===")
	nthWeekdayDemo()
}

func testingHelpersDemo() {
	// 1. SetTestNow - Mock current time for deterministic tests
	fmt.Println("1. SetTestNow - Mock current time")
	chronogo.SetTestNow(chronogo.Date(2024, 12, 25, 10, 0, 0, 0, time.UTC))
	fmt.Printf("   Now(): %v\n", chronogo.Now().Format("2006-01-02 15:04:05"))
	chronogo.ClearTestNow()
	fmt.Printf("   After ClearTestNow(): %v\n", chronogo.Now().Format("2006-01-02 15:04:05"))

	// 2. FreezeTime - Stop time at current moment
	fmt.Println("\n2. FreezeTime - Stop time completely")
	chronogo.FreezeTime()
	frozenTime := chronogo.Now()
	fmt.Printf("   Frozen at: %v\n", frozenTime.Format("2006-01-02 15:04:05"))
	time.Sleep(10 * time.Millisecond)
	fmt.Printf("   After sleep: %v (same)\n", chronogo.Now().Format("2006-01-02 15:04:05"))
	chronogo.UnfreezeTime()

	// 3. TravelTo - Jump to a specific point in time
	fmt.Println("\n3. TravelTo - Time travel to specific date")
	chronogo.TravelTo(chronogo.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	fmt.Printf("   Traveled to: %v\n", chronogo.Now().Format("2006-01-02"))
	chronogo.ClearTestNow()

	// 4. WithTestNow - Scoped time mocking
	fmt.Println("\n4. WithTestNow - Scoped time mocking")
	chronogo.WithTestNow(chronogo.Date(2024, 7, 4, 12, 0, 0, 0, time.UTC), func() {
		fmt.Printf("   Inside scope: %v\n", chronogo.Now().Format("2006-01-02"))
	})
	fmt.Printf("   Outside scope: real time restored\n")
}

func weekdayNavigationDemo() {
	// Start date: Monday, January 15, 2024
	dt := chronogo.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	fmt.Printf("Starting date: %v (%v)\n", dt.Format("2006-01-02"), dt.Weekday())

	// 1. NextWeekday
	nextWed := dt.NextWeekday(time.Wednesday)
	fmt.Printf("\n1. NextWeekday(Wednesday): %v (%v)", nextWed.Format("2006-01-02"), nextWed.Weekday())

	// 2. PreviousWeekday
	prevFri := dt.PreviousWeekday(time.Friday)
	fmt.Printf("\n2. PreviousWeekday(Friday): %v (%v)", prevFri.Format("2006-01-02"), prevFri.Weekday())

	// 3. ClosestWeekday
	closest := dt.ClosestWeekday(time.Thursday)
	fmt.Printf("\n3. ClosestWeekday(Thursday): %v (%v)", closest.Format("2006-01-02"), closest.Weekday())

	// 4. FarthestWeekday
	farthest := dt.FarthestWeekday(time.Thursday)
	fmt.Printf("\n4. FarthestWeekday(Thursday): %v (%v)", farthest.Format("2006-01-02"), farthest.Weekday())

	// 5. NextOrSameWeekday
	sameDay := dt.NextOrSameWeekday(time.Monday)
	fmt.Printf("\n5. NextOrSameWeekday(Monday): %v (same day)", sameDay.Format("2006-01-02"))

	nextDay := dt.NextOrSameWeekday(time.Tuesday)
	fmt.Printf("\n6. NextOrSameWeekday(Tuesday): %v", nextDay.Format("2006-01-02"))
}

func nthWeekdayDemo() {
	// March 2024
	dt := chronogo.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	fmt.Printf("Working with: March 2024\n")

	// 1. FirstWeekdayOf
	firstMonday := dt.FirstWeekdayOf(time.Monday)
	fmt.Printf("\n1. FirstWeekdayOf(Monday): %v", firstMonday.Format("2006-01-02"))

	// 2. LastWeekdayOf
	lastFriday := dt.LastWeekdayOf(time.Friday)
	fmt.Printf("\n2. LastWeekdayOf(Friday): %v", lastFriday.Format("2006-01-02"))

	// 3. NthWeekdayOfMonth
	secondTuesday := dt.NthWeekdayOfMonth(2, time.Tuesday)
	fmt.Printf("\n3. NthWeekdayOfMonth(2, Tuesday): %v", secondTuesday.Format("2006-01-02"))

	thirdWednesday := dt.NthWeekdayOfMonth(3, time.Wednesday)
	fmt.Printf("\n4. NthWeekdayOfMonth(3, Wednesday): %v", thirdWednesday.Format("2006-01-02"))

	// 5. NthWeekdayOfYear
	yearStart := chronogo.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	tenthMonday := yearStart.NthWeekdayOfYear(10, time.Monday)
	fmt.Printf("\n5. NthWeekdayOfYear(10, Monday) in 2024: %v", tenthMonday.Format("2006-01-02"))

	// 6. WeekdayOccurrenceInMonth
	specificDate := chronogo.Date(2024, 3, 11, 0, 0, 0, 0, time.UTC)
	occurrence := specificDate.WeekdayOccurrenceInMonth()
	fmt.Printf("\n6. %v is the %d%s %v of the month",
		specificDate.Format("2006-01-02"),
		occurrence,
		ordinalSuffix(occurrence),
		specificDate.Weekday())

	// 7. IsNthWeekdayOf
	isSecond := specificDate.IsNthWeekdayOf(2, "month")
	fmt.Printf("\n7. Is March 11, 2024 the 2nd Monday? %v", isSecond)

	// 8. Practical example: Thanksgiving (4th Thursday of November)
	thanksgiving := chronogo.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC).NthWeekdayOfMonth(4, time.Thursday)
	fmt.Printf("\n\n📅 Practical Example - Thanksgiving 2024: %v", thanksgiving.Format("Monday, January 2, 2006"))
}

func ordinalSuffix(n int) string {
	if n >= 11 && n <= 13 {
		return "th"
	}
	switch n % 10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}
