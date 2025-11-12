package filesystemV2

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"sync"
)

type File struct {
	err       error        // 错误信息
	name      string       // 文件名
	basePath  string       // 基础路径
	fullPath  string       // 完整路径
	extension string       // 文件扩展名
	size      int64        // 文件大小
	fileInfo  os.FileInfo  // 文件信息
	mode      os.FileMode  // 文件权限
	exist     bool         // 文件是否存在
	mu        sync.RWMutex // 读写锁
	mime      string       // 文件Mime类型
}

var FileApp File

const (
	DefaultCreateMode = os.O_CREATE | os.O_RDWR
	DefaultReadMode   = os.O_RDWR
)

// WhoAmI 获取当前类型：（文件）
func (*File) WhoAmI() FilesystemV2Type { return FilesystemV2File }

// NewByAbs 实例化：通过绝对路径
func (*File) NewByAbs(paths ...string) *File {
	ins := &File{fullPath: filepath.Join(paths...)}
	ins.refresh()

	return ins
}

// NewByRel 实例化：通过相对路径
func (*File) NewByRel(paths ...string) *File {
	ins := &File{fullPath: getRootPath(filepath.Join(paths...))}
	ins.refresh()

	return ins
}

// Lock 加锁：写锁
func (my *File) Lock() *File {
	my.mu.Lock()
	return my
}

// Unlock 解锁：写锁
func (my *File) Unlock() *File {
	my.mu.Unlock()
	return my
}

// RLock 加锁：读锁
func (my *File) RLock() *File {
	my.mu.RLock()
	return my
}

// RUnlock 解锁：读锁
func (my *File) RUnlock() *File {
	my.mu.RUnlock()
	return my
}

// getName 获取文件名
func (my *File) getName() string { return my.name }

// GetName 获取文件名
func (my *File) GetName() string {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getName()
}

// getBasePath 获取基础路径
func (my *File) getBasePath() string { return my.basePath }

// GetBasePath 获取基础路径
func (my *File) GetBasePath() string {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getBasePath()
}

// getFullPath 获取完整路径
func (my *File) getFullPath() string { return my.fullPath }

// GetFullPath 获取完整路径
func (my *File) GetFullPath() string {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getFullPath()
}

// getExtension 获取文件扩展名
func (my *File) getExtension() string { return my.extension }

// GetExtension 获取文件扩展名
func (my *File) GetExtension() string {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getExtension()
}

// getSize 获取文件大小
func (my *File) getSize() int64 { return my.size }

// GetSize 获取文件大小
func (my *File) GetSize() int64 {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getSize()
}

// getInfo 获取文件信息
func (my *File) getInfo() os.FileInfo { return my.fileInfo }

// GetInfo 获取文件信息
func (my *File) GetInfo() os.FileInfo {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getInfo()
}

// getMode 获取文件权限
func (my *File) getMode() os.FileMode { return my.mode }

// GetMode 获取文件权限
func (my *File) GetMode() os.FileMode {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getMode()
}

// getExist 获取文件是否存在
func (my *File) getExist() bool { return my.exist }

// GetExist 获取文件是否存在
func (my *File) GetExist() bool {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getExist()
}

// getMime 获取文件Mime类型
func (my *File) getMime() string { return my.mime }

// GetMime 获取文件Mime类型
func (my *File) GetMime() string {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getMime()
}

func (my *File) join(path string) { my.fullPath = filepath.Join(my.getFullPath(), path) }

// Join 多路径拼接
func (my *File) Join(path string) *File {
	my.mu.Lock()
	defer my.mu.Unlock()
	defer my.refresh()

	my.join(path)

	return my
}

func (my *File) joins(paths ...string) {
	allPaths := []string{my.getFullPath()}
	allPaths = append(allPaths, paths...)
	my.fullPath = filepath.Join(allPaths...)
}

// Joins 多路径拼接
func (my *File) Joins(paths ...string) *File {
	my.mu.Lock()
	defer my.mu.Unlock()
	defer my.refresh()

	my.joins(paths...)

	return my
}

// Error 获取错误
func (my *File) Error() error { return my.err }

// copy 复制当前对象
func (my *File) copy() *File { return FileApp.NewByAbs(my.fullPath) }

// Copy 复制当前对象
func (my *File) Copy() *File {
	my.mu.Lock()
	defer my.mu.Unlock()

	return my.copy()
}

// refresh 刷新文件信息
func (my *File) refresh() {
	if my.fullPath != "" {
		if my.fileInfo, my.err = os.Stat(my.fullPath); my.err != nil {
			if os.IsNotExist(my.err) {
				my.name = ""
				my.size = 0
				my.mode = 0
				my.basePath = filepath.Dir(my.fullPath)
				my.extension = filepath.Ext(my.fullPath)
				my.exist = false
				my.err = nil
				return
			} else {
				my.err = FileInitErr.Wrap(my.err)
				return
			}
		}

		my.name = my.fileInfo.Name()
		my.size = my.fileInfo.Size()
		my.mode = my.fileInfo.Mode()
		my.basePath = path.Dir(my.fullPath)
		my.extension = path.Ext(my.fullPath)
		my.exist = true
		my.err = nil
	} else {
		my.err = FileFullPathEmptyErr.New("")
	}
}

