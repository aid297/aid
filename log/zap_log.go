package log

import (
	"fmt"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/aid297/aid/v2/filesystem"
	"github.com/aid297/aid/v2/operation"
)

// ZapProvider Zap日志服务提供者
type (
	ZapLog struct {
		Level       zapcore.Level
		Path        string
		MaxSize     int
		MaxBackup   int
		MaxDay      int
		Compress    bool
		InConsole   bool
		Extension   string
		EncoderType ZapLogEncoderType
	}

	ZapLogEncoderType string

	ZapLogger interface {
		SetLevel(level zapcore.Level) (err error)
		SetPath(path string) (err error)
		SetMaxSize(maxSize int) (err error)
		SetMaxBackup(maxBackup int) (err error)
		SetMaxDay(maxDay int) (err error)
		SetCompress(compress bool) (err error)
		SetInConsole(inConsole bool) (err error)
		SetExtension(extension string) (err error)
		SetEncoderType(encoderType ZapLogEncoderType) (err error)
		GetLevel() zapcore.Level
		GetPath() string
		GetMaxSize() int
		GetMaxBackup() int
		GetMaxDay() int
		GetCompress() bool
		GetInConsole() bool
		GetExtension() string
		SetAttrs(attrs ...ZapLogAttr) (err error)
	}
)

const (
	EncoderTypeConsole ZapLogEncoderType = "CONSOLE"
	EncoderTypeJson    ZapLogEncoderType = "JSON"
)

func (zapLog *ZapLog) SetLevel(level zapcore.Level) (err error) { zapLog.Level = level; return nil }
func (zapLog *ZapLog) SetPath(path string) (err error)          { zapLog.Path = path; return nil }
func (zapLog *ZapLog) SetMaxSize(maxSize int) (err error)       { zapLog.MaxSize = maxSize; return nil }
func (zapLog *ZapLog) SetMaxBackup(maxBackup int) (err error) {
	zapLog.MaxBackup = maxBackup
	return nil
}
func (zapLog *ZapLog) SetMaxDay(maxDay int) (err error)      { zapLog.MaxDay = maxDay; return nil }
func (zapLog *ZapLog) SetCompress(compress bool) (err error) { zapLog.Compress = compress; return nil }
func (zapLog *ZapLog) SetInConsole(inConsole bool) (err error) {
	zapLog.InConsole = inConsole
	return nil
}
func (zapLog *ZapLog) SetExtension(extension string) (err error) {
	zapLog.Extension = extension
	return nil
}
func (zapLog *ZapLog) SetEncoderType(encoderType ZapLogEncoderType) (err error) {
	zapLog.EncoderType = encoderType
	return nil
}
func (zapLog *ZapLog) GetLevel() zapcore.Level { return zapLog.Level }
func (zapLog *ZapLog) GetPath() string         { return zapLog.Path }
func (zapLog *ZapLog) GetMaxSize() int         { return zapLog.MaxSize }
func (zapLog *ZapLog) GetMaxBackup() int       { return zapLog.MaxBackup }
func (zapLog *ZapLog) GetMaxDay() int          { return zapLog.MaxDay }
func (zapLog *ZapLog) GetCompress() bool       { return zapLog.Compress }
func (zapLog *ZapLog) GetInConsole() bool      { return zapLog.InConsole }
func (zapLog *ZapLog) GetExtension() string    { return zapLog.Extension }

// getWriteSync 获取 zapcore.WriteSync
func getWriteSync(zapLog ZapLog, path string) zapcore.WriteSyncer {
	fileWriter := &lumberjack.Logger{
		Filename:   path,             // 日志文件名称
		MaxSize:    zapLog.MaxSize,   // 文件大小限制,单位MB
		MaxBackups: zapLog.MaxBackup, // 最大保留日志文件数量
		MaxAge:     zapLog.MaxDay,    // 日志文件保留天数
		Compress:   zapLog.Compress,  // 是否压缩处理,压缩以后文件为xxxxx.gz
	}

	if zapLog.InConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(fileWriter), zapcore.AddSync(os.Stdout))
	} else {
		return zapcore.AddSync(fileWriter)
	}
}

func NewZapLog(attrs ...ZapLogAttr) (*zap.Logger, error) {
	var (
		err             error
		fs              filesystem.Filesystem
		zapLogger       *zap.Logger
		zapCores        = make([]zapcore.Core, 0, 7)
		zapLoggerConfig = zapcore.EncoderConfig{
			MessageKey:    "message",
			LevelKey:      "logLevel",
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
		encoderTypes = map[ZapLogEncoderType]func(cfg zapcore.EncoderConfig) zapcore.Encoder{
			EncoderTypeJson:    zapcore.NewJSONEncoder,
			EncoderTypeConsole: zapcore.NewConsoleEncoder,
		}

		ins = &ZapLog{
			Level:       zapcore.DebugLevel,
			Path:        ".",
			MaxSize:     1,
			MaxBackup:   5,
			MaxDay:      30,
			Compress:    false,
			InConsole:   false,
			Extension:   ".log",
			EncoderType: EncoderTypeConsole,
		}
	)

	if err = ins.SetAttrs(attrs...); err != nil {
		return nil, err
	}

	if fs = filesystem.NewFile(filesystem.Auto(ins.Path)); !fs.GetExist() {
		if err = fs.Create().GetError(); err != nil {
			return nil, fmt.Errorf("创建日志目录失败：%w", err)
		}
	}

	if ins.Level < zapcore.DebugLevel {
		ins.Level = zapcore.DebugLevel
	}

	if ins.Level > zapcore.FatalLevel {
		ins.Level = zapcore.FatalLevel
	}

	// for logLevel := ins.Level; logLevel <= zapcore.FatalLevel; logLevel++ {
	// 	writer := getWriteSync(*ins, fs.Copy().Join(fmt.Sprintf("%s%s", logLevel.String(), ins.Extension)).GetFullPath())
	// 	zapCores = append(zapCores, zapcore.NewCore(encoderTypes[ins.EncoderType](zapLoggerConfig), writer, logLevel))
	// }

	for _, logLevel := range []zapcore.Level{zapcore.DebugLevel, zapcore.InfoLevel, zapcore.WarnLevel, zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel} {
		if ins.Level >= logLevel {
			writer := getWriteSync(*ins, fs.Copy().Join(fmt.Sprintf("%s%s", time.Now().Format("2006_01_02"), ins.Extension)).GetFullPath())
			zapCores = append(zapCores, zapcore.NewCore(encoderTypes[ins.EncoderType](zapLoggerConfig), writer, logLevel))
		}
	}

	zapLogger = zap.New(zapcore.NewTee(zapCores...))
	if ins.InConsole {
		zapLogger = zapLogger.WithOptions(zap.AddCaller())
	}

	defer func() {
		_ = operation.NewTernary(operation.TrueFn(func() error { return nil }), operation.FalseFn(func() error { return zapLogger.Sync() })).GetByValue(ins.InConsole)
		// if config.InConsole {
		// 	return
		// }
		// err = zapLogger.Sync()
	}()

	return zapLogger, nil
}

func (my *ZapLog) SetAttrs(attrs ...ZapLogAttr) (err error) {
	for idx := range attrs {
		if err = attrs[idx](my); err != nil {
			return err
		}
	}
	return nil
}
