package str

// Struct 传入结构体或结构体指针作为模板数据源
func Struct[T any](s T) TemplateV2Attr { return func(t *TemplateV2) { t.s = s } }

// Map 传入 map[string]any 作为模板数据源
func Map[T ~map[string]any](s T) TemplateV2Attr { return func(t *TemplateV2) { t.s = s } }
