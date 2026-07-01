package initialize

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/aid297/aid/v2/logs"
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

	if global.LOG, err = logs.NewZapLog(
		logs.Level(zapLevels[global.CONFIG.Log.Zap.Level]),
		logs.EncoderType(logs.ZapLogEncoderType(global.CONFIG.Log.Zap.EncoderType)),
		logs.InConsole(global.CONFIG.System.Debug || global.CONFIG.Log.Zap.InConsole),
		logs.MaxSize(global.CONFIG.Log.Zap.MaxSize),
		logs.MaxDay(global.CONFIG.Log.Zap.MaxDay),
		logs.Filename(global.CONFIG.Log.Zap.Filename),
	); err != nil {
		log.Fatalf("【启动日志失败】 %s", err.Error())
	}
}
