# simpleDB driver

`driver` 是 `simpleDB/kernal` 的上层封装，提供稳定、统一、易用的 API。

## 设计目标

- 定义统一驱动 API
- 映射核心层事务能力（MVCC）
- 统一错误模型（`DriverError` + `ErrorCode`）
- 支持事务隔离级别

## 快速开始

```go
package main

import (
    "fmt"

    "github.com/aid297/aid/simpleDB/driver"
)

func main() {
    db, err := driver.New.DB("demo", "users")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    err = db.Configure(driver.TableSchema{Columns: []driver.Column{
        {Name: "id", Type: "uuid:v7", PrimaryKey: true, AutoIncrement: true},
        {Name: "username", Type: "string", Unique: true, Required: true},
    }})
    if err != nil {
        panic(err)
    }

    tx, err := db.BeginTxWithOptions(driver.TxOptions{Isolation: driver.TxIsolationSnapshot})
    if err != nil {
        panic(err)
    }

    _, err = tx.InsertRow(driver.Row{"username": "alice"})
    if err != nil {
        _ = tx.Rollback()
        panic(err)
    }

    if err = tx.Commit(); err != nil {
        panic(err)
    }

    row, ok, err := db.FindOne(driver.QueryCondition{Field: "username", Operator: driver.QueryOpEQ, Value: "alice"})
    if err != nil {
        panic(err)
    }
    fmt.Println(ok, row)
}
```

## 错误模型

所有返回错误都会包装为 `DriverError`：

- `invalid_argument`
- `not_found`
- `conflict`
- `closed`
- `read_only`
- `internal`

可用 `errors.As(err, *DriverError)` 获取错误码。
