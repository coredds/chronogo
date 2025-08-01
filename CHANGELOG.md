# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
  - `WeekOfYear()` - Get week number within the year

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

**Note**: This changelog will be updated as new versions are released.