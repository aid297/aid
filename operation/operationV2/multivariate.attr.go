package operationV2

type (
	MultivariateAttributer[T any] interface{ Register(ma *MultivariateAttr[T]) }

	MultivariateAttr[T any] struct {
		Items   []T
		Default T
		HitFunc func(idx int, item T)
	}

	AttrItems[T any]   struct{ items []T }
	AttrHitFunc[T any] struct{ hitFunc func(idx int, item T) }
)

func NewMultivariateAttr[T any](attrs ...MultivariateAttributer[T]) *MultivariateAttr[T] {
	return (&MultivariateAttr[T]{}).SetAttrs(attrs...)
}

func (my *MultivariateAttr[T]) SetAttrs(attrs ...MultivariateAttributer[T]) *MultivariateAttr[T] {
	for _, attr := range attrs {
		attr.Register(my)
	}

	return my
}

func Items[T any](items ...T) MultivariateAttributer[T]   { return &AttrItems[T]{items: items} }
func (my *AttrItems[T]) Register(ma *MultivariateAttr[T]) { ma.Items = my.items }

func HitFunc[T any](hitFunc func(idx int, item T)) MultivariateAttributer[T] {
	return &AttrHitFunc[T]{hitFunc}
}
func (my *AttrHitFunc[T]) Register(ma *MultivariateAttr[T]) { ma.HitFunc = my.hitFunc }
