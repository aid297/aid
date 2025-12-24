package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Cors() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "*",                                                                       // 允许所有域名，生产环境请修改为具体域名
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",                                  // 允许的HTTP方法
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",               // 允许的请求头
		AllowCredentials: false,                                                                     // 是否允许发送cookie，如果设置为true，AllowOrigins不能为*
		ExposeHeaders:    "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers", // 暴露的响应头
		MaxAge:           86400,                                                                     // 预检请求的缓存时间（秒）
	})
}
