package operationV2

import "testing"

func TestTernary1(t *testing.T) {
	t1 := NewTernary(TrueValue("真"), FalseValue("假"))
	a1 := t1.GetByValue(true)
	if a1 != "真" {
		t.Errorf("错误1")
	}

	a2 := NewTernary(TrueFn(func() string { return "真" })).GetByFunc(func() bool { return false })
	if a2 != "" {
		t.Errorf("错误2")
	}
}
