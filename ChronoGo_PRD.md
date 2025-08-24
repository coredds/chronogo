# Product Requirements Document (PRD)

## ChronoGo v0.5.0 - A Go Implementation with Advanced Parsing — A Powerful and Easy-to-Use Datetime and Timezone Library

## Purpose

Develop a Go library that mirrors the key features and API philosophy of the Python Pendulum library, enabling developers to work with dates, times, and timezones effectively and intuitively, with built-in support for timezone conversions, daylight saving time (DST), and human-friendly date/time manipulation.

## Background

Python's Pendulum is a popular datetime library that enhances and simplifies Python's built-in datetime functionalities. It fixes many timezone handling issues, especially with DST, and provides a rich, fluent API for date/time creation, manipulation, formatting, and comparison. A comparable Go library focusing on usability, correctness, and performance for timezone-aware datetime handling is valuable for the Go ecosystem.

## Key Features and Requirements
### 1. Drop-in Replacement Design

- Provide types (e.g., DateTime) that have similar methods and behavior to Go time.Time, facilitating easy adoption
- The core types should embed or wrap Go's standard time.Time for compatibility

### 2. Timezones & DST Handling

- Support IANA timezone database for timezone-aware datetimes
- Default storage in UTC but allow easy conversion to/from any timezone
- Correct handling of ambiguous times (e.g., during DST transitions)
- Provide utilities to query if a datetime is in DST, get timezone name, offset, etc.

### 3. Date and Time Creation

- Fluent constructors for datetime with timezone-aware defaults
- Convenient static methods for today, tomorrow, yesterday with timezone
- Create datetime from Unix timestamp, RFC 3339, ISO 8601, and common string formats

### 4. Date/Time Manipulation

- Methods to add/subtract years, months, days, hours, minutes, seconds, milliseconds
- Support for relative adjustments that correctly handle month overflow and leap years
- Methods for changing components (e.g., SetYear, SetMonth)

### 5. Comparison and Diff Utilities
- Compare two datetime objects including across different timezones

- Methods like .IsPast(), .IsFuture(), .IsLeapYear() as convenience

- Compute differences between datetimes returning human-readable periods (e.g., days, months, hours)

Format differences in a “diff for humans” style string.

### 6. Formatting and Localization
- Rich formatting methods inspired by Pendulum for pattern-based output

- Support for some localization of date/time strings (optional or pluggable)

### 7. Parsing
- Parsing methods for common date/time formats without needing explicit format strings

- Strict/lenient parsing options

### 8. Period and Interval Types
- Represent time intervals or ranges (optionally)

- Support iteration or human-friendly interval manipulation

### 9. Thread Safety and Performance
- Ensure safe concurrency where applicable

- Efficient implementation leveraging Go idioms and built-in time package optimizations

## Non-Functional Requirements
- Well-documented public API with examples, similar to Python Pendulum docs for ease of onboarding

- Unit tests covering edge cases around timezones and DST, leap seconds, and date arithmetic

- Semantic versioning, maintainable and extensible codebase

- Compatible with recent Go versions (Go 1.21+)

- Support for modules and integration with Go package managers

## Out of Scope for Initial Release
- Complex calendaring features beyond basic date/time manipulation

- Advanced localization beyond essential formatting

- GUI or interactive components

## Success Metrics
API adoption and positive developer feedback emphasizing ease of use and correctness compared to Go’s standard time package.

Coverage with thorough test suite, particularly handling critical timezone edge cases.

Performance comparable to or better than equivalent operations with Go's time.



## Key Features and API Breakdown
### 1. Core Types

```go
type DateTime struct {
    time.Time
}
```
- DateTime wraps Go's time.Time to extend functionality
- Internally stores time in UTC with conversion utilities for timezones
- Exposes fluent methods for operations and queries

### 2. Construction and Initialization

```go
func Now() DateTime
func Today(location *time.Location) DateTime
func Date(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) DateTime
func Parse(layout, value string) (DateTime, error)
func ParseISO(value string) (DateTime, error)  // parse ISO 8601 strings
func FromUnix(sec int64, nsec int64, loc *time.Location) DateTime
```
- Now() returns current datetime in local timezone

