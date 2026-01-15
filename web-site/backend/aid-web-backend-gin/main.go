package main

import (
	"flag"
	"strings"

	"github.com/aid297/aid/web-site/backend/aid-web-backend-gin/command"
	`github.com/aid297/aid/web-site/backend/aid-web-backend-gin/src/global`
	initialize2 `github.com/aid297/aid/web-site/backend/aid-web-backend-gin/src/initialize`

	"github.com/spf13/cast"

	"github.com/aid297/aid/daemon"
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
	flag.StringVar(&daemonStr, "D", "", "是否开启守护进程")
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

	initialize2.Config.Launch(consoleArgs.configPath)
	initialize2.Zap.Launch()
	initialize2.Timezone.Launch()

	launch(consoleArgs)
}

// launch 启动程序
func launch(consoleArgs consoleArgs) {

	if (consoleArgs.daemonStr == "" && global.CONFIG.System.Daemon) || (consoleArgs.daemonStr != "" && cast.ToBool(consoleArgs.daemonStr)) {
		daemon.APP.Main.Once().
			SetTitle("启动程序").
			SetLog(global.CONFIG.Log.Daemon.Dir, global.CONFIG.Log.Daemon.Filename).
			SetLogEnable(true).
			Launch() // 通过守护进程启动
	}

	switch consoleArgs.commandName {
	case "help":
		command.APP.Help.Launch()
	case "web-service":
		command.APP.WebService.Launch()
	case "sftp-service":
		command.APP.SFTP.Launch()
	default:
		command.APP.WebService.Launch()
	}
}
