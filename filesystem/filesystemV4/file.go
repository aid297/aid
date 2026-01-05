package filesystemV4

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sync"
)

type File struct {
	Error     error         // 错误信息
	Name      string        // 文件名
	BasePath  string        // 基础路径
	FullPath  string        // 完整路径
	Size      int64         // 文件大小
	Info      os.FileInfo   // 文件信息
	Mode      os.FileMode   // 文件权限
	Exist     bool          // 文件是否存在
	mu        *sync.RWMutex // 读写锁
	Extension string        // 文件扩展名
	FileInfo  os.FileInfo   // 文件信息
	Mime      string        // 文件 Mime 类型
}

func (File) Abs(dirs ...string) File {
	return File{mu: &sync.RWMutex{}, FullPath: getRootPath(path.Join(dirs...))}.refresh()
}

func (my File) Rel(dirs ...string) File {
	return File{mu: &sync.RWMutex{}, FullPath: path.Join(dirs...)}.refresh()
}

// refresh 刷新文件信息
func (my File) refresh() File {
	var err error

	if my.FullPath != "" {
		if my.FileInfo, err = os.Stat(my.FullPath); err != nil {
			if os.IsNotExist(err) {
				my.Name = ""
				my.Size = 0
				my.Mode = 0
				my.BasePath = filepath.Dir(my.FullPath)
				my.Extension = filepath.Ext(my.FullPath)
				my.Exist = false
				my.Error = nil
				return my
			} else {
				my.Error = fmt.Errorf("%w:%w", ErrInit, err)
				return my
			}
		}

		my.Name = my.FileInfo.Name()
		my.Size = my.FileInfo.Size()
		my.Mode = my.FileInfo.Mode()
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
func (my File) Lock() File {
	my.mu.Lock()
	return my
}

// Unlock 解锁 → 写
func (my File) Unlock() File {
	my.mu.Unlock()
	return my
}

// RLock 加锁 → 读
func (my File) RLock() File {
	my.mu.RLock()
	return my
}

// RUnlock 解锁 → 读
func (my File) RUnlock() File {
	my.mu.RUnlock()
	return my
}

func (my File) Join(dirs ...string) File {
	my.FullPath = path.Join(append([]string{my.FullPath}, dirs...)...)
	return my.refresh()
}

// Create 创建文件
func (my File) Create(attrs ...FileOperationAttributer) File {
	if dir := APP.Dir.Abs(my.BasePath); !dir.Exist {
		if err := dir.Create().Error; err != nil {
			my.Error = fmt.Errorf("%w:%w", ErrCreateDir, err)
		}
	}

	return my.Write(nil, attrs...)
}

// 向文件内写入内容
func (my File) Write(content []byte, attrs ...FileOperationAttributer) File {
	var (
		err           error
		fileOperation = APP.FileOperation.New(attrs...)
		file          *os.File
	)

	if dir := APP.Dir.Abs(my.BasePath); !dir.Exist {
		if err := dir.Create().Error; err != nil {
			my.Error = fmt.Errorf("%w:%w", ErrCreateDir, err)
		}
	}

	if file, err = os.OpenFile(my.FullPath, fileOperation.FileFlag, fileOperation.FileMode); err != nil {
		my.Error = fmt.Errorf("%w:%w", ErrWriteFile, err)
		return my
	}
	defer func() { _ = file.Close() }()

	if _, err = file.Write(content); err != nil {
		my.Error = fmt.Errorf("%w:%w", ErrWriteFile, err)
		return my
	}

	return my.refresh()
}

// Rename 重命名文件
func (my File) Rename(newName string) File {
	var (
		err     error
		newFile = APP.File.Abs(my.BasePath, newName)
	)

	if err = os.Rename(my.FullPath, newFile.FullPath); err != nil {
		my.Error = fmt.Errorf("%w：%w", ErrRename, err)
		return my
	}

	return newFile.refresh()
}

// Remove 删除文件
func (my File) Remove() File {
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

// Read 读取文件内容
func (my File) Read(attrs ...FileOperationAttributer) ([]byte, error) {
	var (
		fileOperation = APP.FileOperation.New(attrs...)
		file          *os.File
		content       []byte
		err           error
	)

	if file, err = os.OpenFile(my.FullPath, fileOperation.FileFlag, fileOperation.FileMode); err != nil {
		return []byte{}, fmt.Errorf("%w:%w", ErrReadFile, err)
	}
	defer func() { _ = file.Close() }()

	if content, err = io.ReadAll(file); err != nil {
		return []byte{}, fmt.Errorf("%w:%w", ErrReadFile, err)
	}

	return content, nil
}

// CopyTo 复制文件到指定路径
func (my File) CopyTo(isRel bool, dstPaths ...string) File {
	if my.FullPath == "" {
		my.Error = ErrMissFullPath
		return my
	}

	if !my.Exist {
		my.Error = ErrFileNotExist
		return my
	}

	a := APP.Dir.Rel(dstPaths...).Up()
	print(a.FullPath)

	b := APP.Dir.Abs(my.BasePath)
	print(b.FullPath)

	if my.Error = APP.Dir.Rel(dstPaths...).Up().Create(APP.DirOperAttr.Mode.Set(APP.Dir.Abs(my.BasePath).Info.Mode())).Error; my.Error != nil {
		return my
	}

	dst := APP.File.Rel(dstPaths...)
	my.Error = copyFileTo(my.FullPath, dst.FullPath)

	return my
}

// Copy 复制文件实例
func (my File) Copy() File { return APP.File.Abs(my.FullPath) }
