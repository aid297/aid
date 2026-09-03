# JSON 结构化日志字段

## 概述

`JsonObject` 是一个用于日志的结构化 JSON 字段类型,它会自动解析 JSON 字符串并以美观的格式输出,让人类更容易阅读。

## 解决的问题

### 传统方式的问题

```go
// 使用 zap.Any - 输出为 Go 的 map 表示
logger.Info("response", zap.Any("data", responseData))
// 输出: data={map[code:200 message:success]}

// 使用 zap.String - 输出为一行原始字符串
logger.Info("response", zap.String("data", jsonString))
// 输出: data={"code":200,"message":"success","data":{"id":1}}
```

这两种方式对于人类读取都很不友好。

### 使用 JsonObject

```go
logger.Info("response", logs.JsonObject{rawData: jsonString})
// 输出:
// data: {
//   "code": 200,
//   "message": "success",
//   "data": {
//     "id": 1
//   }
// }
```

## 基本用法

```go
package main

import (
    "github.com/aid297/aid/v2/logs"
    "go.uber.org/zap"
)

func main() {
    logger, _ := zap.NewDevelopment()
    
    // 简单的 JSON 对象
    jsonStr := `{"name":"张三","age":25,"city":"北京"}`
    obj := logs.NewJsonObject(jsonStr)
    logger.Info("user info", zap.Stringer("userInfo", obj))
    
    // 嵌套对象
    nestedJson := `{"user":{"name":"李四","contact":{"email":"lisi@example.com"}}}`
    nestedObj := logs.NewJsonObject(nestedJson)
    logger.Info("nested data", zap.Stringer("data", nestedObj))
    
    // 数组
    arrayJson := `[{"id":1,"name":"A"},{"id":2,"name":"B"}]`
    arrayObj := logs.NewJsonObject(arrayJson)
    logger.Info("items list", zap.Stringer("items", arrayObj))
    
    // 复杂混合结构
    complexJson := `{"code":200,"data":{"items":[1,2,3]},"success":true}`
    complexObj := logs.NewJsonObject(complexJson)
    logger.Info("api response", zap.Stringer("response", complexObj))
}
```

## 特性

### 1. 自动格式化

- 自动解析 JSON 字符串
- 以缩进格式输出,层次清晰
- 支持多层嵌套(最多10层)

### 2. 错误处理

```go
invalidObj := logs.NewJsonObject(`{invalid json}`)
logger.Info("error log", zap.Stringer("data", invalidObj))
// 输出:
// data: [JSON Parse Error: invalid character 'i' looking for beginning of object key string] {invalid json}
```

### 3. 深度保护

防止过深的嵌套导致输出过长:

```go
deepObj := logs.NewJsonObject(`{"a":{"b":{"c":{"d":{"e":{"f":{"g":{"h":{"i":{"j":{"k":"deep"}}}}}}}}}}}`)
// 输出会在第10层截断为 {...}
```

### 4. 方法接口

```go
obj := logs.NewJsonObject(jsonStr)

// 获取格式化后的字符串
str := obj.String() // 用于日志输出

// 获取解析后的值
value := obj.Value() // interface{}

// 获取原始字符串
raw := obj.Raw() // string

// 获取解析错误(如果有)
err := obj.Error() // error
```

## 与 Zap 日志集成

```go
// 创建 zap logger
logger, _ := logs.NewZapLog()

// 方式1: 使用 zap.Stringer
jsonStr := `{"key":"value"}`
obj := logs.NewJsonObject(jsonStr)
logger.Info("message", zap.Stringer("jsonData", obj))

// 方式2: 直接作为 field (需要自定义)
// 你可以扩展 JsonObject 实现 zap.ObjectEncoder
```

## 性能考虑

- JSON 解析在创建时完成
- 后续多次调用 `String()` 不会重复解析
- 对于大量日志场景,建议只在 debug 级别使用

## 注意事项

1. **输入必须是合法的 JSON**: 如果不是合法 JSON,会显示原始字符串和错误信息
2. **空字符串处理**: 返回空字符串,不会报错
3. **深度限制**: 最多支持10层嵌套,超过部分显示为 `{...}`
4. **数字类型**: JSON 中的数字会被正确识别为整数或浮点数

## 测试

运行单元测试:

```bash
go test -v ./logs -run TestJsonObject
```
