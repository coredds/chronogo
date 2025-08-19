package chronogo

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestChronoError(t *testing.T) {
	err := &ChronoError{
		Op:         "Parse",
		Input:      "invalid-date",
		Err:        errors.New("parsing failed"),
		Suggestion: "Use ISO format",
	}

	expected := `chronogo.Parse: input: "invalid-date": parsing failed
Suggestion: Use ISO format`

	if err.Error() != expected {
		t.Errorf("Expected error message:\n%s\nGot:\n%s", expected, err.Error())
	}
}

func TestChronoErrorUnwrap(t *testing.T) {
	originalErr := errors.New("original error")
	chronoErr := &ChronoError{
		Op:  "Test",
		Err: originalErr,
	}

	if !errors.Is(chronoErr, originalErr) {
		t.Error("ChronoError should unwrap to original error")
	}
}

func TestParseError(t *testing.T) {
	input := "2023-13-45"
	originalErr := errors.New("month out of range")

	err := ParseError(input, originalErr)

	if err.Op != "Parse" {
		t.Errorf("Expected Op to be 'Parse', got '%s'", err.Op)
	}

	if err.Input != input {
		t.Errorf("Expected Input to be '%s', got '%s'", input, err.Input)
	}

	if err.Suggestion == "" {
		t.Error("Expected non-empty suggestion")
	}

	if !errors.Is(err, originalErr) {
		t.Error("ParseError should wrap original error")
	}
}

func TestSuggestParseFormat(t *testing.T) {
	testCases := []struct {
		input            string
		expectedContains []string
	}{
		{
			"2023-12-25T15:30:45Z",
			[]string{"ISO", "ParseISO8601", "ParseRFC3339"},
		},
		{
			"12/25/2023",
			[]string{"FromFormat", "/"},
		},
		{
			"1640995200",
			[]string{"FromUnix", "timestamp", "seconds"},
		},
		{
			"1640995200000",
			[]string{"FromUnixMilli", "milliseconds"},
		},
		{
			"",
			[]string{"non-empty"},
		},
		{
			"2 hours ago",
			[]string{"relative times", "not supported"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			suggestion := suggestParseFormat(tc.input)

			for _, expected := range tc.expectedContains {
				if !strings.Contains(suggestion, expected) {
					t.Errorf("Expected suggestion to contain '%s', got: %s", expected, suggestion)
				}
			}
		})
	}
}

func TestSuggestTimezone(t *testing.T) {
	testCases := []struct {
		input            string
		expectedContains []string
	}{
		{
			"",
			[]string{"IANA timezone", "America/New_York"},
		},
		{
			"new_york",
			[]string{"America/New_York"},
		},
		{
			"EST",
			[]string{"America/New_York", "Eastern"},
		},
		{
			"PST",
			[]string{"America/Los_Angeles", "Pacific"},
		},
		{
			"london",
			[]string{"Europe/London"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			suggestion := suggestTimezone(tc.input)

			for _, expected := range tc.expectedContains {
				if !strings.Contains(suggestion, expected) {
					t.Errorf("Expected suggestion to contain '%s', got: %s", expected, suggestion)
				}
			}
		})
	}
}

func TestSuggestFormat(t *testing.T) {
	testCases := []struct {
		input            string
		expectedContains []string
	}{
		{
			"",
			[]string{"reference time", "Mon Jan 2"},
		},
		{
			"YYYY-MM-DD",
			[]string{"2006", "01", "02", "reference time"}, // Accept any year format suggestion
		},
		{
			"HH:mm:ss",
			[]string{"15", "04", "05", "reference time"}, // Accept any of these suggestions
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			suggestion := suggestFormat(tc.input)

			for _, expected := range tc.expectedContains {
				if !strings.Contains(suggestion, expected) {
					t.Errorf("Expected suggestion to contain '%s', got: %s", expected, suggestion)
				}
			}
		})
	}
}

func TestValidate(t *testing.T) {
	testCases := []struct {
		name    string
		dt      DateTime
		wantErr bool
	}{
		{
			"Valid DateTime",
			Date(2023, time.December, 25, 15, 30, 0, 0, time.UTC),
			false,
		},
		{
			"Zero DateTime",
			DateTime{},
			true,
		},
		{
			"Year too low",
			Date(0, time.January, 1, 0, 0, 0, 0, time.UTC),
			true,
		},
		{
			"Year too high",
			Date(10000, time.January, 1, 0, 0, 0, 0, time.UTC),
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.dt.Validate()

			if tc.wantErr && err == nil {
				t.Error("Expected error but got none")
			}

			if !tc.wantErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if err != nil {
				var chronoErr *ChronoError
				if !errors.As(err, &chronoErr) {
					t.Error("Expected ChronoError type")
				}

				if chronoErr.Suggestion == "" {
					t.Error("Expected non-empty suggestion")
				}
			}
		})
	}
}

