package filesystem

import (
	"path/filepath"
	"runtime"
	"strings"
)

type (
	PathAttributer interface {
		Joins(paths ...string) PathAttributer
		Register(f Filesystem)
		GetPath() string
	}

	AttrPath struct{ path string }
)

func Rel(paths ...string) PathAttributer { return AttrPath{path: getRootPath(filepath.Join(paths...))} }
func Abs(paths ...string) PathAttributer { return AttrPath{path: filepath.Join(paths...)} }
func Auto(paths ...string) PathAttributer {
	if len(paths) > 0 {
		if isAbs(paths[0]) {
			return Abs(paths...)
		} else {
			return Rel(paths...)
		}
	}

	return Rel(".")
}

func (my AttrPath) Register(f Filesystem) { f.SetFullPathForAttr(my.path) }
func (my AttrPath) Joins(paths ...string) PathAttributer {
	my.path = filepath.Join(append([]string{my.path}, paths...)...)
	return my
}
func (my AttrPath) GetPath() string { return my.path }

// isAbs 自动判断路径是绝对路径还是相对路径
// 支持 Windows、Linux、Unix 等所有平台
// 返回 true 表示绝对路径，false 表示相对路径
func isAbs(path string) bool {
	if path == "" {
		return false
	}

	// 使用 Go 标准库的 filepath.IsAbs，它会自动根据运行平台判断
	// Windows: 检查是否以盘符开头（如 C:\）或 UNC 路径（如 \\server\share）
	// Unix/Linux: 检查是否以 / 开头
	if filepath.IsAbs(path) {
		return true
	}

	// 额外的兼容性检查：处理一些特殊情况
	// 1. Windows UNC 路径（\\server\share）
	if runtime.GOOS == "windows" {
		// 检查 UNC 路径格式
		if strings.HasPrefix(path, `\\`) || strings.HasPrefix(path, `//`) {
			return true
		}
		// 检查带盘符的路径（如 C:/ 或 C:\）
		if len(path) >= 2 && path[1] == ':' {
			return true
		}
	}

	// 2. Unix 风格的绝对路径（已经在 filepath.IsAbs 中处理）
	// 3. 其他情况都是相对路径
	return false
}
