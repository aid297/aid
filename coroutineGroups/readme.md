
### CoroutineGroups 使用说明

泛型协程组，支持并发限制、批量分轮执行、独立超时重试。

---

#### 核心类型

| 类型 | 说明 |
|---|---|
| `CoroutineGrouper[T]` | 协程组接口 |
| `Func[T]` | 任务函数签名 `func() Result[T]` |
| `Result[T]` | 结果接口 |
| `ResultImpl[T]` | 结果实现（Error / Skip / Value） |
| `RetryConfig` | 重试配置（内部类型） |
| `CoroutineGroupRetryAttr` | 重试选项函数 `func(*RetryConfig)` |
| `BackoffStrategy` | 退避策略枚举 |

---

#### 1. 初始化
---
```go
package main

import (
    "fmt"
    "github.com/aid297/aid/v3/coroutineGroups"
)

func main() {
    // 创建协程组，并发限制为4
    g := coroutineGroups.New[int](4)

    // limit=0 时默认为4
    g2 := coroutineGroups.New[int](0)

    // 创建时直接配置重试（等价于 SetRetry）
    g3 := coroutineGroups.New[int](4,
        coroutineGroups.WithRetryAttempts(3),
        coroutineGroups.WithTimeout(5*time.Second),
    )

    // 通过实例方法创建新协程组
    g4 := g.New(8) // 并发限制为8

    _ = g2
    _ = g3
    _ = g4
    fmt.Println("created")
}
```
---

#### 2. Result 结果对象
---
```go
package main

import (
    "errors"
    "fmt"
    "github.com/aid297/aid/v3/coroutineGroups"
)

func main() {
    // 成功结果
    r := &coroutineGroups.ResultImpl[int]{Value: 42}
    fmt.Println(r.IsOK())      // true
    fmt.Println(r.GetValue())  // 42

    // 错误结果
    r2 := &coroutineGroups.ResultImpl[int]{Error: errors.New("fail")}
    fmt.Println(r2.IsOK())           // false
    fmt.Println(r2.GetError())       // fail

    // 跳过结果（Skip 不等于错误）
    r3 := &coroutineGroups.ResultImpl[int]{Value: 0, Skip: true}
    fmt.Println(r3.IsOK())    // false（Skip 时 IsOK 为 false）
    fmt.Println(r3.IsSkip())  // true

    // SetOK 清除错误
    r2.SetOK(true)
    fmt.Println(r2.GetError()) // <nil>

    // SetSkip 设置跳过状态
    r.SetSkip(true)
    fmt.Println(r.IsSkip()) // true
}
```
---

#### 3. GO 批量并发执行
---
```go
package main

import (
    "fmt"
    "github.com/aid297/aid/v3/coroutineGroups"
)

func main() {
    results := coroutineGroups.New[int](4).
        GO(
            func() coroutineGroups.Result[int] { return &coroutineGroups.ResultImpl[int]{Value: 1} },
            func() coroutineGroups.Result[int] { return &coroutineGroups.ResultImpl[int]{Value: 2} },
            func() coroutineGroups.Result[int] { return &coroutineGroups.ResultImpl[int]{Value: 3} },
        )

    // 结果长度等于函数数量（顺序不保证，并发执行）
    fmt.Println(results.Length()) // 3

    // 遍历结果
    for _, r := range results.ToSlice() {
        fmt.Printf("value=%d, ok=%v\n", r.GetValue(), r.IsOK())
    }
}
```

> **注意**：`GO` 通过信号量（semaphore）限制并发数。传入新的 `funcs` 会覆盖 `SetFunc` 预设的函数。结果顺序不保证，因为 goroutine 并发完成。
---

#### 4. SetFunc 预设函数
---
```go
package main

import (
    "fmt"
    "github.com/aid297/aid/v3/coroutineGroups"
)

func main() {
    g := coroutineGroups.New[int](4).SetFunc(
        func() coroutineGroups.Result[int] { return &coroutineGroups.ResultImpl[int]{Value: 10} },
        func() coroutineGroups.Result[int] { return &coroutineGroups.ResultImpl[int]{Value: 20} },
    )

    // GO 传入新函数时会覆盖 SetFunc 的预设
    results := g.GO(
        func() coroutineGroups.Result[int] { return &coroutineGroups.ResultImpl[int]{Value: 100} },
    )

    fmt.Println(results.Length()) // 1
}
```
---

#### 5. GOBatch 分轮批量执行
---
```go
package main

import (
    "fmt"
    "github.com/aid297/aid/v3/coroutineGroups"
)

func main() {
    // total=6, capacity=3 → 2批次，每批3个，共6个结果
    results, err := coroutineGroups.New[int](4).GOBatch(6, 3, func(batch, capacity uint) coroutineGroups.Result[int] {
        return &coroutineGroups.ResultImpl[int]{Value: int(batch*3 + capacity)}
    })
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println(results.Length()) // 6
    // 每批内并发执行，批与批之间串行（等待上一批全部完成）
}
```

