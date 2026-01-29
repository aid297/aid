package operationV2

type MultivariateAttr[T any] struct {
	Item    T
	HitFunc func(idx int, item T)
}

func NewMultivariateAttr[T any](item T) *MultivariateAttr[T] {
	return &MultivariateAttr[T]{Item: item}
}
func (my *MultivariateAttr[T]) SetItems(item T) *MultivariateAttr[T] { my.Item = item; return my }
func (my *MultivariateAttr[T]) SetHitFunc(fn func(idx int, item T)) *MultivariateAttr[T] {
	my.HitFunc = fn
	return my
}
