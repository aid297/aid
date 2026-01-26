package filesystemV4

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/aid297/aid/operation/operationV2"
)

type Dir struct {
	Error    error          `json:"error"`    // 错误信息
	Name     string         `json:"name"`     // 文件名
	BasePath string         `json:"basePath"` // 基础路径
	FullPath string         `json:"fullPath"` // 完整路径
	Size     int64          `json:"size"`     // 文件大小
	Info     os.FileInfo    `json:"info"`     // 文件信息
	Mode     os.FileMode    `json:"mode"`     // 文件权限
	Exist    bool           `json:"exist"`    // 文件是否存在
	mu       sync.RWMutex   // 读写锁
	Files    []Filesystemer `json:"files"` // 目录下的文件列表
	Dirs     []Filesystemer `json:"dirs"`  // 子目录列表
	Kind     string         `json:"kind"`  // 类型
}

func NewDir(attrs ...PathAttributer) Filesystemer {
	return (&Dir{mu: sync.RWMutex{}}).SetAttrs(attrs...).refresh()
}

func (my *Dir) GetExist() bool           { return my.Exist }
func (my *Dir) GetError() error          { return my.Error }
func (my *Dir) GetFullPath() string      { return my.FullPath }
func (my *Dir) GetInfo() os.FileInfo     { return my.Info }
func (my *Dir) GetDirs() []Filesystemer  { return my.Dirs }
func (my *Dir) GetFiles() []Filesystemer { return my.Files }

func (my *Dir) SetAttrs(attrs ...PathAttributer) Filesystemer {
	my.mu.Lock()
	defer my.mu.Unlock()
	for idx := range attrs {
		attrs[idx].Register(my)
	}
	return my
}

func (my *Dir) SetFullPathForAttr(path string) Filesystemer { my.FullPath = path; return my }

func (my *Dir) SetFullPathByAttr(attrs ...PathAttributer) Filesystemer {
	return my.SetAttrs(attrs...).refresh()
}

func (my *Dir) refresh() Filesystemer {
	var err error
	if my.FullPath != "" {
		if my.Info, err = os.Stat(my.FullPath); err != nil {
			if os.IsNotExist(err) {
				my.Name = ""
				my.Size = 0
				my.Mode = 0
				my.BasePath = path.Dir(my.FullPath)
				my.Exist = false
				my.Error = nil
				return my
			} else {
				my.Error = fmt.Errorf("%w:%w", ErrInit, err)
				return my
			}
		}

		my.Name = my.Info.Name()
		my.Size = my.Info.Size()
		my.Mode = my.Info.Mode()
		my.BasePath = path.Dir(my.FullPath)
		my.Exist = true
		my.Error = nil
	} else {
		my.Error = ErrMissFullPath
		return my
	}

	return my
}

// Lock 加锁 → 写
func (my *Dir) Lock() Filesystemer { my.mu.Lock(); return my }

// Unlock 解锁 → 写
func (my *Dir) Unlock() Filesystemer { my.mu.Unlock(); return my }

// RLock 加锁 → 读
func (my *Dir) RLock() Filesystemer { my.mu.RLock(); return my }

// RUnlock 解锁 → 读
func (my *Dir) RUnlock() Filesystemer { my.mu.RUnlock(); return my }

func (my *Dir) Join(paths ...string) Filesystemer {
	my.FullPath = filepath.Join(append([]string{my.FullPath}, paths...)...)
	return my.refresh()
}

// Create 创建多级目录
func (my *Dir) Create(attrs ...OperationAttributer) Filesystemer {
	var (
		err       error
		operation = NewOperation(attrs...)
	)

	if my.FullPath == "" {
		my.Error = ErrMissFullPath
		return my
	}

	if err = os.MkdirAll(my.FullPath, operationV2.NewTernary(operationV2.TrueFn(func() os.FileMode { return operation.Mode }), operationV2.FalseValue(os.FileMode(0777))).GetByValue(operation.Mode != 0)); err != nil {
		my.Error = fmt.Errorf("%w:%w", ErrCreateDir, err)
		return my
	}

	return my.refresh()
}

// Rename 重命名目录
func (my *Dir) Rename(newName string) Filesystemer {
	var err error

	if my.FullPath == "" {
		my.Error = ErrMissFullPath
		return my
	}

	newPath := filepath.Join(filepath.Dir(my.FullPath), newName)
	if err = os.Rename(my.FullPath, newPath); err != nil {
		my.Error = fmt.Errorf("%w:%w", ErrRename, err)
		return my
	}

	return NewDir(Abs(newPath))
}

