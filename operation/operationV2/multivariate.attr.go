package operationV2

type MultivariateAttr[T any] struct {
	Item    T
	HitFunc func(idx int, item T)
}
