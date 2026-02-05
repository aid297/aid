package filesystemV4

import (
	"testing"
)

func Test1(t *testing.T) {
	file := NewFile(Rel("./1.txt"))
	t.Logf("路径：%s", file.GetFullPath())
	t.Logf("错误：%v", file.Write([]byte("abc")).GetError())
}