// create 创建文件
func (my *File) create(mode os.FileMode, operations ...int) {
	var operation = DefaultCreateMode | os.O_TRUNC

	if len(operations) > 0 {
		operation = operations[0]
	}

	if !my.exist {
		var newFile *os.File
		newFile, my.err = os.OpenFile(my.getFullPath(), operation, mode)
		if my.err != nil {
			my.err = CreateFileErr.Wrap(my.err)
		}
		defer func() { _ = newFile.Close() }()
	}
}

// Create 创建文件
func (my *File) Create(mode os.FileMode, operations ...int) *File {
	my.mu.Lock()
	defer my.mu.Unlock()
	defer my.refresh()

	my.create(mode, operations...)

	return my
}

func (my *File) CreateDefaultMode(operations ...int) *File {
	my.mu.Lock()
	defer my.mu.Unlock()
	defer my.refresh()

	my.create(os.ModePerm, operations...)

	return my
}

// rename 修改文件名
func (my *File) rename(newName string) string {
	newPath := path.Join(path.Dir(my.fullPath), newName)
	my.err = os.Rename(my.fullPath, newPath)
	if my.err != nil {
		my.err = RenameFileErr.Wrap(my.err)
		return my.fullPath
	}

	return newPath
}

// Rename 修改文件名
func (my *File) Rename(newName string) *File {
	my.mu.Lock()
	defer my.mu.Unlock()
	defer my.refresh()

	my.fullPath = my.rename(newName)
	return my
}

// remove 删除文件
func (my *File) remove() { my.err = os.Remove(my.fullPath) }

// Remove 删除文件
func (my *File) Remove() *File {
	my.mu.Lock()
	defer my.mu.Unlock()
	defer my.refresh()

	my.remove()

	return my
}

// checkPermission 检查文件权限
func (my *File) checkPermission(operations ...int) {
	var (
		file       *os.File
		permission = os.O_RDONLY
	)

	if len(operations) > 0 {
		permission = operations[0]
	}

	file, my.err = os.OpenFile(my.fullPath, permission, 0666)
	_ = file.Close()
}

// CheckPermission 检查文件权限
func (my *File) CheckPermission(permissions ...int) *File {
	my.mu.RLock()
	defer my.mu.RUnlock()

	my.checkPermission(permissions...)

	return my
}

// write 写入文件
func (my *File) write(content []byte, mode os.FileMode, operations ...int) int {
	var (
		operation    = DefaultCreateMode
		file         *os.File
		bytesWritten int
	)

	if len(operations) > 0 {
		operation = operations[0]
	}

	file, my.err = os.OpenFile(my.getFullPath(), operation, mode)
	defer func() { _ = file.Close() }()
	if my.err != nil {
		my.err = WriteFileErr.Wrap(my.err)
		return 0
	}
	bytesWritten, my.err = file.Write(content)
	if my.err != nil {
		my.err = WriteFileErr.Wrap(my.err)
		return 0
	}

	return bytesWritten
}

// Write 写入文件
func (my *File) Write(content []byte, mode os.FileMode, operations ...int) *File {
	my.mu.Lock()
	defer my.mu.Unlock()
	defer my.refresh()

	_ = my.write(content, mode, operations...)

	return my
}

// copyTo 复制文件
func (my *File) copyTo(dstFilename string) int64 {
	var (
		src, dst     *os.File
		bytesWritten int64
	)
	src, my.err = os.Open(my.fullPath)
	defer func() { _ = src.Close() }()
	if my.err != nil {
		my.err = CopyFileSrcErr.Wrap(my.err)
		return bytesWritten
	}

	dstFile := FileApp.NewByAbs(dstFilename)
	dstDir := DirApp.NewByAbs(dstFile.GetBasePath())
	if !dstDir.GetExist() {
		dstDir.CreateDefaultMode()
	}
	dstFile.Create(my.getMode())
	if dstFile.Error() != nil {
		my.err = CopyFileDstErr.Wrap(dstFile.Error())
		return bytesWritten
	}

	dst, my.err = os.Create(dstFile.getFullPath())
	defer func() { _ = dst.Close() }()
	if my.err != nil {
		my.err = CopyFileDstErr.Wrap(my.err)
		return bytesWritten
	}

	bytesWritten, my.err = io.Copy(dst, src)
	if my.err != nil {
		my.err = CopyFileErr.Wrap(my.err)
		return bytesWritten
	}

	return bytesWritten
}

// CopyTo 复制文件
func (my *File) CopyTo(dstFilename string) *File {
	my.mu.Lock()
	defer my.mu.Unlock()
	defer my.refresh()

	my.copyTo(dstFilename)

	return my
}

// read 读取文件
func (my *File) read() []byte {
	var (
		file    *os.File
		content []byte
	)

	file, my.err = os.OpenFile(my.getFullPath(), os.O_RDWR, 0666)
	if my.err != nil {
		my.err = ReadFileErr.Wrap(my.err)
		return nil
	}
	defer func() { _ = file.Close() }()
	content, my.err = io.ReadAll(file)
	if my.err != nil {
		my.err = ReadFileErr.Wrap(my.err)
		return nil
	}

	return content
}

// Read 读取文件
func (my *File) Read() []byte {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.read()
}
