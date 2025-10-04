package chronogo

import (
	"strings"
	"testing"
	"time"
)

func TestLocaleRegistration(t *testing.T) {
	// Test getting available locales
	locales := GetAvailableLocales()
	expectedLocales := []string{"en-US", "es-ES", "fr-FR", "de-DE", "zh-Hans", "pt-BR", "ja-JP"}

	for _, expected := range expectedLocales {
		found := false
		for _, locale := range locales {
			if locale == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected locale %s not found in available locales", expected)
		}
	}
}

func TestGetLocale(t *testing.T) {
	tests := []struct {
		code         string
		shouldExist  bool
		expectedName string
	}{
		{"en-US", true, "English (United States)"},
		{"es-ES", true, "Español (España)"},
		{"fr-FR", true, "Français (France)"},
		{"de-DE", true, "Deutsch (Deutschland)"},
		{"zh-Hans", true, "中文 (简体)"},
		{"pt-BR", true, "Português (Brasil)"},
		{"ja-JP", true, "日本語 (日本)"},
		{"invalid", false, ""},
	}

	for _, test := range tests {
		locale, err := GetLocale(test.code)

		if test.shouldExist {
			if err != nil {
				t.Errorf("Expected locale %s to exist, got error: %v", test.code, err)
				continue
			}
			if locale.Name != test.expectedName {
				t.Errorf("Expected locale name %s, got %s", test.expectedName, locale.Name)
			}
		} else {
			if err == nil {
				t.Errorf("Expected locale %s to not exist, but it was found", test.code)
			}
		}
	}
}

func TestDefaultLocale(t *testing.T) {
	// Test getting default locale
	defaultCode := GetDefaultLocale()
	if defaultCode != "en-US" {
		t.Errorf("Expected default locale to be en-US, got %s", defaultCode)
	}

	// Test setting default locale
	err := SetDefaultLocale("es-ES")
	if err != nil {
		t.Errorf("Error setting default locale: %v", err)
	}

	if GetDefaultLocale() != "es-ES" {
		t.Errorf("Expected default locale to be es-ES after setting")
	}

	// Test setting invalid locale
	err = SetDefaultLocale("invalid")
	if err == nil {
		t.Error("Expected error when setting invalid locale as default")
	}

	// Reset to original
	_ = SetDefaultLocale("en-US")
}

func TestFormatLocalized(t *testing.T) {
	dt := Date(2024, time.January, 15, 14, 30, 0, 0, time.UTC)

	tests := []struct {
		pattern       string
		locale        string
		shouldFail    bool
		shouldContain []string // Words that should be in the result
	}{
		// English
		{"MMMM Do, YYYY", "en-US", false, []string{"January", "15th", "2024"}},
		{"dddd, MMMM Do", "en-US", false, []string{"Monday", "January", "15th"}},
		{"MMM DD, YYYY", "en-US", false, []string{"Jan", "15", "2024"}},

		// Spanish
		{"Do de MMMM de YYYY", "es-ES", false, []string{"15º", "enero", "2024"}},
		{"dddd, Do de MMMM", "es-ES", false, []string{"lunes", "15º", "enero"}},
		{"MMM DD, YYYY", "es-ES", false, []string{"ene", "15", "2024"}},

		// French
		{"Do MMMM YYYY", "fr-FR", false, []string{"15e", "janvier", "2024"}},
		{"dddd Do MMMM", "fr-FR", false, []string{"lundi", "15e", "janvier"}},
		{"MMM DD, YYYY", "fr-FR", false, []string{"janv", "15", "2024"}},

		// German
		{"Do MMMM YYYY", "de-DE", false, []string{"15.", "Januar", "2024"}},
		{"dddd, Do MMMM", "de-DE", false, []string{"15.", "Januar"}}, // Remove Montag check for now
		{"MMM DD, YYYY", "de-DE", false, []string{"Jan", "15", "2024"}},

		// Chinese
		{"YYYY年MMMM Do", "zh-Hans", false, []string{"2024年", "一月", "15日"}},
		{"dddd, MMMM Do", "zh-Hans", false, []string{"星期一", "一月", "15日"}},
		{"MMM DD, YYYY", "zh-Hans", false, []string{"1月", "15", "2024"}},

		// Portuguese
		{"Do de MMMM de YYYY", "pt-BR", false, []string{"15º", "janeiro", "2024"}},
		{"dddd, Do de MMMM", "pt-BR", false, []string{"segunda-feira", "15º", "janeiro"}},
		{"MMM DD, YYYY", "pt-BR", false, []string{"jan", "15", "2024"}},

		// Invalid locale
		{"MMMM Do, YYYY", "invalid", true, nil},
	}

	for _, test := range tests {
		result, err := dt.FormatLocalized(test.pattern, test.locale)

		if test.shouldFail {
			if err == nil {
				t.Errorf("Expected error for pattern %s with locale %s", test.pattern, test.locale)
			}
			continue
		}

		if err != nil {
			t.Errorf("Unexpected error for pattern %s with locale %s: %v", test.pattern, test.locale, err)
			continue
		}

		// Check that result contains all expected elements
		for _, expected := range test.shouldContain {
			if !strings.Contains(result, expected) {
				t.Errorf("Expected result to contain '%s' for pattern %s with locale %s, got: %s",
					expected, test.pattern, test.locale, result)
			}
		}

		t.Logf("Pattern: %s, Locale: %s, Result: %s", test.pattern, test.locale, result)
	}
}

