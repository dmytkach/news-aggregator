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

// Validate checks that the given sources exist.
func (v sourceValidator) Validate() bool {
	sourcesList := strings.Split(v.sources, ",")
	if len(sourcesList) == 0 {
		println("Please provide at least one source using the --sources flag.")
		return false
	}
	if len(v.availableSources) == 0 {
		println("Not found available sources.")
		return false
	}
	for _, source := range sourcesList {
		if !v.exist(source) {
			fmt.Println("Source not available:", source)
			return false
		}
	}
	return v.baseValidator.Validate()
}

// exist check elements existing in slice.
func (v sourceValidator) exist(source string) bool {
	for _, s := range v.availableSources {
		if strings.Contains(s, strings.ToLower(source)) {
			return true
		}
	}
	return false
}
