### 日志

#### Zap 日志

```go
package logger_test

import (
	"testing"

	"github.com/aid297/aid/v2/logs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestZapLog(t *testing.T) {
	t.Run("创建日志", func(t *testing.T) {
		var (
			zapLogger *zap.Logger
			err       error
		)

		if zapLogger, err = logs.NewZapLog(
			logs.Level(zapcore.DebugLevel),
			logs.Filename("."),
			logs.InConsole(false),
			logs.EncoderType(logs.EncoderTypeConsole),
			logs.Compress(true),
			logs.MaxBackup(30),
			logs.MaxSize(10),
			logs.MaxDay(30),
		); err != nil {
			t.Fatal(err)
		}

		zapLogger.Info("test-info", zap.String("a", "b"))
		zapLogger.Debug("test-debug", zap.String("c", "d"))
		zapLogger.Warn("test-warning", zap.Any("any", []any{"haha", "hehe", 1, 2, 3, 4}))
		zapLogger.Error("test-error", zap.Errors("errors", []error{errors.New("err1"), errors.New("err2"), errors.New("err3")}))
	})
}
```

#### JsonObject - 结构化 JSON 日志字段

当记录 JSON 数据时,使用 `JsonObject` 可以获得更易读的结构化输出:

```go
package logger_test

import (
	"testing"

	"github.com/aid297/aid/v2/logs"
	"go.uber.org/zap"
)

func TestJsonObject(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	// 简单的 JSON 对象
	t.Run("简单对象", func(t *testing.T) {
		jsonStr := `{"name":"张三","age":25,"city":"北京"}`
		obj := logs.NewJsonObject(jsonStr)
		logger.Info("用户信息", zap.Stringer("userInfo", obj))
		// 输出格式化的 JSON,而不是原始字符串或 map
	})

	// 嵌套对象
	t.Run("嵌套对象", func(t *testing.T) {
		jsonStr := `{"user":{"name":"李四","contact":{"email":"lisi@example.com"}}}`
		obj := logs.NewJsonObject(jsonStr)
		logger.Info("联系方式", zap.Stringer("data", obj))
	})

	// API 响应
	t.Run("API 响应", func(t *testing.T) {
		jsonStr := `{"code":200,"message":"success","data":{"total":100,"page":1}}`
		obj := logs.NewJsonObject(jsonStr)
		logger.Info("API响应", zap.Stringer("response", obj))
	})

	// 错误处理
	t.Run("无效 JSON", func(t *testing.T) {
		jsonStr := `{invalid json}`
		obj := logs.NewJsonObject(jsonStr)
		logger.Warn("无效JSON", zap.Stringer("data", obj))
		// 输出: [JSON Parse Error: ...] {invalid json}
	})
}
```

**与传统方式的对比:**

| 方式 | 输出效果 |
|------|----------|
| `zap.Any` | `{map[code:200 message:success]}` |
| `zap.String` | `{"code":200,"message":"success"}` (一行) |
| **JsonObject** | 格式化的多行结构,易于阅读 |

详细说明请参考: [README_JSON_FIELD.md](./README_JSON_FIELD.md) |  [COMPARISON.md](./COMPARISON.md)

