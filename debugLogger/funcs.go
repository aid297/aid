package debugLogger

import "os"

func Print(v ...any)                 { Once.DebugLogger().Print(v...) }
func Printf(format string, v ...any) { Once.DebugLogger().Print(append([]any{format}, v...)...) }
func Error(v ...any)                 { Once.DebugLogger().Error(v...) }
func Errorf(format string, v ...any) { Once.DebugLogger().Error(append([]any{format}, v...)...) }
func Fatal(v ...any)                 { Once.DebugLogger().Error(v...); os.Exit(1) }
func Fatalf(format string, v ...any) {
	Once.DebugLogger().Error(append([]any{format}, v...)...)
	os.Exit(1)
}
