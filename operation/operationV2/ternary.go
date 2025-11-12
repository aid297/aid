package operationV2

type (
	Ternary[T any] struct {
		trueFn  func() T
		falseFn func() T
	}
)

func NewTernary[T any](attrs ...TernaryAttributer[T]) Ternary[T] { return Ternary[T]{}.Set(attrs...) }

func (t Ternary[T]) Set(attrs ...TernaryAttributer[T]) Ternary[T] {
	if len(attrs) > 0 {
		for idx := range attrs {
			attrs[idx].Register(&t)
		}
	}
	return t
}

// DoByValue 执行回调 → 通过值
func (t Ternary[T]) DoByValue(condition bool) {
	if condition {
		if t.trueFn != nil {
			t.trueFn()
		}
	} else {
		if t.falseFn != nil {
			t.falseFn()
		}
	}
}

// DoByFunc 执行回调 → 通过函数
func (t Ternary[T]) DoByFunc(condition func() bool) {
	if condition() {
		if t.trueFn != nil {
			t.trueFn()
		}
	} else {
		if t.falseFn != nil {
			t.falseFn()
		}
	}
}

// GetByValue 获取值 → 通过值
func (t Ternary[T]) GetByValue(condition bool) T {
	var empty T
	if condition {
		if t.trueFn != nil {
			return t.trueFn()
		} else {
			return empty
		}
	} else {
		if t.falseFn != nil {
			return t.falseFn()
		} else {
			return empty
		}
	}
}

// GetByFunc 获取值 → 通过函数
func (t Ternary[T]) GetByFunc(condition func() bool) T {
	var empty T
	if condition() {
		if t.trueFn != nil {
			return t.trueFn()
		} else {
			return empty
		}
	} else {
		if t.falseFn != nil {
			return t.falseFn()
		} else {
			return empty
		}
	}
}

// OrError 三元运算符 → 处理错误
func OrError(target bool, trueValue, falseValue error) error {
	return NewTernary(TrueValue(trueValue), FalseValue(falseValue)).GetByValue(target)
}