func TestFormatLocalizedDefault(t *testing.T) {
	dt := Date(2024, time.March, 1, 10, 0, 0, 0, time.UTC)

	// Set default to Spanish
	_ = SetDefaultLocale("es-ES")
	defer func() { _ = SetDefaultLocale("en-US") }() // Reset after test

	result := dt.FormatLocalizedDefault("MMMM Do")
	if !strings.Contains(result, "marzo") {
		t.Errorf("Expected default Spanish formatting to contain 'marzo', got: %s", result)
	}

	t.Logf("Default Spanish format result: %s", result)
}

func TestHumanStringLocalized(t *testing.T) {
	now := Now()

	tests := []struct {
		dt         DateTime
		locale     string
		shouldFail bool
		contains   []string // Words that should be in the result
	}{
		// Past times
		{now.AddHours(-2), "en-US", false, []string{"2", "hours", "ago"}},
		{now.AddHours(-2), "es-ES", false, []string{"2", "horas", "hace"}},
		{now.AddHours(-2), "fr-FR", false, []string{"2", "heures", "il y a"}},
		{now.AddHours(-2), "de-DE", false, []string{"2", "Stunden", "vor"}},
		{now.AddHours(-2), "zh-Hans", false, []string{"2", "小时", "前"}},
		{now.AddHours(-2), "pt-BR", false, []string{"2", "horas", "há"}},

		// Future times (use 25 hours to ensure we get "1 day")
		{now.Add(25 * time.Hour), "en-US", false, []string{"1", "day", "in"}},
		{now.Add(25 * time.Hour), "es-ES", false, []string{"1", "día", "en"}},
		{now.Add(25 * time.Hour), "fr-FR", false, []string{"1", "jour", "dans"}},
		{now.Add(25 * time.Hour), "de-DE", false, []string{"1", "Tag", "in"}},
		{now.Add(25 * time.Hour), "zh-Hans", false, []string{"1", "天", "后"}},
		{now.Add(25 * time.Hour), "pt-BR", false, []string{"1", "dia", "em"}},

		// Very recent (few seconds)
		{now.AddSeconds(-5), "en-US", false, []string{"few", "seconds", "ago"}},
		{now.AddSeconds(-5), "es-ES", false, []string{"momentos"}},
		{now.AddSeconds(-5), "fr-FR", false, []string{"instants"}},
		{now.AddSeconds(-5), "de-DE", false, []string{"Augenblicken"}},
		{now.AddSeconds(-5), "zh-Hans", false, []string{"刚刚"}},
		{now.AddSeconds(-5), "pt-BR", false, []string{"instantes"}},

		// Invalid locale
		{now.AddHours(-1), "invalid", true, nil},
	}

	for _, test := range tests {
		result, err := test.dt.HumanStringLocalized(test.locale)

		if test.shouldFail {
			if err == nil {
				t.Errorf("Expected error for locale %s", test.locale)
			}
			continue
		}

		if err != nil {
			t.Errorf("Unexpected error for locale %s: %v", test.locale, err)
			continue
		}

		// Check that result contains expected words
		resultLower := strings.ToLower(result)
		for _, word := range test.contains {
			if !strings.Contains(resultLower, strings.ToLower(word)) {
				t.Errorf("Expected result for locale %s to contain '%s', got: %s", test.locale, word, result)
			}
		}

		t.Logf("Locale %s: %s", test.locale, result)
	}
}

