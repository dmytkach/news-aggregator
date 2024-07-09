package validator

import (
	"testing"
)

func TestDateEndValidator_Validate(t *testing.T) {
	tests := []struct {
		name     string
		dateEnd  string
		expected bool
	}{
		{"Valid end date", "2023-12-31", true},
		{"Empty end date", "", true},
		{"Invalid end date format", "31-12-2023", false},
		{"Invalid end date string", "invalid-date", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := dateEndValidator{
				dateEnd: tt.dateEnd,
			}
			result := validator.Validate()
			if result != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}
