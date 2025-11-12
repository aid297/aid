package v1HTTPRoute

import (
	"hr-fiber/api/httpAPI/v1HttpAPI"

	"github.com/gofiber/fiber/v2"
)

type UUIDRout struct{}

var UUID UUIDRout

func (UUIDRout) Register(app *fiber.App) {
	uuidGroup := app.Group("api/v1/uuid")
	{
		uuidGroup.Post("generate", v1HttpAPI.UUID.Generate)
		uuidGroup.Post("versions", v1HttpAPI.UUID.Versions)
	}
}
