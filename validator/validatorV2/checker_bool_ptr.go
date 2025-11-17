package validatorV2

import (
	"fmt"
)

type CheckerBoolPtr struct {
	original *bool
	required bool
	noZero   bool
	eq       *bool
	notEq    *bool
}

func (my CheckerBoolPtr) Check() error {
	if my.required {
		if my.original == nil {
			return ErrRequired
		}
	}

	if my.noZero {
		if !*my.original {
			return ErrNoZero
		}
	}

	if my.eq != nil {
		if *my.original != *my.eq {
			return fmt.Errorf("%w:%v", ErrEq, *my.eq)
		}
	}

	if my.notEq != nil {
		if *my.original == *my.notEq {
			return fmt.Errorf("%w:%v", ErrNotEq, *my.notEq)
		}
	}

	return nil
}
