package filesystemV2

import (
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/aid297/aid/array"
)

type (
	Dir struct {
		err      error                  // 错误信息
		name     string                 // 文件名
		basePath string                 // 基础路径
		fullPath string                 // 完整路径
		size     int64                  // 文件大小
		info     os.FileInfo            // 文件信息
		mode     os.FileMode            // 文件权限
		exist    bool                   // 文件是否存在
		mu       sync.RWMutex           // 读写锁
		files    *array.AnyArray[*File] // 目录下的文件
		dirs     *array.AnyArray[*Dir]  // 目录下的文件夹
	}
	DirCollection struct {
		Dirs  []*Dir
		Files []*File
	}
)

var DirApp Dir

// WhoAmI 获取当前类型：（目录）
func (*Dir) WhoAmI() FilesystemV2Type { return FilesystemV2Dir }

// NewByAbs 实例化：通过绝对路径
func (*Dir) NewByAbs(paths ...string) *Dir {
	ins := &Dir{fullPath: filepath.Join(paths...), files: array.Make[*File](0), dirs: array.Make[*Dir](0)}
	ins.refresh()

	return ins
}

// NewByRel 实例化：通过相对路径
func (*Dir) NewByRel(paths ...string) *Dir {
	ins := &Dir{fullPath: getRootPath(filepath.Join(paths...)), files: array.Make[*File](0), dirs: array.Make[*Dir](0)}
	ins.refresh()

	return ins
}

// Lock 加锁：写锁
func (my *Dir) Lock() *Dir {
	my.mu.Lock()
	return my
}

// Unlock 解锁：写锁
func (my *Dir) Unlock() *Dir {
	my.mu.Unlock()
	return my
}

// RLock 加锁：读锁
func (my *Dir) RLock() *Dir {
	my.mu.RLock()
	return my
}

// RUnlock 解锁：读锁
func (my *Dir) RUnlock() *Dir {
	my.mu.RUnlock()
	return my
}

// getName 获取文件名
func (my *Dir) getName() string { return my.name }

// GetName 获取文件名
func (my *Dir) GetName() string {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getName()
}

// getBasePath 获取基础路径
func (my *Dir) getBasePath() string { return my.basePath }

// GetBasePath 获取基础路径
func (my *Dir) GetBasePath() string {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getBasePath()
}

// getFullPath 获取完整路径
func (my *Dir) getFullPath() string { return my.fullPath }

// GetFullPath 获取完整路径
func (my *Dir) GetFullPath() string {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getFullPath()
}

// getSize 获取文件夹大小
func (my *Dir) getSize() int64 { return my.size }

// GetSize 获取文件夹大小
func (my *Dir) GetSize() int64 {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getSize()
}

// getInfo 获取文件夹信息
func (my *Dir) getInfo() os.FileInfo { return my.info }

// GetInfo 获取文件夹信息
func (my *Dir) GetInfo() os.FileInfo {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getInfo()
}

// getMode 获取文件夹权限
func (my *Dir) getMode() os.FileMode { return my.mode }

// GetMode 获取文件夹权限
func (my *Dir) GetMode() os.FileMode {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getMode()
}

// getExist 获取目录是否存在
func (my *Dir) getExist() bool { return my.exist }

// GetExist 获取目录是否存在
func (my *Dir) GetExist() bool {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getExist()
}

// Error 获取错误
func (my *Dir) Error() error { return my.err }

// copy 复制当前对象
func (my *Dir) copy() *Dir { return DirApp.NewByAbs(my.fullPath) }

// Copy 复制当前对象
func (my *Dir) Copy() *Dir {
	my.mu.Lock()
	defer my.mu.Unlock()

	return my.copy()
}

// refresh 刷新目录信息
func (my *Dir) refresh() {
	if my.fullPath != "" {
		if my.info, my.err = os.Stat(my.fullPath); my.err != nil {
			if os.IsNotExist(my.err) {
				my.name = ""
				my.size = 0
				my.mode = 0
				my.basePath = path.Dir(my.fullPath)
				my.exist = false
				my.err = nil
				return
			} else {
				my.err = DirInitErr.Wrap(my.err)
				return
			}
		}

		my.name = my.info.Name()
		my.size = my.info.Size()
		my.mode = my.info.Mode()
		my.basePath = path.Dir(my.fullPath)
		my.exist = true
		my.err = nil
	} else {
		my.err = DirFullPathEmptyErr.New("")
	}
}

// create 新建目录
func (my *Dir) create(mode os.FileMode) {
	if my.fullPath == "" {
		my.err = DirFullPathEmptyErr.New("")
		return
	}

	err := os.MkdirAll(my.fullPath, mode)
	if err != nil {
		my.err = err
		return
	}

	my.err = nil
}

// Create 新建目录
func (my *Dir) Create(mode os.FileMode) *Dir {
	my.mu.Lock()
	defer my.mu.Unlock()
	defer my.refresh()

	my.create(mode)

	return my
}

// CreateDefaultMode 新建目录：默认权限
func (my *Dir) CreateDefaultMode() *Dir {
	my.mu.Lock()
	defer my.mu.Unlock()
	defer my.refresh()

	my.create(os.ModePerm)

	return my
}

