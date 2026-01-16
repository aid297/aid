package command

import (
	"log"

	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/global"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/route"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/aid297/aid/str"
)

type WebServiceCommand struct{}

func (WebServiceCommand) Launch() {
	if global.CONFIG.System.Debug {
		gin.SetMode(gin.DebugMode)
		global.LOG.Warn("当前运行在 Debug 模式，Gin 使用 Debug 模式")
	} else {
		gin.SetMode(gin.ReleaseMode)
		global.LOG.Warn("当前运行在 非 Debug 模式，Gin 使用 Release 模式")
	}

	app := gin.Default()
	route.Index.Register(app)

	// 启动web-service服务
	if err := app.Run(str.APP.Buffer.JoinString(":", global.CONFIG.WebService.Port)); err != nil {
		global.LOG.Error("启动web服务", zap.Error(err))
		log.Fatalf("【启动web服务错误】%s", err.Error())
	}
}
