package filesystemV4

import "path/filepath"

type (
	PathAttributer interface{ Register(f Filesystemer) }

	AttrPath struct{ path string }
)

func Abs(paths ...string) PathAttributer { return AttrPath{path: getRootPath(filepath.Join(paths...))} }
func Rel(paths ...string) PathAttributer { return AttrPath{path: filepath.Join(paths...)} }

func (my AttrPath) Register(f Filesystemer) { f.SetFullPathForAttr(my.path) }
func (my AttrPath) Joins(paths ...string) PathAttributer {
	my.path = filepath.Join(append([]string{my.path}, paths...)...)
	return my
}
