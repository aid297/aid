package filesystemV4

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	dir := NewDir(Rel("./a"))
	dir.Create()
	for i := range 5 {
		file := NewFile(Rel("./a", fmt.Sprintf("file-%d.txt", i+1)))
		file.Create()
		file.Write([]byte(fmt.Sprintf("file %d", i+1)))
	}
}

func Test2(t *testing.T) {
	dir := NewDir(Rel("./a"))
	err := dir.CopyTo(true, "./b").GetError()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	t.Logf("passed")
}

func Test3(t *testing.T) {
	file := NewFile(Rel("./a/file-1.txt"))
	file.CopyTo(true, "./b/file-1.txt")
}
