package v1HTTPRoute

import (
	"github.com/gin-gonic/gin"

	"github.com/aid297/aid/web-site/backend/aid-web-backend-gin/api/httpAPI/v1HTTPAPI"
)

type UploadRoute struct{}

func (*UploadRoute) Register(app *gin.RouterGroup) {
	r := app.Group("upload")
	{
		r.POST("/single", v1HTTPAPI.APP.Upload.Single)
		r.POST("/multiple", v1HTTPAPI.APP.Upload.Multiple)
	}
}
