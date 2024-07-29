package validator

import (
	"errors"
	"testing"
)

func TestDateStartValidator_Validate(t *testing.T) {
	tests := []struct {
		name      string
		dateStart string
		expected  error
	}{
		{"Valid start date", "2023-01-01", nil},
		{"Empty start date", "", nil},
		{"Invalid start date format", "01-01-2023", errors.New("invalid start date format. Please use YYYY-MM-DD")},
		{"Invalid start date string", "invalid-date", errors.New("invalid start date format. Please use YYYY-MM-DD")},
		{"Future start date", "2999-12-31", errors.New("news for this period is not available")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := dateStartValidator{
				dateStart: tt.dateStart,
			}
			result := validator.Validate()
			if tt.expected != nil {
				if result == nil || result.Error() != tt.expected.Error() {
					t.Errorf("Expected %v, but got %v", tt.expected, result)
				}
			}
		})
	}
}
