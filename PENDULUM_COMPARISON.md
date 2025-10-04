# chronogo vs Pendulum Feature Comparison

## Overview

chronogo is inspired by Python's Pendulum library ([https://pendulum.eustace.io/](https://pendulum.eustace.io/)). This document provides a comprehensive comparison between chronogo's current implementation and Pendulum's features to identify gaps and opportunities for enhancement.

## Feature Comparison Matrix

### Core Features

| Feature | Pendulum | chronogo | Status | Notes |
|---------|----------|----------|--------|-------|
| **DateTime Type** | ✅ Drop-in replacement for datetime | ✅ Wraps time.Time | ✅ Complete | Well implemented |
| **Timezone Support** | ✅ Full IANA support | ✅ Full IANA support | ✅ Complete | Excellent DST handling |
| **Immutable Operations** | ✅ All operations immutable | ✅ All operations immutable | ✅ Complete | Thread-safe |
| **Fluent API** | ✅ Method chaining | ✅ Method chaining | ✅ Complete | Well designed |
| **Human-Readable Output** | ✅ "2 hours ago" format | ✅ "2 hours ago" format | ✅ Complete | Good localization |

### DateTime Creation & Instantiation

| Feature | Pendulum | chronogo | Status | Notes |
|---------|----------|----------|--------|-------|
| `now()` | ✅ | ✅ Now() | ✅ Complete | |
| `today()` | ✅ | ✅ Today() | ✅ Complete | |
| `yesterday()` | ✅ | ✅ Yesterday() | ✅ Complete | |
| `tomorrow()` | ✅ | ✅ Tomorrow() | ✅ Complete | |
| `from_timestamp()` | ✅ | ✅ FromUnix() | ✅ Complete | |
| `from_format()` | ✅ | ✅ FromFormat() | ✅ Complete | |
| `instance()` | ✅ Creates from another datetime | ✅ Instance() | ✅ Complete | |
| **Naive/Aware detection** | ✅ | ⚠️ Partial | ⚠️ Missing | No explicit naive datetime support |

### Parsing

| Feature | Pendulum | chronogo | Status | Notes |
|---------|----------|----------|--------|-------|
| ISO 8601 | ✅ | ✅ | ✅ Complete | |
| RFC 3339 | ✅ | ✅ | ✅ Complete | |
| Common formats | ✅ | ✅ | ✅ Complete | |
| Auto-detection | ✅ | ✅ | ✅ Complete | Good implementation |
| Custom formats | ✅ | ✅ | ✅ Complete | |
| **Strict mode** | ✅ | ✅ ParseStrict() | ✅ Complete | |
| **Ordinal dates** | ✅ YYYY-DDD | ❌ | ❌ Missing | ISO 8601 ordinal format |
| **Week dates** | ✅ YYYY-Www-D | ❌ | ❌ Missing | ISO 8601 week format |
| **Natural language** | ✅ "next monday" | ❌ | ❌ Missing | No natural language parsing |

### Localization

| Feature | Pendulum | chronogo | Status | Notes |
|---------|----------|----------|--------|-------|
| Multi-language support | ✅ | ✅ 6 locales | ✅ Complete | en, es, fr, de, zh, pt |
| Localized formatting | ✅ | ✅ | ✅ Complete | |
| Localized humanization | ✅ | ✅ | ✅ Complete | |
| Month names | ✅ | ✅ | ✅ Complete | |
| Weekday names | ✅ | ✅ | ✅ Complete | |
| Ordinals | ✅ | ✅ | ✅ Complete | 1st, 2nd, etc. |
| **Extended locale count** | ✅ 80+ locales | ⚠️ 6 locales | ⚠️ Limited | Pendulum has many more |

### Timezone Operations

| Feature | Pendulum | chronogo | Status | Notes |
|---------|----------|----------|--------|-------|
| Conversion | ✅ in_timezone() | ✅ InTimezone() | ✅ Complete | |
| UTC conversion | ✅ in_utc() | ✅ UTC() | ✅ Complete | |
| DST detection | ✅ is_dst() | ✅ IsDST() | ✅ Complete | Optimized with caching |
| Offset information | ✅ | ✅ | ✅ Complete | |
| **Timezone listing** | ✅ timezones | ✅ AvailableTimezones() | ✅ Complete | |
| **Timezone validation** | ✅ | ✅ IsValidTimezone() | ✅ Complete | |
| **Local timezone detection** | ✅ | ✅ | ✅ Complete | |

### Arithmetic & Manipulation

