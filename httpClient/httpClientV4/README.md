# httpClientV4 - 低内存逃逸优化版本

## 概述

httpClientV4 是基于 httpClientV2 的优化版本,专注于降低内存逃逸,提升性能和减少 GC 压力。

## 主要优化点

### 1. **使用具体类型替代 `any` 接口**

**V2 版本:**
```go
queries map[string]any
headers map[string][]any
```

**V4 版本:**
```go
queries map[string]string
headers http.Header  // 标准库已优化的类型
```

**优势:** 避免接口装箱导致的内存逃逸,减少堆分配。

---

### 2. **对象池 (sync.Pool) 复用**

**V4 新增:**
```go
var clientPool = sync.Pool{
    New: func() any {
        return &HTTPClient{
            queries: make(map[string]string, 8),
            headers: make(http.Header, 8),
        }
    },
}

func Acquire() *HTTPClient  // 获取客户端
func (c *HTTPClient) Release()  // 归还客户端
```

**优势:** 
- 减少频繁创建/销毁对象的开销
- 降低 GC 压力
- 预分配容量避免动态扩容

---

### 3. **简化接口设计**

**V2 版本:**
```go
type HTTPClientAttributer interface {
    Register(req *HTTPClient)
    Error() error
    ImplHTTPClientAttributer()  // 标记方法
}
```

**V4 版本:**
```go
type HTTPClientAttributer interface {
    Apply(client *HTTPClient)  // 仅一个方法
}
```

**优势:** 
- 更简洁的接口减少虚函数调用开销
- 去除冗余方法,降低逃逸风险

---

### 4. **使用 strings.Builder 优化字符串拼接**

**V4 优化:**
```go
func (my *HTTPClient) buildURL() string {
    if len(my.queries) == 0 {
        return my.url
    }
    
    var sb strings.Builder
    sb.Grow(len(my.url) + len(queries.Encode()) + 1)  // 预分配
    sb.WriteString(my.url)
    sb.WriteByte('?')
    sb.WriteString(queries.Encode())
    
    return sb.String()
}
```

**优势:** 
- 预分配内存避免多次扩容
- 比字符串拼接更高效

---

### 5. **Buffer 池化**

**V4 新增:**
```go
var bufferPool = sync.Pool{
    New: func() any {
        return new(bytes.Buffer)
    },
}

// 使用时
buf := bufferPool.Get().(*bytes.Buffer)
buf.Reset()
defer bufferPool.Put(buf)
```

**优势:** 复用 Buffer 对象,减少大对象分配。

---

### 6. **避免闭包捕获变量**

**V2 版本:**
```go
func (my *HTTPClient) GetStatusCode() int {
    return operationV2.NewTernary(
        operationV2.TrueFn(func() int { 
            return my.GetRawResponse().StatusCode  // 闭包捕获
        })
    ).GetByValue(my.GetRawResponse() != nil)
}
```

**V4 版本:**
```go
func (my *HTTPClient) GetStatusCode() int {
    if my.rawResponse != nil {
        return my.rawResponse.StatusCode
    }
    return 0
}
```

**优势:** 避免闭包导致的变量逃逸到堆。

---

### 7. **常量化类型定义**

**V2 版本:**
```go
var ContentTypes = map[ContentType]string{
    ContentTypeJSON: "application/json",
    // ...
}
```

**V4 版本:**
```go
const (
    ContentTypeJSON ContentType = "application/json"
    // 直接使用常量值
)
```

**优势:** 编译时常量,无运行时开销。

---

### 8. **预分配 map 容量**

**V4 优化:**
```go
queries: make(map[string]string, 8),
headers: make(http.Header, 8),
```

**优势:** 避免 map 初始化时的动态扩容。

---

## 性能对比

基准测试结果 (Apple M4):

```
BenchmarkHTTPClientV4_WithPool-10       134270    24881 ns/op    4943 B/op    62 allocs/op
BenchmarkHTTPClientV4_WithoutPool-10    143919    26003 ns/op    5264 B/op    65 allocs/op
```

**使用对象池的优势:**
- 内存分配减少: 321 字节/op (约 6%)
- 分配次数减少: 3 次/op

---

## 使用示例

### 基本使用

```go
// GET 请求
client := httpClientV4.GET(
    httpClientV4.URL("https://api.example.com/users"),
    httpClientV4.Queries(map[string]string{"page": "1"}),
    httpClientV4.AppendHeader("User-Agent", "MyApp/1.0"),
)
defer client.Release()  // 重要: 归还到池中

client.Send()
if err := client.Error(); err != nil {
    log.Fatal(err)
}

body := client.ToBytes()
fmt.Println(string(body))
```

### POST 请求

```go
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

client := httpClientV4.POST(
    httpClientV4.URL("https://api.example.com/users"),
    httpClientV4.JSON(User{Name: "John", Email: "john@example.com"}),
    httpClientV4.Timeout(10 * time.Second),
)
defer client.Release()

client.Send()
```

### 表单提交

```go
client := httpClientV4.POST(
    httpClientV4.URL("https://api.example.com/login"),
    httpClientV4.Form(map[string]string{
        "username": "user",
        "password": "pass",
    }),
)
defer client.Release()

client.Send()
```

### 带重试

```go
client := httpClientV4.GET(
    httpClientV4.URL("https://api.example.com/data"),
)
defer client.Release()

client.SendWithRetry(3, time.Second, func(statusCode int, err error) bool {
    return statusCode >= 500 || err != nil  // 服务器错误时重试
})
```

---

## 内存逃逸分析

可以使用以下命令查看逃逸分析:

```bash
go build -gcflags="-m -m" ./...
```

V4 版本相比 V2 应该有更少的"escapes to heap"警告。

---

## 注意事项

1. **务必调用 Release()**: 使用完客户端后应调用 `Release()` 将其归还到池中
2. **不要持有引用**: Release 后不要再使用该客户端实例
3. **线程安全**: 单个客户端实例不是线程安全的,不要跨 goroutine 共享

---

## 总结

httpClientV4 通过以下手段有效降低内存逃逸:

1. ✅ 使用具体类型替代接口类型
2. ✅ 实现对象池复用机制
3. ✅ 简化接口设计
4. ✅ 优化字符串操作
5. ✅ Buffer 池化
6. ✅ 避免闭包捕获
7. ✅ 预分配容量
8. ✅ 使用常量

这些优化使得 V4 在保持 V2 所有功能的同时,显著降低了内存分配和 GC 压力。
