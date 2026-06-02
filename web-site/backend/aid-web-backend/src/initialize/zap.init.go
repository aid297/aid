package initialize

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/aid297/aid/v2/logger"
	"github.com/aid297/aid/v2/web-site/backend/aid-web-backend/src/global"
)

type ZapInitialize struct{}

func (*ZapInitialize) Boot() {
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

	if global.LOG, err = logger.NewZapLog(
		logger.Level(zapLevels[global.CONFIG.Log.Zap.Level]),
		logger.EncoderType(logger.ZapLogEncoderType(global.CONFIG.Log.Zap.EncoderType)),
		logger.Extension(global.CONFIG.Log.Zap.Extension),
		logger.InConsole(global.CONFIG.System.Debug || global.CONFIG.Log.Zap.InConsole),
		logger.MaxSize(global.CONFIG.Log.Zap.MaxSize),
		logger.MaxDay(global.CONFIG.Log.Zap.MaxDay),
		logger.Path(global.CONFIG.Log.Zap.Dir),
	); err != nil {
		log.Fatalf("【启动日志失败】 %s", err.Error())
	}
}