| Feature | Pendulum | chronogo | Status | Notes |
|---------|----------|----------|--------|-------|
| Add years/months/days | ✅ | ✅ | ✅ Complete | |
| Add hours/minutes/seconds | ✅ | ✅ | ✅ Complete | |
| Subtract time units | ✅ | ✅ | ✅ Complete | |
| **next(day_of_week)** | ✅ next(pendulum.MONDAY) | ✅ NextWeekday() | ✅ Complete | Get next occurrence of weekday |
| **previous(day_of_week)** | ✅ previous(pendulum.FRIDAY) | ✅ PreviousWeekday() | ✅ Complete | Get previous occurrence of weekday |
| **closest(day_of_week)** | ✅ | ✅ ClosestWeekday() | ✅ Complete | Get closest occurrence |
| **furthest(day_of_week)** | ✅ | ✅ FarthestWeekday() | ✅ Complete | Get furthest occurrence |
| **next_or_same()** | ❌ | ✅ NextOrSameWeekday() | ✅ chronogo Extra | Next or current if same |
| **previous_or_same()** | ❌ | ✅ PreviousOrSameWeekday() | ✅ chronogo Extra | Previous or current if same |
| Start of period | ✅ start_of() | ✅ StartOf...() | ✅ Complete | Day, week, month, year, quarter |
| End of period | ✅ end_of() | ✅ EndOf...() | ✅ Complete | Day, week, month, year, quarter |

### Modifiers

| Feature | Pendulum | chronogo | Status | Notes |
|---------|----------|----------|--------|-------|
| **set(component)** | ✅ | ✅ Set() fluent | ✅ Complete | |
| **on(year, month, day)** | ✅ | ⚠️ | ⚠️ Partial | Have Set().Year().Month().Day() |
| **at(hour, minute, second)** | ✅ | ⚠️ | ⚠️ Partial | Have Set().Hour().Minute() |
| **first_of(unit)** | ✅ first_of('month') | ✅ StartOfMonth() | ✅ Complete | |
| **last_of(unit)** | ✅ last_of('month') | ✅ EndOfMonth() | ✅ Complete | |
| **nth_of(unit, nth, day)** | ✅ nth_of('month', 2, pendulum.MONDAY) | ✅ NthWeekdayOf() | ✅ Complete | Get nth occurrence in period |
| **first_of(unit, day)** | ⚠️ Partial | ✅ FirstWeekdayOf() | ✅ Complete | First occurrence of weekday |
| **last_of(unit, day)** | ⚠️ Partial | ✅ LastWeekdayOf() | ✅ Complete | Last occurrence of weekday |
| **nth_of_month()** | ❌ | ✅ NthWeekdayOfMonth() | ✅ chronogo Extra | Convenience for month |
| **nth_of_year()** | ❌ | ✅ NthWeekdayOfYear() | ✅ chronogo Extra | Convenience for year |
| **is_nth_of()** | ❌ | ✅ IsNthWeekdayOf() | ✅ chronogo Extra | Check if nth occurrence |
| **weekday_occurrence()** | ❌ | ✅ WeekdayOccurrenceInMonth() | ✅ chronogo Extra | Get occurrence number |
| **average(dt)** | ✅ | ❌ | ❌ Missing | Get datetime between two points |

### Comparison

| Feature | Pendulum | chronogo | Status | Notes |
|---------|----------|----------|--------|-------|
| Before/After | ✅ | ✅ | ✅ Complete | |
| Equal | ✅ | ✅ | ✅ Complete | |
| Between | ✅ | ✅ Between() | ✅ Complete | |
| **is_birthday()** | ✅ | ❌ | ❌ Missing | Check if same month/day |
| **is_same_day()** | ✅ | ⚠️ | ⚠️ Partial | Can implement with comparison |
| **closest()/farthest()** | ✅ | ❌ | ❌ Missing | Find closest/farthest from list |
| Min/Max | ✅ | ✅ | ✅ Complete | |

### Difference

| Feature | Pendulum | chronogo | Status | Notes |
|---------|----------|----------|--------|-------|
| diff() method | ✅ | ⚠️ | ⚠️ Partial | Have Sub() but not full diff |
| **diff_for_humans()** | ✅ | ✅ DiffForHumans() | ✅ Complete | Excellent implementation |
| **in_words()** | ✅ | ✅ HumanString() | ✅ Complete | |
| **diff.in_days()** | ✅ | ⚠️ | ⚠️ Partial | Period has Days() |
| **diff.in_hours()** | ✅ | ⚠️ | ⚠️ Partial | Can get from Duration |
| **absolute differences** | ✅ diff(absolute=True) | ⚠️ | ⚠️ Partial | Have to calculate manually |

### Duration Type

