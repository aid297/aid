package filesystemV2

import (
	"os"
	"path"
	"testing"
)

func Test1(t *testing.T) {
	t.Run("test1 create file", func(t *testing.T) {
		file := FileApp.NewByRel("./test.txt")
		if file.Error() != nil {
			t.Fatalf("init file error: %v", file.Error())
		}
		t.Logf("file exist: %v", file.GetExist())

		if file.CreateDefaultMode().Error() != nil {
			t.Fatalf("create file error: %v", file.Error())
		}

		t.Logf("file info: %v", file.GetInfo())
	})
}

func Test2(t *testing.T) {
	t.Run("test2 rename", func(t *testing.T) {
		file := FileApp.NewByRel("./test.txt")
		if file.Error() != nil {
			t.Fatalf("init file error: %v", file.Error())
		}
		t.Logf("file exist: %v", file.GetExist())

		basePath := file.GetBasePath()
		t.Logf("file base path: %s", basePath)

		file.Rename("test2.txt")
		if file.GetFullPath() != path.Join(basePath, file.GetName()) {
			t.Fatal("rename fail: wrong name")
		}
	})
}

func Test3(t *testing.T) {
	t.Run("test3 remove", func(t *testing.T) {
		file := FileApp.NewByRel("./test2.txt")
		if file.Error() != nil {
			t.Fatalf("init file error: %v", file.Error())
		}
		t.Logf("file exist: %v", file.GetExist())

		file.Remove()
		t.Logf("file info: %v", file.GetInfo())
	})
}

func Test4(t *testing.T) {
	t.Run("test4 check permission", func(t *testing.T) {
		file := FileApp.NewByRel("./test.txt")
		if file.Error() != nil {
			t.Fatalf("init file error: %v", file.Error())
		}
		t.Logf("file exist: %v", file.GetExist())

		file.Remove()
		t.Logf("file info: %v", file.GetInfo())
	})
}
func Test5(t *testing.T) {
	t.Run("test5 write", func(t *testing.T) {
		file := FileApp.NewByRel("./test.txt")
		// file.CreateDefaultMode()
		if file.Error() != nil {
			t.Fatalf("init file error: %v", file.Error())
		}
		t.Logf("file exist: %v -> %v", file.GetExist(), file.GetFullPath())

		if file.Write([]byte("this is test content..."), os.ModePerm, DefaultCreateMode|os.O_APPEND).Error() != nil {
			t.Fatalf("write file: %v", file.Error())
		}
	})
}

func Test6(t *testing.T) {
	t.Run("test6 copy file", func(t *testing.T) {
		file := FileApp.NewByRel("./test.txt")
		if file.Error() != nil {
			t.Fatalf("init file error: %v", file.Error())
		}
		t.Logf("file exist: %v -> %v", file.GetExist(), file.GetFullPath())

		if file.CopyTo("./test2.txt").Error() != nil {
			t.Fatalf("copy file error: %v", file.Error())
		}
	})
}

func Test7(t *testing.T) {
	t.Run("test7 read file", func(t *testing.T) {
		file := FileApp.NewByRel("./test.txt")
		if file.Error() != nil {
			t.Fatalf("init file error: %v", file.Error())
		}
		t.Logf("file exist: %v -> %v", file.GetExist(), file.GetFullPath())

		if string(file.Read()) != `this is test content...` {
			t.Fatal("file content wrong")
		}
		if file.Error() != nil {
			t.Fatalf("read file wrong: %v", file.Error())
		}
	})
}
