package validator

type validator struct {
	Sources          string
	AvailableSources map[string]string
	DateStart        string
	DateEnd          string
}

func NewValidator(sources string, availableSources map[string]string, dateStart, dateEnd string) ValidatingComponent {
	return &validator{
		Sources:          sources,
		AvailableSources: availableSources,
		DateStart:        dateStart,
		DateEnd:          dateEnd,
	}
}

// ValidatingComponent defines the contract for components responsible for performing validation.
type ValidatingComponent interface {
	Validate() bool
}

// Validate all implementations with the chain of responsibility pattern.
func (v validator) Validate() bool {
	sourceValidator := &sourceValidator{sources: v.Sources, availableSources: v.AvailableSources}
	dateStartValidator := &dateStartValidator{dateStart: v.DateStart}
	dateEndValidator := &dateEndValidator{dateEnd: v.DateEnd}

	sourceValidator.SetNext(dateStartValidator)
	dateStartValidator.SetNext(dateEndValidator)

	return sourceValidator.Validate()
}
