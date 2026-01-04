package validatorV3

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/cast"

	"github.com/aid297/aid/array/anyArrayV2"
	"github.com/aid297/aid/ptr"
)

type (
	// FieldInfo 保存了字段的相关信息。
	FieldInfo struct {
		Name      string // 字段名
		Value     any    // 实际值
		Kind      reflect.Kind
		Type      reflect.Type
		IsPtr     bool     // 是否是指针
		IsNil     bool     // 是否为空指针
		Required  bool     // 是否必填
		VRuleTags []string // v-rule tag 的值
		VNameTags []string // v-name tag 的值
		wrongs    []error
	}

	FieldRuleBase struct{ Required string }
	FieldRule     struct{ FieldRuleBase }
)

func (my FieldInfo) Wrongs() []error { return my.wrongs }

func (my FieldInfo) getName() string { return strings.Join(my.VNameTags, ".") }

func (my FieldInfo) getRuleType(rules anyArrayV2.AnyArray[string]) (targetType string) {
	// 获取目标类型
	rules.IfIn(func() { targetType = "string" }, "string")
	rules.IfIn(func() { targetType = "datetime" }, "datetime")
	rules.IfIn(func() { targetType = "date" }, "date")
	rules.IfIn(func() { targetType = "time" }, "time")
	rules.IfIn(func() { targetType = "int" }, "int")
	rules.IfIn(func() { targetType = "int8" }, "int8")
	rules.IfIn(func() { targetType = "int16" }, "int16")
	rules.IfIn(func() { targetType = "int32" }, "int32")
	rules.IfIn(func() { targetType = "int64" }, "int64")
	rules.IfIn(func() { targetType = "uint" }, "uint")
	rules.IfIn(func() { targetType = "uint8" }, "uint8")
	rules.IfIn(func() { targetType = "uint16" }, "uint16")
	rules.IfIn(func() { targetType = "uint32" }, "uint32")
	rules.IfIn(func() { targetType = "uint64" }, "uint64")
	rules.IfIn(func() { targetType = "bool" }, "bool")
	rules.IfIn(func() { targetType = "float32" }, "float32")
	rules.IfIn(func() { targetType = "float64" }, "float64")
	rules.IfIn(func() { targetType = "slice" }, "array", "slice")
	rules.IfIn(func() { targetType = "struct" }, "struct")

	return
}

func getRuleRequired(rules anyArrayV2.AnyArray[string]) bool { return rules.In("required") }

func getRuleExFnNames(rule string) []string {
	var (
		value string
		ok    bool
	)

	if value, ok = strings.CutPrefix(rule, "ex:"); ok {
		return strings.Split(value, ",")
	}

	return nil
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

	return nil, false
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

	return nil, false
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

	return in
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

	return notIn
}

func (my FieldInfo) Check() FieldInfo {
	switch my.Kind {
	case reflect.String:
		return my.checkString()
	case reflect.Int:
		return my.checkInt()
	case reflect.Int8:
		return my.checkInt8()
	case reflect.Int16:
		return my.checkInt16()
	case reflect.Int32:
		return my.checkInt32()
	case reflect.Int64:
		return my.checkInt64()
	case reflect.Uint:
		return my.checkUint()
	case reflect.Uint8:
		return my.checkUint8()
	case reflect.Uint16:
		return my.checkUint16()
	case reflect.Uint32:
		return my.checkUint32()
	case reflect.Uint64:
		return my.checkUint64()
	case reflect.Float32:
		return my.checkFloat32()
	case reflect.Float64:
		return my.checkFloat64()
	case reflect.Bool:
		return my.checkBool()
	case reflect.Array, reflect.Slice:
		return my.checkSlice()
	case reflect.Struct:
		if my.Type == reflect.TypeOf(time.Time{}) {
			return my.checkTime()
		}
		return my
	default:
		return my
	}
}
