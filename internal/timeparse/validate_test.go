package timeparse

import "testing"

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"valid time 2h30m", "2h30m", false},
		{"valid time 1.5h", "1.5h", false},
		{"valid time 90m", "90m", false},
		{"invalid format", "invalid", true},
		{"empty string", "", true},
		{"negative", "-2h", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.input)
			if tt.wantError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestRoundToNearest5(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{0, 0},
		{1, 0},
		{2, 0},
		{2.5, 5},
		{3, 5},
		{7, 5},
		{8, 10},
		{12, 10},
		{12.5, 15},
		{13, 15},
		{152, 150},
		{153, 155},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := roundToNearest5(tt.input)
			if result != tt.expected {
				t.Errorf("roundToNearest5(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
