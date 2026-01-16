package initialize

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/aid297/aid/setting"
	"github.com/aid297/aid/str"
	global2 "github.com/aid297/aid/web-site/backend/aid-web-backend/src/global"
)

type ConfigInit struct{}

var Config ConfigInit

func (ConfigInit) Launch(consolePath string) {
	setting.APP.Setting.New(
		setting.ConfigFilename(consolePath),
		setting.EnvName(global2.ENV_CONFIG),
		setting.Content(&global2.CONFIG),
		setting.OnChange(func(v *viper.Viper, e fsnotify.Event) {
			var err error
			global2.LOG.Info(str.APP.Buffer.JoinString("配置文件改变：", e.Name))
			if err = v.Unmarshal(&global2.CONFIG); err != nil {
				global2.LOG.Error(str.APP.Buffer.JoinString("更新配置文件失败：", err.Error()))
			}
		}),
	)
}
