package v1HTTPRequest

import (
	"github.com/aid297/aid/validator"
	"github.com/gofiber/fiber/v2"
)

type (
	UUIDGenerateRequest struct {
		Number         uint16      `json:"number" v-rule:"required;uint16;range=1,100" v-name:"生成数量"`
		NoSubsTractKey bool        `json:"no_subs_tract_key" v-rule:"bool" v-name:"不减去key"`
		IsUpper        bool        `json:"is_upper" v-rule:"bool" v-name:"是否大写"`
		Version        UUIDVersion `json:"version" v-rule:"string;in=v1,v4,v6,v7" v-name:"uuid版本"`
	}
	UUIDVersion string
)

var (
	UUIDVersionV1 UUIDVersion = "v1"
	UUIDVersionV4 UUIDVersion = "v4"
	UUIDVersionV6 UUIDVersion = "v6"
	UUIDVersionV7 UUIDVersion = "v7"

	UUIDGenerate UUIDGenerateRequest
)

func (UUIDGenerateRequest) Bind(c *fiber.Ctx) (UUIDGenerateRequest, error) {
	return validator.WithFiber[UUIDGenerateRequest](c)
}
