package logs

import (
	_sonic "github.com/bytedance/sonic"
	_zap "go.uber.org/zap"
)

// JSONString 从 JSON 字符串创建 zap.Field
func JSONString(key, raw string) _zap.Field {
	if raw == "" {
		return _zap.String(key, raw)
	}

	var data any
	if err := _sonic.UnmarshalString(raw, &data); err != nil {
		return _zap.String(key, raw)
	}

	return _zap.Any(key, data)
}

func JSONBytes(key string, raw []byte) _zap.Field {
	if len(raw) == 0 {
		return _zap.Binary(key, raw)
	}

	var data any
	if err := _sonic.Unmarshal(raw, &data); err != nil {
		return _zap.Binary(key, raw)
	}

	return _zap.Any(key, data)
}
