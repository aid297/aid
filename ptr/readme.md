### 便捷指针

```go
package main

import (
	. `fmt`

	"github.com/aid297/aid/ptr"
)

func main() {
	a := ptr.New(123)
	Printf("%#v\n", a)

	b := ptr.New("张三")
	Printf("%#v\n", b)

	c := ptr.New('a')
	Printf("%#v\n", c)
}
```