**批次计算规则**：
- `total > capacities`：批次数 = `ceil(total / capacities)`
- `total <= capacities`：批次数 = 1
- 每批循环 `capacities` 次，总结果数 = `批次数 × capacities`
---

#### 6. SetRetry 超时重试
---
```go
package main

import (
    "errors"
    "fmt"
    "sync/atomic"
    "time"
    "github.com/aid297/aid/v3/coroutineGroups"
)

func main() {
    var callCount int32

    results := coroutineGroups.New[int](4).SetRetry(
        coroutineGroups.WithRetryAttempts(3),              // 最大执行3次（首次+2次重试）
        coroutineGroups.WithTimeout(5*time.Second),        // 每次执行超时5秒
        coroutineGroups.WithRetryInterval(100*time.Millisecond), // 重试基础间隔
        coroutineGroups.WithBackoff(coroutineGroups.BackoffExponential), // 指数退避
    ).GO(
        func() coroutineGroups.Result[int] {
            n := atomic.AddInt32(&callCount, 1)
            if n < 3 {
                return &coroutineGroups.ResultImpl[int]{Error: errors.New("transient")}
            }
            return &coroutineGroups.ResultImpl[int]{Value: 42}
        },
    )

    r := results.ToSlice()[0]
    fmt.Println(r.IsOK(), r.GetValue()) // true 42
}
```

**重试触发条件**：
- 函数返回 `Error != nil` → 重试
- 函数执行超时 → 重试（结果为 `ErrTimeout`）

**不触发重试**：
- 函数返回 `IsOK() == true` → 直接返回
- 函数返回 `IsSkip() == true` → 直接返回（Skip 不视为错误）

**独立性保证**：每个协程独立维护超时计时器和重试计数，一个协程的超时不会影响其他协程。
---

#### 7. 重试选项详解
---
```go
// 最大执行次数（含首次，小于1自动设为1）
coroutineGroups.WithRetryAttempts(3)        // 首次 + 2次重试

// 每次执行的超时时间（0表示不超时）
coroutineGroups.WithTimeout(5 * time.Second)

// 重试基础间隔时间
coroutineGroups.WithRetryInterval(100 * time.Millisecond)

// 退避策略
coroutineGroups.WithBackoff(coroutineGroups.BackoffFixed)
```
---

#### 8. 退避策略
---
```go
// BackoffFixed：固定间隔
// 第1次重试等待 interval
// 第2次重试等待 interval
// 第3次重试等待 interval

coroutineGroups.WithBackoff(coroutineGroups.BackoffFixed)

// ---

// BackoffExponential：指数退避
// 第1次重试等待 interval × 2^0 = interval
// 第2次重试等待 interval × 2^1 = 2 × interval
// 第3次重试等待 interval × 2^2 = 4 × interval

coroutineGroups.WithBackoff(coroutineGroups.BackoffExponential)

// ---

// BackoffExponentialJitter：指数退避 + 随机抖动
// 第n次重试等待 interval × 2^(n-1) + random[0, base]
// 避免多个协程同时重试导致的惊群效应

coroutineGroups.WithBackoff(coroutineGroups.BackoffExponentialJitter)
```

**示例：固定退避 vs 指数退避耗时对比**
```go
// 固定退避：3次执行 + 2次间隔100ms = 至少200ms
g1 := coroutineGroups.New[int](4).SetRetry(
    coroutineGroups.WithRetryAttempts(3),
    coroutineGroups.WithRetryInterval(100*time.Millisecond),
    coroutineGroups.WithBackoff(coroutineGroups.BackoffFixed),
)

// 指数退避：3次执行 + 100ms + 200ms = 至少300ms
g2 := coroutineGroups.New[int](4).SetRetry(
    coroutineGroups.WithRetryAttempts(3),
    coroutineGroups.WithRetryInterval(100*time.Millisecond),
    coroutineGroups.WithBackoff(coroutineGroups.BackoffExponential),
)
```
---

