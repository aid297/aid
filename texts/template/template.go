package template

import (
	"fmt"
	"reflect"
	"strings"
)

var (
	_ Templater = (*Template)(nil)
)

// Template 基于 struct tag 的模板渲染
// 通过 `template:"name"` 标签定义模板变量，自动将模板中的 {{name}} 替换为字段值
type (
	Templater interface {
		New(content string, attrs ...TemplateAttr) Templater
		rander() Templater
		Error() error
		String() string
		Bytes() []byte
		Content() string
		extractTemplateVars() (map[string]string, error)
	}

	Template struct {
		err     error
		content string
		s       any
		ret     string
	}
)

// New 根据结构体的 template tag 渲染模板
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
func (my *Template) New(content string, attrs ...TemplateAttr) Templater {
	ins := &Template{content: content}

	for idx := range attrs {
		attrs[idx](ins)
	}

	return ins.rander()
}

func (my *Template) rander() Templater {
	vars, err := my.extractTemplateVars()
	if err != nil {
		my.err = fmt.Errorf("提取模板变量失败：%w", err)
		return my
	}

	result := my.content
	for tag, value := range vars {
		placeholder := "{{" + tag + "}}"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	my.ret = result

	return my
}

func (my *Template) Error() error    { return my.err }
func (my *Template) String() string  { return my.ret }
func (my *Template) Bytes() []byte   { return []byte(my.ret) }
func (my *Template) Content() string { return my.content }

// extractTemplateVars 通过反射提取结构体或 map 中的模板变量及其值
func (my *Template) extractTemplateVars() (map[string]string, error) {
	v := reflect.ValueOf(my.s)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	// 处理 map[string]any 及其自定义类型
	if v.Kind() == reflect.Map && v.Type().Key().Kind() == reflect.String {
		vars := make(map[string]string, v.Len())
		iter := v.MapRange()
		for iter.Next() {
			key := iter.Key().String()
			vars[key] = fmt.Sprintf("%v", iter.Value().Interface())
		}
		return vars, nil
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("参数必须是结构体、结构体指针或 map[string]any，当前类型：%s", v.Kind())
	}

	t := v.Type()
	vars := make(map[string]string)

	for i := range t.NumField() {
		field := t.Field(i)
		tag := field.Tag.Get("template")
		if tag == "" || tag == "-" {
			continue
		}
		vars[tag] = fmt.Sprintf("%v", v.Field(i).Interface())
	}

	return vars, nil
}
