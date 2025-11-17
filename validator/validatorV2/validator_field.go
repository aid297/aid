package validatorV2

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/aid297/aid/ptr"
	"github.com/aid297/aid/regexp"
	"github.com/spf13/cast"
)

type (
	ValidatorField struct {
		refValue reflect.Value
		refType  reflect.Type
		refField reflect.StructField
		isPtr    bool
		isNil    bool
		isZero   bool
		indirect reflect.Value
		vType    string
		vRule    string
		vName    string
		Error    error
	}
)

func (ValidatorField) New(value reflect.Value, field reflect.StructField, vName string) ValidatorField {
	var (
		ins ValidatorField = ValidatorField{
			refValue: value,
			refType:  value.Type(),
			refField: field,
			isPtr:    value.Kind() == reflect.Ptr,
			isNil:    value.Kind() == reflect.Ptr && value.IsNil(),
			isZero:   value.IsZero(),
			indirect: reflect.Indirect(value),
			vType:    field.Tag.Get("v-type"),
			vRule:    field.Tag.Get("v-rule"),
			vName:    value.Type().Name(),
		}
		checker Checker
		err     error
	)

	if checker, err = ins.parseRule(); err != nil {
		ins.Error = fmt.Errorf(`"%s"%w`, ins.vName, err)
	} else if checker != nil {
		if err := checker.Check(); err != nil {
			ins.Error = fmt.Errorf(`"%s"%w`, ins.vName, err)
		}
	}
	return ins
}

func (my ValidatorField) parseRuleString(rules []string) CheckerString {
	r := CheckerString{original: my.refValue.String()}

	for idx := range rules {
		if regexp.APP.Regexp.New(`required;`, regexp.TargetString(rules[idx])).Contains() {
			r.required = true
		}

		if regexp.APP.Regexp.New(`not-zero;`, regexp.TargetString(rules[idx])).Contains() {
			r.noZero = true
		}

		if after, ok := strings.CutPrefix(rules[idx], "min:"); ok {
			r.min = ptr.New(cast.ToInt(after))
		}

		if after, ok := strings.CutPrefix(rules[idx], "max:"); ok {
			r.max = ptr.New(cast.ToInt(after))
		}

		if after, ok := strings.CutPrefix(rules[idx], "in:"); ok {
			r.in = strings.Split(after, ",")
		}

		if after, ok := strings.CutPrefix(rules[idx], "not-in:"); ok {
			r.notIn = strings.Split(after, ",")
		}

		if after, ok := strings.CutPrefix(rules[idx], "regex:"); ok {
			r.regex = ptr.New(after)
		}
	}

	return r
}

func (my ValidatorField) parseRuleStringPtr(rules []string) Checker {
	r := CheckerStringPtr{original: my.refValue.Interface().(*string)}

	for idx := range rules {
		if regexp.APP.Regexp.New(`required;`, regexp.TargetString(rules[idx])).Contains() {
			r.required = true
		}

		if regexp.APP.Regexp.New(`not-zero;`, regexp.TargetString(rules[idx])).Contains() {
			r.noZero = true
		}

		if after, ok := strings.CutPrefix(rules[idx], "min:"); ok {
			r.min = ptr.New(cast.ToInt(after))
		}

		if after, ok := strings.CutPrefix(rules[idx], "max:"); ok {
			r.max = ptr.New(cast.ToInt(after))
		}

		if after, ok := strings.CutPrefix(rules[idx], "in:"); ok {
			r.in = strings.Split(after, ",")
		}

		if after, ok := strings.CutPrefix(rules[idx], "not-in:"); ok {
			r.notIn = strings.Split(after, ",")
		}

		if after, ok := strings.CutPrefix(rules[idx], "regex:"); ok {
			r.regex = ptr.New(after)
		}
	}

	return r
}

