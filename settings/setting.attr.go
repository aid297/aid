package settings

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// ******************** 配置文件属性 ********************
type (
	SettingAttributes interface{ Register(settings *Setting) }

	AttrFilename    struct{ filename string }
	AttrEnvName     struct{ envName string }
	AttrDefaultName struct{ name string }
	AttrConfig      struct{ content any }
	AttrOnChange    struct {
		onChange func(v *viper.Viper, e fsnotify.Event)
	}
)

func Filename(filename string) AttrFilename {
	return AttrFilename{filename}
}
func (my AttrFilename) Register(settings *Setting) { settings.configFilename = my.filename }

func EnvName(envName string) AttrEnvName          { return AttrEnvName{envName} }
func (my AttrEnvName) Register(settings *Setting) { settings.envName = my.envName }

func DefaultName(name string) AttrDefaultName         { return AttrDefaultName{name} }
func (my AttrDefaultName) Register(settings *Setting) { settings.defaultName = my.name }

func Content(content any) AttrConfig             { return AttrConfig{content} }
func (my AttrConfig) Register(settings *Setting) { settings.content = my.content }

func OnChange(onChange func(v *viper.Viper, e fsnotify.Event)) AttrOnChange {
	return AttrOnChange{onChange}
}
func (my AttrOnChange) Register(settings *Setting) { settings.onChange = my.onChange }
