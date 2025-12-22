package validatorV3

import (
	"reflect"
	"strings"

	"github.com/aid297/aid/array/anyArrayV2"
)

type (
	// Validator 验证器
	Validator struct {
		data   any
		wrongs []error
	}
)

func (Validator) New(data any) Validator { return Validator{data: data} }

// func (Validator) New(attrs ...ValidatorAttributer) Validator { return Validator{}.SetAttrs(attrs...) }
//
// func (my Validator) SetAttrs(attrs ...ValidatorAttributer) Validator {
// 	for idx := range attrs {
// 		attrs[idx].Register(&my)
// 	}
//
// 	return my
// }

func (my Validator) Wrongs() []error { return my.wrongs }

func (my Validator) Validate() Validator {
	fieldInfos := getStructFieldInfos(my.data, "")
	println(fieldInfos)
	for idx := range fieldInfos {
		if wrongs := fieldInfos[idx].Check().Wrongs(); len(wrongs) > 0 {
			my.wrongs = append(my.wrongs, wrongs...)
		}
	}

	return my
}

func getStructFieldInfos(s any, parentName string) []FieldInfo {
	v := reflect.ValueOf(s)
	t := v.Type()

	var infos []FieldInfo

	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}

	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldValue := v.Field(i)

			// 获取 tags
			vRuleTag := field.Tag.Get("v-rule")
			vNameTag := field.Tag.Get("v-name")

			if vNameTag == "" {
				vNameTag = field.Name
			}

			if vRuleTag == "" || vRuleTag == "-" {
				continue
			}

			isPtr := fieldValue.Kind() == reflect.Ptr
			isNil := isPtr && fieldValue.IsNil()

			// determine the element/type kind safely before dereferencing
			var elemType reflect.Type
			var elemKind reflect.Kind
			if isPtr {
				elemType = fieldValue.Type().Elem()
				elemKind = elemType.Kind()
			} else {
				elemType = fieldValue.Type()
				elemKind = elemType.Kind()
			}

			var value any = nil
			// only dereference if pointer and not nil, otherwise keep value as nil or the concrete value
			if isPtr && !isNil {
				fieldValue = fieldValue.Elem()
				value = fieldValue.Interface()
			} else if !isPtr {
				value = fieldValue.Interface()
			}

			// If the declared/element kind is struct, recurse.
			// For a nil pointer-to-struct, pass a zero value of the element type so reflection works.
			if elemKind == reflect.Struct {
				var recurseArg any
				if isPtr && isNil {
					recurseArg = reflect.Zero(elemType).Interface()
				} else {
					recurseArg = value
				}
				infos = append(infos, getStructFieldInfos(recurseArg, vNameTag)...)
				continue
			}

			infos = append(infos, FieldInfo{
				Name:      field.Name,
				Value:     value,
				Kind:      elemKind,
				IsPtr:     isPtr,
				IsNil:     isNil,
				VRuleTags: strings.Split(vRuleTag, ";"),
				VNameTags: anyArrayV2.NewItems(parentName, vNameTag).RemoveEmpty().ToSlice(),
			})
		}
	}

	return infos
}
