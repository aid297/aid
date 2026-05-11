package anySlice

type AnySlicerAttr[T any] func(as AnySlicer[T])

func List[T any](list []T) AnySlicerAttr[T] {
	return func(as AnySlicer[T]) { as.SetData(list) }
}

func Items[T any](items ...T) AnySlicerAttr[T] {
	return func(as AnySlicer[T]) { as.SetData(items) }
}

func Len[T any](length int) AnySlicerAttr[T] {
	return func(as AnySlicer[T]) { as.SetData(make([]T, length)) }
}

func Cap[T any](cap int) AnySlicerAttr[T] {
	return func(as AnySlicer[T]) { as.SetData(make([]T, 0, cap)) }
}

func Empty[T any]() AnySlicerAttr[T] {
	return func(as AnySlicer[T]) { as.SetData(make([]T, 0)) }
}
