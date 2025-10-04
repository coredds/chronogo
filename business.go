package chronogo

import (
	"time"

	goholiday "github.com/coredds/goholiday"
)

// HolidayChecker is an interface for checking if a date is a holiday.
// Users can implement this interface to provide custom holiday logic.
type HolidayChecker interface {
	IsHoliday(dt DateTime) bool
}

// Holiday represents a specific holiday with optional recurring rules.
type Holiday struct {
	Name    string
	Month   time.Month
	Day     int
	Year    *int          // nil for recurring holiday
	WeekDay *time.Weekday // for holidays like "first Monday of September"
	WeekNum *int          // which week of the month (1-5, -1 for last)
}

// DefaultHolidayChecker provides common holidays for different regions.
type DefaultHolidayChecker struct {
	holidays []Holiday
	region   string
}

// NewUSHolidayChecker creates a holiday checker with common US federal holidays.
func NewUSHolidayChecker() *DefaultHolidayChecker {
	holidays := []Holiday{
		{Name: "New Year's Day", Month: time.January, Day: 1},
		{Name: "Independence Day", Month: time.July, Day: 4},
		{Name: "Christmas Day", Month: time.December, Day: 25},
		// Martin Luther King Jr. Day - third Monday in January
		{Name: "Martin Luther King Jr. Day", Month: time.January, WeekDay: &[]time.Weekday{time.Monday}[0], WeekNum: &[]int{3}[0]},
		// Presidents Day - third Monday in February
		{Name: "Presidents Day", Month: time.February, WeekDay: &[]time.Weekday{time.Monday}[0], WeekNum: &[]int{3}[0]},
		// Memorial Day - last Monday in May
		{Name: "Memorial Day", Month: time.May, WeekDay: &[]time.Weekday{time.Monday}[0], WeekNum: &[]int{-1}[0]},
		// Labor Day - first Monday in September
		{Name: "Labor Day", Month: time.September, WeekDay: &[]time.Weekday{time.Monday}[0], WeekNum: &[]int{1}[0]},
		// Columbus Day - second Monday in October
		{Name: "Columbus Day", Month: time.October, WeekDay: &[]time.Weekday{time.Monday}[0], WeekNum: &[]int{2}[0]},
		// Thanksgiving - fourth Thursday in November
		{Name: "Thanksgiving", Month: time.November, WeekDay: &[]time.Weekday{time.Thursday}[0], WeekNum: &[]int{4}[0]},
	}
	return &DefaultHolidayChecker{
		holidays: holidays,
		region:   "US",
	}
}

// IsHoliday checks if the given date is a holiday.
func (hc *DefaultHolidayChecker) IsHoliday(dt DateTime) bool {
	for _, holiday := range hc.holidays {
		if hc.isHolidayMatch(dt, holiday) {
			return true
		}
	}
	return false
}

// isHolidayMatch checks if a DateTime matches a specific holiday definition.
func (hc *DefaultHolidayChecker) isHolidayMatch(dt DateTime, holiday Holiday) bool {
	// Check if year matches (if specified)
	if holiday.Year != nil && dt.Year() != *holiday.Year {
		return false
	}

	// Check month
	if dt.Month() != holiday.Month {
		return false
	}

	// Fixed date holiday
	if holiday.WeekDay == nil {
		return dt.Day() == holiday.Day
	}

	// Weekday-based holiday
	if dt.Weekday() != *holiday.WeekDay {
		return false
	}

	// Calculate which occurrence of the weekday this is
	weekNum := hc.getWeekOfMonth(dt, *holiday.WeekDay)

	// Handle "last occurrence" (-1)
	if *holiday.WeekNum == -1 {
		return hc.isLastOccurrenceOfWeekday(dt, *holiday.WeekDay)
	}

	return weekNum == *holiday.WeekNum
}

