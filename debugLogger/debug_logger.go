package debugLogger

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/spf13/cast"

	"github.com/aid297/aid/str"
)

type DebugLogger struct {
	mu           sync.Mutex
	debugTrigger string
	isDebug      bool
	mode         OutputMode
	printLoggers map[bool]*log.Logger
	errorLoggers map[bool]*log.Logger
	filePrint    map[bool]*log.Logger
	fileError    map[bool]*log.Logger
	filePath     string
	fileHandle   *os.File
	color        string
}

type OutputMode uint8

const (
	ModeConsole OutputMode = 1 << iota
	ModeWriteToFile
	ModeConsoleAndWriteToFile = ModeConsole | ModeWriteToFile
)

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
	ansiColorRegexp = regexp.MustCompile(`\x1b\[[0-9;]*m`)
)

const ENV_TAG = "AID_DEBUG_LOG"

// Print 输出调试日志
func (my *DebugLogger) Print(v ...any) {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.refreshDebugFlag()

	format, values := my.processArgs(v...)
	consoleFormat := format
	if my.color != "" {
		consoleFormat = my.color + format + COLOR_RESET
	}

	my.writePrint(consoleFormat, format, values...)

	my.color = ""
}

func (my *DebugLogger) Printf(format string, v ...any) {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.refreshDebugFlag()

	consoleFormat := format
	if my.color != "" {
		consoleFormat = my.color + format + COLOR_RESET
	}

	my.writePrint(consoleFormat, format, v...)
	my.color = ""
}

// Error 输出错误日志
func (my *DebugLogger) Error(v ...any) {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.refreshDebugFlag()

	format, values := my.processArgs(v...)
	my.writeError(COLOR_RED+format+COLOR_RESET, format, values...)

	my.color = ""
}

func (my *DebugLogger) Errorf(format string, v ...any) {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.refreshDebugFlag()
	my.writeError(COLOR_RED+format+COLOR_RESET, format, v...)

	my.color = ""
}

func (my *DebugLogger) setMode(mode OutputMode) {
	my.mu.Lock()
	defer my.mu.Unlock()
	my.mode = mode
}

func (my *DebugLogger) enableMode(mode OutputMode) {
	my.mu.Lock()
	defer my.mu.Unlock()
	my.mode |= mode
}

func (my *DebugLogger) configureWriteToFile(path string) {
	my.mu.Lock()
	defer my.mu.Unlock()

	if path == "" {
		return
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return
	}

	if my.filePath == absPath && my.fileHandle != nil {
		my.mode |= ModeWriteToFile
		return
	}

	if my.fileHandle != nil {
		_ = my.fileHandle.Close()
	}

	f, err := os.OpenFile(absPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		my.filePath = ""
		my.fileHandle = nil
		my.filePrint[true] = log.New(os.Stderr, "[DEBUG-FILE-ERROR]", log.LstdFlags|log.Lshortfile)
		my.filePrint[false] = log.New(os.Stderr, "[DEBUG-FILE-ERROR]", log.LstdFlags|log.Lshortfile)
		my.fileError[true] = log.New(os.Stderr, "[ERROR-FILE-ERROR]", log.LstdFlags|log.Lshortfile)
		my.fileError[false] = log.New(os.Stderr, "[ERROR-FILE-ERROR]", log.LstdFlags|log.Lshortfile)
		my.mode &^= ModeWriteToFile
		return
	}

	flag := log.LstdFlags | log.Lshortfile
	my.filePath = absPath
	my.fileHandle = f
	my.filePrint[true] = log.New(f, "[DEBUG]", flag)
	my.filePrint[false] = log.New(io.Discard, "", 0)
	my.fileError[true] = log.New(f, "[ERROR]", flag)
	my.fileError[false] = log.New(io.Discard, "", 0)
	my.mode |= ModeWriteToFile
}

func (my *DebugLogger) refreshDebugFlag() {
	if my.debugTrigger != os.Getenv(ENV_TAG) {
		my.debugTrigger = os.Getenv(ENV_TAG)
		my.isDebug = cast.ToBool(my.debugTrigger)
	}
}

func (my *DebugLogger) writePrint(consoleFormat, fileFormat string, values ...any) {
	if my.mode&ModeConsole != 0 {
		my.printLoggers[my.isDebug].Printf(consoleFormat, values...)
	}
	if my.mode&ModeWriteToFile != 0 {
		my.filePrint[my.isDebug].Printf(stripANSI(fileFormat), values...)
	}
}

func (my *DebugLogger) writeError(consoleFormat, fileFormat string, values ...any) {
	if my.mode&ModeConsole != 0 {
		my.errorLoggers[my.isDebug].Printf(consoleFormat, values...)
	}
	if my.mode&ModeWriteToFile != 0 {
		my.fileError[my.isDebug].Printf(stripANSI(fileFormat), values...)
	}
}

func stripANSI(s string) string {
	return ansiColorRegexp.ReplaceAllString(s, "")
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
