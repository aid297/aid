package validatorV3

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/cast"

	"github.com/aid297/aid/array/anyArrayV2"
)

// checkUint 检查正整数，支持：required、min>、min>=、max<、max<=、in、not-in、size=、size!=、ex:
func (my FieldInfo) checkUint() FieldInfo {
	var (
		min, max, size *uint
		include, eq    bool
		in             []string
		notIn          []string
		value          uint
	)

	if my.Kind != reflect.Uint {
		my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：正整数", my.getName(), ErrInvalidType))
		return my
	}

	if getRuleRequired(my.VRuleTags) && my.IsPtr && my.IsNil {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrRequired)}
		return my
	}

	value, _ = my.Value.(uint)

	my.VRuleTags.Each(func(_ int, rule string) {
		if strings.HasPrefix(rule, "min") {
			if min, include = getRuleUintMin(rule); min != nil {
				if include {
					if !(cast.ToUint(value) >= *min) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：>= %d", my.getName(), ErrInvalidLength, *min))
					}
				} else {
					if !(cast.ToUint(value) > *min) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：> %d", my.getName(), ErrInvalidLength, *min))
					}
				}
			}
		} else if strings.HasPrefix(rule, "max") {
			if max, include = getRuleUintMax(rule); max != nil {
				if include {
					if !(cast.ToUint(value) <= *max) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：<= %d", my.getName(), ErrInvalidLength, *max))
					}
				} else {
					if !(cast.ToUint(value) < *max) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：< %d", my.getName(), ErrInvalidLength, *max))
					}
				}
			}
		} else if strings.HasPrefix(rule, "in") {
			if in = getRuleIn(rule); len(in) > 0 {
				anyArrayV2.NewList(in).IfNotIn(func() {
					my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之中", my.getName(), ErrInvalidValue, in))
				}, cast.ToString(value))
			}
		} else if strings.HasPrefix(rule, "not-in") {
			if notIn = getRuleNotIn(rule); len(notIn) > 0 {
				anyArrayV2.NewList(notIn).IfIn(func() {
					my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之外", my.getName(), ErrInvalidValue, notIn))
				}, cast.ToString(value))
			}
		} else if strings.HasPrefix(rule, "size") {
			if size, eq = getRuleUintSize(rule); size != nil {
				if eq {
					if !(cast.ToUint(value) == *size) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：不等于 %d", my.getName(), ErrInvalidLength, *size))
					}
				} else {
					if !(cast.ToUint(value) != *size) {
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

// checkUint8 检查正整数#8位，支持：required、min>、min>=、max<、max<=、in、not-in、size=、size!=、ex:
func (my FieldInfo) checkUint8() FieldInfo {
	var (
		min, max, size *uint
		include, eq    bool
		in             []string
		notIn          []string
		value          uint8
	)

	if my.Kind != reflect.Uint8 {
		my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：正整数#8位", my.getName(), ErrInvalidType))
		return my
	}

	if getRuleRequired(my.VRuleTags) && my.IsPtr && my.IsNil {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrRequired)}
		return my
	}

	value, _ = my.Value.(uint8)

	my.VRuleTags.Each(func(_ int, rule string) {
		if strings.HasPrefix(rule, "min") {
			if min, include = getRuleUintMin(rule); min != nil {
				if include {
					if !(cast.ToUint(value) >= *min) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：>= %d", my.getName(), ErrInvalidLength, *min))
					}
				} else {
					if !(cast.ToUint(value) > *min) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：> %d", my.getName(), ErrInvalidLength, *min))
					}
				}
			}
		} else if strings.HasPrefix(rule, "max") {
			if max, include = getRuleUintMax(rule); max != nil {
				if include {
					if !(cast.ToUint(value) <= *max) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：<= %d", my.getName(), ErrInvalidLength, *max))
					}
				} else {
					if !(cast.ToUint(value) < *max) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：< %d", my.getName(), ErrInvalidLength, *max))
					}
				}
			}
		} else if strings.HasPrefix(rule, "in") {
			if in = getRuleIn(rule); len(in) > 0 {
				anyArrayV2.NewList(in).IfNotIn(func() {
					my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之中", my.getName(), ErrInvalidValue, in))
				}, cast.ToString(value))
			}
		} else if strings.HasPrefix(rule, "not-in") {
			if notIn = getRuleNotIn(rule); len(notIn) > 0 {
				anyArrayV2.NewList(notIn).IfIn(func() {
					my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之外", my.getName(), ErrInvalidValue, notIn))
				}, cast.ToString(value))
			}
		} else if strings.HasPrefix(rule, "size") {
			if size, eq = getRuleUintSize(rule); size != nil {
				if eq {
					if !(cast.ToUint(value) == *size) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：不等于 %d", my.getName(), ErrInvalidLength, *size))
					}
				} else {
					if !(cast.ToUint(value) != *size) {
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

// checkUint16 检查正整数#16位，支持：required、min>、min>=、max<、max<=、in、not-in、size=、size!=、ex:
func (my FieldInfo) checkUint16() FieldInfo {
	var (
		min, max, size *uint
		include, eq    bool
		in             []string
		notIn          []string
		value          uint16
	)

	if my.Kind != reflect.Uint16 {
		my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：正整数#16位", my.getName(), ErrInvalidType))
		return my
	}

	if getRuleRequired(my.VRuleTags) && my.IsPtr && my.IsNil {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrRequired)}
		return my
	}

	value, _ = my.Value.(uint16)

	my.VRuleTags.Each(func(_ int, rule string) {
		if strings.HasPrefix(rule, "min") {
			if min, include = getRuleUintMin(rule); min != nil {
				if include {
					if !(cast.ToUint(value) >= *min) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：>= %d", my.getName(), ErrInvalidLength, *min))
					}
				} else {
					if !(cast.ToUint(value) > *min) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：> %d", my.getName(), ErrInvalidLength, *min))
					}
				}
			}
		} else if strings.HasPrefix(rule, "max") {
			if max, include = getRuleUintMax(rule); max != nil {
				if include {
					if !(cast.ToUint(value) <= *max) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：<= %d", my.getName(), ErrInvalidLength, *max))
					}
				} else {
					if !(cast.ToUint(value) < *max) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：< %d", my.getName(), ErrInvalidLength, *max))
					}
				}
			}
		} else if strings.HasPrefix(rule, "in") {
			if in = getRuleIn(rule); len(in) > 0 {
				anyArrayV2.NewList(in).IfNotIn(func() {
					my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之中", my.getName(), ErrInvalidValue, in))
				}, cast.ToString(value))
			}
		} else if strings.HasPrefix(rule, "not-in") {
			if notIn = getRuleNotIn(rule); len(notIn) > 0 {
				anyArrayV2.NewList(notIn).IfIn(func() {
					my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之外", my.getName(), ErrInvalidValue, notIn))
				}, cast.ToString(value))
			}
		} else if strings.HasPrefix(rule, "size") {
			if size, eq = getRuleUintSize(rule); size != nil {
				if eq {
					if !(cast.ToUint(value) == *size) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：不等于 %d", my.getName(), ErrInvalidLength, *size))
					}
				} else {
					if !(cast.ToUint(value) != *size) {
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

// checkUint32 检查正整数#32位，支持：required、min>、min>=、max<、max<=、in、not-in、size=、size!=、ex:
func (my FieldInfo) checkUint32() FieldInfo {
	var (
		min, max, size *uint
		include, eq    bool
		in             []string
		notIn          []string
		value          uint32
	)

	if my.Kind != reflect.Uint32 {
		my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：正整数#32位", my.getName(), ErrInvalidType))
		return my
	}

	if getRuleRequired(my.VRuleTags) && my.IsPtr && my.IsNil {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrRequired)}
		return my
	}

	value, _ = my.Value.(uint32)

	my.VRuleTags.Each(func(_ int, rule string) {
		if strings.HasPrefix(rule, "min") {
			if min, include = getRuleUintMin(rule); min != nil {
				if include {
					if !(cast.ToUint(value) >= *min) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：>= %d", my.getName(), ErrInvalidLength, *min))
					}
				} else {
					if !(cast.ToUint(value) > *min) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：> %d", my.getName(), ErrInvalidLength, *min))
					}
				}
			}
		} else if strings.HasPrefix(rule, "max") {
			if max, include = getRuleUintMax(rule); max != nil {
				if include {
					if !(cast.ToUint(value) <= *max) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：<= %d", my.getName(), ErrInvalidLength, *max))
					}
				} else {
					if !(cast.ToUint(value) < *max) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：< %d", my.getName(), ErrInvalidLength, *max))
					}
				}
			}
		} else if strings.HasPrefix(rule, "in") {
			if in = getRuleIn(rule); len(in) > 0 {
				anyArrayV2.NewList(in).IfNotIn(func() {
					my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之中", my.getName(), ErrInvalidValue, in))
				}, cast.ToString(value))
			}
		} else if strings.HasPrefix(rule, "not-in") {
			if notIn = getRuleNotIn(rule); len(notIn) > 0 {
				anyArrayV2.NewList(notIn).IfIn(func() {
					my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之外", my.getName(), ErrInvalidValue, notIn))
				}, cast.ToString(value))
			}
		} else if strings.HasPrefix(rule, "size") {
			if size, eq = getRuleUintSize(rule); size != nil {
				if eq {
					if !(cast.ToUint(value) == *size) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：不等于 %d", my.getName(), ErrInvalidLength, *size))
					}
				} else {
					if !(cast.ToUint(value) != *size) {
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

// checkUint64 检查正整数#64位，支持：required、not-empty、[uint64|u64]、min>、min>=、max<、max<=、in、not-in、size=、size!=、ex:
func (my FieldInfo) checkUint64() FieldInfo {
	var (
		min, max, size *uint
		include, eq    bool
		in             []string
		notIn          []string
		value          uint64
	)

	if my.Kind != reflect.Uint64 {
		my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：正整数#64位", my.getName(), ErrInvalidType))
		return my
	}

	if getRuleRequired(my.VRuleTags) && my.IsPtr && my.IsNil {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrRequired)}
		return my
	}

	value, _ = my.Value.(uint64)

	my.VRuleTags.Each(func(_ int, rule string) {
		if strings.HasPrefix(rule, "min") {
			if min, include = getRuleUintMin(rule); min != nil {
				if include {
					if !(cast.ToUint(value) >= *min) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：>= %d", my.getName(), ErrInvalidLength, *min))
					}
				} else {
					if !(cast.ToUint(value) > *min) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：> %d", my.getName(), ErrInvalidLength, *min))
					}
				}
			}
		} else if strings.HasPrefix(rule, "max") {
			if max, include = getRuleUintMax(rule); max != nil {
				if include {
					if !(cast.ToUint(value) <= *max) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：<= %d", my.getName(), ErrInvalidLength, *max))
					}
				} else {
					if !(cast.ToUint(value) < *max) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：< %d", my.getName(), ErrInvalidLength, *max))
					}
				}
			}
		} else if strings.HasPrefix(rule, "in") {
			if in = getRuleIn(rule); len(in) > 0 {
				anyArrayV2.NewList(in).IfNotIn(func() {
					my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之中", my.getName(), ErrInvalidValue, in))
				}, cast.ToString(value))
			}
		} else if strings.HasPrefix(rule, "not-in") {
			if notIn = getRuleNotIn(rule); len(notIn) > 0 {
				anyArrayV2.NewList(notIn).IfIn(func() {
					my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之外", my.getName(), ErrInvalidValue, notIn))
				}, cast.ToString(value))
			}
		} else if strings.HasPrefix(rule, "size") {
			if size, eq = getRuleUintSize(rule); size != nil {
				if eq {
					if !(cast.ToUint(value) == *size) {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：不等于 %d", my.getName(), ErrInvalidLength, *size))
					}
				} else {
					if !(cast.ToUint(value) != *size) {
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
