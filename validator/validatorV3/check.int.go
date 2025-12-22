package validatorV3

import (
	"fmt"
	"strings"

	"github.com/spf13/cast"

	"github.com/aid297/aid/array/anyArrayV2"
)

// checkInt 检查整数，支持：required、[int|i]、min>、min>=、max<、max<=、in、not-in、size:
func (my FieldInfo) checkInt() FieldInfo {
	var (
		rules          = anyArrayV2.NewList(my.VRuleTags)
		ruleType       = my.getRuleType(rules)
		min, max, size *int
		include        bool
		in             []string
		notIn          []string
		value          int
		ok             bool
	)

	if value, ok = my.Value.(int); !ok {
		my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：整数", my.getName(), ErrInvalidType))
		return my
	}

	for idx := range my.VRuleTags {
		if my.VRuleTags[idx] == "required" && my.IsPtr && my.IsNil {
			my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w", my.getName(), ErrRequired))
			return my
		}

		switch ruleType {
		case "", "int", "i":
			if strings.HasPrefix(my.VRuleTags[idx], "min") {
				if min, include = getRuleIntMin(my.VRuleTags[idx]); min != nil {
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
			if strings.HasPrefix(my.VRuleTags[idx], "max") {
				if max, include = getRuleIntMax(my.VRuleTags[idx]); max != nil {
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
				if size = getRuleIntSize(my.VRuleTags[idx]); size != nil {
					if value != *size {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：= %d", my.getName(), ErrInvalidLength, *size))
					}
				}
			}
		}

		if strings.HasPrefix(my.VRuleTags[idx], "ex") {
			if exFnNames := getRuleExFnNames(my.VRuleTags[idx]); len(exFnNames) > 0 {
				for idx2 := range exFnNames {
					if fn := APP.Validator.Ins().Get(exFnNames[idx2]); fn != nil {
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

// checkInt8 检查整数#8位，支持：required、[int8|i8]、min>、min>=、max<、max<=、in、not-in、size:
func (my FieldInfo) checkInt8() FieldInfo {
	var (
		rules          = anyArrayV2.NewList(my.VRuleTags)
		ruleType       = my.getRuleType(rules)
		min, max, size *int
		include        bool
		in             []string
		notIn          []string
		value          int8
		ok             bool
	)

	if value, ok = my.Value.(int8); !ok {
		my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：小数#32位", my.getName(), ErrInvalidType))
		return my
	}

	if getRuleRequired(rules) && my.IsPtr && my.IsNil {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrRequired)}
		return my
	}

	for idx := range my.VRuleTags {
		switch ruleType {
		case "", "int8", "i8":
			if strings.HasPrefix(my.VRuleTags[idx], "min") {
				if min, include = getRuleIntMin(my.VRuleTags[idx]); min != nil {
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
			if strings.HasPrefix(my.VRuleTags[idx], "max") {
				if max, include = getRuleIntMax(my.VRuleTags[idx]); max != nil {
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
				if size = getRuleIntSize(my.VRuleTags[idx]); size != nil {
					if cast.ToInt(value) != *size {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：= %d", my.getName(), ErrInvalidLength, *size))
					}
				}
			}
		}

		if strings.HasPrefix(my.VRuleTags[idx], "ex") {
			if exFnNames := getRuleExFnNames(my.VRuleTags[idx]); len(exFnNames) > 0 {
				for idx2 := range exFnNames {
					if fn := APP.Validator.Ins().Get(exFnNames[idx2]); fn != nil {
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

// checkInt16 检查整数#16位，支持：required、[int16|i16]、min>、min>=、max<、max<=、in、not-in、size:
func (my FieldInfo) checkInt16() FieldInfo {
	var (
		rules          = anyArrayV2.NewList(my.VRuleTags)
		ruleType       = my.getRuleType(rules)
		min, max, size *int
		include        bool
		in             []string
		notIn          []string
		value          int16
		ok             bool
	)
	if value, ok = my.Value.(int16); !ok {
		my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：整数#16位", my.getName(), ErrInvalidType))
		return my
	}

	if getRuleRequired(rules) && my.IsPtr && my.IsNil {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrRequired)}
		return my
	}

	for idx := range my.VRuleTags {
		switch ruleType {
		case "", "int16", "i16":
			if strings.HasPrefix(my.VRuleTags[idx], "min") {
				if min, include = getRuleIntMin(my.VRuleTags[idx]); min != nil {
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
			if strings.HasPrefix(my.VRuleTags[idx], "max") {
				if max, include = getRuleIntMax(my.VRuleTags[idx]); max != nil {
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
				if size = getRuleIntSize(my.VRuleTags[idx]); size != nil {
					if cast.ToInt(value) != *size {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：= %d", my.getName(), ErrInvalidLength, *size))
					}
				}
			}
		}

		if strings.HasPrefix(my.VRuleTags[idx], "ex") {
			if exFnNames := getRuleExFnNames(my.VRuleTags[idx]); len(exFnNames) > 0 {
				for idx2 := range exFnNames {
					if fn := APP.Validator.Ins().Get(exFnNames[idx2]); fn != nil {
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

// checkInt32 检查整数#32位，支持：required、[int32|i32]、min>、min>=、max<、max<=、in、not-in、size:
func (my FieldInfo) checkInt32() FieldInfo {
	var (
		rules          = anyArrayV2.NewList(my.VRuleTags)
		ruleType       = my.getRuleType(rules)
		min, max, size *int
		include        bool
		in             []string
		notIn          []string
		value          int32
		ok             bool
	)

	if value, ok = my.Value.(int32); !ok {
		my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：整数#32位", my.getName(), ErrInvalidType))
		return my
	}

	if getRuleRequired(rules) && my.IsPtr && my.IsNil {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrRequired)}
		return my
	}

	for idx := range my.VRuleTags {
		switch ruleType {
		case "", "int32", "i32":
			if strings.HasPrefix(my.VRuleTags[idx], "min") {
				if min, include = getRuleIntMin(my.VRuleTags[idx]); min != nil {
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
			if strings.HasPrefix(my.VRuleTags[idx], "max") {
				if max, include = getRuleIntMax(my.VRuleTags[idx]); max != nil {
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
				if size = getRuleIntSize(my.VRuleTags[idx]); size != nil {
					if cast.ToInt(value) != *size {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：= %d", my.getName(), ErrInvalidLength, *size))
					}
				}
			}
		}

		if strings.HasPrefix(my.VRuleTags[idx], "ex") {
			if exFnNames := getRuleExFnNames(my.VRuleTags[idx]); len(exFnNames) > 0 {
				for idx2 := range exFnNames {
					if fn := APP.Validator.Ins().Get(exFnNames[idx2]); fn != nil {
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

// checkInt64 检查整数#64位，支持：required、[int64|i64]、min>、min>=、max<、max<=、in、not-in、size:
func (my FieldInfo) checkInt64() FieldInfo {
	var (
		rules          = anyArrayV2.NewList(my.VRuleTags)
		ruleType       = my.getRuleType(rules)
		min, max, size *int
		include        bool
		in             []string
		notIn          []string
		value          int64
		ok             bool
	)

	if value, ok = my.Value.(int64); !ok {
		my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：整数#64位", my.getName(), ErrInvalidType))
		return my
	}

	if getRuleRequired(rules) && my.IsPtr && my.IsNil {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrRequired)}
		return my
	}

	for idx := range my.VRuleTags {
		switch ruleType {
		case "", "int64", "i64":
			if strings.HasPrefix(my.VRuleTags[idx], "min") {
				if min, include = getRuleIntMin(my.VRuleTags[idx]); min != nil {
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
			if strings.HasPrefix(my.VRuleTags[idx], "max") {
				if max, include = getRuleIntMax(my.VRuleTags[idx]); max != nil {
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
				if size = getRuleIntSize(my.VRuleTags[idx]); size != nil {
					if cast.ToInt(value) != *size {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：= %d", my.getName(), ErrInvalidLength, *size))
					}
				}
			}
		}

		if strings.HasPrefix(my.VRuleTags[idx], "ex") {
			if exFnNames := getRuleExFnNames(my.VRuleTags[idx]); len(exFnNames) > 0 {
				for idx2 := range exFnNames {
					if fn := APP.Validator.Ins().Get(exFnNames[idx2]); fn != nil {
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
