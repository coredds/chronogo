# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.6.1] - 2025-08-27

### Added
- Enhanced Business Day Calculator with optimized performance and custom weekend support
- Holiday-Aware Scheduler for intelligent scheduling that respects holidays and business days
- Holiday Calendar integration with month/year views and upcoming holiday tracking
- Support for custom weekend days (e.g., Friday-Saturday for Middle Eastern countries)
- Recurring scheduling with holiday avoidance (daily, weekly, monthly, quarterly)
- End-of-month business day scheduling with automatic holiday adjustment
- Calendar entries with comprehensive holiday information and formatting
- Upcoming holiday tracking with configurable count limits
- Business day-only scheduling functionality

### Enhanced
- GoHoliday dependency updated to latest version with expanded country support
- Documentation website with new demo cards showcasing enhanced features
- README with comprehensive examples of new business operations
- Test coverage for all new enhanced features

### Performance
- Optimized business day calculations using GoHoliday's enhanced algorithms
- Improved holiday lookup performance with integrated caching
- Enhanced scheduling algorithms for better performance with large date ranges

## [0.6.0] - 2025-08-26

### Added
- Comprehensive GitHub Actions security workflow with automated vulnerability scanning
- CodeQL security analysis for enhanced code security
- Dependency review automation to catch vulnerable dependencies
- Security policy documentation (SECURITY.md)
- Local security testing scripts for PowerShell and Bash
- govulncheck integration for Go vulnerability detection
- Hardcoded secret detection in security workflow

### Enhanced
- Repository security posture with enterprise-grade scanning
- Documentation with security badges and vulnerability reporting procedures
- Developer workflow with automated security checks

## [0.5.0] - 2025-08-24

### Added
- Advanced parsing functions supporting multiple datetime formats
- ISO 8601 ordinal date parsing (YYYY-DDD format)
- ISO 8601 week date parsing (YYYY-Www-D format)
- ISO 8601 interval parsing (start/end, start/duration, duration/end)
- Duration parsing for ISO 8601 format (P1Y2M3DT4H5M6S)
- Token-based format parsing with `FromFormatTokens()` function
- Fallback parsing with `ParseWithFallback()` for human-friendly formats
- Multiple format parsing with `ParseMultiple()` function
- Maximum leniency parsing with `ParseAny()` function
- Optimized parsing with format detection heuristics
- Parse options support (strict/lenient modes)
- Comprehensive validation functions (`IsValidDateTimeString()`)
- Support for Unix timestamp parsing (seconds, milliseconds, microseconds, nanoseconds)

### Changed
- Removed references to external library names for generic branding
- Improved error handling in parsing functions
- Enhanced code documentation and comments

### Fixed
- Token replacement conflicts in format conversion
- Error checking compliance in benchmark functions
- Code formatting and linting issues

## [0.4.3] - 2025-08-23

### Added
- Comprehensive optimization function tests with 91.7% total coverage
- Safety mechanisms for RangeByUnitSlice and FastRangeDays functions
- DST optimization with IsDSTOptimized and related caching functions
- Parse optimization with ParseOptimized and detectLayout functions
- Enhanced error handling for edge cases and invalid inputs

### Fixed
- TestRangeByUnitSlice crash prevention with iteration and capacity limits
- DST detection logic for proper daylight saving time identification
- Compact date format parsing (e.g., "20231225") prioritization
- Layout detection for space-separated datetime formats
- Empty string handling in utility functions

### Changed
- Improved README structure and cohesion without redundant information
- Enhanced code formatting and lint compliance
- Optimized timezone operations with caching mechanisms

### Security
- Added comprehensive safety checks to prevent infinite loops and memory exhaustion
- Implemented iteration limits (max 1000) and capacity constraints
- Enhanced input validation for all optimization functions

## [0.4.2] - 2025-08-20

### Added
- Comprehensive test coverage improvements for period.go
- Enhanced negative period handling tests
- Context cancellation tests for range methods
- String formatting tests for various period representations
- Edge case validation tests

### Changed
- Improved overall test coverage from 90.7% to 92.4%
- Enhanced period.go test coverage with 100% coverage for most functions

## [0.4.1] - 2025-08-19

### Added
- Codecov badge in README for test coverage visibility
- Enhanced release workflow with linting, race testing, and binary optimization
- Manual release workflow for triggered releases from GitHub UI
- SHA256 checksums generation for release binaries

