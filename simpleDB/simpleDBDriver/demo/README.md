# simpleDBDriver demo

这个 demo 展示了当前 `simpleDBDriver` 的主要能力：

- 原始 append-only KV：`Put` / `Update` / `Delete` / `Query` / `Compact`
- 结构化表：`Configure` / `InsertRow` / `FindByUnique` / `FindByIndex`
- 条件查询：索引优先，非索引字段自动回退过滤
- JSON 查询：`FindByConditionsJSON`
- 外键级联查询：`QueryCascadeJSON`
- 嵌套结果：主表、子表、孙表按 JSON 嵌套返回

## 运行方式

在仓库根目录执行：

```bash
go run ./simpleDB/simpleDBDriver/demo
```

## demo 结构

`main.go` 中按方法拆分：

- `demoKVCRUD()`：演示基础 KV 与压缩
- `demoSchemaAndQueries()`：演示 schema、唯一索引、普通索引、条件查询
- `demoJSONQuery()`：演示 JSON 查询返回
- `demoCascadeQuery()`：演示多表外键级联与嵌套 JSON

## 说明

- demo 会使用 `demo_runtime/` 目录保存临时数据
- 程序运行前后都会自动清理示例数据
- 级联查询示例中：
  - `users_rel_demo` 是主表
  - `orders_rel_demo` 通过 `user_id -> users_rel_demo.id` 关联用户
  - `order_items_rel_demo` 通过 `order_id -> orders_rel_demo.id` 关联订单
