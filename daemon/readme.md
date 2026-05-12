### Daemon 使用说明

```go
package main

import (
	"github.com/aid297/aid/v2/daemon"
)

// 主程序
// 通过 go main.go -D=true|false启动
func main() {
	d := flag.Bool("D", false, "daemon")
	flag.Parse()
	debugLogger.Print("启动参数：D %v", *d)

	if *d {
		daemon.OnceDaemon().SetLogEnable(true).SetLogDir(".").Launch()
		debugLogger.Print("daemon 启动")
	} else {
		debugLogger.Printf("daemon 启动")
	}
}
```