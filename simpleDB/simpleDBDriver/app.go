package simpleDBDriver

var New app

type app struct{}

// SimpleDB 初始化数据库
func (*app) SimpleDB(database, table string) (*SimpleDB, error) { return newSimpleDB(database, table) }
