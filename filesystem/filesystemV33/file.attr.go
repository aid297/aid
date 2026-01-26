package filesystemV33

import (
	"path/filepath"

	"github.com/aid297/aid/operation/operationV2"
)

// ******************** 管理器属性 ******************** //
type (
	FileAttributer interface{ Register(dir *File) }

	AttrFilePath  struct{ dirs []string }
	AttrFileIsRel struct{ isRel bool }
)

func (AttrFilePath) Set(vals ...string) AttrFilePath { return AttrFilePath{vals} }
func (my AttrFilePath) Join(dirs ...string) AttrFilePath {
	my.dirs = append(my.dirs, dirs...)
	return my
}
func (my AttrFilePath) Register(file *File) {
	file.FullPath = operationV2.NewTernary(
		operationV2.TrueFn(func() string { return getRootPath(filepath.Join(my.dirs...)) }),
		operationV2.FalseFn(func() string { return filepath.Join(my.dirs...) }),
	).GetByValue(file.IsRel)
}

func (my AttrFileIsRel) Set(isRel bool) AttrFileIsRel { return AttrFileIsRel{isRel: isRel} }
func (AttrFileIsRel) SetAbs() AttrFileIsRel           { return AttrFileIsRel{isRel: false} }
func (AttrFileIsRel) SetRel() AttrFileIsRel           { return AttrFileIsRel{isRel: true} }
func (my AttrFileIsRel) Register(dir *File)           { dir.IsRel = my.isRel }