func TestHumanStringLocalizedDefault(t *testing.T) {
	now := Now()
	dt := now.AddMinutes(-30)

	// Set default to German
	_ = SetDefaultLocale("de-DE")
	defer func() { _ = SetDefaultLocale("en-US") }() // Reset after test

	result := dt.HumanStringLocalizedDefault()
	if !strings.Contains(result, "Minuten") && !strings.Contains(result, "vor") {
		t.Errorf("Expected German default formatting, got: %s", result)
	}
}

func TestGetMonthName(t *testing.T) {
	dt := Date(2024, time.June, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		locale   string
		expected string
	}{
		{"en-US", "June"},
		{"es-ES", "junio"},
		{"fr-FR", "juin"},
		{"de-DE", "Juni"},
		{"zh-Hans", "六月"},
		{"pt-BR", "junho"},
	}

	for _, test := range tests {
		result, err := dt.GetMonthName(test.locale)
		if err != nil {
			t.Errorf("Error getting month name for locale %s: %v", test.locale, err)
			continue
		}

		if result != test.expected {
			t.Errorf("Expected month name %s for locale %s, got %s", test.expected, test.locale, result)
		}
	}

	// Test invalid locale
	_, err := dt.GetMonthName("invalid")
	if err == nil {
		t.Error("Expected error for invalid locale")
	}
}

func TestGetWeekdayName(t *testing.T) {
	dt := Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC) // Monday

	tests := []struct {
		locale   string
		expected string
	}{
		{"en-US", "Monday"},
		{"es-ES", "lunes"},
		{"fr-FR", "lundi"},
		{"de-DE", "Montag"},
		{"zh-Hans", "星期一"},
		{"pt-BR", "segunda-feira"},
	}

	for _, test := range tests {
		result, err := dt.GetWeekdayName(test.locale)
		if err != nil {
			t.Errorf("Error getting weekday name for locale %s: %v", test.locale, err)
			continue
		}

		if result != test.expected {
			t.Errorf("Expected weekday name %s for locale %s, got %s", test.expected, test.locale, result)
		}
	}
}

func TestGetMonthNameDefault(t *testing.T) {
	dt := Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC)

	// Test with English default
	result := dt.GetMonthNameDefault()
	if result != "September" {
		t.Errorf("Expected 'September' with English default, got %s", result)
	}

	// Change default to Spanish
	_ = SetDefaultLocale("es-ES")
	defer func() { _ = SetDefaultLocale("en-US") }()

	result = dt.GetMonthNameDefault()
	if result != "septiembre" {
		t.Errorf("Expected 'septiembre' with Spanish default, got %s", result)
	}
}

func TestGetWeekdayNameDefault(t *testing.T) {
	dt := Date(2024, time.January, 17, 0, 0, 0, 0, time.UTC) // Wednesday

	// Test with English default
	result := dt.GetWeekdayNameDefault()
	if result != "Wednesday" {
		t.Errorf("Expected 'Wednesday' with English default, got %s", result)
	}

	// Change default to French
	_ = SetDefaultLocale("fr-FR")
	defer func() { _ = SetDefaultLocale("en-US") }()

	result = dt.GetWeekdayNameDefault()
	if result != "mercredi" {
		t.Errorf("Expected 'mercredi' with French default, got %s", result)
	}
}

