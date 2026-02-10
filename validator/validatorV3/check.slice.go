package validatorV3

import (
	"fmt"
	"strings"
)

// checkSlice 检查数组、切片，支持：required、min>、min>=、max<、max<=、size=、size!=、ex:
func (my FieldInfo) checkSlice() FieldInfo {
	var (
		min, max, size *int
		include, eq    bool
		value          []any
		ok             bool
	)

	if value, ok = my.Value.([]any); !ok {
		my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：数组", my.getName(), ErrInvalidType))
		return my
	}

	if getRuleRequired(my.VRuleTags) {
		if my.IsPtr && (my.IsNil || len(value) == 0) {
			my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrRequired)}
			return my
		} else if !my.IsPtr && len(value) == 0 {
			my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrNotEmpty)}
			return my
		}
	}

	my.VRuleTags.Each(func(_ int, rule string) {
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
		} else if strings.HasPrefix(rule, "max") {
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
		} else if strings.HasPrefix(rule, "size") {
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
		} else if rule == "ex" {
			if exFnNames := getRuleExFnNames(rule); len(exFnNames) > 0 {
				for idx2 := range exFnNames {
					if fn := APP.Validator.Once().GetExFn(exFnNames[idx2]); fn != nil {
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
