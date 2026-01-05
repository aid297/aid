package filesystemV4

import "os"

// ******************** 目录操作属性 ******************** //
type DirOperation struct {
	DirFlag int
	DirMode os.FileMode
}

func (DirOperation) New(attrs ...DirOperationAttributer) DirOperation {
	return DirOperation{DirFlag: 0, DirMode: 0777}.SetAttrs(attrs...)
}

func (my DirOperation) SetAttrs(attrs ...DirOperationAttributer) DirOperation {
	for idx := range attrs {
		attrs[idx].Register(&my)
	}
	return my
}

type (
	DirOperationAttributer interface{ Register(o *DirOperation) }

	AttrDirOperationFlag struct{ flag int }
	AttrDirOperationMode struct{ mode os.FileMode }
)

func (AttrDirOperationFlag) Set(flag int) DirOperationAttributer {
	return AttrDirOperationFlag{flag: flag}
}
func (my AttrDirOperationFlag) Register(o *DirOperation) { o.DirFlag = my.flag }

func (AttrDirOperationMode) Set(mode os.FileMode) DirOperationAttributer {
	return AttrDirOperationMode{mode: mode}
}
func (my AttrDirOperationMode) Register(o *DirOperation) { o.DirMode = my.mode }