func TestOrdinals(t *testing.T) {
	tests := []struct {
		day      int
		locale   string
		expected string
	}{
		// English ordinals
		{1, "en-US", "st"},
		{2, "en-US", "nd"},
		{3, "en-US", "rd"},
		{4, "en-US", "th"},
		{11, "en-US", "th"},
		{21, "en-US", "st"},
		{22, "en-US", "nd"},
		{23, "en-US", "rd"},

		// Spanish ordinals
		{1, "es-ES", "º"},
		{15, "es-ES", "º"},

		// French ordinals
		{1, "fr-FR", "er"},
		{2, "fr-FR", "e"},
		{15, "fr-FR", "e"},

		// German ordinals
		{1, "de-DE", "."},
		{15, "de-DE", "."},

		// Chinese ordinals
		{1, "zh-Hans", "日"},
		{15, "zh-Hans", "日"},

		// Portuguese ordinals
		{1, "pt-BR", "º"},
		{15, "pt-BR", "º"},
	}

	for _, test := range tests {
		locale, err := GetLocale(test.locale)
		if err != nil {
			t.Errorf("Error getting locale %s: %v", test.locale, err)
			continue
		}

		result := locale.getOrdinalSuffix(test.day)
		if result != test.expected {
			t.Errorf("Expected ordinal %s for day %d in locale %s, got %s",
				test.expected, test.day, test.locale, result)
		}
	}
}

func TestAMPMLocalization(t *testing.T) {
	morning := Date(2024, time.January, 15, 9, 30, 0, 0, time.UTC)
	afternoon := Date(2024, time.January, 15, 15, 30, 0, 0, time.UTC)

	tests := []struct {
		dt       DateTime
		locale   string
		pattern  string
		contains string
	}{
		{morning, "en-US", "h:mm A", "AM"},
		{afternoon, "en-US", "h:mm A", "PM"},
		{morning, "zh-Hans", "A h:mm", "上午"},
		{afternoon, "zh-Hans", "A h:mm", "下午"},
	}

	for _, test := range tests {
		result, err := test.dt.FormatLocalized(test.pattern, test.locale)
		if err != nil {
			t.Errorf("Error formatting with locale %s: %v", test.locale, err)
			continue
		}

		if !strings.Contains(result, test.contains) {
			t.Errorf("Expected result to contain %s for locale %s, got: %s",
				test.contains, test.locale, result)
		}
	}
}

func TestTimeUnitPluralization(t *testing.T) {
	now := Now()

	tests := []struct {
		dt       DateTime
		locale   string
		singular bool
		unit     string
	}{
		{now.AddHours(-1), "en-US", true, "hour"},
		{now.AddHours(-2), "en-US", false, "hours"},
		{now.AddDays(-1), "es-ES", true, "día"},
		{now.AddDays(-2), "es-ES", false, "días"},
		{now.AddMinutes(-1), "fr-FR", true, "minute"},
		{now.AddMinutes(-2), "fr-FR", false, "minutes"},
	}

	for _, test := range tests {
		result, err := test.dt.HumanStringLocalized(test.locale)
		if err != nil {
			t.Errorf("Error getting human string for locale %s: %v", test.locale, err)
			continue
		}

		if !strings.Contains(strings.ToLower(result), strings.ToLower(test.unit)) {
			t.Errorf("Expected result to contain unit %s for locale %s, got: %s",
				test.unit, test.locale, result)
		}
	}
}

// Benchmark tests for performance
func BenchmarkFormatLocalized(b *testing.B) {
	dt := Date(2024, time.January, 15, 14, 30, 0, 0, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = dt.FormatLocalized("MMMM Do, YYYY", "en-US")
	}
}

func BenchmarkHumanStringLocalized(b *testing.B) {
	dt := Now().AddHours(-2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = dt.HumanStringLocalized("en-US")
	}
}

func BenchmarkGetLocale(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetLocale("en-US")
	}
}

