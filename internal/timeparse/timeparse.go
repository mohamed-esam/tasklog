package timeparse

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var (
	// Regex patterns for parsing time formats
	// Supports: 2h 30m, 2.5h, 150m, 2h30m, etc.
	hoursMinutesPattern = regexp.MustCompile(`(?i)^(\d+)\s*h(?:ours?)?\s*(\d+)\s*m(?:in(?:utes?)?)?$`)
	hoursOnlyPattern    = regexp.MustCompile(`(?i)^(\d+(?:\.\d+)?)\s*h(?:ours?)?$`)
	minutesOnlyPattern  = regexp.MustCompile(`(?i)^(\d+)\s*m(?:in(?:utes?)?)?$`)
)

// Parse parses a time string and returns the duration in seconds
// Supports formats: "2h 30m", "2.5h", "150m", "2h30m"
// Rounds to the nearest 5 minutes
func Parse(input string) (int, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return 0, fmt.Errorf("empty time input")
	}

	var totalMinutes float64

	// Try matching "Xh Ym" or "XhYm" format
	if matches := hoursMinutesPattern.FindStringSubmatch(input); matches != nil {
		hours, err := strconv.Atoi(matches[1])
		if err != nil {
			return 0, fmt.Errorf("invalid hours value: %w", err)
		}
		minutes, err := strconv.Atoi(matches[2])
		if err != nil {
			return 0, fmt.Errorf("invalid minutes value: %w", err)
		}
		totalMinutes = float64(hours*60 + minutes)
	} else if matches := hoursOnlyPattern.FindStringSubmatch(input); matches != nil {
		// Try matching "X.Yh" or "Xh" format
		hours, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid hours value: %w", err)
		}
		totalMinutes = hours * 60
	} else if matches := minutesOnlyPattern.FindStringSubmatch(input); matches != nil {
		// Try matching "Xm" format
		minutes, err := strconv.Atoi(matches[1])
		if err != nil {
			return 0, fmt.Errorf("invalid minutes value: %w", err)
		}
		totalMinutes = float64(minutes)
	} else {
		return 0, fmt.Errorf("invalid time format: %s (expected formats: 2h 30m, 2.5h, 150m)", input)
	}

	if totalMinutes <= 0 {
		return 0, fmt.Errorf("time must be positive")
	}

	// Round to nearest 5 minutes
	roundedMinutes := roundToNearest5(totalMinutes)

	return int(roundedMinutes * 60), nil
}

// roundToNearest5 rounds a number to the nearest 5
func roundToNearest5(minutes float64) float64 {
	return math.Round(minutes/5) * 5
}

// Format formats seconds into a human-readable time string
func Format(seconds int) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60

	if hours > 0 && minutes > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dm", minutes)
}

// Validate checks if a time string is valid
func Validate(input string) error {
	_, err := Parse(input)
	return err
}
