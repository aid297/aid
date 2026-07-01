package template

type TemplateAttr func(t Templater)

// Struct 传入结构体或结构体指针作为模板数据源
func Struct(s any) TemplateAttr { return func(t Templater) { t.setS(s) } }

// Map 传入 map[string]any 作为模板数据源
func Map(s map[string]any) TemplateAttr { return func(t Templater) { t.setS(s) } }
