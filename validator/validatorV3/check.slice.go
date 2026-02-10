package validatorV3

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/aid297/aid/operation/operationV2"
)

// checkSlice 检查数组、切片，支持：required、min>、min>=、max<、max<=、size=、size!=、ex:
func (my FieldInfo) checkSlice() FieldInfo {
	var (
		min, max, size *int
		include, eq    bool
	)

	if my.Kind != reflect.Slice && my.Kind != reflect.Array {
		my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：切片或数组", my.getName(), ErrInvalidType))
		return my
	}

	if getRuleRequired(my.VRuleTags) {
		if my.IsPtr && (my.IsNil || my.RefValue.Len() == 0) {
			l := my.RefValue.Len()
			println(l)
			my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrRequired)}
			return my
		} else if !my.IsPtr && my.RefValue.Len() == 0 {
			l := my.RefValue.Len()
			println(l)
			my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrNotEmpty)}
			return my
		}
	}

	my.VRuleTags.Each(func(_ int, rule string) {
		if strings.HasPrefix(rule, "min") {
			if min, include = getRuleIntMin(rule); min != nil {
				if include {
					if !(my.RefValue.Len() >= *min) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：>= %d", my.getName(), ErrInvalidLength, *min))
					}
				} else {
					if !(my.RefValue.Len() > *min) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：> %d", my.getName(), ErrInvalidLength, *min))
					}
				}
			}
		} else if strings.HasPrefix(rule, "max") {
			if max, include = getRuleIntMax(rule); max != nil {
				if include {
					if !(my.RefValue.Len() <= *max) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：<= %d", my.getName(), ErrInvalidLength, *max))
					}
				} else {
					if !(my.RefValue.Len() < *max) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：< %d", my.getName(), ErrInvalidLength, *max))
					}
				}
			}
		} else if strings.HasPrefix(rule, "size") {
			if size, eq = getRuleIntSize(rule); size != nil {
				if eq {
					if !(my.RefValue.Len() == *size) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：不等于 %d", my.getName(), ErrInvalidLength, *size))
					}
				} else {
					if !(my.RefValue.Len() != *size) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：等于 %d", my.getName(), ErrInvalidLength, *size))
					}
				}
			}
		} else if strings.HasPrefix(rule, "ex") {
			if exFnNames := getRuleExFnNames(rule); len(exFnNames) > 0 {
				for idx2 := range exFnNames {
					if fn := APP.Validator.Once().GetExFn(exFnNames[idx2]); fn != nil {
						if err := fn(operationV2.NewTernary(operationV2.TrueFn(my.RefValue.Interface)).GetByValue(my.RefValue.CanInterface())); err != nil {
							my.wrongs = append(my.wrongs, err)
						}
					}
				}
			}
		}
	})

	return my
}
