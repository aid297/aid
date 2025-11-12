package filesystem

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

type (
	// FileSystem 文件系统
	FileSystem struct {
		dir     string
		IsExist bool
		IsDir   bool
		IsFile  bool
	}

	// FileSystemCopyFilesTarget 拷贝文件目标
	FileSystemCopyFilesTarget struct {
		Src         *FileSystem
		DstFilename string
	}
)

var FileSystemApp FileSystem

func (*FileSystem) NewByRel(dir string) *FileSystem { return FileSystemApp.NewByRelative(dir) }
func (*FileSystem) NewByAbs(dir string) *FileSystem { return FileSystemApp.NewByAbsolute(dir) }

// NewByRelative 实例化：文件系统（相对路径）
func (*FileSystem) NewByRelative(dir string) *FileSystem {
	ins := &FileSystem{dir: filepath.Clean(filepath.Join(getRootPath(), dir))}

	return ins.init()
}

// NewByAbsolute 实例化：文件系统（绝对路径）
func (*FileSystem) NewByAbsolute(dir string) *FileSystem {
	ins := &FileSystem{dir: dir}

	return ins.init()
}

// NewFileSystemByRelative 实例化：文件系统（相对路径）
//
//go:fix 推荐使用NewByRelative方法
func NewFileSystemByRelative(dir string) *FileSystem {
	ins := &FileSystem{dir: filepath.Clean(filepath.Join(getRootPath(), dir))}

	return ins.init()
}

// NewFileSystemByAbsolute 实例化：文件系统（绝对路径）
//
//go:fix 推荐使用NewByAbsolute方法
func NewFileSystemByAbsolute(dir string) *FileSystem {
	ins := &FileSystem{dir: dir}

	return ins.init()
}

// Copy 复制一个新的对象
func (my *FileSystem) Copy() *FileSystem {
	copied := *my

	return &copied
}

// SetDirByRelative 设置路径：相对路径
func (my *FileSystem) SetDirByRelative(dir string) *FileSystem {
	my.dir = filepath.Clean(filepath.Join(getRootPath(), dir))
	my.init()

	return my
}

// SetDirByAbs 设置路径：绝对路径
func (my *FileSystem) SetDirByAbs(dir string) *FileSystem {
	my.dir = dir
	my.init()

	return my
}

func (my *FileSystem) Join(dir string) *FileSystem {
	my.dir = filepath.Join(my.dir, dir)
	my.init()

	return my
}

// Joins 增加若干路径
func (my *FileSystem) Joins(dir ...string) *FileSystem {
	for _, v := range dir {
		my.Join(v)
	}

	my.init()

	return my
}

func getRootPath() string {
	rootPath, _ := filepath.Abs(".")

	return rootPath
}

// getCurrentPath 最终方案-全兼容
func getCurrentPath(paths ...string) string {
	dir := getGoBuildPath()

	if strings.Contains(dir, getTmpDir()) {
		return getGoRunPath()
	}

	return dir
}

// 获取系统临时目录，兼容go run
func getTmpDir() string {
	dir := os.Getenv("TEMP")

	if dir == "" {
		dir = os.Getenv("TMP")
	}
	res, _ := filepath.EvalSymlinks(dir)

	return res
}

// 获取当前执行文件绝对路径
func getGoBuildPath() string {
	exePath, err := os.Executable()

	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))

	return res
}

// 获取当前执行文件绝对路径（go run）
func getGoRunPath() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)

	if ok {
		abPath = path.Dir(filename)
	}

	return abPath
}

// 初始化
func (my *FileSystem) init() *FileSystem {
	var e error
	my.IsExist, e = my.Exist() // 检查文件是否存在
	if e != nil {
		panic(fmt.Errorf("检查路径错误：%s", e.Error()))
	}

	if my.IsExist {
		e = my.CheckPathType() // 检查路径类型
		if e != nil {
			panic(fmt.Errorf("检查路径类型错误：%s", e.Error()))
		}
	}

	return my
}

