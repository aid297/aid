package terminalPrinter

import (
	`fmt`
	`os`
	`time`
)

var (
	TerminalPrint TerminalPrinter = (*TerminalPrinterImpl)(nil)
)

type (
	TerminalPrinter interface {
		New(format string) TerminalPrinter
		Print(color TerminalPrinterColor, v ...any)
		Default(v ...any)
		Info(v ...any)
		Success(v ...any)
		Warning(v ...any)
		Wrong(v ...any)
		Fatal(v ...any)
	}

	TerminalPrinterImpl struct{ format string }

	TerminalPrinterColor string

	TerminalPrinterAttr func(tp TerminalPrinter)
)

const (
	TerminalPrinterColorBlack   TerminalPrinterColor = "\033[30m"
	TerminalPrinterColorRed     TerminalPrinterColor = "\033[31m"
	TerminalPrinterColorGreen   TerminalPrinterColor = "\033[32m"
	TerminalPrinterColorYellow  TerminalPrinterColor = "\033[33m"
	TerminalPrinterColorBlue    TerminalPrinterColor = "\033[34m"
	TerminalPrinterColorMagenta TerminalPrinterColor = "\033[35m"
	TerminalPrinterColorCyan    TerminalPrinterColor = "\033[36m"
	TerminalPrinterColorWhite   TerminalPrinterColor = "\033[37m"
	TerminalPrinterColorReset   TerminalPrinterColor = "\033[0m"
)

func (*TerminalPrinterImpl) New(format string) TerminalPrinter {
	return &TerminalPrinterImpl{format: format}
}

// Print 自定义颜色打印
func (my *TerminalPrinterImpl) Print(color TerminalPrinterColor, v ...any) {
	fmt.Printf("%v「%s」%v\n", color, time.Now().Format(time.DateTime), TerminalPrinterColorReset)
	fmt.Printf(fmt.Sprintf("%v>> %s%v\n\n", color, my.format, TerminalPrinterColorReset), v...)
}

// Default 打印日志行
func (my *TerminalPrinterImpl) Default(v ...any) { my.Print(TerminalPrinterColorReset, v...) }

// Info 打印日志行
func (my *TerminalPrinterImpl) Info(v ...any) { my.Print(TerminalPrinterColorBlue, v...) }

// Success 打印成功
func (my *TerminalPrinterImpl) Success(v ...any) { my.Print(TerminalPrinterColorGreen, v...) }

// Warning 警告
func (my *TerminalPrinterImpl) Warning(v ...any) { my.Print(TerminalPrinterColorYellow, v...) }

// Wrong 错误
func (my *TerminalPrinterImpl) Wrong(v ...any) { my.Print(TerminalPrinterColorRed, v...) }

// Fatal 错误并终止程序
func (my *TerminalPrinterImpl) Fatal(v ...any) { my.Print(TerminalPrinterColorRed, v...); os.Exit(-1) }
