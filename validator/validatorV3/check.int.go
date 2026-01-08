package validatorV3

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/cast"

	"github.com/aid297/aid/array/anyArrayV2"
)

// checkInt 检查整数，支持：required、not-empty、[int|i]、min>、min>=、max<、max<=、in、not-in、size=、size!=、ex:
func (my FieldInfo) checkInt() FieldInfo {
	var (
		min, max, size *int
		include, eq    bool
		in             []string
		notIn          []string
		value          int
	)

	if my.Kind != reflect.Int {
		my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：整数", my.getName(), ErrInvalidType))
		return my
	}

	if getRuleRequired(my.VRuleTags) && my.IsPtr && my.IsNil {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrRequired)}
		return my
	}

	value, _ = my.Value.(int)

	if getRuleNotEmpty(my.VRuleTags) && my.IsZero {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrNotEmpty)}
		return my
	}

	my.VRuleTags.Each(func(_ int, rule string) {
		switch rule {
		case "", "int", "i":
			if strings.HasPrefix(rule, "min") {
				if min, include = getRuleIntMin(rule); min != nil {
					if include {
						if !(value >= *min) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：>= %d", my.getName(), ErrInvalidLength, *min))
						}
					} else {
						if !(value > *min) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：> %d", my.getName(), ErrInvalidLength, *min))
						}
					}
				}
			}
			if strings.HasPrefix(rule, "max") {
				if max, include = getRuleIntMax(rule); max != nil {
					if include {
						if !(value <= *max) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：<= %d", my.getName(), ErrInvalidLength, *max))
						}
					} else {
						if !(value < *max) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：< %d", my.getName(), ErrInvalidLength, *max))
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
				if size, eq = getRuleIntSize(rule); size != nil {
					if eq {
						if !(cast.ToInt(value) == *size) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：不等于 %d", my.getName(), ErrInvalidLength, *size))
						}
					} else {
						if !(cast.ToInt(value) != *size) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：等于 %d", my.getName(), ErrInvalidLength, *size))
						}
					}
				}
			}
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

// checkInt8 检查整数#8位，支持：required、not-empty、[int8|i8]、min>、min>=、max<、max<=、in、not-in、size=、size!=、ex:
func (my FieldInfo) checkInt8() FieldInfo {
	var (
		min, max, size *int
		include, eq    bool
		in             []string
		notIn          []string
		value          int8
	)

	if my.Kind != reflect.Int8 {
		my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：小数#32位", my.getName(), ErrInvalidType))
		return my
	}

	if getRuleRequired(my.VRuleTags) && my.IsPtr && my.IsNil {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrRequired)}
		return my
	}

	value, _ = my.Value.(int8)

	if getRuleNotEmpty(my.VRuleTags) && my.IsZero {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrNotEmpty)}
		return my
	}

	my.VRuleTags.Each(func(_ int, rule string) {
		switch rule {
		case "", "int8", "i8":
			if strings.HasPrefix(rule, "min") {
				if min, include = getRuleIntMin(rule); min != nil {
					if include {
						if !(cast.ToInt(value) >= *min) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：>= %d", my.getName(), ErrInvalidLength, *min))
						}
					} else {
						if !(cast.ToInt(value) > *min) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：> %d", my.getName(), ErrInvalidLength, *min))
						}
					}
				}
			}
			if strings.HasPrefix(rule, "max") {
				if max, include = getRuleIntMax(rule); max != nil {
					if include {
						if !(cast.ToInt(value) <= *max) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：<= %d", my.getName(), ErrInvalidLength, *max))
						}
					} else {
						if !(cast.ToInt(value) < *max) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：< %d", my.getName(), ErrInvalidLength, *max))
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
				if size, eq = getRuleIntSize(rule); size != nil {
					if eq {
						if !(cast.ToInt(value) == *size) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：不等于 %d", my.getName(), ErrInvalidLength, *size))
						}
					} else {
						if !(cast.ToInt(value) != *size) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：等于 %d", my.getName(), ErrInvalidLength, *size))
						}
					}
				}
			}
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

