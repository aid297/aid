package initialize

import (
	"log"
	"strings"

	"hr-fiber/global"

	"github.com/aid297/aid/zapProvider"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapInit struct{}

var Zap ZapInit

func (ZapInit) Launch() {
	var (
		err       error
		zapLevels = map[string]zapcore.Level{
			"debug": zap.DebugLevel,
			"info":  zap.InfoLevel,
			"warn":  zap.WarnLevel,
			"error": zap.ErrorLevel,
			"panic": zap.PanicLevel,
			"fatal": zap.FatalLevel,
		}
	)

	if global.LOGGER, err = zapProvider.ZapProviderApp.New(
		zapProvider.ZapProviderConfig.New(zapLevels[global.SETTING.Log.Zap.Level]).
			SetEncoderType(zapProvider.EncoderType(strings.ToUpper(global.SETTING.Log.Zap.EncoderType))).
			SetExtension(global.SETTING.Log.Zap.Extension).
			SetInConsole(global.SETTING.Log.Zap.InConsole).
			SetMaxSize(global.SETTING.Log.Zap.MaxSize).
			SetMaxDay(global.SETTING.Log.Zap.MaxDay).
			SetPathAbs(global.SETTING.Log.Zap.DirAbs).
			SetPath(global.SETTING.Log.Zap.Dir),
	); err != nil {
		log.Fatalf("【启动日志失败】 %s", err.Error())
	}
}
