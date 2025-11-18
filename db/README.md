# Simple MySQL-like Database

一个使用 Go 开发的简单关系型数据库，支持基本的 SQL 语句。

## 功能特性

- 支持创建表 (CREATE TABLE)
- 支持插入数据 (INSERT INTO)
- 支持查询数据 (SELECT)
- 支持多种数据类型：INT, VARCHAR, TEXT, FLOAT, BOOL
- 数据持久化存储到文件
- 表结构保存和加载

## 支持的 SQL 语句

### 创建表
```sql
CREATE TABLE users (id INT, name VARCHAR, age INT)
```

### 插入数据
```sql
INSERT INTO users VALUES (1, 'Alice', 25)
INSERT INTO users VALUES (2, 'Bob', 30)
```

### 查询数据
```sql
SELECT * FROM users
```

## 交互式命令

除了 SQL 语句外，还支持以下命令：

- `LIST TABLES` - 列出所有表
- `DESC table_name` - 显示表结构
- `HELP` - 显示帮助信息
- `EXIT` - 退出数据库

## 使用方法

### 运行交互式模式
```bash
go run .
```

### 运行测试
```bash
go test -v
```

### 运行演示
```bash
chmod +x demo.sh
./demo.sh
```

## 项目结构

- `main.go` - 主程序，包含数据库核心逻辑
- `types.go` - 定义数据类型和表结构
- `parser.go` - SQL 语句解析器
- `storage.go` - 存储引擎，负责数据持久化
- `main_test.go` - 测试用例
- `demo.sh` - 演示脚本

## 存储格式

数据以两种文件格式存储：
- `table_name.schema.json` - 存储表结构
- `table_name.data` - 存储表数据（二进制格式）

## 局限性

当前版本有一些局限性：
- 不支持 WHERE 条件查询
- 不支持 UPDATE 和 DELETE 语句
- 不支持事务
- 不支持索引
- 不支持多表连接

这些功能可以在后续版本中逐步添加。