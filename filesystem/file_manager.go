package filesystem

import (
	"errors"
	"fmt"

	"github.com/aid297/aid/httpClient"
)

type (
	FileManager struct {
		Err            error
		dstDir, srcDir string
		fileBytes      []byte
		fileSize       int64
		config         *FileManagerConfig
		uploadMethods  map[FileManagerConfigDriver]func() (int64, error)
		deleteMethods  map[FileManagerConfigDriver]func() error
	}
	FileManagerConfigDriver string
	FileManagerConfig       struct {
		Username  string
		Password  string
		AuthTitle string
		Driver    FileManagerConfigDriver
	}
)

var FileManagerApp FileManager

func (*FileManager) New(config *FileManagerConfig) *FileManager { return NewFileManager(config) }
func (*FileManager) NewByLocalFile(srcDir, dstDir string, config *FileManagerConfig) (*FileManager, error) {
	return NewFileManagerByLocalFile(srcDir, dstDir, config)
}
func (*FileManager) NewByBytes(srcFileBytes []byte, dstDir string, config *FileManagerConfig) *FileManager {
	return NewFileManagerByBytes(srcFileBytes, dstDir, config)
}

const (
	FileManagerConfigDriverLocal FileManagerConfigDriver = "LOCAL"
	FileManagerConfigDriverNexus FileManagerConfigDriver = "NEXUS"
	FileManagerConfigDriverOss   FileManagerConfigDriver = "OSS"
)

// NewFileManager 初始化：文件管理
//
//go:fix 推荐使用 New方法
func NewFileManager(config *FileManagerConfig) *FileManager {
	fm := &FileManager{
		config:        config,
		uploadMethods: make(map[FileManagerConfigDriver]func() (int64, error)),
		deleteMethods: make(map[FileManagerConfigDriver]func() error),
	}

	fm.uploadMethods = map[FileManagerConfigDriver]func() (int64, error){
		FileManagerConfigDriverLocal: fm.uploadToLocal,
		FileManagerConfigDriverNexus: fm.uploadToNexus,
		FileManagerConfigDriverOss:   fm.uploadToOss,
	}

	fm.deleteMethods = map[FileManagerConfigDriver]func() error{
		FileManagerConfigDriverLocal: fm.deleteFromLocal,
		FileManagerConfigDriverNexus: fm.deleteFromNexus,
		FileManagerConfigDriverOss:   fm.deleteFromOss,
	}

	return fm
}

// NewFileManagerByLocalFile 初始化：文件管理器（通过本地文件）
//
//go:fix 推荐使用NewByLocalFile方法
func NewFileManagerByLocalFile(srcDir, dstDir string, config *FileManagerConfig) (*FileManager, error) {
	fs := NewFileSystemByAbsolute(srcDir)
	if !fs.IsExist {
		return nil, errors.New("目标文件不存在")
	}

	fileBytes, err := fs.Read()
	if err != nil {
		return nil, err
	}

	return &FileManager{dstDir: dstDir, srcDir: srcDir, fileBytes: fileBytes, fileSize: int64(len(fileBytes)), config: config}, nil
}

// NewFileManagerByBytes 实例化：文件管理器（通过字节）
//
//go:fix 推荐使用NewByBytes方法
func NewFileManagerByBytes(srcFileBytes []byte, dstDir string, config *FileManagerConfig) *FileManager {
	return &FileManager{dstDir: dstDir, fileBytes: srcFileBytes, fileSize: int64(len(srcFileBytes)), config: config}
}

// SetSrcDir 设置源文件
func (my *FileManager) SetSrcDir(srcDir string) (*FileManager, error) {
	fs := NewFileSystemByAbsolute(srcDir)
	if !fs.IsExist {
		return nil, errors.New("目标文件不存在")
	}

	fileBytes, err := fs.Read()
	if err != nil {
		return nil, err
	}
	{
		my.srcDir = fs.GetDir()
		my.fileBytes = fileBytes
		my.fileSize = int64(len(fileBytes))
	}

	return my, nil
}

// SetDstDir 设置目标目录
func (my *FileManager) SetDstDir(dstDir string) *FileManager {
	my.dstDir = dstDir

	return my
}

// Upload 上传文件
func (my *FileManager) Upload() (int64, error) {
	switch my.config.Driver {
	case FileManagerConfigDriverLocal, FileManagerConfigDriverNexus, FileManagerConfigDriverOss:
		return my.uploadMethods[my.config.Driver]()
	default:
		return 0, fmt.Errorf("不支持的驱动类型：%s", my.config.Driver)
	}
}

// Delete 删除文件
func (my *FileManager) Delete() error {
	switch my.config.Driver {
	case FileManagerConfigDriverLocal, FileManagerConfigDriverNexus, FileManagerConfigDriverOss:
		return my.deleteMethods[my.config.Driver]()
	default:
		return fmt.Errorf("不支持的驱动类型：%s", my.config.Driver)
	}
}

// 上传到本地
func (my *FileManager) uploadToLocal() (int64, error) {
	dst := FileSystemApp.NewByAbs(my.dstDir)

	return dst.WriteBytes(my.fileBytes)
}

// 上传到nexus
func (my *FileManager) uploadToNexus() (int64, error) {
	client := httpClient.
		NewPut(my.dstDir).
		SetAuthorization(my.config.Username, my.config.Password, my.config.AuthTitle).
		SetHeaders(map[string][]string{"Content-Length": {fmt.Sprintf("%d", my.fileSize)}}).
		SetBody(my.fileBytes).
		Send()

	if client.Err != nil {
		return 0, client.Err
	}

	return int64(len(my.fileBytes)), nil
}

// 上传到oss
func (my *FileManager) uploadToOss() (int64, error) { return 0, errors.New("暂不支持oss方式") }

// 从本地删除文件
func (my *FileManager) deleteFromLocal() error { return NewFileSystemByAbsolute(my.dstDir).DelFile() }

// 从nexus删除文件
func (my *FileManager) deleteFromNexus() error {
	return httpClient.
		NewDelete(my.dstDir).
		SetAuthorization(my.config.Username, my.config.Password, my.config.AuthTitle).
		Send().Err
}

// 从oss删除文件
func (my *FileManager) deleteFromOss() error { return errors.New("暂不支持oss方式") }
