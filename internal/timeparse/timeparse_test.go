package timeparse

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedSecs  int
		expectedError bool
	}{
		{"hours and minutes", "2h 30m", 9000, false},
		{"hours and minutes no space", "2h30m", 9000, false},
		{"decimal hours", "2.5h", 9000, false},
		{"minutes only", "150m", 9000, false},
		{"hours only", "2h", 7200, false},
		{"round up to 5", "2h 32m", 9000, false},      // 152m -> 150m (rounded down to nearest 5)
		{"round down to 5", "2h 27m", 8700, false},    // 147m -> 145m
		{"single minute rounds to 5", "1m", 0, false}, // 1m -> 0m (rounds down)
		{"3 minutes rounds to 5", "3m", 300, false},   // 3m -> 5m
		{"7 minutes rounds to 5", "7m", 300, false},   // 7m -> 5m
		{"case insensitive hours", "2H 30M", 9000, false},
		{"case insensitive minutes", "150M", 9000, false},
		{"hours with full word", "2 hours 30 minutes", 9000, false},
		{"empty string", "", 0, true},
		{"invalid format", "abc", 0, true},
		{"negative time", "-2h", 0, true},
		{"zero time", "0h", 0, true},
		{"zero minutes", "0m", 0, true},
		{"whitespace only", "   ", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != tt.expectedSecs {
					t.Errorf("expected %d seconds, got %d", tt.expectedSecs, result)
				}
			}
		})
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		name     string
		seconds  int
		expected string
	}{
		{"hours and minutes", 9000, "2h 30m"},
		{"hours only", 7200, "2h"},
		{"minutes only", 1800, "30m"},
		{"zero", 0, "0m"},
		{"5 minutes", 300, "5m"},
		{"1 hour 5 minutes", 3900, "1h 5m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Format(tt.seconds)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
