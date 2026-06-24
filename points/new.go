package points

func New[T any](val T) *T { return &val }

func Value[T any](val *T) T {
	var empty T
	if val == nil {
		return empty
	}

	return *val
}
