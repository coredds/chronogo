package chronogo

import (
	"errors"
	"math"
	"regexp"
	"strconv"
	"time"
)

// ParseISODuration parses an ISO 8601 duration string (e.g., "P1Y2M3DT4H5M6S", "PT15M", "P2W").
//
// Notes:
//   - Years and months are approximated using the same factors as ChronoDuration methods:
//     1 year = 365.25 days, 1 month = 30.44 days.
//   - Weeks are converted as 7 days.
//   - A leading minus sign is supported for negative durations.
func ParseISODuration(s string) (ChronoDuration, error) {
	// ^([+-])?P(?:(\d+)Y)?(?:(\d+)M)?(?:(\d+)W)?(?:(\d+)D)?(?:T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+(?:\.\d+)?)S)?)?$
	isoDurRe := regexp.MustCompile(`^([+-])?P(?:(\d+)Y)?(?:(\d+)M)?(?:(\d+)W)?(?:(\d+)D)?(?:T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+(?:\.\d+)?)S)?)?$`)
	m := isoDurRe.FindStringSubmatch(s)
	if m == nil {
		return ChronoDuration{}, errors.New("invalid ISO 8601 duration")
	}

	sign := 1.0
	if m[1] == "-" {
		sign = -1.0
	}

	parseInt := func(idx int) int64 {
		if idx >= len(m) || m[idx] == "" {
			return 0
		}
		v, _ := strconv.ParseInt(m[idx], 10, 64)
		return v
	}
	parseFloat := func(idx int) float64 {
		if idx >= len(m) || m[idx] == "" {
			return 0
		}
		v, _ := strconv.ParseFloat(m[idx], 64)
		return v
	}

	years := parseInt(2)
	months := parseInt(3)
	weeks := parseInt(4)
	days := parseInt(5)
	hours := parseInt(6)
	minutes := parseInt(7)
	seconds := parseFloat(8)

	// Convert everything to time.Duration using approximations for years/months
	totalSeconds := 0.0
	totalSeconds += float64(weeks*7*24*3600 + days*24*3600 + hours*3600 + minutes*60)
	totalSeconds += seconds
	totalSeconds += float64(months) * 30.44 * 24 * 3600
	totalSeconds += float64(years) * 365.25 * 24 * 3600

	d := time.Duration(sign * totalSeconds * float64(time.Second))
	// Normalize to avoid rounding to zero for tiny negatives due to float math
	if d == 0 && totalSeconds != 0 {
		if sign < 0 {
			d = -time.Nanosecond
		} else {
			d = time.Nanosecond
		}
	}
	// Round to the nearest nanosecond
	d = time.Duration(math.Round(float64(d)))
	return ChronoDuration{d}, nil
}
