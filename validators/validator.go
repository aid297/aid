package validators

import (
	_errors "errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/aid297/aid/v2/anySlices"
)

var (
	defaultSliceSplitChar string = "|"
	defaultErrorSplitChar string = "<br />"
)

type (
	Validator[T any] interface {
		Validate(exCheckers ...func(original *T) (errors []error)) Validator[T]
		Invalid() bool
		GetData() T
		GetErrors() []error
		GetError() error
	}
	ValidatorImpl[T any] struct {
		original T
		errors   []error
		checkers []Checker
	}

	DefaultDataBindFn[T any] func() (form T, err error)
)

func WithData[T any](data T, checkers ...Checker) Validator[T] {
	return &ValidatorImpl[T]{original: data, checkers: checkers}
}

func WithDataBind[T any](fn DefaultDataBindFn[T], checkers ...Checker) Validator[T] {
	if fn == nil {
		return &ValidatorImpl[T]{errors: []error{_errors.New("默认读取数据方法为空")}}
	}

	data, err := fn()
	if err != nil {
		return &ValidatorImpl[T]{errors: []error{fmt.Errorf("读取数据失败：%w", err)}}
	}

	return &ValidatorImpl[T]{original: data, checkers: checkers}
}

func SetDefaultSliceSplitChar(char string) { defaultSliceSplitChar = char }

func SetDefaultErrorSplitChar(char string) { defaultErrorSplitChar = char }

func (my *ValidatorImpl[T]) Validate(exCheckers ...func(original *T) (errors []error)) Validator[T] {
	v := reflect.ValueOf(any(my.original)) // 显式转 any，避免对 T 直接反射的坑

	// 解引用指针，nil 指针直接返回，避免后续反射 panic
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return my
		}
		v = v.Elem()
	}

	// 非结构体（基础类型、空接口等）无需遍历字段
	if !v.IsValid() || v.Kind() != reflect.Struct {
		return my
	}

	my.validateRecursive(v, my.checkers, &my.errors)

	if len(exCheckers) > 0 {
		for _, checker := range exCheckers {
			if errs := checker(&my.original); len(errs) > 0 {
				my.errors = append(my.errors, errs...)
			}
		}
	}

	return my
}

// validateRecursive recursively validates a struct and its nested structs.
// It supports dotted field names like "Address.City.Name" for nested validation.
func (my *ValidatorImpl[T]) validateRecursive(v reflect.Value, checkers []Checker, errors *[]error) {
	if !v.IsValid() || v.Kind() != reflect.Struct {
		return
	}

	t := v.Type()

	// Validate top-level fields first
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		val := v.Field(i)

		// Check if there's an exact-match checker for this field
		for _, checker := range checkers {
			if checker.GetField() == f.Name {
				if errs := checker.check(val).GetErrors(); len(errs) > 0 {
					*errors = append(*errors, errs...)
				}
			}
		}

		// If the field value is a struct or pointer to struct, recurse into it
		var fieldVal reflect.Value
		switch val.Kind() {
		case reflect.Pointer:
			if val.IsNil() {
				continue
			}
			fieldVal = val.Elem()
		default:
			fieldVal = val
		}

		if fieldVal.Kind() == reflect.Struct {
			// Find nested checkers: check if there are checkers whose field name matches
			// a pattern like "NestedField.SubField"
			nestedCheckers := make([]Checker, 0)
			for _, checker := range checkers {
				fieldName := checker.GetField()
				if strings.HasPrefix(fieldName, f.Name+".") {
					// This checker is for a field inside the current nested struct
					// Create a new checker with adjusted field name
					adjustedFieldName := strings.TrimPrefix(fieldName, f.Name+".")
					ci := checker.(*CheckerImpl)
					adjusted := NewChecker(adjustedFieldName, ci.name)
					// Copy other fields from original checker
					if ci.required {
						adjusted.Required()
					}
					if ci.min != nil {
						adjusted.Min(*ci.min)
					}
					if ci.max != nil {
						adjusted.Max(*ci.max)
					}
					if ci.size != nil {
						adjusted.Size(*ci.size)
					}
					if ci.in != nil {
						adjusted.In(*ci.in)
					}
					if ci.regex != nil {
						adjusted.Regex(*ci.regex)
					}
					if ci.format != nil {
						adjusted.Format(*ci.format)
					}
					if ci.boolean != nil {
						if *ci.boolean {
							adjusted.True()
						} else {
							adjusted.False()
						}
					}
					if ci.errMsg != nil {
						adjusted.ErrMsg(*ci.errMsg)
					}
					nestedCheckers = append(nestedCheckers, adjusted)
				}
			}

			if len(nestedCheckers) > 0 {
				my.validateRecursive(fieldVal, nestedCheckers, errors)
			}
		}
	}
}

func (my *ValidatorImpl[T]) Invalid() bool { return len(my.errors) > 0 }

func (my *ValidatorImpl[T]) GetData() T { return my.original }

func (my *ValidatorImpl[T]) GetErrors() []error { return my.errors }

func (my *ValidatorImpl[T]) GetError() error {
	return _errors.New(anySlices.FillFunc(my.errors, func(_ int, err error) string { return err.Error() }).JoinNotEmpty(defaultErrorSplitChar))
}
