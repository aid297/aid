package anySlice

import (
	"testing"
)

func Test1(t *testing.T) {
	var a AnySlicer[[]int] = NewItems[[]int]()
	t.Log(a.ToRaw())
}
