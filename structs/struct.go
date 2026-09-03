package structs

import (
	"reflect"
	"time"
)

var timeType = reflect.TypeFor[time.Time]()

// Copy 使用 dst 的字段覆盖 src 中同名字段（要求字段类型一致），仅处理第一层字段，不做递归
// 嵌套结构体字段（非time.Time，含指针形式）一律直接跳过；skipFields 中列出的第一层字段名也会被跳过
// src 必须是非 nil 的结构体指针，否则 panic；dst 可为结构体值或结构体指针
func Copy(src, dst any, skipFields ...string) {
	var (
		srcValue = reflect.ValueOf(src)
		dstValue = reflect.ValueOf(dst)
	)

	if srcValue.Kind() != reflect.Pointer || srcValue.IsNil() {
		panic("structs.Copy: src必须是结构体指针")
	}

	if srcValue.Elem().Kind() != reflect.Struct {
		panic("structs.Copy: src必须是结构体指针")
	}

	if dstValue.Kind() == reflect.Pointer {
		if dstValue.IsNil() {
			return
		}
		dstValue = dstValue.Elem()
	}

	if dstValue.Kind() != reflect.Struct {
		return
	}

	copyFields(srcValue.Elem(), dstValue, skipFields)
}

// copyFields 遍历 src 的可设置字段，用 dst 中同名同类型字段整体覆盖
func copyFields(src, dst reflect.Value, skipFields []string) {
	var srcType = src.Type()

	for i := 0; i < srcType.NumField(); i++ {
		var fieldName = srcType.Field(i).Name

		if inSlice(fieldName, skipFields) {
			continue // 跳过字段列表中的字段
		}

		var srcField = src.Field(i)
		if !srcField.CanSet() {
			continue // 未导出字段无法覆盖
		}

		if isNestedStruct(srcField.Type()) {
			continue // 嵌套结构体字段（非time.Time）直接跳过，避免浅拷贝共享底层对象
		}

		var dstField = dst.FieldByName(fieldName)
		if !dstField.IsValid() || dstField.Type() != srcField.Type() {
			continue // dst没有同名字段或类型不一致
		}

		srcField.Set(dstField) // 同名同类型字段整体覆盖（浅拷贝，不递归）
	}
}

// isNestedStruct 判断字段类型是否是嵌套结构体（非time.Time），指针会解引用到底后判断
func isNestedStruct(fieldType reflect.Type) bool {
	for fieldType.Kind() == reflect.Pointer {
		fieldType = fieldType.Elem()
	}

	return fieldType.Kind() == reflect.Struct && fieldType != timeType
}

// inSlice 判断字符串是否在切片中
func inSlice(target string, values []string) bool {
	for idx := range values {
		if values[idx] == target {
			return true
		}
	}

	return false
}
