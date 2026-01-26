package anyArrayV3

import (
	"testing"
)

func Test1(t *testing.T) {
	var a AnyArrayer[[]int] = NewItems[[]int]()
	t.Log(a.ToSlice())
}