### Changed
- Updated softprops/action-gh-release to v2
- Improved changelog extraction in release workflows

## [0.4.0] - 2025-08-19

### Added
- GitHub Actions workflows for CI, linting, and automated releases
- Comprehensive linting configuration with golangci-lint
- Dependabot configuration for automated dependency updates
- Enhanced test coverage to 90.7%
- Business date operations with holiday checking
- Error handling improvements with helpful suggestions
- JSON/Text marshalers and SQL `driver.Valuer`/`sql.Scanner` for `DateTime`
- Unix time helpers: `UnixMilli/UnixMicro/UnixNano` and `FromUnixMilli/Micro/Nano`
- Week-of-month helpers: `WeekOfMonthISO()` and `WeekOfMonthWithStart(start)`
- Makefile with common developer tasks
- Rounding and range utilities: `Truncate(unit)`, `Round(unit)`, `Clamp(min,max)`, `Between(a,b,inclusive)`; typed iteration `RangeByUnit(unit, step...)`
- Parsing: `ParseStrict`, `ParseStrictInLocation`, and ISO 8601 duration parsing via `ParseISODuration`

### Changed
- `IsDST()` now determines standard offset via minimum offset observed in the year for the location (robust across hemispheres)
- README updates: CI badge, Go Reference badge, Unix helpers, serialization/DB docs, examples section, Makefile usage
- README docs for rounding/range utilities and DST notes
- README docs for strict parsing and ISO 8601 duration parsing
- Improved README cohesion and organization

### Fixed
- Repository hygiene: added `.gitignore` for binaries/coverage; removed committed demo binary
- Removed unnecessary debug files and improved code organization

## [0.2.2] - 2025-01-08

### Added
- **Additional Date Utility Methods**: Added comprehensive date checking and calculation methods
  - `IsFirstDayOfMonth()` - Check if the date is the first day of the month
  - `IsLastDayOfMonth()` - Check if the date is the last day of the month (handles leap years)
  - `IsFirstDayOfYear()` - Check if the date is January 1st
  - `IsLastDayOfYear()` - Check if the date is December 31st
  - `WeekOfMonth()` - Get week number within month (1-6 based on days 1-7, 8-14, etc.)
  - `DaysInMonth()` - Get number of days in current month (properly handles leap years)
  - `DaysInYear()` - Get 365 or 366 based on leap year calculation

### Enhanced
- **Documentation**: Updated README with new utility methods examples and API reference
- **Test Coverage**: Added comprehensive tests for all new utility methods with edge case coverage
- **API Consistency**: New methods follow existing naming conventions and return patterns

## [0.2.1] - 2025-01-08

### Fixed
- **IsDST() Method**: Fixed incorrect logic that compared against UTC offset instead of standard winter offset
  - Now properly detects daylight saving time across all IANA timezones
  - Added comprehensive tests for multiple timezones and seasons (NY, London, UTC)

- **FluentDuration Accuracy**: Fixed inaccurate year/month approximations in calendar arithmetic
  - Calendar units (years/months) now stored separately from time units
  - Follows Go's time package behavior for month/year overflow and leap year handling
  - Eliminates approximation errors in business logic calculations

- **Period.Range Memory Safety**: Added context cancellation to prevent goroutine leaks
  - New `RangeWithContext()` method with proper cancellation support
  - Maintains backward compatibility with existing `Range()` method
  - Prevents memory leaks in long-running or abandoned iterations

### Enhanced
- **Test Organization**: Reorganized critical fix tests into appropriate test files by functionality
- **Code Quality**: Applied go fmt formatting and ensured GitHub CI compliance
- **Documentation**: Updated method documentation for fixed implementations

## [0.2.0] - 2025-01-08

### Added
- **Enhanced Utility Methods**: Added comprehensive start/end operations
  - `StartOfDay()`, `EndOfDay()` - Set time to beginning/end of day
  - `StartOfMonth()`, `EndOfMonth()` - Set to beginning/end of month
  - `StartOfWeek()`, `EndOfWeek()` - Set to beginning/end of week (Monday-Sunday)
  - `StartOfYear()`, `EndOfYear()` - Set to beginning/end of year
  - `StartOfQuarter()`, `EndOfQuarter()` - Set to beginning/end of quarter

