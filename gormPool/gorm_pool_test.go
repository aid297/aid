package gormPool

import (
	"log"
	"testing"

	"github.com/aid297/aid/filesystem/filesystemV2"
)

type (
	TestTable1 struct {
		ID   int    `gorm:"column:id;primaryKey;autoIncrement"`
		Name string `gorm:"column:name;type:varchar(255);not null;default:'';comment:名称"`
	}
)

func Test1(t *testing.T) {
	dbSetting, err := APP.DBSetting.New(filesystemV2.APP.File.NewByRel("./db.yaml").GetFullPath())
	if err != nil {
		t.Fatalf("读取配置文件失败：%v", err)
	}
	pool := APP.MySQLPool.Once(dbSetting)
	conn := pool.GetConn()
	conn.AutoMigrate(&TestTable1{})

	names := []string{"1", "2", "3"}
	testTables := []TestTable1{}
	APP.Finder.New(conn.Model(&TestTable1{})).WhenIn(len(names) > 0, "name", names).GetDB().Find(&testTables)

	log.Printf("%v", testTables)
}
