package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/aid297/aid/v2/anySlice"
	"github.com/spf13/cast"
)

// checkFloat32 检查小数#32位，支持：required、min>、min>=、max<、max<=、in、not-in、size=、size!=、ex:
func (my FieldInfo) checkFloat32() FieldInfo {
	var (
		min, max, size *float64
		include, eq    bool
		in             []string
		notIn          []string
		value          float32
	)

	if my.Kind != reflect.Float32 {
		my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：小数#32位", my.getName(), ErrInvalidType))
		return my
	}

	if getRuleRequired(my.VRuleTags) && my.IsPtr && my.IsNil {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrRequired)}
		return my
	}

	value, _ = my.Value.(float32)

	my.VRuleTags.Each(func(_ int, rule string) (isBreak bool) {
		if strings.HasPrefix(rule, "min") {
			if min, include = getRuleFloatMin(rule); min != nil {
				if include {
					if !(cast.ToFloat64(value) >= *min) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：>= %f", my.getName(), ErrInvalidLength, *min))
					}
				} else {
					if !(cast.ToFloat64(value) > *min) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：> %f", my.getName(), ErrInvalidLength, *min))
					}
				}
			}
		} else if strings.HasPrefix(rule, "max") {
			if max, include = getRuleFloatMax(rule); max != nil {
				if include {
					if !(cast.ToFloat64(value) <= *max) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：<= %f", my.getName(), ErrInvalidLength, *max))
					}
				} else {
					if !(cast.ToFloat64(value) < *max) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：< %f", my.getName(), ErrInvalidLength, *max))
					}
				}
			}
		} else if strings.HasPrefix(rule, "in") {
			if in = getRuleIn(rule); len(in) > 0 {
				anySlice.New(anySlice.List(in)).IfNotIn(func(_ anySlice.AnySlicer[string]) {
					my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之中", my.getName(), ErrInvalidValue, in))
				}, cast.ToString(value))
			}
		} else if strings.HasPrefix(rule, "not-in") {
			if notIn = getRuleNotIn(rule); len(notIn) > 0 {
				anySlice.New(anySlice.List(notIn)).IfIn(func(_ anySlice.AnySlicer[string]) {
					my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之外", my.getName(), ErrInvalidValue, notIn))
				}, cast.ToString(value))
			}
		} else if strings.HasPrefix(rule, "size") {
			if size, eq = getRuleFloatSize(rule); size != nil {
				if eq {
					if !(cast.ToFloat64(value) == *size) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：不等于 %f", my.getName(), ErrInvalidLength, *size))
					}
				} else {
					if !(cast.ToFloat64(value) != *size) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：等于 %f", my.getName(), ErrInvalidLength, *size))
					}
				}
			}
		} else if strings.HasPrefix(rule, "ex") {
			if exFnNames := getRuleExFnNames(rule); len(exFnNames) > 0 {
				for idx2 := range exFnNames {
					if fn := OnceValidator().GetExFn(exFnNames[idx2]); fn != nil {
						if err := fn(value); err != nil {
							my.wrongs = append(my.wrongs, err)
						}
					}
				}
			}
		}

		return
	})

	return my
}

// checkFloat64 检查小数#64位，支持：required、min>、min>=、max<、max<=、in、not-in、size=、size!=、ex:
func (my FieldInfo) checkFloat64() FieldInfo {
	var (
		min, max, size *float64
		include, eq    bool
		in             []string
		notIn          []string
		value          float64
	)

	if my.Kind != reflect.Float64 {
		my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：小数#64位", my.getName(), ErrInvalidType))
		return my
	}

	if getRuleRequired(my.VRuleTags) && my.IsPtr && my.IsNil {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrRequired)}
		return my
	}

	value, _ = my.Value.(float64)

	my.VRuleTags.Each(func(_ int, rule string) (isBreak bool) {
		if strings.HasPrefix(rule, "min") {
			if min, include = getRuleFloatMin(rule); min != nil {
				if include {
					if !(cast.ToFloat64(value) >= *min) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：>= %f", my.getName(), ErrInvalidLength, *min))
					}
				} else {
					if !(cast.ToFloat64(value) > *min) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：> %f", my.getName(), ErrInvalidLength, *min))
					}
				}
			}
		} else if strings.HasPrefix(rule, "max") {
			if max, include = getRuleFloatMax(rule); max != nil {
				if include {
					if !(cast.ToFloat64(value) <= *max) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：<= %f", my.getName(), ErrInvalidLength, *max))
					}
				} else {
					if !(cast.ToFloat64(value) < *max) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：< %f", my.getName(), ErrInvalidLength, *max))
					}
				}
			}
		} else if strings.HasPrefix(rule, "in") {
			if in = getRuleIn(rule); len(in) > 0 {
				anySlice.New(anySlice.List(in)).IfNotIn(func(_ anySlice.AnySlicer[string]) {
					my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之中", my.getName(), ErrInvalidValue, in))
				}, cast.ToString(value))
			}
		} else if strings.HasPrefix(rule, "not-in") {
			if notIn = getRuleNotIn(rule); len(notIn) > 0 {
				anySlice.NewList(notIn).IfIn(func(_ anySlice.AnySlicer[string]) {
					my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之外", my.getName(), ErrInvalidValue, notIn))
				}, cast.ToString(value))
			}
		} else if strings.HasPrefix(rule, "size") {
			if size, eq = getRuleFloatSize(rule); size != nil {
				if eq {
					if !(cast.ToFloat64(value) == *size) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：不等于 %f", my.getName(), ErrInvalidLength, *size))
					}
				} else {
					if !(cast.ToFloat64(value) != *size) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：等于 %f", my.getName(), ErrInvalidLength, *size))
					}
				}
			}
		} else if strings.HasPrefix(rule, "ex") {
			if exFnNames := getRuleExFnNames(rule); len(exFnNames) > 0 {
				for idx2 := range exFnNames {
					if fn := OnceValidator().GetExFn(exFnNames[idx2]); fn != nil {
						if err := fn(value); err != nil {
							my.wrongs = append(my.wrongs, err)
						}
					}
				}
			}
		}

		return
	})

	return my
}
