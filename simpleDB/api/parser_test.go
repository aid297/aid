//go:build never

package api
package sqlapi

import "testing"

func TestParse_CreateAlter(t *testing.T) {
	stmt, err := Parse("CREATE TABLE users (id int PRIMARY KEY AUTO_INCREMENT, username string UNIQUE REQUIRED, age int DEFAULT 0)")
	if err != nil {
		t.Fatalf("parse create: %v", err)
	}
	create, ok := stmt.(CreateTableStmt)


































}	}		t.Fatalf("parse select: %v", err)	if _, err := Parse("SELECT * FROM users WHERE age>=18"); err != nil {	}		t.Fatalf("parse delete: %v", err)	if _, err := Parse("DELETE FROM users WHERE age>=18 AND username!='x'"); err != nil {	}		t.Fatalf("parse update: %v", err)	if _, err := Parse("UPDATE users SET age=21 WHERE id=1"); err != nil {	}		t.Fatalf("parse insert: %v", err)	if _, err := Parse("INSERT INTO users (username, age) VALUES ('alice', 20)"); err != nil {func TestParse_DML(t *testing.T) {}	}		t.Fatalf("unexpected alter plan: %+v", alter.Plan)	if len(alter.Plan.AddColumns) != 1 || len(alter.Plan.AddUniques) != 1 || len(alter.Plan.AddIndexes) != 1 {	}		t.Fatalf("unexpected alter stmt: %#v", stmt)	if !ok {	alter, ok := stmt.(AlterTableStmt)	}		t.Fatalf("parse alter: %v", err)	if err != nil {	stmt, err = Parse("ALTER TABLE users ADD COLUMN email string DEFAULT '', ADD UNIQUE(email), ADD INDEX(email)")	}		t.Fatalf("expected 3 columns, got %d", len(create.Schema.Columns))	if len(create.Schema.Columns) != 3 {	}		t.Fatalf("unexpected create stmt: %#v", stmt)	if !ok || create.Table != "users" {