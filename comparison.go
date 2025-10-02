package chronogo

// IsBirthday checks if the given DateTime represents the same birthday (month and day).
// This is useful for checking if a date is someone's birthday, regardless of the year.
//
// Example:
//   birthday := chronogo.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)
//   today := chronogo.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC)
//   today.IsBirthday(birthday) // Returns true
func (dt DateTime) IsBirthday(other DateTime) bool {
	return dt.Month() == other.Month() && dt.Day() == other.Day()
}

// IsAnniversary checks if the given DateTime represents the same anniversary (month and day).
// This is an alias for IsBirthday and can be used for any recurring annual event.
//
// Example:
//   anniversary := chronogo.Date(2010, 6, 20, 0, 0, 0, 0, time.UTC)
//   today := chronogo.Date(2024, 6, 20, 10, 0, 0, 0, time.UTC)
//   today.IsAnniversary(anniversary) // Returns true
func (dt DateTime) IsAnniversary(other DateTime) bool {
	return dt.IsBirthday(other)
}

// IsSameDay checks if the given DateTime is on the same calendar day.
// This compares year, month, and day, ignoring time components.
//
// Example:
//   dt1 := chronogo.Date(2024, 5, 15, 8, 0, 0, 0, time.UTC)
//   dt2 := chronogo.Date(2024, 5, 15, 20, 0, 0, 0, time.UTC)
//   dt1.IsSameDay(dt2) // Returns true
func (dt DateTime) IsSameDay(other DateTime) bool {
	return dt.Year() == other.Year() &&
		dt.Month() == other.Month() &&
		dt.Day() == other.Day()
}

// IsSameMonth checks if the given DateTime is in the same month and year.
//
// Example:
//   dt1 := chronogo.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
//   dt2 := chronogo.Date(2024, 5, 31, 0, 0, 0, 0, time.UTC)
//   dt1.IsSameMonth(other) // Returns true
func (dt DateTime) IsSameMonth(other DateTime) bool {
	return dt.Year() == other.Year() && dt.Month() == other.Month()
}

// IsSameYear checks if the given DateTime is in the same year.
//
// Example:
//   dt1 := chronogo.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
//   dt2 := chronogo.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
//   dt1.IsSameYear(other) // Returns true
func (dt DateTime) IsSameYear(other DateTime) bool {
	return dt.Year() == other.Year()
}

// Average returns the DateTime that is exactly halfway between this DateTime and another.
// This is useful for finding the midpoint between two dates.
//
// Example:
//   start := chronogo.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
//   end := chronogo.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
//   midpoint := start.Average(end) // Returns 2024-01-16
func (dt DateTime) Average(other DateTime) DateTime {
	// Calculate the duration between the two times
	duration := other.Time.Sub(dt.Time)
	
	// Add half the duration to the earlier time
	midpoint := dt.Add(duration / 2)
	
	return midpoint
}

// Closest returns the closest DateTime from a list of DateTimes.
// Returns zero DateTime if the list is empty.
//
// Example:
//   dt := chronogo.Now()
//   dates := []chronogo.DateTime{
//       chronogo.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
//       chronogo.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
//       chronogo.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
//   }
//   closest := dt.Closest(dates...)
func (dt DateTime) Closest(dates ...DateTime) DateTime {
	if len(dates) == 0 {
		return DateTime{}
	}
	
	closest := dates[0]
	minDuration := dt.Time.Sub(dates[0].Time)
	if minDuration < 0 {
		minDuration = -minDuration
	}
	
	for i := 1; i < len(dates); i++ {
		duration := dt.Time.Sub(dates[i].Time)
		if duration < 0 {
			duration = -duration
		}
		
		if duration < minDuration {
			minDuration = duration
			closest = dates[i]
		}
	}
	
	return closest
}

// Farthest returns the farthest DateTime from a list of DateTimes.
// Returns zero DateTime if the list is empty.
//
// Example:
//   dt := chronogo.Now()
//   dates := []chronogo.DateTime{
//       chronogo.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
//       chronogo.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
//       chronogo.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
//   }
//   farthest := dt.Farthest(dates...)
func (dt DateTime) Farthest(dates ...DateTime) DateTime {
	if len(dates) == 0 {
		return DateTime{}
	}
	
	farthest := dates[0]
	maxDuration := dt.Time.Sub(dates[0].Time)
	if maxDuration < 0 {
		maxDuration = -maxDuration
	}
	
	for i := 1; i < len(dates); i++ {
		duration := dt.Time.Sub(dates[i].Time)
		if duration < 0 {
			duration = -duration
		}
		
		if duration > maxDuration {
			maxDuration = duration
			farthest = dates[i]
		}
	}
	
	return farthest
}

// IsSameQuarter checks if the given DateTime is in the same quarter and year.
//
// Example:
//   dt1 := chronogo.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC) // Q1
//   dt2 := chronogo.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC) // Q1
//   dt1.IsSameQuarter(dt2) // Returns true
func (dt DateTime) IsSameQuarter(other DateTime) bool {
	return dt.Year() == other.Year() && dt.Quarter() == other.Quarter()
}

// IsSameWeek checks if the given DateTime is in the same ISO week and year.
//
// Example:
//   dt1 := chronogo.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
//   dt2 := chronogo.Date(2024, 1, 17, 0, 0, 0, 0, time.UTC)
//   dt1.IsSameWeek(dt2) // Returns true if in same ISO week
func (dt DateTime) IsSameWeek(other DateTime) bool {
	year1, week1 := dt.ISOWeek()
	year2, week2 := other.ISOWeek()
	return year1 == year2 && week1 == week2
}

