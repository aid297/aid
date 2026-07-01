package template

type TemplateAttr func(t *Template)

// Struct 传入结构体或结构体指针作为模板数据源
func Struct[T any](s T) TemplateAttr { return func(t *Template) { t.s = s } }

// Map 传入 map[string]any 作为模板数据源
func Map[T ~map[string]any](s T) TemplateAttr { return func(t *Template) { t.s = s } }
