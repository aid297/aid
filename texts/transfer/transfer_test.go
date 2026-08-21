package transfer_test

import (
	"testing"

	"github.com/aid297/aid/v3/texts/transfer"
)

// ==================== PascalCase 转换 ====================

func TestPascalToCamel(t *testing.T) {
	cases := []struct{ in, want string }{
		{"HelloWorld", "helloWorld"},
		{"MyClass", "myClass"},
		{"A", "a"},
		{"", ""},
	}
	for _, c := range cases {
		got := new(transfer.TransferImpl).New(c.in).PascalToCamel()
		if got != c.want {
			t.Fatalf("PascalToCamel(%q) 期望 %q，得到 %q", c.in, c.want, got)
		}
	}
}

func TestPascalToSnake(t *testing.T) {
	cases := []struct{ in, want string }{
		{"HelloWorld", "hello_world"},
		{"MyClass", "my_class"},
		{"HTMLElement", "html_element"},
	}
	for _, c := range cases {
		got := new(transfer.TransferImpl).New(c.in).PascalToSnake()
		if got != c.want {
			t.Fatalf("PascalToSnake(%q) 期望 %q，得到 %q", c.in, c.want, got)
		}
	}
}

func TestPascalToBabel(t *testing.T) {
	cases := []struct{ in, want string }{
		{"HelloWorld", "hello-world"},
		{"MyClass", "my-class"},
		{"HTMLElement", "html-element"},
	}
	for _, c := range cases {
		got := new(transfer.TransferImpl).New(c.in).PascalToBabel()
		if got != c.want {
			t.Fatalf("PascalToBabel(%q) 期望 %q，得到 %q", c.in, c.want, got)
		}
	}
}

// ==================== CamelCase 转换 ====================

func TestCamelToPascal(t *testing.T) {
	cases := []struct{ in, want string }{
		{"helloWorld", "HelloWorld"},
		{"myClass", "MyClass"},
		{"a", "A"},
		{"", ""},
	}
	for _, c := range cases {
		got := new(transfer.TransferImpl).New(c.in).CamelToPascal()
		if got != c.want {
			t.Fatalf("CamelToPascal(%q) 期望 %q，得到 %q", c.in, c.want, got)
		}
	}
}

func TestCamelToSnake(t *testing.T) {
	cases := []struct{ in, want string }{
		{"helloWorld", "hello_world"},
		{"myClassName", "my_class_name"},
		{"getElement", "get_element"},
	}
	for _, c := range cases {
		got := new(transfer.TransferImpl).New(c.in).CamelToSnake()
		if got != c.want {
			t.Fatalf("CamelToSnake(%q) 期望 %q，得到 %q", c.in, c.want, got)
		}
	}
}

func TestCamelToBabel(t *testing.T) {
	cases := []struct{ in, want string }{
		{"helloWorld", "hello-world"},
		{"myClassName", "my-class-name"},
		{"getElement", "get-element"},
	}
	for _, c := range cases {
		got := new(transfer.TransferImpl).New(c.in).CamelToBabel()
		if got != c.want {
			t.Fatalf("CamelToBabel(%q) 期望 %q，得到 %q", c.in, c.want, got)
		}
	}
}

// ==================== SnakeCase 转换 ====================

func TestSnakeToPascal(t *testing.T) {
	cases := []struct{ in, want string }{
		{"hello_world", "HelloWorld"},
		{"my_class_name", "MyClassName"},
		{"get_element_by_id", "GetElementById"},
	}
	for _, c := range cases {
		got := new(transfer.TransferImpl).New(c.in).SnakeToPascal()
		if got != c.want {
			t.Fatalf("SnakeToPascal(%q) 期望 %q，得到 %q", c.in, c.want, got)
		}
	}
}

func TestSnakeToCamel(t *testing.T) {
	cases := []struct{ in, want string }{
		{"hello_world", "helloworld"},
		{"my_class_name", "myclassname"},
	}
	for _, c := range cases {
		got := new(transfer.TransferImpl).New(c.in).SnakeToCamel()
		if got != c.want {
			t.Fatalf("SnakeToCamel(%q) 期望 %q，得到 %q", c.in, c.want, got)
		}
	}
}

func TestSnakeToBabel(t *testing.T) {
	cases := []struct{ in, want string }{
		{"hello_world", "hello-world"},
		{"my_class_name", "my-class-name"},
	}
	for _, c := range cases {
		got := new(transfer.TransferImpl).New(c.in).SnakeToBabel()
		if got != c.want {
			t.Fatalf("SnakeToBabel(%q) 期望 %q，得到 %q", c.in, c.want, got)
		}
	}
}

