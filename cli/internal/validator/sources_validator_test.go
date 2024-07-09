package validator

import "testing"

func TestSourceValidator_Validate(t *testing.T) {
	tests := []struct {
		name             string
		sources          string
		availableSources map[string]string
		expected         bool
	}{
		{"No sources", "", map[string]string{}, false},
		{"Valid source", "source1", map[string]string{"source1": "Source 1"}, true},
		{"Invalid source", "source2", map[string]string{"source1": "Source 1"}, false},
		{"Mixed sources", "source1.source2", map[string]string{"source1": "Source 1"}, false},
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
