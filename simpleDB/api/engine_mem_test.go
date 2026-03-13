package api

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEngine_MemoryEngine(t *testing.T) {
	dbDir := t.TempDir()
	eng := NewEngine(dbDir, BackendDriver)

	// 1. 创建内存表
	if _, err := eng.Execute("CREATE TABLE mem_table (id int PRIMARY KEY, name string) WITH (engine=mem)"); err != nil {
		t.Fatalf("failed to create mem table: %v", err)
	}

	// 2. 插入数据
	if _, err := eng.Execute("INSERT INTO mem_table (id, name) VALUES (1, 'memory')"); err != nil {
		t.Fatalf("failed to insert into mem table: %v", err)
	}

	// 3. 验证数据存在
	res, err := eng.Execute("SELECT * FROM mem_table WHERE id=1")
	if err != nil {
		t.Fatalf("failed to select from mem table: %v", err)
	}
	if len(res.Rows) != 1 || res.Rows[0]["name"] != "memory" {
		t.Fatalf("expected 1 row with name 'memory', got %v", res.Rows)
	}

	// 4. 验证磁盘上没有 .tbl 文件
	tblPath := filepath.Join(dbDir, "mem_table", "mem_table.tbl")
	if _, err := os.Stat(tblPath); !os.IsNotExist(err) {
		t.Errorf("expected .tbl file to NOT exist for mem engine, but it does")
	}

	// 5. 验证禁止修改引擎
	_, err = eng.Execute("CREATE TABLE mem_table (id int PRIMARY KEY, name string) WITH (engine=disk)")
	if err == nil {
		t.Fatal("expected error when trying to change engine from mem to disk via CREATE TABLE (SyncSchema), but got nil")
	}
}

func TestEngine_DiskEngineDefault(t *testing.T) {
	dbDir := t.TempDir()
	eng := NewEngine(dbDir, BackendDriver)
	defer eng.Close()

	// 1. 创建磁盘表 (显式指定，因为现在默认是 mem)
	if _, err := eng.Execute("CREATE TABLE disk_table (id int PRIMARY KEY, name string) WITH (engine=disk)"); err != nil {
		t.Fatalf("failed to create disk table: %v", err)
	}

	db, err := eng.open("disk_table")
	if err != nil {
		t.Fatalf("failed to open table: %v", err)
	}
	tblPath := db.GetPath()

	// 2. 验证磁盘上有 .tbl 文件
	if _, err := os.Stat(tblPath); err != nil {
		t.Errorf("expected .tbl file to exist for disk engine, but got error: %v", err)
	}
}
