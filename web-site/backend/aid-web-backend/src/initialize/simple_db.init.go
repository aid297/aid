package initialize

import (
	"log"

	"github.com/aid297/aid/simpleDB/simpleDBDriver"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/global"
)

type SDB struct{}

func (*SDB) Boot() {
	var err error

	if global.DB, err = simpleDBDriver.New.SimpleDB("./data/aid-db", "message-boards"); err != nil {
		log.Panicf("数据库启动失败：%v", err)
	}

	global.DB.SetConfig(simpleDBDriver.DatabaseConfig{DefaultUUIDVersion: 6})

	if err = global.DB.Configure(simpleDBDriver.TableSchema{
		Columns: []simpleDBDriver.Column{
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
		log.Printf("数据库配置失败：%v", err)
	}

}