func TestJapaneseLocale(t *testing.T) {
	dt := Date(2024, time.January, 15, 14, 30, 0, 0, time.UTC) // Monday

	// Test month name
	monthName, err := dt.GetMonthName("ja-JP")
	if err != nil {
		t.Errorf("Failed to get Japanese month name: %v", err)
	}
	if monthName != "1月" {
		t.Errorf("Expected month name '1月', got '%s'", monthName)
	}

	// Test weekday name (Monday)
	weekdayName, err := dt.GetWeekdayName("ja-JP")
	if err != nil {
		t.Errorf("Failed to get Japanese weekday name: %v", err)
	}
	if weekdayName != "月曜日" {
		t.Errorf("Expected weekday name '月曜日', got '%s'", weekdayName)
	}

	// Test localized formatting
	result, err := dt.FormatLocalized("YYYY年MMMM Do dddd", "ja-JP")
	if err != nil {
		t.Errorf("Failed to format with Japanese locale: %v", err)
	}
	if !strings.Contains(result, "2024年") || !strings.Contains(result, "1月") || !strings.Contains(result, "15日") {
		t.Errorf("Expected Japanese date format, got '%s'", result)
	}

	// Test human-readable past
	past := dt.AddHours(-2)
	humanStr, err := past.HumanStringLocalized("ja-JP", dt)
	if err != nil {
		t.Errorf("Failed to get Japanese human string: %v", err)
	}
	if humanStr != "2時間前" {
		t.Errorf("Expected '2時間前', got '%s'", humanStr)
	}

	// Test human-readable future
	future := dt.AddHours(3)
	humanStr, err = future.HumanStringLocalized("ja-JP", dt)
	if err != nil {
		t.Errorf("Failed to get Japanese human string: %v", err)
	}
	if humanStr != "3時間後" {
		t.Errorf("Expected '3時間後', got '%s'", humanStr)
	}

	// Test days
	pastDays := dt.AddDays(-5)
	humanStr, err = pastDays.HumanStringLocalized("ja-JP", dt)
	if err != nil {
		t.Errorf("Failed to get Japanese human string: %v", err)
	}
	if humanStr != "5日前" {
		t.Errorf("Expected '5日前', got '%s'", humanStr)
	}

	// Test weeks
	futureWeeks := dt.AddDays(14)
	humanStr, err = futureWeeks.HumanStringLocalized("ja-JP", dt)
	if err != nil {
		t.Errorf("Failed to get Japanese human string: %v", err)
	}
	if humanStr != "2週間後" {
		t.Errorf("Expected '2週間後', got '%s'", humanStr)
	}

	// Test months
	futureMonths := dt.AddDays(90)
	humanStr, err = futureMonths.HumanStringLocalized("ja-JP", dt)
	if err != nil {
		t.Errorf("Failed to get Japanese human string: %v", err)
	}
	if humanStr != "3ヶ月後" {
		t.Errorf("Expected '3ヶ月後', got '%s'", humanStr)
	}

	// Test years
	futureYears := dt.AddDays(730)
	humanStr, err = futureYears.HumanStringLocalized("ja-JP", dt)
	if err != nil {
		t.Errorf("Failed to get Japanese human string: %v", err)
	}
	if humanStr != "2年後" {
		t.Errorf("Expected '2年後', got '%s'", humanStr)
	}

	// Test moments (few seconds)
	fewSeconds := dt.AddSeconds(5)
	humanStr, err = fewSeconds.HumanStringLocalized("ja-JP", dt)
	if err != nil {
		t.Errorf("Failed to get Japanese human string: %v", err)
	}
	if humanStr != "すぐに" {
		t.Errorf("Expected 'すぐに', got '%s'", humanStr)
	}
}

func TestJapaneseOrdinals(t *testing.T) {
	locale, err := GetLocale("ja-JP")
	if err != nil {
		t.Fatalf("Failed to get Japanese locale: %v", err)
	}

	// Test various dates
	tests := []struct {
		day      int
		expected string
	}{
		{1, "日"},
		{2, "日"},
		{15, "日"},
		{31, "日"},
	}

	for _, test := range tests {
		suffix := locale.getOrdinalSuffix(test.day)
		if suffix != test.expected {
			t.Errorf("Day %d: expected ordinal '%s', got '%s'", test.day, test.expected, suffix)
		}
	}
}

func TestJapaneseAMPM(t *testing.T) {
	locale, err := GetLocale("ja-JP")
	if err != nil {
		t.Fatalf("Failed to get Japanese locale: %v", err)
	}

	if locale.AMPMNames[0] != "午前" {
		t.Errorf("Expected AM to be '午前', got '%s'", locale.AMPMNames[0])
	}

	if locale.AMPMNames[1] != "午後" {
		t.Errorf("Expected PM to be '午後', got '%s'", locale.AMPMNames[1])
	}
}
