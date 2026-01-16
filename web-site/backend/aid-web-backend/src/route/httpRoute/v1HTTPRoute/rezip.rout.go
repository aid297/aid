package v1HTTPRoute

import (
	"github.com/gin-gonic/gin"

	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/api/httpAPI/v1HTTPAPI"
)

type RezipRoute struct{}

func (*RezipRoute) Register(app *gin.RouterGroup) {
	r := app.Group("rezip")
	{
		r.POST("/upload", v1HTTPAPI.APP.Rezip.Upload)
	}
}
