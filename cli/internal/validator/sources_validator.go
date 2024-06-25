package validator

import (
	"fmt"
	"strings"
)

// sourceValidator checks the sources field
type sourceValidator struct {
	baseValidator
	availableSources []string
	sources          string
}

func (v sourceValidator) Validate() bool {
	sourcesList := strings.Split(v.sources, ",")
	if len(sourcesList) == 0 {
		println("Please provide at least one source using the --sources flag.")
		return false
	}
	for _, source := range sourcesList {
		if !contains(v.availableSources, source) {
			fmt.Println("Source not available:", source)
			return false
		}
	}
	return v.baseValidator.Validate()
}

// contains check elements existing in slice.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.Contains(s, item) {
			return true
		}
	}
	return false
}
