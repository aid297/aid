package validations

import (
	"sync"
)

type (
	Validation struct {
		data                  map[string]GlobalExCheckFunc
		defaultSliceSplitChar string
		defaultErrorSplitChar string
	}
	GlobalExCheckFunc func(fieldName string, origin any) (err error)
	ExCheckFunc       func(origin any) (errs error)

	Validator interface {
		DefaultSliceSplitChar(char string) Validator
		DefaultErrorSplitChar(char string) Validator
		RegisterExFn(key string, fn GlobalExCheckFunc) Validator
		GetExFn(key string) GlobalExCheckFunc
		Checker(data any) Checker
	}
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
			defaultSliceSplitChar: ",",
			defaultErrorSplitChar: "<br />",
		}
	})
	return validatorExIns
}

func (*Validation) DefaultSliceSplitChar(char string) Validator {
	validatorExIns.defaultSliceSplitChar = char
	return validatorExIns
}

func (*Validation) DefaultErrorSplitChar(char string) Validator {
	validatorExIns.defaultErrorSplitChar = char
	return validatorExIns
}

func (*Validation) RegisterExFn(key string, fn GlobalExCheckFunc) Validator {
	validatorExIns.data[key] = fn
	return validatorExIns
}

func (*Validation) GetExFn(key string) GlobalExCheckFunc { return validatorExIns.data[key] }

func (*Validation) Checker(data any) Checker { return NewCheck(data) }
