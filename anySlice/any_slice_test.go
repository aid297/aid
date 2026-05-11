package anySlice

import (
	"testing"
)

func Test1(t *testing.T) {
	var a AnySlicer[int] = New(List([]int{1, 2, 3, 4, 5}))
	t.Log(a.ToSlice())
}

func Test2(t *testing.T) {
	var a AnySlicer[string] = New(Cap[string](5))
	t.Logf("%#v\n", a.GetValueOrDefault(0, "default"))
}

func Test3(t *testing.T) {
	var a AnySlicer[int] = New(List([]int{1, 2, 3, 4, 5}))
	t.Logf("%#v\n", a.RemoveByIndex(0, 1, 2).ToSlice())
}
