package str

import (
	"fmt"
	"reflect"
	"strings"
)

// TemplateV2 基于 struct tag 的模板渲染
// 通过 `template:"name"` 标签定义模板变量，自动将模板中的 {{name}} 替换为字段值
type TemplateV2[T any] struct {
	err     error
	content string
	s       T
	ret     string
}

// NewTemplateV2 根据结构体的 template tag 渲染模板
// content 中使用 {{tagName}} 作为占位符，tagName 对应结构体字段上的 `template:"tagName"` 标签值
//
// 示例：
//
//	type Data struct {
//	    Name string `template:"name"`
//	    Age  int    `template:"age"`
//	}
//	tpl := "Hello {{name}}, you are {{age}} years old"
//	result := NewTemplateV2(tpl, Data{Name: "Tom", Age: 18}).String()
//	// => "Hello Tom, you are 18 years old"
func NewTemplateV2[T any](content string, s T) *TemplateV2[T] {
	t := &TemplateV2[T]{content: content, s: s}

	vars, err := extractTemplateVars(s)
	if err != nil {
		t.err = fmt.Errorf("提取模板变量失败：%w", err)
		return t
	}

	result := content
	for tag, value := range vars {
		placeholder := "{{" + tag + "}}"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	t.ret = result

	return t
}

func (my *TemplateV2[T]) Error() error    { return my.err }
func (my *TemplateV2[T]) String() string  { return my.ret }
func (my *TemplateV2[T]) Bytes() []byte   { return []byte(my.ret) }
func (my *TemplateV2[T]) Content() string { return my.content }

// extractTemplateVars 通过反射提取结构体上所有带 template tag 的字段及其值
func extractTemplateVars(s any) (map[string]string, error) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("参数必须是结构体或结构体指针，当前类型：%s", v.Kind())
	}

	t := v.Type()
	vars := make(map[string]string)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("template")
		if tag == "" || tag == "-" {
			continue
		}
		vars[tag] = fmt.Sprintf("%v", v.Field(i).Interface())
	}

	return vars, nil
}
