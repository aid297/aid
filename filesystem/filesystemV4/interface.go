package filesystemV4

import "os"

type Filesystemer interface {
	GetName() string
	GetExist() bool
	GetError() error
	GetBasePath() string
	GetFullPath() string
	GetInfo() os.FileInfo
	GetDirs() []Filesystemer
	GetFiles() []Filesystemer
	SetAttrs(attrs ...PathAttributer) Filesystemer
	SetFullPathForAttr(path string) Filesystemer
	SetFullPathByAttr(attrs ...PathAttributer) Filesystemer
	refresh() Filesystemer
	Lock() Filesystemer
	Unlock() Filesystemer
	RLock() Filesystemer
	RUnlock() Filesystemer
	Join(paths ...string) Filesystemer
	Create(attrs ...OperationAttributer) Filesystemer
	Rename(newName string) Filesystemer
	Remove() Filesystemer
	RemoveAll() Filesystemer
	Write(content []byte, attrs ...OperationAttributer) Filesystemer
	Read(attrs ...OperationAttributer) ([]byte, error)
	CopyTo(isRel bool, dstPaths ...string) Filesystemer
	Copy() Filesystemer
	Up() Filesystemer
	LS() Filesystemer
}
