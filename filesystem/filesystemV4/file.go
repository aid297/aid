package filesystemV4

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/aid297/aid/operation/operationV2"
)

type File struct {
	mu       sync.RWMutex `json:"-"`                              // 读写锁
	Error    error        `json:"error" swaggertype:"string"`     // 错误信息
	Name     string       `json:"name" swaggertype:"string"`      // 文件名
	BasePath string       `json:"basePath" swaggertype:"string"`  // 基础路径
	FullPath string       `json:"fullPath" swaggertype:"string"`  // 完整路径
	Size     int64        `json:"size" swaggertype:"integer"`     // 文件大小
	Info     os.FileInfo  `json:"info" swaggertype:"string"`      // 文件信息
	Mode     os.FileMode  `json:"mode" swaggertype:"string"`      // 文件权限
	Exist    bool         `json:"exist" swaggertype:"boolean"`    // 文件是否存在
	Ext      string       `json:"extension" swaggertype:"string"` // 文件扩展名
	Mime     string       `json:"mime" swaggertype:"string"`      // 文件 Mime 类型
	Kind     string       `json:"kind" swaggertype:"string"`      // 类型
}

var (
	DefaultCreateMode = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	DefaultReadMode   = os.O_RDWR
)

func NewFile(attrs ...PathAttributer) IFilesystem {
	return (&File{mu: sync.RWMutex{}, Kind: "FILE"}).SetAttrs(attrs...).refresh()
}

func (my *File) GetName() string         { return my.Name }
func (my *File) GetExist() bool          { return my.Exist }
func (my *File) GetError() error         { return my.Error }
func (my *File) GetBasePath() string     { return my.BasePath }
func (my *File) GetFullPath() string     { return my.FullPath }
func (my *File) GetInfo() os.FileInfo    { return my.Info }
func (my *File) GetDirs() []IFilesystem  { return nil }
func (my *File) GetFiles() []IFilesystem { return nil }
func (my *File) GetKind() string         { return my.Kind }

func (my *File) SetAttrs(attrs ...PathAttributer) IFilesystem {
	my.mu.Lock()
	defer my.mu.Unlock()
	for idx := range attrs {
		attrs[idx].Register(my)
	}
	return my
}

func (my *File) SetFullPathForAttr(path string) IFilesystem { my.FullPath = path; return my }

func (my *File) SetFullPathByAttr(attrs ...PathAttributer) IFilesystem {
	return my.SetAttrs(attrs...).refresh()
}