func (my ValidatorField) parseRuleInt64(rules []string) Checker {
	r := CheckerInt64{original: my.refValue.Int()}

	for idx := range rules {
		if regexp.APP.Regexp.New(`required;`, regexp.TargetString(rules[idx])).Contains() {
			r.required = true
		}

		if regexp.APP.Regexp.New(`not-zero;`, regexp.TargetString(rules[idx])).Contains() {
			r.noZero = true
		}

		if after, ok := strings.CutPrefix(rules[idx], "min:"); ok {
			r.min = ptr.New(cast.ToInt64(after))
		}

		if after, ok := strings.CutPrefix(rules[idx], "max:"); ok {
			r.max = ptr.New(cast.ToInt64(after))
		}

		if after, ok := strings.CutPrefix(rules[idx], "in:"); ok {
			r.in = cast.ToInt64Slice(strings.Split(after, ","))
		}

		if after, ok := strings.CutPrefix(rules[idx], "not-in:"); ok {
			r.notIn = cast.ToInt64Slice(strings.Split(after, ","))
		}
	}

	return r
}

func (my ValidatorField) parseRuleInt64Ptr(rules []string) Checker {
	r := CheckerInt64Ptr{original: my.refValue.Interface().(*int64)}

	for idx := range rules {
		if regexp.APP.Regexp.New(`required;`, regexp.TargetString(rules[idx])).Contains() {
			r.required = true
		}

		if regexp.APP.Regexp.New(`not-zero;`, regexp.TargetString(rules[idx])).Contains() {
			r.noZero = true
		}

		if after, ok := strings.CutPrefix(rules[idx], "min:"); ok {
			r.min = ptr.New(cast.ToInt64(after))
		}

		if after, ok := strings.CutPrefix(rules[idx], "max:"); ok {
			r.max = ptr.New(cast.ToInt64(after))
		}

		if after, ok := strings.CutPrefix(rules[idx], "in:"); ok {
			r.in = cast.ToInt64Slice(strings.Split(after, ","))
		}

		if after, ok := strings.CutPrefix(rules[idx], "not-in:"); ok {
			r.notIn = cast.ToInt64Slice(strings.Split(after, ","))
		}
	}

	return r
}

func (my ValidatorField) parseRuleUint64(rules []string) Checker {
	r := CheckerUint64{original: my.refValue.Uint()}

	for idx := range rules {
		if regexp.APP.Regexp.New(`required;`, regexp.TargetString(rules[idx])).Contains() {
			r.required = true
		}

		if regexp.APP.Regexp.New(`not-zero;`, regexp.TargetString(rules[idx])).Contains() {
			r.noZero = true
		}

		if after, ok := strings.CutPrefix(rules[idx], "min:"); ok {
			r.min = ptr.New(cast.ToUint64(after))
		}

		if after, ok := strings.CutPrefix(rules[idx], "max:"); ok {
			r.max = ptr.New(cast.ToUint64(after))
		}

		if after, ok := strings.CutPrefix(rules[idx], "in:"); ok {
			list := strings.Split(after, ",")
			r.in = make([]uint64, 0, len(list))
			for idx := range list {
				r.in = append(r.in, cast.ToUint64(list[idx]))
			}
		}

		if after, ok := strings.CutPrefix(rules[idx], "not-in:"); ok {
			list := strings.Split(after, ",")
			r.notIn = make([]uint64, 0, len(list))
			for idx := range list {
				r.notIn = append(r.notIn, cast.ToUint64(list[idx]))
			}
		}
	}

	return r
}

func (my ValidatorField) parseRuleUint64Ptr(rules []string) Checker {
	r := CheckerUint64Ptr{original: ptr.New(my.refValue.Uint())}

	for idx := range rules {
		if regexp.APP.Regexp.New(`required;`, regexp.TargetString(rules[idx])).Contains() {
			r.required = true
		}

		if regexp.APP.Regexp.New(`not-zero;`, regexp.TargetString(rules[idx])).Contains() {
			r.noZero = true
		}

		if after, ok := strings.CutPrefix(rules[idx], "min:"); ok {
			r.min = ptr.New(cast.ToUint64(after))
		}

		if after, ok := strings.CutPrefix(rules[idx], "max:"); ok {
			r.max = ptr.New(cast.ToUint64(after))
		}

		if after, ok := strings.CutPrefix(rules[idx], "in:"); ok {
			list := strings.Split(after, ",")
			r.in = make([]uint64, 0, len(list))
			for idx := range list {
				r.in = append(r.in, cast.ToUint64(list[idx]))
			}
		}

		if after, ok := strings.CutPrefix(rules[idx], "not-in:"); ok {
			list := strings.Split(after, ",")
			r.notIn = make([]uint64, 0, len(list))
			for idx := range list {
				r.notIn = append(r.notIn, cast.ToUint64(list[idx]))
			}
		}
	}

	return r
}

