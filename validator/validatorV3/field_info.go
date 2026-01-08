package validatorV3

import (
	"reflect"
	"time"

	"github.com/aid297/aid/array/anyArrayV2"
)

type (
	// FieldInfo 保存了字段的相关信息。
	FieldInfo struct {
		Name      string // 字段名
		Value     any    // 实际值
		Kind      reflect.Kind
		Type      reflect.Type
		IsPtr     bool                        // 是否是指针
		IsNil     bool                        // 是否为空指针
		IsZero    bool                        // 是否是零值
		Required  bool                        // 是否必填
		VRuleTags anyArrayV2.AnyArray[string] // v-rule tag 的值
		VNameTags anyArrayV2.AnyArray[string] // v-name tag 的值
		wrongs    []error
	}

	FieldRuleBase struct{ Required string }
	FieldRule     struct{ FieldRuleBase }
)

func (my FieldInfo) Wrongs() []error { return my.wrongs }

func (my FieldInfo) getName() string { return my.VNameTags.Join(".") }

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
