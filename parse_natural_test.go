package chronogo

import (
	"strings"
	"testing"
	"time"
)

// TestParseNaturalEnglish tests natural language parsing in English
func TestParseNaturalEnglish(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		input       string
		wantDayDiff int // Days difference from today (positive = future, negative = past)
		allowDelta  int // Allow +/- days for fuzzy matching
	}{
		{"Tomorrow", "tomorrow", 1, 0},
		{"Yesterday", "yesterday", -1, 0},
		{"Today", "today", 0, 0},
		{"Next week", "next week", 7, 0},
		{"Last week", "last week", -7, 0},
		{"In 3 days", "in 3 days", 3, 0},
		{"3 days ago", "3 days ago", -3, 0},
		{"Next Monday", "next Monday", 0, 7}, // Within a week
		{"Last Friday", "last Friday", 0, 7},
		{"In 2 weeks", "in 2 weeks", 14, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v", tt.input, err)
			}

			// Calculate expected date
			expected := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, tt.wantDayDiff)
			actualDays := result.Format("2006-01-02")
			expectedDays := DateTime{expected}.Format("2006-01-02")

			// For fuzzy matching (like "next Monday"), check if within range
			if tt.allowDelta > 0 {
				diff := result.Sub(DateTime{now}).Hours() / 24
				if diff < 0 {
					diff = -diff
				}
				if diff > float64(tt.allowDelta) {
					t.Errorf("Parse(%q) = %v, day diff = %.1f, want within %d days",
						tt.input, actualDays, diff, tt.allowDelta)
				}
			} else {
				if actualDays != expectedDays {
					t.Errorf("Parse(%q) = %v, want %v", tt.input, actualDays, expectedDays)
				}
			}
		})
	}
}

// TestParseMultiLanguage tests natural language parsing in multiple languages
func TestParseMultiLanguage(t *testing.T) {
	now := time.Now()
	tomorrow := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1)
	tomorrowStr := DateTime{tomorrow}.Format("2006-01-02")

	tests := []struct {
		name  string
		input string
		lang  string
	}{
		{"Spanish - mañana", "mañana", "es"},
		{"Portuguese - amanhã", "amanhã", "pt"},
		{"French - demain", "demain", "fr"},
		{"German - morgen", "morgen", "de"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v", tt.input, err)
			}

			actualDays := result.Format("2006-01-02")
			if actualDays != tomorrowStr {
				t.Errorf("Parse(%q) = %v, want %v (tomorrow)", tt.input, actualDays, tomorrowStr)
			}
		})
	}
}

