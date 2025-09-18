package chronogo

import (
	"time"

	goholidays "github.com/coredds/GoHoliday"
)

// EnhancedBusinessDayCalculator wraps GoHoliday's optimized BusinessDayCalculator
// to provide high-performance business day calculations with extensive holiday support.
type EnhancedBusinessDayCalculator struct {
	calculator *goholidays.BusinessDayCalculator
	country    *goholidays.Country
}

// NewEnhancedBusinessDayCalculator creates a new enhanced business day calculator
// for the specified country with optimized performance.
func NewEnhancedBusinessDayCalculator(countryCode string) *EnhancedBusinessDayCalculator {
	country := goholidays.NewCountry(countryCode)
	calculator := goholidays.NewBusinessDayCalculator(country)

	return &EnhancedBusinessDayCalculator{
		calculator: calculator,
		country:    country,
	}
}

// SetCustomWeekends allows setting custom weekend days (e.g., Friday-Saturday for some regions)
func (ebc *EnhancedBusinessDayCalculator) SetCustomWeekends(weekends []time.Weekday) {
	ebc.calculator.SetWeekends(weekends)
}

// IsBusinessDay checks if a date is a business day using optimized algorithms
func (ebc *EnhancedBusinessDayCalculator) IsBusinessDay(dt DateTime) bool {
	return ebc.calculator.IsBusinessDay(dt.Time)
}

// AddBusinessDays adds business days using optimized calculation
func (ebc *EnhancedBusinessDayCalculator) AddBusinessDays(dt DateTime, days int) DateTime {
	result := ebc.calculator.AddBusinessDays(dt.Time, days)
	return DateTime{result}
}

// NextBusinessDay returns the next business day using optimized lookup
func (ebc *EnhancedBusinessDayCalculator) NextBusinessDay(dt DateTime) DateTime {
	result := ebc.calculator.NextBusinessDay(dt.Time)
	return DateTime{result}
}

// PreviousBusinessDay returns the previous business day using optimized lookup
func (ebc *EnhancedBusinessDayCalculator) PreviousBusinessDay(dt DateTime) DateTime {
	result := ebc.calculator.PreviousBusinessDay(dt.Time)
	return DateTime{result}
}

// BusinessDaysBetween calculates business days between dates with optimized algorithms
func (ebc *EnhancedBusinessDayCalculator) BusinessDaysBetween(start, end DateTime) int {
	return ebc.calculator.BusinessDaysBetween(start.Time, end.Time)
}

// IsEndOfMonth checks if a date is the last business day of the month.
// New in GoHoliday v0.6.3+ - useful for end-of-month business processes.
func (ebc *EnhancedBusinessDayCalculator) IsEndOfMonth(dt DateTime) bool {
	return ebc.calculator.IsEndOfMonth(dt.Time)
}

// Convenience methods for DateTime that use enhanced calculator
func (dt DateTime) WithEnhancedBusinessDays(countryCode string) *EnhancedBusinessDayCalculator {
	return NewEnhancedBusinessDayCalculator(countryCode)
}
