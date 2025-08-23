package chronogo

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input    string
		expected string // Expected ISO format
		hasError bool
	}{
		{"2023-12-25T15:30:45Z", "2023-12-25T15:30:45Z", false},
		{"2023-12-25T15:30:45+00:00", "2023-12-25T15:30:45Z", false},
		{"2023-12-25 15:30:45", "2023-12-25T15:30:45Z", false},
		{"2023-12-25", "2023-12-25T00:00:00Z", false},
		{"2023/12/25", "2023-12-25T00:00:00Z", false},
		{"2023-1-2 3:04:05", "2023-01-02T03:04:05Z", false},
		{"20231225", "2023-12-25T00:00:00Z", false},
		{"15:30:45", "", false}, // Time only, will have today's date
		{"invalid", "", true},
		{"", "", true},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			dt, err := Parse(test.input)

			if test.hasError {
				if err == nil {
					t.Errorf("Expected error for input '%s', but got none", test.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for input '%s': %v", test.input, err)
				return
			}

			if test.expected != "" && test.input != "15:30:45" { // Skip time-only test
				result := dt.UTC().ToISO8601String()
				if result != test.expected {
					t.Errorf("Parse('%s'): expected %s, got %s", test.input, test.expected, result)
				}
			}
		})
	}
}

func TestParseStrict(t *testing.T) {
	valid := []string{
		"2023-12-25T15:30:45Z",
		"2023-12-25T15:30:45",
	}
	for _, s := range valid {
		if _, err := ParseStrict(s); err != nil {
			t.Fatalf("ParseStrict failed for %q: %v", s, err)
		}
	}

	invalid := []string{"2023/12/25", "2023-12-25 15:30:45", "20231225"}
	for _, s := range invalid {
		if _, err := ParseStrict(s); err == nil {
			t.Fatalf("ParseStrict should reject %q", s)
		}
	}
}

func TestParseInLocation(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip("Could not load America/New_York timezone")
	}

	dt, err := ParseInLocation("2023-12-25 15:30:45", loc)
	if err != nil {
		t.Fatalf("ParseInLocation failed: %v", err)
	}

	if dt.Location() != loc {
		t.Errorf("ParseInLocation should use specified location")
	}
}

func TestParseISO8601(t *testing.T) {
	tests := []struct {
		input     string
		expectErr bool
	}{
		{"2023-12-25T15:30:45Z", false},
		{"2023-12-25T15:30:45+00:00", false},
		{"2023-12-25T15:30:45.123Z", false},
		{"2023-12-25T15:30:45.123456Z", false},
		{"2023-12-25 15:30:45", true}, // Not ISO 8601
		{"invalid", true},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			_, err := ParseISO8601(test.input)

			if test.expectErr && err == nil {
				t.Errorf("Expected error for input '%s'", test.input)
			}

			if !test.expectErr && err != nil {
				t.Errorf("Unexpected error for input '%s': %v", test.input, err)
			}
		})
	}
}

func TestParseRFC3339(t *testing.T) {
	input := "2023-12-25T15:30:45Z"
	dt, err := ParseRFC3339(input)
	if err != nil {
		t.Fatalf("ParseRFC3339 failed: %v", err)
	}

	expected := "2023-12-25T15:30:45Z"
	result := dt.ToISO8601String()
	if result != expected {
		t.Errorf("ParseRFC3339: expected %s, got %s", expected, result)
	}
}

func TestFromFormat(t *testing.T) {
	tests := []struct {
		input    string
		layout   string
		expected string
	}{
		{"25/12/2023 15:30", "02/01/2006 15:04", "2023-12-25T15:30:00Z"},
		{"2023-Dec-25", "2006-Jan-02", "2023-12-25T00:00:00Z"},
		{"15:30:45", "15:04:05", "0000-01-01T15:30:45Z"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			dt, err := FromFormat(test.input, test.layout)
			if err != nil {
				t.Fatalf("FromFormat failed: %v", err)
			}

			result := dt.ToISO8601String()
			if result != test.expected {
				t.Errorf("FromFormat('%s', '%s'): expected %s, got %s",
					test.input, test.layout, test.expected, result)
			}
		})
	}
}

func TestFromFormatInLocation(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip("Could not load America/New_York timezone")
	}

	dt, err := FromFormatInLocation("25/12/2023 15:30", "02/01/2006 15:04", loc)
	if err != nil {
		t.Fatalf("FromFormatInLocation failed: %v", err)
	}

	if dt.Location() != loc {
		t.Errorf("FromFormatInLocation should use specified location")
	}
}

