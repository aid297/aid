package validator

import (
	"reflect"
	"strings"

	"github.com/spf13/cast"

	"github.com/aid297/aid/array"
	"github.com/aid297/aid/operation"

	"github.com/aid297/aid/common"
)

// checkFloat32 验证：float32 -> 支持的规则 required、size<、size<=、size>、size>=、range=
func (my *ValidatorApp[T]) checkFloat32(rule, fieldName string, value any) error {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		if reflect.ValueOf(value).IsNil() {
			return operation.Ternary(rule == "required", RequiredErr.New(fieldName), nil)
		}
		value = reflect.ValueOf(value).Elem().Interface()
	}

	switch {
	case strings.HasPrefix(rule, "size<="):
		small := strings.TrimPrefix(rule, "size<=")
		if value.(float32) <= common.ToFloat32(small) {
			return LengthErr.NewFormat("[%s]长度不能小于等于[%f]", fieldName, common.ToFloat32(small))
		}
	case strings.HasPrefix(rule, "size<"):
		small := strings.TrimPrefix(rule, "size<")
		if value.(float32) < common.ToFloat32(small) {
			return LengthErr.NewFormat("[%s]长度不能小于[%f]", fieldName, common.ToFloat32(small))
		}
	case strings.HasPrefix(rule, "size>="):
		large := strings.TrimPrefix(rule, "size>=")
		if value.(float32) >= common.ToFloat32(large) {
			return LengthErr.NewFormat("[%s]长度不能大于等于[%f]", fieldName, common.ToFloat32(large))
		}
	case strings.HasPrefix(rule, "size>"):
		large := strings.TrimPrefix(rule, "size>")
		if value.(float32) > common.ToFloat32(large) {
			return LengthErr.NewFormat("[%s]长度不能大于[%f]", fieldName, common.ToFloat32(large))
		}
	case strings.HasPrefix(rule, "size="):
		size := strings.TrimPrefix(rule, "size=")
		if value.(float32) != common.ToFloat32(size) {
			return LengthErr.NewFormat("[%s]长度必须等于[%f]", fieldName, common.ToFloat32(size))
		}
	case strings.HasPrefix(rule, "size!="):
		size := strings.TrimPrefix(rule, "size!=")
		if value.(float32) == common.ToFloat32(size) {
			return LengthErr.NewFormat("[%s]长度必须不等于[%f]", fieldName, common.ToFloat32(size))
		}
	case strings.HasPrefix(rule, "range="):
		between := strings.TrimPrefix(rule, "range=")
		betweenRange := strings.Split(between, ",")
		if len(betweenRange) != 2 {
			return RuleErr.NewFormat("[%s]规则定义错误#1，规则格式：f,f", fieldName)
		}
		small := common.ToFloat32(betweenRange[0])
		large := common.ToFloat32(betweenRange[1])
		if value.(float32) < small || value.(float32) > large {
			return LengthErr.NewFormat("[%s]长度必须在：%f~%f之间", fieldName, small, large)
		}
	case strings.HasPrefix(rule, "in="):
		inValuesStr := strings.TrimPrefix(rule, "in=")
		inValuesArr := array.Cast(array.New(strings.Split(inValuesStr, ",")), func(value string) float32 { return cast.ToFloat32(value) })
		if !inValuesArr.In(value.(float32)) {
			return ValidateErr.NewFormat("[%s]值必须在[%s]中", fieldName, inValuesStr)
		}
	case strings.HasPrefix(rule, "not-in="):
		inValuesStr := strings.TrimPrefix(rule, "not-in=")
		inValuesArr := array.Cast(array.New(strings.Split(inValuesStr, ",")), func(value string) float32 { return cast.ToFloat32(value) })
		if inValuesArr.In(value.(float32)) {
			return ValidateErr.NewFormat("[%s]值不可为以下内容：[%s]", fieldName, inValuesStr)
		}
	}

	return nil
}

// checkFloat64 验证：float64 -> 支持的规则 required、size<、size<=、size>、size>=、range=
func (my *ValidatorApp[T]) checkFloat64(rule, fieldName string, value any) error {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		if reflect.ValueOf(value).IsNil() {
			return operation.Ternary(rule == "required", RequiredErr.New(fieldName), nil)
		}
		value = reflect.ValueOf(value).Elem().Interface()
	}

	switch {
	case strings.HasPrefix(rule, "size<="):
		small := strings.TrimPrefix(rule, "size<=")
		if !(value.(float64) <= common.ToFloat64(small)) {
			return LengthErr.NewFormat("[%s]长度不能小于等于[%f]", fieldName, common.ToFloat64(small))
		}
	case strings.HasPrefix(rule, "size<"):
		small := strings.TrimPrefix(rule, "size<")
		if !(value.(float64) < common.ToFloat64(small)) {
			return LengthErr.NewFormat("[%s]长度不能小于[%f]", fieldName, common.ToFloat64(small))
		}
	case strings.HasPrefix(rule, "size>="):
		large := strings.TrimPrefix(rule, "size>=")
		if !(value.(float64) >= common.ToFloat64(large)) {
			return LengthErr.NewFormat("[%s]长度不能大于等于[%f]", fieldName, common.ToFloat64(large))
		}
	case strings.HasPrefix(rule, "size>"):
		large := strings.TrimPrefix(rule, "size>")
		if !(value.(float64) > common.ToFloat64(large)) {
			return LengthErr.NewFormat("[%s]长度不能大于[%f]", fieldName, common.ToFloat64(large))
		}
	case strings.HasPrefix(rule, "size="):
		size := strings.TrimPrefix(rule, "size=")
		if value.(float64) != common.ToFloat64(size) {
			return LengthErr.NewFormat("[%s]长度必须等于[%f]", fieldName, common.ToFloat64(size))
		}
	case strings.HasPrefix(rule, "size!="):
		size := strings.TrimPrefix(rule, "size!=")
		if value.(float64) == common.ToFloat64(size) {
			return LengthErr.NewFormat("[%s]长度必须不等于[%f]", fieldName, common.ToFloat64(size))
		}
	case strings.HasPrefix(rule, "range="):
		between := strings.TrimPrefix(rule, "range=")
		betweenRange := strings.Split(between, ",")
		if len(betweenRange) != 2 {
			return RuleErr.NewFormat("[%s]规则定义错误#1，规则格式：f,f", fieldName)
		}
		small := common.ToFloat64(betweenRange[0])
		large := common.ToFloat64(betweenRange[1])
		if value.(float64) < small || value.(float64) > large {
			return LengthErr.NewFormat("[%s]长度必须在：%f~%f之间", fieldName, small, large)
		}
	}

	return nil
}
