package str

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cast"
)

type (
	Str struct{ original string }

	TerminalLog struct {
		format string
		enable bool
	}

	TerminalLogColor string
)

var (
	StrApp         Str
	TerminalLogApp TerminalLog
)

const (
	TerminalLogColorBlack   TerminalLogColor = "\033[30m"
	TerminalLogColorRed     TerminalLogColor = "\033[31m"
	TerminalLogColorGreen   TerminalLogColor = "\033[32m"
	TerminalLogColorYellow  TerminalLogColor = "\033[33m"
	TerminalLogColorBlue    TerminalLogColor = "\033[34m"
	TerminalLogColorMagenta TerminalLogColor = "\033[35m"
	TerminalLogColorCyan    TerminalLogColor = "\033[36m"
	TerminalLogColorWhite   TerminalLogColor = "\033[37m"
	TerminalLogColorReset   TerminalLogColor = "\033[0m"
)

func (*Str) New(original string) *Str { return &Str{original: original} }

// NewStr 实例化：字符串
//
//go:fix 推荐使用：New方法
func NewStr(original string) *Str { return &Str{original: original} }

// PadLeftZeros 前置补零
func (my *Str) PadLeftZeros(length int) (string, error) {
	var (
		err error
		res strings.Builder = strings.Builder{}
	)

	if len(my.original) >= length {
		return my.original, nil
	}

	for i := 0; i < length-len(my.original); i++ {
		res.WriteRune('0')
	}

	if _, err = res.WriteString(my.original); err != nil {
		return "", err
	}

	return res.String(), nil
}

// PadRightZeros 后置补零
func (my *Str) PadRightZeros(length int) (string, error) {
	var (
		err error
		res strings.Builder = strings.Builder{}
	)

	if len(my.original) >= length {
		return my.original, nil
	}

	if _, err = res.WriteString(my.original); err != nil {
		return "", err
	}

	for i := 0; i < length-len(my.original); i++ {
		res.WriteRune('0')
	}

	return res.String(), nil
}

// PadRight 后置填充
func (my *Str) PadRight(length int, s string) string {
	my.original += strings.Repeat(s, length-(len(my.original)%length))

	return my.original
}

// PadLeft 前置补充
func (my *Str) PadLeft(length int, s string) string {
	my.original = strings.Repeat(s, length-(len(my.original)%length)) + s

	return my.original
}

// ToLong 超出长度截断
func (my *Str) ToLong(length int, ends ...string) string {
	end := "..."
	if len(ends) > 0 {
		end = ends[0]
	}

	if len(my.original) > length {
		return my.original[:length] + end
	}

	return my.original
}

// New 实例化：控制台日志
func (*TerminalLog) New(format ...string) *TerminalLog {
	var f string
	for _, v := range format {
		f += v
	}

	return &TerminalLog{format: f, enable: cast.ToBool(os.Getenv("AID__STR__TERMINAL_LOG__ENABLE"))}
}

// NewTerminalLog 实例化：控制台日志
func NewTerminalLog(format ...string) *TerminalLog {
	var f string
	for _, v := range format {
		f += v
	}

	return &TerminalLog{format: f, enable: cast.ToBool(os.Getenv("AID__STR__TERMINAL_LOG__ENABLE"))}
}

// Default 打印日志行
func (r *TerminalLog) Default(v ...any) {
	if !r.enable {
		return
	}

	fmt.Printf("%v「%s」%v\n", TerminalLogColorReset, time.Now().Format(time.DateTime), TerminalLogColorReset)
	fmt.Printf(fmt.Sprintf("%v>> %s%v\n\n", TerminalLogColorReset, r.format, TerminalLogColorReset), v...)
}

// Info 打印日志行
func (r *TerminalLog) Info(v ...any) {
	if !r.enable {
		return
	}

	fmt.Printf("%v「%s」%v\n", TerminalLogColorMagenta, time.Now().Format(time.DateTime), TerminalLogColorReset)
	fmt.Printf(fmt.Sprintf("%v>> %s%v\n\n", TerminalLogColorMagenta, r.format, TerminalLogColorReset), v...)
}

// Success 打印成功
func (r *TerminalLog) Success(v ...any) {
	if !r.enable {
		return
	}

	fmt.Printf("%v「%s」%v\n", TerminalLogColorGreen, time.Now().Format(time.DateTime), TerminalLogColorReset)
	fmt.Printf(fmt.Sprintf("%v>> %s%v\n\n", TerminalLogColorGreen, r.format, TerminalLogColorReset), v...)
}

// Wrong 打印错误
func (r *TerminalLog) Wrong(v ...any) {
	if !r.enable {
		return
	}

	fmt.Printf("%v「%s」%v\n", TerminalLogColorRed, time.Now().Format(time.DateTime), TerminalLogColorReset)
	fmt.Printf(fmt.Sprintf("%v>> %s%v\n\n", TerminalLogColorRed, r.format, TerminalLogColorReset), v...)
}

// Error 打印错误并终止程序
func (r *TerminalLog) Error(v ...any) {
	if !r.enable {
		return
	}

	fmt.Printf("%v「%s」%v\n", TerminalLogColorRed, time.Now().Format(time.DateTime), TerminalLogColorReset)
	fmt.Printf(fmt.Sprintf("%v>> %s%v\n\n", TerminalLogColorRed, r.format, TerminalLogColorReset), v...)
	os.Exit(-1)
}
