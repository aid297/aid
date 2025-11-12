package gormPool

import (
	"github.com/aid297/aid/honestMan"
)

type (
	DbSetting struct {
		Common    *Common           `yaml:"common,omitempty"`
		MySql     *MySqlSetting     `yaml:"mysql,omitempty"`
		Postgres  *PostgresSetting  `yaml:"postgres,omitempty"`
		SqlServer *SqlServerSetting `yaml:"sqlServer,omitempty"`
		CbitSql   *CbitSqlSetting   `yaml:"cbitSql,omitempty"`
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

	MySqlSetting struct {
		Database  string                      `yaml:"database"`
		Charset   string                      `yaml:"charset"`
		Collation string                      `yaml:"collation"`
		Rws       bool                        `yaml:"rws"`
		Main      *MySqlConnection            `yaml:"main"`
		Sources   map[string]*MySqlConnection `yaml:"sources"`
		Replicas  map[string]*MySqlConnection `yaml:"replicas"`
	}

	MySqlConnection struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     uint16 `yaml:"port"`
	}

	CbitSqlSetting struct {
		Database string                        `yaml:"database"`
		Rws      bool                          `yaml:"rws"`
		Main     *MySqlConnection              `yaml:"main"`
		Sources  map[string]*CbitSqlConnection `yaml:"sources"`
		Replicas map[string]*CbitSqlConnection `yaml:"replicas"`
	}

	CbitSqlConnection struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     uint16 `yaml:"port"`
	}

	PostgresSetting struct {
		Main *PostgresConnection `yaml:"main"`
	}

	PostgresConnection struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     uint16 `yaml:"port"`
		Database string `yaml:"database"`
		TimeZone string `yaml:"timezone"`
		SslMode  string `yaml:"sslmode"`
	}

	SqlServerSetting struct {
		Main *SqlServerConnection `yaml:"main"`
	}

	SqlServerConnection struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     uint16 `yaml:"port"`
		Database string `yaml:"database"`
	}
)

var DbSettingApp DbSetting

// New 初始化：数据库配置
func (*DbSetting) New(path string) *DbSetting {
	var dbSetting *DbSetting = &DbSetting{}

	if err := honestMan.HonestManApp.New(path).LoadYaml(dbSetting); err != nil {
		return nil
	}

	return dbSetting
}

func (*DbSetting) ExampleYaml() string {
	return `common:
  driver: "mysql"
  maxOpenConns: 100
  maxIdleConns: 20
  maxLifetime: 100
  maxIdleTime: 10
cbitSql:
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
sqlServer:
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
