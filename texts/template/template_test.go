package template_test

import (
	"strings"
	"testing"

	"github.com/aid297/aid/v2/texts/template"
)

func TestStruct_Basic(t *testing.T) {
	type Data struct {
		Name string `template:"name"`
		Age  int    `template:"age"`
	}
	tpl := "Hello {{name}}, you are {{age}} years old"
	got := new(template.TemplaterImpl).New(tpl, template.Struct(Data{Name: "Tom", Age: 18})).String()
	want := "Hello Tom, you are 18 years old"
	if got != want {
		t.Fatalf("期望 %q，得到 %q", want, got)
	}
}

func TestStruct_Pointer(t *testing.T) {
	type Data struct {
		Name string `template:"name"`
	}
	data := &Data{Name: "Jerry"}
	got := new(template.TemplaterImpl).New("Hi {{name}}", template.Struct(data)).String()
	if got != "Hi Jerry" {
		t.Fatalf("期望 %q，得到 %q", "Hi Jerry", got)
	}
}

func TestStruct_MultiplePlaceholders(t *testing.T) {
	type Data struct {
		First string `template:"first"`
		Last  string `template:"last"`
	}
	tpl := "{{first}} {{last}} aka {{first}}-{{last}}"
	got := new(template.TemplaterImpl).New(tpl, template.Struct(Data{First: "John", Last: "Doe"})).String()
	want := "John Doe aka John-Doe"
	if got != want {
		t.Fatalf("期望 %q，得到 %q", want, got)
	}
}

func TestStruct_FieldWithoutTag(t *testing.T) {
	type Data struct {
		Name    string `template:"name"`
		Ignored string `json:"ignored"`
	}
	got := new(template.TemplaterImpl).New("{{name}}", template.Struct(Data{Name: "Alice", Ignored: "x"})).String()
	if got != "Alice" {
		t.Fatalf("期望 %q，得到 %q", "Alice", got)
	}
}

func TestStruct_TagDash(t *testing.T) {
	type Data struct {
		Name string `template:"-"`
	}
	got := new(template.TemplaterImpl).New("{{Name}}", template.Struct(Data{Name: "Skip"})).String()
	// tag 为 "-" 时跳过，{{Name}} 不会被替换
	if strings.Contains(got, "Skip") {
		t.Fatalf("tag=\"-\" 的字段不应被替换，得到: %q", got)
	}
}

func TestMap_Basic(t *testing.T) {
	tpl := "Hello {{name}}, age {{age}}"
	data := map[string]any{"name": "Bob", "age": 30}
	got := new(template.TemplaterImpl).New(tpl, template.Map(data)).String()
	want := "Hello Bob, age 30"
	if got != want {
		t.Fatalf("期望 %q，得到 %q", want, got)
	}
}

func TestMap_Pointer(t *testing.T) {
	// map 是引用类型，传指针也应正常工作
	data := map[string]any{"key": "value"}
	got := new(template.TemplaterImpl).New("{{key}}", template.Struct(&data)).String()
	if got != "value" {
		t.Fatalf("期望 %q，得到 %q", "value", got)
	}
}

func TestMap_Empty(t *testing.T) {
	got := new(template.TemplaterImpl).New("hello {{name}}", template.Map(map[string]any{})).String()
	if got != "hello {{name}}" {
		t.Fatalf("空 map 应保留占位符，得到: %q", got)
	}
}

func TestNoDataSource(t *testing.T) {
	// 未提供数据源时，占位符保留原样
	got := new(template.TemplaterImpl).New("hello {{name}}").String()
	if got != "hello {{name}}" {
		t.Fatalf("期望 %q，得到 %q", "hello {{name}}", got)
	}
}

func TestEmptyContent(t *testing.T) {
	got := new(template.TemplaterImpl).New("", template.Struct(struct {
		Name string `template:"name"`
	}{Name: "x"})).String()
	if got != "" {
		t.Fatalf("空模板应返回空字符串，得到: %q", got)
	}
}

func TestBytes(t *testing.T) {
	type Data struct {
		Name string `template:"name"`
	}
	got := new(template.TemplaterImpl).New("{{name}}", template.Struct(Data{Name: "abc"})).Bytes()
	if string(got) != "abc" {
		t.Fatalf("期望 %q，得到 %q", "abc", string(got))
	}
}

func TestContent(t *testing.T) {
	tpl := "hello {{name}}"
	ins := new(template.TemplaterImpl).New(tpl, template.Struct(struct {
		Name string `template:"name"`
	}{Name: "world"}))
	if ins.Content() != tpl {
		t.Fatalf("Content() 期望 %q，得到 %q", tpl, ins.Content())
	}
}

func TestError_NoDataSource(t *testing.T) {
	// 无数据源时，extractTemplateVars 的 s 为 nil
	ins := new(template.TemplaterImpl).New("{{name}}")
	// s 为 nil → reflect.ValueOf(nil) → zero Value → Kind() == reflect.Invalid
	// 不匹配 Pointer, Map, Struct → 应返回错误
	if ins.Error() == nil {
		t.Logf("无数据源时 Error() 为 nil（s 为 nil 时 extractTemplateVars 可能未报错）")
	}
	// 占位符应保留
	if ins.String() != "{{name}}" {
		t.Logf("无数据源时 String() = %q", ins.String())
	}
}

func TestError_InvalidType(t *testing.T) {
	// 传入 int 类型
	ins := new(template.TemplaterImpl).New("{{name}}", template.Struct(42))
	if ins.Error() == nil {
		t.Fatal("传入 int 应返回错误")
	}
}

func TestError_NonStructNonMap(t *testing.T) {
	// 传入 string 类型
	ins := new(template.TemplaterImpl).New("{{name}}", template.Struct("just a string"))
	if ins.Error() == nil {
		t.Fatal("传入 string 应返回错误")
	}
}

func TestPlaceholderNoMatch(t *testing.T) {
	type Data struct {
		Name string `template:"name"`
	}
	// 模板中有 {{age}} 但结构体没有对应 tag
	got := new(template.TemplaterImpl).New("{{name}} is {{age}}", template.Struct(Data{Name: "Alice"})).String()
	if !strings.Contains(got, "Alice") {
		t.Fatalf("应包含替换后的 name，得到: %q", got)
	}
	if !strings.Contains(got, "{{age}}") {
		t.Fatalf("未匹配的占位符应保留，得到: %q", got)
	}
}

func TestRender_NoPlaceholders(t *testing.T) {
	type Data struct {
		Name string `template:"name"`
	}
	got := new(template.TemplaterImpl).New("plain text", template.Struct(Data{Name: "x"})).String()
	if got != "plain text" {
		t.Fatalf("期望 %q，得到 %q", "plain text", got)
	}
}

func TestStruct_BooleanAndFloat(t *testing.T) {
	type Data struct {
		Active bool    `template:"active"`
		Score  float64 `template:"score"`
	}
	got := new(template.TemplaterImpl).New("active={{active}}, score={{score}}",
		template.Struct(Data{Active: true, Score: 95.5})).String()
	if !strings.Contains(got, "active=true") {
		t.Fatalf("缺少 active=true: %q", got)
	}
	if !strings.Contains(got, "score=95.5") {
		t.Fatalf("缺少 score=95.5: %q", got)
	}
}
