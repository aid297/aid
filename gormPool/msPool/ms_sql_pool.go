package msPool

import (
	"fmt"
	"time"

	"github.com/aid297/aid/v2/gormPool"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type (
	MSSQLPool interface {
		SetAttrs(attrs ...MSSQLPoolAttr)
		SetHost(host string)
		SetPort(port uint16)
		SetMaxIdleTime(maxIdleTime int)
		SetMaxLifetime(maxLifetime int)
		SetMaxIdleConns(maxIdleConns int)
		SetMaxOpenConns(maxOpenConns int)
		SetSources(sources map[string]*gormPool.MSSQLEndpoint)
		SetReplicas(replicas map[string]*gormPool.MSSQLEndpoint)
		GetConn() *gorm.DB
		getRws() *gorm.DB
		Close() error
	}

	MSSQLPoolImpl struct {
		username     string
		password     string
		host         string
		port         uint16
		database     string
		maxIdleTime  int
		maxLifetime  int
		maxIdleConns int
		maxOpenConns int
		mainDsn      *gormPool.DSN
		mainConn     *gorm.DB
		sources      map[string]*gormPool.MSSQLEndpoint
		replicas     map[string]*gormPool.MSSQLEndpoint
	}
)

var MSSQLDSNFormat = "sqlserver://%s:%s@%s:?%d?database=%s"

func New(database, username, password string, rws bool, attrs ...MSSQLPoolAttr) (MSSQLPool, error) {
	var (
		err error
		ins = &MSSQLPoolImpl{
			username: "sa",
			password: "123456",
			host:     "localhost",
			port:     1433,
		}
	)

	ins.SetAttrs(attrs...)

	// 主库配置
	ins.mainDsn = &gormPool.DSN{
		Name: "main",
		Content: fmt.Sprintf(
			MSSQLDSNFormat,
			ins.username, ins.password, ins.host, ins.port, ins.database,
		),
	}

	// 连接主库
	if ins.mainConn, err = gorm.Open(sqlserver.Open(ins.mainDsn.Content), &gorm.Config{
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

func (my *MSSQLPoolImpl) SetAttrs(attrs ...MSSQLPoolAttr) {
	for _, attr := range attrs {
		attr(my)
	}
}

func (my *MSSQLPoolImpl) SetHost(host string)                                   { my.host = host }
func (my *MSSQLPoolImpl) SetPort(port uint16)                                   { my.port = port }
func (my *MSSQLPoolImpl) SetMaxIdleTime(maxIdleTime int)                        { my.maxIdleTime = maxIdleTime }
func (my *MSSQLPoolImpl) SetMaxLifetime(maxLifetime int)                        { my.maxLifetime = maxLifetime }
func (my *MSSQLPoolImpl) SetMaxIdleConns(maxIdleConns int)                      { my.maxIdleConns = maxIdleConns }
func (my *MSSQLPoolImpl) SetMaxOpenConns(maxOpenConns int)                      { my.maxOpenConns = maxOpenConns }
func (my *MSSQLPoolImpl) SetSources(sources map[string]*gormPool.MSSQLEndpoint) { my.sources = sources }
func (my *MSSQLPoolImpl) SetReplicas(replicas map[string]*gormPool.MSSQLEndpoint) {
	my.replicas = replicas
}

// GetConn 获取主数据库链接
func (my *MSSQLPoolImpl) GetConn() *gorm.DB { return my.mainConn }

// getRws 获取带有读写分离的数据库链接
func (my *MSSQLPoolImpl) getRws() *gorm.DB {
	var (
		err                                 error
		sourceDialectors, replicaDialectors []gorm.Dialector
		sources                             []*gormPool.DSN
		replicas                            []*gormPool.DSN
	)
	// 配置写库
	if len(my.sources) > 0 {
		sources = make([]*gormPool.DSN, 0)
		for idx, item := range my.sources {
			sources = append(sources, &gormPool.DSN{
				Name: idx,
				Content: fmt.Sprintf(
					MSSQLDSNFormat,
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
		replicas = make([]*gormPool.DSN, 0)
		for idx, item := range my.replicas {
			replicas = append(replicas, &gormPool.DSN{
				Name: idx,
				Content: fmt.Sprintf(
					MSSQLDSNFormat,
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
func (my *MSSQLPoolImpl) Close() error {
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
