package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"os"

	"github.com/aid297/aid/operation/operationV2"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/command"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/global"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/initialize"

	"github.com/spf13/cast"

	"github.com/aid297/aid/daemon"
)

type ConsoleArgs struct {
	cmdAPP     string
	configPath string
	daemonStr  string
	cmdParams  []string
}

func parseArgs() ConsoleArgs {
	var (
		originalCmd, cmdAPP, configPath string
		cmdParams, originalCmds         = make([]string, 0), make([]string, 0)
		daemonStr                       string
	)
	flag.StringVar(&configPath, "C", "", "配置文件路径") // 默认配置文件路径：终端命令(C) > 环境变量(AID-BACKEND-CONFIG) > 默认值(config.yaml)
	flag.StringVar(&originalCmd, "M", "", "命令终端参数")
	flag.StringVar(&daemonStr, "D", "", "是否开启守护进程")
	flag.Parse()

	if originalCmd != "" {
		originalCmds = strings.Split(originalCmd, " ")
		cmdAPP = originalCmds[0]
		cmdParams = originalCmds[1:]
	}

	_, configPath = operationV2.NewMultivariate[string]().
		Append(operationV2.MultivariateAttr[string]{Item: configPath, HitFunc: func(_ int, item string) { fmt.Printf("使用终端参数：%s读取配置", item) }}).
		Append(operationV2.MultivariateAttr[string]{Item: os.Getenv("AID-BACKEND-CONFIG"), HitFunc: func(idx int, item string) { fmt.Printf("使用环境变量：%s读取配置", item) }}).
		SetDefault(operationV2.MultivariateAttr[string]{Item: "config.yaml", HitFunc: func(idx int, item string) { fmt.Printf("使用默认参数：%s读取配置", item) }}).
		Finally(func(item string) bool { return item != "" })

	return ConsoleArgs{
		cmdAPP:     cmdAPP,
		configPath: configPath,
		daemonStr:  daemonStr,
		cmdParams:  cmdParams,
	}
}

func main() {
	var consoleArgs = parseArgs()

	initialize.Elements.Config.Launch(consoleArgs.configPath)
	initialize.Elements.Zap.Launch()
	initialize.Elements.Timezone.Launch()

	launch(consoleArgs)
}

// launch 启动程序
func launch(consoleArgs ConsoleArgs) {
	// 守护进程是否开启：终端命令(D) | 配置文件(system.daemon)
	if cast.ToBool(consoleArgs.daemonStr) || global.CONFIG.System.Daemon {
		daemon.GetDaemonOnce().
			SetTitle("启动程序").
			SetLog(global.CONFIG.Log.Daemon.Dir, global.CONFIG.Log.Daemon.Filename).
			SetLogEnable(true).
			Launch() // 通过守护进程启动
	}

	// if (consoleArgs.daemonStr == "" && global.CONFIG.System.Daemon) || (consoleArgs.daemonStr != "" && cast.ToBool(consoleArgs.daemonStr)) {
	// 	daemon.GetDaemonOnce().
	// 		SetTitle("启动程序").
	// 		SetLog(global.CONFIG.Log.Daemon.Dir, global.CONFIG.Log.Daemon.Filename).
	// 		SetLogEnable(true).
	// 		Launch() // 通过守护进程启动
	// }

	switch consoleArgs.cmdAPP {
	case "help":
		command.Elements.Help.Launch()
	case "web-service", "":
		command.Elements.WebService.Launch()
	case "sftp-service":
		command.Elements.SFTPService.Launch()
	default:
		log.Fatalf("启动失败：启动模式不支持：%s", consoleArgs.cmdAPP)
	}
}
