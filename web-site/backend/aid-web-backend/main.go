package main

import (
	"flag"
	"log"
	"strings"

	"github.com/aid297/aid/web-site/backend/aid-web-backend/command"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/global"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/initialize"

	"github.com/spf13/cast"

	"github.com/aid297/aid/daemon"
)

type consoleArgs struct {
	cmdAPP     string
	configPath string
	daemonStr  string
	cmdParams  []string
}

func parseArgs() consoleArgs {
	var (
		originalCmd, cmdAPP, configPath string
		cmdParams, originalCmds         []string
		daemonStr                       string
	)
	flag.StringVar(&configPath, "C", "config.yaml", "配置文件路径")
	flag.StringVar(&originalCmd, "M", "", "命令终端参数")
	flag.StringVar(&daemonStr, "D", "", "是否开启守护进程")
	flag.Parse()

	cmdAPP = ""
	cmdParams = make([]string, 0)

	if originalCmd != "" {
		originalCmds = strings.Split(originalCmd, " ")
		cmdAPP = originalCmds[0]
		cmdParams = originalCmds[1:]
	}

	return consoleArgs{
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
func launch(consoleArgs consoleArgs) {
	if (consoleArgs.daemonStr == "" && global.CONFIG.System.Daemon) || (consoleArgs.daemonStr != "" && cast.ToBool(consoleArgs.daemonStr)) {
		daemon.GetDaemonOnce().
			SetTitle("启动程序").
			SetLog(global.CONFIG.Log.Daemon.Dir, global.CONFIG.Log.Daemon.Filename).
			SetLogEnable(true).
			Launch() // 通过守护进程启动
	}

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
