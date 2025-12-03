package validatorV2

import (
	"fmt"
	"strings"

	"github.com/aid297/aid/array/anyArrayV2"
	"github.com/spf13/cast"
)

type CheckerSlice struct {
	original []any
	required bool
	noZero   bool
	eq       *int
	notEq    *int
	min      *int
	max      *int
	in       []any
	notIn    []any
}

func (my CheckerSlice) Check() error {
	if my.required {
		if len(my.original) == 0 {
			return ErrRequired
		}
	} else {
		if len(my.original) == 0 {
			return nil
		}
	}

	if my.noZero {
		if len(my.original) == 0 {
			return ErrNoZero
		}
	} else {
		if len(my.original) == 0 {
			return nil
		}
	}

	if my.min != nil {
		if len(my.original) < *my.min {
			return fmt.Errorf(TooShort+"：%d", *my.min)
		}
	}

	if my.max != nil {
		if len(my.original) > *my.max {
			return fmt.Errorf(TooLong+"：%d", *my.max)
		}
	}

	if my.eq != nil {
		if len(my.original) != *my.eq {
			return fmt.Errorf("%w：%v", ErrEq, *my.eq)
		}
	}

	if my.notEq != nil {
		if len(my.original) == *my.notEq {
			return fmt.Errorf("%w：%v", ErrNotEq, *my.notEq)
		}
	}

	if len(my.in) > 0 {
		if !anyArrayV2.NewList(my.in).In(my.original...) {
			return fmt.Errorf("%w：%v", ErrIn, strings.Join(cast.ToStringSlice(my.in), ","))
		}
	}

	if len(my.notIn) > 0 {
		if anyArrayV2.NewList(my.notIn).In(my.original...) {
			return fmt.Errorf("%w：%v", ErrNotIn, strings.Join(cast.ToStringSlice(my.in), ","))
		}
	}

	return nil
}
