package v1HTTPRoute

import (
	"github.com/gin-gonic/gin"

	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/api/httpAPI/v1HTTPAPI"
)

type MessageBoardRoute struct{}

func (*MessageBoardRoute) Register(app *gin.RouterGroup) {
	r := app.Group("messageBoard") // 留言板
	{
		r.POST("/list", v1HTTPAPI.New.MessageBoard().List)       // 列表
		r.POST("/store", v1HTTPAPI.New.MessageBoard().Store)     // 保存
		r.POST("/destroy", v1HTTPAPI.New.MessageBoard().Destroy) // 删除
	}
}
