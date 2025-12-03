package validatorV2

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/aid297/aid/regexp"

	"github.com/aid297/aid/array/anyArrayV2"
)

type CheckerStringPtr struct {
	original *string
	required bool
	noZero   bool
	min      *int
	max      *int
	eq       *string
	notEq    *string
	in       []string
	notIn    []string
	regex    *string
	vType    string
}

func (my CheckerStringPtr) Check() error {
	if my.required {
		if my.original == nil {
			return ErrRequired
		}
	} else {
		if my.original == nil {
			return nil
		}
	}

	if my.noZero {
		if *my.original == "" {
			return ErrNoZero
		}
	} else {
		if *my.original == "" {
			return nil
		}
	}

	if my.min != nil {
		if utf8.RuneCountInString(*my.original) < *my.min {
			return fmt.Errorf(TooShort+"：%d", *my.min)
		}
	}

	if my.max != nil {
		if utf8.RuneCountInString(*my.original) > *my.max {
			return fmt.Errorf(TooLong+"：%d", *my.max)
		}
	}

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

	if len(my.in) > 0 {
		if !anyArrayV2.NewList(my.in).In(*my.original) {
			return fmt.Errorf("%w：%v", ErrIn, strings.Join(my.in, ","))
		}
	}

	if len(my.notIn) > 0 {
		if anyArrayV2.NewList(my.notIn).In(*my.original) {
			return fmt.Errorf("%w：%v", ErrNotIn, strings.Join(my.notIn, ","))
		}
	}

	if my.regex != nil {
		if regexp.APP.Regexp.New(*my.regex, regexp.TargetString(*my.original)).MatchFirst() == "" {
			return fmt.Errorf("%w：%v", ErrRegex, *my.regex)
		}
	}

	return nil
}
