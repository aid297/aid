package v1HTTPService

import (
	"github.com/aid297/aid/simpleDB/simpleDBDriver"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/global"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/module/httpModule/v1HTTPModule/request"
)

// MessageBoardService 服务：留言板
type MessageBoardService struct{}

// List 留言板服务：获取信息列表
func (*MessageBoardService) List() (data []simpleDBDriver.Row, err error) { return global.DB.Find() }

// Store 留言板服务：保存信息
func (*MessageBoardService) Store(form *request.MessageBoardStoreRequest) (err error) {
	rows := []simpleDBDriver.Row{{"content": form.Content}}
	_, err = global.DB.InsertRows(rows)
	return
}

// Destroy 留言板服务：删除信息
func (*MessageBoardService) Destroy(form *request.MessageBoardDestroyRequest) (err error) {
	_, err = global.DB.RemoveByCondition(simpleDBDriver.QueryCondition{
		Field:    "id",
		Operator: simpleDBDriver.QueryOpEQ,
		Value:    form.ID,
	})
	return
}
