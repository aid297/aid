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
	ExCheckFn[T any]         = func(original *T) (errors []error)
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

func (my *ValidatorImpl[T]) Validate(exCheckers ...ExCheckFn[T]) Validator[T] {
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
// It supports dotted field names like "Address.City.Name" for nested validation,
// and also recurses into slices of structs (e.g. "Items.Name" validates each element's Name).
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

		// Dereference pointer fields
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

		switch fieldVal.Kind() {
		case reflect.Struct:
			// Nested struct: find dotted-notation checkers and recurse
			nestedCheckers := findNestedCheckers(checkers, f.Name)
			if len(nestedCheckers) > 0 {
				my.validateRecursive(fieldVal, nestedCheckers, errors)
			}
		case reflect.Slice, reflect.Array:
			// Slice/array: check if element type is struct, then recurse into each element
			elemType := fieldVal.Type().Elem()
			if elemType.Kind() == reflect.Pointer {
				elemType = elemType.Elem()
			}
			if elemType.Kind() == reflect.Struct {
				nestedCheckers := findNestedCheckers(checkers, f.Name)
				if len(nestedCheckers) > 0 {
					for j := 0; j < fieldVal.Len(); j++ {
						elemVal := fieldVal.Index(j)
						if elemVal.Kind() == reflect.Pointer {
							if elemVal.IsNil() {
								continue
							}
							elemVal = elemVal.Elem()
						}
						my.validateRecursive(elemVal, nestedCheckers, errors)
					}
				}
			}
		}
	}
}

// findNestedCheckers finds checkers whose field name starts with the given prefix
// followed by a dot (e.g. prefix="Items" matches "Items.Name"), and creates
// adjusted checkers with the prefix stripped ("Name").
func findNestedCheckers(checkers []Checker, prefix string) []Checker {
	nestedCheckers := make([]Checker, 0)
	for _, checker := range checkers {
		fieldName := checker.GetField()
		if strings.HasPrefix(fieldName, prefix+".") {
			adjustedFieldName := strings.TrimPrefix(fieldName, prefix+".")
			ci := checker.(*CheckerImpl)
			adjusted := NewChecker(adjustedFieldName, ci.name)
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
	return nestedCheckers
}

func (my *ValidatorImpl[T]) Invalid() bool { return len(my.errors) > 0 }

func (my *ValidatorImpl[T]) GetData() T { return my.original }

func (my *ValidatorImpl[T]) GetErrors() []error { return my.errors }

func (my *ValidatorImpl[T]) GetError() error {
	return _errors.New(anySlices.FillFunc(my.errors, func(_ int, err error) string { return err.Error() }).JoinNotEmpty(defaultErrorSplitChar))
}
