package validator

import (
	"reflect"
	"strings"

	"github.com/spf13/cast"

	"github.com/aid297/aid/array"
	"github.com/aid297/aid/operation"

	"github.com/aid297/aid/common"
)

// checkInt 验证：int -> 支持的规则 required、size<、size<=、size>、size>=、range=
func (my *ValidatorApp[T]) checkInt(rule, fieldName string, value any) error {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		if reflect.ValueOf(value).IsNil() {
			return operation.Ternary(rule == "required", RequiredErr.New(fieldName), nil)
		}
		value = reflect.ValueOf(value).Elem().Interface()
	}

	switch {
	case strings.HasPrefix(rule, "size<="):
		small := strings.TrimPrefix(rule, "size<=")
		if value.(int) <= common.ToInt(small) {
			return LengthErr.NewFormat("[%s]长度不能小于等于：[%d]", fieldName, common.ToInt(small))
		}
	case strings.HasPrefix(rule, "size<"):
		small := strings.TrimPrefix(rule, "size<")
		if value.(int) < common.ToInt(small) {
			return LengthErr.NewFormat("[%s]长度不能小于：[%d]", fieldName, common.ToInt(small))
		}
	case strings.HasPrefix(rule, "size>="):
		large := strings.TrimPrefix(rule, "size>=")
		if value.(int) >= common.ToInt(large) {
			return LengthErr.NewFormat("[%s]长度不能大于等于：[%d]", fieldName, common.ToInt(large))
		}
	case strings.HasPrefix(rule, "size>"):
		large := strings.TrimPrefix(rule, "size>")
		if value.(int) > common.ToInt(large) {
			return LengthErr.NewFormat("[%s]长度不能大于：[%d]", fieldName, common.ToInt(large))
		}
	case strings.HasPrefix(rule, "size="):
		size := strings.TrimPrefix(rule, "size=")
		if value.(int) != common.ToInt(size) {
			return LengthErr.NewFormat("[%s]长度必须等于：[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "size!="):
		size := strings.TrimPrefix(rule, "size!=")
		if value.(int) == common.ToInt(size) {
			return LengthErr.NewFormat("[%s]长度必须不等于：[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "range="):
		between := strings.TrimPrefix(rule, "range=")
		betweenRange := strings.Split(between, ",")
		if len(betweenRange) != 2 {
			return RuleErr.NewFormat("[%s]规则定义错误#2，规则格式[d,d]", fieldName)
		}
		small := common.ToInt(betweenRange[0])
		large := common.ToInt(betweenRange[1])
		if value.(int) < small || value.(int) > large {
			return LengthErr.NewFormat("[%s]长度必须在[%d~%d]之间", fieldName, small, large)
		}
	}

	return nil
}

// checkInt8 验证：int8 -> 支持的规则 required、size<、size<=、size>、size>=、range=
func (my *ValidatorApp[T]) checkInt8(rule, fieldName string, value any) error {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		if reflect.ValueOf(value).IsNil() {
			return operation.Ternary(rule == "required", RequiredErr.New(fieldName), nil)
		}
		value = reflect.ValueOf(value).Elem().Interface()
	}

	switch {
	case strings.HasPrefix(rule, "size<="):
		small := strings.TrimPrefix(rule, "size<=")
		if value.(int8) <= common.ToInt8(small) {
			return LengthErr.NewFormat("[%s]长度不能小于等于：[%d]", fieldName, common.ToInt(small))
		}
	case strings.HasPrefix(rule, "size<"):
		small := strings.TrimPrefix(rule, "size<")
		if value.(int8) < common.ToInt8(small) {
			return LengthErr.NewFormat("[%s]长度不能小于：[%d]", fieldName, common.ToInt(small))
		}
	case strings.HasPrefix(rule, "size>="):
		large := strings.TrimPrefix(rule, "size>=")
		if value.(int8) >= common.ToInt8(large) {
			return LengthErr.NewFormat("[%s]长度不能大于等于：[%d]", fieldName, common.ToInt(large))
		}
	case strings.HasPrefix(rule, "size>"):
		large := strings.TrimPrefix(rule, "size>")
		if value.(int8) > common.ToInt8(large) {
			return LengthErr.NewFormat("[%s]长度不能大于：[%d]", fieldName, common.ToInt(large))
		}
	case strings.HasPrefix(rule, "size="):
		size := strings.TrimPrefix(rule, "size=")
		if value.(int8) != common.ToInt8(size) {
			return LengthErr.NewFormat("[%s]长度必须等于：[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "size!="):
		size := strings.TrimPrefix(rule, "size!=")
		if value.(int8) == common.ToInt8(size) {
			return LengthErr.NewFormat("[%s]长度必须不等于：[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "range="):
		between := strings.TrimPrefix(rule, "range=")
		betweenRange := strings.Split(between, ",")
		if len(betweenRange) != 2 {
			return RuleErr.NewFormat("[%s]规则定义错误#2，规则格式[d,d]", fieldName)
		}
		small := common.ToInt8(betweenRange[0])
		large := common.ToInt8(betweenRange[1])
		if value.(int8) < small || value.(int8) > large {
			return LengthErr.NewFormat("[%s]长度必须在[%d~%d]之间", fieldName, small, large)
		}
	case strings.HasPrefix(rule, "in="):
		inValuesStr := strings.TrimPrefix(rule, "in=")
		inValuesArr := array.Cast(array.New(strings.Split(inValuesStr, ",")), func(value string) int { return cast.ToInt(value) })
		if !inValuesArr.In(value.(int)) {
			return ValidateErr.NewFormat("[%s]值必须在[%s]中", fieldName, inValuesStr)
		}
	case strings.HasPrefix(rule, "not-in="):
		inValuesStr := strings.TrimPrefix(rule, "not-in=")
		inValuesArr := array.Cast(array.New(strings.Split(inValuesStr, ",")), func(value string) int { return cast.ToInt(value) })
		if inValuesArr.In(value.(int)) {
			return ValidateErr.NewFormat("[%s]值不可为以下内容：[%s]", fieldName, inValuesStr)
		}
	}

	return nil
}

// checkInt16 验证：int16 -> 支持的规则 required、size<、size<=、size>、size>=、range=
func (my *ValidatorApp[T]) checkInt16(rule, fieldName string, value any) error {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		if reflect.ValueOf(value).IsNil() {
			return operation.Ternary(rule == "required", RequiredErr.New(fieldName), nil)
		}
		value = reflect.ValueOf(value).Elem().Interface()
	}

	switch {
	case strings.HasPrefix(rule, "size<="):
		small := strings.TrimPrefix(rule, "size<=")
		if value.(int16) <= common.ToInt16(small) {
			return LengthErr.NewFormat("[%s]长度必须小于等于：[%d]", fieldName, common.ToInt(small))
		}
	case strings.HasPrefix(rule, "size<"):
		small := strings.TrimPrefix(rule, "size<")
		if value.(int16) < common.ToInt16(small) {
			return LengthErr.NewFormat("[%s]长度必须小于：[%d]", fieldName, common.ToInt(small))
		}
	case strings.HasPrefix(rule, "size>="):
		large := strings.TrimPrefix(rule, "size>=")
		if value.(int16) >= common.ToInt16(large) {
			return LengthErr.NewFormat("[%s]长度必须大于等于：[%d]", fieldName, common.ToInt(large))
		}
	case strings.HasPrefix(rule, "size>"):
		large := strings.TrimPrefix(rule, "size>")
		if !(value.(int16) > common.ToInt16(large)) {
			return LengthErr.NewFormat("[%s]长度必须大于：[%d]", fieldName, common.ToInt(large))
		}
	case strings.HasPrefix(rule, "size="):
		size := strings.TrimPrefix(rule, "size=")
		if !(value.(int16) != common.ToInt16(size)) {
			return LengthErr.NewFormat("[%s]长度必须等于：[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "size!="):
		size := strings.TrimPrefix(rule, "size!=")
		if !(value.(int16) == common.ToInt16(size)) {
			return LengthErr.NewFormat("[%s]长度必须不等于：[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "range="):
		between := strings.TrimPrefix(rule, "range=")
		betweenRange := strings.Split(between, ",")
		if len(betweenRange) != 2 {
			return RuleErr.NewFormat("[%s]规则定义错误#2，规则格式[d,d]", fieldName)
		}
		small := common.ToInt16(betweenRange[0])
		large := common.ToInt16(betweenRange[1])
		if value.(int16) < small || value.(int16) > large {
			return LengthErr.NewFormat("[%s]长度必须在[%d~%d]之间", fieldName, small, large)
		}
	}

	return nil
}

// checkInt32 验证：int32 -> 支持的规则 required、size<、size<=、size>、size>=、range=
func (my *ValidatorApp[T]) checkInt32(rule, fieldName string, value any) error {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		if reflect.ValueOf(value).IsNil() {
			return operation.Ternary(rule == "required", RequiredErr.New(fieldName), nil)
		}
		value = reflect.ValueOf(value).Elem().Interface()
	}

	switch {
	case strings.HasPrefix(rule, "size<="):
		small := strings.TrimPrefix(rule, "size<=")
		if value.(int32) <= common.ToInt32(small) {
			return LengthErr.NewFormat("[%s]长度不能小于等于：[%d]", fieldName, common.ToInt(small))
		}
	case strings.HasPrefix(rule, "size<"):
		small := strings.TrimPrefix(rule, "size<")
		if value.(int32) < common.ToInt32(small) {
			return LengthErr.NewFormat("[%s]长度不能小于：[%d]", fieldName, common.ToInt(small))
		}
	case strings.HasPrefix(rule, "size>="):
		large := strings.TrimPrefix(rule, "size>=")
		if value.(int32) >= common.ToInt32(large) {
			return LengthErr.NewFormat("[%s]长度不能大于等于：[%d]", fieldName, common.ToInt(large))
		}
	case strings.HasPrefix(rule, "size>"):
		large := strings.TrimPrefix(rule, "size>")
		if value.(int32) > common.ToInt32(large) {
			return LengthErr.NewFormat("[%s]长度不能大于：[%d]", fieldName, common.ToInt(large))
		}
	case strings.HasPrefix(rule, "size="):
		size := strings.TrimPrefix(rule, "size=")
		if value.(int32) != common.ToInt32(size) {
			return LengthErr.NewFormat("[%s]长度必须等于：[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "size!="):
		size := strings.TrimPrefix(rule, "size!=")
		if value.(int32) == common.ToInt32(size) {
			return LengthErr.NewFormat("[%s]长度必须不等于#2：[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "range="):
		between := strings.TrimPrefix(rule, "range=")
		betweenRange := strings.Split(between, ",")
		if len(betweenRange) != 2 {
			return RuleErr.NewFormat("[%s]规则定义错误，规则格式[d,d]", fieldName)
		}
		small := common.ToInt32(betweenRange[0])
		large := common.ToInt32(betweenRange[1])
		if value.(int32) < small || value.(int32) > large {
			return LengthErr.NewFormat("[%s]长度必须在[%d~%d]之间", fieldName, small, large)
		}
	}

	return nil
}

// checkInt64 验证：int64 -> 支持的规则 required、size<、size<=、size>、size>=、range=
func (my *ValidatorApp[T]) checkInt64(rule, fieldName string, value any) error {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		if reflect.ValueOf(value).IsNil() {
			return operation.Ternary(rule == "required", RequiredErr.New(fieldName), nil)
		}
		value = reflect.ValueOf(value).Elem().Interface()
	}

	switch {
	case strings.HasPrefix(rule, "size<="):
		small := strings.TrimPrefix(rule, "size<=")
		if value.(int64) <= common.ToInt64(small) {
			return LengthErr.NewFormat("[%s]长度不能小于等于[%d]", fieldName, common.ToInt(small))
		}
	case strings.HasPrefix(rule, "size<"):
		small := strings.TrimPrefix(rule, "size<")
		if value.(int64) < common.ToInt64(small) {
			return LengthErr.NewFormat("[%s]长度不能小于[%d]", fieldName, common.ToInt(small))
		}
	case strings.HasPrefix(rule, "size>="):
		large := strings.TrimPrefix(rule, "size>=")
		if value.(int64) >= common.ToInt64(large) {
			return LengthErr.NewFormat("[%s]长度不能大于等于[%d]", fieldName, common.ToInt(large))
		}
	case strings.HasPrefix(rule, "size>"):
		large := strings.TrimPrefix(rule, "size>")
		if value.(int64) > common.ToInt64(large) {
			return LengthErr.NewFormat("[%s]长度不能大于[%d]", fieldName, common.ToInt(large))
		}
	case strings.HasPrefix(rule, "size="):
		size := strings.TrimPrefix(rule, "size=")
		if value.(int64) != common.ToInt64(size) {
			return LengthErr.NewFormat("[%s]长度必须等于[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "size!="):
		size := strings.TrimPrefix(rule, "size!=")
		if value.(int64) == common.ToInt64(size) {
			return LengthErr.NewFormat("[%s]长度必须不等于[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "range="):
		between := strings.TrimPrefix(rule, "range=")
		betweenRange := strings.Split(between, ",")
		if len(betweenRange) != 2 {
			return RuleErr.NewFormat("[%s]规则定义错误#2，规则格式[d,d]", fieldName)
		}
		small := common.ToInt64(betweenRange[0])
		large := common.ToInt64(betweenRange[1])
		if value.(int64) < small || value.(int64) > large {
			return LengthErr.NewFormat("[%s]长度必须在[%d~%d]之间", fieldName, small, large)
		}
	}

	return nil
}