| Feature | Pendulum | chronogo | Status | Notes |
|---------|----------|----------|--------|-------|
| Duration class | ✅ | ✅ ChronoDuration | ✅ Complete | |
| **years property** | ✅ duration.years | ✅ Years() | ✅ Complete | |
| **months property** | ✅ duration.months | ✅ Months() | ✅ Complete | |
| **weeks property** | ✅ duration.weeks | ✅ Weeks() | ✅ Complete | |
| **days property** | ✅ duration.days | ✅ Days() | ✅ Complete | |
| **hours property** | ✅ duration.hours | ✅ Hours() | ✅ Complete | |
| **in_weeks()** | ✅ | ✅ Weeks() | ✅ Complete | |
| **in_days()** | ✅ | ✅ Days() | ✅ Complete | |
| **in_hours()** | ✅ | ✅ Hours() | ✅ Complete | |
| **in_words()** | ✅ | ✅ HumanString() | ✅ Complete | |
| ISO 8601 parsing | ✅ P1Y2M3D | ✅ ParseISODuration() | ✅ Complete | |

### Period/Interval Type

| Feature | Pendulum | chronogo | Status | Notes |
|---------|----------|----------|--------|-------|
| Period class | ✅ | ✅ Period | ✅ Complete | |
| Iteration | ✅ for dt in period | ✅ Range() | ✅ Complete | Excellent with context |
| **range(unit, amount)** | ✅ period.range('days') | ✅ RangeDays(), RangeHours() | ✅ Complete | |
| Contains check | ✅ dt in period | ✅ Contains() | ✅ Complete | |
| Years/Months/Days | ✅ | ✅ | ✅ Complete | |
| **overlaps(period)** | ✅ | ❌ | ❌ Missing | Check if two periods overlap |
| **as_interval()** | ✅ | ❌ | ❌ Missing | Convert to ISO 8601 interval |

### Attributes & Properties

| Feature | Pendulum | chronogo | Status | Notes |
|---------|----------|----------|--------|-------|
| Year, month, day | ✅ | ✅ | ✅ Complete | |
| Hour, minute, second | ✅ | ✅ | ✅ Complete | |
| Microsecond | ✅ | ✅ Nanosecond() | ✅ Complete | |
| Day of week | ✅ day_of_week | ✅ Weekday() | ✅ Complete | |
| Day of year | ✅ day_of_year | ✅ DayOfYear() | ✅ Complete | |
| Week of year | ✅ week_of_year | ✅ ISOWeek() | ✅ Complete | |
| Week of month | ✅ week_of_month | ✅ WeekOfMonth() | ✅ Complete | |
| Days in month | ✅ days_in_month | ✅ DaysInMonth() | ✅ Complete | |
| Quarter | ✅ quarter | ✅ Quarter() | ✅ Complete | |
| **age** | ✅ age (returns years as int) | ✅ Age() | ✅ Complete | |
| **is_leap_year()** | ✅ | ✅ IsLeapYear() | ✅ Complete | |
| **is_long_year()** | ✅ ISO week year with 53 weeks | ❌ | ❌ Missing | ISO 8601 long year check |
| **is_same_day(dt)** | ✅ | ⚠️ | ⚠️ Partial | Can compare dates |
| **is_anniversary(dt)** | ✅ | ❌ | ❌ Missing | Check if same month/day |

### String Formatting

| Feature | Pendulum | chronogo | Status | Notes |
|---------|----------|----------|--------|-------|
| format() | ✅ | ✅ Format() | ✅ Complete | Go layout format |
| to_iso8601_string() | ✅ | ✅ ToISO8601() | ✅ Complete | |
| to_rfc3339_string() | ✅ | ✅ ToRFC3339() | ✅ Complete | |
| to_datetime_string() | ✅ | ✅ ToDateTimeString() | ✅ Complete | |
| to_date_string() | ✅ | ✅ ToDateString() | ✅ Complete | |
| to_time_string() | ✅ | ✅ ToTimeString() | ✅ Complete | |
| **to_cookie_string()** | ✅ | ✅ ToCookieString() | ✅ Complete | HTTP cookie format (RFC1123) |
| **to_rss_string()** | ✅ | ✅ ToRSSString() | ✅ Complete | RSS feed format (RFC1123Z) |
| **to_w3c_string()** | ✅ | ✅ ToW3CString() | ✅ Complete | W3C datetime format (RFC3339) |
| **to_atom_string()** | ⚠️ | ✅ ToAtomString() | ✅ Complete | Atom feed format |

### Testing Helpers

