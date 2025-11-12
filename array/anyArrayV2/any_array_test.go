package anyArrayV2

import (
	"log"
	"testing"
)

func Test1(t *testing.T) {
	a := New(Items(1, 2, 3, 4, 5))
	if a.Empty() {
		t.Error("Expected array to be not empty")
	}
	if !a.NotEmpty() {
		t.Error("Expected array to be not empty")
	}
	if a.Length() != 5 {
		t.Errorf("Expected length to be 5, got %d", a.Length())
	}
	a = a.Append(6, 7)
	if a.Length() != 7 {
		t.Errorf("Expected length to be 7 after append, got %d", a.Length())
	}
	if val := a.GetValue(0); val != 1 {
		t.Errorf("Expected first element to be 1, got %d", val)
	}
	a = a.RemoveByIndex(0)
	if a.Length() != 6 {
		t.Errorf("Expected length to be 7 after remove, got %d", a.Length())
	}
	if val := a.GetValue(0); val != 2 {
		t.Errorf("Expected first element to be 1 after remove, got %d", val)
	}
	a = a.Clean()
	if !a.Empty() {
		t.Error("Expected array to be empty after clear")
	}
}

func Test2(t *testing.T) {
	a := New(Items("a", "a", "b", "b", "c", "c"))
	if a.Unique().ToString() != "[a b c]" {
		t.Errorf("错误1")
	}

	log.Println(a.IntersectionBySlice("a", "b").ToString())

	log.Println(a.DifferenceBySlice("a", "b").ToString())

	b := New(Len[string](1))
	log.Println(b.ToString())
}