func TestLoadLocation(t *testing.T) {
	// Test valid location
	loc, err := LoadLocation("America/New_York")
	if err != nil {
		t.Errorf("LoadLocation failed for valid timezone: %v", err)
	}
	if loc == nil {
		t.Errorf("LoadLocation returned nil for valid timezone")
	}

	// Test "local" special case
	loc, err = LoadLocation("local")
	if err != nil {
		t.Errorf("LoadLocation failed for 'local': %v", err)
	}
	if loc != time.Local {
		t.Errorf("LoadLocation('local') should return time.Local")
	}

	// Test invalid location
	_, err = LoadLocation("Invalid/Timezone")
	if err == nil {
		t.Errorf("LoadLocation should fail for invalid timezone")
	}
}

func TestInstance(t *testing.T) {
	stdTime := time.Date(2023, time.December, 25, 15, 30, 45, 0, time.UTC)
	dt := Instance(stdTime)

	if !dt.Time.Equal(stdTime) {
		t.Errorf("Instance should wrap the provided time.Time")
	}
}

func TestTryParseUnix(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"1640995200", "2022-01-01T00:00:00Z", false},          // Valid Unix timestamp
		{"1640995200000", "2022-01-01T00:00:00Z", false},       // Unix timestamp with milliseconds
		{"1640995200000000", "2022-01-01T00:00:00Z", false},    // Unix timestamp with microseconds
		{"1640995200000000000", "2022-01-01T00:00:00Z", false}, // Unix timestamp with nanoseconds
		{"abc", "", true},                     // Invalid
		{"99999999999999", "", true},          // Invalid length (14)
		{"-1", "1969-12-31T23:59:59Z", false}, // Negative seconds
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			dt, err := TryParseUnix(test.input)

			if test.hasError {
				if err == nil {
					t.Errorf("Expected error for input '%s'", test.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for input '%s': %v", test.input, err)
				return
			}

			result := dt.ToISO8601String()
			if result != test.expected {
				t.Errorf("TryParseUnix('%s'): expected %s, got %s", test.input, test.expected, result)
			}
		})
	}
}

func TestAvailableTimezones(t *testing.T) {
	timezones := AvailableTimezones()

	// Should return a non-empty slice
	if len(timezones) == 0 {
		t.Error("AvailableTimezones() should return non-empty slice")
	}

	// Should contain common timezones
	found := false
	for _, tz := range timezones {
		if tz == "UTC" {
			found = true
			break
		}
	}
	if !found {
		t.Error("AvailableTimezones() should include 'UTC'")
	}

	// Check that we have a reasonable number of timezones
	if len(timezones) < 10 {
		t.Errorf("Expected more timezones, got %d", len(timezones))
	}
}

func TestIsValidTimezone(t *testing.T) {
	tests := []struct {
		timezone string
		valid    bool
	}{
		{"UTC", true},
		{"America/New_York", true},
		{"Europe/London", true},
		{"Invalid/Timezone", false},
		{"Not_A_Real_Zone", false},
	}

	for _, test := range tests {
		t.Run(test.timezone, func(t *testing.T) {
			result := IsValidTimezone(test.timezone)
			if result != test.valid {
				t.Errorf("IsValidTimezone('%s'): expected %v, got %v", test.timezone, test.valid, result)
			}
		})
	}
}

func TestParseOptimized(t *testing.T) {
	tests := []struct {
		input    string
		expected string // Expected ISO format
		hasError bool
	}{
		{"2023-12-25T15:30:45Z", "2023-12-25T15:30:45Z", false},
		{"2023-12-25T15:30:45+00:00", "2023-12-25T15:30:45Z", false},
		{"2023-12-25 15:30:45", "2023-12-25T15:30:45Z", false},
		{"2023-12-25", "2023-12-25T00:00:00Z", false},
		{"2023/12/25", "2023-12-25T00:00:00Z", false},
		{"20231225", "2023-12-25T00:00:00Z", false},
		{"1640995200", "2022-01-01T00:00:00Z", false},    // Unix timestamp
		{"1640995200000", "2022-01-01T00:00:00Z", false}, // Unix timestamp milliseconds
		{"15:30:45", "", false},                          // Time only
		{"invalid", "", true},
		{"", "", true},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			dt, err := ParseOptimized(test.input)

			if test.hasError {
				if err == nil {
					t.Errorf("Expected error for input '%s', but got none", test.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for input '%s': %v", test.input, err)
				return
			}

			if test.expected != "" && test.input != "15:30:45" { // Skip time-only test
				result := dt.UTC().ToISO8601String()
				if result != test.expected {
					t.Errorf("ParseOptimized('%s'): expected %s, got %s", test.input, test.expected, result)
				}
			}
		})
	}
}

