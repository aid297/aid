package mysqlPool

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	"github.com/aid297/aid/v2/gormPool"
)

type (
	MySQLPool interface {
		GetConn() *gorm.DB
		getRws() *gorm.DB
		Close() error
		SetAttrs(attrs ...MySQLPoolAttr)
		SetUsername(username string)
		SetPassword(password string)
		SetHost(host string)
		SetPort(port uint16)
		SetSources(sources map[string]gormPool.MySQLEndpoint)
		SetReplicas(replicas map[string]gormPool.MySQLEndpoint)
		SetMaxIdleTime(maxIdleTime int)
		SetMaxLifetime(maxLifetime int)
		SetMaxIdleConnections(maxIdleConnections int)
		SetMaxOpenConnections(maxOpenConnections int)
		SetCostomDSNFormat(costomDSNFormat string)
	}

	MySQLPoolImpl struct {
		username           string
		password           string
		host               string
		port               uint16
		database           string
		charset            string
		sources            map[string]gormPool.MySQLEndpoint
		replicas           map[string]gormPool.MySQLEndpoint
		mainDsn            *gormPool.DSN
		mainConn           *gorm.DB
		maxIdleTime        int
		maxLifetime        int
		maxIdleConnections int
		maxOpenConnections int
		costomDSNFormat    string
	}
)

var (
	_              MySQLPool = (*MySQLPoolImpl)(nil)
	mysqlDSNFormat           = "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True"
)

func New(database, charset string, rws bool, attrs ...MySQLPoolAttr) (MySQLPool, error) {
	var (
		err      error
		dbConfig *gorm.Config
		ins      = &MySQLPoolImpl{
			username: "root",
			password: "",
			host:     "localhost",
			port:     3306,
			database: database,
			charset:  charset,
			sources:  nil,
			replicas: nil,
		}
	)

	ins.SetAttrs(attrs...)

	// 配置主库
	ins.mainDsn = &gormPool.DSN{Name: "main", Content: fmt.Sprintf(
		mysqlDSNFormat+ins.costomDSNFormat,
		ins.username, ins.password, ins.host, ins.port, ins.database, ins.charset,
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
	if ins.mainConn, err = gorm.Open(mysql.Open(ins.mainDsn.Content), dbConfig); err != nil {
		panic(fmt.Errorf("配置主库失败：%w", err))
	}

	ins.mainConn = ins.mainConn.Session(&gorm.Session{})
	{
		sqlDB, err := ins.mainConn.DB()
		if err != nil {
			return nil, err
		}
		if ins.maxIdleTime > 0 {
			sqlDB.SetConnMaxIdleTime(time.Duration(ins.maxIdleTime) * time.Hour)
		}
		if ins.maxLifetime > 0 {
			sqlDB.SetConnMaxLifetime(time.Duration(ins.maxLifetime) * time.Hour)
		}
		if ins.maxIdleConnections > 0 {
			sqlDB.SetMaxIdleConns(ins.maxIdleConnections)
		}
		if ins.maxOpenConnections > 0 {
			sqlDB.SetMaxOpenConns(ins.maxOpenConnections)
		}
	}

	return ins, nil
}

func (my *MySQLPoolImpl) SetAttrs(attrs ...MySQLPoolAttr) {
	for _, attr := range attrs {
		attr(my)
	}
}

func (my *MySQLPoolImpl) SetUsername(username string)                          { my.username = username }
func (my *MySQLPoolImpl) SetPassword(password string)                          { my.password = password }
func (my *MySQLPoolImpl) SetHost(host string)                                  { my.host = host }
func (my *MySQLPoolImpl) SetPort(port uint16)                                  { my.port = port }
func (my *MySQLPoolImpl) SetSources(sources map[string]gormPool.MySQLEndpoint) { my.sources = sources }
func (my *MySQLPoolImpl) SetReplicas(replicas map[string]gormPool.MySQLEndpoint) {
	my.replicas = replicas
}
func (my *MySQLPoolImpl) SetMaxIdleTime(maxIdleTime int) { my.maxIdleTime = maxIdleTime }
func (my *MySQLPoolImpl) SetMaxLifetime(maxLifetime int) { my.maxLifetime = maxLifetime }
func (my *MySQLPoolImpl) SetMaxIdleConnections(maxIdleConnections int) {
	my.maxIdleConnections = maxIdleConnections
}
func (my *MySQLPoolImpl) SetMaxOpenConnections(maxOpenConnections int) {
	my.maxOpenConnections = maxOpenConnections
}
func (my *MySQLPoolImpl) SetCostomDSNFormat(costomDSNFormat string) {
	my.costomDSNFormat = costomDSNFormat
}

// GetConn 获取主数据库链接
func (my *MySQLPoolImpl) GetConn() *gorm.DB { my.getRws(); return my.mainConn }

// getRws 获取带有读写分离的数据库链接
func (my *MySQLPoolImpl) getRws() *gorm.DB {
	var (
		err                                  error
		sourceDialectors, replicaDialectores []gorm.Dialector
		sources                              []*gormPool.DSN
		replicas                             []*gormPool.DSN
	)
	// 配置写库
	if len(my.sources) > 0 {
		sources = make([]*gormPool.DSN, 0)
		for idx, item := range my.sources {
			sources = append(sources, &gormPool.DSN{
				Name: idx,
				Content: fmt.Sprintf(
					mysqlDSNFormat+my.costomDSNFormat,
					item.Username,
					item.Password,
					item.Host,
					item.Port,
					my.database,
					my.charset,
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
					mysqlDSNFormat+my.costomDSNFormat,
					item.Username,
					item.Password,
					item.Host,
					item.Port,
					my.database,
					my.charset,
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

	err = my.mainConn.Use(
		dbresolver.Register(dbresolver.Config{
			Sources:           sourceDialectors,          // 写库
			Replicas:          replicaDialectores,        // 读库
			Policy:            dbresolver.RandomPolicy{}, // 策略
			TraceResolverMode: true,
		}).
			SetConnMaxIdleTime(time.Duration(my.maxIdleTime) * time.Hour).
			SetConnMaxLifetime(time.Duration(my.maxLifetime) * time.Hour).
			SetMaxIdleConns(my.maxIdleConnections).
			SetMaxOpenConns(my.maxOpenConnections),
	)
	if err != nil {
		panic(fmt.Errorf("数据库链接错误：%s", err.Error()))
	}

	return my.mainConn
}

// Close 关闭数据库链接
func (my *MySQLPoolImpl) Close() error {
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
