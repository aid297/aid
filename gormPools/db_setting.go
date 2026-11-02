package gormPools

import "github.com/aid297/aid/v2/settings"

type (
	DBSetting struct {
		Common    *Common       `json:"common,omitempty" yaml:"common,omitempty"`
		MySQL     *MySQLSetting `json:"mysql" yaml:"mysql,omitempty"`
		Postgres  *PGSetting    `json:"postgres" yaml:"postgres,omitempty"`
		SQLServer *MSSQLSetting `json:"sql-server" yaml:"sql-server,omitempty"`
		ArSQL     *ArSQLSetting `json:"ar-sql" yaml:"ar-sql,omitempty"`
	}

	Common struct {
		Driver             string `json:"driver" yaml:"driver"`
		MaxOpenConnections int    `json:"max-open-connections" yaml:"max-open-connections"`
		MaxIdleConnections int    `json:"max-idle-connections" yaml:"max-idle-connections"`
		MaxLifetime        int    `json:"max-lifetime" yaml:"max-lifetime"`
		MaxIdleTime        int    `json:"max-idle-timers" yaml:"max-idle-timers"`
	}

	DSN struct {
		Name    string
		Content string
	}

	MySQLSetting struct {
		Database  string                   `json:"database" yaml:"database"`
		Charset   string                   `json:"charset" yaml:"charset"`
		Collation string                   `json:"collation" yaml:"collation"`
		Rws       bool                     `json:"rws" yaml:"rws"`
		Main      *MySQLEndpoint           `json:"main" yaml:"main"`
		Loc       string                   `json:"loc" yaml:"loc"`
		Sources   map[string]MySQLEndpoint `json:"sources" yaml:"sources"`
		Replicas  map[string]MySQLEndpoint `json:"replicas" yaml:"replicas"`
	}

	MySQLEndpoint struct {
		Username string `json:"username" yaml:"username"`
		Password string `json:"password" yaml:"password"`
		Host     string `json:"host" yaml:"host"`
		Port     uint16 `json:"port" yaml:"port"`
	}

	ArSQLSetting struct {
		Database string                    `json:"database" yaml:"database"`
		Rws      bool                      `json:"rws" yaml:"rws"`
		Main     *MySQLEndpoint            `json:"main" yaml:"main"`
		Sources  map[string]*ArSQLEndpoint `json:"sources" yaml:"sources"`
		Replicas map[string]*ArSQLEndpoint `json:"replicas" yaml:"replicas"`
	}

	ArSQLEndpoint struct {
		Username string `json:"username" yaml:"username"`
		Password string `json:"password" yaml:"password"`
		Host     string `json:"host" yaml:"host"`
		Port     uint16 `json:"port" yaml:"port"`
	}

	PGSetting struct {
		Main *PGEndpoint `json:"main" yaml:"main"`
	}

	PGEndpoint struct {
		Username string `json:"username" yaml:"username"`
		Password string `json:"password" yaml:"password"`
		Host     string `json:"host" yaml:"host"`
		Port     uint16 `json:"port" yaml:"port"`
		Database string `json:"database" yaml:"database"`
		TimeZone string `json:"timezone" yaml:"timezone"`
		SSLMode  string `json:"ssl-mode" yaml:"ssl-mode"`
	}

	MSSQLSetting struct {
		Main *PGEndpoint `json:"main" yaml:"main"`
	}

	MSSQLEndpoint struct {
		Username string `json:"username" yaml:"username"`
		Password string `json:"password" yaml:"password"`
		Host     string `json:"host" yaml:"host"`
		Port     uint16 `json:"port" yaml:"port"`
		Database string `json:"database" yaml:"database"`
	}
)

// New 初始化：数据库配置
func (*DBSetting) New(path string) (dbSetting *DBSetting, err error) {
	if _, err = settings.APP.New(settings.Filename(path), settings.Content(dbSetting)); err != nil {
		return nil, err
	}
	return dbSetting, nil
}

func (*DBSetting) ExampleYaml() string {
	return `common:
  driver: "mysql"
  max-open-connections: 100
  max-idle-connections: 20
  max-lifetime: 100
  max-idle-timers: 10
ar-sql:
  database: "cbit_db"
  rws: false
  main:
    username: "yjz"
    password: "123123"
    host: 127.0.0.1
    port: 12344
  sources:
  replicas:
mysql:
  database: "tbl_test"
  charset: "utf8mb4"
  collation: "utf8mb4_general_ci"
  loc: "UTC"
  rws: true
  main:
    username: "root"
    password: "root"
    host: 127.0.0.1
    port: 3308
  sources:
    conn1:
      username: "root"
      password: "root"
      host: 127.0.0.1
      port: 3308
    conn2:
      username: "root"
      password: "root"
      host: 127.0.0.1
      port: 3308
  replicas:
    conn3:
      username: "root"
      password: "root"
      host: 127.0.0.1
      port: 3308
    conn4:
      username: "root"
      password: "root"
      host: 127.0.0.1
      port: 3308
    conn5:
      username: "root"
      password: "root"
      host: 127.0.0.1
      port: 3308
postgres:
  main:
    username: "postgres"
    password: "postgres"
    host: 127.0.0.1
    port: 5432
    database: "tbl_test"
    ssl-mode: "disable"
    timezone: "Asia/Shanghai"
sql-server:
  main:
    username: "admin"
    password: "Admin@1234"
    host: 127.0.0.1
    port: 9930
    database: "tbl_test"`
}
