package validatorV2

import (
	"fmt"
	"strings"

	"github.com/aid297/aid/array/anyArrayV2"
	"github.com/spf13/cast"
)

type CheckerUint64 struct {
	original uint64
	required bool
	noZero   bool
	eq       *uint64
	notEq    *uint64
	min      *uint64
	max      *uint64
	in       []uint64
	notIn    []uint64
}

func (my CheckerUint64) Check() error {
	if my.required {
		if my.original == 0 {
			return ErrRequired
		}
	}

	if my.noZero {
		if my.original == 0 {
			return ErrNoZero
		}
	}

	if my.min != nil {
		if my.original < *my.min {
			return fmt.Errorf("长度不能小于：%d", *my.min)
		}
	}

	if my.max != nil {
		if my.original > *my.max {
			return fmt.Errorf("长度不能大于：%d", *my.max)
		}
	}

	if my.eq != nil {
		if my.original != *my.eq {
			return fmt.Errorf("%w:%v", ErrEq, *my.eq)
		}
	}

	if my.notEq != nil {
		if my.original == *my.notEq {
			return fmt.Errorf("%w:%v", ErrNotEq, *my.notEq)
		}
	}

	if len(my.in) > 0 {
		if !anyArrayV2.NewList(my.in).In(my.original) {
			return fmt.Errorf("%w:%v", ErrIn, strings.Join(cast.ToStringSlice(my.in), ","))
		}
	}

	if len(my.notIn) > 0 {
		if anyArrayV2.NewList(my.notIn).In(my.original) {
			return fmt.Errorf("%w:%v", ErrNotIn, strings.Join(cast.ToStringSlice(my.in), ","))
		}
	}

	return nil
}
