package setting

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aid297/aid/str"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// ******************** 配置文件 ********************
// Setting
type Setting struct {
	configFilename string
	envName        string
	content        any
	onChange       func(v *viper.Viper, e fsnotify.Event)
}

func NewSetting(attrs ...SettingAttributes) (v *viper.Viper) { return Setting{}.New(attrs...) }

func (Setting) New(attrs ...SettingAttributes) (v *viper.Viper) {
	var (
		err        error
		configPath string
		configEnv  string
		ins        = Setting{}.SetAttrs(attrs...)
	)

	if configEnv = os.Getenv(ins.envName); configEnv != "" {
		configPath = configEnv
	}

	if ins.configFilename != "" {
		configPath = ins.configFilename
	}

	v = viper.New()
	{
		v.SetConfigFile(configPath)
		v.SetConfigType(filepath.Ext(configPath)[1:])
		if err = v.ReadInConfig(); err != nil {
			panic(fmt.Sprintf("读取配置文件失败: %s \n", err))
		}
		v.WatchConfig()

		v.OnConfigChange(func(e fsnotify.Event) {
			if ins.onChange != nil {
				ins.onChange(v, e)
			} else {
				log.Println(str.APP.Buffer.JoinString("配置文件改变：", e.Name))
				if err = v.Unmarshal(ins.content); err != nil {
					log.Println(str.APP.Buffer.JoinString("更新配置文件失败：", err.Error()))
				}
			}
		})

		if err = v.Unmarshal(ins.content); err != nil {
			panic(err)
		}
	}

	return
}

func (my Setting) SetAttrs(attrs ...SettingAttributes) Setting {
	if len(attrs) == 0 {
		return my
	}

	for _, attr := range attrs {
		attr.Register(&my)
	}

	return my
}
