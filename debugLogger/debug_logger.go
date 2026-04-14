package debugLogger

import (
	"log"
	"os"
	"sync"

	"github.com/spf13/cast"

	"github.com/aid297/aid/str"
)

type DebugLogger struct {
	mu           sync.Mutex
	debugTrigger string
	isDebug      bool
	printLoggers map[bool]*log.Logger
	errorLoggers map[bool]*log.Logger
	color        string
}

const (
	COLOR_RED    = "\033[31m"
	COLOR_GREEN  = "\033[32m"
	COLOR_YELLOW = "\033[33m"
	COLOR_BLUE   = "\033[34m"
	COLOR_CYAN   = "\033[36m"
	COLOR_RESET  = "\033[0m"
)

var (
	debugLoggerOnce sync.Once
	debugLoggerIns  *DebugLogger
)

const ENV_TAG = "AID_DEBUG_LOG"

// Print 输出调试日志
func (my *DebugLogger) Print(v ...any) {
	my.mu.Lock()
	defer my.mu.Unlock()

	if my.debugTrigger != os.Getenv(ENV_TAG) {
		my.debugTrigger = os.Getenv(ENV_TAG)
		my.isDebug = cast.ToBool(my.debugTrigger)
	}

	format, values := my.processArgs(v...)
	if my.color != "" {
		my.printLoggers[my.isDebug].Printf(my.color+format+COLOR_RESET, values...)
	} else {
		my.printLoggers[my.isDebug].Printf(format, values...)
	}

	my.color = ""
}

func (my *DebugLogger) Printf(format string, v ...any) {
	my.printLoggers[my.isDebug].Printf(format, v...)
	my.color = ""
}

// Error 输出错误日志
func (my *DebugLogger) Error(v ...any) {
	if my.debugTrigger != os.Getenv(ENV_TAG) {
		my.debugTrigger = os.Getenv(ENV_TAG)
		my.isDebug = cast.ToBool(my.debugTrigger)
	}

	format, values := my.processArgs(v...)
	my.errorLoggers[my.isDebug].Printf(COLOR_RED+format+COLOR_RESET, values...)

	my.color = ""
}

func (my *DebugLogger) Errorf(format string, v ...any) {
	my.errorLoggers[my.isDebug].Printf(COLOR_RED+format+COLOR_RESET, v...)

	my.color = ""
}

// processArgs 处理日志参数，支持可选的格式字符串
func (my *DebugLogger) processArgs(v ...any) (format string, values []any) {
	if len(v) > 0 {
		format = str.APP.Buffer.JoinString(cast.ToString(v[0]), "\n")
		values = v[1:]
	} else {
		format = "\n"
		values = []any{}
	}

	return
}