// checkInt16 检查整数#16位，支持：required、not-empty、[int16|i16]、min>、min>=、max<、max<=、in、not-in、size=、size!=、ex:
func (my FieldInfo) checkInt16() FieldInfo {
	var (
		min, max, size *int
		include, eq    bool
		in             []string
		notIn          []string
		value          int16
	)
	if my.Kind != reflect.Int16 {
		my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：整数#16位", my.getName(), ErrInvalidType))
		return my
	}

	if getRuleRequired(my.VRuleTags) && my.IsPtr && my.IsNil {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrRequired)}
		return my
	}

	value, _ = my.Value.(int16)

	if getRuleNotEmpty(my.VRuleTags) && my.IsZero {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrNotEmpty)}
		return my
	}

	my.VRuleTags.Each(func(_ int, rule string) {
		switch rule {
		case "", "int16", "i16":
			if strings.HasPrefix(rule, "min") {
				if min, include = getRuleIntMin(rule); min != nil {
					if include {
						if !(cast.ToInt(value) >= *min) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：>= %d", my.getName(), ErrInvalidLength, *min))
						}
					} else {
						if !(cast.ToInt(value) > *min) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：> %d", my.getName(), ErrInvalidLength, *min))
						}
					}
				}
			}
			if strings.HasPrefix(rule, "max") {
				if max, include = getRuleIntMax(rule); max != nil {
					if include {
						if !(cast.ToInt(value) <= *max) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：<= %d", my.getName(), ErrInvalidLength, *max))
						}
					} else {
						if !(cast.ToInt(value) < *max) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：< %d", my.getName(), ErrInvalidLength, *max))
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
				if size, eq = getRuleIntSize(rule); size != nil {
					if eq {
						if !(cast.ToInt(value) == *size) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：不等于 %d", my.getName(), ErrInvalidLength, *size))
						}
					} else {
						if !(cast.ToInt(value) != *size) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：等于 %d", my.getName(), ErrInvalidLength, *size))
						}
					}
				}
			}
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

// checkInt32 检查整数#32位，支持：required、not-empty、[int32|i32]、min>、min>=、max<、max<=、in、not-in、size=、size!=、ex:
func (my FieldInfo) checkInt32() FieldInfo {
	var (
		min, max, size *int
		include, eq    bool
		in             []string
		notIn          []string
		value          int32
	)

	if my.Kind != reflect.Int32 {
		my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：整数#32位", my.getName(), ErrInvalidType))
		return my
	}

	if getRuleRequired(my.VRuleTags) && my.IsPtr && my.IsNil {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrRequired)}
		return my
	}

	value, _ = my.Value.(int32)

	if getRuleNotEmpty(my.VRuleTags) && my.IsZero {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrNotEmpty)}
		return my
	}

	my.VRuleTags.Each(func(_ int, rule string) {
		switch rule {
		case "", "int32", "i32":
			if strings.HasPrefix(rule, "min") {
				if min, include = getRuleIntMin(rule); min != nil {
					if include {
						if !(cast.ToInt(value) >= *min) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：>= %d", my.getName(), ErrInvalidLength, *min))
						}
					} else {
						if !(cast.ToInt(value) > *min) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：> %d", my.getName(), ErrInvalidLength, *min))
						}
					}
				}
			}
			if strings.HasPrefix(rule, "max") {
				if max, include = getRuleIntMax(rule); max != nil {
					if include {
						if !(cast.ToInt(value) <= *max) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：<= %d", my.getName(), ErrInvalidLength, *max))
						}
					} else {
						if !(cast.ToInt(value) < *max) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：< %d", my.getName(), ErrInvalidLength, *max))
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
				if size, eq = getRuleIntSize(rule); size != nil {
					if eq {
						if !(cast.ToInt(value) == *size) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：不等于 %d", my.getName(), ErrInvalidLength, *size))
						}
					} else {
						if !(cast.ToInt(value) != *size) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：等于 %d", my.getName(), ErrInvalidLength, *size))
						}
					}
				}
			}
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

// checkInt64 检查整数#64位，支持：required、not-empty、[int64|i64]、min>、min>=、max<、max<=、in、not-in、size=、size!=、ex:
func (my FieldInfo) checkInt64() FieldInfo {
	var (
		min, max, size *int
		include, eq    bool
		in             []string
		notIn          []string
		value          int64
	)

	if my.Kind != reflect.Int64 {
		my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：整数#64位", my.getName(), ErrInvalidType))
		return my
	}

	if getRuleRequired(my.VRuleTags) && my.IsPtr && my.IsNil {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrRequired)}
		return my
	}

	value, _ = my.Value.(int64)

	if getRuleNotEmpty(my.VRuleTags) && my.IsZero {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrNotEmpty)}
		return my
	}

	my.VRuleTags.Each(func(_ int, rule string) {
		switch rule {
		case "", "int64", "i64":
			if strings.HasPrefix(rule, "min") {
				if min, include = getRuleIntMin(rule); min != nil {
					if include {
						if !(cast.ToInt(value) >= *min) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：>= %d", my.getName(), ErrInvalidLength, *min))
						}
					} else {
						if !(cast.ToInt(value) > *min) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：> %d", my.getName(), ErrInvalidLength, *min))
						}
					}
				}
			}
			if strings.HasPrefix(rule, "max") {
				if max, include = getRuleIntMax(rule); max != nil {
					if include {
						if !(cast.ToInt(value) <= *max) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：<= %d", my.getName(), ErrInvalidLength, *max))
						}
					} else {
						if !(cast.ToInt(value) < *max) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：< %d", my.getName(), ErrInvalidLength, *max))
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
				if size, eq = getRuleIntSize(rule); size != nil {
					if eq {
						if !(cast.ToInt(value) == *size) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：不等于 %d", my.getName(), ErrInvalidLength, *size))
						}
					} else {
						if !(cast.ToInt(value) != *size) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：等于 %d", my.getName(), ErrInvalidLength, *size))
						}
					}
				}
			}

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
