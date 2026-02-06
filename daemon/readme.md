### Daemon 使用说明

```go
package main

import (
	"github.com/aid297/aid/daemon"
)

func main() {
	consoleArg := false
	configArg := true

	if consoleArg || configArg {
		daemon.OnceDaemon().
			SetTitle("启动程序").             // 程序标题
			SetLog("logs", "deamon.log"). // 日志文件路径和文件名
			SetLogEnable(true).           // 是否记录日志
			Launch()                      // 通过守护进程启动
	}
}
```