package httpMiddleware

import (
	"runtime/debug"

	"github.com/aid297/aid/web/backend/aid-web-backend-gin/module/httpModule"

	"github.com/gin-gonic/gin"
)

func RecoverHandler(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			switch err := r.(type) {
			case error:
				httpModule.NewInternalServerError().SetErrorf("意外错误：%w", err).WithAccept(c)
			}
			debug.PrintStack()

			c.Abort()
		}
	}()

	c.Next()
}
