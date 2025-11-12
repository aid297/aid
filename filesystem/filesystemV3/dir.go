package filesystemV3

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/aid297/aid/operation/operationV2"
)

type Dir struct {
	IsRel    bool         // 是否使用相对路径
	Error    error        // 错误信息
	Name     string       // 文件名
	BasePath string       // 基础路径
	FullPath string       // 完整路径
	Size     int64        // 文件大小
	Info     os.FileInfo  // 文件信息
	Mode     os.FileMode  // 文件权限
	Exist    bool         // 文件是否存在
	mu       sync.RWMutex // 读写锁
	Files    []*File      // 目录下的文件列表
	Dirs     []*Dir       // 子目录列表
}

// New 实例化
func (*Dir) New(attrs ...DirAttributer) *Dir {
	return (&Dir{mu: sync.RWMutex{}, Files: make([]*File, 0), Dirs: make([]*Dir, 0)}).SetAttrs(attrs...).refresh()
}

// Abs 实例化：绝对路径
func (*Dir) Abs(attrs ...DirAttributer) *Dir {
	return APP.Dir.New(DirIsAbs()).SetAttrs(attrs...).refresh()
}

// Rel 实例化：相对路径
func (*Dir) Rel(attrs ...DirAttributer) *Dir {
	return APP.Dir.New(DirIsRel()).SetAttrs(append([]DirAttributer{DirPath(".")}, attrs...)...).refresh()
}

// SetAttrs 设置属性
func (my *Dir) SetAttrs(attrs ...DirAttributer) *Dir {
	for idx := range attrs {
		attrs[idx].Register(my)
	}
	return my
}

// Lock 加锁 → 写
func (my *Dir) Lock() *Dir {
	my.mu.Lock()
	return my
}

// Unlock 解锁 → 写
func (my *Dir) Unlock() *Dir {
	my.mu.Unlock()
	return my
}

// RLock 加锁 → 读
func (my *Dir) RLock() *Dir {
	my.mu.RLock()
	return my
}

// RUnlock 解锁 → 读
func (my *Dir) RUnlock() *Dir {
	my.mu.RUnlock()
	return my
}

// refresh 刷新目录信息
func (my *Dir) refresh() *Dir {
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
		my.Error = fmt.Errorf("%w:%w", ErrMissFullPath, err)
		return my
	}

	return my
}

// Join 追加目录
func (my *Dir) Join(dirs ...string) *Dir {
	return my.SetAttrs(DirIsAbs(), DirPath(append([]string{my.FullPath}, dirs...)...)).refresh()
}

// Create 创建多级目录
func (my *Dir) Create(attrs ...DirOperationAttributer) *Dir {
	var (
		err       error
		operation = new(DirOperation).SetAttrs(attrs...)
	)

	if my.FullPath == "" {
		my.Error = ErrMissFullPath
		return my
	}

	if err = os.MkdirAll(my.FullPath, operationV2.NewTernary(operationV2.TrueFn(func() os.FileMode { return operation.DirMode }), operationV2.FalseValue(os.FileMode(0777))).GetByValue(operation.DirMode != 0)); err != nil {
		my.Error = fmt.Errorf("%w:%w", ErrCreateDir, err)
		return my
	}

	return my.refresh()
}

// Rename 重命名目录
func (my *Dir) Rename(newName string) *Dir {
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

	return APP.Dir.Abs(DirPath(newPath))
}

