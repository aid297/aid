# httpClientV4 快速参考

## 基本用法

```go
import "github.com/aid297/aid/httpClient/httpClientV4"

// GET 请求
client := httpClientV4.GET(httpClientV4.URL("https://api.example.com/users"))
defer client.Release()  // ⚠️ 重要: 归还到对象池
client.Send()
body := client.ToBytes()

// POST JSON
client := httpClientV4.POST(
    httpClientV4.URL("https://api.example.com/users"),
    httpClientV4.JSON(map[string]string{"name": "John"}),
)
defer client.Release()
client.Send()
```

## HTTP 方法

```go
httpClientV4.GET(attrs...)
httpClientV4.POST(attrs...)
httpClientV4.PUT(attrs...)
httpClientV4.PATCH(attrs...)
httpClientV4.DELETE(attrs...)
httpClientV4.HEAD(attrs...)
httpClientV4.OPTIONS(attrs...)
httpClientV4.TRACE(attrs...)
```

## 属性设置

### URL & Queries

```go
// URL
httpClientV4.URL("https://api.example.com")
httpClientV4.URL("https://", "api.example.com", "/users")  // 拼接

// Queries (必须是 string 类型)
httpClientV4.Queries(map[string]string{
    "page": "1",
    "size": "20",
})
```

### Headers

```go
// 单个 header
httpClientV4.AppendHeader("User-Agent", "MyApp/1.0")
httpClientV4.SetHeader("Authorization", "Bearer token")

// 多个 headers
headers := http.Header{
    "User-Agent": []string{"MyApp/1.0"},
    "Accept":     []string{"application/json"},
}
httpClientV4.AppendHeaders(headers)

// Content-Type 和 Accept
httpClientV4.ContentType_(httpClientV4.ContentTypeJSON)
httpClientV4.Accept_(httpClientV4.AcceptJSON)

// Authorization
httpClientV4.Authorization("username", "password", "Basic")
```

### Body

```go
// JSON
httpClientV4.JSON(struct{ Name string }{"John"})

// XML
httpClientV4.XML(data)

// Form
httpClientV4.Form(map[string]string{
    "username": "user",
    "password": "pass",
})

// FormData (multipart)
httpClientV4.FormData(
    map[string]string{"field": "value"},  // fields
    map[string]string{"file": "/path/to/file.txt"},  // files
)

// Plain text
httpClientV4.Plain("Hello World")

// HTML
httpClientV4.HTML("<html>...</html>")

// 原始字节
httpClientV4.Bytes([]byte{0x01, 0x02})

// 从文件
httpClientV4.File("/path/to/file.txt")

// 从 Reader
httpClientV4.Reader(reader)
```

### 其他配置

```go
// Timeout
httpClientV4.Timeout(10 * time.Second)

// Transport
httpClientV4.Transport(customTransport)
httpClientV4.TransportDefault()  // 使用默认配置

// TLS 证书
httpClientV4.Cert(certBytes)

// 自动复制响应体 (允许多次读取)
httpClientV4.AutoCopy(true)
```

## 发送请求

```go
// 普通发送
client.Send()

// 带重试
client.SendWithRetry(
    3,              // 重试次数
    time.Second,    // 重试间隔
    func(statusCode int, err error) bool {
        return statusCode >= 500 || err != nil  // 重试条件
    },
)
```

## 处理响应

```go
// 检查错误
if err := client.Error(); err != nil {
    log.Fatal(err)
}

// 状态码和状态
statusCode := client.GetStatusCode()  // 200
status := client.GetStatus()          // "200 OK"

// 获取原始对象
req := client.GetRawRequest()
resp := client.GetRawResponse()

// 读取 body
bytes := client.ToBytes()

// 解析 JSON
var result map[string]interface{}
client.ToJSON(&result)

// 解析 XML
var data MyStruct
client.ToXML(&data)

// 写入 Writer
client.ToWriter(httpResponseWriter)
```

