# ChronoGo

[![Version](https://img.shields.io/badge/version-v0.6.1-green.svg)](https://github.com/coredds/ChronoGo/releases)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.23-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/coredds/ChronoGo/actions/workflows/ci.yml/badge.svg)](https://github.com/coredds/ChronoGo/actions/workflows/ci.yml)
[![Security](https://github.com/coredds/ChronoGo/actions/workflows/security.yml/badge.svg)](https://github.com/coredds/ChronoGo/actions/workflows/security.yml)
[![Codecov](https://codecov.io/gh/coredds/ChronoGo/branch/main/graph/badge.svg)](https://codecov.io/gh/coredds/ChronoGo)
[![Go Reference](https://pkg.go.dev/badge/github.com/coredds/ChronoGo.svg)](https://pkg.go.dev/github.com/coredds/ChronoGo)

ChronoGo is a comprehensive Go datetime library inspired by Python's Pendulum. It provides a powerful, fluent API that enhances Go's standard time package with better timezone handling, human-friendly operations, and extensive business date functionality.

## Key Features

- **Enhanced DateTime Type**: Drop-in enhancement of Go's time.Time with extended functionality
- **Robust Timezone Support**: Proper DST handling with optimized timezone operations
- **Fluent API**: Method chaining for intuitive date/time manipulation
- **Human-Readable Output**: Time differences like "2 hours ago" and "in 3 days"
- **Immutable Operations**: All methods return new instances for thread safety
- **Period and Duration Types**: Time intervals with powerful iteration capabilities
- **Comprehensive Parsing**: Support for common datetime formats with intelligent detection
- **Business Date Operations**: Holiday checking, business day calculations, and working day arithmetic with GoHoliday integration
- **Serialization Support**: Built-in JSON/Text marshalers and SQL driver integration
- **High Performance**: Optimized operations with extensive test coverage (91.7%)

## Installation

```bash
go get github.com/coredds/ChronoGo
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/coredds/ChronoGo"
)

func main() {
    // Create and manipulate datetime instances
    dt := ChronoGo.Now().AddDays(3).InTimezone("America/New_York")
    fmt.Println(dt.HumanString()) // "in 3 days"
    
    // Business date calculations with enhanced performance
    calc := ChronoGo.NewEnhancedBusinessDayCalculator("US")
    workday := calc.AddBusinessDays(ChronoGo.Today(), 5)
    fmt.Println(workday.Format("2006-01-02"))
    
    // Holiday-aware scheduling
    scheduler := ChronoGo.NewHolidayAwareScheduler("US")
    meetings := scheduler.ScheduleRecurring(ChronoGo.Now(), 24*time.Hour, 10)
    
    // Holiday calendar integration
    calendar := ChronoGo.NewHolidayCalendar("US")
    upcoming := calendar.GetUpcomingHolidays(ChronoGo.Now(), 5)
    for _, holiday := range upcoming {
        fmt.Printf("Upcoming: %s\n", holiday.String())
    }
    
    // Multi-country holiday checking (GoHoliday v0.5.3+ supports 33 countries)
    usChecker := ChronoGo.NewGoHolidayChecker("US")
    brChecker := ChronoGo.NewGoHolidayChecker("BR") // Brazil
    inChecker := ChronoGo.NewGoHolidayChecker("IN") // India
    
    newYear := ChronoGo.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    if usChecker.IsHoliday(newYear) {
        fmt.Println("US Holiday:", usChecker.GetHolidayName(newYear))
    }
    if brChecker.IsHoliday(newYear) {
        fmt.Println("Brazil Holiday:", brChecker.GetHolidayName(newYear))
    }
    
    // Enhanced holiday operations with GoHoliday v0.5.3+
    // Get all holidays in a date range
    start := ChronoGo.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    end := ChronoGo.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)
    holidays := usChecker.GetHolidaysInRange(start, end)
    fmt.Printf("Q1 2024 US holidays: %d\n", len(holidays))
    
    // Batch holiday checking for performance
    dates := []ChronoGo.DateTime{
        ChronoGo.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),  // New Year's Day
        ChronoGo.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC),  // Independence Day
    }
    results := usChecker.AreHolidays(dates)
    fmt.Printf("Batch check results: %v\n", results)
    
    // Period iteration
    period := ChronoGo.NewPeriod(ChronoGo.Now(), ChronoGo.Now().AddDays(7))
    for _, day := range period.Days() {
        fmt.Println(day.Format("Monday, January 2"))
    }
}
```

## Core Components

### DateTime Operations
- **Creation**: Now(), Today(), Date(), FromUnix(), Parse()
- **Manipulation**: Add/Subtract time units with fluent API
- **Formatting**: Standard Go layouts plus human-readable output
- **Timezone**: Convert between timezones with proper DST handling
- **Comparison**: Before(), After(), Between(), Equal() methods

### Business Date Support
- **Holiday Management**: Integrated GoHoliday v0.5.3+ library with comprehensive multi-country holiday data (based on Vacanza holidays v0.80+)
- **Supported Countries**: 33 countries with 500+ regional subdivisions (US, GB, CA, AU, NZ, DE, FR, JP, IN, BR, MX, IT, ES, NL, KR, PT, PL, RU, CN, TH, SG, MY, ID, PH, VN, TW, HK, ZA, EG, NG, KE, GH, MA, TN)
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
- **Intelligent Parsing**: Automatic format detection for common datetime patterns
- **Multiple Formats**: ISO8601, RFC3339, Unix timestamps, and custom formats
- **JSON/SQL Support**: Built-in marshaling for database and API integration

## Dependencies

ChronoGo integrates with the following libraries:

- **GoHoliday**: Enterprise-grade holiday data library providing comprehensive holiday support for multiple countries with optimized lookup performance.

## Documentation

For detailed API documentation and examples, visit [pkg.go.dev/github.com/coredds/ChronoGo](https://pkg.go.dev/github.com/coredds/ChronoGo).

## Testing

Run the test suite with coverage:
```bash
go test -cover ./...
```

Current test coverage: 91.7% with comprehensive safety checks and edge case handling.

## Security

ChronoGo includes comprehensive security scanning and monitoring:

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

- **v0.6.1**: GoHoliday v0.5.3+ integration with 33 countries support, enhanced business operations, holiday-aware scheduling
- **v0.6.0**: Security hardening with comprehensive vulnerability scanning and dependency review automation
- **v0.5.0**: Advanced parsing functions, GoHoliday integration for enterprise holiday support
- **v0.4.3**: Enhanced test coverage (91.7%), improved safety checks, optimized DST handling
- **v0.4.2**: GitHub Actions CI/CD, comprehensive linting, automated dependency management
- **v0.4.0**: Business day operations, enhanced error handling, developer documentation
- **v0.3.0**: Holiday support, must functions, comprehensive validation
- **v0.2.0**: Fluent API, enhanced utilities, duration improvements

See [CHANGELOG.md](CHANGELOG.md) for detailed release notes.
