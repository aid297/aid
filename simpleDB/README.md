# simpleDB

`simpleDB` 是一个 Go 实现的轻量级数据库内核 + SQL 执行层 + HTTP/WebSocket 服务。你可以：

- 作为嵌入式库在 Go 程序中直接打开表并 CRUD（见 driver）。
- 作为独立服务启动，通过 HTTP 或 WebSocket 执行 SQL 与鉴权管理。

## 功能概览

- 存储引擎：`mem` / `disk`
  - `engine=mem`：内存引擎；可选开启落盘（批量落盘 + 阈值清内存 + 按需回灌）
  - `engine=disk`：磁盘引擎；每次写入追加到数据文件并 `fsync`
- SQL：DDL / DML / 查询
  - 参数化 SQL：`:name`（paramMap）与 `?`（paramList）混用
  - 批量写入/更新：`INSERT ... VALUES (...), (...)`、`UPDATE ... WHERE id IN (...)`
  - 安全限制：`DELETE` 必须带 `WHERE`
- Auth：登录/注册/刷新/注销、激活/停用用户、角色授权、权限授权
- Transport：HTTP + WebSocket
  - HTTP 支持 `Accept` 内容协商（json/xml/yaml/toml，默认 json）
  - HTTP 请求体按 `Content-Type` 解析（json/xml/yaml/toml，默认 json）
  - WebSocket 始终使用 JSON 消息体
- 系统字段：每张表自动包含 `_version`（UUID v7），写入/更新时自动生成且不可由用户覆盖/删除

## 快速开始（启动服务）

在仓库根目录执行：

1) 生成默认配置（可覆盖）

```bash
go run ./simpleDB init-config -force
```

2) 启动服务

```bash
go run ./simpleDB serve -config simpleDB/config.yaml
```

默认地址：`http://127.0.0.1:18080`

## “连接建立”有三种方式

### 方式 A：Go 代码直接打开表（嵌入式）

建议使用 driver 层（稳定 API + 统一错误模型 + 事务）：

- 文档与示例：[driver/README.md](./driver/README.md)
- 典型用法：

```go
db, err := driver.New.DB("demo", "users")
defer db.Close()

err = db.Configure(driver.TableSchema{
  Columns: []driver.Column{
    {Name: "id", Type: "uuid:v7", PrimaryKey: true, AutoIncrement: true},
    {Name: "username", Type: "string", Unique: true, Required: true},
  },
})
```

### 方式 B：HTTP（REST）

先登录获取 token，再用 `Authorization: Bearer <token>` 调用需要鉴权的接口：

- `POST /auth/login` → 返回 `token.accessToken`
- `POST /sql/execute` / `GET /me` / `POST /sql/grant` 等 → 需要 Bearer token

### 方式 C：WebSocket（JSON 消息）

1) 建立 WS 连接：`GET /ws`（路由可在配置里改）
2) 首条发送登录消息（route 固定 `/login`，payload 携带 token）：

```json
{"route":"/login","token":"<accessToken>"}
```

3) 之后才能发送 SQL 执行消息（route 为配置的 SQL 路由，默认 `/sql/execute`）：

```json
{
  "route": "/sql/execute",
  "sql": "SELECT * FROM users WHERE id IN :ids AND username = ?",
  "paramMap": {"ids":[1,2,3]},
  "paramList": ["alice"]
}
```

服务端会按 `heartbeatInterval` 定时发送 Ping 帧；大多数 WS 客户端会自动回 Pong。

## 存储模式：disk 与 mem（含 mem 落盘）

存储模式是“表级别”的：每张表在 schema 里指定 `engine` 与 `disk`。

### disk 模式

- `engine=disk`
- 每次写入直接追加到 `<database>/<table>.tbl` 并同步落盘
- 适合对单次写入持久性要求更高的场景

### mem 模式（可选落盘）

- `engine=mem`
- `disk=false`（默认纯内存）：不写文件
- `disk=true`（mem+disk）：启用批量落盘机制
  - `windowSeconds/windowBytes`：满足任一则把内存日志批量落盘
  - `threshold`：超过则落盘并清空内存态；之后读/写前会“按需从磁盘回灌”

当 `engine=mem` 且 `disk=true` 时，如果你没有在建表 schema/SQL 的 `WITH (...)` 里显式填写 `windowSeconds/windowBytes/threshold`，会从配置文件 `engine.persistence` 里继承默认值。

### 通过 SQL 选择引擎

`CREATE TABLE` 支持 `WITH (...)` 参数：

```sql
CREATE TABLE users (
  id uuid:v7 primary key auto_increment,
  username string unique required,
  age int required
) WITH (engine=mem, disk=true, windowSeconds=10, windowBytes=10mb, threshold=100mb)
```

也可以显式选择磁盘引擎：

```sql
CREATE TABLE users (
  id uuid:v7 primary key auto_increment,
  username string unique required
) WITH (engine=disk)
```

## CLI 使用方式

CLI 入口：`go run ./simpleDB <command> [flags]`（无 command 默认等价于 `serve`）。

- `serve -config <path>`：读取配置并启动 HTTP/WS 服务
- `init-config -config <path> [-force]`：写入默认配置
- `print-config -config <path>`：打印“加载后 + 自动补默认值”的最终配置
- `gen-key [-config <path>] [-algo <aes>] [-len <bytes>] [-format <hex|base64>]`：生成加密密钥

示例：

```bash
go run ./simpleDB init-config -force
go run ./simpleDB print-config -config simpleDB/config.yaml
go run ./simpleDB serve -config simpleDB/config.yaml
go run ./simpleDB gen-key -algo aes -len 32 -format base64
```