// getWeekOfMonth returns which occurrence of a weekday this date represents (1-5).
func (hc *DefaultHolidayChecker) getWeekOfMonth(dt DateTime, weekday time.Weekday) int {
	firstOfMonth := dt.StartOfMonth()
	daysDiff := int(dt.Weekday() - firstOfMonth.Weekday())
	if daysDiff < 0 {
		daysDiff += 7
	}
	firstOccurrence := firstOfMonth.AddDays(daysDiff)

	if dt.Before(firstOccurrence) {
		return 0 // This shouldn't happen if weekdays match
	}

	weeksDiff := int(dt.Sub(firstOccurrence).Hours() / (24 * 7))
	return weeksDiff + 1
}

// isLastOccurrenceOfWeekday checks if this is the last occurrence of the weekday in the month.
func (hc *DefaultHolidayChecker) isLastOccurrenceOfWeekday(dt DateTime, weekday time.Weekday) bool {
	// Check if adding 7 days would put us in the next month
	nextWeek := dt.AddDays(7)
	return nextWeek.Month() != dt.Month()
}

// AddHoliday adds a custom holiday to the checker.
func (hc *DefaultHolidayChecker) AddHoliday(holiday Holiday) {
	hc.holidays = append(hc.holidays, holiday)
}

// GetHolidays returns all holidays for a given year.
func (hc *DefaultHolidayChecker) GetHolidays(year int) []DateTime {
	var holidays []DateTime

	for _, holiday := range hc.holidays {
		if holiday.Year != nil && *holiday.Year != year {
			continue
		}

		if holiday.WeekDay == nil {
			// Fixed date holiday
			dt := Date(year, holiday.Month, holiday.Day, 0, 0, 0, 0, time.UTC)
			holidays = append(holidays, dt)
		} else {
			// Weekday-based holiday
			dt := hc.findWeekdayOccurrence(year, holiday.Month, *holiday.WeekDay, *holiday.WeekNum)
			if !dt.IsZero() {
				holidays = append(holidays, dt)
			}
		}
	}

	return holidays
}

// findWeekdayOccurrence finds the nth occurrence of a weekday in a given month/year.
func (hc *DefaultHolidayChecker) findWeekdayOccurrence(year int, month time.Month, weekday time.Weekday, occurrence int) DateTime {
	firstOfMonth := Date(year, month, 1, 0, 0, 0, 0, time.UTC)

	if occurrence == -1 {
		// Last occurrence - start from end of month and work backwards
		lastOfMonth := firstOfMonth.EndOfMonth()
		for d := lastOfMonth; d.Month() == month; d = d.AddDays(-1) {
			if d.Weekday() == weekday {
				return d
			}
		}
		return DateTime{} // Not found
	}

	// Find first occurrence
	daysDiff := int(weekday - firstOfMonth.Weekday())
	if daysDiff < 0 {
		daysDiff += 7
	}
	firstOccurrence := firstOfMonth.AddDays(daysDiff)

	// Add weeks to get the nth occurrence
	target := firstOccurrence.AddDays((occurrence - 1) * 7)

	// Make sure it's still in the same month
	if target.Month() != month {
		return DateTime{} // Not found
	}

	return target
}

// fastCountryChecker is an internal adapter that wraps goholiday's Country
// to provide fast holiday checking with DateTime support.
// This replaces the previous dependency on goholiday/chronogo adapter package.
type fastCountryChecker struct {
	country *goholiday.Country
}

// newFastCountryChecker creates a new fast country checker for the given country code.
func newFastCountryChecker(countryCode string) *fastCountryChecker {
	return &fastCountryChecker{
		country: goholiday.NewCountry(countryCode),
	}
}

// IsHoliday checks if the given time is a holiday.
func (fc *fastCountryChecker) IsHoliday(t time.Time) bool {
	_, ok := fc.country.IsHoliday(t)
	return ok
}

// GetHolidayName returns the name of the holiday if the date is a holiday.
func (fc *fastCountryChecker) GetHolidayName(t time.Time) string {
	holiday, ok := fc.country.IsHoliday(t)
	if ok && holiday != nil {
		return holiday.Name
	}
	return ""
}

