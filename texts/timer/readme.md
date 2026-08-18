# 智能时间解析与格式化库

强大的 Go 时间字符串解析和格式化工具库，支持组合时间格式、纯数字默认值和人性化中文输出。

## ✨ 核心特性

- 🔄 **组合时间解析**：支持 `1w2d3h`、`3d2h15m` 等多单位组合格式
- 🔢 **纯数字支持**：纯数字默认按秒处理，如 `"100"` → 100秒
- 🇨🇳 **中文输出**：将 duration 转换为人友好的中文时间字符串
- 🎯 **智能分级**：自动降级过大数值，如 `80分钟` → `1小时20分`
- ✅ **格式验证**：提供格式验证函数，提前校验输入合法性

## 📦 安装

```go
import "github.com/aid297/aid/v2/texts"
```

## 🚀 快速开始

### 1. 解析时间字符串

```go
// 单一单位
dur, err := texts.Timer.WhatTimeIsIt("1w")
// dur = 7 * 24 * timer.Hour (1周)

// 组合时间
dur, err := texts.Timer.WhatTimeIsIt("1w2d3h")
// dur = 9天3小时

// 纯数字（默认秒）
dur, err := texts.Timer.WhatTimeIsIt("100")
// dur = 100 * timer.Second
```

### 2. 格式化时间为中文

```go
// 单个单位
output := texts.Timer.DurationToString(1 * time.Hour)
// output = "1小时"

// 组合单位
output := texts.Timer.DurationToString(1*time.Hour + 20*time.Minute)
// output = "1小时+20分"

// 分级处理
output := texts.Timer.DurationToString(80 * time.Minute)
// output = "1小时+20分"
```

### 3. 格式验证

```go
valid := texts.Timer.IsValidTimeFormat("1w2d")     // true
valid := texts.Timer.IsValidTimeFormat("100")      // true (纯数字)
valid := texts.Timer.IsValidTimeFormat("abc")      // false
valid := texts.Timer.IsValidTimeFormat("1.5w")     // false
```

### 4. 默认值处理

```go
// 格式无效时返回默认值
dur := texts.Timer.WhatTimeIsItDefault("invalid", 1*time.Hour)
// dur = 1小时
```

## 📖 API 参考

### 公开函数

#### `WhatTimeIsIt(origin string) (time.Duration, error)`

解析时间字符串并返回 `time.Duration`。

**支持的格式：**
- 单一单位：`1s`、`5m`、`2h`、`3d`、`1w`
- 组合单位：`1w2d`、`3d2h15m`、`1h20m30s`
- 纯数字：`100`（默认按秒处理）

**示例：**
```go
// 单一单位
dur, err := texts.Timer.WhatTimeIsIt("1w")    // 7天
dur, err := texts.Timer.WhatTimeIsIt("2h30m") // 2小时30分

// 组合单位
dur, err := texts.Timer.WhatTimeIsIt("1w2d3h") // 9天3小时

// 纯数字
dur, err := texts.Timer.WhatTimeIsIt("3600")   // 3600秒 = 1小时
```

---

#### `DurationToString(dur time.Duration) string`

将 `time.Duration` 转换为人性化的中文时间字符串。

**特点：**
- 从大到小排序：周、天、时、分、秒
- 跳过为0的单位
- 用 `+` 连接多个单位
- 自动分级处理（见下方规则）

**示例：**
```go
texts.Timer.DurationToString(100 * time.Second)        // "1分+40秒"
texts.Timer.DurationToString(2 * time.Hour)            // "2小时"
texts.Timer.DurationToString(9*24*time.Hour + 3*time.Hour) // "9天+3小时"
```

---

#### `IsValidTimeFormat(origin string) bool`

验证字符串是否符合支持的时间格式。

**支持的格式：**
- 组合时间：`1w2d`、`3d2h15m`
- 纯数字：`100`、`3600`
- 单一单位：`1s`、`5m`、`2h`、`3d`、`1w`

**示例：**
```go
texts.Timer.IsValidTimeFormat("1w2d")    // true
texts.Timer.IsValidTimeFormat("100")     // true
texts.Timer.IsValidTimeFormat("abc")     // false
texts.Timer.IsValidTimeFormat("")        // false
texts.Timer.IsValidTimeFormat("1.5w")    // false
```

---

#### `WhatTimeIsItDefault(origin string, def time.Duration) time.Duration`

解析时间字符串，如果格式无效则返回默认值。

**示例：**
```go
// 有效格式
dur := texts.Timer.WhatTimeIsItDefault("1w", 1*time.Hour)
// dur = 7天

// 无效格式
dur := texts.Timer.WhatTimeIsItDefault("invalid", 1*time.Hour)
// dur = 1小时（使用默认值）
```