// Remove 删除目录
func (my *Dir) Remove() Filesystemer {
	var err error

	if my.FullPath == "" {
		my.Error = ErrMissFullPath
		return my
	}

	if err = os.Remove(my.FullPath); err != nil {
		my.Error = fmt.Errorf("%w:%w", ErrRemove, err)
		return my
	}

	return my.refresh()
}

// RemoveAll 递归删除目录
func (my *Dir) RemoveAll() Filesystemer {
	var err error

	if my.FullPath == "" {
		my.Error = ErrMissFullPath
		return my
	}

	if err = os.RemoveAll(my.FullPath); err != nil {
		my.Error = fmt.Errorf("%w:%w", ErrRemove, err)
		return my
	}

	return my.refresh()
}

func (my *Dir) Write(content []byte, attrs ...OperationAttributer) Filesystemer { return my }

func (my *Dir) Read(attrs ...OperationAttributer) ([]byte, error) { return nil, nil }

// LS 列出当前目录下的所有文件和子目录
func (my *Dir) LS() Filesystemer {
	var (
		err     error
		entries []os.DirEntry
	)

	if my.FullPath == "" {
		my.Error = ErrMissFullPath
		return my
	}

	if entries, err = os.ReadDir(my.FullPath); err != nil {
		my.Error = fmt.Errorf("%w:%w", ErrReadDir, err)
		return my
	}

	for _, entry := range entries {
		fmt.Println(entry.Name(), entry.IsDir())
		if entry.IsDir() {
			d := NewDir(Abs(my.FullPath, entry.Name()))
			my.Dirs = append(my.Dirs, d)
		} else {
			my.Files = append(my.Files, NewFile(Abs(my.FullPath, entry.Name())))
		}
	}

	return my
}

// CopyFilesTo 复制当前目录下的所有文件到目标路径
func (my *Dir) CopyFilesTo(isRel bool, dstPaths ...string) *Dir {
	var (
		err error
		dst Filesystemer
	)

	if my.FullPath == "" {
		my.Error = ErrMissFullPath
		return my
	}

	if !my.Exist {
		my.Error = ErrDirNotExist
		return my
	}

	if my.LS().GetError() != nil {
		return my
	}

	if dst = NewDir(Rel(dstPaths...)).Create(Mode(my.Mode)); dst.GetError() != nil {
		my.Error = dst.GetError()
		return my
	}

	for idx := range my.Files {
		if err = my.Files[idx].CopyTo(isRel, dstPaths...).GetError(); err != nil {
			my.Error = err
			return my
		}
	}

	return my
}

// CopyTo 复制当前目录下的所有子目录到目标路径
func (my *Dir) CopyDirsTo(isRel bool, dstPaths ...string) *Dir {
	if my.FullPath == "" {
		my.Error = ErrMissFullPath
		return my
	}

	if !my.Exist {
		my.Error = ErrDirNotExist
		return my
	}

	if my.LS().GetError() != nil {
		return my
	}

	if len(my.Dirs) > 0 {
		my.Error = copyDirTo(my.FullPath, NewDir(Rel(dstPaths...)).GetFullPath())
	}

	return my
}

// CopyAllTo 复制当前目录下的所有文件和子目录到目标路径
func (my *Dir) CopyTo(isRel bool, dstPaths ...string) Filesystemer {
	var (
		err error
		dst Filesystemer
	)

	if my.FullPath == "" {
		my.Error = ErrMissFullPath
		return my
	}

	if !my.Exist {
		my.Error = ErrDirNotExist
		return my
	}

	if my.LS().GetError() != nil {
		return my
	}

	if dst = NewDir(Rel(dstPaths...)); !dst.GetExist() {
		if err = dst.Create(Mode(my.Mode)).GetError(); err != nil {
			my.Error = err
			return my
		}
	}

	if len(my.Files) > 0 {
		my.CopyFilesTo(isRel, dstPaths...)
	}

	if len(my.Dirs) > 0 {
		my.CopyDirsTo(isRel, dstPaths...)

	}

	return my
}

// Copy 复制当前对象
func (my *Dir) Copy() Filesystemer { return NewDir(Abs(my.FullPath)) }

// Up 向上一级目录
func (my *Dir) Up() Filesystemer { my.FullPath = my.BasePath; return my.refresh() }
