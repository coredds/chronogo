package chronogo

import (
	"testing"
	"time"
)

func TestSetTestNow(t *testing.T) {
	// Clean up after test
	defer ClearTestNow()

	testTime := Date(2024, time.January, 15, 12, 0, 0, 0, time.UTC)
	SetTestNow(testTime)

	now := Now()
	if !now.Equal(testTime) {
		t.Errorf("Expected Now() to return %v, got %v", testTime, now)
	}

	// Verify NowUTC also respects test time
	nowUTC := NowUTC()
	if !nowUTC.Equal(testTime) {
		t.Errorf("Expected NowUTC() to return %v, got %v", testTime, nowUTC)
	}
}

func TestClearTestNow(t *testing.T) {
	testTime := Date(2024, time.January, 15, 12, 0, 0, 0, time.UTC)
	SetTestNow(testTime)

	now1 := Now()
	if !now1.Equal(testTime) {
		t.Errorf("Expected Now() to return test time")
	}

	ClearTestNow()

	now2 := Now()
	// After clearing, we should get real time (which will be different from test time)
	if now2.Equal(testTime) {
		t.Errorf("Expected Now() to return real time after ClearTestNow()")
	}
}

func TestFreezeTime(t *testing.T) {
	defer UnfreezeTime()

	FreezeTime()

	now1 := Now()
	time.Sleep(10 * time.Millisecond)
	now2 := Now()

	if !now1.Equal(now2) {
		t.Errorf("Expected frozen time to remain constant. Got %v and %v", now1, now2)
	}

	if !IsFrozen() {
		t.Error("Expected IsFrozen() to return true")
	}
}

func TestFreezeTimeAt(t *testing.T) {
	defer UnfreezeTime()

	frozenTime := Date(2024, time.December, 25, 0, 0, 0, 0, time.UTC)
	FreezeTimeAt(frozenTime)

	now := Now()
	if !now.Equal(frozenTime) {
		t.Errorf("Expected Now() to return %v, got %v", frozenTime, now)
	}

	time.Sleep(10 * time.Millisecond)
	now2 := Now()
	if !now2.Equal(frozenTime) {
		t.Errorf("Expected frozen time to remain constant")
	}
}

func TestUnfreezeTime(t *testing.T) {
	FreezeTime()
	if !IsFrozen() {
		t.Error("Expected time to be frozen")
	}

	UnfreezeTime()
	if IsFrozen() {
		t.Error("Expected time to be unfrozen")
	}

	if IsTestMode() {
		t.Error("Expected test mode to be disabled after UnfreezeTime()")
	}
}

func TestTravelTo(t *testing.T) {
	defer ClearTestNow()

	destination := Date(2024, time.July, 4, 0, 0, 0, 0, time.UTC)
	TravelTo(destination)

	now := Now()
	if !now.Equal(destination) {
		t.Errorf("Expected Now() to return %v, got %v", destination, now)
	}

	if IsFrozen() {
		t.Error("TravelTo should not freeze time")
	}
}

func TestTravelBack(t *testing.T) {
	defer ClearTestNow()

	// Start at a known time
	start := Date(2024, time.January, 15, 12, 0, 0, 0, time.UTC)
	SetTestNow(start)

	// Travel back 7 days
	TravelBack(7 * 24 * time.Hour)

	now := Now()
	expected := start.AddDays(-7)

	if !now.Equal(expected) {
		t.Errorf("Expected Now() to return %v, got %v", expected, now)
	}
}

func TestTravelForward(t *testing.T) {
	defer ClearTestNow()

	// Start at a known time
	start := Date(2024, time.January, 15, 12, 0, 0, 0, time.UTC)
	SetTestNow(start)

	// Travel forward 30 days
	TravelForward(30 * 24 * time.Hour)

	now := Now()
	expected := start.AddDays(30)

	if !now.Equal(expected) {
		t.Errorf("Expected Now() to return %v, got %v", expected, now)
	}
}

