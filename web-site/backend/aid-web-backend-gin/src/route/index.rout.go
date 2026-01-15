package route

import (
	"net/http"

	"github.com/aid297/aid/web-site/backend/aid-web-backend-gin/src/api/httpAPI"
	`github.com/aid297/aid/web-site/backend/aid-web-backend-gin/src/global`
	`github.com/aid297/aid/web-site/backend/aid-web-backend-gin/src/middleware/httpMiddleware`
	v1HTTPMiddleware2 `github.com/aid297/aid/web-site/backend/aid-web-backend-gin/src/middleware/httpMiddleware/v1HTTPMiddleware`
	`github.com/aid297/aid/web-site/backend/aid-web-backend-gin/src/route/httpRoute/v1HTTPRoute`

	"github.com/gin-gonic/gin"
)

type IndexRoute struct{}

var Index IndexRoute

func (*IndexRoute) Register(app *gin.Engine) {
	if global.CONFIG.WebService.Cors {
		app.Use(v1HTTPMiddleware2.Cors())
	}
	app.Use(httpMiddleware.RecoverHandler)

	apiRout := app.Group("api")
	v1Rout := apiRout.Group("v1")
	v1Rout.Use(httpMiddleware.RecoverHandler, v1HTTPMiddleware2.Timeout(120))

	{
		app.Any("/health/json", httpAPI.APP.Health.JSON)
		app.Any("/health/yaml", httpAPI.APP.Health.YAML)
		app.Any("/health/toml", httpAPI.APP.Health.TOML)
		app.Any("/health", httpAPI.APP.Health.TOML)
		app.StaticFS("/upload/rezip", http.Dir(global.CONFIG.Rezip.OutDir)) // 静态资源 (压缩包)

		v1HTTPRoute.APP.Rezip.Register(v1Rout)
		v1HTTPRoute.APP.UUID.Register(v1Rout)
		v1HTTPRoute.APP.Upload.Register(v1Rout)

		for idx := range global.CONFIG.WebService.StaticDirs {
			app.Static(global.CONFIG.WebService.StaticDirs[idx].URL, global.CONFIG.WebService.StaticDirs[idx].Dir) // 静态资源路由
		}
	}
}
