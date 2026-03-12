package api

import "testing"

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
