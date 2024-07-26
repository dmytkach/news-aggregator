package validator

import (
	"errors"
	"time"
)

// dateStartValidator checks the start date format.
type dateStartValidator struct {
	baseValidator
	dateStart string
}

// Validate checks that the start date is in the correct format and no later than today.
func (d dateStartValidator) Validate() error {
	if d.dateStart != "" {
		startDate, err := time.Parse(DateFormat, d.dateStart)
		if err != nil {
			//log.Println("Invalid start date format. Please use YYYY-MM-DD.")
			return errors.New("invalid start date format. Please use YYYY-MM-DD")
		}
		if startDate.After(time.Now()) {
			//log.Println("News for this period is not available.")
			return errors.New("news for this period is not available")
		}
	}
	return d.baseValidator.Validate()
}
