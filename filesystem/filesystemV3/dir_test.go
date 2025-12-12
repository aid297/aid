package filesystemV3

import (
	"testing"
)

func TestDir1(t *testing.T) {
	t.Run("查看当前目录", func(t *testing.T) {
		t.Logf("当前目录：%s", NewDirRel().FullPath)
	})
}

func TestDir2(t *testing.T) {
	t.Run("目录追加", func(t *testing.T) {
		dir := NewDirRel()
		t.Logf("追加目录：%s\n", dir.Join("test-a", "test-b", "test-c").FullPath)
		t.Logf("查看目录是否存在：%v", dir.Exist)
	})
}

func TestDir3(t *testing.T) {
	t.Run("创建多级目录", func(t *testing.T) {
		dir := NewDirRel()
		t.Logf("追加目录：%s\n", dir.Join("test-a", "test-b", "test-c").FullPath)
		if err := dir.Create().Error; err != nil {
			t.Fatalf("创建目录失败：%s", err)
		}
		t.Logf("查看目录是否存在：%v", dir.Exist)
	})
}

func TestDir4(t *testing.T) {
	t.Run("删除单个目录", func(t *testing.T) {
		dir := NewDirRel(APP.DirAttr.Path.Set("test-a", "test-b", "test-c"))
		if err := dir.Remove().Error; err != nil {
			t.Fatalf("删除目录失败：%s", err)
		}
	})
}

func TestDir5(t *testing.T) {
	t.Run("删除多级目录", func(t *testing.T) {
		dir := NewDirRel(APP.DirAttr.Path.Set("test-a"))
		if err := dir.RemoveAll().Error; err != nil {
			t.Fatalf("删除失败：%s", err)
		}
	})
}

func TestDir6(t *testing.T) {
	t.Run("列出当前目录下的所有文件和子目录", func(t *testing.T) {
		dir := NewDirRel(APP.DirAttr.Path.Set("test-a1", "test-a2", "test-a3"))

		if err := dir.Create().Error; err != nil {
			t.Fatalf("创建目录失败：%s", err)
		}

		dir.Up().Up().LS()

		t.Logf("%+v", dir.Dirs)

		for idx := range dir.Dirs {
			t.Logf("子目录1：%s\n", dir.Dirs[idx].Name)

			for idx2 := range dir.Dirs[idx].Dirs {
				t.Logf("子目录2：%s/%s\n", dir.Dirs[idx].Name, dir.Dirs[idx].Dirs[idx2].Name)
			}
		}

		dir.RemoveAll()
	})
}

func TestDir7(t *testing.T) {
	t.Run("复制目录", func(t *testing.T) {
		dir := NewDirRel(APP.DirAttr.Path.Set("test-a1", "test-a2", "test-a3"))

		if err := dir.Create().Error; err != nil {
			t.Fatalf("创建目录失败：%s", err)
		}

		dir.Up().Up()

		if err := dir.CopyAllTo(true, "test-b1").Error; err != nil {
			t.Fatalf("复制目录失败：%s", err)
		}
	})
}

func TestDir8(t *testing.T) {
	t.Run("文件夹改名", func(t *testing.T) {
		dir := NewDirRel(APP.DirAttr.Path.Set("test-a1", "test-a2", "test-a3")).Create(DirMode(0777))
		if dir.Error != nil {
			t.Fatalf("创建目录失败：%s", dir.Error)
		}

		if err := dir.Rename("test-aaa").Error; err != nil {
			t.Fatalf("重命名失败")
		}
	})
}
