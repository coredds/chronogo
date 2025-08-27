package chronogo

import (
	"testing"
	"time"
)

func TestEnhancedBusinessDayCalculator(t *testing.T) {
	calc := NewEnhancedBusinessDayCalculator("US")
	
	// Test with a known business day
	businessDay := Date(2024, time.January, 2, 0, 0, 0, 0, time.UTC) // Tuesday
	if !calc.IsBusinessDay(businessDay) {
		t.Errorf("Expected %v to be a business day", businessDay)
	}
	
	// Test with weekend
	weekend := Date(2024, time.January, 6, 0, 0, 0, 0, time.UTC) // Saturday
	if calc.IsBusinessDay(weekend) {
		t.Errorf("Expected %v to not be a business day", weekend)
	}
	
	// Test adding business days
	start := Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC) // Monday, New Year's Day
	result := calc.AddBusinessDays(start, 1)
	expected := Date(2024, time.January, 2, 0, 0, 0, 0, time.UTC) // Tuesday
	if !result.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestEnhancedBusinessDayCalculatorCustomWeekends(t *testing.T) {
	calc := NewEnhancedBusinessDayCalculator("US")
	
	// Set custom weekends (Friday-Saturday for some Middle Eastern countries)
	calc.SetCustomWeekends([]time.Weekday{time.Friday, time.Saturday})
	
	friday := Date(2024, time.January, 5, 0, 0, 0, 0, time.UTC)
	if calc.IsBusinessDay(friday) {
		t.Errorf("Expected Friday to not be a business day with custom weekends")
	}
	
	sunday := Date(2024, time.January, 7, 0, 0, 0, 0, time.UTC)
	if !calc.IsBusinessDay(sunday) {
		t.Errorf("Expected Sunday to be a business day with custom weekends")
	}
}

func TestHolidayAwareScheduler(t *testing.T) {
	scheduler := NewHolidayAwareScheduler("US")
	
	// Test recurring schedule
	start := Date(2024, time.January, 1, 9, 0, 0, 0, time.UTC)
	daily := scheduler.ScheduleRecurring(start, 24*time.Hour, 5)
	
	if len(daily) != 5 {
		t.Errorf("Expected 5 scheduled dates, got %d", len(daily))
	}
	
	// Verify no weekends or holidays in business day schedule
	businessDays := scheduler.ScheduleBusinessDays(start, 5)
	for _, day := range businessDays {
		if day.IsWeekend() {
			t.Errorf("Business day schedule should not include weekends: %v", day)
		}
	}
}

func TestHolidayAwareSchedulerMonthlyEndOfMonth(t *testing.T) {
	scheduler := NewHolidayAwareScheduler("US")
	
	start := Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	monthly := scheduler.ScheduleMonthlyEndOfMonth(start, 3)
	
	if len(monthly) != 3 {
		t.Errorf("Expected 3 monthly end dates, got %d", len(monthly))
	}
	
	// Each date should be near end of month
	for i, date := range monthly {
		expectedMonth := time.Month(int(time.January) + i)
		if date.Month() != expectedMonth {
			t.Errorf("Expected month %v, got %v for date %v", expectedMonth, date.Month(), date)
		}
		
		// Should be a business day (end of month adjusted for holidays/weekends)
		calc := NewEnhancedBusinessDayCalculator("US")
		if !calc.IsBusinessDay(date) {
			t.Errorf("End of month date should be a business day: %v", date)
		}
	}
}

func TestHolidayCalendar(t *testing.T) {
	calendar := NewHolidayCalendar("US")
	
	// Test generating a month
	entries := calendar.GenerateMonth(2024, time.January)
	if len(entries) != 31 {
		t.Errorf("Expected 31 entries for January, got %d", len(entries))
	}
	
	// January 1st should be a holiday (New Year's Day)
	newYears := entries[0] // January 1st
	if !newYears.IsHoliday {
		t.Errorf("Expected January 1st to be marked as holiday")
	}
	if newYears.HolidayName == "" {
		t.Errorf("Expected holiday name to be set for January 1st")
	}
}

func TestHolidayCalendarGetMonthlyHolidays(t *testing.T) {
	calendar := NewHolidayCalendar("US")
	
	// Get holidays for January 2024
	holidays := calendar.GetMonthlyHolidays(2024, time.January)
	
	// Should have at least New Year's Day and MLK Day
	if len(holidays) < 1 {
		t.Errorf("Expected at least 1 holiday in January 2024, got %d", len(holidays))
	}
	
	// All returned entries should be holidays
	for _, holiday := range holidays {
		if !holiday.IsHoliday {
			t.Errorf("Non-holiday returned in holiday list: %v", holiday)
		}
		if holiday.HolidayName == "" {
			t.Errorf("Holiday missing name: %v", holiday)
		}
	}
}

func TestHolidayCalendarGetUpcomingHolidays(t *testing.T) {
	calendar := NewHolidayCalendar("US")
	
	start := Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	upcoming := calendar.GetUpcomingHolidays(start, 3)
	
	if len(upcoming) < 3 {
		t.Errorf("Expected at least 3 upcoming holidays, got %d", len(upcoming))
	}
	
	// Should be in chronological order
	for i := 1; i < len(upcoming); i++ {
		if upcoming[i].Date.Before(upcoming[i-1].Date) {
			t.Errorf("Upcoming holidays not in chronological order")
		}
	}
	
	// All should be holidays
	for _, holiday := range upcoming {
		if !holiday.IsHoliday {
			t.Errorf("Non-holiday in upcoming holidays: %v", holiday)
		}
	}
}

func TestCalendarEntryString(t *testing.T) {
	entry := CalendarEntry{
		Date:        Date(2024, time.December, 25, 0, 0, 0, 0, time.UTC),
		IsWeekend:   false,
		IsHoliday:   true,
		HolidayName: "Christmas Day",
		IsToday:     false,
	}
	
	result := entry.String()
	expected := "2024-12-25 (Wed) [HOLIDAY: Christmas Day]"
	
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestDateTimeConvenienceMethods(t *testing.T) {
	dt := Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC)
	
	// Test enhanced business day calculator
	calc := dt.WithEnhancedBusinessDays("US")
	if calc == nil {
		t.Error("Expected non-nil enhanced business day calculator")
	}
	
	// Test holiday calendar
	calendar := dt.GetHolidayCalendar("US")
	if calendar == nil {
		t.Error("Expected non-nil holiday calendar")
	}
	
	// Test monthly holidays
	holidays := dt.GetMonthlyHolidays("US")
	if len(holidays) == 0 {
		t.Error("Expected at least one holiday in January")
	}
	
	// Test upcoming holidays
	upcoming := dt.GetUpcomingHolidays("US", 2)
	if len(upcoming) < 2 {
		t.Errorf("Expected at least 2 upcoming holidays, got %d", len(upcoming))
	}
}
