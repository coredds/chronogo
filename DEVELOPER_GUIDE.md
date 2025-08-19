# ChronoGo Developer Guide & Cookbook

This guide provides practical examples and patterns for using ChronoGo effectively in your applications.

## Table of Contents

- [Quick Start](#quick-start)
- [Business Date Operations](#business-date-operations)
- [Error Handling](#error-handling)
- [Performance Best Practices](#performance-best-practices)
- [Common Patterns](#common-patterns)
- [Migration from time Package](#migration-from-time-package)
- [Troubleshooting](#troubleshooting)

## Quick Start

### Basic Operations

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/coredds/ChronoGo"
)

func main() {
    // Current time
    now := chronogo.Now()
    fmt.Printf("Current time: %s\n", now)
    
    // Create specific dates
    christmas := chronogo.Date(2024, time.December, 25, 15, 30, 0, 0, time.UTC)
    fmt.Printf("Christmas 2024: %s\n", christmas)
    
    // Parse from strings
    dt, _ := chronogo.Parse("2024-12-25T15:30:45Z")
    fmt.Printf("Parsed: %s\n", dt)
    
    // Fluent arithmetic
    future := now.AddFluent().Years(1).Months(6).Days(10).To(now)
    fmt.Printf("Future date: %s\n", future)
    
    // Human-readable differences
    fmt.Printf("Christmas is: %s\n", christmas.DiffForHumans())
}
```

## Business Date Operations

### Setting Up Holiday Checkers

```go
// Use built-in US holidays
usHolidays := chronogo.NewUSHolidayChecker()

// Add custom holidays
companyHoliday := chronogo.Holiday{
    Name:  "Company Founding Day",
    Month: time.March,
    Day:   15,
}
usHolidays.AddHoliday(companyHoliday)

// Create custom holiday checker
type CustomHolidayChecker struct{}

func (c *CustomHolidayChecker) IsHoliday(dt chronogo.DateTime) bool {
    // Your custom logic here
    return dt.Month() == time.July && dt.Day() == 4 // Independence Day
}
```

### Business Day Calculations

```go
// Check if a date is a business day
date := chronogo.Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC)
isBusinessDay := date.IsBusinessDay(usHolidays)

// Find next/previous business day
nextBizDay := date.NextBusinessDay(usHolidays)
prevBizDay := date.PreviousBusinessDay(usHolidays)

// Add/subtract business days
futureDate := date.AddBusinessDays(5, usHolidays)
pastDate := date.SubtractBusinessDays(3, usHolidays)

// Count business days between dates
start := chronogo.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
end := chronogo.Date(2024, time.January, 31, 0, 0, 0, 0, time.UTC)
bizDays := start.BusinessDaysBetween(end, usHolidays)

// Count business days in a period
bizDaysInMonth := date.BusinessDaysInMonth(usHolidays)
bizDaysInYear := date.BusinessDaysInYear(usHolidays)
```

### Working with Recurring Holidays

```go
// Martin Luther King Jr. Day - third Monday in January
mlkDay := chronogo.Holiday{
    Name:    "Martin Luther King Jr. Day",
    Month:   time.January,
    WeekDay: &[]time.Weekday{time.Monday}[0],
    WeekNum: &[]int{3}[0],
}

// Memorial Day - last Monday in May
memorialDay := chronogo.Holiday{
    Name:    "Memorial Day",
    Month:   time.May,
    WeekDay: &[]time.Weekday{time.Monday}[0],
    WeekNum: &[]int{-1}[0], // -1 means last occurrence
}

checker := chronogo.NewUSHolidayChecker()
checker.AddHoliday(mlkDay)
checker.AddHoliday(memorialDay)

// Get all holidays for a year
holidays2024 := checker.GetHolidays(2024)
for _, holiday := range holidays2024 {
    fmt.Printf("Holiday: %s\n", holiday.Format("2006-01-02"))
}
```

## Error Handling

### Graceful Error Handling with Suggestions

```go
// Parse with error handling
dt, err := chronogo.Parse("invalid-date")
if err != nil {
    var chronoErr *chronogo.ChronoError
    if errors.As(err, &chronoErr) {
        fmt.Printf("Error: %s\n", chronoErr.Error())
        // This will include helpful suggestions
    }
}

// Timezone loading with suggestions
loc, err := chronogo.LoadLocation("EST")
if err != nil {
    var chronoErr *chronogo.ChronoError
    if errors.As(err, &chronoErr) {
        fmt.Printf("Timezone error: %s\n", chronoErr.Error())
        // Will suggest "Try 'America/New_York' for Eastern Time"
    }
}

// Validation with helpful messages
dt := chronogo.DateTime{} // zero value
if err := dt.Validate(); err != nil {
    fmt.Printf("Validation error: %s\n", err.Error())
    // Will suggest how to create valid DateTime
}
```

### Using Must Functions for Constants

```go
// Use Must* functions when you're sure the input is valid
var (
    AppLaunchDate = chronogo.MustParse("2024-01-01T00:00:00Z")
    UTCLocation   = chronogo.MustLoadLocation("UTC")
    ESTLocation   = chronogo.MustLoadLocation("America/New_York")
)

// These will panic if the input is invalid, which is appropriate for
// compile-time constants where you want to catch errors early
```

### Error Type Checking

```go
_, err := chronogo.Parse("bad-input")
if err != nil {
    // Check for specific error types
    if errors.Is(err, chronogo.ErrInvalidFormat) {
        // Handle format error specifically
    }
    
    // Check for ChronoError
    var chronoErr *chronogo.ChronoError
    if errors.As(err, &chronoErr) {
        if chronoErr.Op == "Parse" {
            // Handle parse-specific error
        }
    }
}
```

## Performance Best Practices

### Timezone Caching

```go
// Load timezones once and reuse
var (
    utc = chronogo.MustLoadLocation("UTC")
    est = chronogo.MustLoadLocation("America/New_York")
    pst = chronogo.MustLoadLocation("America/Los_Angeles")
)

// Use cached locations
dt := chronogo.NowIn(est)
```

### Efficient Business Day Calculations

```go
// Create holiday checker once and reuse
checker := chronogo.NewUSHolidayChecker()

// Add custom holidays once
for _, holiday := range customHolidays {
    checker.AddHoliday(holiday)
}

// Use the same checker for multiple calculations
for _, date := range dates {
    if date.IsBusinessDay(checker) {
        // Process business day
    }
}
```

### Batch Operations

```go
// Instead of multiple individual operations
dates := []chronogo.DateTime{date1, date2, date3}
results := make([]chronogo.DateTime, len(dates))

for i, date := range dates {
    results[i] = date.AddBusinessDays(5, checker)
}
```

## Common Patterns

### Working Days in a Month

```go
func getWorkingDaysInMonth(year int, month time.Month) []chronogo.DateTime {
    checker := chronogo.NewUSHolidayChecker()
    start := chronogo.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
    end := start.EndOfMonth()
    
    var workingDays []chronogo.DateTime
    current := start
    
    for !current.After(end) {
        if current.IsBusinessDay(checker) {
            workingDays = append(workingDays, current)
        }
        current = current.AddDays(1)
    }
    
    return workingDays
}
```

### Project Deadline Calculator

```go
func calculateProjectDeadline(startDate chronogo.DateTime, workingDays int) chronogo.DateTime {
    checker := chronogo.NewUSHolidayChecker()
    return startDate.AddBusinessDays(workingDays, checker)
}

// Usage
startDate := chronogo.Today()
deadline := calculateProjectDeadline(startDate, 20) // 20 working days
fmt.Printf("Project deadline: %s\n", deadline.Format("2006-01-02"))
```

### Time Range Iterator

```go
func processDateRange(start, end chronogo.DateTime, businessDaysOnly bool) {
    checker := chronogo.NewUSHolidayChecker()
    current := start
    
    for !current.After(end) {
        if !businessDaysOnly || current.IsBusinessDay(checker) {
            // Process date
            fmt.Printf("Processing: %s\n", current.Format("2006-01-02"))
        }
        current = current.AddDays(1)
    }
}
```

### SLA Calculation

```go
type SLA struct {
    ResponseHours int
    ResolutionDays int
}

func calculateSLADeadlines(ticketCreated chronogo.DateTime, sla SLA) (response, resolution chronogo.DateTime) {
    checker := chronogo.NewUSHolidayChecker()
    
    // Response time (in hours, including non-business hours)
    response = ticketCreated.AddHours(sla.ResponseHours)
    
    // Resolution time (in business days)
    resolution = ticketCreated.AddBusinessDays(sla.ResolutionDays, checker)
    
    return
}
```

## Migration from time Package

### Common Replacements

```go
// time package -> ChronoGo
time.Now()                    // -> chronogo.Now()
time.Date(...)                // -> chronogo.Date(...)
time.Parse(layout, value)     // -> chronogo.FromFormat(value, layout)
time.LoadLocation(name)       // -> chronogo.LoadLocation(name)

// New capabilities
dt.AddMonths(2)               // Calendar-aware arithmetic
dt.DiffForHumans()           // Human-readable differences
dt.IsBusinessDay()           // Business day checking
dt.AddBusinessDays(5)        // Business date arithmetic
```

### Wrapper for Existing Code

```go
// Create wrapper functions for gradual migration
func parseTime(value string) (time.Time, error) {
    dt, err := chronogo.Parse(value)
    return dt.Time, err
}

func addBusinessDays(t time.Time, days int) time.Time {
    dt := chronogo.Instance(t)
    checker := chronogo.NewUSHolidayChecker()
    return dt.AddBusinessDays(days, checker).Time
}
```

## Troubleshooting

### Common Issues

#### Parse Errors
```go
// Problem: Custom format not working
dt, err := chronogo.Parse("25/12/2023")
if err != nil {
    // Solution: Use FromFormat for custom formats
    dt, err = chronogo.FromFormat("25/12/2023", "02/01/2006")
}
```

#### Timezone Issues
```go
// Problem: Ambiguous timezone abbreviation
loc, err := chronogo.LoadLocation("EST")
if err != nil {
    // Solution: Use IANA timezone names
    loc, err = chronogo.LoadLocation("America/New_York")
}
```

#### Business Day Calculations
```go
// Problem: Unexpected business day count
checker := chronogo.NewUSHolidayChecker()

// Make sure to account for:
// 1. Weekends are excluded
// 2. Holidays are excluded  
// 3. The calculation is exclusive of the start date

start := chronogo.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC) // Holiday
end := chronogo.Date(2024, time.January, 5, 0, 0, 0, 0, time.UTC)   // Friday

// This counts business days BETWEEN the dates (exclusive)
bizDays := start.BusinessDaysBetween(end, checker)
```

### Debug Helpers

```go
// Check if date is valid
if err := dt.Validate(); err != nil {
    fmt.Printf("Invalid date: %s\n", err)
}

// Check timezone availability
if chronogo.IsValidTimezone("America/New_York") {
    fmt.Println("Timezone is valid")
}

// List available timezones
for _, tz := range chronogo.AvailableTimezones() {
    fmt.Printf("Available: %s\n", tz)
}
```

### Performance Profiling

```go
// Benchmark business day operations
func BenchmarkBusinessDays(b *testing.B) {
    checker := chronogo.NewUSHolidayChecker()
    start := chronogo.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        start.AddBusinessDays(10, checker)
    }
}
```