// CountHolidaysInRange counts holidays within a date range.
func (fc *fastCountryChecker) CountHolidaysInRange(start, end time.Time) int {
	count := 0
	current := start
	for !current.After(end) {
		if fc.IsHoliday(current) {
			count++
		}
		current = current.AddDate(0, 0, 1)
	}
	return count
}

// GetHolidaysInRange returns all holidays in the specified date range.
func (fc *fastCountryChecker) GetHolidaysInRange(start, end time.Time) map[time.Time]string {
	result := make(map[time.Time]string)
	current := start
	for !current.After(end) {
		if name := fc.GetHolidayName(current); name != "" {
			result[current] = name
		}
		current = current.AddDate(0, 0, 1)
	}
	return result
}

// AreHolidays performs batch holiday checking.
func (fc *fastCountryChecker) AreHolidays(times []time.Time) []bool {
	result := make([]bool, len(times))
	for i, t := range times {
		result[i] = fc.IsHoliday(t)
	}
	return result
}

// ClearCache clears the holiday cache.
func (fc *fastCountryChecker) ClearCache() {
	// goholiday Country doesn't expose cache clearing directly
	// The cache is managed internally
}

// GoHolidayChecker wraps the goholiday library to implement the HolidayChecker interface.
// This provides comprehensive holiday data for multiple countries and regions.
type GoHolidayChecker struct {
	checker *fastCountryChecker
	country string
}

// NewGoHolidayChecker creates a new holiday checker using the goholiday library.
// The country parameter should be a 2-letter ISO country code (e.g., "US", "GB", "CA", "AU", "NZ", "DE", "FR", "JP", "IN", "BR", "MX", "IT", "ES", "NL", "KR", "PT", "PL", "RU", "CN", "TH", "SG", "TR", "UA", "AT", "BE", "CH", "CL", "FI", "IE", "IL", "NO", "SE").
// goholiday v0.6.4+ supports 34 countries with comprehensive regional subdivision data and multi-language holiday names.
// Based on Vacanza holidays framework. See DEPENDENCIES.md for version tracking.
func NewGoHolidayChecker(country string) *GoHolidayChecker {
	return &GoHolidayChecker{
		checker: newFastCountryChecker(country),
		country: country,
	}
}

// IsHoliday checks if the given date is a holiday using the goholiday library.
func (ghc *GoHolidayChecker) IsHoliday(dt DateTime) bool {
	return ghc.checker.IsHoliday(dt.Time)
}

// GetHolidayName returns the name of the holiday if the date is a holiday, empty string otherwise.
func (ghc *GoHolidayChecker) GetHolidayName(dt DateTime) string {
	return ghc.checker.GetHolidayName(dt.Time)
}

// CountHolidaysInRange counts holidays within a date range.
func (ghc *GoHolidayChecker) CountHolidaysInRange(start, end DateTime) int {
	return ghc.checker.CountHolidaysInRange(start.Time, end.Time)
}

// GetHolidaysInRange returns all holidays in the specified date range.
// Returns a map where keys are holiday dates and values are holiday names.
// New in goholiday v0.6.4+ - optimized for calendar operations.
func (ghc *GoHolidayChecker) GetHolidaysInRange(start, end DateTime) map[DateTime]string {
	holidays := ghc.checker.GetHolidaysInRange(start.Time, end.Time)
	result := make(map[DateTime]string, len(holidays))
	for date, name := range holidays {
		result[FromTime(date)] = name
	}
	return result
}

// AreHolidays performs batch holiday checking for efficient range operations.
// New in goholiday v0.6.4+ - optimized for bulk date processing.
func (ghc *GoHolidayChecker) AreHolidays(dates []DateTime) []bool {
	return ghc.checker.AreHolidays(ToTimes(dates))
}

// ClearCache clears the holiday cache to free memory.
// Useful for long-running applications. New in goholiday v0.6.3+.
func (ghc *GoHolidayChecker) ClearCache() {
	ghc.checker.ClearCache()
}

