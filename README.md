# ChronoGo

[![Version](https://img.shields.io/badge/version-v0.2.1-green.svg)](https://github.com/coredds/ChronoGo/releases)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**ChronoGo** is a Go implementation inspired by Python's [Pendulum](https://pendulum.eustace.io/) library. It provides a powerful and easy-to-use datetime and timezone library that enhances Go's standard `time` package with a fluent API, better timezone handling, and human-friendly datetime operations.

## Features

- 🕒 **Drop-in enhancement** of Go's `time.Time` with extended functionality
- 🌍 **Robust timezone support** with proper DST handling
- 🔗 **Fluent API** with method chaining for intuitive date/time manipulation
- 📝 **Human-readable** time differences ("2 hours ago", "in 3 days")
- 🔄 **Immutable** datetime operations (methods return new instances)
- 📋 **Period and Duration** types for time intervals with iteration support
- 🎯 **Comprehensive parsing** for common datetime formats
- ✅ **Thread-safe** operations
- 🧪 **Well-tested** with extensive unit test coverage

## Installation

```bash
go get github.com/coredds/ChronoGo
```

## Quick Start

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
    fmt.Println("Now:", now)

    // Create specific datetime
    dt := chronogo.Date(2023, time.December, 25, 15, 30, 0, 0, time.UTC)
    fmt.Println("Christmas:", dt)

    // Fluent API with method chaining
    future := dt.AddYears(1).AddMonths(6).AddDays(15)
    fmt.Println("Future date:", future)

    // Human-friendly differences
    fmt.Println("Difference:", future.DiffForHumans(now))

    // Timezone conversions
    ny, _ := chronogo.LoadLocation("America/New_York")
    nyTime := dt.In(ny)
    fmt.Println("In New York:", nyTime)

    // Period iteration
    start := chronogo.Today()
    end := start.AddDays(7)
    period := chronogo.NewPeriod(start, end)
    
    for date := range period.RangeDays() {
        fmt.Println("Date:", date.ToDateString())
    }
}
```

## Core Features

### DateTime Creation

```go
// Current time
now := chronogo.Now()
nowUTC := chronogo.NowUTC()
nowInTZ := chronogo.NowIn(timezone)

// Today, tomorrow, yesterday
today := chronogo.Today()
tomorrow := chronogo.Tomorrow()
yesterday := chronogo.Yesterday()

// Specific date/time
dt := chronogo.Date(2023, time.December, 25, 15, 30, 45, 0, time.UTC)
utcTime := chronogo.UTC(2023, time.December, 25, 15, 30, 45, 0)

// From Unix timestamp
fromUnix := chronogo.FromUnix(1640995200, 0, time.UTC)
```

### Parsing

```go
// Parse common formats automatically
dt1, _ := chronogo.Parse("2023-12-25T15:30:45Z")
dt2, _ := chronogo.Parse("2023-12-25 15:30:45")
dt3, _ := chronogo.Parse("2023-12-25")

// Parse specific formats
dt4, _ := chronogo.ParseISO8601("2023-12-25T15:30:45Z")
dt5, _ := chronogo.ParseRFC3339("2023-12-25T15:30:45Z")

// Parse with custom format
dt6, _ := chronogo.FromFormat("25/12/2023 15:30", "02/01/2006 15:04")

// Parse Unix timestamp
dt7, _ := chronogo.TryParseUnix("1640995200")
```

### Timezone Support

```go
// Load timezones
ny, _ := chronogo.LoadLocation("America/New_York")
tokyo, _ := chronogo.LoadLocation("Asia/Tokyo")

// Convert between timezones
dt := chronogo.Now()
nyTime := dt.In(ny)
tokyoTime := dt.In(tokyo)
utcTime := dt.UTC()

// Check timezone properties
fmt.Println("Is DST:", dt.IsDST())
fmt.Println("Is UTC:", dt.IsUTC())
fmt.Println("Is Local:", dt.IsLocal())
```

### Date/Time Arithmetic

```go
dt := chronogo.Now()

// Add/subtract time units
future := dt.AddYears(1).AddMonths(6).AddDays(15)
past := dt.SubtractYears(1).SubtractMonths(3)

// Individual units
dt.AddHours(5)
dt.AddMinutes(30)
dt.AddSeconds(45)

// Set specific components (returns new instance)
newYear := dt.SetYear(2025)
newTime := dt.SetHour(18).SetMinute(30).SetSecond(0)
```

### Comparisons

```go
dt1 := chronogo.Date(2023, time.January, 15, 12, 0, 0, 0, time.UTC)
dt2 := chronogo.Date(2023, time.January, 16, 12, 0, 0, 0, time.UTC)

// Basic comparisons
fmt.Println(dt1.Before(dt2))  // true
fmt.Println(dt1.After(dt2))   // false
fmt.Println(dt1.Equal(dt2))   // false

// Convenience methods
fmt.Println(dt1.IsPast())     // true (if called after dt1)
fmt.Println(dt2.IsFuture())   // true (if called before dt2)
fmt.Println(dt1.IsLeapYear()) // false (2023 is not a leap year)
```

### Human-Friendly Formatting

```go
now := chronogo.Now()

// Relative to now
past := now.AddHours(-2)
future := now.AddDays(3)

fmt.Println(past.DiffForHumans())   // "2 hours ago"
fmt.Println(future.DiffForHumans()) // "in 3 days"

// Relative to another time
fmt.Println(future.DiffForHumans(past)) // "3 days after"

// Age calculation
birthdate := chronogo.Date(1990, time.May, 15, 0, 0, 0, 0, time.UTC)
fmt.Println(birthdate.Age()) // "33 years old" (approximate)
```

### Periods and Durations

```go
start := chronogo.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
end := chronogo.Date(2023, time.January, 10, 0, 0, 0, 0, time.UTC)

period := chronogo.NewPeriod(start, end)

// Period properties
fmt.Println("Duration:", period.Duration())
fmt.Println("Days:", period.Days())
fmt.Println("Hours:", period.Hours())
fmt.Println("Human format:", period.String())

// Check if date is in period
middle := chronogo.Date(2023, time.January, 5, 0, 0, 0, 0, time.UTC)
fmt.Println("Contains:", period.Contains(middle))

// Iterate over period
for date := range period.RangeDays() {
    fmt.Println("Date:", date.ToDateString())
}

// Custom iteration
for dt := range period.Range("hours", 6) { // Every 6 hours
    fmt.Println("Time:", dt.String())
}
```

### String Formatting

```go
dt := chronogo.Date(2023, time.December, 25, 15, 30, 45, 0, time.UTC)

// Built-in formats
fmt.Println(dt.ToDateString())     // "2023-12-25"
fmt.Println(dt.ToTimeString())     // "15:30:45"
fmt.Println(dt.ToDateTimeString()) // "2023-12-25 15:30:45"
fmt.Println(dt.ToISO8601String())  // "2023-12-25T15:30:45Z"
fmt.Println(dt.String())           // Same as ToISO8601String()

// Custom formats (using Go's time format)
fmt.Println(dt.Format("Monday, January 2, 2006"))     // "Monday, December 25, 2023"
fmt.Println(dt.Format("02/01/2006 15:04"))            // "25/12/2023 15:30"
fmt.Println(dt.Format("Jan 2, 2006 at 3:04 PM"))      // "Dec 25, 2023 at 3:30 PM"
```

### Enhanced Utility Methods

```go
dt := chronogo.Date(2023, time.October, 15, 14, 30, 45, 0, time.UTC)

// Start/End operations
fmt.Println(dt.StartOfDay())     // "2023-10-15T00:00:00Z"
fmt.Println(dt.EndOfDay())       // "2023-10-15T23:59:59.999999999Z"
fmt.Println(dt.StartOfMonth())   // "2023-10-01T00:00:00Z"
fmt.Println(dt.EndOfMonth())     // "2023-10-31T23:59:59.999999999Z"
fmt.Println(dt.StartOfWeek())    // "2023-10-09T00:00:00Z" (Monday)
fmt.Println(dt.EndOfWeek())      // "2023-10-15T23:59:59.999999999Z" (Sunday)
fmt.Println(dt.StartOfYear())    // "2023-01-01T00:00:00Z"
fmt.Println(dt.EndOfYear())      // "2023-12-31T23:59:59.999999999Z"

// Quarter operations
fmt.Println(dt.Quarter())           // 4 (Q4)
fmt.Println(dt.StartOfQuarter())    // "2023-10-01T00:00:00Z"
fmt.Println(dt.EndOfQuarter())      // "2023-12-31T23:59:59.999999999Z"

// Weekday/Weekend detection
fmt.Println(dt.IsWeekend())      // true (Sunday)
fmt.Println(dt.IsWeekday())      // false

// Day and week information
fmt.Println(dt.DayOfYear())      // 288
year, week := dt.ISOWeek()
fmt.Printf("ISO Week: %d-%d\n", year, week) // "ISO Week: 2023-41"
fmt.Println(dt.ISOWeekYear())    // 2023
   fmt.Println(dt.ISOWeekNumber())  // 41
```

### Fluent API for Enhanced Readability

```go
// Fluent duration building and application
now := chronogo.Now()

// Build complex durations with method chaining
future := now.AddFluent().
    Years(1).
    Months(2).
    Days(10).
    Hours(5).
    Minutes(30).
    Seconds(45).
    To(now)

// Or subtract from a date
past := now.AddFluent().
    Days(30).
    Hours(5).
    From(now)

// Fluent setting of date/time components
result := now.Set().
    Year(2024).
    Month(time.December).
    Day(25).
    Hour(15).
    Minute(30).
    Second(0).
    Timezone(timezone).
    Build()
```

### Enhanced Duration Type

```go
// Create enhanced duration
duration := chronogo.NewDuration(25*time.Hour + 30*time.Minute + 45*time.Second)
duration2 := chronogo.NewDurationFromComponents(2, 15, 30) // 2h 15m 30s

// Human-readable operations
fmt.Println(duration.String())           // "25h30m45s"
fmt.Println(duration.HumanString())      // "1 day"
fmt.Printf("Days: %.2f\n", duration.Days())     // "Days: 1.06"
fmt.Printf("Weeks: %.2f\n", duration.Weeks())   // "Weeks: 0.15"
fmt.Printf("Months: %.2f\n", duration.Months()) // "Months: 0.03"
fmt.Printf("Years: %.2f\n", duration.Years())   // "Years: 0.003"

// Duration arithmetic
sum := duration.Add(duration2)
diff := duration.Subtract(duration2)
product := duration.Multiply(2.5)
quotient := duration.Divide(2)

// Duration properties
fmt.Println(duration.IsPositive()) // true
fmt.Println(duration.IsNegative()) // false
fmt.Println(duration.IsZero())     // false
fmt.Println(duration.Abs())        // Absolute value
```

## Quick Reference

### New API Methods (v0.2.0)

| Category | Methods | Description |
|----------|---------|-------------|
| **Start/End** | `StartOfDay()`, `EndOfDay()`, `StartOfMonth()`, `EndOfMonth()`, `StartOfWeek()`, `EndOfWeek()`, `StartOfYear()`, `EndOfYear()`, `StartOfQuarter()`, `EndOfQuarter()` | Set DateTime to beginning or end of time periods |
| **Weekend/Weekday** | `IsWeekend()`, `IsWeekday()` | Check if date falls on weekend or weekday |
| **Quarter** | `Quarter()`, `StartOfQuarter()`, `EndOfQuarter()` | Quarter-based operations (Q1-Q4) |
| **ISO Week** | `ISOWeek()`, `ISOWeekYear()`, `ISOWeekNumber()` | ISO 8601 week operations |
| **Date Info** | `DayOfYear()` | Additional date information |
| **Fluent API** | `AddFluent()`, `Set()` | Method chaining for complex operations |
| **Enhanced Duration** | `NewDuration()`, `NewDurationFromComponents()` | Enhanced duration type with human-readable operations |

### Duration Operations
| Method | Description |
|--------|-------------|
| `Days()`, `Weeks()`, `Months()`, `Years()` | Get duration in different units |
| `HumanString()` | Human-readable representation |
| `Add()`, `Subtract()`, `Multiply()`, `Divide()` | Duration arithmetic |
| `IsPositive()`, `IsNegative()`, `IsZero()`, `Abs()` | Duration properties |

## API Compatibility

ChronoGo's `DateTime` type embeds Go's standard `time.Time`, making it a drop-in replacement in most cases. You can use all standard `time.Time` methods while gaining access to ChronoGo's enhanced functionality.

```go
dt := chronogo.Now()

// Standard time.Time methods work
fmt.Println(dt.Year(), dt.Month(), dt.Day())
fmt.Println(dt.Hour(), dt.Minute(), dt.Second())
fmt.Println(dt.Weekday())
fmt.Println(dt.Unix())

// Plus ChronoGo enhancements
fmt.Println(dt.DiffForHumans())
fmt.Println(dt.AddDays(5).SetHour(14))

// New v0.2.0 methods
fmt.Println(dt.StartOfDay())
fmt.Println(dt.IsWeekend())
fmt.Println(dt.Quarter())
```

## Error Handling

ChronoGo provides detailed error information for parsing operations:

```go
dt, err := chronogo.Parse("invalid-date")
if err != nil {
    // err is of type ParseError with details
    parseErr := err.(chronogo.ParseError)
    fmt.Printf("Failed to parse '%s': %s\n", parseErr.Input, parseErr.Reason)
}
```

## Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

Run benchmarks:

```bash
go test -bench=. ./...
```

## Contributing

Contributions are welcome! Please read our [contributing guidelines](CONTRIBUTING.md) and submit pull requests to our repository.

### Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/coredds/ChronoGo.git
   cd ChronoGo
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run tests:
   ```bash
   go test ./...
   ```

## Inspiration

ChronoGo is inspired by Python's [Pendulum](https://pendulum.eustace.io/) library, which provides an excellent API for datetime manipulation. We've adapted its concepts to Go's type system and idioms while maintaining the intuitive and powerful interface that makes Pendulum so popular.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Related Projects

- [Pendulum](https://pendulum.eustace.io/) - The original Python library that inspired ChronoGo
- [Carbon](https://github.com/golang-module/carbon) - Another Go datetime library
- [Now](https://github.com/jinzhu/now) - Go package for time manipulation

## Roadmap

### Completed in v0.2.0 ✅
- [x] Enhanced utility methods (StartOfDay, EndOfDay, etc.)
- [x] Weekend and weekday detection
- [x] Quarter operations and ISO week support
- [x] Fluent API for method chaining
- [x] Enhanced duration type with human-readable operations

### Planned Features
- [ ] Localization support for human-readable strings
- [ ] Business day calculations
- [ ] Recurrence rules (RRULE support)
- [ ] Holiday calculations
- [ ] Duration parsing from strings
- [ ] More comprehensive DST transition handling
- [ ] Performance optimizations
- [ ] Benchmark tests and performance profiling

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a detailed history of changes.