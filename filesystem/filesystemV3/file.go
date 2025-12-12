package filesystemV3

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/aid297/aid/operation/operationV2"
)

type File struct {
	IsRel     bool         // 是否使用相对路径
	Error     error        // 错误信息
	Name      string       // 文件名
	BasePath  string       // 基础路径
	FullPath  string       // 完整路径
	Size      int64        // 文件大小
	Info      os.FileInfo  // 文件信息
	Mode      os.FileMode  // 文件权限
	Exist     bool         // 文件是否存在
	mu        sync.RWMutex // 读写锁
	Extension string       // 文件扩展名
	Fileinfo  os.FileInfo  // 文件信息
	Mime      string       // 文件Mime类型
}

var (
	DefaultCreateMode = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	DefaultReadMode   = os.O_RDWR
)

// NewFile 实例化
func NewFile(attrs ...FileAttributer) *File {
	return (&File{mu: sync.RWMutex{}}).setAttrs(attrs...).refresh()
}

// NewFileAbs 实例化：绝对路径
func NewFileAbs(attrs ...FileAttributer) *File {
	return NewFile(append(attrs, APP.FileAttr.IsRel.SetAbs())...)
}

// NewFileRel 实例化：相对路径
func NewFileRel(attrs ...FileAttributer) *File {
	return NewFile(append([]FileAttributer{APP.FileAttr.IsRel.SetRel()}, attrs...)...)
}

// New 实例化
func (*File) New(attrs ...FileAttributer) *File {
	return (&File{mu: sync.RWMutex{}}).setAttrs(attrs...).refresh()
}

// Abs 实例化：绝对路径
func (*File) Abs(attrs ...FileAttributer) *File {
	return APP.File.New(append(attrs, APP.FileAttr.IsRel.SetAbs())...)
}

// Rel 实例化：相对路径
func (*File) Rel(attrs ...FileAttributer) *File {
	return APP.File.New(append([]FileAttributer{APP.FileAttr.IsRel.SetRel()}, attrs...)...)
}

// SetAttrs 设置属性
func (my *File) SetAttrs(attrs ...FileAttributer) *File {
	my.mu.Lock()
	defer my.mu.Unlock()
	return my.setAttrs(attrs...)
}

// setAttrs 设置属性
func (my *File) setAttrs(attrs ...FileAttributer) *File {
	for idx := range attrs {
		attrs[idx].Register(my)
	}
	return my
}

// refresh 刷新文件信息
func (my *File) refresh() *File {
	var err error

	if my.FullPath != "" {
		if my.Fileinfo, err = os.Stat(my.FullPath); err != nil {
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

		my.Name = my.Fileinfo.Name()
		my.Size = my.Fileinfo.Size()
		my.Mode = my.Fileinfo.Mode()
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
func (my *File) Lock() *File {
	my.mu.Lock()
	return my
}

// Unlock 解锁 → 写
func (my *File) Unlock() *File {
	my.mu.Unlock()
	return my
}

// RLock 加锁 → 读
func (my *File) RLock() *File {
	my.mu.RLock()
	return my
}

// RUnlock 解锁 → 读
func (my *File) RUnlock() *File {
	my.mu.RUnlock()
	return my
}

// Join 连接路径
func (my *File) Join(dirs ...string) *File {
	return my.setAttrs(APP.FileAttr.IsRel.SetAbs(), APP.FileAttr.Path.Set(append([]string{my.FullPath}, dirs...)...)).refresh()
}

// Create 创建文件
func (my *File) Create(attrs ...FileOperationAttributer) *File {
	if dir := NewDirAbs(APP.DirAttr.Path.Set(my.BasePath)); !dir.Exist {
		if err := dir.Create().Error; err != nil {
			my.Error = fmt.Errorf("%w:%w", ErrCreateDir, err)
		}
	}

	my.Write(nil, attrs...)

	return nil
}

// 向文件内写入内容
func (my *File) Write(content []byte, attrs ...FileOperationAttributer) *File {
	var (
		err           error
		fileOperation = new(FileOperation).SetAttrs(attrs...)
		file          *os.File
	)

	if dir := NewDirAbs(APP.DirAttr.Path.Set(my.BasePath)); !dir.Exist {
		if err := dir.Create().Error; err != nil {
			my.Error = fmt.Errorf("%w:%w", ErrCreateDir, err)
		}
	}

	if file, err = os.OpenFile(
		my.FullPath,
		operationV2.NewTernary(operationV2.TrueValue(fileOperation.FileFlag), operationV2.FalseValue(DefaultCreateMode)).GetByValue(fileOperation.FileFlag != 0),
		operationV2.NewTernary(operationV2.TrueValue(fileOperation.FileMode), operationV2.FalseValue(os.FileMode(0777))).GetByValue(fileOperation.FileMode != 0),
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
func (my *File) Rename(newName string) *File {
	var (
		err     error
		newFile = NewFile(APP.FileAttr.IsRel.SetAbs(), APP.FileAttr.Path.Set(my.BasePath, newName))
	)

	if err = os.Rename(my.FullPath, newFile.FullPath); err != nil {
		my.Error = fmt.Errorf("%w：%w", ErrRename, err)
		return my
	}

	return newFile.refresh()
}

// Remove 删除文件
func (my *File) Remove() *File {
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
func (my *File) Read(attrs ...FileOperationAttributer) ([]byte, error) {
	var (
		fileOperation = new(FileOperation).SetAttrs(attrs...)
		file          *os.File
		content       []byte
		err           error
	)

	if file, err = os.OpenFile(
		my.FullPath,
		operationV2.NewTernary(operationV2.TrueFn(func() int { return fileOperation.FileFlag }), operationV2.FalseFn(func() int { return DefaultReadMode })).GetByValue(fileOperation.FileFlag != 0),
		operationV2.NewTernary(operationV2.TrueFn(func() os.FileMode { return fileOperation.FileMode }), operationV2.FalseFn(func() os.FileMode { return os.FileMode(0777) })).GetByValue(fileOperation.FileMode != 0),
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
func (my *File) CopyTo(isRel bool, dstPaths ...string) *File {
	if my.FullPath == "" {
		my.Error = ErrMissFullPath
		return my
	}

	if !my.Exist {
		my.Error = ErrFileNotExist
		return my
	}

	a := NewDir(APP.DirAttr.IsRel.SetRel(), APP.DirAttr.Path.Set(dstPaths...)).Up()
	print(a.FullPath)

	b := NewDirAbs(APP.DirAttr.Path.Set(my.BasePath))
	print(b.FullPath)

	if my.Error = NewDir(APP.DirAttr.IsRel.SetRel(), APP.DirAttr.Path.Set(dstPaths...)).Up().Create(DirMode(NewDirAbs(APP.DirAttr.Path.Set(my.BasePath)).Info.Mode())).Error; my.Error != nil {
		return my
	}

	dst := NewFile(APP.FileAttr.IsRel.Set(isRel), APP.FileAttr.Path.Set(dstPaths...)).FullPath
	my.Error = copyFileTo(my.FullPath, dst)

	return my
}

// Copy 复制文件实例
func (my *File) Copy() *File { return NewFileAbs(APP.FileAttr.Path.Set(my.FullPath)) }
