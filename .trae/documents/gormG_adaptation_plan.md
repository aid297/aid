# gormPools 适配 gorm.G（v1.30.0+ 泛型 API）增量方案

## 1. 背景与现状

### 1.1 当前实现

`gormPools` 包提供了基于传统 `*gorm.DB` 的模型绑定与查询辅助：

| 现有 API | 位置 | 作用 |
|---|---|---|
| `ToModel[model Modeler](db, attrs...) *gorm.DB` | [model.attr.go](file:///Users/yujizhou/development/projects/go/aid/gormPools/model.attr.go#L39-L48) | 绑定模型 + 应用 `ModelAttr`，返回 `*gorm.DB` |
| `ToFinder[model Modeler](db, attrs...) Finder` | [model.attr.go](file:///Users/yujizhou/development/projects/go/aid/gormPools/model.attr.go#L50-L52) | `ToModel` 后包装为 `Finder` |
| `ModelAttr` 接口 + `Table/Joins/Preload/PreloadMap/Select/Distinct/Scopes` | [model.attr.go](file:///Users/yujizhou/development/projects/go/aid/gormPools/model.attr.go#L17-L37) | 通过 `Register(model, *gorm.DB) *gorm.DB` 链式应用属性 |
| `Finder` 接口 + `FinderImpl` | [db_finder.go](file:///Users/yujizhou/development/projects/go/aid/gormPools/db_finder.go#L18-L55) | `When*` 条件、分页、排序、预加载等便捷方法 |

当前 `go.mod` 使用 `gorm.io/gorm v1.25.12`，尚不具备 `gorm.G`。

### 1.2 gorm v1.30.0 泛型 API 要点

`gorm.G[T](db, opts...)` 返回 `gorm.Interface[T]`，形成接口嵌套链：

```
Interface[T]
  └─ CreateInterface[T]   (Table, Create, CreateInBatches)
       └─ ChainInterface[T]  (Where, Not, Or, Limit, Offset, Joins, Preload, Select, Omit, Distinct, Group, Having, Order, Scopes, Delete, Update, Updates, Count)
            └─ ExecInterface[T]  (First, Last, Take, Find, FindInBatches, Scan, Row, Rows)
```

**关键差异与约束**（相对传统 `*gorm.DB`）：

1. **类型收窄**：链式方法返回更窄的接口。`Where/Select/...` 返回 `ChainInterface[T]`，不再具备 `Create` 能力；`Table` 仅存在于 `CreateInterface[T]`/`Interface[T]`，**不在 `ChainInterface[T]` 上**。
2. **Joins 签名变更**：`Joins(query clause.JoinTarget, on func(db JoinBuilder, joinTable, curTable clause.Table) error)`，不再接受字符串 join SQL。`clause.JoinTarget` 由 `clause.Has(name)` / `clause.LeftJoin.Association(name)` 等构造。
3. **Preload 签名变更**：`Preload(association string, query func(db PreloadBuilder) error)`，条件通过 `PreloadBuilder`（Where/Select/Limit/Order/LimitPerRecord）传入，不再接受 `...any`。
4. **Scopes 签名变更**：`Scopes(scopes ...func(db *Statement))`，回调参数为 `*gorm.Statement` 而非 `*gorm.DB`。
5. **Finisher 方法**：`Find(ctx) ([]T, error)`、`First(ctx) (T, error)`、`Count(ctx, column) (int64, error)`，需传入 `context.Context`，直接返回错误。
6. **移除 API**：无 `Save`、`FirstOrCreate`。
7. **向后兼容**：传统 `*gorm.DB` API 完全保留，可混用。

## 2. 设计目标

- **只做增量**：不修改 `ToModel`、`ToFinder`、`ModelAttr`、`Finder` 及其实现；新增文件承载泛型适配代码。
- **类型安全**：复用 `gorm.G[T]` 的泛型能力，返回类型化结果。
- **API 风格一致**：新增 `ToModelG`/`ToFinderG`/`G*` 属性构造函数，命名与现有体系对齐。

## 3. 方案设计

### 3.1 新增文件

`gormPools/model_g.attr.go`（不改动 `model.attr.go`）。

### 3.2 泛型属性 `GenericAttr[T]`

由于 `Table` 不在 `ChainInterface[T]` 上，而所有其他链式方法均返回 `ChainInterface[T]`，为统一属性接口，选择以 `gorm.ChainInterface[T]` 为基线：

```go
// GenericAttr 泛型模型属性：作用于 gorm.ChainInterface[T]
type GenericAttr[T any] func(db gorm.ChainInterface[T]) gorm.ChainInterface[T]
```

`GTable` 通过 `Scopes` 修改 `Statement.TableExpr` 实现表别名（与 `db.Table("x as y")` 内部机制一致）：

```go
func GTable[T Modeler](alias string) GenericAttr[T] {
    var m T
    tableName := m.TableName()
    expr := str.APP.Buffer.JoinStringLimit(" ", tableName, "as", alias)
    return func(db gorm.ChainInterface[T]) gorm.ChainInterface[T] {
        return db.Scopes(func(stmt *gorm.Statement) {
            stmt.TableExpr = &clause.Expr{SQL: expr}
        })
    }
}
```

其余属性直接映射到 `ChainInterface[T]` 方法：

| 构造函数 | 对应链式调用 | 说明 |
|---|---|---|
| `GSelect[T](query, args...)` | `db.Select(query, args...)` | 同传统 |
| `GDistinct[T](args...)` | `db.Distinct(args...)` | 同传统 |
| `GJoins[T](jt clause.JoinTarget, on func(gorm.JoinBuilder, clause.Table, clause.Table) error)` | `db.Joins(jt, on)` | 新签名，用 `clause.Has/LeftJoin.Association` 构造 |
| `GPreload[T](assoc string, fn func(gorm.PreloadBuilder) error)` | `db.Preload(assoc, fn)` | 新签名，`fn=nil` 表示无额外条件 |
| `GPreloadMap[T](map[string]func(gorm.PreloadBuilder) error)` | 循环 `db.Preload` | 适配新签名 |
| `GScopes[T](...func(*gorm.Statement))` | `db.Scopes(...)` | 新签名，参数为 `*Statement` |

> 注：`Joins`/`Preload`/`Scopes` 因 gorm 泛型 API 签名变化，无法与传统 `ModelAttr` 的参数保持一致，属于 API 层面的必然差异。

### 3.3 `ToModelG[T]`

```go
func ToModelG[T Modeler](db *gorm.DB, attrs ...GenericAttr[T]) gorm.ChainInterface[T] {
    var chain gorm.ChainInterface[T] = gorm.G[T](db) // Interface[T] 满足 ChainInterface[T]
    for _, attr := range attrs {
        chain = attr(chain)
    }
    return chain
}
```

- 返回 `gorm.ChainInterface[T]`：覆盖 `Where/Order/Limit/...` 与 `Find/First/Count` 等查询能力；不提供 `Create`（gorm 泛型 API 设计上 `Create` 仅在未加条件前可用）。
- 如需 `Create`，调用方直接使用 `gorm.G[T](db).Create(ctx, &obj)`。

### 3.4 `FinderG[T]`（对应 `Finder` 的泛型版）

新增 `FinderG[T]` 结构体包装 `gorm.ChainInterface[T]`，提供与 `Finder` 对齐的便捷方法，finisher 直接返回 `(结果, error)`：

```go
type FinderG[T any] struct {
    chain          gorm.ChainInterface[T]
    total          int64
    needPagination bool
}

func ToFinderG[T Modeler](db *gorm.DB, attrs ...GenericAttr[T]) *FinderG[T] {
    return &FinderG[T]{chain: ToModelG[T](db, attrs...)}
}
```

方法清单（对齐 `Finder`，但错误处理遵循泛型 API 风格）：

| 方法 | 签名 | 说明 |
|---|---|---|
| `Where/Not/Or` | `(query, args...) *FinderG[T]` | 透传到 chain |
| `When` | `(cond bool, query, args...) *FinderG[T]` | 条件成立时 `Where` |
| `WhenIn/WhenNotIn` | `(cond bool, query, args any) *FinderG[T]` | 指针/非指针对齐现有逻辑 |
| `WhenBetween/WhenNotBetween` | `(cond bool, query, args...) *FinderG[T]` | |
| `WhenLike/WhenLikeLeft/WhenLikeRight` | `(cond bool, query, arg any) *FinderG[T]` | 复用 `EscapeLikeValue` |
| `WhenFunc` | `(cond bool, fn func(*gorm.Statement)) *FinderG[T]` | 用 `Scopes` 承载 |
| `TryPagination` | `(page, size int) *FinderG[T]` | 分页，需 Count |
| `TryOrder` | `(orders ...string) *FinderG[T]` | 复用 `validateOrder` |
| `TryPreload` | `(assocs ...string) *FinderG[T]` | 无额外条件的 Preload |
| `Find` | `(ctx) ([]T, error)` | 调用 `chain.Find(ctx)` |
| `First` | `(ctx) (T, error)` | 调用 `chain.First(ctx)` |
| `Count` | `(ctx) (int64, error)` | 调用 `chain.Count(ctx, "*")` |
| `GetTotal` | `() int64` | |

> `FinderG` 的 `WhenIn` 等方法不再需要 `reflect` 判断指针，因为泛型 API 的 `Where` 直接处理切片；但为与现有 `Finder` 行为一致，保留对指针参数的解包。

### 3.5 gorm 版本升级

需将 `go.mod` 中 `gorm.io/gorm` 从 `v1.25.12` 升级至 `v1.30.0+`。官方声明传统 API 完全向后兼容，但需验证：
- `clause` 包是否有不兼容变更（`clause.Expr`、`clause.Table` 等）
- `gorm.Statement` 字段 `TableExpr` 在 v1.30.0 是否仍存在
- 现有 `condition.go` 中 `gorm/schema.Parse` 用法是否兼容

## 4. 改动清单

| 文件 | 操作 | 内容 |
|---|---|---|
| `go.mod` / `go.sum` | 修改 | 升级 `gorm.io/gorm` 至 `v1.30.0+`（执行 `go get gorm.io/gorm@v1.30.0`） |
| `gormPools/model_g.attr.go` | 新增 | `GenericAttr[T]`、`GTable/GSelect/GDistinct/GJoins/GPreload/GPreloadMap/GScopes`、`ToModelG[T]` |
| `gormPools/finder_g.go` | 新增 | `FinderG[T]` 结构体及方法、`ToFinderG[T]` |

不修改：`model.attr.go`、`db_finder.go`、`interface.go`、`condition.go`。

## 5. 风险与注意事项

1. **`GTable` 依赖 `Statement.TableExpr` 内部实现**：若 gorm 未来版本调整 `Statement` 结构，需同步更新。当前 v1.30.0 中 `db.Table(name)` 内部即设置 `Statement.TableExpr = &clause.Expr{SQL: name}`，机制一致。
2. **`Joins`/`Preload` API 不兼容**：泛型版参数与传统版不同，无法做到签名一致，调用方需按新签名使用。
3. **类型收窄限制**：`ToModelG` 返回 `ChainInterface[T]`，无法继续调用 `Create`。如需创建，走 `gorm.G[T](db).Create(...)`。
4. **gorm 升级回归**：升级后需运行 `go build ./...` 与现有测试，确认传统 API 无破坏。
5. **`FinderG` 与 `Finder` 行为差异**：`Finder` 的 finisher 不返回 error（通过 `Error()` 获取），`FinderG` 的 finisher 直接返回 `(result, error)`，更贴合泛型 API 习惯。

## 6. 验证方式

1. `go build ./...` 编译通过。
2. 现有 `gormPools` 测试（`db_finder_test.go`）继续通过，证明传统 API 未破坏。
3. 新增 `model_g.attr` 编译期类型校验：`ToModelG[TestTable1](db, GTable[TestTable1]("t"), GSelect[TestTable1]("id,name")).Where("id=?", 1).Find(ctx)` 能正确推导 `[]TestTable1`。
