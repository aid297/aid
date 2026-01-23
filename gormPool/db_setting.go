package gormPool

import (
	"github.com/aid297/aid/setting"
)

type (
	DBSetting struct {
		Common    *Common           `yaml:"common,omitempty"`
		MySql     *MySQLSetting     `yaml:"mysql,omitempty"`
		Postgres  *PGSetting        `yaml:"postgres,omitempty"`
		SqlServer *SQLServerSetting `yaml:"sql-server,omitempty"`
		CbitSql   *ArSQLSetting     `yaml:"ar-sql,omitempty"`
	}

	Common struct {
		Driver             string `yaml:"driver"`
		MaxOpenConnections int    `yaml:"maxOpenConns"`
		MaxIdleConnections int    `yaml:"maxIdleConns"`
		MaxLifetime        int    `yaml:"maxLifetime"`
		MaxIdleTime        int    `yaml:"maxIdleTime"`
	}

	Dsn struct {
		Name    string
		Content string
	}

	MySQLSetting struct {
		Database  string                      `yaml:"database"`
		Charset   string                      `yaml:"charset"`
		Collation string                      `yaml:"collation"`
		Rws       bool                        `yaml:"rws"`
		Main      *MySQLConnection            `yaml:"main"`
		Sources   map[string]*MySQLConnection `yaml:"sources"`
		Replicas  map[string]*MySQLConnection `yaml:"replicas"`
	}

	MySQLConnection struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     uint16 `yaml:"port"`
	}

	ArSQLSetting struct {
		Database string                      `yaml:"database"`
		Rws      bool                        `yaml:"rws"`
		Main     *MySQLConnection            `yaml:"main"`
		Sources  map[string]*ArSQLConnection `yaml:"sources"`
		Replicas map[string]*ArSQLConnection `yaml:"replicas"`
	}

	ArSQLConnection struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     uint16 `yaml:"port"`
	}

	PGSetting struct {
		Main *PGConnection `yaml:"main"`
	}

	PGConnection struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     uint16 `yaml:"port"`
		Database string `yaml:"database"`
		TimeZone string `yaml:"timezone"`
		SslMode  string `yaml:"sslmode"`
	}

	SQLServerSetting struct {
		Main *SQLServerConnection `yaml:"main"`
	}

	SQLServerConnection struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     uint16 `yaml:"port"`
		Database string `yaml:"database"`
	}
)

// New 初始化：数据库配置
func (*DBSetting) New(path string) *DBSetting {
	var dbSetting = &DBSetting{}
	setting.NewSetting(setting.Filename(path), setting.Content(dbSetting))
	return dbSetting
}

func (*DBSetting) ExampleYaml() string {
	return `common:
  driver: "mysql"
  maxOpenConns: 100
  maxIdleConns: 20
  maxLifetime: 100
  maxIdleTime: 10
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
    sslmode: "disable"
    timezone: "Asia/Shanghai"
sql-server:
  maxOpenConns: 100
  maxIdleConns: 20
  maxLifetime: 100
  maxIdleTime: 10
  main:
    username: "admin"
    password: "Admin@1234"
    host: 127.0.0.1
    port: 9930
    database: "tbl_test"`
}
