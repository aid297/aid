package initialize

import (
	"fmt"
	"hr-fiber/global"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type SettingInit struct{}

var Setting SettingInit

func (SettingInit) Launch(path string) {
	var (
		err         error
		settingPath = path
	)

	if settingEnv := os.Getenv(global.ENV_SETTING); settingEnv != "" {
		settingPath = settingEnv
	}

	global.VIPER = viper.New()
	global.VIPER.SetConfigFile(settingPath)
	global.VIPER.SetConfigType("json")
	if err = global.VIPER.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("读取配置文件失败: %s \n", err))
	}
	global.VIPER.WatchConfig()

	global.VIPER.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件改变:", e.Name)
		if err = global.VIPER.Unmarshal(&global.SETTING); err != nil {
			fmt.Printf("更新配置文件失败：%s\n", err)
		}
	})
	if err = global.VIPER.Unmarshal(&global.SETTING); err != nil {
		panic(err)
	}
}
