package validatorV3

import (
	"sync"
)

type Validator struct{ data map[string]func(any) error }

var (
	validatorExOnce sync.Once
	validatorExIns  *Validator
)

func (*Validator) Once() *Validator {
	validatorExOnce.Do(func() { validatorExIns = &Validator{data: make(map[string]func(any) (err error))} })
	return validatorExIns
}

func (*Validator) RegisterExFn(key string, fn func(any) (err error)) *Validator {
	validatorExIns.data[key] = fn
	return validatorExIns
}

func (*Validator) GetExFn(key string) func(any) (err error) { return validatorExIns.data[key] }

func (*Validator) Checker(data any) Checker { return NewCheck(data) }
