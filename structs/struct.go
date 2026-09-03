package structs

import (
	"reflect"
	"time"
)

var timeType = reflect.TypeFor[time.Time]()

// Copy 使用 dst 的字段覆盖 src 中同名字段（要求字段类型一致），嵌套结构体会递归替换
// src 必须是非 nil 的结构体指针，否则 panic；dst 可为结构体值或结构体指针
func Copy(src, dst any) {
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

	copyFields(srcValue.Elem(), dstValue)
}

// copyFields 遍历 src 的可设置字段，用 dst 中同名同类型字段覆盖
func copyFields(src, dst reflect.Value) {
	var srcType = src.Type()

	for i := 0; i < srcType.NumField(); i++ {
		var srcField = src.Field(i)
		if !srcField.CanSet() {
			continue // 未导出字段无法覆盖
		}

		var dstField = dst.FieldByName(srcType.Field(i).Name)
		if !dstField.IsValid() || dstField.Type() != srcField.Type() {
			continue // dst没有同名字段或类型不一致
		}

		copyValue(srcField, dstField)
	}
}

// copyValue 按类型分发覆盖
func copyValue(src, dst reflect.Value) {
	switch src.Kind() {
	case reflect.Struct:
		if src.Type() == timeType {
			src.Set(dst) // time.Time字段全部私有，递归无意义，整体覆盖
			return
		}
		copyFields(src, dst) // 嵌套结构体递归替换
	case reflect.Pointer:
		if src.IsNil() || dst.IsNil() {
			src.Set(dst) // 任一方为nil无法递归，直接整体覆盖
			return
		}
		copyValue(src.Elem(), dst.Elem()) // 解引用后继续处理，天然支持多级指针
	default:
		src.Set(dst)
	}
}
