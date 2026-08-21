package gormPools_test

import (
	"testing"

	"github.com/aid297/aid/v3/debugLogs"
	"github.com/aid297/aid/v3/filesystems"
	"github.com/aid297/aid/v3/gormPools"
	"github.com/aid297/aid/v3/gormPools/mysqlPool"
)

type (
	TestTable1 struct {
		ID   int    `gorm:"column:id;primaryKey;autoIncrement"`
		Name string `gorm:"column:name;type:varchar(255);not null;default:'';comment:名称"`
	}
)

func Test1(t *testing.T) {
	dbSetting, err := gormPools.APP.DBSetting.New(filesystems.NewFile(filesystems.Rel("./db.yaml")).GetFullPath())
	if err != nil {
		t.Fatalf("读取配置文件失败：%v", err)
	}
	pool, err := mysqlPool.New(
		dbSetting.MySQL.Database,
		dbSetting.MySQL.Charset,
		dbSetting.MySQL.Rws,
		mysqlPool.Username(dbSetting.MySQL.Main.Username),
		mysqlPool.Password(dbSetting.MySQL.Main.Password),
		mysqlPool.Host(dbSetting.MySQL.Main.Host),
		mysqlPool.Port(dbSetting.MySQL.Main.Port),
		mysqlPool.Sources(dbSetting.MySQL.Sources),
		mysqlPool.Replicas(dbSetting.MySQL.Replicas),
		mysqlPool.MaxIdleTime(dbSetting.Common.MaxIdleTime),
		mysqlPool.MaxLifetime(dbSetting.Common.MaxLifetime),
		mysqlPool.MaxIdleConnections(dbSetting.Common.MaxIdleConnections),
		mysqlPool.MaxOpenConnections(dbSetting.Common.MaxOpenConnections),
	)
	conn := pool.GetConn()
	conn.AutoMigrate(&TestTable1{})

	names := []string{"1", "2", "3"}
	testTables := []TestTable1{}
	gormPools.NewFinder(conn.Model(&TestTable1{})).WhenIn(len(names) > 0, "name", names).GetDB().Find(&testTables)

	debugLogs.Print("%v", testTables)
}
