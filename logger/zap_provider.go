package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/aid297/aid/filesystem"
	"github.com/aid297/aid/operation"
)

// ZapProvider Zap日志服务提供者
type (
	ZapProvider struct{}

	EncoderType = string
)

const (
	EncoderTypeConsole EncoderType = "CONSOLE"
	EncoderTypeJson    EncoderType = "JSON"
)

// getWriteSync 获取 zapcore.WriteSync
func getWriteSync(config *zapConfig, path string) zapcore.WriteSyncer {
	fileWriter := &lumberjack.Logger{
		Filename:   path,                // 日志文件名称
		MaxSize:    config.MaxSize,      // 文件大小限制,单位MB
		MaxBackups: config.MaxBackup,    // 最大保留日志文件数量
		MaxAge:     config.MaxDay,       // 日志文件保留天数
		Compress:   config.NeedCompress, // 是否压缩处理,压缩以后文件为xxxxx.gz
	}

	if config.InConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(fileWriter), zapcore.AddSync(os.Stdout))
	} else {
		return zapcore.AddSync(fileWriter)
	}
}

// New 实例化：Zap日志服务提供者
func (*ZapProvider) New(config *zapConfig) (*zap.Logger, error) { return NewZapProvider(config) }

// NewZapProvider 实例化：Zap日志服务提供者
//
//go:fix 推荐使用：New方法
func NewZapProvider(config *zapConfig) (*zap.Logger, error) {
	var (
		err             error
		fs              *filesystem.FileSystem
		zapLogger       *zap.Logger
		zapCores        = make([]zapcore.Core, 0, 7)
		zapLoggerConfig = zapcore.EncoderConfig{
			MessageKey:    "message",
			LevelKey:      "level",
			TimeKey:       "time",
			NameKey:       "logger",
			CallerKey:     "caller",
			StacktraceKey: "stacktrace",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   zapcore.LowercaseLevelEncoder,
			EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
				encoder.AppendString(t.Format(time.DateTime + ".000"))
			},
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
		}
		encoderTypes = map[EncoderType]func(cfg zapcore.EncoderConfig) zapcore.Encoder{
			EncoderTypeJson:    zapcore.NewJSONEncoder,
			EncoderTypeConsole: zapcore.NewConsoleEncoder,
		}
	)

	fs = operation.TernaryFuncAll(func() bool { return config.PathAbs }, func() *filesystem.FileSystem { return filesystem.FileSystemApp.NewByAbs(config.Path) }, func() *filesystem.FileSystem { return filesystem.FileSystemApp.NewByRel(config.Path) })
	if !fs.IsExist {
		if err = fs.MkDir(); err != nil {
			return nil, fmt.Errorf("创建日志目录失败：%w", err)
		}
	}

	if config.Level < zapcore.DebugLevel {
		config.Level = zapcore.DebugLevel
	}

	if config.Level > zapcore.FatalLevel {
		config.Level = zapcore.FatalLevel
	}

	for _, logLevel := range []zapcore.Level{zapcore.DebugLevel, zapcore.InfoLevel, zapcore.WarnLevel, zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel} {
		if config.Level >= logLevel {
			writer := getWriteSync(config, fs.Copy().Join(logLevel.String()+config.Extension).GetDir())
			zapCores = append(zapCores, zapcore.NewCore(encoderTypes[config.EncoderType](zapLoggerConfig), writer, logLevel))
		}
	}

	zapLogger = zap.New(zapcore.NewTee(zapCores...))

	defer func() {
		if config.InConsole {
			return
		}
		err = zapLogger.Sync()
	}()

	return zapLogger, err
}
