package main

import (
	"flag"
	"log"

	"github.com/aid297/aid/daemon"
)

// 主程序
// 通过 go main.go -D=true|false启动
func main() {
	d := flag.Bool("D", false, "daemon")
	flag.Parse()
	log.Printf("参数：D %v", *d)

	if *d {
		daemon.OnceDaemon().SetTitle("daemon-test").SetLogDir(".").SetLogEnable(true).Launch()
		log.Printf("daemon 模式启动")
	} else {
		log.Printf("非 daemon 模式启动")
	}
}
