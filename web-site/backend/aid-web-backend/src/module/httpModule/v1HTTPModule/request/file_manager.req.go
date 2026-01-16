package request

import (
	"github.com/aid297/aid/validator/validatorV3"
	"github.com/gin-gonic/gin"
)

type (
	FileListRequest struct {
		Path string `json:"path" yaml:"path" toml:"path" v-rule:"string" v-name:"文件路径"`
	}
)

var (
	FileList FileListRequest
)

func (FileListRequest) Bind(c *gin.Context) (FileListRequest, validatorV3.Checker) {
	return validatorV3.WithGin[FileListRequest](c)
}
