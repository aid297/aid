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

func Default[T any](val *T, def T) T {
	if val == nil {
		return def
	}
	if reflect.ValueOf(Value(val)).IsZero() {
		return def
	}

	return Value(val)
}

func DefaultNil[T any](val *T, def T) T {
	if val == nil {
		return def
	}

	return *val
}