#### 9. 超时重试 + GOBatch
---
```go
package main

import (
    "fmt"
    "sync/atomic"
    "time"
    "github.com/aid297/aid/v3/coroutineGroups"
)

func main() {
    var callCount int32

    results, err := coroutineGroups.New[int](4).SetRetry(
        coroutineGroups.WithRetryAttempts(3),
        coroutineGroups.WithTimeout(2*time.Second),
        coroutineGroups.WithRetryInterval(50*time.Millisecond),
        coroutineGroups.WithBackoff(coroutineGroups.BackoffExponential),
    ).GOBatch(10, 5, func(batch, capacity uint) coroutineGroups.Result[int] {
        n := atomic.AddInt32(&callCount, 1)
        if n < 3 {
            return &coroutineGroups.ResultImpl[int]{Error: fmt.Errorf("batch %d cap %d fail", batch, capacity)}
        }
        return &coroutineGroups.ResultImpl[int]{Value: int(batch*5 + capacity)}
    })
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println("results:", results.Length())
}
```
---

#### 10. 错误定义
---
```go
var (
    ErrBatchInvalid    = errors.New("轮数不能为0")     // GOBatch 批次数为0
    ErrCapacityInvalid = errors.New("每轮循环数不能为0") // GOBatch 容量为0（注：当前会先触发除零panic）
    ErrEmptyFuncs      = errors.New("没有任何待执行任务") // GOBatch total=0
    ErrTimeout         = errors.New("协程执行超时")      // 重试超时
)

// 使用示例
results, err := coroutineGroups.New[int](4).GOBatch(0, 5, fn)
if errors.Is(err, coroutineGroups.ErrEmptyFuncs) {
    // total=0
}
```
---

#### 11. 全局变量
---
```go
// app.go 中定义了全局协程组实例（any 类型）
var CoroutineGroup CoroutineGrouper[any] = (*CoroutineGroupImpl[any])(nil)

// 可用于全局注册和依赖注入
```
---

#### 12. 完整示例
---
```go
package main

import (
    "errors"
    "fmt"
    "sync/atomic"
    "time"
    "github.com/aid297/aid/v3/coroutineGroups"
)

func main() {
    var callCount int32

    // 创建协程组：并发限制4，最多重试3次，每次超时2秒，指数退避
    g := coroutineGroups.New[int](4,
        coroutineGroups.WithRetryAttempts(3),
        coroutineGroups.WithTimeout(2*time.Second),
        coroutineGroups.WithRetryInterval(100*time.Millisecond),
        coroutineGroups.WithBackoff(coroutineGroups.BackoffExponential),
    )

    results := g.GO(
        // 任务1：前2次失败，第3次成功
        func() coroutineGroups.Result[int] {
            n := atomic.AddInt32(&callCount, 1)
            if n < 3 {
                return &coroutineGroups.ResultImpl[int]{Error: errors.New("transient")}
            }
            return &coroutineGroups.ResultImpl[int]{Value: 100}
        },
        // 任务2：一次成功
        func() coroutineGroups.Result[int] {
            return &coroutineGroups.ResultImpl[int]{Value: 200}
        },
        // 任务3：跳过（不重试）
        func() coroutineGroups.Result[int] {
            return &coroutineGroups.ResultImpl[int]{Value: 0, Skip: true}
        },
    )

    for _, r := range results.ToSlice() {
        switch {
        case r.IsSkip():
            fmt.Println("skipped")
        case r.IsOK():
            fmt.Printf("ok: %d\n", r.GetValue())
        default:
            fmt.Printf("error: %v\n", r.GetError())
        }
    }
    // ok: 100
    // ok: 200
    // skipped
}
```
---

#### 接口方法一览

| 方法 | 说明 |
|---|---|
| `New(limit uint16, attrs ...CoroutineGroupRetryAttr)` | 创建协程组，可选重试配置 |
| `New(limit uint16) CoroutineGrouper[T]` | 实例方法创建新协程组 |
| `SetFunc(funcs ...Func[T])` | 预设任务函数 |
| `SetRetry(attrs ...CoroutineGroupRetryAttr)` | 设置/更新超时重试配置 |
| `GO(funcs ...Func[T]) AnySlicer[Result[T]]` | 并发执行，信号量限制并发数 |
| `GOBatch(total, capacities int, fn) (AnySlicer[Result[T]], error)` | 分轮批量执行 |

#### Result 方法一览

| 方法 | 说明 |
|---|---|
| `IsOK() bool` | 无错误且未跳过 |
| `IsSkip() bool` | 是否跳过 |
| `GetError() error` | 获取错误 |
| `GetValue() T` | 获取值 |
| `SetOK(ok bool)` | 清除错误（设为nil） |
| `SetSkip(skip bool)` | 设置跳过状态 |

#### 重试选项一览

| 选项 | 说明 |
|---|---|
| `WithRetryAttempts(n int)` | 最大执行次数（含首次） |
| `WithTimeout(d time.Duration)` | 每次执行超时时间 |
| `WithRetryInterval(d time.Duration)` | 重试基础间隔 |
| `WithBackoff(s BackoffStrategy)` | 退避策略 |

---
