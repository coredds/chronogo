# ChronoGo

[![Version](https://img.shields.io/badge/version-v0.2.2-green.svg)](https://github.com/coredds/ChronoGo/releases)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/coredds/ChronoGo/actions/workflows/ci.yml/badge.svg)](https://github.com/coredds/ChronoGo/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/coredds/ChronoGo.svg)](https://pkg.go.dev/github.com/coredds/ChronoGo)

**ChronoGo** is a Go implementation inspired by Python's [Pendulum](https://pendulum.eustace.io/) library. It provides a powerful and easy-to-use datetime and timezone library that enhances Go's standard `time` package with a fluent API, better timezone handling, and human-friendly datetime operations.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [API Reference](#api-reference)
- [Advanced Usage](#advanced-usage)
- [API Compatibility](#api-compatibility)
- [Error Handling](#error-handling)
- [Testing](#testing)
- [Contributing](#contributing)
- [License & Related Projects](#license--related-projects)

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
- 🔌 **Serialization-ready**: JSON/Text marshalers and SQL driver integration
- ⏱️ **Unix helpers**: conversions and constructors for seconds/ms/µs/ns

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

## API Reference

### New API Methods (v0.2.0+)

| Category | Methods | Description |
|----------|---------|-------------|
| **Start/End** | `StartOfDay()`, `EndOfDay()`, `StartOfMonth()`, `EndOfMonth()`, `StartOfWeek()`, `EndOfWeek()`, `StartOfYear()`, `EndOfYear()`, `StartOfQuarter()`, `EndOfQuarter()` | Set DateTime to beginning or end of time periods |
| **Weekend/Weekday** | `IsWeekend()`, `IsWeekday()` | Check if date falls on weekend or weekday |
| **Quarter** | `Quarter()`, `StartOfQuarter()`, `EndOfQuarter()` | Quarter-based operations (Q1-Q4) |
| **ISO Week** | `ISOWeek()`, `ISOWeekYear()`, `ISOWeekNumber()` | ISO 8601 week operations |
| **Date Info** | `DayOfYear()` | Additional date information |
| **Date Utilities** | `IsFirstDayOfMonth()`, `IsLastDayOfMonth()`, `IsFirstDayOfYear()`, `IsLastDayOfYear()`, `WeekOfMonth()`, `WeekOfMonthISO()`, `WeekOfMonthWithStart(start time.Weekday)`, `DaysInMonth()`, `DaysInYear()` | Additional date utility methods for common date checks and calculations |
| **Fluent API** | `AddFluent()`, `Set()` | Method chaining for complex operations |
| **Enhanced Duration** | `NewDuration()`, `NewDurationFromComponents()` | Enhanced duration type with human-readable operations |
| **Unix Helpers** | `UnixMilli()`, `UnixMicro()`, `UnixNano()`, `FromUnixMilli()`, `FromUnixMicro()`, `FromUnixNano()` | Convert to/from various Unix time resolutions |
| **Serialization/DB** | `MarshalJSON()`, `UnmarshalJSON()`, `MarshalText()`, `UnmarshalText()`, `Value()`, `Scan()` | Seamless JSON/Text/SQL integration |

### Duration Operations
| Method | Description |
|--------|-------------|
| `Days()`, `Weeks()`, `Months()`, `Years()` | Get duration in different units |
| `HumanString()` | Human-readable representation |
| `Add()`, `Subtract()`, `Multiply()`, `Divide()` | Duration arithmetic |
| `IsPositive()`, `IsNegative()`, `IsZero()`, `Abs()` | Duration properties |

## Advanced Usage

### Parsing Different Formats

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

// Parse Unix timestamp (supports seconds/ms/µs/ns; signed)
dt7, _ := chronogo.TryParseUnix("1640995200")      // seconds
dt8, _ := chronogo.TryParseUnix("1640995200000")   // milliseconds
dt9, _ := chronogo.TryParseUnix("1640995200000000") // microseconds
dt10, _ := chronogo.TryParseUnix("1640995200000000000") // nanoseconds
```

### Examples

Explore rich examples in `example_test.go` which are runnable via `go test`. You can also build and run the demo:

```bash
go run ./cmd/chrono-demo
```

### String Formatting Options

```go
dt := chronogo.Date(2023, time.December, 25, 15, 30, 45, 0, time.UTC)

// Built-in formats
fmt.Println(dt.ToDateString())     // "2023-12-25"
fmt.Println(dt.ToTimeString())     // "15:30:45"
fmt.Println(dt.ToDateTimeString()) // "2023-12-25 15:30:45"
fmt.Println(dt.ToISO8601String())  // "2023-12-25T15:30:45Z"

// Custom formats (using Go's time format)
fmt.Println(dt.Format("Monday, January 2, 2006"))     // "Monday, December 25, 2023"
// Unix helpers
fmt.Println(dt.Unix())       // seconds
fmt.Println(dt.UnixMilli())  // milliseconds
fmt.Println(dt.UnixMicro())  // microseconds
fmt.Println(dt.UnixNano())   // nanoseconds

// Construct from Unix values in a specific location
fmt.Println(chronogo.FromUnixMilli(1703516445000, time.UTC))
fmt.Println(chronogo.FromUnixMicro(1703516445000000, time.UTC))
fmt.Println(chronogo.FromUnixNano(1703516445000000000, time.UTC))
fmt.Println(dt.Format("02/01/2006 15:04"))            // "25/12/2023 15:30"
fmt.Println(dt.Format("Jan 2, 2006 at 3:04 PM"))      // "Dec 25, 2023 at 3:30 PM"
```

### Date Comparisons and Checks

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

// Quarter operations
fmt.Println(dt.Quarter())           // 4 (Q4)
fmt.Println(dt.StartOfQuarter())    // "2023-10-01T00:00:00Z"
fmt.Println(dt.EndOfQuarter())      // "2023-12-31T23:59:59.999999999Z"

// Weekend/weekday detection
fmt.Println(dt.IsWeekend())      // true (Sunday)
fmt.Println(dt.IsWeekday())      // false

// Day and week information
fmt.Println(dt.DayOfYear())      // 288
year, week := dt.ISOWeek()
fmt.Printf("ISO Week: %d-%d\n", year, week) // "ISO Week: 2023-41"
fmt.Println(dt.ISOWeekYear())    // 2023
fmt.Println(dt.ISOWeekNumber())  // 41

// Additional date utilities
fmt.Println(dt.IsFirstDayOfMonth()) // false (15th is not first day)
fmt.Println(dt.IsLastDayOfMonth())  // false (15th is not last day)
fmt.Println(dt.WeekOfMonth())       // 3 (simple grouping by days 1-7, 8-14, ...)
fmt.Println(dt.WeekOfMonthISO())    // ISO-style week-of-month (Mon-start)
fmt.Println(dt.WeekOfMonthWithStart(time.Sunday)) // Sunday-start week-of-month
fmt.Println(dt.DaysInMonth())       // 31 (October has 31 days)
fmt.Println(dt.DaysInYear())        // 365 (2023 is not a leap year)
```

### Fluent API for Complex Operations

```go
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

// Fluent setting of date/time components
result := now.Set().
    Year(2024).
    Month(time.December).
    Day(25).
    Hour(15).
    Minute(30).
    Second(0).
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

// Duration arithmetic
sum := duration.Add(duration2)
diff := duration.Subtract(duration2)
product := duration.Multiply(2.5)

// Duration properties
fmt.Println(duration.IsPositive()) // true
fmt.Println(duration.IsNegative()) // false
fmt.Println(duration.IsZero())     // false
```

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

### Serialization & Database Integration

```go
dt := chronogo.Date(2023, time.December, 25, 15, 30, 45, 0, time.UTC)

// JSON
b, _ := dt.MarshalJSON()
var dt2 chronogo.DateTime
_ = dt2.UnmarshalJSON(b)

// Text
txt, _ := dt.MarshalText()
var dt3 chronogo.DateTime
_ = dt3.UnmarshalText(txt)

// SQL (database/sql)
// DateTime implements driver.Valuer and sql.Scanner
// so it can be used as a struct field mapped to TIMESTAMP columns.
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
make test
```

Run tests with coverage:

```bash
make cover
```

Run benchmarks:

```bash
make bench
```

### Continuous Integration

This repository includes GitHub Actions CI with a cross-platform, multi-version matrix:
- OS: Ubuntu, macOS, Windows
- Go: 1.21.x, 1.22.x

CI runs vet, unit tests, race tests where CGO is available, and publishes coverage as an artifact.

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
   make test
   ```

## License & Related Projects

### License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

### Inspiration
ChronoGo is inspired by Python's [Pendulum](https://pendulum.eustace.io/) library, which provides an excellent API for datetime manipulation. We've adapted its concepts to Go's type system and idioms while maintaining the intuitive and powerful interface that makes Pendulum so popular.

### Related Projects
- [Pendulum](https://pendulum.eustace.io/) - The original Python library that inspired ChronoGo
- [Carbon](https://github.com/golang-module/carbon) - Another Go datetime library
- [Now](https://github.com/jinzhu/now) - Go package for time manipulation

### Roadmap
**Completed in v0.2.0+ ✅**
- Enhanced utility methods (StartOfDay, EndOfDay, etc.)
- Weekend and weekday detection  
- Quarter operations and ISO week support
- Fluent API for method chaining
- Enhanced duration type with human-readable operations
- Additional date utility methods (v0.2.2)

**Planned Features**
- Localization support for human-readable strings
- Business day calculations
- Recurrence rules (RRULE support)
- Holiday calculations
- Duration parsing from strings
- More comprehensive DST transition handling
- Performance optimizations
- Benchmark tests and performance profiling

### Changelog
See [CHANGELOG.md](CHANGELOG.md) for a detailed history of changes.
