package template

import (
	"fmt"
	"reflect"
	"strings"
)

// TemplaterImpl 基于 struct tag 的模板渲染
// 通过 `template:"name"` 标签定义模板变量，自动将模板中的 {{name}} 替换为字段值
type (
	Templater interface {
		New(content string, attrs ...TemplateAttr) Templater
		Error() error
		String() string
		Bytes() []byte
		Content() string
		setS(s any)
		render() Templater
		extractTemplateVars() (map[string]string, error)
	}

	TemplaterImpl struct {
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
func (my *TemplaterImpl) New(content string, attrs ...TemplateAttr) Templater {
	ins := &TemplaterImpl{content: content}

	for idx := range attrs {
		attrs[idx](ins)
	}

	return ins.render()
}

func (my *TemplaterImpl) Error() error    { return my.err }
func (my *TemplaterImpl) String() string  { return my.ret }
func (my *TemplaterImpl) Bytes() []byte   { return []byte(my.ret) }
func (my *TemplaterImpl) Content() string { return my.content }

func (my *TemplaterImpl) setS(s any) { my.s = s }

func (my *TemplaterImpl) render() Templater {
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

// extractTemplateVars 通过反射提取结构体或 map 中的模板变量及其值
func (my *TemplaterImpl) extractTemplateVars() (map[string]string, error) {
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
