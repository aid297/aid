### Compression

1. Zlib 压缩与解压缩
   ```go
   package main
   
   import (
   	. "fmt"
   
   	"github.com/aid297/aid/compression"
   )
   
   func main() {
   	unzip := []byte("abc")
   
   	zipper := compression.NewZlib()
   
   	zipped, err := zipper.Compress(unzip)
   	if err != nil {
   		panic(err)
   	}
   
   	Printf("Zipped: %v\n", zipped) // Zipped: [120 156 74 76 74 6 4 0 0 255 255 2 77 1 39]
   
   	unzipped, err := zipper.Decompress(zipped)
   	if err != nil {
   		panic(err)
   	}
   
   	Println("Unzipped:", string(unzipped)) // Unzipped: abc
   }
   ```