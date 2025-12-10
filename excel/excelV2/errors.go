package excelV2

import "errors"

var (
	ErrRead              = errors.New("读取文件失败")
	ErrFilenameRequired  = errors.New("文件名不能为空")
	ErrSheetNameRequired = errors.New("工作表名称不能为空")
	ErrCreateSheet       = errors.New("创建工作表失败")
	ErrSheetNotFound     = errors.New("工作表不存在")
)
