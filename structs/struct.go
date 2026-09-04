package structs

import (
	"fmt"
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
		dst      T
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

// CoverBySkip 黑名单方式覆盖：跳过 skipFields 中列出的第一层字段，覆盖其余符合条件的字段，返回修改后的src
func (my *Struct[T]) CoverBySkip(skipFields ...string) *Struct[T] {
	return my.cover(true, skipFields)
}

// CoverByAssign 白名单方式覆盖：仅覆盖 assignFields 中列出的第一层字段（仍需同名同类型等条件），返回修改后的src
func (my *Struct[T]) CoverByAssign(assignFields ...string) *Struct[T] {
	return my.cover(false, assignFields)
}

// cover 覆盖核心：skip 为 true 时 fields 是黑名单，为 false 时 fields 是白名单；在src副本上执行覆盖并返回
func (my *Struct[T]) cover(skip bool, fields []string) *Struct[T] {
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

	my.dst = src.Interface().(T) // src副本由src的类型构造，断言回T必定成功

	return my
}

func (my *Struct[T]) Value() T { return my.dst }

func (my *Struct[T]) Pointer() *T { return &my.dst }

// CompareBySkip 黑名单方式比较：跳过 skipFields 中列出的第一层字段，比较其余同名同类型字段的值内容，返回存在差异的字段记录
func (my *Struct[T]) CompareBySkip(skipFields ...string) []string {
	return my.compare(true, skipFields)
}

// CompareByAssign 白名单方式比较：仅比较 assignFields 中列出的第一层字段（仍需同名、解引用后同类型等条件），返回存在差异的字段记录
func (my *Struct[T]) CompareByAssign(assignFields ...string) []string {
	return my.compare(false, assignFields)
}

// compare 比较核心：逐字段比较src与dst的值内容（指针解引用取值，nil指针取零值，结构体及其指针跳过），返回差异记录，格式为"name标签值(或字段名):src的字段值 -> dst的字段值"
func (my *Struct[T]) compare(skip bool, fields []string) []string {
	var diffs []string
	var srcType = my.srcValue.Type()

	for i := 0; i < srcType.NumField(); i++ {
		var field = srcType.Field(i)
		var fieldName = field.Name

		if slices.Contains(fields, fieldName) == skip {
			continue // 黑名单：列表中的跳过；白名单：不在列表中的跳过
		}

		if !field.IsExported() {
			continue // 未导出字段不比较
		}

		var dstField = my.dstValue.FieldByName(fieldName)
		if !dstField.IsValid() || !dstField.CanInterface() {
			continue // dst没有同名字段或字段不可导出
		}

		var srcField = my.srcValue.Field(i)
		if isNestedStruct(srcField.Type()) || isNestedStruct(dstField.Type()) {
			continue // 结构体或结构体指针字段（非time.Time）跳过
		}

		if derefType(srcField.Type()) != derefType(dstField.Type()) {
			continue // 解引用到底后类型仍不同，无法比较
		}

		var srcValue = derefValue(srcField)
		var dstValue = derefValue(dstField)
		if !srcValue.Equal(dstValue) {
			diffs = append(diffs, fieldNameOf(field)+":"+fmt.Sprint(srcValue.Interface())+" -> "+fmt.Sprint(dstValue.Interface()))
		}
	}

	return diffs
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

// derefValue 解引用到底取值：逐层解引用指针，遇到nil指针取对应元素类型的零值
func derefValue(fieldValue reflect.Value) reflect.Value {
	for fieldValue.Kind() == reflect.Pointer {
		if fieldValue.IsNil() {
			return reflect.Zero(fieldValue.Type().Elem())
		}
		fieldValue = fieldValue.Elem()
	}

	return fieldValue
}

// fieldNameOf 取字段显示名：优先读tag name的值，没有则直接用字段名
func fieldNameOf(field reflect.StructField) string {
	if name := field.Tag.Get("name"); name != "" {
		return name
	}

	return field.Name
}

// isNestedStruct 判断字段类型是否是嵌套结构体（非time.Time），指针会解引用到底后判断
func isNestedStruct(fieldType reflect.Type) bool {
	var finalType = derefType(fieldType)

	return finalType.Kind() == reflect.Struct && finalType != timeType
}
