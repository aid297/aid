package filesystemV4

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func getRootPath(dir string) string {
	rootPath, _ := filepath.Abs(".")

	return filepath.Clean(filepath.Join(rootPath, dir))
}

func isDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err // 路径不存在或权限不足等错误
	}
	return info.IsDir(), nil
}

func copyFileTo(src, dst string) (err error) {
	var (
		srcFile *os.File
		dstFile *os.File
		srcInfo os.FileInfo
	)

	// 打开源文件
	if srcFile, err = os.Open(src); err != nil {
		return fmt.Errorf("无法打开源文件 %s: %w", src, err)
	}
	defer srcFile.Close()

	// 获取源文件信息（用于保留权限）
	if srcInfo, err = os.Stat(src); err != nil {
		return fmt.Errorf("无法获取源文件信息 %s: %w", src, err)
	}

	// 创建目标文件，保留原始权限
	if dstFile, err = os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode()); err != nil {
		return fmt.Errorf("无法创建目标文件 %s: %w", dst, err)
	}
	defer dstFile.Close()

	// 复制文件内容
	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("复制文件内容失败 %s -> %s: %w", src, dst, err)
	}

	return
}

func copyDirTo(src, dst string) (err error) {
	var (
		srcInfo          os.FileInfo
		entries          []os.DirEntry
		entry            os.DirEntry
		srcPath, dstPath string
	)

	// 获取源目录信息
	if srcInfo, err = os.Stat(src); err != nil {
		return fmt.Errorf("无法获取源目录信息: %w", err)
	}

	// 确认源路径是目录
	if !srcInfo.IsDir() {
		return fmt.Errorf("源路径 %s 不是一个目录", src)
	}

	// 创建目标目录（保留原始权限）
	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("无法创建目标目录 %s: %w", dst, err)
	}

	// 读取源目录内容
	if entries, err = os.ReadDir(src); err != nil {
		return fmt.Errorf("无法读取源目录 %s: %w", src, err)
	}

	// 遍历并处理每个条目
	for _, entry = range entries {
		srcPath = filepath.Join(src, entry.Name())
		dstPath = filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// 递归处理子目录
			if err = copyDirTo(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// 复制文件
			if err = copyFileTo(srcPath, dstPath); err != nil {
				return fmt.Errorf("复制文件 %s 失败: %w", srcPath, err)
			}
		}
	}

	return nil
}
