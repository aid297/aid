package logs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/aid297/aid/v2/filesystems"
)

// ZapProvider Zap日志服务提供者
type (
	ZapLog struct {
		Level     zapcore.Level
		Filename  string
		MaxSize   int
		MaxBackup int
		MaxDay    int
		Compress  bool
		InConsole bool
		// Extension   string
		EncoderType ZapLogEncoderType
	}

	ZapLogEncoderType string

	ZapLogger interface {
		SetLevel(level zapcore.Level) (err error)
		SetFilename(path string) (err error)
		SetMaxSize(maxSize int) (err error)
		SetMaxBackup(maxBackup int) (err error)
		SetMaxDay(maxDay int) (err error)
		SetCompress(compress bool) (err error)
		SetInConsole(inConsole bool) (err error)
		SetEncoderType(encoderType ZapLogEncoderType) (err error)
		GetLevel() zapcore.Level
		GetFilename() string
		GetMaxSize() int
		GetMaxBackup() int
		GetMaxDay() int
		GetCompress() bool
		GetInConsole() bool
		SetAttrs(attrs ...ZapLogAttr) (err error)
	}
)

const (
	EncoderTypeConsole ZapLogEncoderType = "CONSOLE"
	EncoderTypeJson    ZapLogEncoderType = "JSON"
)

func (my *ZapLog) SetLevel(level zapcore.Level) (err error) { my.Level = level; return nil }
func (my *ZapLog) SetFilename(filename string) (err error) {
	my.Filename = filename
	return nil
}
func (my *ZapLog) SetMaxSize(maxSize int) (err error) { my.MaxSize = maxSize; return nil }
func (my *ZapLog) SetMaxBackup(maxBackup int) (err error) {
	my.MaxBackup = maxBackup
	return nil
}
func (my *ZapLog) SetMaxDay(maxDay int) (err error)      { my.MaxDay = maxDay; return nil }
func (my *ZapLog) SetCompress(compress bool) (err error) { my.Compress = compress; return nil }
func (my *ZapLog) SetInConsole(inConsole bool) (err error) {
	my.InConsole = inConsole
	return nil
}

//	func (zapLog *ZapLog) SetExtension(extension string) (err error) {
//		zapLog.Extension = extension
//		return nil
//	}
func (my *ZapLog) SetEncoderType(encoderType ZapLogEncoderType) (err error) {
	my.EncoderType = encoderType
	return nil
}
func (my *ZapLog) GetLevel() zapcore.Level { return my.Level }
func (my *ZapLog) GetFilename() string     { return my.Filename }
func (my *ZapLog) GetMaxSize() int         { return my.MaxSize }
func (my *ZapLog) GetMaxBackup() int       { return my.MaxBackup }
func (my *ZapLog) GetMaxDay() int          { return my.MaxDay }
func (my *ZapLog) GetCompress() bool       { return my.Compress }
func (my *ZapLog) GetInConsole() bool      { return my.InConsole }

type dailyRotateWriteSyncer struct {
	mu          sync.Mutex
	zapLog      ZapLog
	target      logTarget
	dateLayout  string
	currentDate string
	fileWriter  *lumberjack.Logger
}

type logTarget struct {
	dirPath   string
	filename  string
	extension string
}

func buildLogTarget(path string) logTarget {
	cleanPath := filepath.Clean(path)
	pathExt := filepath.Ext(cleanPath)

	// 完整路径模式：Path("/path/client.log") -> client.2026_01_02.log
	return logTarget{
		dirPath:   filepath.Dir(cleanPath),
		filename:  strings.TrimSuffix(filepath.Base(cleanPath), pathExt),
		extension: pathExt,
	}
}

func newLumberjackWriter(zapLog ZapLog, path string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   path,             // 日志文件名称
		MaxSize:    zapLog.MaxSize,   // 文件大小限制,单位MB
		MaxBackups: zapLog.MaxBackup, // 最大保留日志文件数量
		MaxAge:     zapLog.MaxDay,    // 日志文件保留天数
		Compress:   zapLog.Compress,  // 是否压缩处理,压缩以后文件为xxxxx.gz
	}
}

func (syncer *dailyRotateWriteSyncer) rotateIfNeeded(now time.Time) {
	date := now.Format(syncer.dateLayout)
	if syncer.fileWriter != nil && date == syncer.currentDate {
		return
	}

	syncer.currentDate = date
	filePath := filepath.Join(syncer.target.dirPath, fmt.Sprintf("%s%s", date, syncer.target.extension))
	if syncer.target.filename != "" {
		filePath = filepath.Join(syncer.target.dirPath, fmt.Sprintf("%s.%s%s", syncer.target.filename, date, syncer.target.extension))
	}
	syncer.fileWriter = newLumberjackWriter(syncer.zapLog, filePath)
}

func (syncer *dailyRotateWriteSyncer) Write(p []byte) (n int, err error) {
	syncer.mu.Lock()
	defer syncer.mu.Unlock()

	syncer.rotateIfNeeded(time.Now())
	n, err = syncer.fileWriter.Write(p)
	if err != nil {
		return n, err
	}

	if syncer.zapLog.InConsole {
		_, _ = os.Stdout.Write(p)
	}

	return n, nil
}

func (syncer *dailyRotateWriteSyncer) Sync() error { return nil }

func newDailyRotateWriteSyncer(zapLog ZapLog, target logTarget) zapcore.WriteSyncer {
	return &dailyRotateWriteSyncer{
		zapLog:     zapLog,
		target:     target,
		dateLayout: "2006_01_02",
	}
}

func NewZapLog(attrs ...ZapLogAttr) (*zap.Logger, error) {
	var (
		err             error
		d               filesystems.Filesystem
		target          logTarget
		zapLoggerConfig = zapcore.EncoderConfig{
			MessageKey:    "message",
			LevelKey:      "logLevel",
			TimeKey:       "timers",
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
			Level:     zapcore.DebugLevel,
			Filename:  ".",
			MaxSize:   1,
			MaxBackup: 5,
			MaxDay:    30,
			Compress:  false,
			InConsole: false,
			// Extension:   ".log",
			EncoderType: EncoderTypeConsole,
		}
	)

	if err = ins.SetAttrs(attrs...); err != nil {
		return nil, err
	}

	target = buildLogTarget(ins.Filename)
	if d = filesystems.NewDir(filesystems.Auto(target.dirPath)); !d.GetExist() {
		if err = d.Create().GetError(); err != nil {
			return nil, fmt.Errorf("创建日志目录失败：%w", err)
		}
	}

	if ins.Level < zapcore.DebugLevel {
		ins.Level = zapcore.DebugLevel
	}

	if ins.Level > zapcore.FatalLevel {
		ins.Level = zapcore.FatalLevel
	}

	writer := newDailyRotateWriteSyncer(*ins, target)
	levelEnabler := zap.LevelEnablerFunc(func(logLevel zapcore.Level) bool {
		return logLevel >= ins.Level
	})

	t := encoderTypes[ins.EncoderType](zapLoggerConfig)

	core := zapcore.NewCore(t, writer, levelEnabler)
	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0)), nil
}

func (my *ZapLog) SetAttrs(attrs ...ZapLogAttr) (err error) {
	for idx := range attrs {
		if err = attrs[idx](my); err != nil {
			return err
		}
	}
	return nil
}
