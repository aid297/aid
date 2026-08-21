package v1HTTPMiddleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/aid297/aid/v3/operations"
)

func Timeout(second time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		keepalive := cast.ToDuration(ctx.GetHeader("keep-alive"))
		timeout.New(
			timeout.WithTimeout(operations.NewTernary(operations.TrueValue(keepalive), operations.FalseValue(second)).GetByValue(keepalive > 0)*time.Second),
			timeout.WithResponse(func(c *gin.Context) {
				c.JSON(http.StatusRequestTimeout, gin.H{"code": 0, "status": http.StatusRequestTimeout, "data": nil, "msg": "请求超时"})
			}),
		)
	}
}
