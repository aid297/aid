package validatorV2

import (
	"errors"
	"reflect"
)

type (
	Validator struct{}
)

func NewValidator(original any) (*OriginalDatum, error) {
	ref := reflect.ValueOf(original)

	if ref.Kind() != reflect.Struct {
		return nil, errors.New("待验证数据必须是结构体或它的指针")
	}

	for i := range ref.NumField() {
		field := ref.Type().Field(i)
		val := ref.Field(i)
		// if field.Anonymous {
		// 	// 递归验证嵌套字段
		// 	if err := New(ref.Field(i).Interface(), my.prefixNames...).Validate(); err != nil {
		// 		return err
		// 	}
		// 	continue
		// }

		vRule := field.Tag.Get("v-rule")
		if vRule == "" || vRule == "-" {
			continue
		}

		vName := field.Tag.Get("v-name")
		if vName == "" || vName == "-" {
			continue
		}

		originalDatum := NewOriginalDatum(val, field)
		return originalDatum, nil
	}

	return nil, nil
}
