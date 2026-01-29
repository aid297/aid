package operationV2

type MultivariateAttr[T any] struct {
	Items   []T
	HitFunc func(idx int, item T)
}

func NewMultivariateAttr[T any](items ...T) *MultivariateAttr[T] {
	return &MultivariateAttr[T]{Items: items}
}
func (my *MultivariateAttr[T]) SetItems(items ...T) *MultivariateAttr[T] { my.Items = items; return my }
func (my *MultivariateAttr[T]) SetHitFunc(fn func(idx int, item T)) *MultivariateAttr[T] {
	my.HitFunc = fn
	return my
}
