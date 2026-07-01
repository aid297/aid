package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/spf13/cast"

	"github.com/aid297/aid/v2/daemons"
	"github.com/aid297/aid/v2/debugLogs"
	"github.com/aid297/aid/v2/operations"
	"github.com/aid297/aid/v2/web-site/backend/aid-web-backend/command"
	"github.com/aid297/aid/v2/web-site/backend/aid-web-backend/src/global"
	"github.com/aid297/aid/v2/web-site/backend/aid-web-backend/src/initialize"
)

type ConsoleArgs struct {
	cmdAPP     string
	configPath string
	daemonsStr string
	cmdParams  []string
}

func parseArgs() ConsoleArgs {
	var (
		originalCmd, cmdAPP, configPath string
		cmdParams, originalCmds         = make([]string, 0), make([]string, 0)
		daemonsStr                      string
	)
	flag.StringVar(&configPath, "C", "", "配置文件路径") // 默认配置文件路径：终端命令(C) > 环境变量(AID-BACKEND-CONFIG) > 默认值(config.yaml)
	flag.StringVar(&originalCmd, "M", "", "命令终端参数")
	flag.StringVar(&daemonsStr, "D", "", "是否开启守护进程")
	flag.Parse()

	if originalCmd != "" {
		originalCmds = strings.Split(originalCmd, " ")
		cmdAPP = originalCmds[0]
		cmdParams = originalCmds[1:]
	}

	_, configPath = operations.NewMultivariate[string]().
		Append(operations.MultivariateAttr[string]{Item: configPath, HitFunc: func(_ int, item string) { debugLogs.Print("使用终端参数：%s读取配置", item) }}).
		Append(operations.MultivariateAttr[string]{Item: os.Getenv("AID-BACKEND-CONFIG"), HitFunc: func(idx int, item string) { debugLogs.Print("使用环境变量：%s读取配置", item) }}).
		SetDefault(operations.MultivariateAttr[string]{Item: "config.yaml", HitFunc: func(idx int, item string) { debugLogs.Print("使用默认参数：%s读取配置", item) }}).
		Finally(func(item string) bool { return item != "" })

	return ConsoleArgs{
		cmdAPP:     cmdAPP,
		configPath: configPath,
		daemonsStr: daemonsStr,
		cmdParams:  cmdParams,
	}
}

// @title           Aid Web Backend API
// @version         1.0
// @description     Aid Web Backend API 服务
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html
// @host            localhost:19900
// @BasePath        /api/v1
func main() {
	var consoleArgs = parseArgs()

	initialize.New.Config().Boot(consoleArgs.configPath)
	initialize.New.Zap().Boot()
	initialize.New.Timezone().Boot()
	initialize.New.FileManager().Boot()

	launch(consoleArgs)
}

// launch 启动程序
func launch(consoleArgs ConsoleArgs) {
	// 守护进程是否开启：终端命令(D) | 配置文件(system.daemons)
	if cast.ToBool(consoleArgs.daemonsStr) || global.CONFIG.System.Daemon {
		daemons.OnceDaemon().
			SetTitle("启动程序").
			SetLog(global.CONFIG.Log.Daemon.Dir, global.CONFIG.Log.Daemon.Filename).
			SetLogEnable(true).
			Launch() // 通过守护进程启动
	}

	switch consoleArgs.cmdAPP {
	case "help":
		command.Catalog.Help.Launch()
	case "web-service", "":
		command.Catalog.WebService.Launch()
	case "sftp-service":
		command.Catalog.SFTPService.Launch()
	default:
		log.Fatalf("启动失败：启动模式不支持：%s", consoleArgs.cmdAPP)
	}
}
