package setting

import "github.com/spf13/viper"

var APP app

type app struct{}

func (*app) New(attrs ...SettingAttributes) (*viper.Viper, error) { return newSetting(attrs...) }
