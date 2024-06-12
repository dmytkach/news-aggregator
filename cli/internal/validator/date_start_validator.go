package validator

import (
	"time"
)

type DateStartValidator struct {
	BaseValidator
	DateStart string
}

func (d DateStartValidator) Validate() bool {
	if d.DateStart != "" {
		startDate, err := time.Parse("2006-01-02", d.DateStart)
		if err != nil {
			println("Invalid start date format. Please use YYYY-MM-DD.")
			return false
		}
		if startDate.After(time.Now()) {
			println("News for this period is not available.")
			return false
		}
	}
	return d.BaseValidator.Validate()
}