func (my ValidatorField) parseRuleFloat64(rules []string) Checker {
	r := CheckerFloat64{original: my.refValue.Float()}

	for idx := range rules {
		if regexp.APP.Regexp.New(`required;`, regexp.TargetString(rules[idx])).Contains() {
			r.required = true
		}

		if regexp.APP.Regexp.New(`not-zero;`, regexp.TargetString(rules[idx])).Contains() {
			r.noZero = true
		}

		if after, ok := strings.CutPrefix(rules[idx], "min:"); ok {
			r.min = ptr.New(cast.ToFloat64(after))
		}

		if after, ok := strings.CutPrefix(rules[idx], "max:"); ok {
			r.max = ptr.New(cast.ToFloat64(after))
		}

		if after, ok := strings.CutPrefix(rules[idx], "in:"); ok {
			r.in = append(r.in, cast.ToFloat64(strings.Split(after, ",")))
		}

		if after, ok := strings.CutPrefix(rules[idx], "not-in:"); ok {
			r.notIn = append(r.notIn, cast.ToFloat64(strings.Split(after, ",")))
		}
	}

	return r
}

func (my ValidatorField) parseRuleFloat64Ptr(rules []string) Checker {
	r := CheckerFloat64Ptr{original: ptr.New(my.refValue.Float())}

	for idx := range rules {
		if regexp.APP.Regexp.New(`required;`, regexp.TargetString(rules[idx])).Contains() {
			r.required = true
		}

		if regexp.APP.Regexp.New(`not-zero;`, regexp.TargetString(rules[idx])).Contains() {
			r.noZero = true
		}

		if after, ok := strings.CutPrefix(rules[idx], "min:"); ok {
			r.min = ptr.New(cast.ToFloat64(after))
		}

		if after, ok := strings.CutPrefix(rules[idx], "max:"); ok {
			r.max = ptr.New(cast.ToFloat64(after))
		}

		if after, ok := strings.CutPrefix(rules[idx], "in:"); ok {
			r.in = append(r.in, cast.ToFloat64(strings.Split(after, ",")))
		}

		if after, ok := strings.CutPrefix(rules[idx], "not-in:"); ok {
			r.notIn = append(r.notIn, cast.ToFloat64(strings.Split(after, ",")))
		}
	}

	return r
}

func (my ValidatorField) parseRuleBool(rules []string) Checker {
	r := CheckerBool{original: cast.ToBool(my.refValue.Interface())}

	for idx := range rules {
		if regexp.APP.Regexp.New(`required;`, regexp.TargetString(rules[idx])).Contains() {
			r.required = true
		}

		if regexp.APP.Regexp.New(`not-zero;`, regexp.TargetString(rules[idx])).Contains() {
			r.noZero = true
		}
	}

	return r
}

func (my ValidatorField) parseRuleBoolPtr(rules []string) Checker {
	r := CheckerBoolPtr{original: ptr.New(cast.ToBool(my.refValue.Interface()))}

	for idx := range rules {
		if regexp.APP.Regexp.New(`required;`, regexp.TargetString(rules[idx])).Contains() {
			r.required = true
		}

		if regexp.APP.Regexp.New(`not-zero;`, regexp.TargetString(rules[idx])).Contains() {
			r.noZero = true
		}
	}

	return r
}

func (my ValidatorField) parseRuleSlice(rules []string, ref reflect.Value) Checker {
	target := make([]any, 0, ref.Len())

	for idx := range ref.Len() {
		target = append(target, ref.Index(idx).Interface())
	}

	r := CheckerSlice{original: target}

	for idx := range rules {
		if regexp.APP.Regexp.New(`required;`, regexp.TargetString(rules[idx])).Contains() {
			r.required = true
		}

		if regexp.APP.Regexp.New(`not-zero;`, regexp.TargetString(rules[idx])).Contains() {
			r.noZero = true
		}

		if after, ok := strings.CutPrefix(rules[idx], "min:"); ok {
			r.min = ptr.New(cast.ToInt(after))
		}

		if after, ok := strings.CutPrefix(rules[idx], "max:"); ok {
			r.max = ptr.New(cast.ToInt(after))
		}

		if after, ok := strings.CutPrefix(rules[idx], "in:"); ok {
			list := strings.Split(after, ",")
			r.in = make([]any, 0, len(list))
			for idx := range list {
				r.in = append(r.in, any(list[idx]))
			}
		}

		if after, ok := strings.CutPrefix(rules[idx], "not-in:"); ok {
			list := strings.Split(after, ",")
			r.notIn = make([]any, 0, len(list))
			for idx := range list {
				r.notIn = append(r.notIn, any(list[idx]))
			}
		}
	}

	return r
}