## 配置文件说明（simpleDB/config.yaml）

默认配置样例见：[config.yaml](./config.yaml)

### database

- `database.path`：数据库根路径（目录名）。每张表会在该目录下创建 `<table>.tbl` / `<table>.lock`。

### engine

- `engine.persistence.windowSeconds`：mem+disk 时的落盘窗口（秒）
- `engine.persistence.windowBytes`：mem+disk 时的落盘窗口（字节数，支持 `kb/mb/gb`）
- `engine.persistence.threshold`：超过阈值则“落盘并清空内存态”
- `engine.security.compressAlgorithm`：落盘写入时压缩（例如 `gzip`；空表示关闭）
- `engine.security.encryptAlgorithm`：落盘写入时加密（例如 `aes`；空表示关闭）
- `engine.security.encryptKey`：加密密钥（例如 AES-256 的 32 字节 base64）

### transport.http

- `enabled`：是否开启 HTTP（当前配置校验要求必须为 true）
- `address`：监听地址（例如 `:18080`）
- `ginMode`：`release` / `debug`
- `route.*`：各接口路由（可改）
- `limit.*`：基于 token 的简单限流（可配置不需要 token 的路径列表）
- `initPassword`：用于调用 `POST /auth/init-sdb-password` 的“初始化口令”
- `tokenTTL`：token 有效期
- `tokenSecret`：token 签名密钥（空则使用默认值：`simpledb.transport.<database.path>.secret`）
- `enableAdminRoute` / `enableReportRoute`：是否注册 `/admin` / `/reports`
- `sqlAllowedOps`：SQL 操作白名单（空/nil 表示不限制）

### transport.websocket

- `enabled`：是否开启 WS
- `route`：WS 路由（例如 `/ws`）
- `heartbeatInterval`：服务端 ping 周期
- `writeTimeout` / `readTimeout`：读写超时
- `executionTimeout`：单条 SQL 的最大执行时间

## HTTP 接口使用

### 内容类型（重要）

- 请求体：按 `Content-Type` 解析 `json/xml/yaml/toml`，默认 `application/json`
- 响应体：按 `Accept` 返回 `json/xml/yaml/toml`，默认 `application/json`

### 接口列表（默认路由）

Auth：

- `POST /auth/login`
- `POST /auth/register`
- `POST /auth/refresh`
- `POST /auth/logout`
- `POST /auth/activate`（需要 `super_admin`）
- `POST /auth/deactivate`（需要 `super_admin`）
- `POST /auth/assign-roles`（需要 `super_admin`）
- `POST /auth/assign-role-permissions`（需要 `super_admin`）
- `POST /auth/init-sdb-password`（需要 `initPassword`）

SQL：

- `POST /sql/execute`（需要 Bearer token）
- `POST /sql/grant`（需要 Bearer token）
- `POST /sql/revoke`（需要 Bearer token）

Service：

- `GET /health`
- `GET /me`（需要 Bearer token）
- `GET /admin`（可选，且需要 `super_admin`）
- `GET /reports`（可选，且需要 permission `report.read`）

### 登录（获取 token）

```http
POST /auth/login
Content-Type: application/json

{"username":"sdb","password":"simpleDB"}
```

### 初始化/重置 sdb 密码（需要 initPassword）

该接口会把系统管理员 `sdb` 的密码重置为默认值 `simpleDB`，并确保其为 active + admin。

```http
POST /auth/init-sdb-password
Content-Type: application/json

{"password":"<initPassword from config.yaml>"}
```

### 执行 SQL（支持 paramMap + paramList）

```http
POST /sql/execute
Authorization: Bearer <accessToken>
Content-Type: application/json

{
  "sql": "SELECT * FROM users WHERE id IN :ids AND username = ?",
  "paramMap": {"ids":[1,2,3]},
  "paramList": ["alice"]
}
```

返回体（示例字段）：

```json
{
  "success": true,
  "result": {
    "statement": "update",
    "affected": 2,
    "updatedRows": [
      {"id": 1, "username": "alice", "age": 30}
    ]
  }
}
```

## WebSocket 使用

1) 连接：`ws://127.0.0.1:18080/ws`
2) 发送登录消息（必须先 login）：

```json
{"route":"/login","token":"<accessToken>"}
```

3) 执行 SQL（route 必须等于配置的 SQL 路由，默认 `/sql/execute`）：

```json
{"route":"/sql/execute","sql":"SELECT 1","paramMap":{},"paramList":[]}
```

WebSocket 返回体始终是 JSON。

## Demo 与请求文件

- 全流程脚本：`simpleDB/demo_http_sql.sh`
- REST Client 文件：
  - `simpleDB/demo_http_sql.rest`（VS Code REST Client）
  - `simpleDB/demo_for_vscode.rest` / `simpleDB/demo_for_idea.rest`

## 常见问题（FAQ）

### JSON 大整数精度丢失

当你把 JSON 数字解到 `any` / `map[string]any` 时，默认会变成 `float64`，超过 `2^53-1` 的整数会丢精度。

本项目中：

- HTTP/WS 的 SQL 参数 `paramMap/paramList` 已启用 `UseNumber`，可以安全传输超大整数并保持精度。
- 其它接口如果需要传输超大整数，建议客户端用字符串表示（例如 `"id":"9223372036854775807"`），或确保数值不超过 `2^53-1`。

### 常见错误码

- `401 unauthorized`：未登录或 token 无效/过期
- `403 forbidden`：访问系统表但非 `super_admin`，或缺少 required permission/role
- `400 bad_request`：SQL 语法不支持、参数缺失，或 `DELETE` 无 `WHERE`
