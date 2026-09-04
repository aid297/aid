package structs

import (
	"reflect"
	"slices"
	"time"
)

var timeType = reflect.TypeFor[time.Time]()

type (
	// Struct 结构体字段覆盖器：使用 dst 的字段覆盖 src 中同名字段（要求字段类型一致）
	Struct[T any] struct {
		srcValue reflect.Value
		dstValue reflect.Value
	}
)

// NewStruct 实例化：src 与 dst 都必须是非 nil 的结构体值（传指针或其他类型会 panic）；T 由 src 推断，dst 可为不同结构体类型
func NewStruct[T any](src T, dst any) *Struct[T] {
	var (
		srcValue = reflect.ValueOf(src)
		dstValue = reflect.ValueOf(dst)
	)

	if srcValue.Kind() == reflect.Pointer || srcValue.Kind() != reflect.Struct {
		panic("structs.NewStruct: src必须是非指针的结构体")
	}

	if dstValue.Kind() == reflect.Pointer || dstValue.Kind() != reflect.Struct {
		panic("structs.NewStruct: dst必须是非指针的结构体")
	}

	return &Struct[T]{srcValue: srcValue, dstValue: dstValue}
}

// Cover 覆盖所有符合条件的字段（同名同类型、非嵌套结构体、可设置），返回修改后的src
func (my *Struct[T]) Cover() T {
	return my.cover(true, nil)
}

// CoverBySkip 黑名单方式覆盖：跳过 skipFields 中列出的第一层字段，覆盖其余符合条件的字段，返回修改后的src
func (my *Struct[T]) CoverBySkip(skipFields ...string) T {
	return my.cover(true, skipFields)
}

// CoverByAssign 白名单方式覆盖：仅覆盖 assignFields 中列出的第一层字段（仍需同名同类型等条件），返回修改后的src
func (my *Struct[T]) CoverByAssign(assignFields ...string) T {
	return my.cover(false, assignFields)
}

// cover 覆盖核心：skip 为 true 时 fields 是黑名单，为 false 时 fields 是白名单；在src副本上执行覆盖并返回
func (my *Struct[T]) cover(skip bool, fields []string) T {
	// src以值传入，先拷贝到可设置的副本上再修改，返回值即为修改后的src
	var src = reflect.New(my.srcValue.Type()).Elem()
	src.Set(my.srcValue)

	var srcType = src.Type()

	for i := 0; i < srcType.NumField(); i++ {
		var fieldName = srcType.Field(i).Name

		if slices.Contains(fields, fieldName) == skip {
			continue // 黑名单：列表中的跳过；白名单：不在列表中的跳过
		}

		var srcField = src.Field(i)
		if !srcField.CanSet() {
			continue // 未导出字段无法覆盖
		}

		var dstField = my.dstValue.FieldByName(fieldName)
		if !dstField.IsValid() {
			continue // dst没有同名字段
		}

		if dstField.Type() == srcField.Type() {
			if isNestedStruct(srcField.Type()) {
				continue // 嵌套结构体字段（非time.Time）直接跳过，避免浅拷贝共享底层对象
			}

			srcField.Set(dstField) // 类型完全一致：整体覆盖（浅拷贝，不递归）
			continue
		}

		my.coverByDeref(srcField, dstField) // 类型不一致：尝试指针解引用匹配
	}

	return src.Interface().(T) // src副本由src的类型构造，断言回T必定成功
}

// coverByDeref 指针/值双向解引用覆盖：当src与dst的同名字段一方是指针、另一方是值（解引用到底后类型相同）时生效
// src为*T、dst为T时把dst的值写入src指向的对象（src为nil则逐层新建独立对象）；
// src为T、dst为*T时用dst指向的值覆盖src（dst为nil则用零值覆盖）
func (my *Struct[T]) coverByDeref(srcField, dstField reflect.Value) {
	if derefType(srcField.Type()) != derefType(dstField.Type()) {
		return // 解引用到底后类型仍不同，无法匹配
	}

	// 归一化src：逐层解引用，遇到nil则新建独立对象并让指针指向它
	var srcTarget = srcField
	for srcTarget.Kind() == reflect.Pointer {
		if srcTarget.IsNil() {
			var newPtr = reflect.New(srcTarget.Type().Elem())
			srcTarget.Set(newPtr)
		}
		srcTarget = srcTarget.Elem()
	}

	// 归一化dst：逐层解引用，遇到nil则取零值（dst侧只读，不新建对象）
	var dstTarget = dstField
	for dstTarget.Kind() == reflect.Pointer {
		if dstTarget.IsNil() {
			dstTarget = reflect.Zero(dstTarget.Type().Elem())
			break
		}
		dstTarget = dstTarget.Elem()
	}

	srcTarget.Set(dstTarget)
}

// derefType 返回解引用到底后的类型
func derefType(fieldType reflect.Type) reflect.Type {
	for fieldType.Kind() == reflect.Pointer {
		fieldType = fieldType.Elem()
	}

	return fieldType
}

// isNestedStruct 判断字段类型是否是嵌套结构体（非time.Time），指针会解引用到底后判断
func isNestedStruct(fieldType reflect.Type) bool {
	var finalType = derefType(fieldType)

	return finalType.Kind() == reflect.Struct && finalType != timeType
}
