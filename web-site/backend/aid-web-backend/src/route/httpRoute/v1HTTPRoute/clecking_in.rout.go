package v1HTTPRoute

import (
	"github.com/gin-gonic/gin"

	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/api/httpAPI/v1HTTPAPI"
)

type CheckingInRoute struct{}

func (*CheckingInRoute) Register(app *gin.RouterGroup) {
	r := app.Group("checkingIn")
	{
		r.POST("/cal", v1HTTPAPI.New.CheckingIn().Cal) // 计算考勤
	}
}
