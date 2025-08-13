//go:build go1.18
// +build go1.18

package chronogo

import "testing"

// FuzzParse fuzzes the Parse function to ensure it doesn't panic on random inputs.
func FuzzParse(f *testing.F) {
	seeds := []string{
		"2023-12-25T15:30:45Z",
		"2023-12-25 15:30:45",
		"2023-12-25",
		"2023/12/25",
		"2023-1-2 3:04:05",
		"20231225",
		"15:30:45",
		"1640995200",
		"",
		"invalid",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		_, _ = Parse(s)
	})
}
