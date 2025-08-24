package chronogo

import (
	"testing"
	"time"
)

func TestParseWithOptions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		options  ParseOptions
		expected string
		wantErr  bool
	}{
		{
			name:     "strict parsing - valid RFC3339",
			input:    "2023-12-25T10:30:45Z",
			options:  ParseOptions{Strict: true},
			expected: "2023-12-25T10:30:45Z",
			wantErr:  false,
		},
		{
			name:     "strict parsing - invalid format",
			input:    "2023/12/25 10:30:45",
			options:  ParseOptions{Strict: true},
			expected: "",
			wantErr:  true,
		},
		{
			name:     "lenient parsing - slash format",
			input:    "2023/12/25 10:30:45",
			options:  ParseOptions{Strict: false},
			expected: "2023-12-25T10:30:45Z",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input, tt.options)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Parse() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Parse() unexpected error: %v", err)
				return
			}

			if result.Format(time.RFC3339) != tt.expected {
				t.Errorf("Parse() = %v, want %v", result.Format(time.RFC3339), tt.expected)
			}
		})
	}
}

func TestFromFormatTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		format   string
		expected string
		wantErr  bool
	}{
		{
			name:     "YYYY-MM-DD format",
			input:    "2023-12-25",
			format:   "YYYY-MM-DD",
			expected: "2023-12-25T00:00:00Z",
			wantErr:  false,
		},
		{
			name:     "DD/MM/YYYY HH:mm format",
			input:    "25/12/2023 14:30",
			format:   "DD/MM/YYYY HH:mm",
			expected: "2023-12-25T14:30:00Z",
			wantErr:  false,
		},
		{
			name:     "MMMM D, YYYY format",
			input:    "December 25, 2023",
			format:   "MMMM D, YYYY",
			expected: "2023-12-25T00:00:00Z",
			wantErr:  false,
		},
		{
			name:     "Invalid format",
			input:    "2023-12-25",
			format:   "DD/MM/YYYY",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FromFormatTokens(tt.input, tt.format)

			if tt.wantErr {
				if err == nil {
					t.Errorf("FromFormatTokens() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("FromFormatTokens() unexpected error: %v", err)
				return
			}

			if result.Format(time.RFC3339) != tt.expected {
				t.Errorf("FromFormatTokens() = %v, want %v", result.Format(time.RFC3339), tt.expected)
			}
		})
	}
}

func TestParseOrdinalDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "ordinal date with separator",
			input:    "2023-359",
			expected: "2023-12-25T00:00:00Z",
			wantErr:  false,
		},
		{
			name:     "compact ordinal date",
			input:    "2023359",
			expected: "2023-12-25T00:00:00Z",
			wantErr:  false,
		},
		{
			name:     "leap year ordinal date",
			input:    "2024-366",
			expected: "2024-12-31T00:00:00Z",
			wantErr:  false,
		},
		{
			name:     "invalid day of year",
			input:    "2023-366",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "invalid day of year - too small",
			input:    "2023-000",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Parse() expected error for ordinal date, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Parse() unexpected error for ordinal date: %v", err)
				return
			}

			if result.Format(time.RFC3339) != tt.expected {
				t.Errorf("Parse() ordinal date = %v, want %v", result.Format(time.RFC3339), tt.expected)
			}
		})
	}
}

func TestParseWeekDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "week date with day",
			input:    "2023-W52-1", // Monday of week 52, 2023
			expected: "2023-12-25T00:00:00Z",
			wantErr:  false,
		},
		{
			name:     "compact week date with day",
			input:    "2023W521",
			expected: "2023-12-25T00:00:00Z",
			wantErr:  false,
		},
		{
			name:     "week date without day (defaults to Monday)",
			input:    "2023-W52",
			expected: "2023-12-25T00:00:00Z",
			wantErr:  false,
		},
		{
			name:     "compact week date without day",
			input:    "2023W52",
			expected: "2023-12-25T00:00:00Z",
			wantErr:  false,
		},
		{
			name:     "invalid week number",
			input:    "2023-W54-1",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "invalid day of week",
			input:    "2023-W52-8",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Parse() expected error for week date, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Parse() unexpected error for week date: %v", err)
				return
			}

			if result.Format(time.RFC3339) != tt.expected {
				t.Errorf("Parse() week date = %v, want %v", result.Format(time.RFC3339), tt.expected)
			}
		})
	}
}

