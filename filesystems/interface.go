package filesystems

import (
	"os"

	"github.com/aid297/aid/v3/operations"
)

type Filesystem interface {
	GetName() string
	GetExist() bool
	GetError() error
	GetBasePath() string
	GetFullPath() string
	GetInfo() os.FileInfo
	GetDirs() []Filesystem
	GetFiles() []Filesystem
	GetKind() string
	SetAttrs(attrs ...PathAttributer) Filesystem
	SetFullPathForAttr(path string) Filesystem
	SetFullPathByAttr(attrs ...PathAttributer) Filesystem
	refresh() Filesystem
	Lock() Filesystem
	Unlock() Filesystem
	RLock() Filesystem
	RUnlock() Filesystem
	Join(paths ...string) Filesystem
	Create(attrs ...OperationAttributer) Filesystem
	Rename(newName string) Filesystem
	Remove() Filesystem
	RemoveAll() Filesystem
	Write(content []byte, attrs ...OperationAttributer) Filesystem
	Read(attrs ...OperationAttributer) ([]byte, error)
	CopyTo(isRel bool, dstPaths ...string) Filesystem
	Copy() Filesystem
	Up() Filesystem
	LS() Filesystem
	Zip() Filesystem
}

func New(attr PathAttributer) (Filesystem, error) {
	isDir, err := isDir(attr.GetPath())
	if err != nil {
		return nil, err
	}
	return operations.NewTernary(operations.TrueFn(func() Filesystem { return NewFile(attr) }), operations.FalseFn(func() Filesystem { return NewDir(attr) })).GetByValue(!isDir), nil
}