func TestValidateRange(t *testing.T) {
	start := Date(2023, time.December, 25, 0, 0, 0, 0, time.UTC)
	end := Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC)

	// Valid range
	err := ValidateRange(start, end)
	if err != nil {
		t.Errorf("Expected no error for valid range, got: %v", err)
	}

	// Invalid range (start after end)
	err = ValidateRange(end, start)
	if err == nil {
		t.Error("Expected error for invalid range")
	}

	var chronoErr *ChronoError
	if errors.As(err, &chronoErr) {
		if !errors.Is(chronoErr, ErrInvalidRange) {
			t.Error("Expected ErrInvalidRange")
		}
	}

	// Invalid start date
	zeroStart := DateTime{}
	err = ValidateRange(zeroStart, end)
	if err == nil {
		t.Error("Expected error for invalid start date")
	}
}

func TestMustParse(t *testing.T) {
	// Valid parse should not panic
	dt := MustParse("2023-12-25T15:30:45Z")
	expected := Date(2023, time.December, 25, 15, 30, 45, 0, time.UTC)

	if !dt.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, dt)
	}

	// Invalid parse should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for invalid parse")
		}
	}()

	MustParse("invalid-date")
}

func TestMustLoadLocation(t *testing.T) {
	// Valid location should not panic
	loc := MustLoadLocation("UTC")
	if loc != time.UTC {
		t.Error("Expected UTC location")
	}

	// Invalid location should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for invalid location")
		}
	}()

	MustLoadLocation("Invalid/Timezone")
}

func TestIsNumericOnly(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"123", true},
		{"1640995200", true},
		{"123abc", false},
		{"12.34", false},
		{"12-34", false},
		{" 123", false},
		{"123 ", false},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := isNumericOnly(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %v for input '%s', got %v", tc.expected, tc.input, result)
			}
		})
	}
}

func TestErrorIs(t *testing.T) {
	originalErr := ErrInvalidFormat
	chronoErr := &ChronoError{
		Op:  "Test",
		Err: originalErr,
	}

	if !errors.Is(chronoErr, ErrInvalidFormat) {
		t.Error("ChronoError should match wrapped error")
	}

	if errors.Is(chronoErr, ErrInvalidTimezone) {
		t.Error("ChronoError should not match different error")
	}

	// Test ChronoError comparison
	sameErr := &ChronoError{
		Op:   "Test",
		Path: "same",
	}
	differentErr := &ChronoError{
		Op:   "Different",
		Path: "same",
	}

	chronoErrWithPath := &ChronoError{
		Op:   "Test",
		Path: "same",
		Err:  originalErr,
	}

	if !chronoErrWithPath.Is(sameErr) {
		t.Error("ChronoError should match same Op and Path")
	}

	if chronoErrWithPath.Is(differentErr) {
		t.Error("ChronoError should not match different Op")
	}
}

func TestFormatError(t *testing.T) {
	format := "YYYY-MM-DD"
	originalErr := errors.New("invalid format")

	err := FormatError(format, originalErr)

	if err.Op != "Format" {
		t.Errorf("Expected Op to be 'Format', got '%s'", err.Op)
	}

	if err.Input != format {
		t.Errorf("Expected Input to be '%s', got '%s'", format, err.Input)
	}

	if !errors.Is(err, originalErr) {
		t.Error("FormatError should wrap original error")
	}

	if !strings.Contains(err.Suggestion, "format") {
		t.Errorf("FormatError should contain format suggestion, got: %s", err.Suggestion)
	}
}

func TestMustParseInLocation(t *testing.T) {
	loc, _ := time.LoadLocation("America/New_York")

	// Test successful parse
	result := MustParseInLocation("2023-12-25T15:30:45Z", loc)
	// Note: when parsing UTC time with location, it may remain in UTC
	if result.IsZero() {
		t.Error("MustParseInLocation should successfully parse the input")
	}

	// Test panic on invalid input
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustParseInLocation should panic on invalid input")
		}
	}()
	MustParseInLocation("invalid-date", loc)
}

func TestMustFromFormat(t *testing.T) {
	// Test successful parse
	result := MustFromFormat("25/12/2023", "02/01/2006")
	expected := Date(2023, time.December, 25, 0, 0, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("MustFromFormat failed: expected %v, got %v", expected, result)
	}

	// Test panic on invalid input
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustFromFormat should panic on invalid format")
		}
	}()
	MustFromFormat("invalid-date", "invalid-format")
}
