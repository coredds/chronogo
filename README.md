# chronogo

[![Version](https://img.shields.io/badge/version-v0.6.7-green.svg)](https://github.com/coredds/chronogo/releases)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.23-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/coredds/chronogo/actions/workflows/ci.yml/badge.svg)](https://github.com/coredds/chronogo/actions/workflows/ci.yml)
[![Security](https://github.com/coredds/chronogo/actions/workflows/security.yml/badge.svg)](https://github.com/coredds/chronogo/actions/workflows/security.yml)
[![Codecov](https://codecov.io/gh/coredds/chronogo/branch/main/graph/badge.svg)](https://codecov.io/gh/coredds/chronogo)
[![Go Reference](https://pkg.go.dev/badge/github.com/coredds/chronogo.svg)](https://pkg.go.dev/github.com/coredds/chronogo)

chronogo is a comprehensive Go datetime library inspired by Python's Pendulum. It provides a powerful, fluent API that enhances Go's standard time package with better timezone handling, human-friendly operations, and extensive business date functionality.

## Key Features

- **Enhanced DateTime Type**: Drop-in enhancement of Go's time.Time with extended functionality
- **Robust Timezone Support**: Proper DST handling with optimized timezone operations
- **Fluent API**: Method chaining for intuitive date/time manipulation
- **Human-Readable Output**: Time differences like "2 hours ago" and "in 3 days"
- **Localization Support**: Multi-language formatting and human-readable differences (6 locales)
- **Immutable Operations**: All methods return new instances for thread safety
- **Period and Duration Types**: Time intervals with powerful iteration capabilities
- **Natural Language Parsing**: Parse "tomorrow", "next Monday", "3 days ago" in multiple languages (EN, ES, PT, FR, DE, ZH, JA)
- **Business Date Operations**: Holiday checking, business day calculations, and working day arithmetic with goholiday integration
- **Serialization Support**: Built-in JSON/Text marshalers and SQL driver integration
- **High Performance**: Optimized operations with extensive test coverage (91.7%)

## Installation

```bash
go get github.com/coredds/chronogo
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/coredds/chronogo"
)

func main() {
    // Natural language parsing in multiple languages
    dt1, _ := chronogo.Parse("tomorrow")                // English
    dt2, _ := chronogo.Parse("next Monday")             // Relative dates
    dt3, _ := chronogo.Parse("3 days ago")              // Quantity expressions
    dt4, _ := chronogo.Parse("mañana")                  // Spanish for "tomorrow"
    dt5, _ := chronogo.Parse("明天")                     // Chinese for "tomorrow"
    
    // Traditional datetime parsing still works
    dt6, _ := chronogo.Parse("2024-01-15T14:30:00Z")   // ISO 8601
    dt7, _ := chronogo.Parse("1705329000")              // Unix timestamp
    
    // Create and manipulate datetime instances
    dt := chronogo.Now().AddDays(3).InTimezone("America/New_York")
    fmt.Println(dt.HumanString()) // "in 3 days"
    
    // Convenience methods for quick date/time modifications
    meeting := chronogo.Now().On(2024, time.June, 15).At(14, 30, 0)  // Set to June 15, 2024 at 14:30
    deadline := chronogo.Today().On(2024, time.December, 31).At(23, 59, 59)  // End of year deadline
    
    // Check if a year has 53 ISO weeks (long year)
    if chronogo.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).IsLongYear() {
        fmt.Println("2020 is a long year with 53 weeks")
    }
    
    // Rich Diff type for datetime differences
    start := chronogo.Date(2023, time.January, 15, 10, 0, 0, 0, time.UTC)
    end := chronogo.Date(2024, time.March, 20, 14, 30, 0, 0, time.UTC)
    diff := end.Diff(start)
    fmt.Printf("Calendar-aware: %d years, %d months\n", diff.Years(), diff.Months())
    fmt.Printf("Precise: %.2f days, %.2f hours\n", diff.InDays(), diff.InHours())
    fmt.Printf("Human-readable: %s\n", diff.ForHumans())
    fmt.Printf("Compact: %s\n", diff.CompactString())
    
    // Business date calculations with enhanced performance
    calc := chronogo.NewEnhancedBusinessDayCalculator("US")
    workday := calc.AddBusinessDays(chronogo.Today(), 5)
    fmt.Println(workday.Format("2006-01-02"))
    
    // Holiday-aware scheduling
    scheduler := chronogo.NewHolidayAwareScheduler("US")
    meetings := scheduler.ScheduleRecurring(chronogo.Now(), 24*time.Hour, 10)
    
    // Holiday calendar integration
    calendar := chronogo.NewHolidayCalendar("US")
    upcoming := calendar.GetUpcomingHolidays(chronogo.Now(), 5)
    for _, holiday := range upcoming {
        fmt.Printf("Upcoming: %s\n", holiday.String())
    }
    
    // Multi-country holiday checking (goholiday v0.6.3+ supports 34 countries)
    usChecker := chronogo.NewGoHolidayChecker("US")
    brChecker := chronogo.NewGoHolidayChecker("BR") // Brazil
    trChecker := chronogo.NewGoHolidayChecker("TR") // Turkey (new in v0.6.3)
    uaChecker := chronogo.NewGoHolidayChecker("UA") // Ukraine (new in v0.6.3)
    
    newYear := chronogo.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    if usChecker.IsHoliday(newYear) {
        fmt.Println("US Holiday:", usChecker.GetHolidayName(newYear))
    }
    if brChecker.IsHoliday(newYear) {
        fmt.Println("Brazil Holiday:", brChecker.GetHolidayName(newYear))
    }
    
    // Enhanced holiday operations with goholiday v0.6.3+
    
    // New features in v0.6.3: subdivision support, holiday categories, language support
    subdivisions := usChecker.GetSubdivisions()
    categories := usChecker.GetHolidayCategories()
    language := usChecker.GetLanguage()
    holidayCount, _ := usChecker.GetHolidayCount(2024)
    // Get all holidays in a date range
    start := chronogo.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    end := chronogo.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)
    holidays := usChecker.GetHolidaysInRange(start, end)
    fmt.Printf("Q1 2024 US holidays: %d\n", len(holidays))
    
    // Batch holiday checking for performance
    dates := []chronogo.DateTime{
        chronogo.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),  // New Year's Day
        chronogo.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC),  // Independence Day
    }
    results := usChecker.AreHolidays(dates)
    fmt.Printf("Batch check results: %v\n", results)
    
    // Period iteration
    period := chronogo.NewPeriod(chronogo.Now(), chronogo.Now().AddDays(7))
    for _, day := range period.Days() {
        fmt.Println(day.Format("Monday, January 2"))
    }
}
```

## Localization Support

chronogo provides comprehensive localization support for formatting dates and human-readable time differences in multiple languages.

### Supported Locales

- **en-US** (English - United States)
- **es-ES** (Spanish - Spain) 
- **fr-FR** (French - France)
- **de-DE** (German - Germany)
- **zh-Hans** (Chinese - Simplified)
- **pt-BR** (Portuguese - Brazil)

### Localized Formatting

```go
dt := chronogo.Date(2024, time.January, 15, 14, 30, 0, 0, time.UTC)

// English formatting
result, _ := dt.FormatLocalized("dddd, MMMM Do YYYY", "en-US")
fmt.Println(result) // "Monday, January 15th 2024"

// Spanish formatting  
result, _ = dt.FormatLocalized("dddd, Do de MMMM de YYYY", "es-ES")
fmt.Println(result) // "lunes, 15º de enero de 2024"

// French formatting
result, _ = dt.FormatLocalized("dddd Do MMMM YYYY", "fr-FR") 
fmt.Println(result) // "lundi 15e janvier 2024"

// German formatting
result, _ = dt.FormatLocalized("dddd, Do MMMM YYYY", "de-DE")
fmt.Println(result) // "Montag, 15. Januar 2024"

// Chinese formatting
result, _ = dt.FormatLocalized("YYYY年MMMM Do dddd", "zh-Hans")
fmt.Println(result) // "2024年一月 15日 星期一"

// Portuguese formatting
result, _ = dt.FormatLocalized("dddd, Do de MMMM de YYYY", "pt-BR")
fmt.Println(result) // "segunda-feira, 15º de janeiro de 2024"
```

### Localized Human-Readable Differences

```go
now := chronogo.Now()
past := now.AddHours(-2)
future := now.AddDays(3)

// English
result, _ := past.HumanStringLocalized("en-US")
fmt.Println(result) // "2 hours ago"

// Spanish  
result, _ = past.HumanStringLocalized("es-ES")
fmt.Println(result) // "hace 2 horas"

// French
result, _ = past.HumanStringLocalized("fr-FR") 
fmt.Println(result) // "il y a 2 heures"

// German
result, _ = past.HumanStringLocalized("de-DE")
fmt.Println(result) // "vor 2 Stunden"

// Chinese
result, _ = past.HumanStringLocalized("zh-Hans")
fmt.Println(result) // "2小时前"

// Portuguese
result, _ = past.HumanStringLocalized("pt-BR")
fmt.Println(result) // "há 2 horas"
```

### Default Locale Management

```go
// Set default locale for the application
chronogo.SetDefaultLocale("es-ES")

// Use default locale formatting
dt := chronogo.Now()
result := dt.FormatLocalizedDefault("dddd, MMMM Do")
fmt.Println(result) // Uses Spanish formatting

// Use default locale for human strings
result = dt.HumanStringLocalizedDefault()
fmt.Println(result) // Uses Spanish phrasing
```

### Locale Information

```go
// Get available locales
locales := chronogo.GetAvailableLocales()
fmt.Println(locales) // ["en-US", "es-ES", "fr-FR", "de-DE", "zh-Hans", "pt-BR"]

// Get month and weekday names
dt := chronogo.Date(2024, time.June, 15, 0, 0, 0, 0, time.UTC) // Monday

monthName, _ := dt.GetMonthName("fr-FR")
fmt.Println(monthName) // "juin"

weekdayName, _ := dt.GetWeekdayName("de-DE") 
fmt.Println(weekdayName) // "Montag"
```

## Core Components

### DateTime Operations
- **Creation**: Now(), Today(), Date(), FromUnix(), Parse()
- **Manipulation**: Add/Subtract time units with fluent API
- **Convenience Methods**: On(), At() for quick date/time setting
- **Differences**: Explicit Diff type with rich methods for time differences
- **Formatting**: Standard Go layouts plus human-readable output
- **Timezone**: Convert between timezones with proper DST handling
- **Comparison**: Before(), After(), Between(), Equal() methods
- **ISO 8601**: IsLongYear() for 53-week year detection

### Business Date Support
- **Holiday Management**: Integrated goholiday v0.6.3+ library with comprehensive multi-country holiday data (34 countries including Turkey and Ukraine)
- **Supported Countries**: 34 countries with comprehensive regional subdivisions (US, GB, CA, AU, NZ, DE, FR, JP, IN, BR, MX, IT, ES, NL, KR, PT, PL, RU, CN, TH, SG, TR, UA, AT, BE, CH, CL, FI, IE, IL, NO, SE, AR, ID)
- **Performance**: Sub-microsecond lookup performance with intelligent caching and thread-safe operations
- **Multi-language Support**: Holiday names available in multiple languages
- **Business Day Calculations**: Working day arithmetic with holiday awareness
- **Custom Holiday Support**: Implement HolidayChecker interface for organization-specific holidays
- **Enhanced Operations**: Holiday-aware scheduling, calendar integration, and recurring schedules

### Period and Duration
- **Period Type**: Represents time intervals between two datetime instances
- **Range Operations**: Iterate over periods by day, hour, or custom units
- **Duration Extensions**: Human-readable duration formatting and calculations

### Parsing and Serialization
- **Natural Language Parsing**: Powered by godateparser with support for 7 languages (English, Spanish, Portuguese, French, German, Chinese, Japanese)
- **Multi-Language NLP**: Parse "tomorrow", "mañana", "demain", "明天" automatically
- **Intelligent Parsing**: Automatic format detection for technical formats (ISO8601, RFC3339, Unix timestamps)
- **Custom Formats**: Token-based and Go layout format parsing
- **JSON/SQL Support**: Built-in marshaling for database and API integration

Example:
```go
// Natural language parsing
dt1, _ := chronogo.Parse("tomorrow")               // Future relative
dt2, _ := chronogo.Parse("3 days ago")             // Past relative
dt3, _ := chronogo.Parse("next Monday")            // Next weekday
dt4, _ := chronogo.Parse("mañana")                 // Spanish
dt5, _ := chronogo.Parse("demain")                 // French
dt6, _ := chronogo.Parse("明天")                    // Chinese

// Technical format parsing
dt7, _ := chronogo.Parse("2024-01-15T14:30:00Z")  // ISO 8601
dt8, _ := chronogo.Parse("1705329000")             // Unix timestamp

// Custom language configuration
chronogo.SetDefaultParseLanguages("en", "es")      // English and Spanish only
config := chronogo.ParseConfig{
    Languages: []string{"en", "fr"},                // Explicit languages
    Location:  time.UTC,
    Strict:    false,                               // Allow natural language
}
dt9, _ := chronogo.ParseWith("demain", config)

// Strict mode: technical formats only
dt10, _ := chronogo.ParseStrict("2024-01-15T14:30:00Z")  // OK
dt11, _ := chronogo.ParseStrict("tomorrow")               // Error: natural language disabled
```

## Dependencies

chronogo integrates with the following libraries:

- **godateparser**: Advanced natural language date parsing library supporting 7 languages with intelligent relative date handling.
- **goholiday**: Enterprise-grade holiday data library providing comprehensive holiday support for multiple countries with optimized lookup performance.

## Documentation

For detailed API documentation and examples, visit [pkg.go.dev/github.com/coredds/chronogo](https://pkg.go.dev/github.com/coredds/chronogo).

## Testing

Run the test suite with coverage:
```bash
go test -cover ./...
```

Current test coverage: 91.7% with comprehensive safety checks and edge case handling.

## Security

chronogo includes comprehensive security scanning and monitoring:

### Automated Security Scanning
- **GitHub Actions**: Automated security workflows on every push and PR
- **CodeQL Analysis**: Static application security testing
- **Vulnerability Scanning**: Go vulnerability database checks with `govulncheck`
- **Dependency Review**: Automated dependency vulnerability scanning
- **License Compliance**: Automated license compliance checking

### Security Tools Integration
- **gosec**: Go security checker for common security issues
- **Semgrep**: Static analysis for security patterns
- **Snyk**: Vulnerability scanning for dependencies
- **OSSF Scorecard**: Open source security best practices assessment

### Local Security Checks
Run security checks locally:
```bash
# Unix/Linux/macOS
make security

# Or run the script directly
./scripts/security-check.sh

# Windows PowerShell
.\scripts\security-check.ps1
```

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass and linting is clean
5. Submit a pull request

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Version History

- **v0.6.5**: Comprehensive localization support with 6 locales (en-US, es-ES, fr-FR, de-DE, zh-Hans, pt-BR), localized date formatting, human-readable differences, ordinal numbers, and AM/PM indicators
- **v0.6.2**: goholiday v0.6.3+ integration with 34 countries support (added Turkey and Ukraine), enhanced business operations, holiday-aware scheduling, subdivision and category support, new APIs for subdivisions and holiday categories
- **v0.6.1**: Maintenance release and documentation improvements
- **v0.6.0**: Security hardening with comprehensive vulnerability scanning and dependency review automation
- **v0.5.0**: Advanced parsing functions, goholiday integration for enterprise holiday support
- **goholiday v0.6.3**: Added Turkey and Ukraine support, subdivision/category APIs, enhanced performance, error handling improvements
- **v0.4.3**: Enhanced test coverage (91.7%), improved safety checks, optimized DST handling
- **v0.4.2**: GitHub Actions CI/CD, comprehensive linting, automated dependency management
- **v0.4.0**: Business day operations, enhanced error handling, developer documentation
- **v0.3.0**: Holiday support, must functions, comprehensive validation
- **v0.2.0**: Fluent API, enhanced utilities, duration improvements

See [CHANGELOG.md](CHANGELOG.md) for detailed release notes.
