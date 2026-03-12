//go:build never

package main

import (
	"encoding/json"
	"fmt"

	"github.com/aid297/aid/simpleDB/api"
)

func main() {
	engine := api.NewEngine("demo_api", api.BackendDriver)

	_, _ = engine.Execute("DROP TABLE orders")
	_, _ = engine.Execute("DROP TABLE users")

	stmt, err := engine.Parse("CREATE TABLE users (id int PRIMARY KEY AUTO_INCREMENT, username string UNIQUE REQUIRED, age int DEFAULT 0)")
	if err != nil {
		panic(err)
	}
	printJSON("parsed create-users AST", stmt)

	mustExecute(engine, "CREATE TABLE users (id int PRIMARY KEY AUTO_INCREMENT, username string UNIQUE REQUIRED, age int DEFAULT 0)")
	mustExecute(engine, "ALTER TABLE users ADD COLUMN email string DEFAULT '', ADD INDEX(email)")
	mustExecute(engine, "INSERT INTO users (username, age, email) VALUES ('alice', 20, 'alice@aid.dev')")
	mustExecute(engine, "INSERT INTO users (username, age, email) VALUES ('bob', 28, 'bob@aid.dev')")
	mustExecute(engine, "UPDATE users SET age=29 WHERE id=2")

	users := mustExecute(engine, "SELECT * FROM users WHERE age>=20")
	printJSON("users", users)

	mustExecute(engine, "CREATE TABLE orders (id int PRIMARY KEY AUTO_INCREMENT, user_id int, amount float DEFAULT 0)")
	mustExecute(engine, "ALTER TABLE orders ADD FOREIGN KEY(user_id) REFERENCES users(id) AS user NAME fk_orders_user")
	mustExecute(engine, "INSERT INTO orders (user_id, amount) VALUES (1, 19.9)")
	mustExecute(engine, "INSERT INTO orders (user_id, amount) VALUES (2, 88.8)")

	orders := mustExecute(engine, "SELECT * FROM orders WHERE amount>=10")
	printJSON("orders", orders)

	mustExecute(engine, "DELETE FROM orders WHERE id=1")
	leftOrders := mustExecute(engine, "SELECT * FROM orders")
	printJSON("orders after delete", leftOrders)

	mustExecute(engine, "TRUNCATE TABLE orders")
	emptyOrders := mustExecute(engine, "SELECT * FROM orders")
	printJSON("orders after truncate", emptyOrders)
}

func mustExecute(engine *api.Engine, sql string) api.ExecResult {
	result, err := engine.Execute(sql)
	if err != nil {
		panic(fmt.Sprintf("execute failed: %s -> %v", sql, err))
	}
	fmt.Printf("SQL: %s\n", sql)
	return result
}

func printJSON(title string, value any) {
	payload, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n== %s ==\n%s\n", title, string(payload))
}





























































}	fmt.Printf("\n== %s ==\n%s\n", title, string(payload))	}		panic(err)	if err != nil {	payload, err := json.MarshalIndent(value, "", "  ")func printJSON(title string, value any) {}	return result	fmt.Printf("SQL: %s\n", sql)	}		panic(fmt.Sprintf("execute failed: %s -> %v", sql, err))	if err != nil {	result, err := engine.Execute(sql)func mustExecute(engine *api.Engine, sql string) api.ExecResult {}	printJSON("orders after truncate", emptyOrders)	emptyOrders := mustExecute(engine, "SELECT * FROM orders")	mustExecute(engine, "TRUNCATE TABLE orders")	printJSON("orders after delete", leftOrders)	leftOrders := mustExecute(engine, "SELECT * FROM orders")	mustExecute(engine, "DELETE FROM orders WHERE id=1")		printJSON("orders", orders)	orders := mustExecute(engine, "SELECT * FROM orders WHERE amount>=10")	mustExecute(engine, "INSERT INTO orders (user_id, amount) VALUES (2, 88.8)")	mustExecute(engine, "INSERT INTO orders (user_id, amount) VALUES (1, 19.9)")	mustExecute(engine, "ALTER TABLE orders ADD FOREIGN KEY(user_id) REFERENCES users(id) AS user NAME fk_orders_user")	mustExecute(engine, "CREATE TABLE orders (id int PRIMARY KEY AUTO_INCREMENT, user_id int, amount float DEFAULT 0)")	printJSON("users", users)	users := mustExecute(engine, "SELECT * FROM users WHERE age>=20")	mustExecute(engine, "UPDATE users SET age=29 WHERE id=2")	mustExecute(engine, "INSERT INTO users (username, age, email) VALUES ('bob', 28, 'bob@aid.dev')")	mustExecute(engine, "INSERT INTO users (username, age, email) VALUES ('alice', 20, 'alice@aid.dev')")	mustExecute(engine, "ALTER TABLE users ADD COLUMN email string DEFAULT '', ADD INDEX(email)")	mustExecute(engine, "CREATE TABLE users (id int PRIMARY KEY AUTO_INCREMENT, username string UNIQUE REQUIRED, age int DEFAULT 0)")	printJSON("parsed create-users AST", stmt)	}		panic(err)	if err != nil {	stmt, err := engine.Parse("CREATE TABLE users (id int PRIMARY KEY AUTO_INCREMENT, username string UNIQUE REQUIRED, age int DEFAULT 0)")	_, _ = engine.Execute("DROP TABLE users")	_, _ = engine.Execute("DROP TABLE orders")	// 先清理旧表，保证 demo 可重复执行。	engine := api.NewEngine("demo_api", api.BackendDriver)func main() {)	"github.com/aid297/aid/simpleDB/api"	"fmt"	"encoding/json"import (