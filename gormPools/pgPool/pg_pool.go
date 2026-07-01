package pgPool

import (
	"fmt"
	"time"

	"github.com/aid297/aid/v2/gormPools"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type (
	PGPool interface {
		SetAttrs(attrs ...PGPoolAttr)
		SetHost(host string)
		SetPort(port uint16)
		SetSSLMode(sslMode bool)
		SetTimeZone(timezone string)
		SetMaxIdleTime(maxIdleTime int)
		SetMaxLifetime(maxLifetime int)
		SetMaxIdleConns(maxIdleConns int)
		SetMaxOpenConns(maxOpenConns int)
		SetSources(sources map[string]*gormPools.PGEndpoint)
		SetReplicas(replicas map[string]*gormPools.PGEndpoint)
		getRws() *gorm.DB
		GetConn() *gorm.DB
		Close() error
	}

	PGPoolImpl struct {
		username     string
		password     string
		host         string
		port         uint16
		database     string
		maxIdleTime  int
		maxLifetime  int
		maxIdleConns int
		maxOpenConns int
		sslMode      bool
		mainDSN      *gormPools.DSN
		mainConn     *gorm.DB
		sources      map[string]*gormPools.PGEndpoint
		replicas     map[string]*gormPools.PGEndpoint
		timezone     string
	}
)

var (
	_           PGPool = (*PGPoolImpl)(nil)
	pgDSNFormat        = "host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s"
)

func New(database, username, password string, rws bool, attrs ...PGPoolAttr) (PGPool, error) {
	var (
		err error
		ins = &PGPoolImpl{
			username: username,
			password: password,
			host:     "localhost",
			port:     5432,
			database: database,
		}
	)

	ins.SetAttrs(attrs...)

	// 配置主库
	ins.mainDSN = &gormPools.DSN{
		Name: "main",
		Content: fmt.Sprintf(
			pgDSNFormat,
			ins.host, ins.username, ins.password, ins.database, ins.port, ins.sslMode, ins.timezone,
		),
	}

	// 连接主库
	if ins.mainConn, err = gorm.Open(postgres.Open(ins.mainDSN.Content), &gorm.Config{
		PrepareStmt:                              true,  // 预编译
		CreateBatchSize:                          500,   // 批量操作
		DisableForeignKeyConstraintWhenMigrating: true,  // 禁止自动创建外键
		SkipDefaultTransaction:                   false, // 开启自动事务
		QueryFields:                              true,  // 查询字段
		AllowGlobalUpdate:                        false, // 不允许全局修改,必须带有条件
	}); err != nil {
		panic(fmt.Sprintf("配置数据库失败：%s", err.Error()))
	}

	ins.mainConn = ins.mainConn.Session(&gorm.Session{})
	{
		sqlDb, _ := ins.mainConn.DB()
		sqlDb.SetConnMaxIdleTime(time.Duration(ins.maxIdleTime) * time.Hour)
		sqlDb.SetConnMaxLifetime(time.Duration(ins.maxLifetime) * time.Hour)
		sqlDb.SetMaxIdleConns(ins.maxIdleConns)
		sqlDb.SetMaxOpenConns(ins.maxOpenConns)
	}

	return ins, nil
}

func (my *PGPoolImpl) SetAttrs(attrs ...PGPoolAttr) {
	for _, attr := range attrs {
		attr(my)
	}
}

func (my *PGPoolImpl) SetHost(host string)                                   { my.host = host }
func (my *PGPoolImpl) SetPort(port uint16)                                   { my.port = port }
func (my *PGPoolImpl) SetSSLMode(sslMode bool)                               { my.sslMode = sslMode }
func (my *PGPoolImpl) SetTimeZone(timezone string)                           { my.timezone = timezone }
func (my *PGPoolImpl) SetMaxIdleTime(maxIdleTime int)                        { my.maxIdleTime = maxIdleTime }
func (my *PGPoolImpl) SetMaxLifetime(maxLifetime int)                        { my.maxLifetime = maxLifetime }
func (my *PGPoolImpl) SetMaxIdleConns(maxIdleConns int)                      { my.maxIdleConns = maxIdleConns }
func (my *PGPoolImpl) SetMaxOpenConns(maxOpenConns int)                      { my.maxOpenConns = maxOpenConns }
func (my *PGPoolImpl) SetSources(sources map[string]*gormPools.PGEndpoint)   { my.sources = sources }
func (my *PGPoolImpl) SetReplicas(replicas map[string]*gormPools.PGEndpoint) { my.replicas = replicas }

// GetConn 获取主数据库链接
func (my *PGPoolImpl) GetConn() *gorm.DB {
	my.getRws()
	return my.mainConn
}

// getRws 获取带有读写分离的数据库链接
func (my *PGPoolImpl) getRws() *gorm.DB {
	var (
		err                                 error
		sourceDialectors, replicaDialectors []gorm.Dialector
		sources                             []*gormPools.DSN
		replicas                            []*gormPools.DSN
	)
	// 配置写库
	if len(my.sources) > 0 {
		sources = make([]*gormPools.DSN, 0)
		for idx, item := range my.sources {
			sources = append(sources, &gormPools.DSN{
				Name: idx,
				Content: fmt.Sprintf(
					pgDSNFormat,
					item.Host,
					item.Username,
					item.Password,
					item.Database,
					item.Port,
					item.SSLMode,
					item.TimeZone,
				),
			})
		}
	}

	// 配置读库
	if len(my.replicas) > 0 {
		replicas = make([]*gormPools.DSN, 0)
		for idx, item := range my.replicas {
			replicas = append(replicas, &gormPools.DSN{
				Name: idx,
				Content: fmt.Sprintf(
					pgDSNFormat,
					item.Host,
					item.Username,
					item.Password,
					item.Database,
					item.Port,
					item.SSLMode,
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
func (my *PGPoolImpl) Close() error {
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