// TestParseTechnicalFormatsStillWork ensures existing technical parsing still works
func TestParseTechnicalFormatsStillWork(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"ISO 8601", "2024-01-15T14:30:00Z", "2024-01-15 14:30:00"},
		{"RFC 3339", "2024-01-15T14:30:00+00:00", "2024-01-15 14:30:00"},
		{"Date only", "2024-01-15", "2024-01-15 00:00:00"},
		{"Unix timestamp (seconds)", "1705329000", "2024-01-15 14:30:00"},
		{"Compact date", "20240115", "2024-01-15 00:00:00"},
		{"ISO Week", "2024-W03-1", "2024-01-15 00:00:00"},
		{"ISO Ordinal", "2024-015", "2024-01-15 00:00:00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v", tt.input, err)
			}

			got := result.Format("2006-01-02 15:04:05")
			if got != tt.want {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// TestParseWithConfig tests the new ParseWith function with custom configurations
func TestParseWithConfig(t *testing.T) {
	t.Run("Strict mode rejects natural language", func(t *testing.T) {
		config := ParseConfig{
			Strict:   true,
			Location: time.UTC,
		}

		_, err := ParseWith("tomorrow", config)
		if err == nil {
			t.Error("ParseWith(\"tomorrow\", strict) should fail but succeeded")
		}
	})

	t.Run("Strict mode accepts technical formats", func(t *testing.T) {
		config := ParseConfig{
			Strict:   true,
			Location: time.UTC,
		}

		result, err := ParseWith("2024-01-15T14:30:00Z", config)
		if err != nil {
			t.Fatalf("ParseWith(ISO8601, strict) error = %v", err)
		}

		want := "2024-01-15 14:30:00"
		got := result.Format("2006-01-02 15:04:05")
		if got != want {
			t.Errorf("ParseWith(ISO8601, strict) = %v, want %v", got, want)
		}
	})

	t.Run("Custom languages", func(t *testing.T) {
		config := ParseConfig{
			Languages: []string{"en", "es"},
			Location:  time.UTC,
		}

		// Test English
		_, err := ParseWith("tomorrow", config)
		if err != nil {
			t.Errorf("ParseWith(\"tomorrow\", en+es) error = %v", err)
		}

		// Test Spanish
		_, err = ParseWith("mañana", config)
		if err != nil {
			t.Errorf("ParseWith(\"mañana\", en+es) error = %v", err)
		}
	})

	t.Run("Location affects relative dates", func(t *testing.T) {
		ny, _ := time.LoadLocation("America/New_York")
		tokyo, _ := time.LoadLocation("Asia/Tokyo")

		configNY := ParseConfig{
			Languages: []string{"en"},
			Location:  ny,
		}

		configTokyo := ParseConfig{
			Languages: []string{"en"},
			Location:  tokyo,
		}

		resultNY, err1 := ParseWith("tomorrow", configNY)
		resultTokyo, err2 := ParseWith("tomorrow", configTokyo)

		if err1 != nil || err2 != nil {
			t.Fatalf("ParseWith errors: NY=%v, Tokyo=%v", err1, err2)
		}

		// Should have different timezones
		if resultNY.Location().String() == resultTokyo.Location().String() {
			t.Errorf("Expected different locations, both are %v", resultNY.Location())
		}
	})
}

// TestSetDefaultParseLanguages tests the language configuration API
func TestSetDefaultParseLanguages(t *testing.T) {
	// Save original
	original := GetDefaultParseLanguages()
	defer func() {
		DefaultParseConfig.Languages = original
	}()

	// Test setting custom languages
	SetDefaultParseLanguages("en", "fr")
	got := GetDefaultParseLanguages()

	if len(got) != 2 || got[0] != "en" || got[1] != "fr" {
		t.Errorf("GetDefaultParseLanguages() = %v, want [en fr]", got)
	}

	// Verify it affects Parse()
	_, err := Parse("demain") // French for "tomorrow"
	if err != nil {
		t.Errorf("Parse(\"demain\") with French enabled error = %v", err)
	}
}

// TestParseStrictOptions tests the Strict option in ParseOptions
func TestParseStrictOptions(t *testing.T) {
	t.Run("Strict via ParseOptions", func(t *testing.T) {
		_, err := Parse("tomorrow", ParseOptions{Strict: true})
		if err == nil {
			t.Error("Parse(\"tomorrow\", Strict:true) should fail but succeeded")
		}
	})

	t.Run("Non-strict via ParseOptions", func(t *testing.T) {
		_, err := Parse("tomorrow", ParseOptions{Strict: false})
		if err != nil {
			t.Errorf("Parse(\"tomorrow\", Strict:false) error = %v", err)
		}
	})
}

// TestParseInLocationWithNaturalLanguage tests location-aware natural language parsing
func TestParseInLocationWithNaturalLanguage(t *testing.T) {
	ny, _ := time.LoadLocation("America/New_York")

	result, err := ParseInLocation("tomorrow", ny)
	if err != nil {
		t.Fatalf("ParseInLocation(\"tomorrow\", NY) error = %v", err)
	}

	if result.Location().String() != "America/New_York" {
		t.Errorf("ParseInLocation location = %v, want America/New_York", result.Location())
	}
}

// TestParseErrors tests error handling for invalid inputs
func TestParseErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Empty string", ""},
		{"Invalid gibberish", "xyzabc123invalidtext"},
		{"Only whitespace", "   "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.input)
			if err == nil {
				t.Errorf("Parse(%q) should fail but succeeded", tt.input)
			}

			// Ensure error is a ChronoError
			if _, ok := err.(*ChronoError); !ok {
				t.Errorf("Parse(%q) error type = %T, want *ChronoError", tt.input, err)
			}
		})
	}
}

// BenchmarkParseNatural benchmarks natural language parsing
func BenchmarkParseNatural(b *testing.B) {
	inputs := []string{
		"tomorrow",
		"yesterday",
		"next Monday",
		"in 3 days",
		"last week",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input := inputs[i%len(inputs)]
		_, err := Parse(input)
		if err != nil {
			b.Fatalf("Parse(%q) error = %v", input, err)
		}
	}
}

// BenchmarkParseTechnical benchmarks technical format parsing
func BenchmarkParseTechnical(b *testing.B) {
	inputs := []string{
		"2024-01-15T14:30:00Z",
		"2024-01-15",
		"20240115",
		"1705329000",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input := inputs[i%len(inputs)]
		_, err := Parse(input)
		if err != nil {
			b.Fatalf("Parse(%q) error = %v", input, err)
		}
	}
}

// BenchmarkParseWith benchmarks the ParseWith function
func BenchmarkParseWith(b *testing.B) {
	config := ParseConfig{
		Languages: []string{"en"},
		Location:  time.UTC,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ParseWith("tomorrow", config)
		if err != nil {
			b.Fatalf("ParseWith error = %v", err)
		}
	}
}

// TestParseBackwardCompatibility ensures we didn't break existing tests
func TestParseBackwardCompatibility(t *testing.T) {
	// Test that old parsing behavior still works
	tests := []struct {
		input string
		want  string
	}{
		{"2006-01-02 15:04:05", "2006-01-02 15:04:05"},
		{"2006/01/02", "2006-01-02 00:00:00"},
		{"20060102", "2006-01-02 00:00:00"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v", tt.input, err)
			}

			got := result.Format("2006-01-02 15:04:05")
			if !strings.HasPrefix(got, strings.TrimSuffix(tt.want, " 00:00:00")) {
				t.Errorf("Parse(%q) = %v, want prefix %v", tt.input, got, tt.want)
			}
		})
	}
}
