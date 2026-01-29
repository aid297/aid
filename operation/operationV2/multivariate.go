package operationV2

import "github.com/aid297/aid/array/anySlice"

type (
	Multivatiater[T any] interface {
		SetItems(priority uint, items ...T) *Multivariate[T]
		Finally(fn func(item T) bool) T
		SetDefault(item T) *Multivariate[T]
	}

	Multivariate[T any] struct {
		Items   anySlice.AnySlicer[[]T]
		Default T
	}
)

// NewMultivariate 实例化：多元运算
func NewMultivariate[T any](cap int) *Multivariate[T] {
	var def T
	return &Multivariate[T]{Items: anySlice.New(anySlice.Len[[]T](cap)), Default: def}
}

// SetItems 设置优先级项
func (my *Multivariate[T]) SetItems(priority uint, items ...T) *Multivariate[T] {
	if my.Items.Has(int(priority)) {
		my.Items.SetValue(int(priority), items)
	}

	return my
}

// SetDefault 设置默认值
func (my *Multivariate[T]) SetDefault(item T) *Multivariate[T] { my.Default = item; return my }

// FinllayFunc 获取优先级选项
func (my *Multivariate[T]) FinallyFunc(condition func(item T) bool) (int, T) {
	for _, items := range my.Items.ToSlice() {
		for idx := range items {
			a := condition(items[idx])
			if a {
				return idx, items[idx]
			}
		}
	}

	return -1, my.Default
}
