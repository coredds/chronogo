package chronogo

import (
	"fmt"
	"time"

	goholidays "github.com/coredds/GoHoliday"
)

// CalendarEntry represents a day in a calendar with holiday information
type CalendarEntry struct {
	Date        DateTime
	IsWeekend   bool
	IsHoliday   bool
	HolidayName string
	IsToday     bool
}

// HolidayCalendar provides calendar functionality with holiday awareness
type HolidayCalendar struct {
	calendar *goholidays.HolidayCalendar
	country  *goholidays.Country
}

// NewHolidayCalendar creates a new holiday calendar for the specified country
func NewHolidayCalendar(countryCode string) *HolidayCalendar {
	country := goholidays.NewCountry(countryCode)
	calendar := goholidays.NewHolidayCalendar(country)
	
	return &HolidayCalendar{
		calendar: calendar,
		country:  country,
	}
}

// GenerateMonth generates calendar entries for a specific month
func (hc *HolidayCalendar) GenerateMonth(year int, month time.Month) []CalendarEntry {
	entries := hc.calendar.GenerateMonth(year, month)
	result := make([]CalendarEntry, len(entries))
	
	today := Now()
	
	for i, entry := range entries {
		dt := DateTime{entry.Date}
		holidayName := ""
		if entry.Holiday != nil {
			holidayName = entry.Holiday.Name
		}
		
		result[i] = CalendarEntry{
			Date:        dt,
			IsWeekend:   entry.IsWeekend,
			IsHoliday:   entry.IsHoliday,
			HolidayName: holidayName,
			IsToday:     dt.Format("2006-01-02") == today.Format("2006-01-02"),
		}
	}
	
	return result
}

// PrintMonth prints a formatted calendar for the specified month
func (hc *HolidayCalendar) PrintMonth(year int, month time.Month) {
	hc.calendar.PrintMonth(year, month)
}

// GetMonthlyHolidays returns all holidays in a specific month
func (hc *HolidayCalendar) GetMonthlyHolidays(year int, month time.Month) []CalendarEntry {
	entries := hc.GenerateMonth(year, month)
	var holidays []CalendarEntry
	
	for _, entry := range entries {
		if entry.IsHoliday {
			holidays = append(holidays, entry)
		}
	}
	
	return holidays
}

// GetYearlyHolidays returns all holidays in a specific year
func (hc *HolidayCalendar) GetYearlyHolidays(year int) []CalendarEntry {
	var allHolidays []CalendarEntry
	
	for month := time.January; month <= time.December; month++ {
		monthlyHolidays := hc.GetMonthlyHolidays(year, month)
		allHolidays = append(allHolidays, monthlyHolidays...)
	}
	
	return allHolidays
}

// GetUpcomingHolidays returns the next N holidays from the given date
func (hc *HolidayCalendar) GetUpcomingHolidays(from DateTime, count int) []CalendarEntry {
	var upcoming []CalendarEntry
	current := from
	
	for len(upcoming) < count && current.Year() <= from.Year()+2 { // Don't search more than 2 years ahead
		entries := hc.GenerateMonth(current.Year(), current.Month())
		
		for _, entry := range entries {
			if entry.IsHoliday && (entry.Date.After(from) || entry.Date.Equal(from)) {
				upcoming = append(upcoming, entry)
				if len(upcoming) >= count {
					break
				}
			}
		}
		
		current = current.AddMonths(1).StartOfMonth()
	}
	
	return upcoming
}

// FormatCalendarEntry formats a calendar entry for display
func (ce CalendarEntry) String() string {
	dateStr := ce.Date.Format("2006-01-02")
	weekday := ce.Date.Weekday().String()[:3]
	
	status := ""
	if ce.IsToday {
		status += " [TODAY]"
	}
	if ce.IsWeekend {
		status += " [WEEKEND]"
	}
	if ce.IsHoliday {
		status += fmt.Sprintf(" [HOLIDAY: %s]", ce.HolidayName)
	}
	
	return fmt.Sprintf("%s (%s)%s", dateStr, weekday, status)
}

// Convenience methods for DateTime
func (dt DateTime) GetHolidayCalendar(countryCode string) *HolidayCalendar {
	return NewHolidayCalendar(countryCode)
}

func (dt DateTime) GetMonthlyHolidays(countryCode string) []CalendarEntry {
	calendar := NewHolidayCalendar(countryCode)
	return calendar.GetMonthlyHolidays(dt.Year(), dt.Month())
}

func (dt DateTime) GetUpcomingHolidays(countryCode string, count int) []CalendarEntry {
	calendar := NewHolidayCalendar(countryCode)
	return calendar.GetUpcomingHolidays(dt, count)
}
