package validator

type Validator struct {
	Sources          string
	AvailableSources []string
	DateStart        string
	DateEnd          string
}

// ValidatingComponent defines the contract for components responsible for performing validation.
type ValidatingComponent interface {
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
