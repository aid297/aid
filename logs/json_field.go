package logs

import (
	"github.com/bytedance/sonic"
	"go.uber.org/zap"
)

// JSONString 从 JSON 字符串创建 JSONStringerImpl
func JSONString(key, raw string) zap.Field {
	if raw == "" {
		return zap.String(key, "")
	}

	var data any
	if err := sonic.Unmarshal([]byte(raw), &data); err != nil {
		return zap.String(key, raw)
	}

	return zap.Any(key, data)
}