// Remove 删除目录
func (my *Dir) Remove() *Dir {
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
func (my *Dir) RemoveAll() *Dir {
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

// LS 列出当前目录下的所有文件和子目录
func (my *Dir) LS() *Dir {
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
		if entry.IsDir() {
			d := APP.Dir.Abs(DirPath(my.FullPath)).Join(entry.Name())
			my.Dirs = append(my.Dirs, d.LS())
		} else {
			my.Files = append(my.Files, APP.File.Abs(FilePath(my.FullPath, entry.Name())))
		}
	}

	return my
}

// CopyFilesTo 复制当前目录下的所有文件到目标路径
func (my *Dir) CopyFilesTo(isRel bool, dstPaths ...string) *Dir {
	var (
		err error
		dst *Dir
	)

	if my.FullPath == "" {
		my.Error = ErrMissFullPath
		return my
	}

	if !my.Exist {
		my.Error = ErrDirNotExist
		return my
	}

	if my.LS().Error != nil {
		return my
	}

	if dst = APP.Dir.New(DirSetRel(isRel), DirPath(dstPaths...)).Create(DirMode(my.Mode)); dst.Error != nil {
		my.Error = dst.Error
		return my
	}

	for idx := range my.Files {
		if err = my.Files[idx].CopyTo(isRel, dstPaths...).Error; err != nil {
			my.Error = err
			return my
		}
	}

	return my
}

// CopyDirsTo 复制当前目录下的所有子目录到目标路径
func (my *Dir) CopyDirsTo(isRel bool, dstPaths ...string) *Dir {
	var dst *Dir

	if my.FullPath == "" {
		my.Error = ErrMissFullPath
		return my
	}

	if !my.Exist {
		my.Error = ErrDirNotExist
		return my
	}

	if my.LS().Error != nil {
		return my
	}

	if dst = APP.Dir.New(DirSetRel(isRel), DirPath(dstPaths...)); dst.Error != nil {
		my.Error = dst.Error
		return my
	}

	if len(my.Dirs) > 0 {
		my.Error = copyDirTo(my.FullPath, dst.FullPath)
	}

	return my
}

// CopyAllTo 复制当前目录下的所有文件和子目录到目标路径
func (my *Dir) CopyAllTo(isRel bool, dstPaths ...string) *Dir {
	var dst *Dir

	if my.FullPath == "" {
		my.Error = ErrMissFullPath
		return my
	}

	if !my.Exist {
		my.Error = ErrDirNotExist
		return my
	}

	if my.LS().Error != nil {
		return my
	}

	if dst = APP.Dir.New(DirSetRel(isRel), DirPath(dstPaths...)); !dst.Exist {
		dst.Create(DirMode(my.Mode))
	}

	// if len(my.Files) > 0 {
	// 	my = my.CopyFilesTo(isRel, dstPaths...)
	// }

	if len(my.Dirs) > 0 {
		my = my.CopyDirsTo(isRel, dstPaths...)

	}

	return my
}

// Copy 复制当前对象
func (my *Dir) Copy() *Dir { return APP.Dir.Abs(DirPath(my.FullPath)) }

// Up 向上一级目录
func (my *Dir) Up() *Dir { return my.SetAttrs(DirIsAbs(), DirPath(my.BasePath)).refresh() }

// ******************** 管理器属性 ******************** //
type (
	DirAttributer interface {
		Register(dir *Dir)
	}
	AttrDirPath  struct{ dirs []string }
	AttrDirIsRel struct{ isRel bool }
)

func DirPath(dirs ...string) AttrDirPath { return AttrDirPath{dirs: dirs} }
func (my AttrDirPath) Register(dir *Dir) {
	dir.FullPath = operationV2.NewTernary(
		operationV2.TrueFn(func() string { return getRootPath(filepath.Join(my.dirs...)) }),
		operationV2.FalseFn(func() string { return filepath.Join(my.dirs...) }),
	).GetByValue(dir.IsRel)
}

func DirSetRel(isRel bool) AttrDirIsRel   { return AttrDirIsRel{isRel: isRel} }
func DirIsAbs() AttrDirIsRel              { return AttrDirIsRel{isRel: false} }
func DirIsRel() AttrDirIsRel              { return AttrDirIsRel{isRel: true} }
func (my AttrDirIsRel) Register(dir *Dir) { dir.IsRel = my.isRel }
