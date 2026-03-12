package kernal

var New app

type app struct{}

// DB 初始化数据库
func (*app) DB(database, table string, attrs ...SchemaAttributer) (*SimpleDB, error) {
	if err := ensureSystemTables(database); err != nil {
		return nil, err
	}
	return newSimpleDB(database, table, attrs...)
}

func (*app) EnsureSystemTables(database string) error {
	return ensureSystemTables(database)
}
