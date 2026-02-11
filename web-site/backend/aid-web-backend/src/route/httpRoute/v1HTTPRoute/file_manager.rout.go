package v1HTTPRoute

import (
	"github.com/gin-gonic/gin"

	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/api/httpAPI/v1HTTPAPI"
)

type FileManagerRoute struct{}

func (*FileManagerRoute) Register(app *gin.RouterGroup) {
	r := app.Group("fileManager")
	{
		r.POST("/upload", v1HTTPAPI.Catalog.FileManager.Upload)     // 上传
		r.POST("/list", v1HTTPAPI.Catalog.FileManager.List)         // 列表
		r.POST("/destroy", v1HTTPAPI.Catalog.FileManager.Destroy)   // 删除
		r.POST("/download", v1HTTPAPI.Catalog.FileManager.Download) // 下载
		r.POST("/zip", v1HTTPAPI.Catalog.FileManager.Zip)           // 压缩
	}
}
