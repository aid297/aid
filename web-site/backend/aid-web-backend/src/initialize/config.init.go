package initialize

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/aid297/aid/setting"
	"github.com/aid297/aid/str"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/global"
)

type ConfigInitialize struct{}

func (*ConfigInitialize) Boot(consolePath string) {
	if _, err := setting.NewSetting(
		setting.Filename(consolePath),
		setting.EnvName(global.ENV_CONFIG),
		setting.Content(&global.CONFIG),
		setting.DefaultName("config.yaml"),
		setting.OnChange(func(v *viper.Viper, e fsnotify.Event) {
			global.LOG.Info(str.APP.Buffer.JoinString("配置文件改变：", e.Name))
			if err := v.Unmarshal(&global.CONFIG); err != nil {
				global.LOG.Error(str.APP.Buffer.JoinString("更新配置文件失败：", err.Error()))
			}
		}),
	); err != nil {
		panic(fmt.Errorf("加载主配置文件错误：%w", err))
	}
}
