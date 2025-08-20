# ChronoGo

[![Version](https://img.shields.io/badge/version-v0.4.2-green.svg)](https://github.com/coredds/ChronoGo/releases)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/coredds/ChronoGo/actions/workflows/ci.yml/badge.svg)](https://github.com/coredds/ChronoGo/actions/workflows/ci.yml)
[![Codecov](https://codecov.io/gh/coredds/ChronoGo/branch/main/graph/badge.svg)](https://codecov.io/gh/coredds/ChronoGo)
[![Go Reference](https://pkg.go.dev/badge/github.com/coredds/ChronoGo.svg)](https://pkg.go.dev/github.com/coredds/ChronoGo)

**ChronoGo** is a Go implementation inspired by Python's [Pendulum](https://pendulum.eustace.io/) library. It provides a powerful and easy-to-use datetime and timezone library that enhances Go's standard `time` package with a fluent API, better timezone handling, and human-friendly datetime operations.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [API Reference](#api-reference)
- [Advanced Usage](#advanced-usage)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)
- [Roadmap](#roadmap)
- [Changelog](#changelog)

## Features

- **Drop-in enhancement** of Go's `time.Time` with extended functionality
- **Robust timezone support** with proper DST handling
- **Fluent API** with method chaining for intuitive date/time manipulation
- **Human-readable** time differences ("2 hours ago", "in 3 days")
- **Immutable** datetime operations (methods return new instances)
- **Period and Duration** types for time intervals with iteration support
- **Comprehensive parsing** for common datetime formats
- **Thread-safe** operations
- **Well-tested** with extensive unit test coverage
- **Serialization-ready**: JSON/Text marshalers and SQL driver integration
- **Unix helpers**: conversions and constructors for seconds/ms/µs/ns
- **Utilities**: Truncate/Round to common units; Clamp/Between range helpers; typed units for safe iteration
- **Business date operations**: Holiday checking, business day calculations, working day arithmetic

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
	// Create a new datetime instance
	dt := ChronoGo.Now().AddDays(3).InTimezone("America/New_York")
	fmt.Println(dt.HumanString()) // Output: "in 3 days"
}
```

## API Reference

Refer to the [Go Reference Documentation](https://pkg.go.dev/github.com/coredds/ChronoGo) for detailed API descriptions and examples.

## Advanced Usage

### Business Date Operations
```go
package main

import (
	"fmt"
	"github.com/coredds/ChronoGo"
)

func main() {
	// Subtract business days
	dt := ChronoGo.Today().SubtractBusinessDays(5)
	fmt.Println(dt) // Output: Date 5 business days ago
}
```

### Serialization
ChronoGo supports JSON and SQL serialization:
```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/coredds/ChronoGo"
)

func main() {
	dt := ChronoGo.Now()
	jsonData, _ := json.Marshal(dt)
	fmt.Println(string(jsonData)) // Output: JSON representation of datetime
}
```

### Period Iteration
```go
package main

import (
	"fmt"
	"github.com/coredds/ChronoGo"
)

func main() {
	period := ChronoGo.NewPeriod(ChronoGo.Now(), ChronoGo.Now().AddDays(10))
	for _, day := range period.Days() {
		fmt.Println(day)
	}
}
```

## Testing

Run tests and check coverage:
```bash
go test ./... -cover
```

## Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Write tests for your changes.
4. Ensure all tests pass.
5. Submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Roadmap
**Completed in v0.2.0+**
- Enhanced utility methods (StartOfDay, EndOfDay, etc.)
- Weekend and weekday detection  
- Quarter operations and ISO week support
- Fluent API for method chaining
- Enhanced duration type with human-readable operations
- Additional date utility methods (v0.2.2)

**Completed in v0.3.0+**
- Business day calculations and holiday support
- Enhanced error handling with helpful suggestions
- Must functions for constants
- Comprehensive developer documentation

**Completed in v0.4.0+**
- GitHub Actions CI/CD pipeline with automated testing and releases
- Comprehensive linting and code quality checks
- Enhanced test coverage (90.7%)
- Automated dependency management with Dependabot

**Planned Features**
- Localization support for human-readable strings
- Recurrence rules (RRULE support)
- More international holiday sets
- Duration parsing from strings
- More comprehensive DST transition handling
- Performance optimizations
- Benchmark tests and performance profiling

## Changelog
See [CHANGELOG.md](CHANGELOG.md) for a detailed history of changes.