- **Weekend and Weekday Detection**
  - `IsWeekend()` - Check if date falls on Saturday or Sunday
  - `IsWeekday()` - Check if date falls on Monday through Friday

- **Quarter Operations**
  - `Quarter()` - Get quarter number (1-4)
  - Quarter-based start/end operations

- **Enhanced Date Information**
  - `DayOfYear()` - Get day number within the year (1-366)
  - `ISOWeek()` - Get ISO 8601 year and week number
  - `ISOWeekYear()` - Get ISO 8601 year for the week
  - `ISOWeekNumber()` - Get ISO 8601 week number (1-53)

- **Fluent API for Enhanced Readability**
  - `AddFluent()` - Returns FluentDuration for chaining time additions
  - `Set()` - Returns FluentDateTime for chaining component setting
  - Method chaining for complex date/time operations
  - Improved code readability with builder pattern

- **Enhanced Duration Type (`ChronoDuration`)**
  - `NewDuration()` - Create from time.Duration
  - `NewDurationFromComponents()` - Create from hours, minutes, seconds
  - Human-readable duration methods: `Days()`, `Weeks()`, `Months()`, `Years()`
  - `HumanString()` - Human-readable string representation
  - Duration arithmetic: `Add()`, `Subtract()`, `Multiply()`, `Divide()`
  - Duration properties: `IsPositive()`, `IsNegative()`, `IsZero()`, `Abs()`

### Enhanced
- **Test Coverage**: Added comprehensive tests for all new API features
- **Demo Application**: Updated to showcase new API capabilities
- **Documentation**: Extensive README updates with examples for new features

## [0.1.1] - 2024-01-15

### Added
- Initial implementation of ChronoGo datetime library
- Core DateTime type with timezone-aware operations
- Comprehensive parsing support for common datetime formats
- Human-readable time differences with `DiffForHumans()`
- Period type for time intervals with iteration capabilities
- Fluent API with method chaining support
- Extensive test coverage with unit tests and examples
- Demo application showcasing library capabilities

### Features
- **DateTime Operations**: Creation, arithmetic, comparison, and formatting
- **Timezone Support**: Full IANA timezone database support with DST handling
- **Parsing**: ISO 8601, RFC 3339, and custom format parsing
- **Human-Friendly**: Age calculation and natural language time differences
- **Period Handling**: Time intervals with range iteration and calculations
- **Immutable API**: All operations return new instances
- **Thread-Safe**: Safe for concurrent use

### Supported Operations
- DateTime creation: `Now()`, `Today()`, `Tomorrow()`, `Yesterday()`, `Date()`
- Timezone operations: `In()`, `UTC()`, `Local()`, `IsDST()`
- Arithmetic: `AddYears()`, `AddMonths()`, `AddDays()`, etc.
- Setters: `SetYear()`, `SetMonth()`, `SetDay()`, etc.
- Comparisons: `Before()`, `After()`, `Equal()`, `IsPast()`, `IsFuture()`
- Formatting: `ToDateString()`, `ToISO8601String()`, custom `Format()`
- Human-readable: `DiffForHumans()`, `Age()`, `Humanize()`

### Documentation
- Comprehensive README with usage examples
- Extensive code documentation and comments
- Example tests demonstrating common use cases
- Product Requirements Document (PRD)

### Requirements
- Go 1.21 or later
- Compatible with standard library `time` package

## [0.1.0] - 2024-01-01

### Added
- Initial project setup and structure

---

[unreleased]: https://github.com/coredds/ChronoGo/compare/v0.6.0...HEAD
[0.6.0]: https://github.com/coredds/ChronoGo/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/coredds/ChronoGo/compare/v0.4.3...v0.5.0
[0.4.3]: https://github.com/coredds/ChronoGo/compare/v0.4.2...v0.4.3
[0.4.2]: https://github.com/coredds/ChronoGo/compare/v0.4.1...v0.4.2
[0.4.1]: https://github.com/coredds/ChronoGo/compare/v0.4.0...v0.4.1
[0.4.0]: https://github.com/coredds/ChronoGo/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/coredds/ChronoGo/compare/v0.2.2...v0.3.0
[0.2.2]: https://github.com/coredds/ChronoGo/compare/v0.2.1...v0.2.2
[0.2.1]: https://github.com/coredds/ChronoGo/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/coredds/ChronoGo/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/coredds/ChronoGo/releases/tag/v0.1.0

**Note**: This changelog will be updated as new versions are released.