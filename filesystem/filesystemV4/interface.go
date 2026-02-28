package filesystemV4

import (
	"os"

	"github.com/aid297/aid/operation/operationV2"
)

type IFilesystem interface {
	GetName() string
	GetExist() bool
	GetError() error
	GetBasePath() string
	GetFullPath() string
	GetInfo() os.FileInfo
	GetDirs() []IFilesystem
	GetFiles() []IFilesystem
	GetKind() string
	SetAttrs(attrs ...PathAttributer) IFilesystem
	SetFullPathForAttr(path string) IFilesystem
	SetFullPathByAttr(attrs ...PathAttributer) IFilesystem
	refresh() IFilesystem
	Lock() IFilesystem
	Unlock() IFilesystem
	RLock() IFilesystem
	RUnlock() IFilesystem
	Join(paths ...string) IFilesystem
	Create(attrs ...OperationAttributer) IFilesystem
	Rename(newName string) IFilesystem
	Remove() IFilesystem
	RemoveAll() IFilesystem
	Write(content []byte, attrs ...OperationAttributer) IFilesystem
	Read(attrs ...OperationAttributer) ([]byte, error)
	CopyTo(isRel bool, dstPaths ...string) IFilesystem
	Copy() IFilesystem
	Up() IFilesystem
	LS() IFilesystem
	Zip() IFilesystem
}

func New(attr PathAttributer) (IFilesystem, error) {
	isDir, err := isDir(attr.GetPath())
	if err != nil {
		return nil, err
	}
	return operationV2.NewTernary(operationV2.TrueFn(func() IFilesystem { return NewFile(attr) }), operationV2.FalseFn(func() IFilesystem { return NewDir(attr) })).GetByValue(!isDir), nil
}
