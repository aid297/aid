package route

import (
	"net/http"

	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/api/httpAPI"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/global"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/middleware/httpMiddleware"
	v1HTTPMiddleware2 "github.com/aid297/aid/web-site/backend/aid-web-backend/src/middleware/httpMiddleware/v1HTTPMiddleware"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/route/httpRoute/v1HTTPRoute"

	_ "github.com/aid297/aid/web-site/backend/aid-web-backend/docs" // 导入生成的 docs
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type IndexRoute struct{}

var Index IndexRoute

func (*IndexRoute) Register(app *gin.Engine) {
	if global.CONFIG.WebService.Cors {
		app.Use(v1HTTPMiddleware2.Cors())
	}

	// 注册 Swagger 路由
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	app.Use(httpMiddleware.RecoverHandler)

	apiRout := app.Group("api")
	v1Rout := apiRout.Group("v1")
	v1Rout.Use(httpMiddleware.RecoverHandler, v1HTTPMiddleware2.Timeout(120))

	{
		app.Any("/health/json", httpAPI.Catalog.Health.JSON)
		app.Any("/health/yaml", httpAPI.Catalog.Health.YAML)
		app.Any("/health/toml", httpAPI.Catalog.Health.TOML)
		app.Any("/health", httpAPI.Catalog.Health.TOML)
		app.StaticFS("/upload/rezip", http.Dir(global.CONFIG.Rezip.OutDir)) // 静态资源 (压缩包)

		v1HTTPRoute.Catalog.Rezip.Register(v1Rout)
		v1HTTPRoute.Catalog.UUID.Register(v1Rout)
		v1HTTPRoute.Catalog.Upload.Register(v1Rout)

		for idx := range global.CONFIG.WebService.StaticDirs {
			app.Static(global.CONFIG.WebService.StaticDirs[idx].URL, global.CONFIG.WebService.StaticDirs[idx].Dir) // 静态资源路由
		}
	}
}
