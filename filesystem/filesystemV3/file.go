package filesystemV3

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/aid297/aid/operation/operationV2"
)

type (
	File struct {
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
)

var (
	DefaultCreateMode = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	DefaultReadMode   = os.O_RDWR
)

func NewFile(attrs ...FileAttributer) *File {
	file := &File{mu: sync.RWMutex{}}
	return file.SetAttrs(attrs...).refresh()
}

func (my *File) SetAttrs(attrs ...FileAttributer) *File {
	for idx := range attrs {
		attrs[idx].Register(my)
	}
	return my
}

func (my *File) refresh() *File {
	var err error

	if my.FullPath != "" {
		if my.Fileinfo, err = os.Stat(my.FullPath); err != nil {
			if os.IsNotExist(my.Error) {
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
		my.Exist = true
		my.Error = nil
	} else {
		my.Error = fmt.Errorf("%w:%w", ErrMissFullPath, err)
		return my
	}

	return my
}

func (my *File) Lock() *File {
	my.mu.Lock()
	return my
}

func (my *File) Unlock() *File {
	my.mu.Unlock()
	return my
}

func (my *File) RLock() *File {
	my.mu.RLock()
	return my
}

func (my *File) RUnlock() *File {
	my.mu.RUnlock()
	return my
}

func (my *File) Join(dirs ...string) *File {
	allDirs := []string{my.FullPath}
	if len(dirs) == 0 {
		return my
	}
	allDirs = append(allDirs, dirs...)

	return NewFile(FileIsAbs(), FilePath(allDirs...))
}

func (my *File) Create(attrs ...FileOperationAttributer) *File {
	if dir := NewDir(DirIsAbs(), DirPath(my.BasePath)); !dir.Exist {
		if err := dir.Create().Error; err != nil {
			my.Error = fmt.Errorf("%w:%w", ErrCreateDir, err)
		}
	}

	my.Write(nil, attrs...)

	return nil
}

func (my *File) Write(content []byte, attrs ...FileOperationAttributer) *File {
	var (
		err           error
		fileOperation = new(FileOperation).SetAttrs(attrs...)
		file          *os.File
	)

	if dir := NewDir(DirIsAbs(), DirPath(my.BasePath)); !dir.Exist {
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

	return my
}

func (my *File) Rename(newName string) *File {
	var (
		err     error
		newFile = NewFile(FileIsAbs(), FilePath(my.BasePath, newName))
	)

	if err = os.Rename(my.FullPath, newFile.FullPath); err != nil {
		my.Error = fmt.Errorf("%w：%w", ErrRename, err)
		return my
	}

	return newFile
}

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

func (my *File) CopyTo(dstPath string) *File {
	var (
		err      error
		src, dst *os.File
	)
	if my.FullPath == "" {
		my.Error = ErrMissFullPath
		return my
	}

	if !my.Exist {
		my.Error = ErrFileNotExist
		return my
	}

	if src, err = os.Open(my.FullPath); err != nil {
		my.Error = fmt.Errorf("%w:%w", ErrOpenFile, err)
		return my
	}
	defer func() { _ = src.Close() }()

	dstFile := NewFile(FileIsAbs(), FilePath(dstPath))
	if err = dstFile.Create(FileMode(my.Mode)).Error; err != nil {
		my.Error = fmt.Errorf("%w:%w", ErrCreateFile, err)
		return my
	}

	if dst, err = os.Create(dstFile.FullPath); err != nil {
		my.Error = fmt.Errorf("%w:%w", ErrOpenFile, err)
		return my
	}
	defer func() { _ = dst.Close() }()

	if _, err = io.Copy(dst, src); err != nil {
		my.Error = fmt.Errorf("%w:%w", ErrWriteFile, err)
		return my
	}

	return my
}

// ******************** 管理器属性 ******************** //
type (
	FileAttributer interface{ Register(dir *File) }
	AttrFilePath   struct{ dirs []string }
	AttrFileIsRel  struct{ isRel bool }
)

func FilePath(dirs ...string) AttrFilePath { return AttrFilePath{dirs: dirs} }
func (my AttrFilePath) Register(dir *File) {
	dir.FullPath = operationV2.NewTernary(
		operationV2.TrueFn(func() string { return getRootPath(filepath.Join(my.dirs...)) }),
		operationV2.FalseFn(func() string { return filepath.Join(my.dirs...) }),
	).GetByValue(dir.IsRel)
}

func FileSetRel(isRel bool) AttrDirIsRel    { return AttrDirIsRel{isRel: isRel} }
func FileIsAbs() AttrFileIsRel              { return AttrFileIsRel{isRel: false} }
func FileIsRel() AttrFileIsRel              { return AttrFileIsRel{isRel: true} }
func (my AttrFileIsRel) Register(dir *File) { dir.IsRel = my.isRel }
