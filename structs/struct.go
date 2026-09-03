package structs

import (
	"reflect"
)

// Copy 使用 dst 的字段覆盖 src 中同名字段（要求字段类型一致），仅处理第一层字段，不做递归
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

// copyFields 遍历 src 的可设置字段，用 dst 中同名同类型字段整体覆盖
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

		srcField.Set(dstField) // 同名同类型字段整体覆盖（浅拷贝，不递归）
	}
}
