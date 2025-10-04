package main

import (
	"fmt"
	"time"

	chronogo "github.com/coredds/chronogo"
)

func main() {
	fmt.Println("chronogo Localization Demo")
	fmt.Println("==========================")

	// Create a sample datetime
	dt := chronogo.Date(2024, time.January, 15, 14, 30, 0, 0, time.UTC)
	fmt.Printf("Base DateTime: %s\n\n", dt.String())

	// Demonstrate localized formatting
	fmt.Println("Localized Formatting Examples:")
	fmt.Println("------------------------------")

	locales := []string{"en-US", "es-ES", "fr-FR", "de-DE", "zh-Hans", "pt-BR"}
	patterns := []string{
		"dddd, MMMM Do YYYY",
		"Do de MMMM de YYYY", // Spanish/Portuguese pattern
		"dddd Do MMMM YYYY",  // French pattern
		"dddd, Do MMMM YYYY", // German pattern
		"YYYY年MMMM Do dddd",  // Chinese pattern
	}

	for _, locale := range locales {
		var pattern string
		switch locale {
		case "es-ES", "pt-BR":
			pattern = patterns[1]
		case "fr-FR":
			pattern = patterns[2]
		case "de-DE":
			pattern = patterns[3]
		case "zh-Hans":
			pattern = patterns[4]
		default:
			pattern = patterns[0]
		}

		result, err := dt.FormatLocalized(pattern, locale)
		if err != nil {
			fmt.Printf("Error formatting for %s: %v\n", locale, err)
			continue
		}
		fmt.Printf("%-8s: %s\n", locale, result)
	}

	fmt.Println("\nHuman-Readable Time Differences:")
	fmt.Println("--------------------------------")

	// Create some time differences
	now := chronogo.Now()
	times := []struct {
		name string
		dt   chronogo.DateTime
	}{
		{"2 hours ago", now.AddHours(-2)},
		{"30 minutes ago", now.AddMinutes(-30)},
		{"in 1 day", now.AddDays(1)},
		{"in 3 weeks", now.AddDays(21)},
		{"few seconds ago", now.AddSeconds(-5)},
	}

	for _, timeExample := range times {
		fmt.Printf("\n%s:\n", timeExample.name)
		for _, locale := range locales {
			result, err := timeExample.dt.HumanStringLocalized(locale)
			if err != nil {
				fmt.Printf("  %-8s: Error - %v\n", locale, err)
				continue
			}
			fmt.Printf("  %-8s: %s\n", locale, result)
		}
	}

	fmt.Println("\nMonth and Weekday Names:")
	fmt.Println("------------------------")

	// Demonstrate month names
	fmt.Println("June in different locales:")
	juneDate := chronogo.Date(2024, time.June, 15, 0, 0, 0, 0, time.UTC)
	for _, locale := range locales {
		monthName, err := juneDate.GetMonthName(locale)
		if err != nil {
			fmt.Printf("  %-8s: Error - %v\n", locale, err)
			continue
		}
		fmt.Printf("  %-8s: %s\n", locale, monthName)
	}

	// Demonstrate weekday names (Monday)
	fmt.Println("\nMonday in different locales:")
	mondayDate := chronogo.Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC) // This is a Monday
	for _, locale := range locales {
		weekdayName, err := mondayDate.GetWeekdayName(locale)
		if err != nil {
			fmt.Printf("  %-8s: Error - %v\n", locale, err)
			continue
		}
		fmt.Printf("  %-8s: %s\n", locale, weekdayName)
	}

	fmt.Println("\nDefault Locale Management:")
	fmt.Println("-------------------------")

	// Demonstrate default locale functionality
	originalDefault := chronogo.GetDefaultLocale()
	fmt.Printf("Original default locale: %s\n", originalDefault)

	// Set to Spanish
	err := chronogo.SetDefaultLocale("es-ES")
	if err != nil {
		fmt.Printf("Error setting default locale: %v\n", err)
	} else {
		fmt.Printf("Changed default locale to: %s\n", chronogo.GetDefaultLocale())

		// Use default locale formatting
		result := dt.FormatLocalizedDefault("dddd, Do de MMMM")
		fmt.Printf("Default format result: %s\n", result)

		// Use default locale for human strings
		pastTime := now.AddHours(-3)
		humanResult := pastTime.HumanStringLocalizedDefault()
		fmt.Printf("Default human string: %s\n", humanResult)
	}

	// Reset to original
	chronogo.SetDefaultLocale(originalDefault)
	fmt.Printf("Reset default locale to: %s\n", chronogo.GetDefaultLocale())

	fmt.Println("\nAvailable Locales:")
	fmt.Println("------------------")
	availableLocales := chronogo.GetAvailableLocales()
	for i, locale := range availableLocales {
		localeInfo, _ := chronogo.GetLocale(locale)
		fmt.Printf("%d. %s - %s\n", i+1, locale, localeInfo.Name)
	}
}
