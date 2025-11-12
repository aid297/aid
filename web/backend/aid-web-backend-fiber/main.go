package main

import (
	"flag"
	"hr-fiber/command"
	"hr-fiber/global"
	"hr-fiber/initialize"
	"strings"

	"github.com/aid297/aid/daemon"
	"github.com/aid297/aid/operation/operationV2"
)

func main() {
	var (
		originalCommand, commandName, settingPath string
		commandParams, originalCommands           []string
		daemon                                    bool
		daemonStr                                 string
	)
	flag.StringVar(&originalCommand, "t", "", "命令终端参数")
	flag.StringVar(&settingPath, "s", "./setting.json", "指定配置文件路径")
	flag.StringVar(&daemonStr, "daemon", "", "是否开启守护进程")
	flag.Parse()

	commandName = ""
	commandParams = make([]string, 0)

	if originalCommand != "" {
		originalCommands = strings.Split(originalCommand, " ")
		commandName = originalCommands[0]
		commandParams = originalCommands[1:]
	}

	initialize.Setting.Launch(settingPath)
	initialize.Zap.Launch()

	daemon = operationV2.NewTernary(operationV2.TrueValue(daemonStr == "true"), operationV2.FalseValue(global.SETTING.System.Daemon)).GetByValue(daemonStr != "")

	launch(commandName, commandParams, originalCommands, daemon)
}

// launch 启动程序
func launch(commandName string, commandParams, originalCommands []string, daemonOpen bool) {
	if daemonOpen {
		daemon.App.New().Launch("启动程序", global.SETTING.Log.Daemon) // 通过守护进程启动
	}

	switch commandName {
	default:
		command.WebService.Launch()
	}
}
