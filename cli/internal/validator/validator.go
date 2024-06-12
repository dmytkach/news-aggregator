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
	sourceValidator := &SourceValidator{Sources: v.Sources, AvailableSources: v.AvailableSources}
	dateStartValidator := &DateStartValidator{DateStart: v.DateStart}
	dateEndValidator := &DateEndValidator{DateEnd: v.DateEnd}

	sourceValidator.SetNext(dateStartValidator)
	dateStartValidator.SetNext(dateEndValidator)

	return sourceValidator.Validate()
}
