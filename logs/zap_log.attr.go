package logs

import "go.uber.org/zap/zapcore"

type ZapLogAttr func(zapLog ZapLogger) (err error)

func Level(level zapcore.Level) ZapLogAttr {
	return func(zapLog ZapLogger) (err error) { zapLog.SetLevel(level); return nil }
}
func Filename(filename string) ZapLogAttr {
	return func(zapLog ZapLogger) (err error) { zapLog.SetFilename(filename); return nil }
}
func MaxSize(maxSize int) ZapLogAttr {
	return func(zapLog ZapLogger) (err error) { zapLog.SetMaxSize(maxSize); return nil }
}
func MaxBackup(maxBackup int) ZapLogAttr {
	return func(zapLog ZapLogger) (err error) { zapLog.SetMaxBackup(maxBackup); return nil }
}
func MaxDay(maxDay int) ZapLogAttr {
	return func(zapLog ZapLogger) (err error) { zapLog.SetMaxDay(maxDay); return nil }
}
func Compress(compress bool) ZapLogAttr {
	return func(zapLog ZapLogger) (err error) { zapLog.SetCompress(compress); return nil }
}
func InConsole(inConsole bool) ZapLogAttr {
	return func(zapLog ZapLogger) (err error) { zapLog.SetInConsole(inConsole); return nil }
}
func EncoderType(encoderType ZapLogEncoderType) ZapLogAttr {
	return func(zapLog ZapLogger) (err error) { zapLog.SetEncoderType(encoderType); return nil }
}
