package v1HTTPRoute

import (
	"github.com/gin-gonic/gin"

	"github.com/aid297/aid/web/backend/aid-web-backend-gin/api/httpAPI/v1HTTPAPI"
)

type RezipRoute struct{}

var Rezip RezipRoute

func (*RezipRoute) Register(app *gin.RouterGroup) {
	r := app.Group("rezip")
	{
		r.POST("/upload", v1HTTPAPI.APP.Rezip.Upload)
	}
}
