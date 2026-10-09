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