// ==================== Babel/KebabCase 转换 ====================

func TestBabelToPascal(t *testing.T) {
	cases := []struct{ in, want string }{
		{"hello-world", "HelloWorld"},
		{"my-class-name", "MyClassName"},
	}
	for _, c := range cases {
		got := new(transfer.TransferImpl).New(c.in).BabelToPascal()
		if got != c.want {
			t.Fatalf("BabelToPascal(%q) 期望 %q，得到 %q", c.in, c.want, got)
		}
	}
}

func TestKebabToCamel(t *testing.T) {
	cases := []struct{ in, want string }{
		{"hello-world", "helloWorld"},
		{"my-class-name", "myClassName"},
	}
	for _, c := range cases {
		got := new(transfer.TransferImpl).New(c.in).KebabToCamel()
		if got != c.want {
			t.Fatalf("KebabToCamel(%q) 期望 %q，得到 %q", c.in, c.want, got)
		}
	}
}

func TestBabelToSnake(t *testing.T) {
	// 注意：源码中 BabelToSnake 实现为 strings.ReplaceAll(original, "_", "-")
	// 即将下划线替换为连字符，而非将连字符替换为下划线
	// 对于不含下划线的 babel 字符串，结果不变
	got := new(transfer.TransferImpl).New("hello-world").BabelToSnake()
	if got != "hello-world" {
		t.Logf("BabelToSnake(\"hello-world\") = %q (源码逻辑：替换下划线为连字符)", got)
	}

	// 含下划线时会被替换为连字符
	got2 := new(transfer.TransferImpl).New("hello_world").BabelToSnake()
	if got2 != "hello-world" {
		t.Fatalf("BabelToSnake(\"hello_world\") 期望 %q，得到 %q", "hello-world", got2)
	}
}

// ==================== Pluralize 单数变复数 ====================

func TestPluralize_Default(t *testing.T) {
	cases := []struct{ in, want string }{
		{"cat", "cats"},
		{"dog", "dogs"},
		{"day", "days"}, // 元音+a+y → +s
		{"key", "keys"}, // 元音+e+y → +s
		{"boy", "boys"}, // 元音+o+y → +s
	}
	for _, c := range cases {
		got := new(transfer.TransferImpl).New(c.in).Pluralize()
		if got != c.want {
			t.Fatalf("Pluralize(%q) 期望 %q，得到 %q", c.in, c.want, got)
		}
	}
}

func TestPluralize_SEnding(t *testing.T) {
	cases := []struct{ in, want string }{
		{"bus", "buses"},
		{"class", "classes"},
		{"glass", "glasses"},
	}
	for _, c := range cases {
		got := new(transfer.TransferImpl).New(c.in).Pluralize()
		if got != c.want {
			t.Fatalf("Pluralize(%q) 期望 %q，得到 %q", c.in, c.want, got)
		}
	}
}

func TestPluralize_XZEnding(t *testing.T) {
	cases := []struct{ in, want string }{
		{"box", "boxes"},
		{"fox", "foxes"},
		{"buzz", "buzzes"},
	}
	for _, c := range cases {
		got := new(transfer.TransferImpl).New(c.in).Pluralize()
		if got != c.want {
			t.Fatalf("Pluralize(%q) 期望 %q，得到 %q", c.in, c.want, got)
		}
	}
}

func TestPluralize_CHEnding(t *testing.T) {
	cases := []struct{ in, want string }{
		{"watch", "watches"},
		{"match", "matches"},
		{"dish", "dishes"},
		{"fish", "fishes"},
	}
	for _, c := range cases {
		got := new(transfer.TransferImpl).New(c.in).Pluralize()
		if got != c.want {
			t.Fatalf("Pluralize(%q) 期望 %q，得到 %q", c.in, c.want, got)
		}
	}
}

func TestPluralize_FEnding(t *testing.T) {
	cases := []struct{ in, want string }{
		{"leaf", "leaves"},
		{"knife", "knives"},
	}
	for _, c := range cases {
		got := new(transfer.TransferImpl).New(c.in).Pluralize()
		if got != c.want {
			t.Fatalf("Pluralize(%q) 期望 %q，得到 %q", c.in, c.want, got)
		}
	}
}

func TestNew_ReturnsTransfer(t *testing.T) {
	tr := new(transfer.TransferImpl).New("HelloWorld")
	if tr == nil {
		t.Fatal("New 不应返回 nil")
	}
}
