package validatorV3

type (
	ValidatorAttributer interface{ Register(validator *Validator) }

	AttrData struct{ data any }
)

func (AttrData) Set(data any) ValidatorAttributer { return AttrData{data} }
func (my AttrData) Register(validator *Validator) { validator.data = my.data }