## 对象池

```go
// 方式 1: 使用便捷方法 (推荐)
client := httpClientV4.GET(...)
defer client.Release()

// 方式 2: 手动管理
client := httpClientV4.Acquire()
client.SetAttrs(
    httpClientV4.URL("https://api.example.com"),
    httpClientV4.Method(http.MethodGet),
)
defer client.Release()
```

## Content-Type 常量

```go
httpClientV4.ContentTypeJSON               // application/json
httpClientV4.ContentTypeXML                // application/xml
httpClientV4.ContentTypeXWwwFormURLencoded // application/x-www-form-urlencoded
httpClientV4.ContentTypeFormData           // multipart/form-data
httpClientV4.ContentTypePlain              // text/plain
httpClientV4.ContentTypeHTML               // text/html
httpClientV4.ContentTypeCSS                // text/css
httpClientV4.ContentTypeJavascript         // text/javascript
httpClientV4.ContentTypeSteam              // application/octet-stream
```

## Accept 常量

```go
httpClientV4.AcceptJSON       // application/json
httpClientV4.AcceptXML        // application/xml
httpClientV4.AcceptPlain      // text/plain
httpClientV4.AcceptHTML       // text/html
httpClientV4.AcceptCSS        // text/css
httpClientV4.AcceptJavascript // text/javascript
httpClientV4.AcceptSteam      // application/octet-stream
httpClientV4.AcceptAny        // */*
```

## 完整示例

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/aid297/aid/httpClient/httpClientV4"
)

func main() {
    // 创建客户端
    client := httpClientV4.POST(
        httpClientV4.URL("https://api.example.com/users"),
        httpClientV4.Queries(map[string]string{"debug": "true"}),
        httpClientV4.JSON(map[string]string{
            "name":  "John Doe",
            "email": "john@example.com",
        }),
        httpClientV4.AppendHeader("User-Agent", "MyApp/1.0"),
        httpClientV4.Timeout(10*time.Second),
    )
    defer client.Release()  // ⚠️ 重要
    
    // 发送请求 (带重试)
    client.SendWithRetry(3, time.Second, func(statusCode int, err error) bool {
        return statusCode >= 500 || err != nil
    })
    
    // 检查错误
    if err := client.Error(); err != nil {
        log.Fatalf("请求失败: %v", err)
    }
    
    // 检查状态
    if client.GetStatusCode() != 200 {
        log.Fatalf("非预期状态码: %d", client.GetStatusCode())
    }
    
    // 解析响应
    var result struct {
        ID   int    `json:"id"`
        Name string `json:"name"`
    }
    client.ToJSON(&result)
    
    if err := client.Error(); err != nil {
        log.Fatalf("解析失败: %v", err)
    }
    
    fmt.Printf("创建成功: ID=%d, Name=%s\n", result.ID, result.Name)
}
```

## 注意事项

⚠️ **必须调用 Release()**: 使用完毕后必须归还到对象池

```go
// ✅ 正确
client := httpClientV4.GET(...)
defer client.Release()

// ❌ 错误 - 内存泄漏
client := httpClientV4.GET(...)
client.Send()
// 忘记 Release()
```

⚠️ **不要在 Release 后使用**

```go
// ❌ 错误
client := httpClientV4.GET(...)
client.Send()
client.Release()
client.GetStatusCode()  // 危险! 可能访问已被重用的对象
```

⚠️ **线程安全**: 单个客户端实例不是线程安全的

```go
// ❌ 错误
client := httpClientV4.GET(...)
defer client.Release()
go client.Send()  // 危险!
go client.Send()  // 危险!
```

⚠️ **Queries 必须是 string**

```go
// ❌ 错误
Queries(map[string]any{"page": 1})

// ✅ 正确
Queries(map[string]string{"page": "1"})
Queries(map[string]string{"page": strconv.Itoa(1)})
```
