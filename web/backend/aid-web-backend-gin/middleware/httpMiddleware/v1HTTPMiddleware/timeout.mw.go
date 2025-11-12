package v1HTTPMiddleware

import (
	"net/http"
	"time"

	"github.com/aid297/aid/operation/operationV2"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func Timeout(second time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		keepalive := cast.ToDuration(ctx.GetHeader("keep-alive"))
		timeout.New(
			timeout.WithTimeout(operationV2.NewTernary(operationV2.TrueValue(keepalive), operationV2.FalseValue(second)).GetByValue(keepalive > 0)*time.Second),
			timeout.WithResponse(func(c *gin.Context) {
				c.JSON(http.StatusRequestTimeout, gin.H{"code": 0, "status": http.StatusRequestTimeout, "data": nil, "msg": "请求超时"})
			}),
		)
	}
}
