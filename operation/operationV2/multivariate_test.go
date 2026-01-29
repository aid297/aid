package operationV2

import (
	"testing"
)

func Test1(t *testing.T) {
	m := NewMultivariate[string](3).
		SetItems(0, "a").      // A 最高优先级：终端命令
		SetItems(1, "c", "b"). // B 次高优先级：全局变量
		SetDefault("d")        // 设置默认值

	_, f := m.FinallyFunc(func(item string) bool { return item != "" })

	if f != "d" {
		t.Fatalf("错误：%s", f)
	}
	t.Logf("成功：%s", f)
}
