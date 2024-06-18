package validator

import (
	"testing"
)

func TestDateStartValidator_Validate(t *testing.T) {
	tests := []struct {
		name      string
		dateStart string
		expected  bool
	}{
		{"Valid start date", "2023-01-01", true},
		{"Empty start date", "", true},
		{"Invalid start date format", "01-01-2023", false},
		{"Invalid start date string", "invalid-date", false},
		{"Future start date", "2999-12-31", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := dateStartValidator{
				dateStart: tt.dateStart,
			}
			result := validator.Validate()
			if result != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}
