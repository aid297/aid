package gormPool

var APP struct {
	MySQLPool     MySQLPool
	PGPool        PGPool
	SQLServerPool SQLServerPool
	DBSetting     DBSetting
	Finder        Finder
}
