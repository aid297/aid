package kernal

var New app

type app struct{}

// DB 初始化数据库
func (*app) DB(database, table string) (*SimpleDB, error) { return newSimpleDB(database, table) }
