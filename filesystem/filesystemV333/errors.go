package filesystemV3

import "fmt"

var (
	ErrInit         = fmt.Errorf("获取文件或目录信息失败")
	ErrMissFullPath = fmt.Errorf("文件或目录完整路径不能为空")
	ErrRename       = fmt.Errorf("修改文件名失败")
	ErrRemove       = fmt.Errorf("删除文件失败")
	ErrFileNotExist = fmt.Errorf("文件不存在")
	ErrCreateFile   = fmt.Errorf("创建文件失败")
	ErrWriteFile    = fmt.Errorf("写入文件失败")
	ErrReadFile     = fmt.Errorf("读取文件失败")
	ErrOpenFile     = fmt.Errorf("打开文件失败")
	ErrDirNotExist  = fmt.Errorf("目录不存在")
	ErrCreateDir    = fmt.Errorf("创建目录失败")
	ErrReadDir      = fmt.Errorf("读取目录失败")
)