func TestIsTestMode(t *testing.T) {
	// Initially not in test mode
	if IsTestMode() {
		t.Error("Expected IsTestMode() to return false initially")
	}

	// Enter test mode
	SetTestNow(Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC))
	if !IsTestMode() {
		t.Error("Expected IsTestMode() to return true after SetTestNow()")
	}

	// Exit test mode
	ClearTestNow()
	if IsTestMode() {
		t.Error("Expected IsTestMode() to return false after ClearTestNow()")
	}
}

func TestGetTestNow(t *testing.T) {
	defer ClearTestNow()

	// No test time set
	if GetTestNow() != nil {
		t.Error("Expected GetTestNow() to return nil when no test time is set")
	}

	// Set test time
	testTime := Date(2024, time.March, 15, 10, 30, 0, 0, time.UTC)
	SetTestNow(testTime)

	retrievedTime := GetTestNow()
	if retrievedTime == nil {
		t.Fatal("Expected GetTestNow() to return a value")
	}

	if !retrievedTime.Equal(testTime) {
		t.Errorf("Expected GetTestNow() to return %v, got %v", testTime, *retrievedTime)
	}
}

func TestWithTestNow(t *testing.T) {
	testTime := Date(2024, time.June, 1, 0, 0, 0, 0, time.UTC)

	WithTestNow(testTime, func() {
		now := Now()
		if !now.Equal(testTime) {
			t.Errorf("Expected Now() to return %v inside WithTestNow, got %v", testTime, now)
		}

		if !IsTestMode() {
			t.Error("Expected IsTestMode() to return true inside WithTestNow")
		}
	})

	// After WithTestNow, test mode should be cleared
	if IsTestMode() {
		t.Error("Expected test mode to be cleared after WithTestNow")
	}
}

func TestWithFrozenTime(t *testing.T) {
	WithFrozenTime(func() {
		now1 := Now()
		time.Sleep(10 * time.Millisecond)
		now2 := Now()

		if !now1.Equal(now2) {
			t.Error("Expected time to be frozen inside WithFrozenTime")
		}

		if !IsFrozen() {
			t.Error("Expected IsFrozen() to return true inside WithFrozenTime")
		}
	})

	// After WithFrozenTime, time should be unfrozen
	if IsFrozen() {
		t.Error("Expected time to be unfrozen after WithFrozenTime")
	}
}

func TestWithFrozenTimeAt(t *testing.T) {
	frozenTime := Date(2024, time.October, 31, 23, 59, 59, 0, time.UTC)

	WithFrozenTimeAt(frozenTime, func() {
		now := Now()
		if !now.Equal(frozenTime) {
			t.Errorf("Expected Now() to return %v, got %v", frozenTime, now)
		}

		time.Sleep(10 * time.Millisecond)
		now2 := Now()
		if !now2.Equal(frozenTime) {
			t.Error("Expected time to remain frozen")
		}
	})

	// Time should be unfrozen after
	if IsFrozen() {
		t.Error("Expected time to be unfrozen after WithFrozenTimeAt")
	}
}

func TestNowInWithTestTime(t *testing.T) {
	defer ClearTestNow()

	testTime := Date(2024, time.January, 15, 12, 0, 0, 0, time.UTC)
	SetTestNow(testTime)

	ny, _ := time.LoadLocation("America/New_York")
	nowInNY := NowIn(ny)

	// Should be the same instant, just different timezone
	if !nowInNY.Equal(testTime) {
		t.Errorf("Expected NowIn() to return equivalent time in different timezone")
	}

	// But different wall clock time
	if nowInNY.Hour() == testTime.Hour() {
		// They might be equal if UTC offset is 0, but in general should differ
		t.Logf("Note: Hour may differ based on timezone offset")
	}
}

func TestTodayWithTestTime(t *testing.T) {
	defer ClearTestNow()

	testTime := Date(2024, time.January, 15, 23, 45, 30, 0, time.UTC)
	SetTestNow(testTime)

	today := Today()
	expected := Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC)

	if !today.Equal(expected) {
		t.Errorf("Expected Today() to return %v, got %v", expected, today)
	}
}

func TestConcurrentTestNowAccess(t *testing.T) {
	defer ClearTestNow()

	testTime := Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	SetTestNow(testTime)

	// Run multiple goroutines accessing Now() concurrently
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				_ = Now()
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should not panic or race
	t.Log("Concurrent access test passed")
}

