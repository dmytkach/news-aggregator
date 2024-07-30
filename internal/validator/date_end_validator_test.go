package validator

import (
	"errors"
	"testing"
)

func TestDateEndValidator_Validate(t *testing.T) {
	tests := []struct {
		name     string
		dateEnd  string
		expected error
	}{
		{"Valid end date", "2023-12-31", nil},
		{"Empty end date", "", nil},
		{"Invalid end date format", "31-12-2023", errors.New("invalid end date format. Please use YYYY-MM-DD")},
		{"Invalid end date string", "invalid-date", errors.New("invalid end date format. Please use YYYY-MM-DD")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := dateEndValidator{
				dateEnd: tt.dateEnd,
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