## 🎯 分级处理规则

`DurationToString` 函数会自动进行分级处理，确保输出的时间单位在合理范围内：

### 规则

| 原始值 | 分级后 | 说明 |
|--------|--------|------|
| 80分钟 | 1小时+20分 | ≥10分钟 → 进位到小时 |
| 90分钟 | 1小时+30分 | ≥10分钟 → 进位到小时 |
| 25小时 | 1天+1小时 | ≥10小时 → 进位到天 |
| 30小时 | 1天+6小时 | ≥10小时 → 进位到天 |
| 10天 | 1周+3天 | ≥10天 → 进位到周 |
| 17天 | 2周+3天 | ≥10天 → 进位到周 |

### 工作原理

```go
// 输入：80分钟
dur := 80 * time.Minute
output := texts.Timer.DurationToString(dur)
// 输出: "1小时+20分"

// 输入：10天
dur := 10 * 24 * time.Hour
output := texts.Timer.DurationToString(dur)
// 输出: "1周+3天"
```

## 📊 时间单位说明

| 符号 | 含义 | 等价时间 |
|------|------|---------|
| `s` | 秒 | 1秒 |
| `m` | 分 | 60秒 |
| `h` | 小时 | 60分钟 |
| `d` | 天 | 24小时 |
| `w` | 周 | 168小时（7天） |

## 💡 完整示例

```go
package main

import (
    "fmt"
    "github.com/aid297/aid/v2/consts/textTime"
)

func main() {
    // 场景1：解析用户输入的组合时间
    userInput := "1w2d3h"
    dur, err := texts.Timer.WhatTimeIsIt(userInput)
    if err != nil {
        fmt.Printf("解析失败: %v\n", err)
    } else {
        output := texts.Timer.DurationToString(dur)
        fmt.Printf("%s → %s\n", userInput, output)
        // 输出: 1w2d3h → 1周+2天+3小时
    }

    // 场景2：纯数字默认按秒处理
    seconds := "100"
    dur, _ = texts.Timer.WhatTimeIsIt(seconds)
    output := texts.Timer.DurationToString(dur)
    fmt.Printf("%s → %s\n", seconds, output)
    // 输出: 100 → 1分+40秒

    // 场景3：格式验证
    if texts.Timer.IsValidTimeFormat("1w2d") {
        fmt.Println("格式合法")
    }

    // 场景4：智能分级处理
    dur = 80 * time.Minute
    output = texts.Timer.DurationToString(dur)
    fmt.Printf("80分钟 → %s\n", output)
    // 输出: 80分钟 → 1小时+20分
}
```

## 🧪 测试

运行测试用例：

```bash
cd consts/textTime
go test -v
```

运行性能基准测试：

```bash
go test -bench=.
```

## 📝 对比示例

| 输入 | 解析结果 | 中文输出 |
|------|---------|---------|
| `"1s"` | 1秒 | `"1秒"` |
| `"5m"` | 5分钟 | `"5分"` |
| `"2h"` | 2小时 | `"2小时"` |
| `"3d"` | 3天 | `"3天"` |
| `"1w"` | 7天 | `"1周"` |
| `"1w2d"` | 9天 | `"1周+2天"` |
| `"1h20m"` | 1小时20分 | `"1小时+20分"` |
| `"100"` | 100秒 | `"1分+40秒"` |
| `"3600"` | 3600秒 | `"1小时"` |
| `"80m"` | 80分钟 | `"1小时+20分"` |

## ⚠️ 注意事项

1. **大小写不敏感**：`1W2D` 和 `1w2d` 会被同样解析
2. **无空格格式**：输入不能包含空格，如 `"1w 2d"` 是非法的
3. **不支持小数**：`"1.5w"` 是非法格式，请使用整数
4. **最小单位秒**：纯数字默认按秒处理，无需显式写出单位
5. **分级阈值**：当某一位 ≥ 10 时，会自动进位到上一级单位

## 🔍 错误处理

当输入格式不合法时，`WhatTimeIsIt` 会返回错误：

```go
dur, err := texts.Timer.WhatTimeIsIt("abc")
// err: "未找到有效的时间单元"

dur, err := texts.Timer.WhatTimeIsIt("1x")
// err: "时间单位不支持：x"

dur, err := texts.Timer.WhatTimeIsIt("")
// err: "未找到有效的时间单元"
```

## 📚 应用场景

- ⏰ **定时任务配置**：解析用户输入的时间配置
- 📅 **日历应用**：将 Duration 转换为用户友好的时间显示
- ⏱️ **倒计时功能**：格式化剩余时间
- 🔧 **参数校验**：验证时间格式合法性
- 📊 **日志分析**：解析和展示时间间隔

## 📄 License

MIT License
