package str

import (
	"testing"
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
	result := NewTemplateV2(tpl, data)

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
	result := NewTemplateV2(tpl, data)

	if result.Error() != nil {
		t.Fatalf("unexpected error: %v", result.Error())
	}

	expected := "Jerry - 25 - Shanghai"
	if result.String() != expected {
		t.Errorf("expected %q, got %q", expected, result.String())
	}
}

func TestNewTemplate2_NonStruct(t *testing.T) {
	result := NewTemplateV2("hello", "not a struct")
	if result.Error() == nil {
		t.Fatal("expected error for non-struct input, got nil")
	}
}

func TestNewTemplate2_NoTag(t *testing.T) {
	type NoTag struct {
		Name string
	}
	result := NewTemplateV2("{{name}}", NoTag{Name: "test"})
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
	result := NewTemplateV2("{{name}}", data)
	if string(result.Bytes()) != "A" {
		t.Errorf("expected %q, got %q", "A", string(result.Bytes()))
	}
}

func TestTemplate2_Content(t *testing.T) {
	data := testTemplate2Data{Name: "A", Age: 1, City: "X"}
	content := "{{name}}-{{age}}"
	result := NewTemplateV2(content, data)
	if result.Content() != content {
		t.Errorf("expected content %q, got %q", content, result.Content())
	}
}
