package initialize

import (
	"github.com/aid297/aid/web/backend/aid-web-backend-gin/global"

	"github.com/aid297/aid/setting"
	"github.com/aid297/aid/str"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type ConfigInit struct{}

var Config ConfigInit

func (ConfigInit) Launch(consolePath string) {
	setting.APP.Setting.New(
		setting.ConfigFilename(consolePath),
		setting.EnvName(global.ENV_CONFIG),
		setting.Content(&global.CONFIG),
		setting.OnChange(func(v *viper.Viper, e fsnotify.Event) {
			var err error
			global.LOG.Info(str.APP.Buffer.JoinString("配置文件改变：", e.Name))
			if err = v.Unmarshal(&global.CONFIG); err != nil {
				global.LOG.Error(str.APP.Buffer.JoinString("更新配置文件失败：", err.Error()))
			}
		}),
	)
}
