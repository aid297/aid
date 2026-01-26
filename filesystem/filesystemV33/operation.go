package filesystemV33

import "os"

// ******************** 目录操作属性 ******************** //
type DirOperation struct {
	DirFlag int
	DirMode os.FileMode
}

func (my *DirOperation) SetAttrs(attrs ...DirOperationAttributer) *DirOperation {
	for idx := range attrs {
		attrs[idx].Register(my)
	}

	return my
}

type (
	DirOperationAttributer interface{ Register(o *DirOperation) }

	AttrDirFlag struct{ flag int }
	AttrDirMode struct{ mode os.FileMode }
)

func DirFlag(flag int) DirOperationAttributer   { return AttrDirFlag{flag: flag} }
func (my AttrDirFlag) Register(o *DirOperation) { o.DirFlag = my.flag }

func DirMode(mode os.FileMode) DirOperationAttributer { return AttrDirMode{mode: mode} }
func (my AttrDirMode) Register(o *DirOperation)       { o.DirMode = my.mode }

// ******************** 文件操作属性 ******************** //
type FileOperation struct {
	FileFlag int
	FileMode os.FileMode
}

func (my *FileOperation) SetAttrs(attrs ...FileOperationAttributer) *FileOperation {
	for idx := range attrs {
		attrs[idx].Register(my)
	}

	return my
}

type (
	FileOperationAttributer interface{ Register(o *FileOperation) }

	AttrFileFlag struct{ flag int }
	AttrFileMode struct{ mode os.FileMode }
)

func FileFlag(flag int) AttrFileFlag              { return AttrFileFlag{flag: flag} }
func (my AttrFileFlag) Register(o *FileOperation) { o.FileFlag = my.flag }

func FileMode(mode os.FileMode) AttrFileMode      { return AttrFileMode{mode: mode} }
func (my AttrFileMode) Register(o *FileOperation) { o.FileMode = my.mode }
