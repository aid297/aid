### digestV2 使用方法

1. Bcrypt
   ```go
   package main
   
   import (
   	. "fmt"
   
   	"github.com/aid297/aid/digest/digestV2"
   )
   
   func main() {
   	b := digestV2.NewBcrypt("123123")
   	hashed := b.Hash()
   	Printf("hashed: %s\n", hashed)                 // hashed: $2a$10$KnHEwMWRx5MQqx5NSTHJQeP82ehaqXedlZu35UM8.ATqxLAgilCgC
   	Printf("check: %v\n", b.Check(string(hashed))) // check: false
   }
   ```
   
2. MD5
   ```go
   package main
   
   import (
   	. "fmt"
   
   	"github.com/aid297/aid/digest/digestV2"
   )
   
   func main() {
   	hashed, e := digestV2.NewMD5("123123").Encode()
   	if e != nil {
   		panic(e)
   	}
   	Printf("md5: %v\n", hashed) // md5: 4297f44b13955235245b2497399d7a93
   }
   ```
   
3. SHA
   ```go
   package main
   
   import (
   	. "fmt"
   
   	"github.com/aid297/aid/digest/digestV2"
   )
   
   func main() {
   	sha := digestV2.NewSHA("123123")
   	hashed, err := sha.Encode256()
   	if err != nil {
   		panic(err)
   	}
   	Printf("hashed: %s\n", hashed) // hashed: 96cae35ce8a9b0244178bf28e4966c2ce1b8385723a96a6b838858cdd6ca0a1e
   
   	sum256 := sha.Encode256Sum256()
   	Printf("sum256: %s\n", sum256) // sum256: 96cae35ce8a9b0244178bf28e4966c2ce1b8385723a96a6b838858cdd6ca0a1e
   }
   ```
   
4. SM3
   ```go
   package main
   
   import (
   	. "fmt"
   
   	"github.com/aid297/aid/digest/digestV2"
   )
   
   func main() {
   	sm3 := digestV2.NewSM3("123123")
   	hased, err := sm3.Encode()
   	if err != nil {
   		panic(err)
   	}
   	Printf("hased: %s\n", hased) // hased: c68ac63173fcfc537bf22f19a425977029d7dd35ddc5d76b36e58af222dfda39
   }
   ```