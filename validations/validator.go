package validations

import (
	"sync"
)

type (
	Validation struct{ data map[string]ExCheckFn }
	ExCheckFn  func(fieldName string, origin any) (err error)

	Validator interface {
		RegisterExFn(key string, fn ExCheckFn) Validator
		GetExFn(key string) ExCheckFn
		Checker(data any) Checker
	}
)

var (
	_               Validator = (*Validation)(nil)
	validatorExOnce sync.Once
	validatorExIns  *Validation
)

func OnceValidator() Validator {
	validatorExOnce.Do(func() { validatorExIns = &Validation{data: make(map[string]ExCheckFn)} })
	return validatorExIns
}

func (*Validation) RegisterExFn(key string, fn ExCheckFn) Validator {
	validatorExIns.data[key] = fn
	return validatorExIns
}

func (*Validation) GetExFn(key string) ExCheckFn { return validatorExIns.data[key] }

func (*Validation) Checker(data any) Checker { return NewCheck(data) }