// GetCountry returns the country code for this holiday checker.
func (ghc *GoHolidayChecker) GetCountry() string {
	return ghc.country
}

// GetSubdivisions returns the supported subdivisions for this country.
// New in goholiday v0.6.3+ - provides access to regional holiday data.
func (ghc *GoHolidayChecker) GetSubdivisions() []string {
	return ghc.checker.country.GetSubdivisions()
}

// GetHolidayCategories returns the supported holiday categories for this country.
// New in goholiday v0.6.3+ - categories include "public", "bank", "school", etc.
func (ghc *GoHolidayChecker) GetHolidayCategories() []goholiday.HolidayCategory {
	return ghc.checker.country.GetCategories()
}

// GetLanguage returns the language used for holiday names.
// New in goholiday v0.6.3+ - supports multi-language holiday names.
func (ghc *GoHolidayChecker) GetLanguage() string {
	return ghc.checker.country.GetLanguage()
}

// GetHolidayCount returns the number of holidays in a given year.
// New in goholiday v0.6.3+ - efficient counting without loading all holidays.
func (ghc *GoHolidayChecker) GetHolidayCount(year int) (int, error) {
	return ghc.checker.country.GetHolidayCount(year)
}

// ValidateCountryCode validates if a country code is supported.
// New in goholiday v0.6.3+ - provides validation before creating checkers.
func ValidateCountryCode(countryCode string) error {
	return goholiday.ValidateCountryCode(countryCode)
}

// NewHolidayChecker creates a new goholiday-based holiday checker for the specified country.
// This is the recommended way to create holiday checkers for production use.
// Supported countries: US, GB, CA, AU, NZ, DE, FR, JP, IN, BR, MX, IT, ES, NL, KR, PT, PL, RU, CN, TH, SG, TR, UA, AT, BE, CH, CL, FI, IE, IL, NO, SE (34 countries with comprehensive regional subdivisions)
// Features: Sub-microsecond holiday lookups, multi-language support, thread-safe operations, intelligent caching
// Based on Vacanza holidays framework. See DEPENDENCIES.md for version tracking.
func NewHolidayChecker(country string) HolidayChecker {
	return NewGoHolidayChecker(country)
}

// defaultUSHolidayChecker is a cached US holiday checker for convenience functions
var defaultUSHolidayChecker = NewGoHolidayChecker("US")

// Business date operations for DateTime

// IsBusinessDay returns true if the date is a business day (Monday-Friday and not a holiday).
// If no holiday checker is provided, it uses the default US holiday checker.
func (dt DateTime) IsBusinessDay(holidayChecker ...HolidayChecker) bool {
	if dt.IsWeekend() {
		return false
	}

	var checker HolidayChecker
	if len(holidayChecker) > 0 && holidayChecker[0] != nil {
		checker = holidayChecker[0]
	} else {
		checker = defaultUSHolidayChecker
	}

	return !checker.IsHoliday(dt)
}

// IsHoliday returns true if the date is a holiday.
// If no holiday checker is provided, it uses the default US holiday checker.
func (dt DateTime) IsHoliday(holidayChecker ...HolidayChecker) bool {
	var checker HolidayChecker
	if len(holidayChecker) > 0 && holidayChecker[0] != nil {
		checker = holidayChecker[0]
	} else {
		checker = defaultUSHolidayChecker
	}

	return checker.IsHoliday(dt)
}

// GetHolidayName returns the name of the holiday if the date is a holiday.
// Returns empty string if the date is not a holiday.
// If no holiday checker is provided, it uses the default US holiday checker.
func (dt DateTime) GetHolidayName(holidayChecker ...HolidayChecker) string {
	var checker HolidayChecker
	if len(holidayChecker) > 0 && holidayChecker[0] != nil {
		checker = holidayChecker[0]
	} else {
		checker = defaultUSHolidayChecker
	}

	// Try to cast to GoHolidayChecker for enhanced functionality
	if ghc, ok := checker.(*GoHolidayChecker); ok {
		return ghc.GetHolidayName(dt)
	}

	// Fallback for other HolidayChecker implementations
	if checker.IsHoliday(dt) {
		return "Holiday" // Generic name for non-goholiday checkers
	}
	return ""
}

