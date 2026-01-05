package filesystemV4

import "os"

// ******************** 文件操作属性 ******************** //
type FileOperation struct {
	FileFlag int
	FileMode os.FileMode
}

var (
	DefaultCreateMode = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	DefaultReadMode   = os.O_RDWR
)

func (my FileOperation) New(attrs ...FileOperationAttributer) FileOperation {
	return FileOperation{FileFlag: DefaultCreateMode, FileMode: os.FileMode(0777)}.SetAttrs(attrs...)
}

func (my FileOperation) SetAttrs(attrs ...FileOperationAttributer) FileOperation {
	for idx := range attrs {
		attrs[idx].Register(&my)
	}

	return my
}

type (
	FileOperationAttributer interface{ Register(o *FileOperation) }

	AttrFileOperationFlag struct{ flag int }
	AttrFileOperationMode struct{ mode os.FileMode }
)

func (AttrFileOperationFlag) Set(flag int) AttrFileOperationFlag {
	return AttrFileOperationFlag{flag: flag}
}
func (my AttrFileOperationFlag) Register(o *FileOperation) { o.FileFlag = my.flag }

func (AttrFileOperationMode) Set(mode os.FileMode) AttrFileOperationMode {
	return AttrFileOperationMode{mode: mode}
}
func (my AttrFileOperationMode) Register(o *FileOperation) { o.FileMode = my.mode }
