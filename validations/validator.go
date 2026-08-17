package validations

import "sync"

type (
	Validation struct {
		data                  map[string]GlobalExCheckFunc
		defaultSliceSplitChar SliceSplitCharTag
		defaultErrorSplitChar ErrorSplitCharTag
	}
	GlobalExCheckFunc func(fieldName string, origin any) (err error)
	ExCheckFunc       func(origin any) (errs error)

	Validator interface {
		DefaultSliceSplitChar(char SliceSplitCharTag) Validator
		DefaultErrorSplitChar(char ErrorSplitCharTag) Validator
		RegisterExFn(key string, fn GlobalExCheckFunc) Validator
		GetExFn(key string) GlobalExCheckFunc
		Checker(data any) Checker
	}

	SliceSplitCharTag string
	ErrorSplitCharTag string
)

const (
	SliceSplitChar        SliceSplitCharTag = ","
	WebErrorSplitChar     ErrorSplitCharTag = "<br />"
	ConsoleErrorSplitChar ErrorSplitCharTag = "\n"
)

var (
	_               Validator = (*Validation)(nil)
	validatorExOnce sync.Once
	validatorExIns  *Validation
)

func OnceValidator() Validator {
	validatorExOnce.Do(func() {
		validatorExIns = &Validation{
			data:                  make(map[string]GlobalExCheckFunc),
			defaultSliceSplitChar: SliceSplitChar,
			defaultErrorSplitChar: ConsoleErrorSplitChar,
		}
	})
	return validatorExIns
}

func (*Validation) DefaultSliceSplitChar(char SliceSplitCharTag) Validator {
	validatorExIns.defaultSliceSplitChar = char
	return validatorExIns
}

func (*Validation) DefaultErrorSplitChar(char ErrorSplitCharTag) Validator {
	validatorExIns.defaultErrorSplitChar = char
	return validatorExIns
}

func (*Validation) RegisterExFn(key string, fn GlobalExCheckFunc) Validator {
	validatorExIns.data[key] = fn
	return validatorExIns
}

func (*Validation) GetExFn(key string) GlobalExCheckFunc { return validatorExIns.data[key] }

func (*Validation) Checker(data any) Checker { return NewCheck(data) }
