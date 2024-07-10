package validator

import (
	"testing"
	"time"
)

func TestValidator_Validate(t *testing.T) {
	tests := []struct {
		name             string
		sources          string
		availableSources []string
		dateStart        string
		dateEnd          string
		expected         bool
	}{
		{"All valid inputs", "source1", []string{"source1", "source2"}, "2023-01-01", "2023-12-31", true},
		{"Invalid source", "source2", []string{"Source 1"}, "2023-01-01", "2023-12-31", false},
		{"Invalid start date format", "source1", []string{"Source 1"}, "01-01-2023", "2023-12-31", false},
		{"Future start date", "source1", []string{"Source 1"}, time.Now().AddDate(1, 0, 0).Format(DateFormat), "2023-12-31", false},
		{"Invalid end date format", "source1", []string{"Source 1"}, "2023-01-01", "31-12-2023", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := validator{
				Sources:          tt.sources,
				AvailableSources: tt.availableSources,
				DateStart:        tt.dateStart,
				DateEnd:          tt.dateEnd,
			}
			result := validator.Validate()
			if result != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}
