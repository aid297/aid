package validatorV2

import (
	"fmt"
	"testing"
	"time"
)

// TestValidator_Basic 覆盖字符串、数字、数组、嵌套和时间规则
func TestValidator_Basic(t *testing.T) {
	type Inner struct {
		Code string `v-rule:"(required)(len:3)" v-name:"编码"`
	}

	now := time.Now()
	// 构造测试结构体
	type Req struct {
		Name     string      `v-rule:"(required)(min>3)(max<10)" v-name:"名称" v-ex:"onlyEnglish"`
		Email    string      `v-rule:"(required)(email)" v-name:"邮箱"`
		Tags     []string    `v-rule:"(min:1)(max:3)" v-name:"标签"`
		Score    int         `v-rule:"(min:1)(max:100)" v-name:"分数"`
		InnerVal Inner       `v-name:"内嵌"`
		Dates    []time.Time `v-rule:"(min:1)" v-name:"日期列表"`
	}

	tests := []struct {
		name string
		in   Req
		want int // 期望错误字段数量
	}{
		{"valid", Req{
			Name:     "abcd",
			Email:    "x@y.com",
			Tags:     []string{"a"},
			Score:    50,
			InnerVal: Inner{Code: "ABC"},
			Dates:    []time.Time{now},
		}, 0},
		{"missing required and bad email", Req{
			Name:     "ab1",
			Email:    "not-an-email",
			Tags:     []string{},
			Score:    0,
			InnerVal: Inner{Code: "AB"},
			Dates:    []time.Time{},
		}, 7},
	}

	// 注册扩展函数 onlyEnglish: 仅允许 A-Z a-z
	RegisterExFun("onlyEnglish", func(val any) error {
		var s string
		switch v := val.(type) {
		case string:
			s = v
		case *string:
			if v == nil {
				return nil
			}
			s = *v
		default:
			return nil
		}
		for _, r := range s {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				continue
			}
			return fmt.Errorf("值只能包含英文字母")
		}
		return nil
	})

	v := Validator{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := v.Validate(tt.in)
			if len(res) != tt.want {
				t.Fatalf("got %d field errors, want %d, details: %+v", len(res), tt.want, res)
			}
		})
	}
}

func TestRemoveRule(t *testing.T) {
	cases := []struct {
		in   string
		key  string
		want string
	}{
		{"(required)(min:3)", "required", "(min:3)"},
		{"(min:3)(required)(max:5)", "required", "(min:3)(max:5)"},
		{"(required)", "required", ""},
		{"", "required", ""},
	}
	for _, c := range cases {
		got := removeRule(c.in, c.key)
		if got != c.want {
			t.Fatalf("removeRule(%q,%q) = %q; want %q", c.in, c.key, got, c.want)
		}
	}
}

func TestRequiredPointerSemantics(t *testing.T) {
	type Req struct {
		P *string `v-rule:"required" v-name:"p"`
	}
	v := Validator{}
	str := ""
	// nil pointer -> error
	res := v.Validate(Req{P: nil})
	if len(res) != 1 {
		t.Fatalf("expected 1 error for nil pointer, got %d: %+v", len(res), res)
	}

	// pointer to empty string -> should pass (presence semantics)
	res = v.Validate(Req{P: &str})
	if len(res) != 0 {
		t.Fatalf("expected 0 errors for pointer to empty string, got %d: %+v", len(res), res)
	}
}

func TestVExPointerInvocation(t *testing.T) {
	defer UnregisterExFun("checkNil")
	// register a v-ex that errors when val is nil
	RegisterExFun("checkNil", func(val any) error {
		if val == nil {
			return fmt.Errorf("val is nil")
		}
		return nil
	})

	type Req struct {
		P *string `v-ex:"checkNil" v-name:"p"`
	}
	v := Validator{}

	// nil pointer -> expect error from v-ex
	res := v.Validate(Req{P: nil})
	if len(res) == 0 {
		t.Fatalf("expected v-ex error for nil pointer, got none")
	}

	// non-nil pointer -> no error from v-ex
	s := ""
	res = v.Validate(Req{P: &s})
	if len(res) != 0 {
		t.Fatalf("expected no v-ex error for non-nil pointer, got: %+v", res)
	}
}