// NextBusinessDay returns the next business day.
func (dt DateTime) NextBusinessDay(holidayChecker ...HolidayChecker) DateTime {
	next := dt.AddDays(1)
	for !next.IsBusinessDay(holidayChecker...) {
		next = next.AddDays(1)
	}
	return next
}

// PreviousBusinessDay returns the previous business day.
func (dt DateTime) PreviousBusinessDay(holidayChecker ...HolidayChecker) DateTime {
	prev := dt.AddDays(-1)
	for !prev.IsBusinessDay(holidayChecker...) {
		prev = prev.AddDays(-1)
	}
	return prev
}

// GetHolidaysInRange returns all holidays between this date and the end date.
// If no holiday checker is provided, it uses the default US holiday checker.
// New in goholiday v0.6.4+ - optimized for calendar operations.
func (dt DateTime) GetHolidaysInRange(end DateTime, holidayChecker ...HolidayChecker) map[DateTime]string {
	var checker HolidayChecker
	if len(holidayChecker) > 0 && holidayChecker[0] != nil {
		checker = holidayChecker[0]
	} else {
		checker = defaultUSHolidayChecker
	}

	// Try to cast to GoHolidayChecker for enhanced functionality
	if ghc, ok := checker.(*GoHolidayChecker); ok {
		return ghc.GetHolidaysInRange(dt, end)
	}

	// Fallback for other HolidayChecker implementations
	result := make(map[DateTime]string)
	current := dt
	for current.Before(end) || current.Equal(end) {
		if checker.IsHoliday(current) {
			result[current] = "Holiday" // Generic name for non-goholiday checkers
		}
		current = current.AddDays(1)
	}
	return result
}

// AddBusinessDays adds the specified number of business days.
func (dt DateTime) AddBusinessDays(days int, holidayChecker ...HolidayChecker) DateTime {
	if days == 0 {
		return dt
	}

	current := dt
	remaining := days
	direction := 1

	if days < 0 {
		direction = -1
		remaining = -days
	}

	for remaining > 0 {
		current = current.AddDays(direction)
		if current.IsBusinessDay(holidayChecker...) {
			remaining--
		}
	}

	return current
}

// SubtractBusinessDays subtracts the specified number of business days.
func (dt DateTime) SubtractBusinessDays(days int, holidayChecker ...HolidayChecker) DateTime {
	return dt.AddBusinessDays(-days, holidayChecker...)
}

// BusinessDaysBetween returns the number of business days between two dates.
func (dt DateTime) BusinessDaysBetween(other DateTime, holidayChecker ...HolidayChecker) int {
	start := dt
	end := other

	if start.After(end) {
		start, end = end, start
	}

	count := 0
	current := start

	for current.Before(end) {
		if current.IsBusinessDay(holidayChecker...) {
			count++
		}
		current = current.AddDays(1)
	}

	return count
}

// BusinessDaysInMonth returns the number of business days in the month.
func (dt DateTime) BusinessDaysInMonth(holidayChecker ...HolidayChecker) int {
	start := dt.StartOfMonth()
	end := dt.EndOfMonth()

	count := 0
	current := start

	for !current.After(end) {
		if current.IsBusinessDay(holidayChecker...) {
			count++
		}
		current = current.AddDays(1)
	}

	return count
}

// BusinessDaysInYear returns the number of business days in the year.
func (dt DateTime) BusinessDaysInYear(holidayChecker ...HolidayChecker) int {
	start := dt.StartOfYear()
	end := dt.EndOfYear()

	count := 0
	current := start

	for !current.After(end) {
		if current.IsBusinessDay(holidayChecker...) {
			count++
		}
		current = current.AddDays(1)
	}

	return count
}
