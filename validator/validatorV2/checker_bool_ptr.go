package validatorV2

import (
	"fmt"
)

type CheckerBoolPtr struct {
	original *bool
	eq       *bool
	notEq    *bool
}

func (my CheckerBoolPtr) Check() error {
	if my.eq != nil {
		if *my.original != *my.eq {
			return fmt.Errorf("%w：%v", ErrEq, *my.eq)
		}
	}

	if my.notEq != nil {
		if *my.original == *my.notEq {
			return fmt.Errorf("%w：%v", ErrNotEq, *my.notEq)
		}
	}

	return nil
}
