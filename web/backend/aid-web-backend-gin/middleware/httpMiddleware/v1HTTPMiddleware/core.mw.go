package v1HTTPMiddleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Cors 返回一个配置好的跨域中间件
func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Kind", "Authorization", "AccessToken", "X-CSRF-Token", "Token", "X-Token", "X-User-Id", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Kind", "Set-Token", "Set-Expires-At"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
