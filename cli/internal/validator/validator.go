package validator

type Validator struct {
	Sources          string
	AvailableSources []string
	DateStart        string
	DateEnd          string
}

type Validatable interface {
	Validate() bool
}

func (v Validator) Validate() bool {
	sourceValidator := sourceValidator{sources: v.Sources, availableSources: v.AvailableSources}
	dateStartValidator := dateStartValidator{dateStart: v.DateStart}
	dateEndValidator := dateEndValidator{dateEnd: v.DateEnd}

	sourceValidator.SetNext(dateStartValidator)
	dateStartValidator.SetNext(dateEndValidator)

	return sourceValidator.Validate()
}
