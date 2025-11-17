package validatorV2

import (
	"reflect"
)

type Validator struct {
	Errors          []error
	ValidatorFields []ValidatorField
	prefixNames     []string
}

func (Validator) New(original any, prefixNames ...string) Validator {
	var (
		ins = Validator{prefixNames: make([]string, 0, len(prefixNames)+1)}
		ref = reflect.ValueOf(original)
	)

	if ref.Kind() == reflect.Ptr {
		ref = ref.Elem()
	}

	if ref.Kind() != reflect.Struct {
		return Validator{}
	}

	ins.ValidatorFields = make([]ValidatorField, 0, ref.NumField())

	for i := range ref.NumField() {
		field := ref.Type().Field(i)
		val := ref.Field(i)
		if field.Anonymous {
			// 递归验证嵌套字段
			v := Validator{}.New(ref.Field(i).Interface(), ins.prefixNames...)
			if len(v.Errors) > 0 {
				ins.Errors = append(ins.Errors, v.Errors...)
			}
			continue
		}

		vRule := field.Tag.Get("v-rule")
		if vRule == "" || vRule == "-" {
			continue
		}

		vName := field.Tag.Get("v-name")
		if vName == "" || vName == "-" {
			continue
		}

		ins.ValidatorFields = append(ins.ValidatorFields, APP.ValidatorField.New(val, field))
	}

	if len(ins.ValidatorFields) > 0 {
		ins.Errors = make([]error, 0, len(ins.ValidatorFields))
		for idx := range ins.ValidatorFields {
			ins.Errors = append(ins.Errors, ins.ValidatorFields[idx].Error)
		}
	}

	return ins
}
