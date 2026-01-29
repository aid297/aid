package operationV2

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	m := NewMultivariate[string]().
		Append(NewMultivariateAttr("a").SetHitFunc(func(_ int, _ string) { fmt.Printf("采用高级") })).       // A 最高优先级：终端命令
		Append(NewMultivariateAttr("b").SetHitFunc(func(idx int, item string) { fmt.Printf("采用次高级") })). // B 次高优先级：全局变量
		SetDefault(NewMultivariateAttr("c"))                                                             // 设置默认值

	_, f := m.Finally(func(item string) bool { return item != "" })

	if f != "a" {
		t.Fatalf("错误：%s", f)
	}
	t.Logf("成功：%s", f)

	// Test default value
	_, f = m.Finally(func(item string) bool { return item == "missing" })
	if f != "d" {
		t.Fatalf("Default value error: %s", f)
	}
	t.Logf("Default value success: %s", f)
}
