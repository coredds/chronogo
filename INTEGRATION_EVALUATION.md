# Integration Evaluation: chronogo ↔ godateparser ↔ goholiday

**Date**: October 4, 2025  
**Version**: chronogo v0.6.8, godateparser v1.3.3, goholiday v0.6.5

## Executive Summary

The current integration between the three codebases is **functional but has several opportunities for optimization**. The main friction points are:

1. **Type conversion overhead** between `time.Time` and `chronogo.DateTime`
2. **Dual import pattern** for goholiday (both base and chronogo adapter)
3. **Missing direct DateTime support** in godateparser and goholiday
4. **API inconsistencies** in how DateTime wrapping is handled

## Current Architecture

### chronogo.DateTime Structure
```go
type DateTime struct {
    time.Time  // Embedded, not wrapped
}
```

**Key Insight**: `DateTime` embeds `time.Time`, making it compatible with any API that accepts `time.Time`, but requiring manual wrapping when converting back.

---

## 1. godateparser Integration

### Current State

**File**: `parse_natural.go`

```go
func parseWithGodateparser(value string, loc *time.Location, languages []string, preferFuture bool) (DateTime, error) {
    settings := &godateparser.Settings{
        Languages: languages,
    }
    
    if loc != nil {
        settings.RelativeBase = time.Now().In(loc)
    } else {
        settings.RelativeBase = time.Now().UTC()
    }
    
    result, err := godateparser.ParseDate(value, settings)
    if err != nil {
        return DateTime{}, ParseError(value, err)
    }
    
    // Manual conversion: time.Time -> DateTime
    if loc != nil && loc != result.Location() {
        result = result.In(loc)
    }
    
    return DateTime{result}, nil
}
```

### Pain Points

1. **Manual wrapping**: Every `godateparser.ParseDate()` call returns `time.Time`, requiring manual `DateTime{result}` wrapping
2. **Dual time.Now() calls**: chronogo has `Now()` but must use `time.Now()` for godateparser's `RelativeBase`
3. **PreferFuture not exposed**: Comment on line 48 indicates godateparser v1.3.3 may not support `PreferFuture`, but chronogo has it in `ParseConfig`
4. **No DateTime-aware API**: godateparser doesn't know about chronogo types

### Recommended Improvements

#### Option A: Add chronogo adapter to godateparser (Recommended)

**In godateparser**:
```go
// chronogo/adapter.go (new file in godateparser)
package chronogo

import (
    "time"
    "github.com/coredds/godateparser"
)

// DateTime interface that chronogo.DateTime satisfies
type DateTime interface {
    Time() time.Time
    Location() *time.Location
}

// ParseDateTime parses and returns a result compatible with chronogo
func ParseDateTime(value string, settings *godateparser.Settings) (time.Time, error) {
    return godateparser.ParseDate(value, settings)
}

// ParseWithRelativeBase is a convenience for chronogo integration
func ParseWithRelativeBase(value string, relativeBase time.Time, languages []string) (time.Time, error) {
    settings := &godateparser.Settings{
        Languages:    languages,
        RelativeBase: relativeBase,
    }
    return godateparser.ParseDate(value, settings)
}
```

**Benefits**:
- Keeps godateparser independent (no chronogo dependency)
- Provides a clear integration point
- Makes the relationship explicit

#### Option B: Add PreferFuture support to godateparser

**In godateparser Settings**:
```go
type Settings struct {
    Languages    []string
    RelativeBase time.Time
    PreferFuture bool  // NEW: When true, ambiguous dates prefer future
}
```

