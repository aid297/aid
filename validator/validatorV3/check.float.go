package validatorV3

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/cast"

	"github.com/aid297/aid/array/anyArrayV2"
)

// checkFloat32 检查小数#32位，支持：required、not-empty、[float32|f32]、min>、min>=、max<、max<=、in、not-in、size=、size!=、ex:
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

	if getRuleNotEmpty(my.VRuleTags) && my.IsZero {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrNotEmpty)}
		return my
	}

	my.VRuleTags.Each(func(_ int, rule string) {
		switch rule {
		case "", "float32", "f32":
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
			}
			if strings.HasPrefix(rule, "max") {
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
			}
			if strings.HasPrefix(rule, "in") {
				if in = getRuleIn(rule); len(in) > 0 {
					anyArrayV2.NewList(in).IfNotIn(func() {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之中", my.getName(), ErrInvalidValue, in))
					}, cast.ToString(value))
				}
			}
			if strings.HasPrefix(rule, "not-in") {
				if notIn = getRuleNotIn(rule); len(notIn) > 0 {
					anyArrayV2.NewList(notIn).IfIn(func() {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之外", my.getName(), ErrInvalidValue, notIn))
					}, cast.ToString(value))
				}
			}
			if strings.HasPrefix(rule, "size") {
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
			}

		case "ex":
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

// checkFloat64 检查小数#64位，支持：required、not-empty、[float|f64]、min>、min>=、max<、max<=、in、not-in、size=、size!=、ex:
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

	if getRuleNotEmpty(my.VRuleTags) && my.IsZero {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrNotEmpty)}
		return my
	}

	my.VRuleTags.Each(func(_ int, rule string) {
		switch rule {
		case "", "float64", "f64":
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
			}
			if strings.HasPrefix(rule, "max") {
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
			}
			if strings.HasPrefix(rule, "in") {
				if in = getRuleIn(rule); len(in) > 0 {
					anyArrayV2.NewList(in).IfNotIn(func() {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之中", my.getName(), ErrInvalidValue, in))
					}, cast.ToString(value))
				}
			}
			if strings.HasPrefix(rule, "not-in") {
				if notIn = getRuleNotIn(rule); len(notIn) > 0 {
					anyArrayV2.NewList(notIn).IfIn(func() {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之外", my.getName(), ErrInvalidValue, notIn))
					}, cast.ToString(value))
				}
			}
			if strings.HasPrefix(rule, "size") {
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
			}
		case "ex":
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
