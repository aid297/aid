package simpleDBDriver

var New app

type app struct{}

// SimpleDB 初始化数据库
func (*app) SimpleDB(path string) (*SimpleDB, error) { return newSimpleDB(path) }
