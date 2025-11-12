package filesystemV2

import (
	"fmt"
	"reflect"

	"github.com/aid297/aid/array"
	"github.com/aid297/aid/myError"
	"github.com/aid297/aid/operation"
)

type (
	FileInitError          struct{ myError.MyError }
	FileNotExistsError     struct{ myError.MyError }
	FileFullPathEmptyError struct{ myError.MyError }
	CreateFileError        struct{ myError.MyError }
	RenameFileError        struct{ myError.MyError }
	RemoveFileError        struct{ myError.MyError }
	PermissionFileError    struct{ myError.MyError }
	CopyFileError          struct{ myError.MyError }
	CopyFileSrcError       struct{ myError.MyError }
	CopyFileDstError       struct{ myError.MyError }
	WriteFileError         struct{ myError.MyError }
	ReadFileError          struct{ myError.MyError }

	DirInitError          struct{ myError.MyError }
	DirNotExistsError     struct{ myError.MyError }
	DirFullPathEmptyError struct{ myError.MyError }
	CreateDirError        struct{ myError.MyError }
	RenameDirError        struct{ myError.MyError }
	RemoveDirError        struct{ myError.MyError }
	PermissionDirError    struct{ myError.MyError }
	CopyDirError          struct{ myError.MyError }
	CopyDirSrcError       struct{ myError.MyError }
	CopyDirDstError       struct{ myError.MyError }
	WriteDirError         struct{ myError.MyError }
	ReadDirError          struct{ myError.MyError }
)

var (
	FileInitErr          FileInitError
	FileNotExistsErr     FileNotExistsError
	FileFullPathEmptyErr FileFullPathEmptyError
	CreateFileErr        CreateFileError
	RenameFileErr        RenameFileError
	RemoveFileErr        RemoveFileError
	PermissionFileErr    PermissionFileError
	CopyFileErr          CopyFileError
	CopyFileSrcErr       CopyFileSrcError
	CopyFileDstErr       CopyFileDstError
	WriteFileErr         WriteFileError
	ReadFileErr          ReadFileError

	DirInitErr          DirInitError
	DirNotExistsErr     DirNotExistsError
	DirFullPathEmptyErr DirFullPathEmptyError
	CreateDirErr        CreateDirError
	RenameDirErr        RenameDirError
	RemoveDirErr        RemoveDirError
	PermissionDirErr    PermissionDirError
	CopyDirErr          CopyDirError
	CopyDirSrcErr       CopyDirSrcError
	CopyDirDstErr       CopyDirDstError
	WriteDirErr         WriteDirError
	ReadDirErr          ReadDirError
)

func (*FileInitError) New(msg string) myError.IMyError {
	return &FileInitError{MyError: myError.MyError{Msg: array.NewDestruction("文件初始化错误", msg).JoinWithoutEmpty("：")}}
}

