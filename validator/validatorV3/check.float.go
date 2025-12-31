package validatorV3

import (
	"fmt"
	"strings"

	"github.com/spf13/cast"

	"github.com/aid297/aid/array/anyArrayV2"
)

// checkFloat32 检查小数#32位，支持：required、[float32|f32]、min>、min>=、max<、max<=、in、not-in、size=、size!=
func (my FieldInfo) checkFloat32() FieldInfo {
	var (
		rules          = anyArrayV2.NewList(my.VRuleTags)
		ruleType       = my.getRuleType(rules)
		min, max, size *float64
		include, eq    bool
		in             []string
		notIn          []string
		value          float32
		ok             bool
	)

	if value, ok = my.Value.(float32); !ok {
		my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：小数#32位", my.getName(), ErrInvalidType))
		return my
	}

	if getRuleRequired(rules) && my.IsPtr && my.IsNil {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrRequired)}
		return my
	}

	for idx := range my.VRuleTags {
		if my.VRuleTags[idx] == "required" && my.IsPtr && my.IsNil {
			my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w", my.getName(), ErrRequired))
		}

		switch ruleType {
		case "", "float32", "f32":
			if strings.HasPrefix(my.VRuleTags[idx], "min") {
				if min, include = getRuleFloatMin(my.VRuleTags[idx]); min != nil {
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
			}
			if strings.HasPrefix(my.VRuleTags[idx], "max") {
				if max, include = getRuleFloatMax(my.VRuleTags[idx]); max != nil {
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
			}
			if strings.HasPrefix(my.VRuleTags[idx], "in") {
				if in = getRuleIn(my.VRuleTags[idx]); len(in) > 0 {
					anyArrayV2.NewList(in).IfNotIn(func() {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之中", my.getName(), ErrInvalidValue, in))
					}, cast.ToString(value))
				}
			}
			if strings.HasPrefix(my.VRuleTags[idx], "not-in") {
				if notIn = getRuleNotIn(my.VRuleTags[idx]); len(notIn) > 0 {
					anyArrayV2.NewList(notIn).IfIn(func() {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之外", my.getName(), ErrInvalidValue, notIn))
					}, cast.ToString(value))
				}
			}
			if strings.HasPrefix(my.VRuleTags[idx], "size") {
				if size, eq = getRuleFloatSize(my.VRuleTags[idx]); size != nil {
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

// checkFloat64 检查小数#64位，支持：required、[float|f64]、min>、min>=、max<、max<=、in、not-in、size:
func (my FieldInfo) checkFloat64() FieldInfo {
	var (
		rules          = anyArrayV2.NewList(my.VRuleTags)
		ruleType       = my.getRuleType(rules)
		min, max, size *float64
		include, eq    bool
		in             []string
		notIn          []string
		value          float64
		ok             bool
	)

	if value, ok = my.Value.(float64); !ok {
		my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：小数#64位", my.getName(), ErrInvalidType))
		return my
	}

	if getRuleRequired(rules) && my.IsPtr && my.IsNil {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrRequired)}
		return my
	}

	for idx := range my.VRuleTags {
		switch ruleType {
		case "", "float64", "f64":
			if strings.HasPrefix(my.VRuleTags[idx], "min") {
				if min, include = getRuleFloatMin(my.VRuleTags[idx]); min != nil {
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
			}
			if strings.HasPrefix(my.VRuleTags[idx], "max") {
				if max, include = getRuleFloatMax(my.VRuleTags[idx]); max != nil {
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
			}
			if strings.HasPrefix(my.VRuleTags[idx], "in") {
				if in = getRuleIn(my.VRuleTags[idx]); len(in) > 0 {
					anyArrayV2.NewList(in).IfNotIn(func() {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之中", my.getName(), ErrInvalidValue, in))
					}, cast.ToString(value))
				}
			}
			if strings.HasPrefix(my.VRuleTags[idx], "not-in") {
				if notIn = getRuleNotIn(my.VRuleTags[idx]); len(notIn) > 0 {
					anyArrayV2.NewList(notIn).IfIn(func() {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之外", my.getName(), ErrInvalidValue, notIn))
					}, cast.ToString(value))
				}
			}
			if strings.HasPrefix(my.VRuleTags[idx], "size") {
				if size, eq = getRuleFloatSize(my.VRuleTags[idx]); size != nil {
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
