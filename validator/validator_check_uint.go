package validator

import (
	"reflect"
	"strings"

	"github.com/spf13/cast"

	"github.com/aid297/aid/array"
	"github.com/aid297/aid/operation"

	"github.com/aid297/aid/common"
)

// checkUint 验证：uint -> 支持的规则 required、size<、size<=、size>、size>=、range=
func (my *ValidatorApp[T]) checkUint(rule, fieldName string, value any) error {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		if reflect.ValueOf(value).IsNil() {
			return operation.Ternary(rule == "required", RequiredErr.New(fieldName), nil)
		}
		value = reflect.ValueOf(value).Elem().Interface()
	}

	switch {
	case strings.HasPrefix(rule, "size<="):
		small := strings.TrimPrefix(rule, "size<=")
		if value.(uint) <= common.ToUint(small) {
			return LengthErr.NewFormat("[%s]长度不能小于等于[%d]", fieldName, common.ToInt(small))
		}
	case strings.HasPrefix(rule, "size<"):
		small := strings.TrimPrefix(rule, "size<")
		if value.(uint) < common.ToUint(small) {
			return LengthErr.NewFormat("[%s]长度不能小于[%d]", fieldName, common.ToInt(small))
		}
	case strings.HasPrefix(rule, "size>="):
		large := strings.TrimPrefix(rule, "size>=")
		if value.(uint) >= common.ToUint(large) {
			return LengthErr.NewFormat("[%s]长度不能大于等于[%d]", fieldName, common.ToInt(large))
		}
	case strings.HasPrefix(rule, "size>"):
		large := strings.TrimPrefix(rule, "size>")
		if value.(uint) > common.ToUint(large) {
			return LengthErr.NewFormat("[%s]长度不能大于[%d]", fieldName, common.ToInt(large))
		}
	case strings.HasPrefix(rule, "size="):
		size := strings.TrimPrefix(rule, "size=")
		if value.(uint) != common.ToUint(size) {
			return LengthErr.NewFormat("[%s]长度必须等于[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "size!="):
		size := strings.TrimPrefix(rule, "size!=")
		if value.(uint) == common.ToUint(size) {
			return LengthErr.NewFormat("[%s]长度必须不等于[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "range="):
		between := strings.TrimPrefix(rule, "range=")
		betweenRange := strings.Split(between, ",")
		if len(betweenRange) != 2 {
			return RuleErr.NewFormat("[%s]规则定义错误#4，规则格式[d,d]", fieldName)
		}
		small := common.ToUint(betweenRange[0])
		large := common.ToUint(betweenRange[1])
		if value.(uint) < small || value.(uint) > large {
			return LengthErr.NewFormat("[%s]长度必须在[%d~%d]之间", fieldName, small, large)
		}
	case strings.HasPrefix(rule, "in="):
		inValuesStr := strings.TrimPrefix(rule, "in=")
		inValuesArr := array.Cast(array.New(strings.Split(inValuesStr, ",")), func(value string) uint { return cast.ToUint(value) })
		if !inValuesArr.In(value.(uint)) {
			return ValidateErr.NewFormat("[%s]值必须在[%s]中", fieldName, inValuesStr)
		}
	case strings.HasPrefix(rule, "not-in="):
		inValuesStr := strings.TrimPrefix(rule, "not-in=")
		inValuesArr := array.Cast(array.New(strings.Split(inValuesStr, ",")), func(value string) uint { return cast.ToUint(value) })
		if inValuesArr.In(value.(uint)) {
			return ValidateErr.NewFormat("[%s]值不可为以下内容：[%s]", fieldName, inValuesStr)
		}
	}

	return nil
}

// checkUint8 验证：uint8 -> 支持的规则 required、size<、size<=、size>、size>=、range=
func (my *ValidatorApp[T]) checkUint8(rule, fieldName string, value any) error {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		if reflect.ValueOf(value).IsNil() {
			return operation.Ternary(rule == "required", RequiredErr.New(fieldName), nil)
		}
		value = reflect.ValueOf(value).Elem().Interface()
	}

	switch {
	case strings.HasPrefix(rule, "size<="):
		small := strings.TrimPrefix(rule, "size<=")
		if value.(uint8) <= common.ToUint8(small) {
			return LengthErr.NewFormat("[%s]长度不能小于等于[%d]", fieldName, common.ToInt(small))
		}
	case strings.HasPrefix(rule, "size<"):
		small := strings.TrimPrefix(rule, "size<")
		if value.(uint8) < common.ToUint8(small) {
			return LengthErr.NewFormat("[%s]长度不能小于[%d]", fieldName, common.ToInt(small))
		}
	case strings.HasPrefix(rule, "size>="):
		large := strings.TrimPrefix(rule, "size>=")
		if value.(uint8) >= common.ToUint8(large) {
			return LengthErr.NewFormat("[%s]长度不能大于等于[%d]", fieldName, common.ToInt(large))
		}
	case strings.HasPrefix(rule, "size>"):
		large := strings.TrimPrefix(rule, "size>")
		if value.(uint8) > common.ToUint8(large) {
			return LengthErr.NewFormat("[%s]长度不能大于[%d]", fieldName, common.ToInt(large))
		}
	case strings.HasPrefix(rule, "size="):
		size := strings.TrimPrefix(rule, "size=")
		if value.(uint8) != common.ToUint8(size) {
			return LengthErr.NewFormat("[%s]长度必须等于[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "size!="):
		size := strings.TrimPrefix(rule, "size!=")
		if value.(uint8) == common.ToUint8(size) {
			return LengthErr.NewFormat("[%s]长度必须不等于[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "range="):
		between := strings.TrimPrefix(rule, "range=")
		betweenRange := strings.Split(between, ",")
		if len(betweenRange) != 2 {
			return RuleErr.NewFormat("[%s]规则定义错误#4，规则格式[d,d]", fieldName)
		}
		small := common.ToUint8(betweenRange[0])
		large := common.ToUint8(betweenRange[1])
		if value.(uint8) < small || value.(uint8) > large {
			return LengthErr.NewFormat("[%s]长度必须在[%d~%d]之间", fieldName, small, large)
		}
	}

	return nil
}

// checkUint16 验证：uint16 -> 支持的规则 required、size<、size<=、size>、size>=、range=
func (my *ValidatorApp[T]) checkUint16(rule, fieldName string, value any) error {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		if reflect.ValueOf(value).IsNil() {
			return operation.Ternary(rule == "required", RequiredErr.New(fieldName), nil)
		}
		value = reflect.ValueOf(value).Elem().Interface()
	}

	switch {
	case strings.HasPrefix(rule, "size<="):
		small := strings.TrimPrefix(rule, "size<=")
		if value.(uint16) <= common.ToUint16(small) {
			return LengthErr.NewFormat("[%s]长度不能小于等于[%d]", fieldName, common.ToInt(small))
		}
	case strings.HasPrefix(rule, "size<"):
		small := strings.TrimPrefix(rule, "size<")
		if value.(uint16) < common.ToUint16(small) {
			return LengthErr.NewFormat("[%s]长度不能小于[%d]", fieldName, common.ToInt(small))
		}
	case strings.HasPrefix(rule, "size>="):
		large := strings.TrimPrefix(rule, "size>=")
		if value.(uint16) >= common.ToUint16(large) {
			return LengthErr.NewFormat("[%s]长度不能大于等于[%d]", fieldName, common.ToInt(large))
		}
	case strings.HasPrefix(rule, "size>"):
		large := strings.TrimPrefix(rule, "size>")
		if value.(uint16) > common.ToUint16(large) {
			return LengthErr.NewFormat("[%s]长度不能大于[%d]", fieldName, common.ToInt(large))
		}
	case strings.HasPrefix(rule, "size="):
		size := strings.TrimPrefix(rule, "size=")
		if value.(uint16) != common.ToUint16(size) {
			return LengthErr.NewFormat("[%s]长度必须等于[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "size!="):
		size := strings.TrimPrefix(rule, "size!=")
		if value.(uint16) == common.ToUint16(size) {
			return LengthErr.NewFormat("[%s]长度必须不等于[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "range="):
		between := strings.TrimPrefix(rule, "range=")
		betweenRange := strings.Split(between, ",")
		if len(betweenRange) != 2 {
			return RuleErr.NewFormat("[%s]规则定义错误#4，规则格式[d,d]", fieldName)
		}
		small := common.ToUint16(betweenRange[0])
		large := common.ToUint16(betweenRange[1])
		if value.(uint16) < small || value.(uint16) > large {
			return LengthErr.NewFormat("[%s]长度必须在[%d~%d]之间", fieldName, small, large)
		}
	}

	return nil
}

// checkUint32 验证：uint32 -> 支持的规则 required、size<、size<=、size>、size>=、range=
func (my *ValidatorApp[T]) checkUint32(rule, fieldName string, value any) error {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		if reflect.ValueOf(value).IsNil() {
			return operation.Ternary(rule == "required", RequiredErr.New(fieldName), nil)
		}
		value = reflect.ValueOf(value).Elem().Interface()
	}

	switch {
	case strings.HasPrefix(rule, "size<="):
		small := strings.TrimPrefix(rule, "size<=")
		if value.(uint32) <= common.ToUint32(small) {
			return LengthErr.NewFormat("[%s]长度不能小于等于[%d]", fieldName, common.ToInt(small))
		}
	case strings.HasPrefix(rule, "size<"):
		small := strings.TrimPrefix(rule, "size<")
		if value.(uint32) < common.ToUint32(small) {
			return LengthErr.NewFormat("[%s]长度不能小于[%d]", fieldName, common.ToInt(small))
		}
	case strings.HasPrefix(rule, "size>="):
		large := strings.TrimPrefix(rule, "size>=")
		if value.(uint32) >= common.ToUint32(large) {
			return LengthErr.NewFormat("[%s]长度不能大于等于[%d]", fieldName, common.ToInt(large))
		}
	case strings.HasPrefix(rule, "size>"):
		large := strings.TrimPrefix(rule, "size>")
		if value.(uint32) > common.ToUint32(large) {
			return LengthErr.NewFormat("[%s]长度不能大于[%d]", fieldName, common.ToInt(large))
		}
	case strings.HasPrefix(rule, "size="):
		size := strings.TrimPrefix(rule, "size=")
		if value.(uint32) != common.ToUint32(size) {
			return LengthErr.NewFormat("[%s]长度必须等于[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "size!="):
		size := strings.TrimPrefix(rule, "size!=")
		if value.(uint32) == common.ToUint32(size) {
			return LengthErr.NewFormat("[%s]长度必须不等于[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "range="):
		between := strings.TrimPrefix(rule, "range=")
		betweenRange := strings.Split(between, ",")
		if len(betweenRange) != 2 {
			return RuleErr.NewFormat("[%s]规则定义错误#4，规则格式[d,d]", fieldName)
		}
		small := common.ToUint32(betweenRange[0])
		large := common.ToUint32(betweenRange[1])
		if value.(uint32) < small || value.(uint32) > large {
			return LengthErr.NewFormat("[%s]长度必须在[%d~%d]之间", fieldName, small, large)
		}
	}

	return nil
}

// checkUint64 验证：uint64 -> 支持的规则 required、size<、size<=、size>、size>=、range=
func (my *ValidatorApp[T]) checkUint64(rule, fieldName string, value any) error {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		if reflect.ValueOf(value).IsNil() {
			return operation.Ternary(rule == "required", RequiredErr.New(fieldName), nil)
		}
		value = reflect.ValueOf(value).Elem().Interface()
	}

	switch {
	case strings.HasPrefix(rule, "size<="):
		small := strings.TrimPrefix(rule, "size<=")
		if value.(uint64) <= common.ToUint64(small) {
			return LengthErr.NewFormat("[%s]长度不能小于等于[%d]", fieldName, common.ToInt(small))
		}
	case strings.HasPrefix(rule, "size<"):
		small := strings.TrimPrefix(rule, "size<")
		if value.(uint64) < common.ToUint64(small) {
			return LengthErr.NewFormat("[%s]长度不能小于[%d]", fieldName, common.ToInt(small))
		}
	case strings.HasPrefix(rule, "size>="):
		large := strings.TrimPrefix(rule, "size>=")
		if value.(uint64) >= common.ToUint64(large) {
			return LengthErr.NewFormat("[%s]长度不能大于等于[%d]", fieldName, common.ToInt(large))
		}
	case strings.HasPrefix(rule, "size>"):
		large := strings.TrimPrefix(rule, "size>")
		if value.(uint64) > common.ToUint64(large) {
			return LengthErr.NewFormat("[%s]长度不能大于[%d]", fieldName, common.ToInt(large))
		}
	case strings.HasPrefix(rule, "size="):
		size := strings.TrimPrefix(rule, "size=")
		if value.(uint64) != common.ToUint64(size) {
			return LengthErr.NewFormat("[%s]长度必须等于[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "size!="):
		size := strings.TrimPrefix(rule, "size!=")
		if value.(uint64) == common.ToUint64(size) {
			return LengthErr.NewFormat("[%s]长度必须不等于[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "range="):
		between := strings.TrimPrefix(rule, "range=")
		betweenRange := strings.Split(between, ",")
		if len(betweenRange) != 2 {
			return RuleErr.NewFormat("[%s]规则定义错误#4，规则格式[d,d]", fieldName)
		}
		small := common.ToUint64(betweenRange[0])
		large := common.ToUint64(betweenRange[1])
		if value.(uint64) < small || value.(uint64) > large {
			return LengthErr.NewFormat("[%s]长度必须在[%d~%d]之间", fieldName, small, large)
		}
	}

	return nil
}
