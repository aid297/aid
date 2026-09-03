package logs

import (
	"bytes"
	"fmt"

	"github.com/bytedance/sonic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ JSONStringer = (*JSONStringerImpl)(nil)

type (
	JSONStringer interface {
		MarshalLogObject(enc zapcore.ObjectEncoder) error
		String() string
		Value() any
		Raw() string
		Error() error
	}

	JSONStringerImpl struct {
		rawData   string
		parsed    any
		formatStr string
		err       error
	}
)

// JSONString 从 JSON 字符串创建 JSONStringerImpl
func JSONString(key, raw string) zap.Field {
	if raw == "" {
		return zap.Object(key, &JSONStringerImpl{})
	}

	my := &JSONStringerImpl{rawData: raw}

	var temp any
	if err := sonic.Unmarshal([]byte(raw), &temp); err != nil {
		my.err = err
		my.formatStr = raw
		return zap.Object(key, my)
	}

	my.parsed = temp
	my.formatStr = my.formatValue(temp, "", 0)
	return zap.Object(key, my)
}

// MarshalLogObject 实现 zapcore.ObjectMarshaler 接口
// 这是让类型支持 zap 的标准方式
func (my *JSONStringerImpl) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if my.err != nil {
		enc.AddString("json", fmt.Sprintf("[ERROR: %v] %s", my.err, my.rawData))
		return nil
	}

	if my.parsed != nil {
		enc.AddString("json", my.formatStr)
	} else {
		enc.AddString("json", "")
	}
	return nil
}

// String 返回格式化字符串
func (my *JSONStringerImpl) String() string {
	if my.err != nil {
		return fmt.Sprintf("[JSON Parse Error: %v] %s", my.err, my.rawData)
	}
	return my.formatStr
}

// Value 返回解析后的值
func (my *JSONStringerImpl) Value() any {
	return my.parsed
}

// Raw 返回原始字符串
func (my *JSONStringerImpl) Raw() string {
	return my.rawData
}

// Error 返回解析错误
func (my *JSONStringerImpl) Error() error {
	return my.err
}

// formatValue 递归格式化 JSON 值
func (my *JSONStringerImpl) formatValue(v any, indent string, depth int) string {
	if depth > 10 {
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

func (my *JSONStringerImpl) formatMap(m map[string]any, indent string, depth int) string {
	if len(m) == 0 {
		return "{}"
	}

	buf := new(bytes.Buffer)
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

func (my *JSONStringerImpl) formatArray(arr []any, indent string, depth int) string {
	if len(arr) == 0 {
		return "[]"
	}

	buf := new(bytes.Buffer)
	buf.WriteString("[\n")
	nextIndent := indent + "  "

	for i, v := range arr {
		buf.WriteString(nextIndent)
		buf.WriteString(my.formatValue(v, nextIndent, depth+1))
		if i < len(arr)-1 {
			buf.WriteString(",")
		}
		buf.WriteString("\n")
	}

	buf.WriteString(indent)
	buf.WriteString("]")
	return buf.String()
}
