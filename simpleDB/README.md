# simpleDB HTTP SQL Demo

## 功能概览

`simpleDB` 现在提供：

- 用户认证：登录、注册、激活、停用、角色授权。
- SQL 执行接口：`POST /sql/execute`。
- 支持参数化 SQL：使用 `:name` 占位，`params` 传值。
- 支持批量写入/批量更新：
  - `INSERT INTO ... VALUES (...), (...)`
  - `UPDATE ... WHERE id IN (...)`
- 安全限制：`DELETE` **必须带 `WHERE`**，否则拒绝执行。

## 快速启动

在仓库根目录执行：

1. 初始化配置（可覆盖）

```bash
go run ./simpleDB init-config -force
```

2. 启动服务

```bash
go run ./simpleDB serve -config simpleDB/config.yaml
```

默认地址：`http://127.0.0.1:18080`

> 首次启动若 `initPassword` 为空，会自动生成并写回配置。

## SQL 执行接口

### 路由

默认：`POST /sql/execute`

可在配置中修改：

```yaml
transport:
  http:
    route:
      sql: /sql/execute
```

### 请求体

```json
{
  "sql": "SELECT * FROM users WHERE id IN :ids AND username = ?",
  "paramMap": {
    "ids": [1, 2, 3]
  },
  "paramList": ["alice"]
}
```

- `sql`：必填。
- `paramMap`：可选，`map` 结构，对应 `:name` 占位符。
- `paramList`：可选，切片结构，对应 `?` 占位符，按出现顺序绑定。
- `paramMap` 和 `paramList` 可以同时使用。
- 兼容旧字段 `params`，建议新调用统一改为 `paramMap`。
- 需要 `Authorization: Bearer <token>`。

### 返回体（示例字段）

```json
{
  "success": true,
  "result": {
    "statement": "update",
    "affected": 2,
    "updatedRows": [
      {"id": 1, "username": "alice", "age": 30},
      {"id": 2, "username": "bob", "age": 30}
    ]
  }
}
```

- 单条插入：`result.inserted`
- 批量插入：`result.insertedRows`
- 单条更新：`result.updated`
- 批量更新：`result.updatedRows`
- 查询：`result.rows`
- 删除：`result.affected`

## 完整 Demo 脚本

脚本：

 - [simpleDB/demo_http_sql.rest](simpleDB/demo_http_sql.rest)

作用：

1. 自动启动服务（后台）
2. 登录 `sdb`
3. 建表
4. 单条插入（参数化）
5. 批量插入
6. 批量更新（参数化 `IN`）
7. 查询结果
8. 验证无条件删除被拒绝
9. 条件删除并再次查询

## .rest 请求文件（REST Client）

如果你更喜欢在编辑器里点请求，使用：

- [simpleDB/demo_http_sql.rest](simpleDB/demo_http_sql.rest)

使用方式：

1. 安装 VS Code 的 REST Client 扩展。
2. 打开 [simpleDB/demo_http_sql.rest](simpleDB/demo_http_sql.rest)。
3. 启动服务后，从上到下点击 `Send Request`。
4. 文件里已包含登录、建表、单条/批量写入、批量更新、查询、删除限制校验全流程。
### 运行

```bash
chmod +x simpleDB/demo_http_sql.sh
./simpleDB/demo_http_sql.sh
```

可选环境变量：

- `BASE_URL`：默认 `http://127.0.0.1:18080`
- `CONFIG_PATH`：默认 `simpleDB/config.yaml`
- `START_SERVER`：默认 `1`（自动启动）；设为 `0` 时只发请求

例如：

```bash
BASE_URL=http://127.0.0.1:18080 START_SERVER=0 ./simpleDB/demo_http_sql.sh
```

## 常见错误

- `401 unauthorized`：未登录或 token 无效。
- `403 forbidden`：访问系统表但非 `super_admin`。
- `400 bad_request`：SQL 语法不支持、参数缺失，或 `DELETE` 无 `WHERE`。
