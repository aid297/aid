package validations

import (
	"sync"
)

type (
	Validation        struct{ data map[string]GlobalExCheckFunc }
	GlobalExCheckFunc func(fieldName string, origin any) (err error)
	ExCheckFunc       func(origin any) (errs error)

	Validator interface {
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
	validatorExOnce.Do(func() { validatorExIns = &Validation{data: make(map[string]GlobalExCheckFunc)} })
	return validatorExIns
}

func (*Validation) RegisterExFn(key string, fn GlobalExCheckFunc) Validator {
	validatorExIns.data[key] = fn
	return validatorExIns
}

func (*Validation) GetExFn(key string) GlobalExCheckFunc { return validatorExIns.data[key] }

func (*Validation) Checker(data any) Checker { return NewCheck(data) }
