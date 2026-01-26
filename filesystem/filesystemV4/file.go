package filesystemV4

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/aid297/aid/operation/operationV2"
)

type (
	File struct {
		Error    error        `json:"error"`    // 错误信息
		Name     string       `json:"name"`     // 文件名
		BasePath string       `json:"basePath"` // 基础路径
		FullPath string       `json:"fullPath"` // 完整路径
		Size     int64        `json:"size"`     // 文件大小
		Info     os.FileInfo  `json:"info"`     // 文件信息
		Mode     os.FileMode  `json:"mode"`     // 文件权限
		Exist    bool         `json:"exist"`    // 文件是否存在
		mu       sync.RWMutex // 读写锁
		Ext      string       `json:"extension"` // 文件扩展名
		FileInfo os.FileInfo  `json:"fileInfo"`  // 文件信息
		Mime     string       `json:"mime"`      // 文件 Mime 类型
		Kind     string       `json:"kind"`      // 类型
	}
)

var (
	DefaultCreateMode = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	DefaultReadMode   = os.O_RDWR
)

func NewFile(attrs ...PathAttributer) Filesystemer {
	return (&File{mu: sync.RWMutex{}}).SetAttrs(attrs...).refresh()
}

func (my *File) GetExist() bool           { return my.Exist }
func (my *File) GetError() error          { return my.Error }
func (my *File) GetFullPath() string      { return my.FullPath }
func (my *File) GetInfo() os.FileInfo     { return my.Info }
func (my *File) GetDirs() []Filesystemer  { return nil }
func (my *File) GetFiles() []Filesystemer { return nil }

func (my *File) SetAttrs(attrs ...PathAttributer) Filesystemer {
	my.mu.Lock()
	defer my.mu.Unlock()
	for idx := range attrs {
		attrs[idx].Register(my)
	}
	return my
}

func (my *File) SetFullPathForAttr(path string) Filesystemer { my.FullPath = path; return my }

func (my *File) SetFullPathByAttr(attrs ...PathAttributer) Filesystemer {
	return my.SetAttrs(attrs...).refresh()
}

// refresh 刷新文件信息
func (my *File) refresh() Filesystemer {
	var err error

	if my.FullPath != "" {
		if my.FileInfo, err = os.Stat(my.FullPath); err != nil {
			if os.IsNotExist(err) {
				my.Name = ""
				my.Size = 0
				my.Mode = 0
				my.BasePath = filepath.Dir(my.FullPath)
				my.Ext = filepath.Ext(my.FullPath)
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
func (my *File) Lock() Filesystemer { my.mu.Lock(); return my }

// Unlock 解锁 → 写
func (my *File) Unlock() Filesystemer { my.mu.Unlock(); return my }

// RLock 加锁 → 读
func (my *File) RLock() Filesystemer { my.mu.RLock(); return my }

// RUnlock 解锁 → 读
func (my *File) RUnlock() Filesystemer { my.mu.RUnlock(); return my }

func (my *File) Join(paths ...string) Filesystemer {
	my.FullPath = filepath.Join(append([]string{my.FullPath}, paths...)...)
	return my.refresh()
}

// Create 创建文件
func (my *File) Create(attrs ...OperationAttributer) Filesystemer {
	if dir := NewDir(Abs(my.BasePath)); !dir.GetExist() {
		if err := dir.Create(attrs...).GetError(); err != nil {
			my.Error = fmt.Errorf("%w:%w", ErrCreateDir, err)
		}
	}

	my.Write(nil, attrs...)

	return my
}

// 向文件内写入内容
func (my *File) Write(content []byte, attrs ...OperationAttributer) Filesystemer {
	var (
		err       error
		operation = NewOperation(attrs...)
		file      *os.File
		dir       Filesystemer
	)

	if dir = NewDir(Abs(my.BasePath)); !dir.GetExist() {
		if err = dir.Create(attrs...).GetError(); err != nil {
			my.Error = fmt.Errorf("%w:%w", ErrCreateDir, err)
			return my
		}
	}

	if file, err = os.OpenFile(
		my.FullPath,
		operationV2.NewTernary(operationV2.TrueValue(operation.Flag), operationV2.FalseValue(DefaultCreateMode)).GetByValue(operation.Flag != 0),
		operationV2.NewTernary(operationV2.TrueValue(operation.Mode), operationV2.FalseValue(os.FileMode(0777))).GetByValue(operation.Mode != 0),
	); err != nil {
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
func (my *File) Rename(newName string) Filesystemer {
	var (
		err     error
		newFile = NewFile(Abs(my.BasePath, newName))
	)

	if err = os.Rename(my.FullPath, newFile.GetFullPath()); err != nil {
		my.Error = fmt.Errorf("%w：%w", ErrRename, err)
		return my
	}

	return newFile.refresh()
}

// Remove 删除文件
func (my *File) Remove() Filesystemer {
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

func (my *File) RemoveAll() Filesystemer { return my }

// Read 读取文件内容
func (my *File) Read(attrs ...OperationAttributer) ([]byte, error) {
	var (
		fileOperation = NewOperation(attrs...)
		file          *os.File
		content       []byte
		err           error
	)

	if file, err = os.OpenFile(
		my.FullPath,
		operationV2.NewTernary(operationV2.TrueFn(func() int { return fileOperation.Flag }), operationV2.FalseFn(func() int { return DefaultReadMode })).GetByValue(fileOperation.Flag != 0),
		operationV2.NewTernary(operationV2.TrueFn(func() os.FileMode { return fileOperation.Mode }), operationV2.FalseFn(func() os.FileMode { return os.FileMode(0777) })).GetByValue(fileOperation.Mode != 0),
	); err != nil {
		return []byte{}, fmt.Errorf("%w:%w", ErrReadFile, err)
	}
	defer func() { _ = file.Close() }()
	if content, err = io.ReadAll(file); err != nil {
		return []byte{}, fmt.Errorf("%w:%w", ErrReadFile, err)
	}

	return content, nil
}

// CopyTo 复制文件到指定路径
func (my *File) CopyTo(isRel bool, dstPaths ...string) Filesystemer {
	if my.FullPath == "" {
		my.Error = ErrMissFullPath
		return my
	}

	if !my.Exist {
		my.Error = ErrFileNotExist
		return my
	}

	a := NewDir(Rel(dstPaths...)).Up()
	b := NewDir(Abs(my.BasePath))

	if my.Error = a.Create(Mode(b.GetInfo().Mode())).GetError(); my.Error != nil {
		return my
	}

	my.Error = copyFileTo(my.FullPath, NewFile(operationV2.NewTernary(operationV2.TrueFn(func() PathAttributer { return Rel(dstPaths...) }), operationV2.FalseFn(func() PathAttributer { return Abs(dstPaths...) })).GetByValue(isRel)).GetFullPath())

	return my
}

func (my *File) Copy() Filesystemer { return NewFile(Abs(my.FullPath)) }

func (my *File) Up() Filesystemer { return my }

func (my *File) LS() Filesystemer { return my }
