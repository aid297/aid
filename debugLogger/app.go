package debugLogger

import (
	"io"
	"log"
	"os"
	"sync"
)

var Once once

type once struct{}

// DebugLogger 单例：debug 日志
func (*once) DebugLogger() *DebugLogger {
	debugLoggerOnce.Do(func() {
		debugLoggerIns = &DebugLogger{
			mu:           sync.Mutex{},
			printLoggers: make(map[bool]*log.Logger, 2),
			errorLoggers: make(map[bool]*log.Logger, 2),
		}

		flag := log.LstdFlags | log.Lshortfile
		silenceLogger := log.New(io.Discard, "", 0)

		debugLoggerIns.printLoggers[true] = log.New(os.Stdout, "[DEBUG]", flag)
		debugLoggerIns.printLoggers[false] = silenceLogger
		debugLoggerIns.errorLoggers[true] = log.New(os.Stderr, "[ERROR]", flag)
		debugLoggerIns.errorLoggers[false] = silenceLogger
	})

	return debugLoggerIns
}
