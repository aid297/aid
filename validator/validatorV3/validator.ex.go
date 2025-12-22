package validatorV3

import (
	"sync"
)

type ValidatorEx struct{ data map[string]func(any) error }

var (
	validatorExOnce sync.Once
	validatorExIns  *ValidatorEx
)

func (*ValidatorEx) New() *ValidatorEx {
	validatorExOnce.Do(func() { validatorExIns = new(ValidatorEx) })
	return validatorExIns
}

func (*ValidatorEx) Register(key string, fn func(any) error) *ValidatorEx {
	validatorExIns.data[key] = fn
	return validatorExIns
}

func (*ValidatorEx) Get(key string) func(any) error {
	return validatorExIns.data[key]
}
