package gormPool

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	"github.com/aid297/aid/operation/operationV2"
)

type MySQLPool struct {
	username  string
	password  string
	host      string
	port      uint16
	database  string
	charset   string
	sources   map[string]*MySQLConnection
	replicas  map[string]*MySQLConnection
	mainDsn   *Dsn
	mainConn  *gorm.DB
	dbSetting *DBSetting
}

var (
	mysqlPoolIns   *MySQLPool
	mysqlPoolOnce  sync.Once
	MySqlDsnFormat = "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local"
	MySqlPoolApp   MySQLPool
)

func (*MySQLPool) Once(dbSetting *DBSetting) *MySQLPool {
	return operationV2.NewTernary(operationV2.TrueValue(OnceMySqlPool(dbSetting)), operationV2.FalseValue(mysqlPoolIns)).GetByValue(dbSetting != nil)
}

// OnceMySqlPool 单例化：mysql链接池
//
//go:fix 推荐使用：Once方法
func OnceMySqlPool(dbSetting *DBSetting) *MySQLPool {
	mysqlPoolOnce.Do(func() {
		mysqlPoolIns = &MySQLPool{
			username:  dbSetting.MySQL.Main.Username,
			password:  dbSetting.MySQL.Main.Password,
			host:      dbSetting.MySQL.Main.Host,
			port:      dbSetting.MySQL.Main.Port,
			database:  dbSetting.MySQL.Database,
			charset:   dbSetting.MySQL.Charset,
			sources:   dbSetting.MySQL.Sources,
			replicas:  dbSetting.MySQL.Replicas,
			dbSetting: dbSetting,
		}
	})

	var (
		err      error
		dbConfig *gorm.Config
	)

	// 配置主库
	mysqlPoolIns.mainDsn = &Dsn{Name: "main", Content: fmt.Sprintf(
		MySqlDsnFormat,
		dbSetting.MySQL.Main.Username,
		dbSetting.MySQL.Main.Password,
		dbSetting.MySQL.Main.Host,
		dbSetting.MySQL.Main.Port,
		dbSetting.MySQL.Database,
		dbSetting.MySQL.Charset,
	)}

	// 数据库配置
	dbConfig = &gorm.Config{
		PrepareStmt:                              true,  // 预编译
		CreateBatchSize:                          500,   // 批量操作
		DisableForeignKeyConstraintWhenMigrating: true,  // 禁止自动创建外键
		SkipDefaultTransaction:                   false, // 开启自动事务
		QueryFields:                              true,  // 查询字段
		AllowGlobalUpdate:                        false, // 不允许全局修改,必须带有条件
	}

	// 配置主库
	if mysqlPoolIns.mainConn, err = gorm.Open(mysql.Open(mysqlPoolIns.mainDsn.Content), dbConfig); err != nil {
		panic(fmt.Errorf("配置主库失败：%w", err))
	}

	mysqlPoolIns.mainConn = mysqlPoolIns.mainConn.Session(&gorm.Session{})
	{
		sqlDb, _ := mysqlPoolIns.mainConn.DB()
		sqlDb.SetConnMaxIdleTime(time.Duration(mysqlPoolIns.dbSetting.Common.MaxIdleTime) * time.Hour)
		sqlDb.SetConnMaxLifetime(time.Duration(mysqlPoolIns.dbSetting.Common.MaxLifetime) * time.Hour)
		sqlDb.SetMaxIdleConns(mysqlPoolIns.dbSetting.Common.MaxIdleConnections)
		sqlDb.SetMaxOpenConns(mysqlPoolIns.dbSetting.Common.MaxOpenConnections)
	}

	return mysqlPoolIns
}

// GetConn 获取主数据库链接
func (*MySQLPool) GetConn() *gorm.DB {
	mysqlPoolIns.getRws()
	return mysqlPoolIns.mainConn
}

// getRws 获取带有读写分离的数据库链接
func (*MySQLPool) getRws() *gorm.DB {
	var (
		err                                  error
		sourceDialectors, replicaDialectores []gorm.Dialector
		sources                              []*Dsn
		replicas                             []*Dsn
	)
	// 配置写库
	if len(mysqlPoolIns.sources) > 0 {
		sources = make([]*Dsn, 0)
		for idx, item := range mysqlPoolIns.sources {
			sources = append(sources, &Dsn{
				Name: idx,
				Content: fmt.Sprintf(
					MySqlDsnFormat,
					item.Username,
					item.Password,
					item.Host,
					item.Port,
					mysqlPoolIns.dbSetting.MySQL.Database,
					mysqlPoolIns.dbSetting.MySQL.Charset,
				),
			})
		}
	}

	// 配置读库
	if len(mysqlPoolIns.replicas) > 0 {
		replicas = make([]*Dsn, 0)
		for idx, item := range mysqlPoolIns.replicas {
			replicas = append(replicas, &Dsn{
				Name: idx,
				Content: fmt.Sprintf(
					MySqlDsnFormat,
					item.Username,
					item.Password,
					item.Host,
					item.Port,
					mysqlPoolIns.dbSetting.MySQL.Database,
					mysqlPoolIns.dbSetting.MySQL.Charset,
				),
			})
		}
	}

	if len(sources) > 0 {
		sourceDialectors = make([]gorm.Dialector, len(sources))
		for i := 0; i < len(sources); i++ {
			sourceDialectors[i] = mysql.Open(sources[i].Content)
		}
	}

	if len(replicas) > 0 {
		replicaDialectores = make([]gorm.Dialector, len(replicas))
		for i := 0; i < len(replicas); i++ {
			replicaDialectores[i] = mysql.Open(replicas[i].Content)
		}
	}

	err = mysqlPoolIns.mainConn.Use(
		dbresolver.Register(dbresolver.Config{
			Sources:           sourceDialectors,          // 写库
			Replicas:          replicaDialectores,        // 读库
			Policy:            dbresolver.RandomPolicy{}, // 策略
			TraceResolverMode: true,
		}).
			SetConnMaxIdleTime(time.Duration(mysqlPoolIns.dbSetting.Common.MaxIdleTime) * time.Hour).
			SetConnMaxLifetime(time.Duration(mysqlPoolIns.dbSetting.Common.MaxLifetime) * time.Hour).
			SetMaxIdleConns(mysqlPoolIns.dbSetting.Common.MaxIdleConnections).
			SetMaxOpenConns(mysqlPoolIns.dbSetting.Common.MaxOpenConnections),
	)
	if err != nil {
		panic(fmt.Errorf("数据库链接错误：%s", err.Error()))
	}

	return mysqlPoolIns.mainConn
}

// Close 关闭数据库链接
func (*MySQLPool) Close() error {
	if mysqlPoolIns.mainConn != nil {
		db, err := mysqlPoolIns.mainConn.DB()
		if err != nil {
			return fmt.Errorf("关闭数据库链接失败：获取数据库链接失败 %s", err.Error())
		}
		err = db.Close()
		if err != nil {
			return fmt.Errorf("关闭数据库连接失败 %s", err.Error())
		}
	}

	return nil
}
