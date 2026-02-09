### 流处理

```go
package main

import (
	`context`
	`fmt`

	`github.com/aid297/aid/steam`
)

func main() {
	p := context.TODO()
	ctx := context.Background()
	steam.APP.Steam.New(
		steam.ReadCloser(nil), // HTTP 响应体或其他流
		steam.CopyFn(func(copied []byte) error {
			// 将 copied 已拷贝体复制到其他地方存储
			ctx = context.WithValue(p, "copy", copied)
			return nil // 当程序时，readCloser 会被重新填充
		}),
	)

	fmt.Printf(ctx.Value("copy").(string))
}
```

