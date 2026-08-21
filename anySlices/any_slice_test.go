package anySlices_test

import (
	"testing"

	"github.com/aid297/aid/v3/anySlices"
)

func Test1(t *testing.T) {
	var a anySlices.AnySlicer[int] = anySlices.New(anySlices.List([]int{1, 2, 3, 4, 5}))
	t.Log(a.ToSlice())
}

func Test2(t *testing.T) {
	var a anySlices.AnySlicer[string] = anySlices.New(anySlices.Cap[string](5))
	t.Logf("%#v\n", a.GetValueOrDefault(0, "default"))
}

func Test3(t *testing.T) {
	var a anySlices.AnySlicer[int] = anySlices.New(anySlices.List([]int{1, 2, 3, 4, 5}))
	t.Logf("%#v\n", a.RemoveByIndex(0, 1, 2).ToSlice())
}
