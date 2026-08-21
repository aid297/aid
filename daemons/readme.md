### Daemon 使用说明

```go
package main

import (
	"github.com/aid297/aid/v3/daemons"
)

// 主程序
// 通过 go main.go -D=true|false启动
func main() {
	d := flag.Bool("D", false, "daemons")
	flag.Parse()
	debugLogs.Print("启动参数：D %v", *d)

	if *d {
		daemons.OnceDaemon().SetLogEnable(true).SetLogDir(".").Launch()
		debugLogs.Print("daemons 启动")
	} else {
		debugLogs.Printf("daemons 启动")
	}
}
```