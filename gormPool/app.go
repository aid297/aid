package gormPool

var APP struct {
	MySqlPool     MySqlPool
	PostgresPool  PostgresPool
	SqlServerPool SqlServerPool
	DBSetting     DbSetting
	Finder        Finder
}
