package filesystemV4

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/aid297/aid/operation/operationV2"
)

type (
	Dir struct {
		Error    error         // 错误信息
		Name     string        // 文件名
		BasePath string        // 基础路径
		FullPath string        // 完整路径
		Size     int64         // 文件大小
		Info     os.FileInfo   // 文件信息
		Mode     os.FileMode   // 文件权限
		Exist    bool          // 文件是否存在
		mu       *sync.RWMutex // 读写锁
		Files    []File        // 目录下的文件列表
		Dirs     []Dir         // 子目录列表
	}
)

func (Dir) Abs(dirs ...string) Dir {
	return Dir{mu: &sync.RWMutex{}, FullPath: getRootPath(path.Join(dirs...))}.refresh()
}

func (my Dir) Rel(dirs ...string) Dir {
	return Dir{mu: &sync.RWMutex{}, FullPath: path.Join(dirs...)}.refresh()
}

// refresh 刷新目录信息
func (my Dir) refresh() Dir {
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
			}
			my.Error = fmt.Errorf("%w:%w", ErrInit, err)
			return my
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
func (my Dir) Lock() Dir {
	my.mu.Lock()
	return my
}

// Unlock 解锁 → 写
func (my Dir) Unlock() Dir {
	my.mu.Unlock()
	return my
}

// RLock 加锁 → 读
func (my Dir) RLock() Dir {
	my.mu.RLock()
	return my
}

// RUnlock 解锁 → 读
func (my Dir) RUnlock() Dir {
	my.mu.RUnlock()
	return my
}

func (my Dir) Join(dirs ...string) Dir {
	my.FullPath = path.Join(append([]string{my.FullPath}, dirs...)...)
	return my.refresh()
}

func (my Dir) Create(attrs ...DirOperationAttributer) Dir {
	var (
		err       error
		operation = APP.DirOperation.New().SetAttrs(attrs...)
	)

	if my.FullPath == "" {
		my.Error = ErrMissFullPath
		return my
	}

	if err = os.MkdirAll(my.FullPath, operation.DirMode); err != nil {
		my.Error = fmt.Errorf("%w:%w", ErrCreateDir, err)
		return my
	}

	return my.refresh()
}

// Rename 重命名目录
func (my Dir) Rename(newName string) Dir {
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

	return APP.Dir.Abs(newPath)
}

// Remove 删除目录
func (my Dir) Remove() Dir {
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
func (my Dir) RemoveAll() Dir {
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
func (my Dir) LS() Dir {
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
			d := APP.Dir.Abs(my.FullPath, entry.Name())
			my.Dirs = append(my.Dirs, d.LS())
		} else {
			my.Files = append(my.Files, APP.File.Abs(my.FullPath, entry.Name()))
		}
	}

	return my
}

// CopyFilesTo 复制当前目录下的所有文件到目标路径
func (my Dir) CopyFilesTo(isRel bool, dstPaths ...string) Dir {
	var (
		err error
		dst Dir
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

	if dst = APP.Dir.Rel(dstPaths...).Create(APP.DirOperAttr.Mode.Set(my.Mode)); dst.Error != nil {
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
func (my Dir) CopyDirsTo(isRel bool, dstPaths ...string) Dir {
	var dst = operationV2.NewTernary(operationV2.TrueFn(func() Dir {
		return APP.Dir.Rel(dstPaths...)
	}), operationV2.FalseFn(func() Dir {
		return APP.Dir.Abs(dstPaths...)
	})).GetByValue(isRel)

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

	if dst = APP.Dir.Rel(dstPaths...); dst.Error != nil {
		my.Error = dst.Error
		return my
	}

	if len(my.Dirs) > 0 {
		my.Error = copyDirTo(my.FullPath, dst.FullPath)
	}

	return my
}

// CopyAllTo 复制当前目录下的所有文件和子目录到目标路径
func (my Dir) CopyAllTo(isRel bool, dstPaths ...string) Dir {
	var dst = operationV2.NewTernary(operationV2.TrueFn(func() Dir {
		return APP.Dir.Rel(dstPaths...)
	}), operationV2.FalseFn(func() Dir {
		return APP.Dir.Abs(dstPaths...)
	})).GetByValue(isRel)

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

	if !dst.Exist {
		dst.Create(APP.DirOperAttr.Mode.Set(my.Mode))
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
func (my Dir) Copy() Dir { return APP.Dir.Abs(my.FullPath) }

// Up 向上一级目录
func (my Dir) Up() Dir {
	my.FullPath = path.Base(my.FullPath)
	return my.refresh()
}