func TestParseOptimizedInLocation(t *testing.T) {
	ny, _ := time.LoadLocation("America/New_York")

	tests := []struct {
		input    string
		location *time.Location
		hasError bool
	}{
		{"2023-12-25T15:30:45Z", ny, false},
		{"2023-12-25 15:30:45", ny, false},
		{"2023-12-25", ny, false},
		{"1640995200", ny, false}, // Unix timestamp
		{"invalid", ny, true},
		{"", ny, true},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			dt, err := ParseOptimizedInLocation(test.input, test.location)

			if test.hasError {
				if err == nil {
					t.Errorf("Expected error for input '%s', but got none", test.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for input '%s': %v", test.input, err)
				return
			}

			// Verify location is correct
			if dt.Location() != test.location && !test.hasError {
				t.Errorf("ParseOptimizedInLocation('%s'): expected location %v, got %v",
					test.input, test.location, dt.Location())
			}
		})
	}
}

func TestDetectLayout(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"2023-12-25T15:30:45Z", "2006-01-02T15:04:05Z"},
		{"2023-12-25T15:30:45", "2006-01-02T15:04:05"},
		{"2023-12-25 15:30:45", "2006-01-02 15:04:05"},
		{"2023-12-25 15:30", "2006-01-02 15:04"},
		{"2023-12-25", "2006-01-02"},
		{"20231225", "20060102"},
		{"2023/12/25 15:30:45", "2006/01/02 15:04:05"},
		{"2023/12/25", "2006/01/02"},
		{"15:30:45", "15:04:05"},
		{"15:30", "15:04"},
		{"invalid", ""},
		{"", ""},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := detectLayout(test.input)
			if result != test.expected {
				t.Errorf("detectLayout('%s'): expected '%s', got '%s'", test.input, test.expected, result)
			}
		})
	}
}

func TestIsAllDigits(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"123", true},
		{"1640995200", true},
		{"123abc", false},
		{"12.34", false},
		{"12-34", false},
		{"", false},
		{"abc", false},
		{"0", true},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := isAllDigits(test.input)
			if result != test.expected {
				t.Errorf("isAllDigits('%s'): expected %v, got %v", test.input, test.expected, result)
			}
		})
	}
}

func TestParseBatch(t *testing.T) {
	inputs := []string{
		"2023-12-25T15:30:45Z",
		"2023-12-25 15:30:45",
		"2023-12-25",
		"invalid",
		"",
		"1640995200",
	}

	results, errors := ParseBatch(inputs, time.UTC)

	if len(results) != len(inputs) {
		t.Errorf("ParseBatch: expected %d results, got %d", len(inputs), len(results))
	}

	if len(errors) != len(inputs) {
		t.Errorf("ParseBatch: expected %d errors, got %d", len(inputs), len(errors))
	}

	// Check that valid inputs parsed correctly
	validInputs := []int{0, 1, 2, 5} // indices of valid inputs
	for _, i := range validInputs {
		if errors[i] != nil {
			t.Errorf("ParseBatch: unexpected error for valid input '%s': %v", inputs[i], errors[i])
		}
		if results[i].IsZero() {
			t.Errorf("ParseBatch: got zero DateTime for valid input '%s'", inputs[i])
		}
	}

	// Check that invalid inputs have errors
	invalidInputs := []int{3, 4} // indices of invalid inputs
	for _, i := range invalidInputs {
		if errors[i] == nil {
			t.Errorf("ParseBatch: expected error for invalid input '%s'", inputs[i])
		}
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		a, b     int
		expected int
	}{
		{1, 2, 1},
		{5, 3, 3},
		{0, 0, 0},
		{-1, 1, -1},
		{10, 5, 5},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			result := min(test.a, test.b)
			if result != test.expected {
				t.Errorf("min(%d, %d): expected %d, got %d", test.a, test.b, test.expected, result)
			}
		})
	}
}