func (my ValidatorField) parseRuleSlicePtr(rules []string, ref reflect.Value) Checker {
	target := make([]any, 0, ref.Len())

	for idx := range ref.Len() {
		target = append(target, ref.Index(idx).Interface())
	}

	r := CheckerSlicePtr{original: &target}

	for idx := range rules {
		if regexp.APP.Regexp.New(`required;`, regexp.TargetString(rules[idx])).Contains() {
			r.required = true
		}

		if regexp.APP.Regexp.New(`not-zero;`, regexp.TargetString(rules[idx])).Contains() {
			r.noZero = true
		}

		if after, ok := strings.CutPrefix(rules[idx], "min:"); ok {
			r.min = ptr.New(cast.ToInt(after))
		}

		if after, ok := strings.CutPrefix(rules[idx], "max:"); ok {
			r.max = ptr.New(cast.ToInt(after))
		}

		if after, ok := strings.CutPrefix(rules[idx], "in:"); ok {
			list := strings.Split(after, ",")
			r.in = make([]any, 0, len(list))
			for idx := range list {
				r.in = append(r.in, any(list[idx]))
			}
		}

		if after, ok := strings.CutPrefix(rules[idx], "not-in:"); ok {
			list := strings.Split(after, ",")
			r.notIn = make([]any, 0, len(list))
			for idx := range list {
				r.notIn = append(r.notIn, any(list[idx]))
			}
		}
	}

	return r
}

func (my ValidatorField) parseRule() (Checker, error) {
	rules := strings.Split(my.vRule, ";")

	if !my.isPtr {
		switch my.refType.Kind() {
		case reflect.String:
			if my.vType != "" && my.vType != "string" {
				return nil, ErrInvalidType
			}
			return my.parseRuleString(rules), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if my.vType != "" && my.vType != "int" {
				return nil, ErrInvalidType
			}
			return my.parseRuleInt64(rules), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if my.vType != "" && my.vType != "uint" {
				return nil, ErrInvalidType
			}
			return my.parseRuleUint64(rules), nil
		case reflect.Float32, reflect.Float64:
			if my.vType != "" && my.vType != "float" {
				return nil, ErrInvalidType
			}
			return my.parseRuleFloat64(rules), nil
		case reflect.Bool:
			if my.vType != "" && my.vType != "bool" {
				return nil, ErrInvalidType
			}
			return my.parseRuleBool(rules), nil
		case reflect.Array, reflect.Slice:
			if my.vType != "" && (my.vType != "array" && my.vType != "slice") {
				return nil, ErrInvalidType
			}
			return my.parseRuleSlice(rules, my.refValue), nil
		case reflect.Map:
		case reflect.Struct:
		case reflect.Func:
		}
	} else {
		switch my.indirect.Kind() {
		case reflect.String:
			if my.vType != "" && my.vType != "string" {
				return nil, ErrInvalidType
			}
			return my.parseRuleStringPtr(rules), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if my.vType != "" && my.vType != "int" {
				return nil, ErrInvalidType
			}
			return my.parseRuleInt64Ptr(rules), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if my.vType != "" && my.vType != "uint" {
				return nil, ErrInvalidType
			}
			return my.parseRuleUint64Ptr(rules), nil
		case reflect.Float32, reflect.Float64:
			if my.vType != "" && my.vType != "float" {
				return nil, ErrInvalidType
			}
			return my.parseRuleFloat64Ptr(rules), nil
		case reflect.Bool:
			if my.vType != "" && my.vType != "bool" {
				return nil, ErrInvalidType
			}
			return my.parseRuleBoolPtr(rules), nil
		case reflect.Array, reflect.Slice:
			if my.vType != "" && (my.vType != "array" && my.vType != "slice") {
				return nil, ErrInvalidType
			}
			return my.parseRuleSlicePtr(rules, my.refValue), nil
		case reflect.Map:
		case reflect.Struct:
		case reflect.Func:
		}
	}

	return nil, nil
}
