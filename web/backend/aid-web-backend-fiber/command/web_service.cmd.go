package command

import (
	"hr-fiber/global"
	"hr-fiber/route"
	"log"

	"github.com/aid297/aid/str"
	"github.com/gofiber/fiber/v2"
	jsonIter "github.com/json-iterator/go"
)

type WebServiceCommand struct{}

var WebService WebServiceCommand

func (WebServiceCommand) Launch() (app *fiber.App) {
	app = fiber.New(fiber.Config{
		// JSONEncoder:  sonic.Marshal, // 换成更快的 JSON 库
		// JSONDecoder:  sonic.Unmarshal,
		JSONEncoder:  jsonIter.Marshal,
		JSONDecoder:  jsonIter.Unmarshal,
		ReadTimeout:  2_000_000_000, // 2s
		WriteTimeout: 2_000_000_000, // 2s
		BodyLimit:    1 << 20,       // 限制请求体 1MB
	})

	route.Index.Register(app) // 注册路由

	if err := app.Listen(str.APP.Buffer.JoinString(":", global.SETTING.WebService.Port)); err != nil {
		log.Fatalf("启动web-service失败：%s\n", err)
	}

	return
}