func TestParseInterval(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "datetime interval",
			input:   "2023-01-01T00:00:00Z/2023-12-31T23:59:59Z",
			wantErr: false,
		},
		{
			name:    "datetime with duration",
			input:   "2023-01-01T00:00:00Z/P1Y",
			wantErr: false,
		},
		{
			name:    "duration with datetime",
			input:   "P1Y/2024-01-01T00:00:00Z",
			wantErr: false,
		},
		{
			name:    "complex duration",
			input:   "2023-01-01T00:00:00Z/P1Y2M3DT4H5M6S",
			wantErr: false,
		},
		{
			name:    "not an interval - single datetime",
			input:   "2023-01-01T00:00:00Z",
			wantErr: false, // This should parse as a regular datetime, not an interval
		},
		{
			name:    "invalid interval parts",
			input:   "invalid/P1Y",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Parse() expected error for interval, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Parse() unexpected error for interval: %v", err)
				return
			}

			// For intervals, we return the start datetime
			if result.IsZero() {
				t.Errorf("Parse() interval returned zero datetime")
			}
		})
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ISO8601Duration
		wantErr  bool
	}{
		{
			name:  "simple year duration",
			input: "P1Y",
			expected: ISO8601Duration{
				Years: 1,
			},
			wantErr: false,
		},
		{
			name:  "complex duration",
			input: "P1Y2M3DT4H5M6.5S",
			expected: ISO8601Duration{
				Years:   1,
				Months:  2,
				Days:    3,
				Hours:   4,
				Minutes: 5,
				Seconds: 6.5,
			},
			wantErr: false,
		},
		{
			name:  "time-only duration",
			input: "PT2H30M",
			expected: ISO8601Duration{
				Hours:   2,
				Minutes: 30,
			},
			wantErr: false,
		},
		{
			name:    "invalid duration",
			input:   "P1Y2M3DT4H5M6.5",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseDuration(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("parseDuration() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("parseDuration() unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("parseDuration() = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

func TestParseWithFallback(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "US date format",
			input:   "12/25/2023",
			wantErr: false,
		},
		{
			name:    "human readable date",
			input:   "Dec 25, 2023",
			wantErr: false,
		},
		{
			name:    "email date format",
			input:   "Mon, 25 Dec 2023 10:30:45 UTC",
			wantErr: false,
		},
		{
			name:    "completely invalid",
			input:   "not a date",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseWithFallback(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseWithFallback() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseWithFallback() unexpected error: %v", err)
				return
			}

			if result.IsZero() {
				t.Errorf("ParseWithFallback() returned zero datetime")
			}
		})
	}
}

func TestParseMultiple(t *testing.T) {
	formats := []string{"YYYY-MM-DD", "DD/MM/YYYY", "MMMM D, YYYY"}

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "first format matches",
			input:   "2023-12-25",
			wantErr: false,
		},
		{
			name:    "second format matches",
			input:   "25/12/2023",
			wantErr: false,
		},
		{
			name:    "third format matches",
			input:   "December 25, 2023",
			wantErr: false,
		},
		{
			name:    "no format matches",
			input:   "25-Dec-2023",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseMultiple(tt.input, formats)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseMultiple() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseMultiple() unexpected error: %v", err)
				return
			}

			if result.IsZero() {
				t.Errorf("ParseMultiple() returned zero datetime")
			}
		})
	}
}

func TestParseAny(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "quoted date",
			input:   `"2023-12-25"`,
			wantErr: false,
		},
		{
			name:    "date with ordinals",
			input:   "December 25th, 2023",
			wantErr: false,
		},
		{
			name:    "date with extra text",
			input:   "Date: 2023-12-25",
			wantErr: true, // This would require more sophisticated parsing
		},
		{
			name:    "Unix timestamp",
			input:   "1703505645",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseAny(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseAny() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseAny() unexpected error: %v", err)
				return
			}

			if result.IsZero() {
				t.Errorf("ParseAny() returned zero datetime")
			}
		})
	}
}

func TestIsValidDateTimeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid ISO date",
			input:    "2023-12-25T10:30:45Z",
			expected: true,
		},
		{
			name:     "valid simple date",
			input:    "2023-12-25",
			expected: true,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "no digits",
			input:    "hello world",
			expected: false,
		},
		{
			name:     "too long",
			input:    "this is way too long to be a reasonable datetime string representation",
			expected: false,
		},
		{
			name:     "too short",
			input:    "123",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidDateTimeString(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidDateTimeString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestConvertTokenFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic date format",
			input:    "YYYY-MM-DD",
			expected: "2006-01-02",
		},
		{
			name:     "datetime format",
			input:    "YYYY-MM-DD HH:mm:ss",
			expected: "2006-01-02 15:04:05",
		},
		{
			name:     "month names",
			input:    "MMMM DD, YYYY",
			expected: "January 02, 2006",
		},
		{
			name:     "12-hour time with AM/PM",
			input:    "MM/DD/YYYY hh:mm A",
			expected: "01/02/2006 03:04 PM",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertTokenFormat(tt.input)
			if result != tt.expected {
				t.Errorf("convertTokenFormat() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetSupportedFormats(t *testing.T) {
	formats := GetSupportedFormats()

	if len(formats) == 0 {
		t.Error("GetSupportedFormats() returned empty slice")
	}

	// Check that some expected formats are present
	expectedFormats := []string{
		"2006-01-02T15:04:05Z",
		"2006-01-02",
		"15:04:05",
	}

	for _, expected := range expectedFormats {
		found := false
		for _, format := range formats {
			if format == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetSupportedFormats() missing expected format: %s", expected)
		}
	}
}

func BenchmarkParseOptimized(b *testing.B) {
	testStrings := []string{
		"2023-12-25T10:30:45Z",
		"2023-12-25 10:30:45",
		"2023-12-25",
		"1703505645",
		"20231225",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, s := range testStrings {
			_, _ = ParseOptimized(s)
		}
	}
}

func BenchmarkParseWithFallback(b *testing.B) {
	testStrings := []string{
		"12/25/2023",
		"Dec 25, 2023",
		"25 December 2023",
		"Mon, 25 Dec 2023 10:30:45 UTC",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, s := range testStrings {
			_, _ = ParseWithFallback(s)
		}
	}
}
