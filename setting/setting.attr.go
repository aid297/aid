package setting

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

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
func (my AttrConfig) Register(setting *Setting) { setting.content = my.content }

func OnChange(onChange func(v *viper.Viper, e fsnotify.Event)) AttrOnChange {
	return AttrOnChange{onChange}
}
func (my AttrOnChange) Register(setting *Setting) { setting.onChange = my.onChange }
