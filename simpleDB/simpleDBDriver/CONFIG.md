# DatabaseConfig 配置指南

## 概述

SimpleDB 提供 `DatabaseConfig` 来配置全局行为，包括 UUID 生成版本和级联查询最大深度。

## 配置项

### 1. DefaultUUIDVersion
**类型：** `int`  
**默认值：** `6`  
**范围：** `1-8`  
**说明：** 指定自动生成 UUID 时使用的版本号

**支持的版本：**
- `1` - 基于时间戳和 MAC 地址
- `4` - 完全随机（大多数场景用）
- `6` - 时间戳优化版（推荐用于 B-tree 索引）
- `7` - 时间戳 + 随机（分布式推荐）

### 2. DefaultCascadeMaxDepth
**类型：** `int`  
**默认值：** `6`  
**范围：** `1-6`  
**说明：** 级联查询的最大递归深度，防止无限递归

## 使用方法

### 基础用法
```go
db, err := simpleDBDriver.New.SimpleDB("mydb", "users")
defer db.Close()

// 查看当前配置
config := db.GetConfig()
fmt.Printf("UUID Version: %d, Cascade Depth: %d\n", 
    config.DefaultUUIDVersion, config.DefaultCascadeMaxDepth)
```

### 修改配置
```go
// 设置 UUID 版本为 7，级联深度为 4
db.SetConfig(simpleDBDriver.DatabaseConfig{
    DefaultUUIDVersion:     7,
    DefaultCascadeMaxDepth: 4,
})
```

### 在 Schema 中指定 UUID 版本

#### 方式一：Type 中指定版本（推荐）
```go
err := db.Configure(simpleDBDriver.TableSchema{
    Columns: []simpleDBDriver.Column{
        {
            Name:          "id",
            Type:          "uuid:v7",  // 使用冒号指定版本
            PrimaryKey:    true,
            AutoIncrement: true,
        },
    },
})
```

#### 方式二：使用斜杠（也支持）
```go
Type: "uuid/v6",  // 等同于 "uuid:v6"
```

#### 方式三：使用默认版本
```go
Type: "uuid",  // 使用 config 中的 DefaultUUIDVersion
```

## 版本选择建议

| 场景 | 推荐版本 | 说明 |
|------|---------|------|
| 数据库主键（单机） | v6 | 时间有序，索引友好 |
| 数据库主键（分布式） | v7 | 时间有序+随机，全局唯一 |
| 完全随机需求 | v4 | 传统 UUID |
| 遗留系统兼容 | v1 | 基于时间戳和 MAC |

## 完整示例

```go
package main

import (
    "fmt"
    "github.com/aid297/aid/simpleDB/simpleDBDriver"
)

func main() {
    db, _ := simpleDBDriver.New.SimpleDB("mydb", "users")
    defer db.Close()

    // 设置配置
    db.SetConfig(simpleDBDriver.DatabaseConfig{
        DefaultUUIDVersion:     7,
        DefaultCascadeMaxDepth: 4,
    })

    // 定义表结构
    db.Configure(simpleDBDriver.TableSchema{
        Columns: []simpleDBDriver.Column{
            {
                Name:          "id",
                Type:          "uuid:v7",
                PrimaryKey:    true,
                AutoIncrement: true,
            },
            {
                Name:     "name",
                Type:     "string",
                Required: true,
            },
        },
    })

    // 插入数据（UUID 自动生成为 v7）
    row, _ := db.InsertRow(simpleDBDriver.Row{
        "name": "Alice",
    })
    fmt.Println(row["id"]) // 输出 UUID v7 格式
}
```

## 验证配置

所有配置值都会在 `SetConfig` 时进行验证：
- UUID 版本必须在 1-8 范围内
- 级联深度必须在 1-6 范围内
- 无效值会被忽略，保留原有配置

```go
db.SetConfig(simpleDBDriver.DatabaseConfig{
    DefaultUUIDVersion: 10,  // 无效，会被忽略
})

config := db.GetConfig()
// config.DefaultUUIDVersion 仍然是前面的值
```

## API 参考

### SetConfig
```go
func (db *SimpleDB) SetConfig(config DatabaseConfig)
```
设置数据库配置（会验证参数有效性）

### GetConfig
```go
func (db *SimpleDB) GetConfig() DatabaseConfig
```
获取当前数据库配置

### 内部获取器
```go
func (db *SimpleDB) getDefaultUUIDVersion() int       // 获取有效的 UUID 版本
func (db *SimpleDB) getDefaultCascadeMaxDepth() int   // 获取有效的级联深度
```
