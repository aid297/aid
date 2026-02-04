package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/aid297/aid/dict/anyMap"
)

var (
	Version   string              // 版本号
	GitCommit string              // Git 提交哈希
	BuildTime string              // 编译时间
	GoVersion = runtime.Version() // Go 运行时版本 (也可以注入，或者直接获取)
)

func main() {
	log.Printf("aaa: %+v\n", anyMap.New(anyMap.Map(map[string]string{
		"error": "1",
		"op":    "",
	})).RemoveEmpty().ToMap())

	if len(os.Args) > 1 && os.Args[1] == "-v" {
		printVersion()
		return
	}
}

func printVersion() {
	fmt.Printf("Application Version: %s\n", Version)
	fmt.Printf("Git Commit: %s\n", GitCommit)
	fmt.Printf("Build Time: %s\n", BuildTime)
	fmt.Printf("Go Version: %s\n", GoVersion)
}
