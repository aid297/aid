package filesystemV33

import (
	"path/filepath"

	"github.com/aid297/aid/operation/operationV2"
)

// ******************** 管理器属性 ******************** //
type (
	DirAttributer interface{ Register(dir *Dir) }

	AttrDirPath  struct{ dirs []string }
	AttrDirIsRel struct{ isRel bool }
)

func (AttrDirPath) Set(vals ...string) AttrDirPath { return AttrDirPath{vals} }
func (my AttrDirPath) Join(dirs ...string) AttrDirPath {
	my.dirs = append(my.dirs, dirs...)
	return my
}
func (my AttrDirPath) Register(dir *Dir) {
	dir.FullPath = operationV2.NewTernary(
		operationV2.TrueFn(func() string { return getRootPath(filepath.Join(my.dirs...)) }),
		operationV2.FalseFn(func() string { return filepath.Join(my.dirs...) }),
	).GetByValue(dir.IsRel)
}

func (AttrDirIsRel) Set(isRel bool) AttrDirIsRel { return AttrDirIsRel{isRel: isRel} }
func (AttrDirIsRel) SetAbs() AttrDirIsRel        { return AttrDirIsRel{isRel: false} }
func (AttrDirIsRel) SetRel() AttrDirIsRel        { return AttrDirIsRel{isRel: true} }
func (my AttrDirIsRel) Register(dir *Dir)        { dir.IsRel = my.isRel }
