package filesystemV2

import (
	"path/filepath"

	"github.com/aid297/aid/array"
)

func getRootPath(dir string) string {
	rootPath, _ := filepath.Abs(".")

	return filepath.Clean(filepath.Join(rootPath, dir))
}

// CopyFiles 批量复制文件
func CopyFiles(srcFiles, dstFiles *array.AnyArray[*File]) {
	srcFiles.Each(func(idx int, item *File) { item.CopyTo(dstFiles.Get(idx).GetFullPath()) })
}

// CopyFilesByDstPath 批量复制文件：通过dst绝对路径（无法指定拷贝后的文件名）
func CopyFilesByDstPath(srcFiles *array.AnyArray[*File], dstPath string) {
	dstFiles := array.Make[*File](srcFiles.Len())
	srcFiles.Each(func(idx int, item *File) { dstFiles.Set(idx, FileApp.NewByAbs(filepath.Join(dstPath, item.GetName()))) })

	CopyFiles(srcFiles, dstFiles)
}

// CopyFilesBy2Path 批量复制文件：通过src绝对路径到dst绝对路径（无法指定拷贝后的文件名）
func CopyFilesBy2Path(srcPath, dstPath string) {
	CopyFilesByDstPath(DirApp.NewByAbs(srcPath).GetFiles(), dstPath)
}
