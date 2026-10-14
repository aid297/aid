package validations_test

import (
	"testing"

	"github.com/aid297/aid/v2/validations"
)

type TrimRequest struct {
	All    string  `v-rule:"(tb)" v-name:"全裁"`
	Left   string  `v-rule:"(tl)" v-name:"左裁"`
	Right  *string `v-rule:"(tr)" v-name:"右裁"`
	Space  string  `v-rule:"(tb{ })" v-name:"仅半角空格"`
	Custom string  `v-rule:"(tb{ |　|\t|\n|})" v-name:"自定义字符集"`
	Order  string  `v-rule:"(tb)(min>=3)" v-name:"先裁后校"`
}

func TestTrim(t *testing.T) {
	right := "x \t\n　"
	form := &TrimRequest{
		All:    "  a b 　 \t\n",
		Left:   " \t\n　a",
		Right:  &right,
		Space:  "  a b　",
		Custom: " \t\n　a b",
		Order:  " ab ",
	}

	checker := validations.Once().Checker(form).Validate()
	if !checker.Invalid() {
		t.Fatalf("期望验证不通过：Order 裁剪后为 2 个字符，不满足 min>=3")
	}

	expectAll := "a b"
	if form.All != expectAll {
		t.Errorf("All = %q, 期望 %q", form.All, expectAll)
	}
	expectLeft := "a"
	if form.Left != expectLeft {
		t.Errorf("Left = %q, 期望 %q", form.Left, expectLeft)
	}
	if form.Right == nil {
		t.Fatal("Right 不应为 nil")
	}
	expectRight := "x"
	if *form.Right != expectRight {
		t.Errorf("Right = %q, 期望 %q", *form.Right, expectRight)
	}
	expectSpace := "a b　"
	if form.Space != expectSpace {
		t.Errorf("Space = %q, 期望 %q", form.Space, expectSpace)
	}
	expectCustom := "a b"
	if form.Custom != expectCustom {
		t.Errorf("Custom = %q, 期望 %q", form.Custom, expectCustom)
	}
	expectOrder := "ab"
	if form.Order != expectOrder {
		t.Errorf("Order = %q, 期望 %q", form.Order, expectOrder)
	}
}

func TestTrimNoChange(t *testing.T) {
	form := &TrimRequest{
		All:   "a",
		Left:  "a",
		Right: nil,
		Order: " abc ",
	}

	checker := validations.Once().Checker(form).Validate()
	for _, wrong := range checker.Errors() {
		t.Logf("%v", wrong)
	}

	if form.All != "a" {
		t.Errorf("All = %q, 期望 %q", form.All, "a")
	}
}

func TestTrimRequired(t *testing.T) {
	type T struct {
		A string  `v-rule:"(required)(tb)" v-name:"A"`
		B *string `v-rule:"(required)(tb)" v-name:"B"`
	}

	right := " 　\t\n"
	form := &T{A: " 　\t\n", B: &right}

	checker := validations.Once().Checker(form).Validate()
	if !checker.Invalid() {
		t.Fatalf("期望验证不通过：A、B 裁剪后为空串，required 应报错")
	}
	for _, wrong := range checker.Errors() {
		t.Logf("%v", wrong)
	}
}

func TestTrimDefault(t *testing.T) {
	type T struct {
		A string  `v-rule:"(min>=1)" v-name:"A"`
		B *string `v-rule:"(min>=1)" v-name:"B"`
	}

	b := " 　\t\nb 　\t\n"
	form := &T{A: "  a b　", B: &b}

	checker := validations.Once().Checker(form).Validate()
	if checker.Invalid() {
		t.Fatalf("期望验证通过：%v", checker.Errors())
	}

	if expect := "a b"; form.A != expect {
		t.Errorf("A = %q, 期望 %q（未声明裁剪规则时应默认按 (tb) 清理）", form.A, expect)
	}
	if form.B == nil {
		t.Fatal("B 不应为 nil")
	}
	if expect := "b"; *form.B != expect {
		t.Errorf("B = %q, 期望 %q（*string 未声明裁剪规则时应默认按 (tb) 清理）", *form.B, expect)
	}
}

func TestTrimNoTrim(t *testing.T) {
	type T struct {
		Keep      string  `v-rule:"(nt)" v-name:"不清理"`
		KeepPtr   *string `v-rule:"(nt)" v-name:"指针不清理"`
		KeepEmpty string  `v-rule:"(required)(nt)" v-name:"必填不清理"`
		Explicit  string  `v-rule:"(nt)(tl)" v-name:"显式规则优先"`
	}

	keepPtr := " 　b "
	form := &T{
		Keep:      "  a 　",
		KeepPtr:   &keepPtr,
		KeepEmpty: " 　",
		Explicit:  "  a  ",
	}

	checker := validations.Once().Checker(form).Validate()
	if checker.Invalid() {
		t.Fatalf("期望验证通过：KeepEmpty 带 (nt) 不清理，required 不应报错：%v", checker.Errors())
	}

	if expect := "  a 　"; form.Keep != expect {
		t.Errorf("Keep = %q, 期望 %q（(nt) 不应清理）", form.Keep, expect)
	}
	if expect := " 　b "; *form.KeepPtr != expect {
		t.Errorf("KeepPtr = %q, 期望 %q（(nt) 不应清理）", *form.KeepPtr, expect)
	}
	if expect := "a  "; form.Explicit != expect {
		t.Errorf("Explicit = %q, 期望 %q（(nt) 与显式 (tl) 同时存在时显式规则优先）", form.Explicit, expect)
	}
}