| Feature | Pendulum | chronogo | Status | Notes |
|---------|----------|----------|--------|-------|
| **set_test_now()** | ✅ | ✅ SetTestNow() | ✅ Complete | Mock current time for testing |
| **travel_to(dt)** | ✅ | ✅ TravelTo() | ✅ Complete | Time travel for testing |
| **travel()** | ✅ Context manager for tests | ✅ WithTestNow() | ✅ Complete | Test time manipulation |
| **freeze()** | ✅ Freeze time at specific point | ✅ FreezeTime() | ✅ Complete | Stop time for testing |
| **travel_back()** | ❌ | ✅ TravelBack() | ✅ chronogo Extra | Go backwards in time |
| **travel_forward()** | ❌ | ✅ TravelForward() | ✅ chronogo Extra | Go forward in time |

### Business Date Features (chronogo Enhancement)

| Feature | Pendulum | chronogo | Status | Notes |
|---------|----------|----------|--------|-------|
| Holiday checking | ❌ | ✅ | ✅ chronogo Extra | 34 countries via goholiday |
| Business day calculations | ❌ | ✅ | ✅ chronogo Extra | Advanced implementation |
| Holiday calendars | ❌ | ✅ | ✅ chronogo Extra | |
| Scheduling | ❌ | ✅ | ✅ chronogo Extra | Holiday-aware scheduling |

## Missing Critical Features

### High Priority (COMPLETED ✅)

1. ✅ **Testing Helpers** - SetTestNow(), TravelTo(), FreezeTime() IMPLEMENTED
2. ✅ **Weekday Navigation** - NextWeekday(), PreviousWeekday(), ClosestWeekday() IMPLEMENTED
3. ✅ **nth_of() Method** - NthWeekdayOf(), FirstWeekdayOf(), LastWeekdayOf() IMPLEMENTED

### High Priority (Remaining)

4. ✅ **is_birthday()/is_anniversary()/comparison methods** - IMPLEMENTED with extras
5. ✅ **Period.overlaps() and related methods** - IMPLEMENTED with Gap(), Encompasses(), Merge()
6. ✅ **Extended String Formats** - ToCookieString(), ToRSSString(), ToW3CString() IMPLEMENTED
7. **Natural Language Parsing** - Pendulum supports "next monday", "last week", etc. (TODO)
8. **Naive/Aware DateTime Distinction** - Pendulum explicitly differentiates timezone-aware and naive datetimes (TODO)

### Medium Priority

9. **ISO 8601 Week Dates** - YYYY-Www-D format parsing
10. **ISO 8601 Ordinal Dates** - YYYY-DDD format parsing
11. **is_long_year()** - ISO 8601 long year detection
12. **average()** - Get midpoint between two datetimes
13. **closest()/farthest()** - Find closest/farthest datetime from list
14. **More Locales** - Pendulum supports 80+ locales vs chronogo's 6

### Low Priority

15. **on()/at() Convenience Methods** - Simpler alternatives to Set()
16. **Explicit Diff Type** - Richer difference object with more methods

## Recommendations

### Quick Wins (✅ ALL COMPLETED)

1. ✅ DONE: `IsBirthday(dt DateTime) bool` method
2. ✅ DONE: `IsAnniversary(dt DateTime) bool` method  
3. ✅ DONE: `IsSameDay(dt DateTime) bool` method + IsSameMonth, IsSameYear, IsSameQuarter, IsSameWeek
4. ✅ DONE: `Average(dt DateTime) DateTime` method
5. ✅ DONE: `ToCookieString()`, `ToRSSString()`, `ToW3CString()`, `ToAtomString()` methods
6. ✅ DONE: Period `Overlaps(p Period) bool` method + Gap, Encompasses, Merge
7. ✅ DONE: `Closest(dates ...DateTime)` and `Farthest(dates ...DateTime)` methods

### Medium Effort (Moderate Implementation)

7. ✅ DONE: Weekday navigation implemented
8. ✅ DONE: NthWeekdayOf implemented with extensive features
9. Add ISO 8601 week date parsing support
10. Add ISO 8601 ordinal date parsing support
11. ✅ DONE: Testing helpers fully implemented

### High Effort (Complex Implementation)

12. Natural language date parsing ("next monday", "last week", "in 3 days")
13. Explicit naive/aware datetime distinction with conversion methods
14. Expand localization to 20+ most common locales
15. Add explicit Diff type with comprehensive methods

## Conclusion

chronogo has implemented the core features of Pendulum very well and has added significant value with business date operations that Pendulum lacks. The main gaps are:

1. **Testing utilities** - Essential for developers
2. **Weekday navigation helpers** - Very common use case
3. **Natural language parsing** - Nice to have but complex
4. **More locales** - Good for international applications

chronogo's business date features (holiday checking, business day calculations, scheduling) are a significant advantage over Pendulum and provide unique value to Go developers.

## References

- Pendulum Documentation: [https://pendulum.eustace.io/](https://pendulum.eustace.io/)
- chronogo Repository: github.com/coredds/chronogo

