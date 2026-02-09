### 日志

```go
package main

import (
	`errors`

	`github.com/aid297/aid/logger`
	`go.uber.org/zap`
	`go.uber.org/zap/zapcore`
)

func main() {
	var (
		zapLogger *zap.Logger
		err       error
	)

	if zapLogger, err = logger.APP.Zap.New(
		logger.APP.ZapConfig.
			New(zapcore.ErrorLevel).
			SetPath(".").
			SetPathAbs(false).
			SetInConsole(false).
			SetEncoderType(logger.EncoderTypeConsole).
			SetNeedCompress(true).
			SetMaxBackup(30).
			SetMaxSize(10).
			SetMaxDay(30),
	); err != nil {
		panic(err)
	}

	zapLogger.Info("test-info", zap.String("a", "b"))
	zapLogger.Debug("test-debug", zap.String("c", "d"))
	zapLogger.Warn("test-warning", zap.Any("any", []any{"haha", "hehe", 1, 2, 3, 4}))
	zapLogger.Error("test-error", zap.Errors("errors", []error{errors.New("err1"), errors.New("err2"), errors.New("err3")}))
}
```

