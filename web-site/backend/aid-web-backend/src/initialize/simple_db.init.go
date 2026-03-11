package initialize

import (
	"log"

	"github.com/aid297/aid/debugLogger"
	"github.com/aid297/aid/simpleDB/kernal"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/global"
)

type SDB struct{}

func (*SDB) Boot() {
	var err error

	if global.DB, err = kernal.New.DB(
		"./data/aid-db",
		"message-boards",
		kernal.UUIDVersion(6),
	); err != nil {
		log.Panicf("数据库启动失败：%v", err)
	}

	if err = global.DB.Configure(kernal.TableSchema{
		Columns: []kernal.Column{
			{
				Name:          "id",
				Type:          "uuid:v6",
				PrimaryKey:    true,
				AutoIncrement: true,
			},
			{
				Name:     "content",
				Type:     "string",
				Required: true,
				Default:  "",
			},
		},
	}); err != nil {
		debugLogger.Print("数据库配置失败：%v", err)
	}

}
