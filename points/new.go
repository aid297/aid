package points

import "reflect"

func New[T any](val T) *T { return &val }

func Value[T any](val *T) T {
	var empty T
	if val == nil {
		return empty
	}

	return *val
}

func DefaultNil[T any](val *T, def T) {
	if val == nil {
		val = &def
		return
	}
}

func Default[T any](val *T, def T) {
	v := reflect.ValueOf(val)
	if v.Elem().IsZero() {
		val = &def
		return
	}
}
