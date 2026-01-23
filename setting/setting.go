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
	defaultName    string
	content        any
	onChange       func(v *viper.Viper, e fsnotify.Event)
}

func NewSetting(attrs ...SettingAttributes) (*viper.Viper, error) { return Setting{}.New(attrs...) }

func (Setting) New(attrs ...SettingAttributes) (v *viper.Viper, err error) {
	var (
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

	if configPath == "" && ins.defaultName != "" {
		configPath = ins.defaultName
	}

	if configPath == "" {
		return nil, ErrConfigNotSet
	}

	v = viper.New()
	{
		v.SetConfigFile(configPath)
		v.SetConfigType(filepath.Ext(configPath)[1:])
		if err = v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("%w：%w", ErrConfigFileNotFound, err)
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
			return nil, fmt.Errorf("%w：%w", ErrConfigReadFailed, err)
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