**Rationale**: This is a common feature in date parsers (e.g., Python's dateutil). When parsing "Friday" on a Wednesday, should it mean last Friday or next Friday?

---

## 2. goholiday Integration

### Current State

**File**: `business.go` (lines 1-8)

```go
import (
    "time"
    
    goholidays "github.com/coredds/goholiday"
    goholiday "github.com/coredds/goholiday/chronogo"  // Adapter package
)
```

**Problem**: This dual import pattern is confusing and error-prone:
- `goholidays` for base types (`Country`, `BusinessDayCalculator`, `HolidayCalendar`)
- `goholiday` for the chronogo adapter (`FastCountryChecker`)

### Pain Points

1. **Confusing import aliases**: Two imports from same package with similar names
2. **Manual DateTime wrapping everywhere**:
   ```go
   // Line 209: business.go
   func (ghc *GoHolidayChecker) IsHoliday(dt DateTime) bool {
       return ghc.checker.IsHoliday(dt.Time)  // Manual unwrap
   }
   
   // Line 225-231: business.go
   func (ghc *GoHolidayChecker) GetHolidaysInRange(start, end DateTime) map[DateTime]string {
       holidays := ghc.checker.GetHolidaysInRange(start.Time, end.Time)  // Manual unwrap
       result := make(map[DateTime]string, len(holidays))
       for date, name := range holidays {
           result[DateTime{Time: date}] = name  // Manual wrap
       }
       return result
   }
   
   // Line 236-241: business.go
   func (ghc *GoHolidayChecker) AreHolidays(dates []DateTime) []bool {
       times := make([]time.Time, len(dates))
       for i, dt := range dates {
           times[i] = dt.Time  // Manual unwrap for each element
       }
       return ghc.checker.AreHolidays(times)
   }
   ```

3. **Performance overhead**: Slice conversions for batch operations
4. **Map key conversions**: `map[time.Time]string` → `map[DateTime]string` requires full iteration

### Recommended Improvements

#### Option A: Generic-based adapter in goholiday (Recommended)

**In goholiday**:
```go
// adapter.go (new file in goholiday)
package goholiday

import "time"

// TimeProvider is an interface for types that can provide time.Time
// This allows goholiday to work with both time.Time and chronogo.DateTime
type TimeProvider interface {
    Time() time.Time
}

// IsHolidayGeneric works with any type that provides Time()
func (fc *FastCountryChecker) IsHolidayGeneric(tp TimeProvider) bool {
    return fc.IsHoliday(tp.Time())
}

// GetHolidaysInRangeGeneric returns holidays with generic time providers
func (fc *FastCountryChecker) GetHolidaysInRangeGeneric(start, end TimeProvider) map[time.Time]string {
    return fc.GetHolidaysInRange(start.Time(), end.Time())
}
```

**In chronogo**: Add method to satisfy interface
```go
// datetime.go
func (dt DateTime) Time() time.Time {
    return dt.Time  // Already embedded, just expose it
}
```

**Benefits**:
- No dependency on chronogo
- Works with any time-providing type
- Reduces conversion overhead
- Clear, idiomatic Go

#### Option B: Move chronogo adapter to chronogo (Cleaner separation)

**Current**: `goholiday/chronogo/adapter.go` (adapter lives in goholiday)  
**Proposed**: `chronogo/goholiday_adapter.go` (adapter lives in chronogo)

**Rationale**:
- goholiday shouldn't need to know about chronogo
- chronogo depends on goholiday, not vice versa
- Keeps goholiday focused on holiday logic
- Reduces maintenance burden on goholiday

**Implementation**:
```go
// chronogo/goholiday_adapter.go
package chronogo

import goholiday "github.com/coredds/goholiday"

// FastCountryChecker wraps goholiday's checker with DateTime support
type FastCountryChecker struct {
    checker *goholiday.FastCountryChecker
}

func NewFastCountryChecker(country string) *FastCountryChecker {
    return &FastCountryChecker{
        checker: goholiday.NewFastCountryChecker(country),
    }
}

func (fc *FastCountryChecker) IsHoliday(dt DateTime) bool {
    return fc.checker.IsHoliday(dt.Time)
}

// ... etc
```

This eliminates the dual import entirely.

---

## 3. Cross-Cutting Improvements

### 3.1 Unified Type Conversion

**Problem**: Manual `DateTime{time.Time}` wrapping is scattered throughout the codebase.

**Solution**: Add helper functions

```go
// datetime.go
func FromTime(t time.Time) DateTime {
    return DateTime{t}
}

func FromTimes(times []time.Time) []DateTime {
    result := make([]DateTime, len(times))
    for i, t := range times {
        result[i] = DateTime{t}
    }
    return result
}

func ToTimes(dates []DateTime) []time.Time {
    result := make([]time.Time, len(dates))
    for i, dt := range dates {
        result[i] = dt.Time
    }
    return result
}
```

**Usage**:
```go
// Before
for date, name := range holidays {
    result[DateTime{Time: date}] = name
}

// After
for date, name := range holidays {
    result[FromTime(date)] = name
}
```

### 3.2 Batch Operations Optimization

**Problem**: Converting `[]DateTime` ↔ `[]time.Time` for batch operations is expensive.

**Solution**: Add zero-copy conversion using unsafe (if acceptable) or optimize the adapter layer.

```go
// OPTION 1: Safe but requires copy
func (ghc *GoHolidayChecker) AreHolidays(dates []DateTime) []bool {
    times := ToTimes(dates)  // Helper function
    return ghc.checker.AreHolidays(times)
}

// OPTION 2: Zero-copy using unsafe (advanced)
func (ghc *GoHolidayChecker) AreHolidays(dates []DateTime) []bool {
    // Since DateTime embeds time.Time with no other fields,
    // we can cast the slice directly (requires careful validation)
    times := *(*[]time.Time)(unsafe.Pointer(&dates))
    return ghc.checker.AreHolidays(times)
}
```

**Note**: Option 2 only works because `DateTime` is defined as `struct { time.Time }` with no additional fields. This should be documented and validated.

### 3.3 Interface-Based Design

**Problem**: Tight coupling between libraries.

**Solution**: Define interfaces in chronogo that external libraries can satisfy.

```go
// chronogo/interfaces.go
package chronogo

// HolidayProvider can check holidays (already exists as HolidayChecker)
type HolidayChecker interface {
    IsHoliday(dt DateTime) bool
}

// DateParser can parse date strings
type DateParser interface {
    Parse(value string) (DateTime, error)
    ParseInLocation(value string, loc *time.Location) (DateTime, error)
}

// BusinessDayCalculator can perform business day calculations
type BusinessDayCalculator interface {
    IsBusinessDay(dt DateTime) bool
    AddBusinessDays(dt DateTime, days int) DateTime
    NextBusinessDay(dt DateTime) DateTime
}
```

This allows:
1. Easy mocking for tests
2. Alternative implementations
3. Clear contracts between libraries

---

## 4. Specific Recommendations by Priority

### High Priority (Immediate Impact)

1. **Move goholiday adapter from goholiday to chronogo**
   - Eliminates dual import confusion
   - Cleaner separation of concerns
   - Reduces goholiday maintenance burden
   - **Effort**: 2-3 hours
   - **Impact**: High (cleaner API, better maintainability)

2. **Add conversion helpers to chronogo**
   - `FromTime()`, `FromTimes()`, `ToTimes()`
   - Reduces boilerplate
   - Centralizes conversion logic
   - **Effort**: 30 minutes
   - **Impact**: Medium (cleaner code, easier to optimize later)

3. **Add PreferFuture to godateparser**
   - Completes the feature chronogo already exposes
   - Common feature in date parsers
   - **Effort**: 1-2 hours
   - **Impact**: Medium (feature completeness)

### Medium Priority (Quality of Life)

4. **Add chronogo adapter package to godateparser**
   - Optional `chronogo/adapter.go` in godateparser
   - Provides convenience functions
   - Documents the integration
   - **Effort**: 1 hour
   - **Impact**: Low (mostly documentation)

5. **Optimize batch operations**
   - Reduce slice conversion overhead
   - Consider zero-copy techniques
   - **Effort**: 2-3 hours
   - **Impact**: Medium (performance for large datasets)

### Low Priority (Nice to Have)

6. **Add interface-based design**
   - Define clear contracts
   - Enable alternative implementations
   - **Effort**: 1-2 hours
   - **Impact**: Low (flexibility for future)

7. **Add benchmarks for integration overhead**
   - Measure conversion costs
   - Identify bottlenecks
   - **Effort**: 2-3 hours
   - **Impact**: Low (visibility into performance)

---

## 5. Proposed Changes by Repository

### godateparser

```go
// settings.go
type Settings struct {
    Languages    []string
    RelativeBase time.Time
    PreferFuture bool  // NEW
}

// chronogo/adapter.go (NEW FILE)
package chronogo

// Optional convenience package for chronogo integration
// See documentation for usage examples
```

**Breaking Changes**: None (additive only)

### goholiday

```go
// REMOVE: chronogo/adapter.go (move to chronogo)
// This package will be deprecated in favor of chronogo managing its own adapter
```

**Breaking Changes**: Yes, but with migration path:
- Deprecate `goholiday/chronogo` package
- Document migration to `chronogo`'s internal adapter
- Provide 1-2 version overlap for transition

### chronogo

```go
// datetime.go
// Add conversion helpers
func FromTime(t time.Time) DateTime
func FromTimes(times []time.Time) []DateTime
func ToTimes(dates []DateTime) []time.Time

// goholiday_adapter.go (NEW FILE - moved from goholiday)
// Contains all goholiday integration code

// business.go
// Simplify imports:
import (
    "time"
    goholiday "github.com/coredds/goholiday"
)
// Remove dual import pattern
```

**Breaking Changes**: None (internal refactoring only)

---

## 6. Migration Path

### Phase 1: Immediate (No Breaking Changes)
1. Add conversion helpers to chronogo
2. Add PreferFuture to godateparser
3. Update chronogo to use new godateparser feature

### Phase 2: Deprecation (1-2 releases)
1. Create new adapter in chronogo for goholiday
2. Deprecate `goholiday/chronogo` package
3. Update documentation with migration guide

### Phase 3: Cleanup (After deprecation period)
1. Remove `goholiday/chronogo` package
2. Update chronogo to use internal adapter
3. Simplify imports

---

## 7. Testing Recommendations

### Integration Tests
Create a new test file: `integration_test.go`

```go
package chronogo_test

import (
    "testing"
    "github.com/coredds/chronogo"
    "github.com/coredds/godateparser"
    "github.com/coredds/goholiday"
)

// Test that all three libraries work together seamlessly
func TestFullIntegration(t *testing.T) {
    // Parse natural language with godateparser
    dt, err := chronogo.Parse("next Monday")
    if err != nil {
        t.Fatal(err)
    }
    
    // Check if it's a holiday with goholiday
    checker := chronogo.NewHolidayChecker("US")
    if dt.IsHoliday(checker) {
        t.Log("Next Monday is a holiday!")
    }
    
    // Add business days
    future := dt.AddBusinessDays(5, checker)
    t.Logf("5 business days from %s is %s", dt, future)
}

// Benchmark conversion overhead
func BenchmarkDateTimeConversion(b *testing.B) {
    times := make([]time.Time, 1000)
    for i := range times {
        times[i] = time.Now()
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        dates := chronogo.FromTimes(times)
        _ = chronogo.ToTimes(dates)
    }
}
```

---

## 8. Documentation Improvements

### README Updates

**chronogo/README.md**:
```markdown
## Integration with godateparser and goholiday

ChronoGo seamlessly integrates with two companion libraries:

- **godateparser**: Natural language date parsing in 7 languages
- **goholiday**: Holiday data for 34 countries

### Natural Language Parsing
ChronoGo uses godateparser internally for `Parse()`:
\`\`\`go
dt, _ := chronogo.Parse("tomorrow at 3pm")
dt, _ := chronogo.Parse("próximo viernes")  // Spanish
dt, _ := chronogo.Parse("来週の月曜日")      // Japanese
\`\`\`

### Holiday Support
ChronoGo uses goholiday for business day calculations:
\`\`\`go
checker := chronogo.NewHolidayChecker("US")
isHoliday := dt.IsHoliday(checker)
nextBizDay := dt.NextBusinessDay(checker)
\`\`\`

See [INTEGRATION.md](INTEGRATION.md) for detailed integration guide.
```

### New File: INTEGRATION.md

Create comprehensive integration guide showing:
- How the three libraries work together
- Performance characteristics
- Best practices
- Common patterns
- Troubleshooting

---

## 9. Conclusion

The current integration is **functional but not optimal**. The main issues are:

1. **Type conversion overhead**: Manual wrapping/unwrapping throughout
2. **Confusing import patterns**: Dual imports for goholiday
3. **Missing features**: PreferFuture in godateparser

**Recommended Action Plan**:

1. **Week 1**: Add conversion helpers to chronogo, add PreferFuture to godateparser
2. **Week 2**: Move goholiday adapter to chronogo, update documentation
3. **Week 3**: Deprecate old goholiday adapter, add integration tests
4. **Week 4**: Optimize batch operations, add benchmarks

**Expected Outcome**:
- Cleaner, more maintainable code
- Better performance for batch operations
- Clearer separation of concerns
- Easier for users to understand and use

**Estimated Total Effort**: 12-16 hours across all three repositories

---

## Appendix: Current Import Graph

```
chronogo
  ├─→ godateparser (for Parse)
  └─→ goholiday (for holidays)
       └─→ goholiday/chronogo (adapter)
            └─→ chronogo (circular!)
```

**Problem**: Circular dependency through adapter package

**Proposed**:
```
chronogo
  ├─→ godateparser (for Parse)
  └─→ goholiday (for holidays)

(No circular dependency)
```

---

**End of Evaluation**
