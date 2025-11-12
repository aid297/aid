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
	config         any
	onChange       func(v *viper.Viper, e fsnotify.Event)
}

func (Setting) New(attrs ...SettingAttributes) (v *viper.Viper) {
	var (
		err        error
		configPath string
		configEnv  string
		ins        = Setting{}
	)

	ins = ins.SetAttrs(attrs...)

	if configEnv = os.Getenv(ins.envName); configEnv != "" {
		configPath = configEnv
	}

	if ins.configFilename != "" {
		configPath = ins.configFilename
	}

	v = viper.New()
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
			if err = v.Unmarshal(ins.config); err != nil {
				log.Println(str.APP.Buffer.JoinString("更新配置文件失败：", err.Error()))
			}
		}
	})

	if err = v.Unmarshal(ins.config); err != nil {
		panic(err)
	}

	return
}

func (my Setting) SetAttrs(attrs ...SettingAttributes) Setting {
	for _, attr := range attrs {
		attr.Register(&my)
	}

	return my
}

// ******************** 配置文件属性 ********************
type (
	SettingAttributes interface{ Register(setting *Setting) }

	AttrConfigFilename struct{ configFilename string }
	AttrEnvName        struct{ envName string }
	AttrConfig         struct{ content any }
	AttrOnChange       struct {
		onChange func(v *viper.Viper, e fsnotify.Event)
	}
)

func ConfigFilename(configFilename string) AttrConfigFilename {
	return AttrConfigFilename{configFilename}
}
func (my AttrConfigFilename) Register(setting *Setting) { setting.configFilename = my.configFilename }

func EnvName(envName string) AttrEnvName         { return AttrEnvName{envName} }
func (my AttrEnvName) Register(setting *Setting) { setting.envName = my.envName }

func Content(content any) AttrConfig            { return AttrConfig{content} }
func (my AttrConfig) Register(setting *Setting) { setting.config = my.content }

func OnChange(onChange func(v *viper.Viper, e fsnotify.Event)) AttrOnChange {
	return AttrOnChange{onChange}
}
func (my AttrOnChange) Register(setting *Setting) { setting.onChange = my.onChange }
