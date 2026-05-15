package operation

// GetIntersection 获取两个数组的交集
func GetIntersection[T comparable](setA, setB []T) (intersection []T) {
	elementMap := make(map[T]bool)

	// 将 setB 的元素存入 map
	for _, b := range setB {
		elementMap[b] = true
	}

	// 遍历 setA，检查是否在 setB 中
	for _, a := range setA {
		if elementMap[a] {
			intersection = append(intersection, a)
		}
	}

	return
}

// GetUnion 获取两个数组的并集
func GetUnion[T comparable](setA, setB []T) (union []T) {
	elementMap := make(map[T]bool)

	// 将 setA 的元素存入 map
	for _, a := range setA {
		if !elementMap[a] {
			union = append(union, a)
			elementMap[a] = true
		}
	}

	// 将 setB 的元素存入 map
	for _, b := range setB {
		if !elementMap[b] {
			union = append(union, b)
			elementMap[b] = true
		}
	}

	return
}

// GetDifference 获取两个数组的差集
func GetDifference[T comparable](setA, setB []T) (difference []T) {
	elementMap := make(map[T]bool)

	// 将 setB 的元素存入 map
	for _, b := range setB {
		elementMap[b] = true
	}

	// 遍历 setA，检查是否不在 setB 中
	for _, a := range setA {
		if !elementMap[a] {
			difference = append(difference, a)
		}
	}

	return
}
