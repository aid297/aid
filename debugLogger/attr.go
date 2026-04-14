package debugLogger

type DebugLoggerAttr func()

func ColorRed() DebugLoggerAttr    { return func() { debugLoggerIns.color = COLOR_RED } }
func ColorGreen() DebugLoggerAttr  { return func() { debugLoggerIns.color = COLOR_GREEN } }
func ColorYellow() DebugLoggerAttr { return func() { debugLoggerIns.color = COLOR_YELLOW } }
func ColorCyan() DebugLoggerAttr   { return func() { debugLoggerIns.color = COLOR_CYAN } }
func ColorBlue() DebugLoggerAttr   { return func() { debugLoggerIns.color = COLOR_BLUE } }
