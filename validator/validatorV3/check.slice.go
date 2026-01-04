package validatorV3

import (
	"fmt"
	"strings"

	"github.com/aid297/aid/array/anyArrayV2"
)

// checkSlice 检查数组、切片，支持：required、not-empty、[array|slice]、min>、min>=、max<、max<=、size=、size!=、ex:
func (my FieldInfo) checkSlice() FieldInfo {
	var (
		rules          = anyArrayV2.NewList(my.VRuleTags)
		min, max, size *int
		include, eq    bool
		value          []any
		ok             bool
	)

	if value, ok = my.Value.([]any); !ok {
		my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：数组", my.getName(), ErrInvalidType))
		return my
	}

	if getRuleRequired(rules) && my.IsPtr && my.IsNil {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrRequired)}
		return my
	}

	if getRuleNotEmpty(rules) && my.IsZero {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrNotEmpty)}
		return my
	}

	rules.Each(func(_ int, rule string) {
		switch rule {
		case "", "array", "slice":
			if strings.HasPrefix(rule, "min") {
				if min, include = getRuleIntMin(rule); min != nil {
					if include {
						if !(len(value) >= *min) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：>= %d", my.getName(), ErrInvalidLength, *min))
						}
					} else {
						if !(len(value) > *min) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：> %d", my.getName(), ErrInvalidLength, *min))
						}
					}
				}
			}
			if strings.HasPrefix(rule, "max") {
				if max, include = getRuleIntMax(rule); max != nil {
					if include {
						if !(len(value) <= *max) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：<= %d", my.getName(), ErrInvalidLength, *max))
						}
					} else {
						if !(len(value) < *max) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：< %d", my.getName(), ErrInvalidLength, *max))
						}
					}
				}
			}
			if strings.HasPrefix(rule, "size") {
				if size, eq = getRuleIntSize(rule); size != nil {
					if eq {
						if !(len(value) == *size) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：不等于 %d", my.getName(), ErrInvalidLength, *size))
						}
					} else {
						if !(len(value) != *size) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：等于 %d", my.getName(), ErrInvalidLength, *size))
						}
					}
				}
			}
			fallthrough
		case "ex":
			if exFnNames := getRuleExFnNames(rule); len(exFnNames) > 0 {
				for idx2 := range exFnNames {
					if fn := APP.Validator.Ins().GetExFn(exFnNames[idx2]); fn != nil {
						if err := fn(value); err != nil {
							my.wrongs = append(my.wrongs, err)
						}
					}
				}
			}
		}
	})

	return my
}
