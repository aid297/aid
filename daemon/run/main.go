package main

import (
	"flag"

	"github.com/aid297/aid/daemon"
	"github.com/aid297/aid/debugLogger"
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

	// daemon.OnceDaemon().SetTitle("daemon-test").SetLogDir(".").SetLogEnable(false).LaunchWithCondition(*d, func() {
	// 	debugLogger.Print("启动成功")
	// })
}
