package validator

import (
	"log"
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
	if len(v.availableSources) == 0 {
		log.Println("Not found available sources.")
		return false
	}
	for _, source := range sourcesList {
		if strings.TrimSpace(source) == "" {
			log.Println("Please provide at least one source.")
			return false
		}
		if !v.exist(source) {
			log.Println("Source not available:", source)
			return false
		}
	}
	return v.baseValidator.Validate()
}

// exist check elements existing in slice.
func (v sourceValidator) exist(source string) bool {
	for _, s := range v.availableSources {
		if strings.EqualFold(s, source) {
			return true
		}
	}
	return false
}
