package points

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	a := New(123)
	t.Logf("%#v", a)
}

func Test2(t *testing.T) {
	var a *int

	b := Value(a)

	fmt.Printf("res: %v\n\n", b)
}

func Test3(t *testing.T) {
	var a *int = New(4)

	b := Default(a, 3)

	fmt.Printf("res: %v\n\n", b)
}

func Test4(t *testing.T) {
	var a *int

	b := DefaultNil(a, 3)

	fmt.Printf("res: %v\n\n", b)
}