// Exist 检查文件是否存在
func (my *FileSystem) Exist() (bool, error) {
	_, err := os.Stat(my.dir)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// IAmDir 检查当前路径是否是目录
func (my *FileSystem) IAmDir() (bool, error) {
	_, err := my.Exist()
	if err != nil {
		return false, err
	}

	info, e := os.Stat(my.dir)
	if e != nil {
		return false, e
	}

	return info.IsDir(), nil
}

// IAmFile 检查当前路径是否是文件
func (my *FileSystem) IAmFile() (bool, error) {
	_, err := my.Exist()
	if err != nil {
		return false, err
	}

	info, e := os.Stat(my.dir)
	if e != nil {
		return false, e
	}

	return !info.IsDir(), nil
}

// MkDir 创建文件夹
func (my *FileSystem) MkDir() error {
	if !my.IsExist {
		if e := os.MkdirAll(my.dir, os.ModePerm); e != nil {
			return e
		}
	}

	return nil
}

// GetDir 获取当前路径
func (my *FileSystem) GetDir() string { return my.dir }

// CheckPathType 判断一个路径是文件还是文件夹
func (my *FileSystem) CheckPathType() error {
	info, e := os.Stat(my.dir)
	if e != nil {
		return e
	}

	if info.IsDir() {
		my.IsDir = true
		my.IsFile = !my.IsDir
	} else {
		my.IsFile = true
		my.IsDir = !my.IsFile
	}

	return nil
}

// Delete 删除文件或文件夹
func (my *FileSystem) Delete() error {
	if my.IsExist {
		if my.IsDir {
			return my.DelDir()
		}
		if my.IsFile {
			return my.DelFile()
		}
	}

	return nil
}

// DelDir 删除文件夹
func (my *FileSystem) DelDir() error {
	err := os.RemoveAll(my.dir)
	if err != nil {
		return err
	}

	return nil
}

// DelFile 删除文件
func (my *FileSystem) DelFile() error {
	e := os.Remove(my.dir)
	if e != nil {
		return e
	}

	return nil
}

// Read 读取文件
func (my *FileSystem) Read() ([]byte, error) { return os.ReadFile(my.dir) }

// RenameFile 修改文件名并获取新的文件对象
func (my *FileSystem) RenameFile(newFilename string, deleteRepetition bool) (*FileSystem, error) {
	dir, _ := filepath.Split(my.GetDir())
	dst := FileSystemApp.NewByAbsolute(path.Join(dir, newFilename))

	if deleteRepetition {
		if dst.IsExist {
			if err := dst.DelFile(); err != nil {
				return nil, err
			}
		}
	}

	if err := os.Rename(my.GetDir(), dst.GetDir()); err != nil {
		return nil, err
	}

	return dst, nil
}

// CopyFile 拷贝单文件
func (my *FileSystem) CopyFile(dstDir, dstFilename string, abs bool) (string, error) {
	var (
		err         error
		srcFile     *os.File
		srcFilename string
		dst         *FileSystem
	)

	// 如果是相对路径
	if !abs {
		dst = FileSystemApp.NewByRelative(dstDir)
	} else {
		dst = FileSystemApp.NewByAbsolute(dstDir)
	}
	// 创建目标文件夹
	if !dst.IsDir {
		if err = dst.MkDir(); err != nil {
			return "", err
		}
	}

	// 判断源是否是文件
	if !my.IsFile {
		return "", fmt.Errorf("源文件不存在：%s", my.GetDir())
	}

	// 打开源文件
	srcFile, err = os.Open(my.GetDir())
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	if dstFilename == "" {
		srcFilename = filepath.Base(my.GetDir())
		dst.Join(srcFilename)
	} else {
		dst.Join(dstFilename)
	}

	// 创建目标文件
	dstFile, err := os.Create(dst.GetDir())
	if err != nil {
		return "", err
	}
	defer dstFile.Close()

	// 拷贝内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return "", err
	}

	// 确保所有内容都已写入磁盘
	err = dstFile.Sync()
	if err != nil {
		return "", err
	}

	return dst.GetDir(), nil
}

