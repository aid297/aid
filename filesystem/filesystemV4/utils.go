package filesystemV4

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aid297/aid/operation/operationV2"
)

func getRootPath(dir string) string {
	rootPath, _ := filepath.Abs(".")

	return filepath.Clean(filepath.Join(rootPath, dir))
}

func copyFileTo(src, dst string) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("无法打开源文件 %s: %w", src, err)
	}
	defer func() { _ = srcFile.Close() }()

	// 获取源文件信息（用于保留权限）
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("无法获取源文件信息 %s: %w", src, err)
	}

	// 创建目标文件，保留原始权限
	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("无法创建目标文件 %s: %w", dst, err)
	}
	defer func() { _ = dstFile.Close() }()

	// 复制文件内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("复制文件内容失败 %s -> %s: %w", src, dst, err)
	}

	return nil
}

func copyDirTo(src, dst string) error {
	// 获取源目录信息
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("无法获取源目录信息: %w", err)
	}

	// 确认源路径是目录
	if !srcInfo.IsDir() {
		return fmt.Errorf("源路径 %s 不是一个目录", src)
	}

	// 创建目标目录（保留原始权限）
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("无法创建目标目录 %s: %w", dst, err)
	}

	// 读取源目录内容
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("无法读取源目录 %s: %w", src, err)
	}

	// 遍历并处理每个条目
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// 递归处理子目录
			if err := copyDirTo(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// 复制文件
			if err := copyFileTo(srcPath, dstPath); err != nil {
				return fmt.Errorf("复制文件 %s 失败: %w", srcPath, err)
			}
		}
	}

	return nil
}

func getDirByCopy(isRel bool, dstPaths ...string) Dir {
	return operationV2.NewTernary(operationV2.TrueFn(func() Dir { return APP.Dir.Rel(dstPaths...) }), operationV2.FalseFn(func() Dir { return APP.Dir.Abs(dstPaths...) })).GetByValue(isRel)
}
