# chronogo

[![Version](https://img.shields.io/badge/version-v0.6.8-green.svg)](https://github.com/coredds/chronogo/releases)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.23-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/coredds/chronogo/actions/workflows/ci.yml/badge.svg)](https://github.com/coredds/chronogo/actions/workflows/ci.yml)
[![Security](https://github.com/coredds/chronogo/actions/workflows/security.yml/badge.svg)](https://github.com/coredds/chronogo/actions/workflows/security.yml)
[![Codecov](https://codecov.io/gh/coredds/chronogo/branch/main/graph/badge.svg)](https://codecov.io/gh/coredds/chronogo)
[![Go Reference](https://pkg.go.dev/badge/github.com/coredds/chronogo.svg)](https://pkg.go.dev/github.com/coredds/chronogo)

A comprehensive Go datetime library inspired by Python's Pendulum. chronogo enhances Go's standard time package with natural language parsing, business date operations, localization support, and a fluent API for intuitive datetime manipulation.

## Features

### Core Capabilities
- **Natural Language Parsing**: Parse dates in 7 languages - English, Spanish, Portuguese, French, German, Chinese, Japanese (powered by godateparser)
- **Enhanced DateTime Type**: Extended functionality built on Go's time.Time
- **Fluent API**: Method chaining for readable, expressive code
- **Convenience Methods**: On() and At() for quick date/time setting
- **Diff Type**: Rich datetime differences with calendar-aware and precise calculations
- **ISO 8601 Support**: Long year detection, ordinal dates, week dates, intervals
- **Immutable Operations**: Thread-safe with all methods returning new instances

### Business Operations
- **Holiday Support**: 34 countries with comprehensive regional data (via goholiday)
- **Business Day Calculations**: Working day arithmetic with automatic holiday awareness
- **Enhanced Calculator**: High-performance operations with custom weekend support
- **Holiday-Aware Scheduler**: Intelligent scheduling that respects holidays and business days
- **Calendar Integration**: Holiday calendars with formatted output and tracking

### Localization
- **7 Locales for Formatting**: en-US, es-ES, fr-FR, de-DE, zh-Hans, pt-BR, ja-JP
- **Localized Date Formatting**: Format dates and times in multiple languages
- **Human-Readable Differences**: "2 hours ago", "hace 2 horas", "il y a 2 heures", "2時間前"
- **Ordinal Numbers**: Language-specific ordinal formatting (1st, 2nd, 3rd, 日, etc.)

### Advanced Features
- **Timezone Operations**: Proper DST handling with optimized conversions
- **Period Type**: Time intervals with powerful iteration capabilities
- **Comparison Methods**: Closest, Farthest, Between, and boolean checks
- **Testing Helpers**: Time mocking and freezing for deterministic tests
- **Serialization**: Built-in JSON/Text marshalers and SQL driver integration
- **High Performance**: Optimized operations with 90% test coverage

## Installation

```bash
go get github.com/coredds/chronogo
```

## Quick Start

```go
package main

import (
    "fmt"
    "time"
    "github.com/coredds/chronogo"
)

func main() {
    // Natural language parsing in multiple languages
    dt1, _ := chronogo.Parse("tomorrow")
    dt2, _ := chronogo.Parse("next Monday")
    dt3, _ := chronogo.Parse("3 days ago")
    dt4, _ := chronogo.Parse("mañana")        // Spanish
    dt5, _ := chronogo.Parse("明天")           // Chinese
    
    // Technical format parsing
    dt6, _ := chronogo.Parse("2024-01-15T14:30:00Z")  // ISO 8601
    dt7, _ := chronogo.Parse("1705329000")             // Unix timestamp
    dt8, _ := chronogo.Parse("2023-359")               // Ordinal date
    dt9, _ := chronogo.Parse("2023-W52-1")             // Week date
    
    // Convenience methods
    meeting := chronogo.Now().On(2024, time.June, 15).At(14, 30, 0)
    
    // ISO 8601 long year detection
    if chronogo.Date(2020, time.January, 1).IsLongYear() {
        fmt.Println("2020 has 53 ISO weeks")
    }
    
    // Rich Diff type
    start := chronogo.Date(2023, time.January, 15)
    end := chronogo.Date(2024, time.March, 20)
    diff := end.Diff(start)
    
    fmt.Printf("%d years, %d months\n", diff.Years(), diff.Months())
    fmt.Printf("%.2f days\n", diff.InDays())
    fmt.Println(diff.ForHumans())  // "1 year from now"
    
    // Business date calculations
    calc := chronogo.NewEnhancedBusinessDayCalculator("US")
    workday := calc.AddBusinessDays(chronogo.Today(), 5)
    
    // Holiday checking
    usChecker := chronogo.NewGoHolidayChecker("US")
    if usChecker.IsHoliday(chronogo.Today()) {
        fmt.Println("Holiday:", usChecker.GetHolidayName(chronogo.Today()))
    }
    
    // Holiday-aware scheduling
    scheduler := chronogo.NewHolidayAwareScheduler("US")
    meetings := scheduler.ScheduleRecurring(chronogo.Now(), 24*time.Hour, 10)
    
    // Period iteration
    period := chronogo.NewPeriod(chronogo.Now(), chronogo.Now().AddDays(7))
    for _, day := range period.Days() {
        fmt.Println(day.Format("Monday, January 2"))
    }
}
```

## Usage Examples

### Natural Language Parsing

```go
// Parse dates in multiple languages
dt1, _ := chronogo.Parse("tomorrow")
dt2, _ := chronogo.Parse("next Monday")
dt3, _ := chronogo.Parse("3 days ago")
dt4, _ := chronogo.Parse("mañana")        // Spanish
dt5, _ := chronogo.Parse("demain")        // French
dt6, _ := chronogo.Parse("明天")           // Chinese

// Configure languages
chronogo.SetDefaultParseLanguages("en", "es")

// Strict mode (technical formats only)
dt7, _ := chronogo.ParseStrict("2024-01-15T14:30:00Z")  // OK
dt8, _ := chronogo.ParseStrict("tomorrow")               // Error
```

### Convenience Methods

```go
// Quick date and time setting
dt := chronogo.Now()
  .On(2024, time.December, 25)
  .At(14, 30, 0)

// ISO 8601 long year detection
if dt.IsLongYear() {
    fmt.Println("This year has 53 ISO weeks")
}

// Boundary operations
start := dt.StartOfMonth()
end := dt.EndOfMonth()
firstMonday := dt.FirstWeekdayOf(time.Monday)
lastFriday := dt.LastWeekdayOf(time.Friday)
```

### Diff Type

```go
// Create difference between dates
start := chronogo.Date(2023, time.January, 15)
end := chronogo.Date(2025, time.March, 20)
diff := start.Diff(end)

// Calendar-aware differences
fmt.Printf("%d years, %d months\n", diff.Years(), diff.Months())

// Precise differences
fmt.Printf("%.2f years\n", diff.InYears())
fmt.Printf("%d total days\n", diff.InDays())

// Human-readable
fmt.Println(diff.ForHumans())      // "2 years from now"
fmt.Println(diff.CompactString())  // "2y 2m"

// Comparisons
if diff.IsPositive() {
    fmt.Println("Future date")
}
```

### Business Date Operations

```go
// Business day calculations
today := chronogo.Today()
deadline := today.AddBusinessDays(5)

// Holiday checking
if today.IsHoliday() {
    fmt.Println("It's a holiday!")
}

// Enhanced calculator with custom weekends
calc := chronogo.NewEnhancedBusinessDayCalculator("US")
calc.SetCustomWeekends([]time.Weekday{time.Friday, time.Saturday})

nextBiz := calc.NextBusinessDay(today)
bizDays := calc.BusinessDaysBetween(today, today.AddDays(30))

// Multi-country support (34 countries)
usChecker := chronogo.NewGoHolidayChecker("US")
ukChecker := chronogo.NewGoHolidayChecker("GB")
jpChecker := chronogo.NewGoHolidayChecker("JP")
```

### Holiday-Aware Scheduler

```go
// Create scheduler
scheduler := chronogo.NewHolidayAwareScheduler("US")

// Schedule recurring meetings (avoids holidays)
start := chronogo.Date(2025, time.September, 1)
meetings := scheduler.ScheduleRecurring(start, 7*24*time.Hour, 8)

// Monthly end-of-month reports
reports := scheduler.ScheduleMonthlyEndOfMonth(start, 6)

// Business days only
bizMeetings := scheduler.ScheduleBusinessDays(start, 10)
```

### Localization

```go
dt := chronogo.Date(2024, time.January, 15, 14, 30, 0, 0, time.UTC)

// Localized formatting (7 locales supported)
result, _ := dt.FormatLocalized("dddd, MMMM Do YYYY", "en-US")
// "Monday, January 15th 2024"

result, _ = dt.FormatLocalized("dddd, Do de MMMM de YYYY", "es-ES")
// "lunes, 15º de enero de 2024"

result, _ = dt.FormatLocalized("YYYY年MMMM Do dddd", "ja-JP")
// "2024年1月 15日 月曜日"

// Human-readable differences
past := chronogo.Now().AddHours(-2)
result, _ = past.HumanStringLocalized("en-US")  // "2 hours ago"
result, _ = past.HumanStringLocalized("es-ES")  // "hace 2 horas"
result, _ = past.HumanStringLocalized("fr-FR")  // "il y a 2 heures"
result, _ = past.HumanStringLocalized("ja-JP")  // "2時間前"

// Default locale
chronogo.SetDefaultLocale("ja-JP")
result = dt.FormatLocalizedDefault("YYYY年MMMM Do")
```

### Timezone Operations

```go
// Create and convert timezones
meeting := chronogo.Parse("2025-01-15 14:00")
  .InTimezone("America/New_York")

tokyo := meeting.InTimezone("Asia/Tokyo")
london := meeting.InTimezone("Europe/London")

// DST-aware conversions
fmt.Println("NYC:", meeting.Format("15:04"))
fmt.Println("Tokyo:", tokyo.Format("15:04"))
fmt.Println("London:", london.Format("15:04"))
```

### Period Operations

```go
// Create period
start := chronogo.Date(2025, time.January, 1)
end := chronogo.Date(2025, time.January, 31)
period := chronogo.NewPeriod(start, end)

// Iterate by day
for _, dt := range period.Days() {
    fmt.Println(dt.Format("Jan 2"))
}

// Custom intervals
for _, dt := range period.RangeByUnitSlice(chronogo.Hour, 6) {
    fmt.Println(dt.Format("Jan 2 15:04"))
}

// Check containment
inRange := period.Contains(chronogo.Date(2025, time.January, 15))
```

### Testing Helpers

```go
// Mock current time
chronogo.SetTestNow(chronogo.Date(2024, time.January, 1))
defer chronogo.ClearTestNow()

// Freeze time
chronogo.FreezeTime(chronogo.Date(2024, time.January, 1))
defer chronogo.UnfreezeTime()

// Time travel
chronogo.TravelTo(chronogo.Date(2024, time.June, 1))
defer chronogo.TravelBack()

// Scoped mocking with auto-cleanup
chronogo.WithTestNow(chronogo.Date(2024, time.January, 1), func() {
    // Test code here
})
```

## Supported Countries

Business date operations support 34 countries via goholiday integration:

US, GB, CA, AU, NZ, DE, FR, JP, IN, BR, MX, IT, ES, NL, KR, PT, PL, RU, CN, TH, SG, TR, UA, AT, BE, CH, CL, FI, IE, IL, NO, SE, AR, ID

## Dependencies

- **godateparser**: Natural language date parsing in 7 languages
- **goholiday**: Holiday data for 34 countries with regional subdivisions

## Documentation

Full API documentation: [pkg.go.dev/github.com/coredds/chronogo](https://pkg.go.dev/github.com/coredds/chronogo)

## Testing

```bash
go test -cover ./...
```

Current coverage: 90% with comprehensive edge case handling

## Security

Automated security scanning with GitHub Actions:
- CodeQL static analysis
- Go vulnerability database checks (govulncheck)
- Dependency vulnerability scanning
- License compliance checking

Security tools: gosec, Semgrep, Snyk, OSSF Scorecard

Run local security checks:
```bash
# Unix/Linux/macOS
make security
./scripts/security-check.sh

# Windows PowerShell
.\scripts\security-check.ps1
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure tests pass and linting is clean
5. Submit a pull request

## License

MIT License - see [LICENSE](LICENSE) file

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for detailed release notes.
