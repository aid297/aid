package route

import (
	"hr-fiber/global"
	"hr-fiber/middleware"
	"hr-fiber/module/httpResponse"
	"hr-fiber/route/httpRoute/v1HTTPRoute"

	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
)

type (
	IndexRout struct{}
)

var (
	Index              IndexRout
	rootRoutMiddleware []any
)

func (IndexRout) Register(app *fiber.App) {
	rootRoutMiddleware = make([]any, 0)

	rootRoutMiddleware = append(rootRoutMiddleware, recover.New(
		recover.Config{
			EnableStackTrace: true,
			StackTraceHandler: func(c *fiber.Ctx, e any) {
				// 自定义 panic 处理逻辑
				global.LOGGER.Error("Panic recovered",
					zap.Any("panic", e),
					zap.String("path", c.Path()),
					zap.String("method", c.Method()),
				)
			},
		},
	))

	if global.SETTING.WebService.Cors {
		rootRoutMiddleware = append(rootRoutMiddleware, middleware.Cors())
	}

	app.Use(rootRoutMiddleware...)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(httpResponse.HealthRes{
			Version: global.SETTING.System.Version,
			Cors:    true,
		})
	})

	v1HTTPRoute.UUID.Register(app)
}
