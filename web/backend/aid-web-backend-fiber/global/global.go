package global

import (
	"hr-fiber/setting"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	VIPER   *viper.Viper
	LOGGER  *zap.Logger
	SETTING setting.Setting
)
