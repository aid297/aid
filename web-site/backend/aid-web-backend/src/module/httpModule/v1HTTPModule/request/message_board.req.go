package request

import (
	"github.com/aid297/aid/validator/validatorV3"
	"github.com/gin-gonic/gin"
)

type (
	// MessageBoardStoreRequest 表单：保存留言板信息
	MessageBoardStoreRequest struct {
		Content string `json:"content" yaml:"content" toml:"content" swaggertype:"string" v-role:"(required)(min>0)" v-name:"消息内容"`
	}

	// MessageBoardDestroyRequest 表单：删除留言板信息
	MessageBoardDestroyRequest struct {
		ID string `json:"id" yaml:"id" toml:"id" swaggertype:"string" v-role:"(required)(min>0)" v-name:"留言ID"`
	}
)

// Bind 表单绑定：保存留言板信息
func (MessageBoardStoreRequest) Bind(c *gin.Context) (MessageBoardStoreRequest, validatorV3.Checker) {
	return validatorV3.WithGin[MessageBoardStoreRequest](c)
}

// Bind 表单绑定：删除留言板信息
func (MessageBoardDestroyRequest) Bind(c *gin.Context) (MessageBoardDestroyRequest, validatorV3.Checker) {
	return validatorV3.WithGin[MessageBoardDestroyRequest](c)
}