Today(location) returns today’s date at midnight in specified timezone.

Date(...) creates a DateTime similar to time.Date().

Parse and ParseISO for parsing string timestamps robustly.

FromUnix for creation from Unix timestamps.

### 3. Timezone and DST Support

```go
func (dt DateTime) In(loc *time.Location) DateTime        // convert to specified timezone
func (dt DateTime) UTC() DateTime                          // convert to UTC
func (dt DateTime) Location() *time.Location               // get current location/timezone
func (dt DateTime) IsDST() bool                            // whether the datetime is in daylight saving time
func LoadLocation(name string) (*time.Location, error)    // load IANA timezone
```
Handle timezone conversion crisply.

Expose DST detection.

Load IANA zones from Go’s tzdata or system.

### 4. Date/Time Arithmetic

```go
func (dt DateTime) AddYears(years int) DateTime
func (dt DateTime) AddMonths(months int) DateTime
func (dt DateTime) AddDays(days int) DateTime
func (dt DateTime) AddHours(hours int) DateTime
func (dt DateTime) AddMinutes(mins int) DateTime
func (dt DateTime) AddSeconds(secs int) DateTime
func (dt DateTime) Sub(dt2 DateTime) time.Duration
```

- Arithmetic methods for each unit

Handle overflow, leap years gracefully.

Sub returns difference as time.Duration.

### 5. Date Component Setters

```go
func (dt DateTime) SetYear(year int) DateTime
func (dt DateTime) SetMonth(month time.Month) DateTime
func (dt DateTime) SetDay(day int) DateTime
func (dt DateTime) SetHour(hour int) DateTime
func (dt DateTime) SetMinute(minute int) DateTime
func (dt DateTime) SetSecond(second int) DateTime
```
- Create new DateTime with updated fields (immutable API)

### 6. Comparison and Queries

```go
func (dt DateTime) Before(other DateTime) bool
func (dt DateTime) After(other DateTime) bool
func (dt DateTime) Equal(other DateTime) bool
func (dt DateTime) IsPast() bool           // before now
func (dt DateTime) IsFuture() bool         // after now
func (dt DateTime) IsLeapYear() bool
```
- Comparison methods for easy relational queries

### 7. Human-Friendly Differences

```go
func (dt DateTime) DiffForHumans(other DateTime) string
```
Returns a human-readable string like “3 hours ago”, “in 2 months”.

### 8. Formatting and String Output

```go
func (dt DateTime) Format(layout string) string
func (dt DateTime) ToISO8601String() string
func (dt DateTime) String() string  // default string representation
```
Use Go-style or custom format layouts inspired by Pendulum.

Standard string outputs include ISO8601.

### 9. Optional: Period Type (Intervals)

```go
type Period struct {
    Start DateTime
    End   DateTime
}

func NewPeriod(start, end DateTime) Period
func (p Period) Duration() time.Duration
func (p Period) Contains(dt DateTime) bool
```
- Represent time ranges and query membership

## Sample Usage Examples
```go
package main

import (
    "fmt"
    "time"
    "pendulumgo" // hypothetical package
)

func main() {
    loc, _ := pendulumgo.LoadLocation("America/New_York")

    now := pendulumgo.Now()
    fmt.Println("Now:", now)

    dt := pendulumgo.Date(2025, time.July, 30, 11, 0, 0, 0, loc)
    fmt.Println("Custom datetime:", dt)

    dt2 := dt.AddDays(5).AddHours(3)
    fmt.Println("After 5 days 3 hours:", dt2)

    fmt.Println("In UTC:", dt2.UTC())

    if dt2.IsDST() {
        fmt.Println("DST is in effect")
    }

    diff := dt2.Sub(now)
    fmt.Println("Duration from now:", diff)

    fmt.Println("Diff for humans:", dt2.DiffForHumans(now))

    fmt.Println("Formatted:", dt2.Format("2006-01-02 15:04:05 MST"))
}
```
## Non-Functional Requirements (Re-summarized)
- Thread-safe and immutable DateTime values

- Extensive unit tests for DST, locality, leap year, edge cases

- Documentation and idiomatic Go API

- Compatible with Go 1.21+

- Modular, pluggable architecture for timezone backend