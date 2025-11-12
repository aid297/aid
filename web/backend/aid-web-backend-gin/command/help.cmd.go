package command

import "fmt"

type HelpCommand struct{}

var (
	Help      HelpCommand
	helpTexts = []string{
		`1. 启动参数说明：`,
		`    1.1. -C`,
		`        是否必填：可选`,
		`        含义：配置文件所在路径`,
		`        配置方法：1、设置环境变量："HR_BACKEND_CONFIG" 。2、配置命令行参数：-C xxx.[json|yaml|toml]。3、默认值：config.yaml`,
		`        优先级：命令行参数 > 环境变量 > 默认值(config.yaml)`,
		`    1.2. -M`,
		`        是否必填：可选`,
		`        含义：启动服务模块。 web-service=开启web服务。默认开启 web-service 服务`,
		`    1.3. -D`,
		`        是否必填：可选`,
		`        含义：是否以守护进程方式启动程序`,
		`        配置方法：1、命令行参数：-D=true|false。2、配置文件 system.daemon 字段指定 true|false。3、默认值：false`,
		`        优先级：命令行参数 > 配置文件 > 默认值(false)`,
		`2. 配置参数说明：`,
		`    2.1. system.debug`,
		`        含义：是否开启 Debug 模式。true=开启 Debug 模式，false=开启 Release 模式。Release 模式下绑定路由不会输出到控制台`,
		`    2.2. log.zap.in-console`,
		`        含义：zap 日志是否输出到控制台。true=输出到控制台，false=不输出到控制台`,
		`        备注：system.debug || log.zap.in-console 为 true 时，日志会输出到控制台。也就是说：如果是 Debug 模式，日志一定会输出到控制台。如果是 Release 模式，则根据 log.zap.in-console 配置决定是否输出到控制台`,
	}
)

func (HelpCommand) Launch() {
	for idx := range helpTexts {
		fmt.Println(helpTexts[idx])
	}
}
