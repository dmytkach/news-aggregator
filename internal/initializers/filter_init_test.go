package initializers

import (
	"testing"
)

func TestInitializeFilters(t *testing.T) {
	keywords := "keyword1,keyword2"
	dateStart := "2024-06-01"
	dateEnd := "2024-06-30"

	filters := InitializeFilters(&keywords, &dateStart, &dateEnd)

	if len(filters) != 3 {
		t.Errorf("Expected 3 filters, got %d", len(filters))
	}

	dateEnd = ""
	filters = InitializeFilters(&keywords, &dateStart, &dateEnd)

	if len(filters) != 2 {
		t.Errorf("Expected 2 filters, got %d", len(filters))
	}

	keywords = ""
	filters = InitializeFilters(&keywords, &dateStart, &dateEnd)

	if len(filters) != 1 {
		t.Errorf("Expected 1 filter, got %d", len(filters))
	}
}

func TestConvertKeywords(t *testing.T) {
	keywords := "keyword1,keyword2"

	result := convertKeywords(&keywords)

	if result == nil {
		t.Errorf("Expected non-nil result for valid input")
	}

	if len(result.Keywords) != 2 {
		t.Errorf("Expected 2 keywords, got %d", len(result.Keywords))
	}

	keywords = ""

	result = convertKeywords(&keywords)

	if result != nil {
		t.Errorf("Expected nil result for empty input")
	}
}

func TestConvertDateStart(t *testing.T) {
	dateStart := "2024-06-01"

	result := convertDateStart(&dateStart)

	if result == nil {
		t.Errorf("Expected non-nil result for valid input")
	}

	dateStart = ""

	result = convertDateStart(&dateStart)

	if result != nil {
		t.Errorf("Expected nil result for empty input")
	}
}

func TestConvertDateEnd(t *testing.T) {
	dateEnd := "2024-06-30"

	result := convertDateEnd(&dateEnd)

	if result == nil {
		t.Errorf("Expected non-nil result for valid input")
	}
	dateEnd = ""

	result = convertDateEnd(&dateEnd)

	if result != nil {
		t.Errorf("Expected nil result for empty input")
	}
}
