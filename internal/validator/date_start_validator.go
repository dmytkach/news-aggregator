package validator

import (
	"log"
	"time"
)

// dateStartValidator checks the start date format.
type dateStartValidator struct {
	baseValidator
	dateStart string
}

// Validate checks that the start date is in the correct format and no later than today.
func (d dateStartValidator) Validate() bool {
	if d.dateStart != "" {
		startDate, err := time.Parse(DateFormat, d.dateStart)
		if err != nil {
			log.Println("Invalid start date format. Please use YYYY-MM-DD.")
			return false
		}
		if startDate.After(time.Now()) {
			log.Println("News for this period is not available.")
			return false
		}
	}
	return d.baseValidator.Validate()
}