// TestValidator_Pointers 覆盖指针字段的校验行为（nil 与 非 nil）
func TestValidator_Pointers(t *testing.T) {
	type Inner struct {
		Code string `v-rule:"(required)(len:3)" v-name:"编码"`
	}

	type ReqPtr struct {
		PName  *string `v-rule:"(required)(max<10)" v-name:"名称" v-ex:"onlyEnglish"`
		PInner *Inner  `v-rule:"(required)" v-name:"内嵌指针"`
	}

	strPtr := func(s string) *string { return &s }

	tests := []struct {
		name string
		in   ReqPtr
		want int
	}{
		{"ptr valid", ReqPtr{PName: strPtr("abc"), PInner: &Inner{Code: "ABC"}}, 0},
		{"ptr empty string passes", ReqPtr{PName: strPtr(""), PInner: &Inner{Code: "ABC"}}, 0},
		{"ptr nil values", ReqPtr{PName: nil, PInner: nil}, 2},
		{"ptr invalid contents", ReqPtr{PName: strPtr("ab1"), PInner: &Inner{Code: "AB"}}, 2},
	}

	// 确保扩展函数已注册（重复注册是安全的）
	RegisterExFun("onlyEnglish", func(val any) error {
		var s string
		switch v := val.(type) {
		case string:
			s = v
		case *string:
			if v == nil {
				return nil
			}
			s = *v
		default:
			return nil
		}
		for _, r := range s {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				continue
			}
			return fmt.Errorf("值只能包含英文字母")
		}
		return nil
	})

	v := Validator{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := v.Validate(tt.in)
			if len(res) != tt.want {
				t.Fatalf("got %d field errors, want %d, details: %+v", len(res), tt.want, res)
			}
		})
	}
}

func TestInNotInAndRegex(t *testing.T) {
	type Req struct {
		Choice string `v-rule:"in=a,b,c" v-name:"choice"`
		Block  string `v-rule:"not-in=x,y" v-name:"block"`
		Phone  string `v-rule:"regex:^\\d{3}-\\d{4}$" v-name:"phone"`
	}

	v := Validator{}

	// valid case
	res := v.Validate(Req{Choice: "a", Block: "z", Phone: "123-4567"})
	if len(res) != 0 {
		t.Fatalf("expected 0 errors for valid case, got: %+v", res)
	}

	// invalid cases
	res = v.Validate(Req{Choice: "d", Block: "x", Phone: "12-34567"})
	if len(res) != 3 {
		t.Fatalf("expected 3 errors for invalid case, got %d: %+v", len(res), res)
	}
}

func TestSliceElementVEx(t *testing.T) {
	// register a v-ex for element validation
	defer UnregisterExFun("elemEnglish")
	RegisterExFun("elemEnglish", func(val any) error {
		// val may be *Element or Element or pointer to string depending on call
		switch v := val.(type) {
		case *string:
			if v == nil {
				return fmt.Errorf("nil")
			}
			for _, r := range *v {
				if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
					return fmt.Errorf("elem only letters")
				}
			}
		case string:
			for _, r := range v {
				if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
					return fmt.Errorf("elem only letters")
				}
			}
		default:
			// ignore other types
		}
		return nil
	})

	type Elem struct {
		Name string `v-ex:"elemEnglish" v-name:"name"`
	}
	type Req struct {
		Items []Elem `v-name:"items"`
	}

	v := Validator{}
	// valid
	res := v.Validate(Req{Items: []Elem{{Name: "Bob"}, {Name: "Alice"}}})
	if len(res) != 0 {
		t.Fatalf("expected 0 errors for valid items, got: %+v", res)
	}

	// invalid (one item contains digit)
	res = v.Validate(Req{Items: []Elem{{Name: "Bob1"}, {Name: "Alice"}}})
	if len(res) == 0 {
		t.Fatalf("expected errors for invalid item, got none")
	}
}
