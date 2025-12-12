package main

import (
	"log"

	"github.com/aid297/aid/filesystem/filesystemV2"
	"github.com/aid297/aid/gormPool"
)

type (
	TestTable1 struct {
		ID   int    `gorm:"column:id;primaryKey;autoIncrement"`
		Name string `gorm:"column:name;type:varchar(255);not null;default:'';comment:名称"`
	}
)

func main() {
	pool := gormPool.APP.MySQLPool.Once(gormPool.APP.DBSetting.New(filesystemV2.APP.File.NewByRel("./db.yaml").GetFullPath()))
	conn := pool.GetConn()
	conn.AutoMigrate(&TestTable1{})

	names := []string{"1", "2", "3"}
	testTables := []TestTable1{}
	gormPool.APP.Finder.New(conn.Model(&TestTable1{})).WhenIn(len(names) > 0, "name", names).GetDB().Find(&testTables)

	log.Printf("%v", testTables)
}
