package filesystemV4

import "os"

type Filesystemer interface {
	ImplFilesystemer()
	GetExist() bool
	GetError() error
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
	CopyTo(isRel bool, dstPaths ...string) Filesystemer
	Copy() Filesystemer
	Up() Filesystemer
	LS() Filesystemer
}
