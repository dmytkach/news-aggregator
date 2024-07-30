package validator

import (
	"errors"
	"testing"
)

func TestSourceValidator_Validate(t *testing.T) {
	tests := []struct {
		name             string
		sources          string
		availableSources []string
		expected         error
	}{
		{"No sources", "", []string{}, errors.New("not found available sources")},
		{"Valid source", "source1", []string{"source1"}, nil},
		{"Invalid source", "source2", []string{"source1"}, errors.New("source not available: source2")},
		{"Mixed sources", "source1.source2", []string{"source1"}, errors.New("source not available: source1.source2")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := sourceValidator{
				sources:          tt.sources,
				availableSources: tt.availableSources,
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
