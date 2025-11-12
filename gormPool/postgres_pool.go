package gormPool

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type PostgresPool struct {
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
	sources      map[string]*PostgresConnection
	replicas     map[string]*PostgresConnection
}

var (
	postgresPoolIns   *PostgresPool
	postgresPoolOnce  sync.Once
	PostgresDsnFormat = "host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s"
	PostgresPoolApp   PostgresPool
)

func (*PostgresPool) Once(dbSetting *DbSetting) GormPool { return OncePostgresPool(dbSetting) }

// OncePostgresPool 单例化：postgres链接池
//
//go:fix 推荐使用Once方法
func OncePostgresPool(dbSetting *DbSetting) GormPool {
	postgresPoolOnce.Do(func() {
		postgresPoolIns = &PostgresPool{
			username:     dbSetting.Postgres.Main.Username,
			password:     dbSetting.Postgres.Main.Password,
			host:         dbSetting.Postgres.Main.Host,
			port:         dbSetting.Postgres.Main.Port,
			database:     dbSetting.Postgres.Main.Database,
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
			PostgresDsnFormat,
			dbSetting.Postgres.Main.Host,
			dbSetting.Postgres.Main.Username,
			dbSetting.Postgres.Main.Password,
			dbSetting.Postgres.Main.Database,
			dbSetting.Postgres.Main.Port,
			dbSetting.Postgres.Main.SslMode,
			dbSetting.Postgres.Main.TimeZone,
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
	postgresPoolIns.mainConn, err = gorm.Open(postgres.Open(postgresPoolIns.mainDsn.Content), dbConfig)
	if err != nil {
		panic(fmt.Sprintf("配置数据库失败：%s", err.Error()))
	}

	postgresPoolIns.mainConn = postgresPoolIns.mainConn.Session(&gorm.Session{})
	{
		sqlDb, _ := postgresPoolIns.mainConn.DB()
		sqlDb.SetConnMaxIdleTime(time.Duration(postgresPoolIns.maxIdleTime) * time.Hour)
		sqlDb.SetConnMaxLifetime(time.Duration(postgresPoolIns.maxLifetime) * time.Hour)
		sqlDb.SetMaxIdleConns(postgresPoolIns.maxIdleConns)
		sqlDb.SetMaxOpenConns(postgresPoolIns.maxOpenConns)
	}

	return postgresPoolIns
}

// GetConn 获取主数据库链接
func (my *PostgresPool) GetConn() *gorm.DB {
	my.getRws()
	return my.mainConn
}

// getRws 获取带有读写分离的数据库链接
func (my *PostgresPool) getRws() *gorm.DB {
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
					PostgresDsnFormat,
					item.Host,
					item.Username,
					item.Password,
					item.Database,
					item.Port,
					item.SslMode,
					item.TimeZone,
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
					PostgresDsnFormat,
					item.Host,
					item.Username,
					item.Password,
					item.Database,
					item.Port,
					item.SslMode,
					item.TimeZone,
				),
			})
		}
	}

	if len(sources) > 0 {
		sourceDialectors = make([]gorm.Dialector, len(sources))
		for i := 0; i < len(sources); i++ {
			sourceDialectors[i] = postgres.Open(sources[i].Content)
		}
	}

	if len(replicas) > 0 {
		replicaDialectors = make([]gorm.Dialector, len(replicas))
		for i := 0; i < len(replicas); i++ {
			replicaDialectors[i] = postgres.Open(replicas[i].Content)
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
func (my *PostgresPool) Close() error {
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
