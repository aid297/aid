package initialize

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	zapLog "github.com/aid297/aid/v2/log"
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

	if global.LOG, err = zapLog.NewZapLog(
		zapLog.Level(zapLevels[global.CONFIG.Log.Zap.Level]),
		zapLog.EncoderType(zapLog.ZapLogEncoderType(global.CONFIG.Log.Zap.EncoderType)),
		zapLog.Extension(global.CONFIG.Log.Zap.Extension),
		zapLog.InConsole(global.CONFIG.System.Debug || global.CONFIG.Log.Zap.InConsole),
		zapLog.MaxSize(global.CONFIG.Log.Zap.MaxSize),
		zapLog.MaxDay(global.CONFIG.Log.Zap.MaxDay),
		zapLog.Filename(global.CONFIG.Log.Zap.Dir),
	); err != nil {
		log.Fatalf("【启动日志失败】 %s", err.Error())
	}
}