func (*FileInitError) Wrap(err error) myError.IMyError {
	return &FileInitError{MyError: myError.MyError{Msg: fmt.Errorf("文件初始化错误"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*FileInitError) Panic() myError.IMyError {
	return &FileInitError{MyError: myError.MyError{Msg: "文件初始化错误"}}
}

func (my *FileInitError) Error() string { return my.Msg }

func (*FileInitError) Is(target error) bool { return reflect.DeepEqual(target, &FileNotExistsErr) }

func (*FileNotExistsError) New(msg string) myError.IMyError {
	return &FileNotExistsError{MyError: myError.MyError{Msg: array.NewDestruction("文件不存在", msg).JoinWithoutEmpty("：")}}
}

func (*FileNotExistsError) Wrap(err error) myError.IMyError {
	return &FileNotExistsError{MyError: myError.MyError{Msg: fmt.Errorf("文件不存在"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*FileNotExistsError) Panic() myError.IMyError {
	return &FileNotExistsError{MyError: myError.MyError{Msg: "文件不存在"}}
}

func (my *FileNotExistsError) Error() string { return my.Msg }

func (*FileNotExistsError) Is(target error) bool { return reflect.DeepEqual(target, &FileNotExistsErr) }

func (*FileFullPathEmptyError) New(msg string) myError.IMyError {
	return &FileFullPathEmptyError{MyError: myError.MyError{Msg: "文件路径不能为空"}}
}

func (*FileFullPathEmptyError) Wrap(err error) myError.IMyError {
	return &FileFullPathEmptyError{MyError: myError.MyError{Msg: "文件路径不能为空"}}
}

func (*FileFullPathEmptyError) Panic() myError.IMyError {
	return &FileFullPathEmptyError{MyError: myError.MyError{Msg: "文件路径不能为空"}}
}

func (my *FileFullPathEmptyError) Error() string { return my.Msg }

func (*FileFullPathEmptyError) Is(target error) bool {
	return reflect.DeepEqual(target, &FileFullPathEmptyErr)
}

func (*CreateFileError) New(msg string) myError.IMyError {
	return &FileFullPathEmptyError{MyError: myError.MyError{Msg: array.NewDestruction("创建文件失败", msg).JoinWithoutEmpty("：")}}
}

func (*CreateFileError) Wrap(err error) myError.IMyError {
	return &CreateFileError{MyError: myError.MyError{Msg: fmt.Errorf("创建文件失败"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*CreateFileError) Panic() myError.IMyError {
	return &CreateFileError{MyError: myError.MyError{Msg: "创建文件失败"}}
}

func (my *CreateFileError) Error() string { return my.Msg }

func (*CreateFileError) Is(target error) bool {
	return reflect.DeepEqual(target, &CreateFileErr)
}

func (*RenameFileError) New(msg string) myError.IMyError {
	return &RenameFileError{MyError: myError.MyError{Msg: array.NewDestruction("修改文件名失败", msg).JoinWithoutEmpty("：")}}
}

func (*RenameFileError) Wrap(err error) myError.IMyError {
	return &RenameFileError{MyError: myError.MyError{Msg: fmt.Errorf("修改文件名失败"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*RenameFileError) Panic() myError.IMyError {
	return &RenameFileError{MyError: myError.MyError{Msg: "修改文件名失败"}}
}

func (my *RenameFileError) Error() string { return my.Msg }

func (*RenameFileError) Is(target error) bool {
	return reflect.DeepEqual(target, &RenameFileErr)
}

func (*RemoveFileError) New(msg string) myError.IMyError {
	return &RemoveFileError{MyError: myError.MyError{Msg: array.NewDestruction("删除文件失败", msg).JoinWithoutEmpty("：")}}
}

func (*RemoveFileError) Wrap(err error) myError.IMyError {
	return &RemoveFileError{MyError: myError.MyError{Msg: fmt.Errorf("删除文件失败"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*RemoveFileError) Panic() myError.IMyError {
	return &RemoveFileError{MyError: myError.MyError{Msg: "删除文件失败"}}
}

func (my *RemoveFileError) Error() string { return my.Msg }

func (*RemoveFileError) Is(target error) bool {
	return reflect.DeepEqual(target, &RemoveFileErr)
}

func (*PermissionFileError) New(msg string) myError.IMyError {
	return &PermissionFileError{MyError: myError.MyError{Msg: array.New([]string{"文件权限错误", msg}).JoinWithoutEmpty("：")}}
}

func (*PermissionFileError) Wrap(err error) myError.IMyError {
	return &PermissionFileError{MyError: myError.MyError{Msg: fmt.Errorf("文件权限错误"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*PermissionFileError) Panic() myError.IMyError {
	return &PermissionFileError{MyError: myError.MyError{Msg: "文件权限错误"}}
}

func (my *PermissionFileError) Error() string { return my.Msg }

func (*PermissionFileError) Is(target error) bool {
	return reflect.DeepEqual(target, &PermissionFileErr)
}

func (*CopyFileSrcError) New(msg string) myError.IMyError {
	return &CopyFileSrcError{MyError: myError.MyError{Msg: array.NewDestruction("复制文件时打开源文件错误", msg).JoinWithoutEmpty("：")}}
}

func (*CopyFileSrcError) Wrap(err error) myError.IMyError {
	return &CopyFileSrcError{MyError: myError.MyError{Msg: fmt.Errorf("复制文件时打开源文件错误"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*CopyFileSrcError) Panic() myError.IMyError {
	return &CopyFileSrcError{MyError: myError.MyError{Msg: "复制文件时打开源文件错误"}}
}

func (my *CopyFileSrcError) Error() string { return my.Msg }

func (*CopyFileSrcError) Is(target error) bool {
	return reflect.DeepEqual(target, &CopyFileSrcErr)
}

func (*CopyFileDstError) New(msg string) myError.IMyError {
	return &CopyFileDstError{MyError: myError.MyError{Msg: array.NewDestruction("复制文件时打开目标文件错误", msg).JoinWithoutEmpty("：")}}
}

func (*CopyFileDstError) Wrap(err error) myError.IMyError {
	return &CopyFileDstError{MyError: myError.MyError{Msg: fmt.Errorf("复制文件时打开目标文件错误"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*CopyFileDstError) Panic() myError.IMyError {
	return &CopyFileDstError{MyError: myError.MyError{Msg: "复制文件时打开目标文件错误"}}
}

func (my *CopyFileDstError) Error() string { return my.Msg }

func (*CopyFileDstError) Is(target error) bool {
	return reflect.DeepEqual(target, &CopyFileDstErr)
}

func (*CopyFileError) New(msg string) myError.IMyError {
	return &CopyFileError{MyError: myError.MyError{Msg: array.NewDestruction("复制文件错误", msg).JoinWithoutEmpty("：")}}
}

func (*CopyFileError) Wrap(err error) myError.IMyError {
	return &CopyFileError{MyError: myError.MyError{Msg: fmt.Errorf("复制文件错误"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*CopyFileError) Panic() myError.IMyError {
	return &CopyFileError{MyError: myError.MyError{Msg: "复制文件错误"}}
}

func (my *CopyFileError) Error() string { return my.Msg }

func (*CopyFileError) Is(target error) bool {
	return reflect.DeepEqual(target, &CopyFileErr)
}

func (*WriteFileError) New(msg string) myError.IMyError {
	return &WriteFileError{MyError: myError.MyError{Msg: array.NewDestruction("写入文件错误", msg).JoinWithoutEmpty("：")}}
}

func (*WriteFileError) Wrap(err error) myError.IMyError {
	return &WriteFileError{MyError: myError.MyError{Msg: fmt.Errorf("写入文件错误"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*WriteFileError) Panic() myError.IMyError {
	return &WriteFileError{MyError: myError.MyError{Msg: "写入文件错误"}}
}

func (my *WriteFileError) Error() string { return my.Msg }

func (*WriteFileError) Is(target error) bool {
	return reflect.DeepEqual(target, &WriteFileErr)
}

func (*ReadFileError) New(msg string) myError.IMyError {
	return &ReadFileError{MyError: myError.MyError{Msg: array.NewDestruction("读取文件错误", msg).JoinWithoutEmpty("：")}}
}

func (*ReadFileError) Wrap(err error) myError.IMyError {
	return &ReadFileError{MyError: myError.MyError{Msg: fmt.Errorf("读取文件错误"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*ReadFileError) Panic() myError.IMyError {
	return &ReadFileError{MyError: myError.MyError{Msg: "读取文件错误"}}
}

func (my *ReadFileError) Error() string { return my.Msg }

func (*ReadFileError) Is(target error) bool {
	return reflect.DeepEqual(target, &ReadFileErr)
}

func (*DirInitError) New(msg string) myError.IMyError {
	return &DirInitError{MyError: myError.MyError{Msg: array.NewDestruction("目录初始化错误", msg).JoinWithoutEmpty("：")}}
}

func (*DirInitError) Wrap(err error) myError.IMyError {
	return &DirInitError{MyError: myError.MyError{Msg: fmt.Errorf("目录初始化错误"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*DirInitError) Panic() myError.IMyError {
	return &DirInitError{MyError: myError.MyError{Msg: "目录初始化错误"}}
}

func (my *DirInitError) Error() string { return my.Msg }

func (*DirInitError) Is(target error) bool { return reflect.DeepEqual(target, &DirNotExistsErr) }

func (*DirNotExistsError) New(msg string) myError.IMyError {
	return &DirNotExistsError{MyError: myError.MyError{Msg: array.NewDestruction("目录不存在", msg).JoinWithoutEmpty("：")}}
}

func (*DirNotExistsError) Wrap(err error) myError.IMyError {
	return &DirNotExistsError{MyError: myError.MyError{Msg: fmt.Errorf("目录不存在"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*DirNotExistsError) Panic() myError.IMyError {
	return &DirNotExistsError{MyError: myError.MyError{Msg: "目录不存在"}}
}

func (my *DirNotExistsError) Error() string { return my.Msg }

func (*DirNotExistsError) Is(target error) bool { return reflect.DeepEqual(target, &DirNotExistsErr) }

func (*DirFullPathEmptyError) New(msg string) myError.IMyError {
	return &DirFullPathEmptyError{MyError: myError.MyError{Msg: "目录路径不能为空"}}
}

func (*DirFullPathEmptyError) Wrap(err error) myError.IMyError {
	return &DirFullPathEmptyError{MyError: myError.MyError{Msg: "目录路径不能为空"}}
}

func (*DirFullPathEmptyError) Panic() myError.IMyError {
	return &DirFullPathEmptyError{MyError: myError.MyError{Msg: "目录路径不能为空"}}
}

func (my *DirFullPathEmptyError) Error() string { return my.Msg }

func (*DirFullPathEmptyError) Is(target error) bool {
	return reflect.DeepEqual(target, &DirFullPathEmptyErr)
}

func (*CreateDirError) New(msg string) myError.IMyError {
	return &DirFullPathEmptyError{MyError: myError.MyError{Msg: array.NewDestruction("创建目录失败", msg).JoinWithoutEmpty("：")}}
}

func (*CreateDirError) Wrap(err error) myError.IMyError {
	return &CreateDirError{MyError: myError.MyError{Msg: fmt.Errorf("创建目录失败"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*CreateDirError) Panic() myError.IMyError {
	return &CreateDirError{MyError: myError.MyError{Msg: "创建目录失败"}}
}

func (my *CreateDirError) Error() string { return my.Msg }

func (*CreateDirError) Is(target error) bool { return reflect.DeepEqual(target, &CreateDirErr) }

func (*RenameDirError) New(msg string) myError.IMyError {
	return &RenameDirError{MyError: myError.MyError{Msg: array.NewDestruction("修改目录名失败", msg).JoinWithoutEmpty("：")}}
}

func (*RenameDirError) Wrap(err error) myError.IMyError {
	return &RenameDirError{MyError: myError.MyError{Msg: fmt.Errorf("修改目录名失败"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*RenameDirError) Panic() myError.IMyError {
	return &RenameDirError{MyError: myError.MyError{Msg: "修改目录名失败"}}
}

func (my *RenameDirError) Error() string { return my.Msg }

func (*RenameDirError) Is(target error) bool { return reflect.DeepEqual(target, &RenameDirErr) }

func (*RemoveDirError) New(msg string) myError.IMyError {
	return &RemoveDirError{MyError: myError.MyError{Msg: array.NewDestruction("删除目录失败", msg).JoinWithoutEmpty("：")}}
}

func (*RemoveDirError) Wrap(err error) myError.IMyError {
	return &RemoveDirError{MyError: myError.MyError{Msg: fmt.Errorf("删除目录失败"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*RemoveDirError) Panic() myError.IMyError {
	return &RemoveDirError{MyError: myError.MyError{Msg: "删除目录失败"}}
}

func (my *RemoveDirError) Error() string { return my.Msg }

func (*RemoveDirError) Is(target error) bool { return reflect.DeepEqual(target, &RemoveDirErr) }

func (*PermissionDirError) New(msg string) myError.IMyError {
	return &PermissionDirError{MyError: myError.MyError{Msg: array.NewDestruction("目录权限错误", msg).JoinWithoutEmpty("：")}}
}

func (*PermissionDirError) Wrap(err error) myError.IMyError {
	return &PermissionDirError{MyError: myError.MyError{Msg: fmt.Errorf("目录权限错误"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*PermissionDirError) Panic() myError.IMyError {
	return &PermissionDirError{MyError: myError.MyError{Msg: "目录权限错误"}}
}

func (my *PermissionDirError) Error() string { return my.Msg }

func (*PermissionDirError) Is(target error) bool { return reflect.DeepEqual(target, &PermissionDirErr) }

func (*CopyDirSrcError) New(msg string) myError.IMyError {
	return &CopyDirSrcError{MyError: myError.MyError{Msg: array.NewDestruction("复制目录时打开源目录错误", msg).JoinWithoutEmpty("：")}}
}

func (*CopyDirSrcError) Wrap(err error) myError.IMyError {
	return &CopyDirSrcError{MyError: myError.MyError{Msg: fmt.Errorf("复制目录时打开源目录错误"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*CopyDirSrcError) Panic() myError.IMyError {
	return &CopyDirSrcError{MyError: myError.MyError{Msg: "复制目录时打开源目录错误"}}
}

func (my *CopyDirSrcError) Error() string { return my.Msg }

func (*CopyDirSrcError) Is(target error) bool { return reflect.DeepEqual(target, &CopyDirSrcErr) }

func (*CopyDirDstError) New(msg string) myError.IMyError {
	return &CopyDirDstError{MyError: myError.MyError{Msg: array.NewDestruction("复制目录时打开目标目录错误", msg).JoinWithoutEmpty("：")}}
}

func (*CopyDirDstError) Wrap(err error) myError.IMyError {
	return &CopyDirDstError{MyError: myError.MyError{Msg: fmt.Errorf("复制目录时打开目标目录错误"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*CopyDirDstError) Panic() myError.IMyError {
	return &CopyDirDstError{MyError: myError.MyError{Msg: "复制目录时打开目标目录错误"}}
}

func (my *CopyDirDstError) Error() string { return my.Msg }

func (*CopyDirDstError) Is(target error) bool { return reflect.DeepEqual(target, &CopyDirDstErr) }

func (*CopyDirError) New(msg string) myError.IMyError {
	return &CopyDirError{MyError: myError.MyError{Msg: array.NewDestruction("复制目录错误", msg).JoinWithoutEmpty("：")}}
}

func (*CopyDirError) Wrap(err error) myError.IMyError {
	return &CopyDirError{MyError: myError.MyError{Msg: fmt.Errorf("复制目录错误"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*CopyDirError) Panic() myError.IMyError {
	return &CopyDirError{MyError: myError.MyError{Msg: "复制目录错误"}}
}

func (my *CopyDirError) Error() string { return my.Msg }

func (*CopyDirError) Is(target error) bool { return reflect.DeepEqual(target, &CopyDirErr) }

func (*WriteDirError) New(msg string) myError.IMyError {
	return &WriteDirError{MyError: myError.MyError{Msg: array.NewDestruction("写入目录错误", msg).JoinWithoutEmpty("：")}}
}

func (*WriteDirError) Wrap(err error) myError.IMyError {
	return &WriteDirError{MyError: myError.MyError{Msg: fmt.Errorf("写入目录错误"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*WriteDirError) Panic() myError.IMyError {
	return &WriteDirError{MyError: myError.MyError{Msg: "写入目录错误"}}
}

func (my *WriteDirError) Error() string { return my.Msg }

func (*WriteDirError) Is(target error) bool { return reflect.DeepEqual(target, &WriteDirErr) }

func (*ReadDirError) New(msg string) myError.IMyError {
	return &ReadDirError{MyError: myError.MyError{Msg: array.NewDestruction("读取目录错误", msg).JoinWithoutEmpty("：")}}
}

func (*ReadDirError) Wrap(err error) myError.IMyError {
	return &ReadDirError{MyError: myError.MyError{Msg: fmt.Errorf("读取目录错误"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*ReadDirError) Panic() myError.IMyError {
	return &ReadDirError{MyError: myError.MyError{Msg: "读取目录错误"}}
}

func (my *ReadDirError) Error() string { return my.Msg }

func (*ReadDirError) Is(target error) bool { return reflect.DeepEqual(target, &ReadDirErr) }
