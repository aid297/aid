### FilesystemV4 用法

1. 初始化`Dir` 和`File`
   ```go
   package main
   
   import (
   	. "fmt"
   
   	"github.com/aid297/aid/filesystem/filesystemV4"
   )
   
   func main() {
   	dir1 := filesystemV4.NewDir(filesystemV4.Rel("."))
   	Printf("init dir by rel: %s\n", dir1.GetFullPath()) // init dir by rel: /Users/yujizhou/development/projects/go/readme
   
   	dir2 := filesystemV4.NewDir(filesystemV4.Abs("/Users/yujizhou/development/projects/go/readme"))
   	Printf("init dir by abs: %s\n", dir2.GetFullPath()) // init dir by abs: /Users/yujizhou/development/projects/go/readme
   
   	file1 := filesystemV4.NewFile(filesystemV4.Rel("./main.go"))
   	Printf("init file by rel: %s\n", file1.GetFullPath()) // init file by rel: /Users/yujizhou/development/projects/go/readme/main.go
   
   	file2 := filesystemV4.NewFile(filesystemV4.Abs("/Users/yujizhou/development/projects/go/readme/main.go"))
   	Printf("init file by abs: %s\n", file2.GetFullPath()) // init file by abs: /Users/yujizhou/development/projects/go/readme/main.go
   }
   ```

2. 创建`文件夹`
   ```go
   package main
   
   import (
   	. "fmt"
   
   	"github.com/aid297/aid/filesystem/filesystemV4"
   )
   
   func main() {
   	dir1 := filesystemV4.NewDir(filesystemV4.Rel("./a"))
   	Printf("./a is exist before create: %v\n", dir1.GetExist()) // ./a is exist before create: true
   
   	err := dir1.Create(filesystemV4.Flag(0644)).GetError()
   	if err != nil {
   		panic(err)
   	}
   
   	Printf("./a is exist after create: %v\n", dir1.GetExist()) // ./a is exist after create: true
   	Printf("./a fullpath is: %s \n", dir1.GetFullPath())       // ./a fullpath is: /Users/yujizhou/development/projects/go/readme/a
     // 注意：创建文件夹之后，当前对象的目录会自动移动到新文件夹下
   }
   ```

3. 创建`文件`，并写入内容
   ```go
   package main
   
   import (
   	. "fmt"
   
   	"github.com/aid297/aid/filesystem/filesystemV4"
   )
   
   func main() {
   	dir1 := filesystemV4.NewDir(filesystemV4.Rel("./a"))
   	file1 := filesystemV4.NewFile(filesystemV4.Abs(dir1.GetFullPath(), "1.txt"))
   	Printf("file1 is exist: %v\n", file1.GetExist()) // file1 is exist: false
   
   	err := file1.Create().GetError()
   	if err != nil {
   		panic(err)
   	}
   
   	Printf("file1 is exist: %v\n", file1.GetExist()) // file1 is exist: true
   	// 注意，如果a文件夹不存在，file1.Create()会直接创建这个文件夹
   
   	err = file1.Write([]byte("Is somthing here ...")).GetError()
   	if err != nil {
   		panic(err)
   	}
   
   	fileContent, err := file1.Read()
   	if err != nil {
   		panic(err)
   	}
   
   	Printf("file1 content: %s\n", string(fileContent)) // file1 content: Is somthing here ...
   }
   ```

4. 复制`文件`
   ```go
   package main
   
   import (
   	. "fmt"
   
   	"github.com/aid297/aid/filesystem/filesystemV4"
   )
   
   func main() {
   	dir1 := filesystemV4.NewDir(filesystemV4.Rel("./a"))
   	file1 := filesystemV4.NewFile(filesystemV4.Abs(dir1.GetFullPath(), "1.txt"))
   	Printf("file1 is exist: %v\n", file1.GetExist()) // file1 is exist: false
   
   	file1.CopyTo(true, "./2.txt") // copy file1 to 2.txt, but file1 is not exist, so 2.txt is not exist
   	file2 := filesystemV4.NewFile(filesystemV4.Rel("./2.txt"))
   	Printf("file1 fullpath: %s\n", file1.GetFullPath()) // 复制不会导致原文件路径发生变化
   	Printf("file2 is exist: %v\n", file2.GetExist())    // file2 is exist: true
   	Printf("file2 fullpath: %s\n", file2.GetFullPath()) // file2 fullpath: /Users/yujizhou/development/projects/go/readme/2.txt
   	// 注意：复制文件时的isRel指的是当前路径的相对路径，不是 file1 的相对路径
   }
   ```

5. 删除`文件`
   ```go
   package main
   
   import (
   	. "fmt"
   
   	"github.com/aid297/aid/filesystem/filesystemV4"
   )
   
   func main() {
   	dir1 := filesystemV4.NewDir(filesystemV4.Rel("./a"))
   	file1 := filesystemV4.NewFile(filesystemV4.Abs(dir1.GetFullPath(), "1.txt"))
   	file1.CopyTo(true, "./a/2.txt")
   
   	file2 := filesystemV4.NewFile(filesystemV4.Abs(dir1.GetFullPath(), "2.txt"))
   	Printf("file2 is exist: %v\n", file2.GetExist()) // file2 is exist: true
   
   	err := file2.Remove().GetError()
   	if err != nil {
   		panic(err)
   	}
   
   	Printf("file2 is exist: %v\n", file2.GetExist()) // file2 is exist: false
   }
   ```

6. 