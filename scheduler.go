package chronogo

import (
	"time"

	goholidays "github.com/coredds/GoHoliday"
)

// HolidayAwareScheduler provides intelligent scheduling that respects holidays and business days
type HolidayAwareScheduler struct {
	scheduler *goholidays.HolidayAwareScheduler
	country   *goholidays.Country
}

// NewHolidayAwareScheduler creates a new scheduler for the specified country
func NewHolidayAwareScheduler(countryCode string) *HolidayAwareScheduler {
	country := goholidays.NewCountry(countryCode)
	scheduler := goholidays.NewHolidayAwareScheduler(country)

	return &HolidayAwareScheduler{
		scheduler: scheduler,
		country:   country,
	}
}

// ScheduleRecurring creates a recurring schedule that avoids holidays and weekends
// frequency: how often to schedule (e.g., 24*time.Hour for daily, 7*24*time.Hour for weekly)
// count: how many occurrences to generate
func (has *HolidayAwareScheduler) ScheduleRecurring(start DateTime, frequency time.Duration, count int) []DateTime {
	times := has.scheduler.ScheduleRecurring(start.Time, frequency, count)
	result := make([]DateTime, len(times))
	for i, t := range times {
		result[i] = DateTime{t}
	}
	return result
}

// ScheduleMonthlyEndOfMonth schedules events at the end of each month, adjusting for holidays
// months: number of months to schedule for
func (has *HolidayAwareScheduler) ScheduleMonthlyEndOfMonth(start DateTime, months int) []DateTime {
	times := has.scheduler.ScheduleMonthlyEndOfMonth(start.Time, months)
	result := make([]DateTime, len(times))
	for i, t := range times {
		result[i] = DateTime{t}
	}
	return result
}

// ScheduleQuarterly schedules events quarterly, avoiding holidays
func (has *HolidayAwareScheduler) ScheduleQuarterly(start DateTime, quarters int) []DateTime {
	return has.ScheduleRecurring(start, 90*24*time.Hour, quarters) // Approximate quarterly
}

// ScheduleWeekly schedules weekly events, moving to next business day if needed
func (has *HolidayAwareScheduler) ScheduleWeekly(start DateTime, weeks int) []DateTime {
	return has.ScheduleRecurring(start, 7*24*time.Hour, weeks)
}

// ScheduleBusinessDays schedules events on business days only
func (has *HolidayAwareScheduler) ScheduleBusinessDays(start DateTime, days int) []DateTime {
	calculator := goholidays.NewBusinessDayCalculator(has.country)
	result := make([]DateTime, 0, days)
	current := start

	for i := 0; i < days; i++ {
		if calculator.IsBusinessDay(current.Time) {
			result = append(result, current)
		}
		current = current.AddDays(1)
		// Skip to next business day if needed
		for !calculator.IsBusinessDay(current.Time) && i < days-1 {
			current = current.AddDays(1)
		}
	}

	return result
}

// Convenience methods for DateTime
func (dt DateTime) ScheduleRecurring(countryCode string, frequency time.Duration, count int) []DateTime {
	scheduler := NewHolidayAwareScheduler(countryCode)
	return scheduler.ScheduleRecurring(dt, frequency, count)
}

func (dt DateTime) ScheduleMonthlyEndOfMonth(countryCode string, months int) []DateTime {
	scheduler := NewHolidayAwareScheduler(countryCode)
	return scheduler.ScheduleMonthlyEndOfMonth(dt, months)
}
