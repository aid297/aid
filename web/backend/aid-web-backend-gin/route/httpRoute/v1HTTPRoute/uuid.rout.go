package v1HTTPRoute

import (
	"github.com/gin-gonic/gin"

	"github.com/aid297/aid/web/backend/aid-web-backend-gin/api/httpAPI/v1HTTPAPI"
)

type (
	UuidRoute struct{}
)

var (
	Uuid UuidRoute
)

func (*UuidRoute) Register(app *gin.RouterGroup) {
	r := app.Group("uuid")
	{
		r.POST("/generate", v1HTTPAPI.APP.UUID.Generate)
		r.POST("/versions", v1HTTPAPI.APP.UUID.Versions)
	}
}
