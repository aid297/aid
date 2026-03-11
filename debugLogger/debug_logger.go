package debugLogger

import (
	"log"
	"os"
	"sync"

	"github.com/aid297/aid/str"
	"github.com/spf13/cast"
)

type DebugLogger struct {
	mu           sync.Mutex
	debugTrigger string
	isDebug      bool
	printLoggers map[bool]*log.Logger
	errorLoggers map[bool]*log.Logger
}

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
	my.printLoggers[my.isDebug].Printf(format, values...)
}

// Error 输出错误日志
func (my *DebugLogger) Error(v ...any) {
	if my.debugTrigger != os.Getenv(ENV_TAG) {
		my.debugTrigger = os.Getenv(ENV_TAG)
		my.isDebug = cast.ToBool(my.debugTrigger)
	}

	format, values := my.processArgs(v...)
	my.errorLoggers[my.isDebug].Printf(format, values...)
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
