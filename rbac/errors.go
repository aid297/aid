package rbac

import "errors"

var (
	ErrDBConnFailed   = errors.New("数据库连接失败")
	ErrCheckRepeat    = errors.New("查重失败")
	ErrRepeatName     = errors.New("名称重复")
	ErrRepeatIdentity = errors.New("名称重复")
	ErrUpdate         = errors.New("编辑失败")
)
