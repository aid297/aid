package validatorV3

import (
	"fmt"
	"strings"

	"github.com/aid297/aid/array/anyArrayV2"
)

// checkSlice 检查数组、切片，支持：required、[array|slice|a|s]、min>、min>=、max<、max<=、size:
func (my FieldInfo) checkSlice() FieldInfo {
	var (
		rules          = anyArrayV2.NewList(my.VRuleTags)
		ruleType       = my.getRuleType(rules)
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

	for idx := range my.VRuleTags {
		switch ruleType {
		case "", "array", "slice", "a", "s":
			if strings.HasPrefix(my.VRuleTags[idx], "min") {
				if min, include = getRuleIntMin(my.VRuleTags[idx]); min != nil {
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
			if strings.HasPrefix(my.VRuleTags[idx], "max") {
				if max, include = getRuleIntMax(my.VRuleTags[idx]); max != nil {
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
			if strings.HasPrefix(my.VRuleTags[idx], "size") {
				if size, eq = getRuleIntSize(my.VRuleTags[idx]); size != nil {
					if eq {
						if !(len(value) == *size) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：不等于 %f", my.getName(), ErrInvalidLength, *size))
						}
					} else {
						if !(len(value) != *size) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：等于 %f", my.getName(), ErrInvalidLength, *size))
						}
					}
				}
			}
		}

		if strings.HasPrefix(my.VRuleTags[idx], "ex") {
			if exFnNames := getRuleExFnNames(my.VRuleTags[idx]); len(exFnNames) > 0 {
				for idx2 := range exFnNames {
					if fn := APP.Validator.Ins().GetExFn(exFnNames[idx2]); fn != nil {
						if err := fn(value); err != nil {
							my.wrongs = append(my.wrongs, err)
						}
					}
				}
			}
		}
	}

	return my
}
