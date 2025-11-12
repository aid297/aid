package initialize

import (
	"log"
	"strings"

	"github.com/aid297/aid/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/aid297/aid/web/backend/aid-web-backend-gin/global"
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

	if global.LOG, err = logger.APP.Zap.New(
		logger.APP.ZapConfig.New(zapLevels[global.CONFIG.Log.Zap.Level]).
			SetEncoderType(logger.EncoderType(strings.ToUpper(global.CONFIG.Log.Zap.EncoderType))).
			SetExtension(global.CONFIG.Log.Zap.Extension).
			SetInConsole(global.CONFIG.System.Debug || global.CONFIG.Log.Zap.InConsole).
			SetMaxSize(global.CONFIG.Log.Zap.MaxSize).
			SetMaxDay(global.CONFIG.Log.Zap.MaxDay).
			SetPathAbs(global.CONFIG.Log.Zap.DirAbs).
			SetPath(global.CONFIG.Log.Zap.Dir),
	); err != nil {
		log.Fatalf("【启动日志失败】 %s", err.Error())
	}
}
