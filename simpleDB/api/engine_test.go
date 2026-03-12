package api

import (
	"strings"
	"testing"
)

func TestEngine_Execute_Smoke(t *testing.T) {
	eng := NewEngine(t.TempDir(), BackendDriver)

	if _, err := eng.Execute("CREATE TABLE users (id int PRIMARY KEY AUTO_INCREMENT, username string UNIQUE REQUIRED, age int DEFAULT 0)"); err != nil {
		t.Fatalf("create: %v", err)
	}

	if _, err := eng.Execute("INSERT INTO users (username, age) VALUES ('alice', 20)"); err != nil {
		t.Fatalf("insert: %v", err)
	}

	res, err := eng.Execute("SELECT * FROM users WHERE age>=20")
	if err != nil {
		t.Fatalf("select: %v", err)
	}
	if len(res.Rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(res.Rows))
	}
}

func TestEngine_Execute_BatchInsertAndUpdate(t *testing.T) {
	eng := NewEngine(t.TempDir(), BackendDriver)

	if _, err := eng.Execute("CREATE TABLE users (id int PRIMARY KEY AUTO_INCREMENT, username string UNIQUE REQUIRED, age int DEFAULT 0)"); err != nil {
		t.Fatalf("create: %v", err)
	}

	insertRes, err := eng.Execute("INSERT INTO users (username, age) VALUES ('alice', 20), ('bob', 21)")
	if err != nil {
		t.Fatalf("batch insert: %v", err)
	}
	if insertRes.Affected != 2 {
		t.Fatalf("insert affected = %d, want 2", insertRes.Affected)
	}
	if len(insertRes.InsertedRows) != 2 {
		t.Fatalf("inserted rows = %d, want 2", len(insertRes.InsertedRows))
	}

	updateRes, err := eng.Execute("UPDATE users SET age=30 WHERE id IN (1,2)")
	if err != nil {
		t.Fatalf("batch update: %v", err)
	}
	if updateRes.Affected != 2 {
		t.Fatalf("update affected = %d, want 2", updateRes.Affected)
	}
	if len(updateRes.UpdatedRows) != 2 {
		t.Fatalf("updated rows = %d, want 2", len(updateRes.UpdatedRows))
	}
}

func TestEngine_Execute_DeleteWithoutWhereRejected(t *testing.T) {
	eng := NewEngine(t.TempDir(), BackendDriver)

	if _, err := eng.Execute("CREATE TABLE users (id int PRIMARY KEY AUTO_INCREMENT, username string UNIQUE REQUIRED)"); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := eng.Execute("INSERT INTO users (username) VALUES ('alice')"); err != nil {
		t.Fatalf("insert: %v", err)
	}

	_, err := eng.Execute("DELETE FROM users")
	if err == nil {
		t.Fatal("expected error when delete has no where condition")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "where") {
		t.Fatalf("unexpected error: %v", err)
	}
}
