package filesystemV3

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
)

func NewDir(attrs ...DirAttributer) *Dir {
	dir := &Dir{mu: sync.RWMutex{}, Files: make([]*File, 0), Dirs: make([]*Dir, 0)}
	return dir.SetAttrs(attrs...).refresh()
}

func (my *Dir) SetAttrs(attrs ...DirAttributer) *Dir {
	for idx := range attrs {
		attrs[idx].Register(my)
	}
	return my
}

func (my *Dir) Lock() *Dir {
	my.mu.Lock()
	return my
}

func (my *Dir) Unlock() *Dir {
	my.mu.Unlock()
	return my
}

func (my *Dir) RLock() *Dir {
	my.mu.RLock()
	return my
}

func (my *Dir) RUnlock() *Dir {
	my.mu.RUnlock()
	return my
}

func (my *Dir) refresh() *Dir {
	var err error
	if my.FullPath != "" {
		if my.Info, err = os.Stat(my.FullPath); err != nil {
			if os.IsNotExist(my.Error) {
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

func (my *Dir) Join(dirs ...string) *Dir {
	allDirs := []string{my.FullPath}
	if len(dirs) == 0 {
		return my
	}
	allDirs = append(allDirs, dirs...)

	return NewDir(DirIsRel(), DirPath(allDirs...))
}

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

	return NewDir(DirIsAbs(), DirPath(newPath))
}

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

	return my
}

func (my *Dir) Ls() *Dir {
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
			d := NewDir(DirIsAbs(), DirPath(my.FullPath)).Join(entry.Name())
			my.Dirs = append(my.Dirs, d.Ls())
		} else {
			my.Files = append(my.Files, NewFile(FilePath(my.FullPath, entry.Name())))
		}
	}

	return my
}

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

	if dst = NewDir(DirSetRel(isRel), DirPath(dstPaths...)).Create(DirMode(my.Mode)); dst.Error != nil {
		my.Error = dst.Error
		return my
	}

	for idx := range my.Files {
		if err = my.Files[idx].CopyTo(my.FullPath).Error; err != nil {
			my.Error = err
			return my
		}
	}

	return my
}

func (my *Dir) CopyDirsTo(isRel bool, dstPaths ...string) *Dir {
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

	if dst = NewDir(DirSetRel(isRel), DirPath(dstPaths...)).Create(DirMode(my.Mode)); dst.Error != nil {
		my.Error = dst.Error
		return my
	}

	for idx := range my.Dirs {
		if err = dst.Join(my.Dirs[idx].Name).Create(DirMode(my.Dirs[idx].Mode)).Error; err != nil {
			my.Error = fmt.Errorf("%w:%w", ErrCreateDir, err)
			return my
		}
	}

	return my
}

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

	if dst = NewDir(DirSetRel(isRel), DirPath(dstPaths...)); !dst.Exist {
		dst.Create(DirMode(my.Mode))
	}

	if len(my.Files) > 0 {
		my = my.CopyFilesTo(isRel, dstPaths...)
	}

	if len(my.Dirs) > 0 {
		my = my.CopyDirsTo(isRel, dstPaths...)
	}

	return my
}

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
