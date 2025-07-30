# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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