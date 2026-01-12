package v1HTTPRoute

import (
	"github.com/gin-gonic/gin"

	"github.com/aid297/aid/web-site/backend/aid-web-backend-gin/api/httpAPI/v1HTTPAPI"
)

type FileManagerRoute struct{}

func (*FileManagerRoute) Register(app *gin.RouterGroup) {
	r := app.Group("fileManager")
	{
		r.POST("/upload", v1HTTPAPI.APP.FileManager.Upload)
		r.POST("/list", v1HTTPAPI.APP.FileManager.List)
	}
}
