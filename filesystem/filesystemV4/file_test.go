package filesystemV4

import "testing"

func TestFile1(t *testing.T) {
	t.Run("新建文件", func(t *testing.T) {
		file := APP.File.Rel("test-file.txt")
		if err := file.Write([]byte("aaa")).Error; err != nil {
			t.Fatalf("创建文件失败：%s", err)
		}
	})
}

func TestFile2(t *testing.T) {
	t.Run("新建文件", func(t *testing.T) {
		file := APP.File.Rel("test-file.txt")
		if err := file.Write([]byte("aaa")).Error; err != nil {
			t.Fatalf("创建文件失败：%s", err)
		}

		t.Logf("查看文件是否存在：%v\n", file.Exist)
	})
}

func TestFile3(t *testing.T) {
	t.Run("拷贝文件", func(t *testing.T) {
		file := APP.File.Rel("test-file.txt")
		if err := file.Write([]byte("aaa")).Error; err != nil {
			t.Fatalf("创建文件失败：%s", err)
		}

		if err := file.CopyTo(true, "test-d", "test-file2.txt").Error; err != nil {
			t.Fatalf("复制文件失败：%s", err)
		}
	})
}

func TestFile4(t *testing.T) {
	t.Run("删除文件", func(t *testing.T) {
		file := APP.File.Rel("test-file.txt")
		if err := file.Write([]byte("aaa")).Error; err != nil {
			t.Fatalf("创建文件失败：%s", err)
		}

		if err := file.Remove().Error; err != nil {
			t.Fatalf("删除文件失败：%s", err)
		}

		t.Logf("文件是否存在：%v", file.Exist)
	})
}

func TestFile5(t *testing.T) {
	t.Run("读取文件内容", func(t *testing.T) {
		var (
			err     error
			content []byte
			file    File
		)

		file = APP.File.Rel("test-file.txt")
		if err = file.Write([]byte("aaa")).Error; err != nil {
			t.Fatalf("创建文件失败：%s", err)
		}

		if content, err = file.Read(); err != nil {
			t.Fatalf("读取文件内容失败：%s", err)
		}
		t.Logf("文件内容：%s", string(content))
	})
}

func TestFile6(t *testing.T) {
	t.Run("文件改名", func(t *testing.T) {
		file := APP.File.Rel("test-file.txt")
		if err := file.Write([]byte("aaa")).Error; err != nil {
			t.Fatalf("创建文件失败：%s", err)
		}

		if err := file.Rename("test-file-rename.txt").Error; err != nil {
			t.Fatalf("文件改名失败：%s", err)
		}

		t.Logf("文件新路径：%s", file.FullPath)
	})
}
