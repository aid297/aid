package validatorV3

import (
	"strings"
	"time"

	"github.com/aid297/aid/array/anyArrayV2"
	"github.com/aid297/aid/ptr"
	"github.com/spf13/cast"
)

// func (my FieldInfo) getRuleType(rules anyArrayV2.AnyArray[string]) (targetType string) {
// 	// 获取目标类型
// 	rules.IfIn(func() { targetType = "string" }, "string")
// 	rules.IfIn(func() { targetType = "datetime" }, "datetime")
// 	rules.IfIn(func() { targetType = "date" }, "date")
// 	rules.IfIn(func() { targetType = "time" }, "time")
// 	rules.IfIn(func() { targetType = "int" }, "int")
// 	rules.IfIn(func() { targetType = "int8" }, "int8")
// 	rules.IfIn(func() { targetType = "int16" }, "int16")
// 	rules.IfIn(func() { targetType = "int32" }, "int32")
// 	rules.IfIn(func() { targetType = "int64" }, "int64")
// 	rules.IfIn(func() { targetType = "uint" }, "uint")
// 	rules.IfIn(func() { targetType = "uint8" }, "uint8")
// 	rules.IfIn(func() { targetType = "uint16" }, "uint16")
// 	rules.IfIn(func() { targetType = "uint32" }, "uint32")
// 	rules.IfIn(func() { targetType = "uint64" }, "uint64")
// 	rules.IfIn(func() { targetType = "bool" }, "bool")
// 	rules.IfIn(func() { targetType = "float32" }, "float32")
// 	rules.IfIn(func() { targetType = "float64" }, "float64")
// 	rules.IfIn(func() { targetType = "slice" }, "array", "slice")
// 	rules.IfIn(func() { targetType = "struct" }, "struct")
// 	rules.IfIn(func() { targetType = "ex" }, "ex")

// 	return
// }

func getRuleRequired(rules anyArrayV2.AnyArray[string]) bool { return rules.In("required") }

func getRuleNotEmpty(rules anyArrayV2.AnyArray[string]) bool { return rules.In("not-empty") }

func getRuleExFnNames(rule string) (exFnNames []string) {
	var (
		value string
		ok    bool
	)

	if value, ok = strings.CutPrefix(rule, "ex:"); ok {
		exFnNames = strings.Split(value, ",")
		return
	}

	return
}

func getRuleUintSize(rule string) (size *uint, eq bool) {
	var s *int
	s, eq = getRuleIntSize(rule)
	size = ptr.New(uint(*s))
	return
}

func getRuleIntSize(rule string) (size *int, eq bool) {
	var (
		value string
		ok    bool
	)

	if value, ok = strings.CutPrefix(rule, "size="); ok {
		size = ptr.New(cast.ToInt(value))
		eq = true
		return
	}

	if value, ok = strings.CutPrefix(rule, "size="); ok {
		size = ptr.New(cast.ToInt(value))
		eq = false
		return
	}

	return
}

func getRuleFloatSize(rule string) (size *float64, eq bool) {
	var (
		value string
		ok    bool
	)

	if value, ok = strings.CutPrefix(rule, "size="); ok {
		size = ptr.New(cast.ToFloat64(value))
		eq = true
		return
	}

	if value, ok = strings.CutPrefix(rule, "size!="); ok {
		size = ptr.New(cast.ToFloat64(value))
		eq = false
		return
	}

	return
}

func getRuleUintMin(rule string) (size *uint, include bool) {
	var s *int
	s, include = getRuleIntMin(rule)
	size = ptr.New(uint(*s))
	return
}

func getRuleUintMax(rule string) (size *uint, include bool) {
	var s *int
	s, include = getRuleIntMax(rule)
	size = ptr.New(uint(*s))
	return
}

func getRuleIntMin(rule string) (size *int, include bool) {
	var (
		value string
		ok    bool
	)

	if value, ok = strings.CutPrefix(rule, "min>="); ok {
		size = ptr.New(cast.ToInt(value))
		include = true
		return
	}
	if value, ok = strings.CutPrefix(rule, "min>"); ok {
		size = ptr.New(cast.ToInt(value))
		return
	}

	return
}

func getRuleIntMax(rule string) (size *int, include bool) {
	var (
		value string
		ok    bool
	)

	if value, ok = strings.CutPrefix(rule, "max<="); ok {
		size = ptr.New(cast.ToInt(value))
		include = true
		return
	}
	if value, ok = strings.CutPrefix(rule, "max<"); ok {
		size = ptr.New(cast.ToInt(value))
		include = false
		return
	}

	return
}

func getRuleFloatMin(rule string) (size *float64, include bool) {
	var (
		value string
		ok    bool
	)

	if value, ok = strings.CutPrefix(rule, "min>="); ok {
		size = ptr.New(cast.ToFloat64(value))
		include = true
		return
	}
	if value, ok = strings.CutPrefix(rule, "min>"); ok {
		size = ptr.New(cast.ToFloat64(value))
		include = false
		return
	}

	return
}

func getRuleFloatMax(rule string) (size *float64, include bool) {
	var (
		value string
		ok    bool
	)

	if value, ok = strings.CutPrefix(rule, "max<="); ok {
		size = ptr.New(cast.ToFloat64(value))
		include = true
		return
	}
	if value, ok = strings.CutPrefix(rule, "max<"); ok {
		size = ptr.New(cast.ToFloat64(value))
		include = false
		return
	}

	return
}

func getRuleIn(rule string) (in []string) {
	var (
		value string
		ok    bool
	)
	if value, ok = strings.CutPrefix(rule, "in:"); ok {
		in = strings.Split(value, ",")
		return
	}

	return
}

func getRuleNotIn(rule string) (notIn []string) {
	var (
		value string
		ok    bool
	)
	if value, ok = strings.CutPrefix(rule, "not-in:"); ok {
		notIn = strings.Split(value, ",")
		return
	}

	return
}

func getRuleTimeMin(rule string) (t *time.Time, include bool) {
	var (
		value string
		ok    bool
	)

	if value, ok = strings.CutPrefix(rule, "min>="); ok {
		t = ptr.New(cast.ToTime(value))
		include = true
		return
	}
	if value, ok = strings.CutPrefix(rule, "min>"); ok {
		t = ptr.New(cast.ToTime(value))
		include = false
		return
	}

	return
}

func getRuleTimeMax(rule string) (t *time.Time, include bool) {
	var (
		value string
		ok    bool
	)

	if value, ok = strings.CutPrefix(rule, "max<="); ok {
		t = ptr.New(cast.ToTime(value))
		include = true
		return
	}
	if value, ok = strings.CutPrefix(rule, "max<"); ok {
		t = ptr.New(cast.ToTime(value))
		include = false
		return
	}

	return
}

func getRuleTimeIn(rule string) (in []time.Time) {
	var (
		value string
		ok    bool
		times []string
	)
	if value, ok = strings.CutPrefix(rule, "in:"); ok {
		times = strings.Split(value, ",")
	}

	if len(times) > 0 {
		in = make([]time.Time, 0, len(times))
		for idx := range times {
			in = append(in, cast.ToTime(times[idx]))
		}
	}

	return
}

func getRuleTimeNotIn(rule string) (notIn []time.Time) {
	var (
		value string
		ok    bool
		times []string
	)
	if value, ok = strings.CutPrefix(rule, "not-in:"); ok {
		times = strings.Split(value, ",")
	}

	if len(times) > 0 {
		notIn = make([]time.Time, 0, len(times))
		for idx := range times {
			notIn = append(notIn, cast.ToTime(times[idx]))
		}
	}

	return
}
