package validators

import (
	_errors "errors"
	"fmt"
	"reflect"

	"github.com/aid297/aid/v2/anySlices"
)

type (
	Validator[T any] interface {
		Validate(exCheckers ...func(original T) (errors []error)) Validator[T]
		Invalid() bool
		GetErrors() []error
	}
	ValidatorImpl[T any] struct {
		original T
		errors   []error
		checkers []Checker
	}

	DefaultDataBindFn[T any] func() (T, error)
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

func (my *ValidatorImpl[T]) Validate(exCheckers ...func(original T) (errors []error)) Validator[T] {
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

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		val := v.Field(i)

		for _, checker := range my.checkers {
			if checker.GetField() == f.Name {
				if errs := checker.Check(val).GetErrors(); len(errs) > 0 {
					my.errors = append(my.errors, errs...)
				}
			}
		}

		if len(exCheckers) > 0 {
			for _, checker := range exCheckers {
				if errs := checker(my.original); len(errs) > 0 {
					my.errors = append(my.errors, errs...)
				}
			}
		}
	}

	return my
}

func (my *ValidatorImpl[T]) Invalid() bool { return len(my.errors) > 0 }

func (my *ValidatorImpl[T]) GetData() T { return my.original }

func (my *ValidatorImpl[T]) GetErrors() []error { return my.errors }

func (my *ValidatorImpl[T]) GetError() error {
	return _errors.New(anySlices.FillFunc(my.errors, func(_ int, err error) string { return err.Error() }).JoinNotEmpty("<br />"))
}
