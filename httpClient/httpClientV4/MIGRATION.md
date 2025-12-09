# httpClientV4 vs httpClientV2 对比

## API 变化对比

### 1. 基本使用

#### V2
```go
client := new(httpClientV2.HTTPClient).GET(
    httpClientV2.URL("https://api.example.com"),
)
```

#### V4
```go
client := httpClientV4.GET(
    httpClientV4.URL("https://api.example.com"),
)
defer client.Release() // 新增: 需要归还到池
```

---

### 2. Queries 参数类型

#### V2
```go
queries := map[string]any{
    "page": 1,           // 可以是任意类型
    "name": "test",
}
```

#### V4
```go
queries := map[string]string{
    "page": "1",         // 必须是 string
    "name": "test",
}
```

**迁移建议:** 使用 `strconv` 或 `fmt.Sprintf` 转换非字符串类型。

---

### 3. Headers 设置

#### V2
```go
// 方式 1: 使用 map[string]any
httpClientV2.AppendHeaderValue(map[string]any{
    "User-Agent": "MyApp",
})

// 方式 2: 使用 map[string][]any
httpClientV2.AppendHeaderValues(map[string][]any{
    "Accept": []any{"application/json"},
})
```

#### V4
```go
// 统一使用 http.Header 或简化方法
httpClientV4.AppendHeader("User-Agent", "MyApp")

// 或者使用 http.Header
headers := http.Header{
    "User-Agent": []string{"MyApp"},
    "Accept":     []string{"application/json"},
}
httpClientV4.AppendHeaders(headers)
```

---

### 4. 接口定义

#### V2
```go
type HTTPClientAttributer interface {
    Register(req *HTTPClient)
    Error() error
    ImplHTTPClientAttributer()
}
```

#### V4
```go
type HTTPClientAttributer interface {
    Apply(client *HTTPClient)  // 简化为单一方法
}
```

**影响:** 如果你自定义了属性,需要调整方法名从 `Register` 到 `Apply`。

---

### 5. 对象生命周期

#### V2
```go
client := new(httpClientV2.HTTPClient).GET(...)
// 使用后自动 GC
```

#### V4
```go
client := httpClientV4.GET(...)
defer client.Release()  // 必须归还到池
```

**重要:** V4 引入对象池,使用完必须调用 `Release()`,否则会造成内存泄漏。

---

## 性能对比

### 内存分配

```bash
# V2 典型场景
~8000 B/op, ~85 allocs/op

# V4 使用对象池
~4943 B/op, ~62 allocs/op

# 改进
- 内存减少: ~38%
- 分配次数减少: ~27%
```

### GC 压力

V4 通过对象池和减少接口装箱,显著降低了 GC 频率。

---

## 功能对比

| 功能 | V2 | V4 | 备注 |
|------|----|----|------|
| 基本 HTTP 方法 | ✅ | ✅ | 完全兼容 |
| Queries 参数 | `map[string]any` | `map[string]string` | V4 类型更严格 |
| Headers | `map[string][]any` | `http.Header` | V4 使用标准库 |
| Body 设置 | ✅ | ✅ | 完全兼容 |
| JSON/XML | ✅ | ✅ | 完全兼容 |
| Form/FormData | ✅ | ✅ | 完全兼容 |
| 重试机制 | ✅ | ✅ | 完全兼容 |
| Timeout | ✅ | ✅ | 完全兼容 |
| Transport | ✅ | ✅ | 完全兼容 |
| TLS/Cert | ✅ | ✅ | 完全兼容 |
| **对象池** | ❌ | ✅ | **V4 新增** |
| **Buffer 池** | ❌ | ✅ | **V4 新增** |

---

## 迁移指南

### 1. 简单迁移 (最小改动)

只需要添加 `defer Release()`:

```go
// V2
client := new(httpClientV2.HTTPClient).GET(...)
client.Send()

// V4
client := httpClientV4.GET(...)
defer client.Release()  // 新增这一行
client.Send()
```

### 2. Queries 类型转换

```go
// V2
queries := map[string]any{
    "page": 1,
    "size": 20,
}

// V4
import "strconv"

queries := map[string]string{
    "page": strconv.Itoa(1),
    "size": strconv.Itoa(20),
}

// 或者使用辅助函数
func toStringMap(m map[string]any) map[string]string {
    result := make(map[string]string, len(m))
    for k, v := range m {
        result[k] = fmt.Sprint(v)
    }
    return result
}
```

### 3. Headers 迁移

```go
// V2
httpClientV2.AppendHeaderValue(map[string]any{
    "Authorization": "Bearer token",
    "Content-Type": "application/json",
})

// V4 方式1: 逐个添加
httpClientV4.AppendHeader("Authorization", "Bearer token")
httpClientV4.AppendHeader("Content-Type", "application/json")

// V4 方式2: 使用 http.Header
headers := http.Header{
    "Authorization": []string{"Bearer token"},
    "Content-Type":  []string{"application/json"},
}
httpClientV4.AppendHeaders(headers)
```

### 4. 自定义属性

如果你实现了自定义属性:

```go
// V2
type MyAttr struct { value string }

func (a *MyAttr) Register(req *httpClientV2.HTTPClient) {
    // ...
}
func (a *MyAttr) Error() error { return nil }
func (a *MyAttr) ImplHTTPClientAttributer() {}

// V4
type MyAttr struct { value string }

func (a *MyAttr) Apply(client *httpClientV4.HTTPClient) {
    // ...
}
```

---

## 何时使用 V4

### 推荐使用 V4 的场景:

- ✅ 高并发场景(对象池优势明显)
- ✅ 内存敏感的应用
- ✅ 需要降低 GC 压力
- ✅ 新项目

### 继续使用 V2 的场景:

- ⚠️ 需要 `map[string]any` 的灵活性
- ⚠️ 迁移成本过高
- ⚠️ 低并发场景(性能差异不明显)

---

## 注意事项

### V4 使用要点:

1. **务必 Release**: 使用 `defer client.Release()` 归还对象
2. **不要跨 goroutine**: 单个客户端实例不是线程安全的
3. **类型转换**: Queries 必须是 `string` 类型
4. **不要持有引用**: Release 后不要再使用该实例

### 常见错误:

```go
// ❌ 错误: 忘记 Release
client := httpClientV4.GET(...)
client.Send()
// 内存泄漏!

// ✅ 正确
client := httpClientV4.GET(...)
defer client.Release()
client.Send()

// ❌ 错误: Release 后继续使用
client := httpClientV4.GET(...)
client.Send()
client.Release()
client.GetStatusCode()  // 危险!

// ✅ 正确
client := httpClientV4.GET(...)
defer client.Release()
client.Send()
statusCode := client.GetStatusCode()
```

---

## 总结

httpClientV4 在保持 V2 核心功能的基础上,通过优化内存管理显著提升了性能:

- **内存分配减少 ~38%**
- **分配次数减少 ~27%**
- **GC 压力大幅降低**

代价是:
- 需要手动管理对象生命周期 (`Release()`)
- Queries 类型限制为 `string`
- Headers 使用标准库 `http.Header`

对于大多数应用,这些代价是值得的。
