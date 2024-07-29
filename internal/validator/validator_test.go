package validator

import (
	"errors"
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
		expected         error
	}{
		{"All valid inputs", "source1", []string{"source1", "source2"}, "2023-01-01", "2023-12-31", nil},
		{"Invalid source", "source2", []string{"source1"}, "2023-01-01", "2023-12-31", errors.New("source not available: source2")},
		{"Invalid start date format", "source1", []string{"source1"}, "01-01-2023", "2023-12-31", errors.New("invalid start date format. Please use YYYY-MM-DD")},
		{"Future start date", "source1", []string{"source1"}, time.Now().AddDate(1, 0, 0).Format(DateFormat), "2023-12-31", errors.New("news for this period is not available")},
		{"Invalid end date format", "source1", []string{"source1"}, "2023-01-01", "31-12-2023", errors.New("invalid end date format. Please use YYYY-MM-DD")},
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
			if tt.expected != nil {
				if result == nil || result.Error() != tt.expected.Error() {
					t.Errorf("Expected %v, but got %v", tt.expected, result)
				}
			}
		})
	}
}
