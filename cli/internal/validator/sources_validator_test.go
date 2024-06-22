package validator

import "testing"

func TestSourceValidator_Validate(t *testing.T) {
	tests := []struct {
		name             string
		sources          string
		availableSources []string
		expected         bool
	}{
		{"No sources", "", []string{}, false},
		{"Valid source", "source1", []string{"source1"}, true},
		{"Invalid source", "source2", []string{"source1"}, false},
		{"Mixed sources", "source1.source2", []string{"source1"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := sourceValidator{
				sources:          tt.sources,
				availableSources: tt.availableSources,
			}
			result := validator.Validate()
			if result != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}
