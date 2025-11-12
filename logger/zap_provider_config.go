package logger

import "go.uber.org/zap/zapcore"

type zapConfig struct {
	Path         string
	PathAbs      bool
	MaxSize      int
	MaxBackup    int
	MaxDay       int
	NeedCompress bool
	InConsole    bool
	Extension    string
	Level        zapcore.Level
	EncoderType  EncoderType
}

// New 实例化：日志配置
func (*zapConfig) New(level zapcore.Level) *zapConfig {
	ins := &zapConfig{
		Path:         ".",
		PathAbs:      false,
		EncoderType:  EncoderTypeConsole,
		Level:        level,
		MaxSize:      1,
		MaxBackup:    5,
		MaxDay:       30,
		NeedCompress: false,
		InConsole:    false,
		Extension:    ".log",
	}

	return ins
}

// SetPath 设置路径
func (my *zapConfig) SetPath(path string) *zapConfig {
	my.Path = path

	return my
}

// SetEncoderType 设置编码类型
func (my *zapConfig) SetEncoderType(encoderType EncoderType) *zapConfig {
	my.EncoderType = encoderType

	return my
}

// SetPathAbs 设置路径是否使用绝对路径
func (my *zapConfig) SetPathAbs(pathAbs bool) *zapConfig {
	my.PathAbs = pathAbs

	return my
}

// SetMaxSize 设置单文件最大存储容量
func (my *zapConfig) SetMaxSize(maxSize int) *zapConfig {
	my.MaxSize = maxSize

	return my
}

// SetMaxBackup 设置最大备份数量
func (my *zapConfig) SetMaxBackup(maxBackup int) *zapConfig {
	my.MaxBackup = maxBackup

	return my
}

// SetMaxDay 设置日志文件最大保存天数
func (my *zapConfig) SetMaxDay(maxDay int) *zapConfig {
	my.MaxDay = maxDay

	return my
}

// SetNeedCompress 设置是否需要压缩
func (my *zapConfig) SetNeedCompress(needCompress bool) *zapConfig {
	my.NeedCompress = needCompress

	return my
}

// SetInConsole 设置是否需要在终端显示
func (my *zapConfig) SetInConsole(InConsole bool) *zapConfig {
	my.InConsole = InConsole

	return my
}

// SetExtension 设置扩展名
func (my *zapConfig) SetExtension(extension string) *zapConfig {
	my.Extension = extension

	return my
}
