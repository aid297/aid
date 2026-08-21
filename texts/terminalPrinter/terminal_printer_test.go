package terminalPrinter_test

import (
	"testing"

	"github.com/aid297/aid/v3/texts/terminalPrinter"
)

func Test1(t *testing.T) {
	terminalPrinter.TerminalPrint.New("[测试] %s").Default("default")
	terminalPrinter.TerminalPrint.New("[测试] %s").Print(terminalPrinter.TerminalPrinterColorCyan, "custom color: cyan")
	terminalPrinter.TerminalPrint.New("[测试] %s").Warning("warning")
	terminalPrinter.TerminalPrint.New("[测试] %s").Wrong("wrong")
	terminalPrinter.TerminalPrint.New("[测试] %s").Info("info")
	terminalPrinter.TerminalPrint.New("[测试] %s").Success("success")
	terminalPrinter.TerminalPrint.New("[测试] %s").Fatal("error")
}
