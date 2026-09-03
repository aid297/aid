package logs

import (
	"github.com/bytedance/sonic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ JSONStringer = (*JSONStringerImpl)(nil)

type (
	JSONStringer interface {
		MarshalLogObject(enc zapcore.ObjectEncoder) error
	}

	JSONStringerImpl struct {
		key         string
		raw         string
		unmarshaled any
		formatStr   string
		err         error
	}
)

// JSONString 从 JSON 字符串创建 JSONStringerImpl
func JSONString(key, raw string) zap.Field {
	if raw == "" {
		return zap.Object(key, &JSONStringerImpl{})
	}

	my := &JSONStringerImpl{raw: raw, key: key}

	if err := sonic.Unmarshal([]byte(raw), &my.unmarshaled); err != nil {
		my.formatStr = raw
		return zap.String(key, raw)
	}

	return zap.Object(key, my)
}

// MarshalLogObject 实现 zapcore.ObjectMarshaler 接口
// 这是让类型支持 zap 的标准方式
func (my *JSONStringerImpl) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if my.unmarshaled != nil {
		_ = enc.AddObject(my.key, my)
	}
	return nil
}
