package filesystemV4

import "path/filepath"

type (
	PathAttributer interface {
		Joins(paths ...string) PathAttributer
		Register(f IFilesystem)
		GetPath() string
	}

	AttrPath struct{ path string }
)

func Rel(paths ...string) PathAttributer { return AttrPath{path: getRootPath(filepath.Join(paths...))} }
func Abs(paths ...string) PathAttributer { return AttrPath{path: filepath.Join(paths...)} }

func (my AttrPath) Register(f IFilesystem) { f.SetFullPathForAttr(my.path) }
func (my AttrPath) Joins(paths ...string) PathAttributer {
	my.path = filepath.Join(append([]string{my.path}, paths...)...)
	return my
}
func (my AttrPath) GetPath() string { return my.path }