// refresh 刷新文件信息
func (my *File) refresh() IFilesystem {
	var err error

	if my.FullPath != "" {
		if my.Info, err = os.Stat(my.FullPath); err != nil {
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
func (my *File) Lock() IFilesystem { my.mu.Lock(); return my }

// Unlock 解锁 → 写
func (my *File) Unlock() IFilesystem { my.mu.Unlock(); return my }

// RLock 加锁 → 读
func (my *File) RLock() IFilesystem { my.mu.RLock(); return my }

// RUnlock 解锁 → 读
func (my *File) RUnlock() IFilesystem { my.mu.RUnlock(); return my }

func (my *File) Join(paths ...string) IFilesystem {
	my.FullPath = filepath.Join(append([]string{my.FullPath}, paths...)...)
	return my.refresh()
}

// Create 创建文件
func (my *File) Create(attrs ...OperationAttributer) IFilesystem {
	if dir := NewDir(Abs(my.BasePath)); !dir.GetExist() {
		if err := dir.Create(attrs...).GetError(); err != nil {
			my.Error = fmt.Errorf("%w:%w", ErrCreateDir, err)
		}
	}

	my.Write(nil, attrs...)

	return my
}

// 向文件内写入内容
func (my *File) Write(content []byte, attrs ...OperationAttributer) IFilesystem {
	var (
		err       error
		operation = NewOperation(attrs...)
		file      *os.File
		dir       IFilesystem
	)

	if dir = NewDir(Abs(my.BasePath)); !dir.GetExist() {
		if err = dir.Create(attrs...).GetError(); err != nil {
			my.Error = fmt.Errorf("%w:%w", ErrCreateDir, err)
			return my
		}
	}

	if file, err = os.OpenFile(
		my.FullPath,
		operationV2.NewTernary(
			operationV2.TrueValue(operation.Flag),
			operationV2.FalseValue(DefaultCreateMode),
		).GetByValue(operation.Flag != 0),
		operationV2.NewTernary(
			operationV2.TrueValue(operation.Mode),
			operationV2.FalseValue(os.FileMode(0777)),
		).GetByValue(operation.Mode != 0),
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
func (my *File) Rename(newName string) IFilesystem {
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
func (my *File) Remove() IFilesystem {
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

func (my *File) RemoveAll() IFilesystem { return my.Remove() }

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
		operationV2.NewTernary(
			operationV2.TrueFn(func() int { return fileOperation.Flag }),
			operationV2.FalseFn(func() int { return DefaultReadMode }),
		).GetByValue(fileOperation.Flag != 0),
		operationV2.NewTernary(
			operationV2.TrueFn(func() os.FileMode { return fileOperation.Mode }),
			operationV2.FalseFn(func() os.FileMode { return os.FileMode(0777) }),
		).GetByValue(fileOperation.Mode != 0),
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
func (my *File) CopyTo(isRel bool, dstPaths ...string) IFilesystem {
	if my.FullPath == "" {
		my.Error = ErrMissFullPath
		return my
	}

	if !my.Exist {
		my.Error = ErrFileNotExist
		return my
	}

	if dstDir := NewDir(operationV2.NewTernary(operationV2.TrueValue(Rel(dstPaths...)), operationV2.FalseValue(Abs(dstPaths...))).GetByValue(isRel)).Up(); dstDir.GetError() != nil {
		my.Error = dstDir.GetError()
		return my
	}

	newPath := NewFile(
		operationV2.NewTernary(
			operationV2.TrueFn(func() PathAttributer { return Rel(dstPaths...) }),
			operationV2.FalseFn(func() PathAttributer { return Abs(dstPaths...) }),
		).GetByValue(isRel),
	)

	my.Error = copyFileTo(my.GetFullPath(), newPath.GetFullPath())

	return my
}

// Zip 压缩文件到 zip 格式
func (my *File) Zip() IFilesystem {
	var (
		err       error
		srcFile   *os.File
		zipFile   *os.File
		zipWriter *zip.Writer
		writer    io.Writer
		zipPath   string
		content   []byte
	)

	if my.FullPath == "" {
		my.Error = ErrMissFullPath
		return my
	}

	if !my.Exist {
		my.Error = ErrFileNotExist
		return my
	}

	// 压缩后的文件名：原文件名 + .zip
	zipPath = my.FullPath + ".zip"

	// 读取源文件内容
	if srcFile, err = os.Open(my.FullPath); err != nil {
		my.Error = fmt.Errorf("打开源文件失败:%w", err)
		return my
	}
	defer srcFile.Close()

	// 创建 zip 文件
	if zipFile, err = os.Create(zipPath); err != nil {
		my.Error = fmt.Errorf("创建 zip 文件失败:%w", err)
		return my
	}
	defer zipFile.Close()

	// 创建 zip writer
	zipWriter = zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 在 zip 中创建文件条目（使用原文件名）
	if writer, err = zipWriter.Create(my.Name); err != nil {
		my.Error = fmt.Errorf("创建 zip 条目失败:%w", err)
		return my
	}

	// 读取源文件内容
	if content, err = io.ReadAll(srcFile); err != nil {
		my.Error = fmt.Errorf("读取源文件失败:%w", err)
		return my
	}

	// 写入到 zip
	if _, err = writer.Write(content); err != nil {
		my.Error = fmt.Errorf("写入 zip 文件失败:%w", err)
		return my
	}

	// 返回指向压缩后文件的新 File 对象
	return NewFile(Abs(zipPath))
}

func (my *File) Copy() IFilesystem { return NewFile(Abs(my.FullPath)) }

func (my *File) Up() IFilesystem { return my }

func (my *File) LS() IFilesystem { return my }
