package excelV2

import "errors"

var (
	ErrRead              = errors.New("读取文件失败")
	ErrFilenameRequired  = errors.New("文件名不能为空")
	ErrSheetNameRequired = errors.New("工作表名称不能为空")
	ErrCreateSheet       = errors.New("创建工作表失败")
	ErrSheetNotFound     = errors.New("工作表不存在")
	ErrSetFont           = errors.New("设置字体失败")
	ErrSetCell           = errors.New("设置单元格失败")
	ErrSetSheet          = errors.New("设置工作表失败")
	ErrWriteCellFormula  = errors.New("写入单元格失败（公式）")
	ErrWriteCellInt      = errors.New("写入单元格失败（整数）")
	ErrWriteCellFloat    = errors.New("写入单元格失败（浮点）")
	ErrWriteCellBool     = errors.New("写入单元格失败（布尔）")
	ErrWriteCellTime     = errors.New("写入单元格失败（时间）")
	ErrWriteCellAny      = errors.New("写入单元格失败（常规）")
	ErrSave              = errors.New("保存文件失败")
	ErrDownload          = errors.New("下载文件失败")
	ErrColumnNumber      = errors.New("错误的列索引")
)