// 修改文件名
func (my *Dir) rename(newName string) {
	if my.fullPath == "" {
		my.err = DirFullPathEmptyErr.New("")
		return
	}

	newPath := filepath.Join(filepath.Dir(my.fullPath), newName)
	err := os.Rename(my.fullPath, newPath)
	if err != nil {
		my.err = err
		return
	}

	my.fullPath = newPath
	my.err = nil
}

// Rename 文件夹重命名
func (my *Dir) Rename(newName string) *Dir {
	my.mu.Lock()
	defer my.mu.Unlock()
	defer my.refresh()

	my.rename(newName)

	return my
}

func (my *Dir) removeAll() { my.err = os.RemoveAll(my.fullPath) }

func (my *Dir) RemoveAll() *Dir {
	my.mu.Lock()
	defer my.mu.Unlock()
	defer my.refresh()

	my.removeAll()

	return my
}

// remove 删除目录
func (my *Dir) remove() { my.err = os.Remove(my.fullPath) }

// Remove 删除目录
func (my *Dir) Remove() *Dir {
	my.mu.Lock()
	defer my.mu.Unlock()
	defer my.refresh()

	my.remove()

	return my
}

// checkPermission 检查目录权限
func (my *Dir) checkPermission(mode os.FileMode) bool { return my.mode == mode }

// CheckPermission 检查目录权限
func (my *Dir) CheckPermission(mode os.FileMode) bool {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.checkPermission(mode)
}

// join 添加路径
func (my *Dir) join(path string) { my.fullPath = filepath.Join(my.fullPath, path) }

// Join 添加路径
func (my *Dir) Join(path string) *Dir {
	my.mu.Lock()
	defer my.mu.Unlock()
	defer my.refresh()

	my.join(path)

	return my
}

// joins 添加多个路径
func (my *Dir) joins(paths ...string) {
	my.fullPath = filepath.Join(my.fullPath, filepath.Join(paths...))
}

// Joins 添加多个路径
func (my *Dir) Joins(paths ...string) *Dir {
	my.mu.Lock()
	defer my.mu.Unlock()
	defer my.refresh()

	my.joins(paths...)

	return my
}

// ls 获取当前目录内容
func (my *Dir) ls() {
	var entries []os.DirEntry

	entries, my.err = os.ReadDir(my.fullPath)
	if my.err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			d := DirApp.NewByAbs(my.fullPath).Join(entry.Name())
			my.dirs.Append(d.Ls())
		} else {
			my.files.Append(FileApp.NewByAbs(filepath.Join(my.fullPath, entry.Name())))
		}
	}
}

// Ls 获取当前目录内容
func (my *Dir) Ls() *Dir {
	my.mu.RLock()
	defer my.mu.RUnlock()

	my.ls()

	return my
}

// getDirs 获取当前目录下所有文件夹
func (my *Dir) getDirs() *array.AnyArray[*Dir] {
	if my.dirs.IsEmpty() {
		my.ls()
	}

	return my.dirs
}

// GetDirs 获取当前目录下所有文件夹
func (my *Dir) GetDirs() *array.AnyArray[*Dir] {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getDirs()
}

// getFiles 获取当前目录下所有文件
func (my *Dir) getFiles() *array.AnyArray[*File] {
	if my.files.IsEmpty() {
		my.ls()
	}

	return my.files
}

// GetFiles 获取当前目录下所有文件
func (my *Dir) GetFiles() *array.AnyArray[*File] {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getFiles()
}

// copyAllFilesTo 复制所有文件
func (my *Dir) copyAllFilesTo(dstDir string) {
	if !my.getExist() {
		return
	}

	dst := DirApp.NewByAbs(dstDir)
	if !dst.getExist() {
		dst.create(os.ModePerm)
	}

	my.getFiles().Each(func(idx int, item *File) { item.CopyTo(dst.Copy().Join(item.getName()).GetFullPath()) })
}

// CopyAllFilesTo 复制所有文件
func (my *Dir) CopyAllFilesTo(dstDir string) *Dir {
	my.mu.Lock()
	defer my.mu.Unlock()
	defer my.refresh()

	my.copyAllFilesTo(dstDir)

	return my
}

// copyAllDirsTo 复制目录下所有文件夹及文件
func (my *Dir) copyAllDirsTo(dstDir string) {
	if !my.getExist() {
		return
	}

	dst := DirApp.NewByAbs(dstDir)
	if !dst.getExist() {
		dst.create(os.ModePerm)
	}

	if my.getFiles().IsNotEmpty() {
		my.copyAllFilesTo(dst.Join(my.getName()).Create(my.getMode()).GetFullPath())
	}

	if my.getDirs().IsNotEmpty() {
		my.dirs.Each(func(idx int, item *Dir) {
			dstFullPath := dst.Copy().Join(item.getName()).Create(item.getMode()).GetFullPath()
			if item.files.IsNotEmpty() {
				item.copyAllFilesTo(dstFullPath)
			}
		})
	}
}

// CopyAllDirsTo 复制目录下所有文件夹及文件
func (my *Dir) CopyAllDirsTo(dstDir string) *Dir {
	my.mu.Lock()
	defer my.mu.Unlock()
	defer my.refresh()

	my.copyAllDirsTo(dstDir)

	return my
}
