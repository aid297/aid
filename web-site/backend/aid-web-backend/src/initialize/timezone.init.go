package initialize

import (
	"log"
	"time"

	"go.uber.org/zap"

	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/global"
)

type TimezoneInit struct{}

var Timezone TimezoneInit

func (TimezoneInit) Launch() {
	if global.CONFIG.System.Timezone != "" {
		if timezoneL, err := time.LoadLocation(global.CONFIG.System.Timezone); err != nil {
			global.LOG.Error("加载时区失败", zap.String("timezone", global.CONFIG.System.Timezone), zap.Error(err))
			log.Fatalf("设置时区失败：%s", err)
		} else {
			time.Local = timezoneL
		}
	}
}
