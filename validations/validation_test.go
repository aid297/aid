package validations_test

import (
	"strings"
	"testing"

	"github.com/aid297/aid/v2/validations"
)

func TestCustomType(t *testing.T) {
	type A string
	var (
		Custom1 A = "CUSTOM-A"
		Custom2 A = "CUSTOM-B"
	)
	type S struct {
		AField A `json:"a_field" v-rule:"(required)(min>0)(in==CUSTOM-A,CUSTOM-B)" v-name:"a_field"`
	}

	_ = Custom2
	var s = S{AField: Custom1}

	checker := validations.Once().Checker(s).Validate()
	if checker.Invalid() {
		t.Errorf("验证不通过：%v", checker.Error())
	}

	t.Logf("OK")
}

// TestNestedPathPreserved 验证 3 层及以上嵌套时，错误信息保留完整点分路径。
// 修复前：递归只会传递直接父级的 v-name，导致最内层字段名丢失外层上下文。
// 修复后：递归携带祖先链，结构体字段与切片元素内部字段均能拼接出完整路径。
func TestNestedPathPreserved(t *testing.T) {
	type L3 struct {
		X string `v-rule:"(required)" v-name:"x"`
	}
	type L2 struct {
		Items []L3 `v-rule:"(required)(min>=1)" v-name:"items"`
	}
	type L1 struct {
		Sub L2 `v-rule:"(required)" v-name:"sub"`
	}

	// 情况 1：非空切片 + 内层字段为空，应得到完整路径 "sub.items.x"
	req1 := &L1{Sub: L2{Items: []L3{{X: ""}}}}
	joined1 := validations.Once().Checker(req1).Validate().ErrorToString("\n")
	if !strings.Contains(joined1, "sub.items.x") {
		t.Errorf("期望错误信息包含完整路径 'sub.items.x'，实际：%s", joined1)
	}

	// 情况 2：空切片，递归零值仍应保留完整路径
	req2 := &L1{Sub: L2{Items: []L3{}}}
	joined2 := validations.Once().Checker(req2).Validate().ErrorToString("\n")
	if !strings.Contains(joined2, "sub.items.x") {
		t.Errorf("空切片递归期望错误信息包含 'sub.items.x'，实际：%s", joined2)
	}

	// 情况 3：3 层嵌套 struct（非切片）也应保留完整路径
	type Inner struct {
		Name string `v-rule:"(required)" v-name:"name"`
	}
	type Middle struct {
		Data Inner `v-rule:"(required)" v-name:"data"`
	}
	type Root struct {
		Sub Middle `v-rule:"(required)" v-name:"sub"`
	}
	req3 := &Root{Sub: Middle{Data: Inner{Name: ""}}}
	joined3 := validations.Once().Checker(req3).Validate().ErrorToString("\n")
	if !strings.Contains(joined3, "sub.data.name") {
		t.Errorf("嵌套 struct 期望错误信息包含 'sub.data.name'，实际：%s", joined3)
	}

	t.Logf("OK: %s", joined1)
}
