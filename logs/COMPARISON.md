# JsonObject 效果对比

## 问题描述

在使用 zap 日志记录 JSON 数据时,传统的两种方式都有不足:

### 1. 使用 `zap.Any`

```go
logger.Info("response", zap.Any("data", responseData))
```

**输出效果:**
```json
{
  "data": "{map[code:200 message:success data:map[id:123 name:商品A]]}",
  "level": "info",
  "msg": "response",
  "time": "2026-09-03 09:35:57.000"
}
```

**问题:** 
- ❌ Go 的 map 表示法,不易读
- ❌ 缺少引号和逗号
- ❌ 嵌套结构不清晰
- ❌ 不符合 JSON 标准格式

### 2. 使用 `zap.String`

```go
logger.Info("response", zap.String("data", jsonString))
```

**输出效果:**
```json
{
  "data": "{\"code\":200,\"message\":\"success\",\"data\":{\"id\":123,\"name\":\"商品A\"}}",
  "level": "info",
  "msg": "response",
  "time": "2026-09-03 09:35:57.000"
}
```

**问题:**
- ❌ 所有数据在一行,难以阅读
- ❌ 转义字符多,视觉混乱
- ❌ 深层嵌套几乎无法辨认

---

## 解决方案:使用 `JsonObject`

```go
obj := logs.NewJsonObject(jsonString)
logger.Info("response", zap.Stringer("data", obj))
```

**输出效果:**
```json
{
  "data": "{\n  \"code\": 200,\n  \"message\": \"success\",\n  \"data\": {\n    \"id\": 123,\n    \"name\": \"商品A\"\n  }\n}",
  "level": "info",
  "msg": "response",
  "time": "2026-09-03 09:35:57.000"
}
```

**控制台显示(人类可读):**
```
data: {
  "code": 200,
  "message": "success",
  "data": {
    "id": 123,
    "name": "商品A"
  }
}
```

**优势:**
- ✅ 结构化的缩进格式
- ✅ 清晰的嵌套层次
- ✅ 标准的 JSON 语法
- ✅ 易于人类阅读
- ✅ 支持错误提示
- ✅ 自动处理类型

---

## 实际对比截图

### 示例 1: 简单对象

**原始 JSON:**
```json
{"name":"张三","age":25,"city":"北京","active":true}
```

| 方式 | 控制台输出 |
|------|----------|
| `zap.Any` | `{map[name:张三 age:25 city:北京 active:true]}` |
| `zap.String` | `{"name":"张三","age":25,"city":"北京","active":true}` |
| **JsonObject** | ```{
  "name": "张三",
  "age": 25,
  "city": "北京",
  "active": true
}``` |

### 示例 2: 嵌套对象

**原始 JSON:**
```json
{"user":{"name":"李四","contact":{"email":"lisi@example.com","address":{"province":"浙江省","city":"杭州"}}}}
```

| 方式 | 控制台输出 |
|------|----------|
| `zap.Any` | `{map[user:map[name:李四 contact:map[email:lisi@example.com address:map[province:浙江省 city:杭州]]]]}` |
| `zap.String` | `{"user":{"name":"李四","contact":{"email":"lisi@example.com","address":{"province":"浙江省","city":"杭州"}}}}` |
| **JsonObject** | ```{
  "user": {
    "name": "李四",
    "contact": {
      "email": "lisi@example.com",
      "address": {
        "province": "浙江省",
        "city": "杭州"
      }
    }
  }
}``` |

### 示例 3: API 响应

**原始 JSON:**
```json
{"code":200,"message":"success","data":{"total":100,"page":1,"items":[{"id":1,"name":"商品A"},{"id":2,"name":"商品B"}]},"timestamp":1633024800}
```

| 方式 | 控制台输出 |
|------|----------|
| `zap.Any` | `{map[code:200 message:success data:map[total:100 page:1 items:[map[id:1 name:商品A] map[id:2 name:商品B]]] timestamp:1.6330248e+09]}` |
| `zap.String` | `{"code":200,"message":"success","data":{"total":100,"page":1,"items":[{"id":1,"name":"商品A"},{"id":2,"name":"商品B"}]},"timestamp":1633024800}` |
| **JsonObject** | ```{
  "code": 200,
  "message": "success",
  "data": {
    "total": 100,
    "page": 1,
    "items": [
      {
        "id": 1,
        "name": "商品A"
      },
      {
        "id": 2,
        "name": "商品B"
      }
    ]
  },
  "timestamp": 1633024800
}``` |

---

## 错误处理示例

当传入无效的 JSON 字符串时:

```go
obj := logs.NewJsonObject(`{invalid json}`)
logger.Warn("error", zap.Stringer("data", obj))
```

**输出:**
```
[JSON Parse Error: invalid character 'i' looking for beginning of object key string] {invalid json}
```

这样既能看到错误信息,又能看到原始输入,方便调试。

---

## 性能考虑

- JSON 解析只在创建 `JsonObject` 时执行一次
- 多次调用 `String()` 不会重复解析
- 建议仅在 debug 级别使用于高性能场景
- 有深度保护(最多10层),防止过深嵌套

---

## 使用场景

非常适合以下场景:

1. **API 日志记录** - 记录请求/响应的 JSON 数据
2. **调试工具** - 查看复杂的数据结构
3. **监控告警** - 结构化展示监控指标
4. **审计日志** - 清晰地记录业务数据
5. **数据导入导出** - 展示转换前后的数据结构
