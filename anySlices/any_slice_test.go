package anySlices_test

import (
	"testing"

	"github.com/aid297/aid/v2/anySlices"
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

func Test4(t *testing.T) {
	// 子级 ("a","b") 中 "b" 超过父级 ["a"] 范围 → NotIn 应返回 true
	if !anySlices.NewItems("a").
		NotIn("a", "b") {
		t.Fatalf("子级角色菜单不能超过父级")
	}
}

func Test5(t *testing.T) {
	if anySlices.NewItems("a").
		In("a", "b") {
		t.Fatalf("子级角色菜单不能超过父级")
	}
}
