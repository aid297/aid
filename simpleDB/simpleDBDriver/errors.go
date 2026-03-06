package simpleDBDriver

import "fmt"

var (
	ErrEmptyKey        = fmt.Errorf("%s 读取错误：key 为空", dbLogTitle)
	ErrKeyNotFound     = fmt.Errorf("%s 读取错误：key 对应数据不存在", dbLogTitle)
	ErrKeyDeleted      = fmt.Errorf("%s 读取错误：key 已经被删除", dbLogTitle)
	ErrDatabaseClosed  = fmt.Errorf("%s 读取错误：数据库已经被关闭", dbLogTitle)
	ErrCorruptedRecord = fmt.Errorf("%s 读取错误：数据记录损坏", dbLogTitle)
	ErrDBPathEmpty     = fmt.Errorf("%s 打开数据库错误：目录为空", dbLogTitle)
	ErrUnkownOperation = fmt.Errorf("%s 未知操作：", dbLogTitle)
)
