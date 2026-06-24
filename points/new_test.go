package points

import (
	"testing"
)

func Test1(t *testing.T) {
	a := New(123)
	t.Logf("%#v", a)
}
