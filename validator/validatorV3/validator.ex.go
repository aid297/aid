package validatorV3

import (
	"sync"
)

type Validator struct{ data map[string]func(any) error }

var (
	validatorExOnce sync.Once
	validatorExIns  *Validator
)

func (*Validator) Ins() *Validator {
	validatorExOnce.Do(func() { validatorExIns = new(Validator) })
	return validatorExIns
}

func (*Validator) Register(key string, fn func(any) error) *Validator {
	validatorExIns.data[key] = fn
	return validatorExIns
}

func (*Validator) Get(key string) func(any) error {
	return validatorExIns.data[key]
}

func (*Validator) Checker(data any) Checker { return APP.Validator.Ins().Checker(data) }
