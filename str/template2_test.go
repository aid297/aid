package str_test

import (
	"testing"

	"github.com/aid297/aid/v2/str"
)

type testTemplate2Data struct {
	Name    string `template:"name"`
	Age     int    `template:"age"`
	City    string `template:"city"`
	Ignored string
	Skipped string `template:"-"`
}

func TestNewTemplate2(t *testing.T) {
	data := testTemplate2Data{
		Name:    "Tom",
		Age:     18,
		City:    "Beijing",
		Ignored: "should not appear",
		Skipped: "should not appear",
	}

	tpl := "Hello {{name}}, you are {{age}} years old, living in {{city}}."
	result := str.NewTemplateV2(tpl, str.Struct(data))

	if result.Error() != nil {
		t.Fatalf("unexpected error: %v", result.Error())
	}

	expected := "Hello Tom, you are 18 years old, living in Beijing."
	if result.String() != expected {
		t.Errorf("expected %q, got %q", expected, result.String())
	}
}

func TestNewTemplate2_Pointer(t *testing.T) {
	data := &testTemplate2Data{
		Name: "Jerry",
		Age:  25,
		City: "Shanghai",
	}

	tpl := "{{name}} - {{age}} - {{city}}"
	result := str.NewTemplateV2(tpl, str.Struct(data))

	if result.Error() != nil {
		t.Fatalf("unexpected error: %v", result.Error())
	}

	expected := "Jerry - 25 - Shanghai"
	if result.String() != expected {
		t.Errorf("expected %q, got %q", expected, result.String())
	}
}

func TestNewTemplate2_NonStruct(t *testing.T) {
	result := str.NewTemplateV2("hello", str.Struct("not a struct"))
	if result.Error() == nil {
		t.Fatal("expected error for non-struct input, got nil")
	}
}

func TestNewTemplate2_NoTag(t *testing.T) {
	type NoTag struct {
		Name string
	}
	result := str.NewTemplateV2("{{name}}", str.Struct(NoTag{Name: "test"}))
	if result.Error() != nil {
		t.Fatalf("unexpected error: %v", result.Error())
	}
	// 没有 template tag，不会替换
	if result.String() != "{{name}}" {
		t.Errorf("expected %q, got %q", "{{name}}", result.String())
	}
}

func TestTemplate2_Bytes(t *testing.T) {
	data := testTemplate2Data{Name: "A", Age: 1, City: "X"}
	result := str.NewTemplateV2("{{name}}", str.Struct(data))
	if string(result.Bytes()) != "A" {
		t.Errorf("expected %q, got %q", "A", string(result.Bytes()))
	}
}

func TestTemplate2_Content(t *testing.T) {
	data := testTemplate2Data{Name: "A", Age: 1, City: "X"}
	content := "{{name}}-{{age}}"
	result := str.NewTemplateV2(content, str.Struct(data))
	if result.Content() != content {
		t.Errorf("expected content %q, got %q", content, result.Content())
	}
}

func TestNewTemplateV2_Map(t *testing.T) {
	data := map[string]any{
		"name": "Alice",
		"age":  30,
		"city": "Shenzhen",
	}

	tpl := "Hello {{name}}, you are {{age}} years old, living in {{city}}."
	result := str.NewTemplateV2(tpl, str.Map(data))

	if result.Error() != nil {
		t.Fatalf("unexpected error: %v", result.Error())
	}

	expected := "Hello Alice, you are 30 years old, living in Shenzhen."
	if result.String() != expected {
		t.Errorf("expected %q, got %q", expected, result.String())
	}
}

type myMapType map[string]any

func TestNewTemplateV2_MapCustomType(t *testing.T) {
	data := myMapType{
		"name": "Bob",
		"age":  42,
	}

	tpl := "{{name}} is {{age}}"
	result := str.NewTemplateV2(tpl, str.Map(data))

	if result.Error() != nil {
		t.Fatalf("unexpected error: %v", result.Error())
	}

	expected := "Bob is 42"
	if result.String() != expected {
		t.Errorf("expected %q, got %q", expected, result.String())
	}
}

func TestNewTemplateV2_MapPartialReplace(t *testing.T) {
	data := map[string]any{
		"name": "Charlie",
	}

	tpl := "{{name}} - {{missing}}"
	result := str.NewTemplateV2(tpl, str.Map(data))

	if result.Error() != nil {
		t.Fatalf("unexpected error: %v", result.Error())
	}

	// 没有对应的 key，占位符保留
	expected := "Charlie - {{missing}}"
	if result.String() != expected {
		t.Errorf("expected %q, got %q", expected, result.String())
	}
}

func TestNewTemplateV2_MapEmpty(t *testing.T) {
	data := map[string]any{}

	tpl := "{{name}}"
	result := str.NewTemplateV2(tpl, str.Map(data))

	if result.Error() != nil {
		t.Fatalf("unexpected error: %v", result.Error())
	}

	if result.String() != "{{name}}" {
		t.Errorf("expected %q, got %q", "{{name}}", result.String())
	}
}

func TestNewTemplateV2_MapNilValue(t *testing.T) {
	data := map[string]any{
		"name": nil,
	}

	tpl := "value: {{name}}"
	result := str.NewTemplateV2(tpl, str.Map(data))

	if result.Error() != nil {
		t.Fatalf("unexpected error: %v", result.Error())
	}

	expected := "value: <nil>"
	if result.String() != expected {
		t.Errorf("expected %q, got %q", expected, result.String())
	}
}
