//go:build never

package sqlapi
package sqlapi

import "testing"

func TestEngine_Execute_Driver(t *testing.T) {
	dbName := t.TempDir()
	eng := NewEngine(dbName, BackendDriver)

	_, err := eng.Execute("CREATE TABLE users (id int PRIMARY KEY AUTO_INCREMENT, username string UNIQUE REQUIRED, age int DEFAULT 0)")
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if _, err = eng.Execute("INSERT INTO users (username, age) VALUES ('alice', 20)"); err != nil {
		t.Fatalf("insert: %v", err)
	}
	if _, err = eng.Execute("UPDATE users SET age=21 WHERE id=1"); err != nil {




























}	}		t.Fatalf("alter add fk: %v", err)	if _, err := eng.Execute("ALTER TABLE orders ADD FOREIGN KEY(user_id) REFERENCES users(id) AS user NAME fk_orders_user"); err != nil {	}		t.Fatalf("create orders: %v", err)	if _, err := eng.Execute("CREATE TABLE orders (id int PRIMARY KEY AUTO_INCREMENT, user_id int, amount float DEFAULT 0)"); err != nil {	}		t.Fatalf("create users: %v", err)	if _, err := eng.Execute("CREATE TABLE users (id int PRIMARY KEY AUTO_INCREMENT, username string UNIQUE REQUIRED)"); err != nil {	eng := NewEngine(dbName, BackendKernal)	dbName := t.TempDir()func TestEngine_Execute_Kernal_WithForeignKeyDDL(t *testing.T) {}	}		t.Fatalf("delete: %v", err)	if _, err = eng.Execute("DELETE FROM users WHERE id=1"); err != nil {	}		t.Fatalf("expected 1 row, got %d", len(res.Rows))	if len(res.Rows) != 1 {	}		t.Fatalf("select: %v", err)	if err != nil {	res, err := eng.Execute("SELECT * FROM users WHERE age>=21")	}		t.Fatalf("update: %v", err)