package filesystemV2

import (
	"testing"

	"github.com/aid297/aid/debugLogger"
)

func Test8(t *testing.T) {
	t.Run("test8 ls dirs", func(t *testing.T) {
		dir := DirApp.NewByRel("../../")
		t.Log(dir.getFullPath())
		dir.Ls()
		for _, dir := range dir.GetDirs().ToSlice() {
			debugLogger.Print("%v\n", dir.getFullPath())
		}
	})
}

func Test9(t *testing.T) {
	t.Run("test9 ls files", func(t *testing.T) {
		dir := DirApp.NewByRel("../..")
		for _, file := range dir.GetFiles().ToSlice() {
			debugLogger.Print("%v -> %v\n", file.GetFullPath(), file.getExtension())
		}
	})
}

func Test10(t *testing.T) {
	t.Run("test10 copy all files", func(t *testing.T) {
		dir := DirApp.NewByRel("../..")
		dst := DirApp.NewByRel("./copyAllFilesTest")
		dir.CopyAllFilesTo(dst.GetFullPath())
		dst.GetFiles().Each(func(idx int, item *File) { debugLogger.Print("%v\n", item.GetName()) })
	})
}

func Test11(t *testing.T) {
	t.Run("test11 copy all dirs", func(t *testing.T) {
		src := DirApp.NewByRel("../../a")
		dst := DirApp.NewByRel("./copyAllDirsTest")
		src.CopyAllDirsTo(dst.GetFullPath())
		dst.GetDirs().Each(func(idx int, item *Dir) { debugLogger.Print("%v\n", item.GetName()) })
	})
}
