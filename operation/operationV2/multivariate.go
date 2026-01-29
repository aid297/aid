package operationV2

import "github.com/aid297/aid/array/anySlice"

type (
	Multivariater[T any] interface {
		Append(items *MultivariateAttr[T]) Multivariater[T]
		Finally(fn func(item T) bool) (int, T)
		SetDefault(item *MultivariateAttr[T]) Multivariater[T]
	}

	Multivariate[T any] struct {
		Items   anySlice.AnySlicer[*MultivariateAttr[T]]
		Default *MultivariateAttr[T]
	}
)

// NewMultivariate 实例化：多元运算
func NewMultivariate[T any]() *Multivariate[T] {
	return &Multivariate[T]{Items: anySlice.New[*MultivariateAttr[T]]()}
}

// Append 添加优先级项
func (my *Multivariate[T]) Append(item *MultivariateAttr[T]) *Multivariate[T] {
	my.Items.Append(item)
	return my
}

// SetDefault 设置默认值
func (my *Multivariate[T]) SetDefault(item *MultivariateAttr[T]) *Multivariate[T] {
	my.Default = item
	return my
}

// Finally 获取优先级选项
func (my *Multivariate[T]) Finally(condition func(item T) bool) (int, T) {
	for idx, item := range my.Items.ToSlice() {
		if condition(item.Item) {
			if item.HitFunc != nil {
				item.HitFunc(idx, item.Item)
			}
			return idx, item.Item
		}
	}

	return -1, my.Default.Item
}
