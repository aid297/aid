package request

import (
	"github.com/gin-gonic/gin"

	"github.com/aid297/aid/validator/validatorV3"
)

type (
	UUIDGenerateRequest struct {
		Number         uint8       `json:"number" yaml:"number" toml:"number" v-rule:"(required)(uint8)(min>0)(max<256)" v-name:"生成数量"`
		NoSubsTractKey bool        `json:"noSubsTractKey" yaml:"noSubsTractKey" toml:"noSubsTractKey" v-rule:"bool" v-name:"不减去key"`
		IsUpper        bool        `json:"isUpper" yaml:"isUpper" toml:"isUpper" v-rule:"bool" v-name:"是否大写"`
		Version        UUIDVersion `json:"version" yaml:"version" toml:"version" v-rule:"(required)(string)(in:v1,v4,v6,v7)" v-name:"uuid版本"`
	}

	UUIDVersion string
)

const (
	UUIDVersionV1 UUIDVersion = "v1"
	UUIDVersionV4 UUIDVersion = "v4"
	UUIDVersionV6 UUIDVersion = "v6"
	UUIDVersionV7 UUIDVersion = "v7"
)

var UUIDGenerate UUIDGenerateRequest

func (UUIDGenerateRequest) Bind(c *gin.Context) (UUIDGenerateRequest, validatorV3.Checker) {
	return validatorV3.WithGin[UUIDGenerateRequest](c)
}
