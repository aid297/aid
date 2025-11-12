package gormPool

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type SqlServerPool struct {
	username     string
	password     string
	host         string
	port         uint16
	database     string
	maxIdleTime  int
	maxLifetime  int
	maxIdleConns int
	maxOpenConns int
	mainDsn      *Dsn
	mainConn     *gorm.DB
	sources      map[string]*SqlServerConnection
	replicas     map[string]*SqlServerConnection
}

var (
	sqlServerPoolIns   *SqlServerPool
	sqlServerPoolOnce  sync.Once
	SqlServerDsnFormat = "sqlserver://%s:%s@%s:?%d?database=%s"
	SqlServerPoolApp   SqlServerPool
)

func (*SqlServerPool) Once(dbSetting *DbSetting) GormPool { return OnceSqlServerPool(dbSetting) }

// OnceSqlServerPool 单例化：sql server连接池
//
//go:fix 推荐使用Once方法
func OnceSqlServerPool(dbSetting *DbSetting) GormPool {
	sqlServerPoolOnce.Do(func() {
		sqlServerPoolIns = &SqlServerPool{
			username:     dbSetting.SqlServer.Main.Username,
			password:     dbSetting.SqlServer.Main.Password,
			host:         dbSetting.SqlServer.Main.Host,
			port:         dbSetting.SqlServer.Main.Port,
			database:     dbSetting.SqlServer.Main.Database,
			maxIdleTime:  dbSetting.Common.MaxIdleTime,
			maxLifetime:  dbSetting.Common.MaxLifetime,
			maxIdleConns: dbSetting.Common.MaxIdleConnections,
			maxOpenConns: dbSetting.Common.MaxOpenConnections,
		}
	})

	var (
		err      error
		dbConfig *gorm.Config
	)

	// 配置主库
	postgresPoolIns.mainDsn = &Dsn{
		Name: "main",
		Content: fmt.Sprintf(
			SqlServerDsnFormat,
			dbSetting.SqlServer.Main.Username,
			dbSetting.SqlServer.Main.Password,
			dbSetting.SqlServer.Main.Host,
			dbSetting.SqlServer.Main.Port,
			dbSetting.SqlServer.Main.Database,
		),
	}

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
	sqlServerPoolIns.mainConn, err = gorm.Open(sqlserver.Open(sqlServerPoolIns.mainDsn.Content), dbConfig)
	if err != nil {
		panic(fmt.Sprintf("配置数据库失败：%s", err.Error()))
	}

	sqlServerPoolIns.mainConn = sqlServerPoolIns.mainConn.Session(&gorm.Session{})
	{
		sqlDb, _ := sqlServerPoolIns.mainConn.DB()
		sqlDb.SetConnMaxIdleTime(time.Duration(sqlServerPoolIns.maxIdleTime) * time.Hour)
		sqlDb.SetConnMaxLifetime(time.Duration(sqlServerPoolIns.maxLifetime) * time.Hour)
		sqlDb.SetMaxIdleConns(sqlServerPoolIns.maxIdleConns)
		sqlDb.SetMaxOpenConns(sqlServerPoolIns.maxOpenConns)
	}

	return sqlServerPoolIns
}

// GetConn 获取主数据库链接
func (my *SqlServerPool) GetConn() *gorm.DB { return my.mainConn }

// getRws 获取带有读写分离的数据库链接
func (my *SqlServerPool) getRws() *gorm.DB {
	var (
		err                                 error
		sourceDialectors, replicaDialectors []gorm.Dialector
		sources                             []*Dsn
		replicas                            []*Dsn
	)
	// 配置写库
	if len(my.sources) > 0 {
		sources = make([]*Dsn, 0)
		for idx, item := range my.sources {
			sources = append(sources, &Dsn{
				Name: idx,
				Content: fmt.Sprintf(
					SqlServerDsnFormat,
					item.Username,
					item.Password,
					item.Host,
					item.Port,
					item.Database,
				),
			})
		}
	}

	// 配置读库
	if len(my.replicas) > 0 {
		replicas = make([]*Dsn, 0)
		for idx, item := range my.replicas {
			replicas = append(replicas, &Dsn{
				Name: idx,
				Content: fmt.Sprintf(
					SqlServerDsnFormat,
					item.Username,
					item.Password,
					item.Host,
					item.Port,
					item.Database,
				),
			})
		}
	}

	if len(sources) > 0 {
		sourceDialectors = make([]gorm.Dialector, len(sources))
		for i := 0; i < len(sources); i++ {
			sourceDialectors[i] = sqlserver.Open(sources[i].Content)
		}
	}

	if len(replicas) > 0 {
		replicaDialectors = make([]gorm.Dialector, len(replicas))
		for i := 0; i < len(replicas); i++ {
			replicaDialectors[i] = sqlserver.Open(replicas[i].Content)
		}
	}

	err = my.mainConn.Use(
		dbresolver.Register(dbresolver.Config{
			Sources:           sourceDialectors,          // 写库
			Replicas:          replicaDialectors,         // 读库
			Policy:            dbresolver.RandomPolicy{}, // 策略
			TraceResolverMode: true,
		}).
			SetConnMaxIdleTime(time.Duration(my.maxIdleTime) * time.Hour).
			SetConnMaxLifetime(time.Duration(my.maxLifetime) * time.Hour).
			SetMaxIdleConns(my.maxIdleConns).
			SetMaxOpenConns(my.maxOpenConns),
	)
	if err != nil {
		panic(fmt.Errorf("数据库链接错误：%s", err.Error()))
	}

	return my.mainConn
}

// Close 关闭数据库链接
func (my *SqlServerPool) Close() error {
	if my.mainConn != nil {
		db, err := my.mainConn.DB()
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
