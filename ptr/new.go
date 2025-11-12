package ptr

func New[T any](val T) *T { return &val }
