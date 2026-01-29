package operationV2

import "github.com/aid297/aid/array/anySlice"

type (
	Multivariater[T any] interface {
		SetItems(priority uint, items ...T) *Multivariate[T]
		Finally(fn func(item T) bool) (int, T)
		SetDefault(item T) *Multivariate[T]
	}

	Multivariate[T any] struct {
		Items   anySlice.AnySlicer[[]T]
		Default T
	}
)

// NewMultivariate 实例化：多元运算
func NewMultivariate[T any]() *Multivariate[T] {
	return &Multivariate[T]{Items: anySlice.New[[]T]()}
}

// Append 添加优先级项
func (my *Multivariate[T]) Append(items ...T) *Multivariate[T] { my.Items.Append(items); return my }

// SetItems 设置优先级项
func (my *Multivariate[T]) Set(priority uint, items ...T) *Multivariate[T] {
	if my.Items.Has(int(priority)) {
		my.Items.SetValue(int(priority), items)
	}

	return my
}

// SetDefault 设置默认值
func (my *Multivariate[T]) SetDefault(item T) *Multivariate[T] { my.Default = item; return my }

// Finally 获取优先级选项
func (my *Multivariate[T]) Finally(condition func(item T) bool) (int, T) {
	for _, items := range my.Items.ToSlice() {
		for idx, item := range items {
			if condition(item) {
				return idx, item
			}
		}
	}

	return -1, my.Default
}
