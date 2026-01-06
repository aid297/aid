package main

import (
	"flag"
	"strings"

	"github.com/aid297/aid/web/backend/aid-web-backend-gin/command"
	"github.com/aid297/aid/web/backend/aid-web-backend-gin/global"
	"github.com/aid297/aid/web/backend/aid-web-backend-gin/initialize"

	"github.com/aid297/aid/daemon"
	"github.com/spf13/cast"
)

type consoleArgs struct {
	commandName   string
	configPath    string
	daemonStr     string
	commandParams []string
}

func parseArgs() consoleArgs {
	var (
		originalCommand, commandName, configPath string
		commandParams, originalCommands          []string
		daemonStr                                string
	)
	flag.StringVar(&configPath, "C", "config.yaml", "配置文件路径")
	flag.StringVar(&originalCommand, "M", "", "命令终端参数")
	flag.StringVar(&daemonStr, "D", "false", "是否开启守护进程")
	flag.Parse()

	commandName = ""
	commandParams = make([]string, 0)

	if originalCommand != "" {
		originalCommands = strings.Split(originalCommand, " ")
		commandName = originalCommands[0]
		commandParams = originalCommands[1:]
	}

	return consoleArgs{
		commandName:   commandName,
		configPath:    configPath,
		daemonStr:     daemonStr,
		commandParams: commandParams,
	}
}

func main() {
	var consoleArgs = parseArgs()

	initialize.Config.Launch(consoleArgs.configPath)
	initialize.Zap.Launch()
	initialize.Timezone.Launch()

	launch(consoleArgs)
}

// launch 启动程序
func launch(consoleArgs consoleArgs) {
	if cast.ToBool(consoleArgs.daemonStr) || global.CONFIG.System.Daemon {
		daemon.APP.Main.Once().
			SetTitle("启动程序").
			SetLog(global.CONFIG.Log.Daemon.Dir, global.CONFIG.Log.Daemon.Filename).
			SetLogEnable(true).
			Launch() // 通过守护进程启动
	}

	switch consoleArgs.commandName {
	case "help":
		command.Help.Launch()
	case "web-service":
		command.WebService.Launch()
	default:
		command.WebService.Launch()
	}
}
