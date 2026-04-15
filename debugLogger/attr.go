package debugLogger

type DebugLoggerAttr func()

func ColorRed() DebugLoggerAttr     { return func() { debugLoggerIns.color = COLOR_RED } }
func ColorGreen() DebugLoggerAttr   { return func() { debugLoggerIns.color = COLOR_GREEN } }
func ColorYellow() DebugLoggerAttr  { return func() { debugLoggerIns.color = COLOR_YELLOW } }
func ColorCyan() DebugLoggerAttr    { return func() { debugLoggerIns.color = COLOR_CYAN } }
func ColorBlue() DebugLoggerAttr    { return func() { debugLoggerIns.color = COLOR_BLUE } }
func ColorDefault() DebugLoggerAttr { return func() { debugLoggerIns.color = COLOR_RESET } }

// SetMode 支持 1(Console)、2(WriteToFile)、3(Console+WriteToFile)
func SetMode(mode OutputMode) DebugLoggerAttr {
	return func() {
		if debugLoggerIns == nil {
			return
		}

		switch mode {
		case ModeConsole, ModeWriteToFile, ModeConsoleAndWriteToFile:
			debugLoggerIns.setMode(mode)
		default:
			debugLoggerIns.setMode(ModeConsole)
		}
	}
}

func Console() DebugLoggerAttr {
	return func() {
		if debugLoggerIns == nil {
			return
		}
		debugLoggerIns.enableMode(ModeConsole)
	}
}

func OnlyConsole() DebugLoggerAttr {
	return SetMode(ModeConsole)
}

func WriteToFile(path string) DebugLoggerAttr {
	return func() {
		if debugLoggerIns == nil {
			return
		}
		debugLoggerIns.configureWriteToFile(path)
	}
}

func OnlyWriteToFile(path string) DebugLoggerAttr {
	return func() {
		if debugLoggerIns == nil {
			return
		}
		SetMode(ModeWriteToFile)()
		debugLoggerIns.configureWriteToFile(path)
	}
}

