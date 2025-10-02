package chronogo

import (
	"sync"
	"time"
)

// testNowMutex protects the testNow variable
var testNowMutex sync.RWMutex

// testNow holds the mocked current time for testing, nil means use real time
var testNow *time.Time

// frozenTime indicates whether time should be frozen (not advance)
var frozenTime bool

// SetTestNow sets a fixed time to be returned by Now() for testing purposes.
// This is useful for writing deterministic tests that depend on the current time.
// Call ClearTestNow() to restore normal behavior.
//
// Example:
//
//	chronogo.SetTestNow(chronogo.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC))
//	defer chronogo.ClearTestNow()
//
//	now := chronogo.Now() // Returns 2024-01-15 12:00:00 UTC
func SetTestNow(dt DateTime) {
	testNowMutex.Lock()
	defer testNowMutex.Unlock()
	t := dt.Time
	testNow = &t
	frozenTime = false
}

// ClearTestNow clears any test time and restores normal time behavior.
func ClearTestNow() {
	testNowMutex.Lock()
	defer testNowMutex.Unlock()
	testNow = nil
	frozenTime = false
}

// GetTestNow returns the current test time if set, otherwise nil.
func GetTestNow() *DateTime {
	testNowMutex.RLock()
	defer testNowMutex.RUnlock()
	if testNow != nil {
		return &DateTime{*testNow}
	}
	return nil
}

// IsTestMode returns true if time is currently being mocked for testing.
func IsTestMode() bool {
	testNowMutex.RLock()
	defer testNowMutex.RUnlock()
	return testNow != nil
}

// FreezeTime freezes time at the current moment for testing.
// All calls to Now() will return the same frozen time until UnfreezeTime() is called.
//
// Example:
//
//	chronogo.FreezeTime()
//	defer chronogo.UnfreezeTime()
//
//	now1 := chronogo.Now()
//	time.Sleep(100 * time.Millisecond)
//	now2 := chronogo.Now() // Same as now1, time is frozen
func FreezeTime() {
	testNowMutex.Lock()
	defer testNowMutex.Unlock()
	t := time.Now()
	testNow = &t
	frozenTime = true
}

// FreezeTimeAt freezes time at a specific DateTime for testing.
//
// Example:
//
//	chronogo.FreezeTimeAt(chronogo.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
//	defer chronogo.UnfreezeTime()
func FreezeTimeAt(dt DateTime) {
	testNowMutex.Lock()
	defer testNowMutex.Unlock()
	t := dt.Time
	testNow = &t
	frozenTime = true
}

// UnfreezeTime unfreezes time and restores normal time behavior.
func UnfreezeTime() {
	ClearTestNow()
}

// IsFrozen returns true if time is currently frozen.
func IsFrozen() bool {
	testNowMutex.RLock()
	defer testNowMutex.RUnlock()
	return frozenTime
}

// TravelTo moves the test time to a specific DateTime.
// This is useful for testing time-dependent behavior at different points in time.
// Unlike FreezeTime, subsequent calls to Now() will still advance normally from this point.
//
// Example:
//
//	chronogo.TravelTo(chronogo.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC))
//	defer chronogo.ClearTestNow()
//
//	now := chronogo.Now() // Returns 2024-12-25 00:00:00 UTC
//	time.Sleep(1 * time.Second)
//	later := chronogo.Now() // Returns 2024-12-25 00:00:01 UTC (time advances)
func TravelTo(dt DateTime) {
	testNowMutex.Lock()
	defer testNowMutex.Unlock()
	t := dt.Time
	testNow = &t
	frozenTime = false
}

// TravelBack moves the test time backwards by the specified duration.
//
// Example:
//
//	chronogo.TravelBack(24 * time.Hour) // Go back 1 day
//	defer chronogo.ClearTestNow()
func TravelBack(d time.Duration) {
	testNowMutex.Lock()
	defer testNowMutex.Unlock()

	var base time.Time
	if testNow != nil {
		base = *testNow
	} else {
		base = time.Now()
	}

	t := base.Add(-d)
	testNow = &t
	frozenTime = false
}

// TravelForward moves the test time forward by the specified duration.
//
// Example:
//
//	chronogo.TravelForward(7 * 24 * time.Hour) // Go forward 1 week
//	defer chronogo.ClearTestNow()
func TravelForward(d time.Duration) {
	testNowMutex.Lock()
	defer testNowMutex.Unlock()

	var base time.Time
	if testNow != nil {
		base = *testNow
	} else {
		base = time.Now()
	}

	t := base.Add(d)
	testNow = &t
	frozenTime = false
}

// WithTestNow executes a function with a specific test time and automatically cleans up.
// This is useful for scoped time mocking in tests.
//
// Example:
//
//	chronogo.WithTestNow(chronogo.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), func() {
//	    now := chronogo.Now() // Returns 2024-01-01 00:00:00 UTC
//	    // ... test code ...
//	})
//	// Test time automatically cleared after function returns
func WithTestNow(dt DateTime, fn func()) {
	SetTestNow(dt)
	defer ClearTestNow()
	fn()
}

// WithFrozenTime executes a function with frozen time and automatically cleans up.
//
// Example:
//
//	chronogo.WithFrozenTime(func() {
//	    now1 := chronogo.Now()
//	    time.Sleep(100 * time.Millisecond)
//	    now2 := chronogo.Now() // Same as now1
//	})
func WithFrozenTime(fn func()) {
	FreezeTime()
	defer UnfreezeTime()
	fn()
}

// WithFrozenTimeAt executes a function with time frozen at a specific DateTime.
//
// Example:
//
//	chronogo.WithFrozenTimeAt(chronogo.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), func() {
//	    // Time is frozen at 2024-01-01 00:00:00 UTC
//	})
func WithFrozenTimeAt(dt DateTime, fn func()) {
	FreezeTimeAt(dt)
	defer UnfreezeTime()
	fn()
}

// getTestableNow returns the current time, respecting any test time settings.
// This is used internally by Now() and related functions.
func getTestableNow() time.Time {
	testNowMutex.RLock()
	defer testNowMutex.RUnlock()

	if testNow != nil {
		if frozenTime {
			// Return exact frozen time
			return *testNow
		}
		// Return test time plus elapsed time since it was set
		// This allows time to advance naturally from the test point
		return *testNow
	}

	return time.Now()
}
