package logs

import (
	"bytes"
	"encoding/json"
	"fmt"
)

var _ JSONStringer = (*JSONStringImpl)(nil)

type (
	JSONStringer interface {
		formatValue(v any, indent string, depth int) string
		formatMap(m map[string]any, indent string, depth int) string
		formatArray(arr []any, indent string, depth int) string
		String() string
		Error() error
		Value() any
		Raw() string
	}

	// JSONStringImpl 是一个用于日志的结构化 JSON 字段类型
	// 它会接收 JSON 字符串并在日志中以格式化后的结构显示
	JSONStringImpl struct {
		rawData   string // 原始 JSON 字符串
		parsed    any    // 解析后的数据
		formatStr string // 格式化后的字符串
		error     error  // 解析错误
	}
)

// JSONString 从 JSON 字符串创建 JSONStringImpl
// 如果字符串不是有效的 JSON,会保留原始字符串并记录错误
func JSONString(jsonStr string) *JSONStringImpl {
	if jsonStr == "" {
		return &JSONStringImpl{
			rawData:   "",
			formatStr: "",
		}
	}

	my := &JSONStringImpl{
		rawData: jsonStr,
	}

	// 尝试解析 JSON
	var temp any
	if err := json.Unmarshal([]byte(jsonStr), &temp); err != nil {
		my.error = err
		my.formatStr = jsonStr // 解析失败时显示原始字符串
		return my
	}

	my.parsed = temp
	my.formatStr = my.formatValue(temp, "", 0)
	return my
}

// formatValue 递归格式化 JSON 值
func (my *JSONStringImpl) formatValue(v any, indent string, depth int) string {
	if depth > 10 { // 防止过深的嵌套
		return "{...}"
	}

	switch val := v.(type) {
	case map[string]any:
		return my.formatMap(val, indent, depth)
	case []any:
		return my.formatArray(val, indent, depth)
	case string:
		return fmt.Sprintf("\"%s\"", val)
	case float64:
		// 处理数字类型(JSON unmarshal后数字都是float64)
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%g", val)
	case bool:
		return fmt.Sprintf("%t", val)
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", val)
	}
}

// formatMap 格式化对象
func (my *JSONStringImpl) formatMap(m map[string]any, indent string, depth int) string {
	if len(m) == 0 {
		return "{}"
	}

	var buf bytes.Buffer
	buf.WriteString("{\n")
	nextIndent := indent + "  "

	idx := 0
	for k, v := range m {
		idx++
		buf.WriteString(nextIndent)
		buf.WriteString(fmt.Sprintf("\"%s\": ", k))
		buf.WriteString(my.formatValue(v, nextIndent, depth+1))
		if idx < len(m) {
			buf.WriteString(",")
		}
		buf.WriteString("\n")
	}

	buf.WriteString(indent)
	buf.WriteString("}")
	return buf.String()
}

// formatArray 格式化数组
func (my *JSONStringImpl) formatArray(arr []any, indent string, depth int) string {
	if len(arr) == 0 {
		return "[]"
	}

	var buf bytes.Buffer
	buf.WriteString("[\n")
	nextIndent := indent + "  "

	for idx, v := range arr {
		buf.WriteString(nextIndent)
		buf.WriteString(my.formatValue(v, nextIndent, depth+1))
		if idx < len(arr)-1 {
			buf.WriteString(",")
		}
		buf.WriteString("\n")
	}

	buf.WriteString(indent)
	buf.WriteString("]")
	return buf.String()
}

// String 返回格式化的 JSON 字符串,用于日志输出
func (my *JSONStringImpl) String() string {
	if my.error != nil {
		return fmt.Sprintf("[JSON Parse Error: %v] %s", my.error, my.rawData)
	}
	if my.formatStr == "" {
		return ""
	}
	return my.formatStr
}

// Error 返回解析错误(如果有)
func (my *JSONStringImpl) Error() error {
	return my.error
}

// Value 返回原始值,可用于进一步处理
func (my *JSONStringImpl) Value() any {
	return my.parsed
}

// Raw 返回原始 JSON 字符串
func (my *JSONStringImpl) Raw() string {
	return my.rawData
}