// CopyFiles 拷贝多个文件
func copyFiles(srcFiles []*FileSystemCopyFilesTarget, dstDir string, abs bool) error {
	var (
		err error
		dst *FileSystem
	)

	if abs {
		dst = FileSystemApp.NewByAbsolute(dstDir)
	} else {
		dst = FileSystemApp.NewByRelative(dstDir)
	}

	if !dst.IsDir {
		if err = dst.MkDir(); err != nil {
			return err
		}
	}

	for _, srcFile := range srcFiles {
		// 获取源文件名
		srcFilename := filepath.Base(srcFile.Src.GetDir())

		// 拷贝文件
		if srcFile.DstFilename != "" {
			_, err = srcFile.Src.CopyFile(dst.GetDir(), srcFile.DstFilename, true)
		} else {
			_, err = srcFile.Src.CopyFile(dst.GetDir(), srcFilename, true)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// CopyDir 拷贝目录
func (my *FileSystem) CopyDir(dstDir string, abs bool) error {
	// 判断是否是目录
	if !my.IsDir {
		return errors.New("源目录不存在")
	}

	// 遍历源目录
	if err := filepath.Walk(my.GetDir(), func(srcPath string, info os.FileInfo, err error) error {
		var (
			src         *FileSystem
			dst         *FileSystem
			srcFilename string
		)

		if abs {
			dst = FileSystemApp.NewByAbsolute(dstDir)
		} else {
			dst = FileSystemApp.NewByRelative(dstDir)
		}

		if !dst.IsDir {
			if err = dst.MkDir(); err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}

		srcFilename = filepath.Base(srcPath)
		src = FileSystemApp.NewByAbsolute(srcPath)

		if src.IsFile {
			if _, err = src.CopyFile(dst.GetDir(), srcFilename, true); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// WriteBytes 写入文件：bytes
func (my *FileSystem) WriteBytes(content []byte) (int64, error) {
	var written int
	// 打开文件
	file, err := os.OpenFile(my.GetDir(), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return 0, err
	}
	defer func(file *os.File) { _ = file.Close() }(file)

	// 写入内容
	written, err = file.Write(content)
	if err != nil {
		return 0, err
	}

	return int64(written), nil
}

// WriteString 写入文件：string
func (my *FileSystem) WriteString(content string) (int64, error) {
	return my.WriteBytes([]byte(content))
}

// WriteIoReader 写入文件：io.Reader
func (my *FileSystem) WriteIoReader(content io.Reader) (written int64, err error) {
	dst, err := os.Create(my.dir)
	if err != nil {
		return 0, err
	}
	defer func(dst *os.File) { _ = dst.Close() }(dst)

	return io.Copy(dst, content)
}

// WriteBytesAppend 追加写入文件：bytes
func (my *FileSystem) WriteBytesAppend(content []byte) (int64, error) {
	var written int
	// Open the file in append mode.
	file, e := os.OpenFile(my.GetDir(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if e != nil {
		return 0, e
	}
	defer func(file *os.File) { _ = file.Close() }(file)

	// 追加写入内容
	written, e = file.Write(content)
	if e != nil {
		return 0, e
	}

	return int64(written), nil
}

// WriteStringAppend 追加写入文件：string
func (my *FileSystem) WriteStringAppend(content string) (int64, error) {
	return my.WriteBytesAppend([]byte(content))
}

// WriteIoReaderAppend 追加写入文件：io.Reader
func (my *FileSystem) WriteIoReaderAppend(content io.Reader) (int64, error) {
	var written int
	dst, err := os.OpenFile(my.dir, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}
	defer func(dst *os.File) { _ = dst.Close() }(dst)

	c, err := io.ReadAll(content)
	if err != nil {
		return 0, err
	}

	written, err = dst.Write(c)
	if err != nil {
		return 0, err
	}

	return int64(written), nil
}
