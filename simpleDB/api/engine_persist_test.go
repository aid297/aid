package api

import (
	"os"
	"testing"
	"time"
)

func TestEngine_MemoryPersist_Window(t *testing.T) {
	dbDir := t.TempDir()
	eng := NewEngine(dbDir, BackendDriver)
	defer eng.Close()

	// 1. 创建内存表，开启落盘
	// 设置较短的窗口期以便测试
	if _, err := eng.Execute("CREATE TABLE persist_table (id int PRIMARY KEY, name string) WITH (engine=mem, disk=true, windowSeconds=1, windowBytes=10mb, threshold=100mb)"); err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	// 2. 获取表路径
	db, err := eng.open("persist_table")
	if err != nil {
		t.Fatalf("failed to open table: %v", err)
	}
	tblPath := db.GetPath()

	if _, err := eng.Execute("INSERT INTO persist_table (id, name) VALUES (1, 'p1')"); err != nil {
		t.Fatalf("failed to insert: %v", err)
	}

	// 此时 .tbl 文件应该已经存在，因为 CreateTable 强制同步了 Schema
	if _, err := os.Stat(tblPath); err != nil {
		t.Errorf("expected .tbl file to exist after CreateTable, but got error: %v", err)
	}

	// 我们想要测试的是数据是否落盘。我们可以通过检查文件大小来验证。
	info, _ := os.Stat(tblPath)
	sizeBefore := info.Size()

	// 等待超过 1 秒
	t.Log("Waiting for window persistence (1s)...")
	time.Sleep(2 * time.Second)

	// 再次触发一个操作以触发 appendRecord 中的 check
	if _, err := eng.Execute("INSERT INTO persist_table (id, name) VALUES (2, 'p2')"); err != nil {
		t.Fatalf("failed to insert: %v", err)
	}

	// 此时文件大小应该增加了
	infoAfter, _ := os.Stat(tblPath)
	if infoAfter.Size() <= sizeBefore {
		t.Errorf("expected file size to increase after window persistence, but it didn't (%d -> %d)", sizeBefore, infoAfter.Size())
	}

	// 此时应该已经落盘
	if _, err := os.Stat(tblPath); err != nil {
		t.Errorf("expected .tbl file to exist after window timeout, but got error: %v", err)
	}

	// 验证数据还在内存中（窗口期落盘不清空内存）
	res, _ := eng.Execute("SELECT * FROM persist_table WHERE id=1")
	if len(res.Rows) != 1 {
		t.Errorf("expected row 1 to be in memory, but it's gone")
	}
}

func TestEngine_MemoryPersist_Threshold(t *testing.T) {
	dbDir := t.TempDir()
	eng := NewEngine(dbDir, BackendDriver)
	defer eng.Close()

	// 1. 创建内存表，开启落盘
	if _, err := eng.Execute("CREATE TABLE threshold_table (id int PRIMARY KEY, name string) WITH (engine=mem, disk=true, threshold=10, windowSeconds=100, windowBytes=100mb)"); err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	// 2. 获取表路径
	db, err := eng.open("threshold_table")
	if err != nil {
		t.Fatalf("failed to open table: %v", err)
	}
	tblPath := db.GetPath()

	// 3. 插入数据触发阈值
	if _, err := eng.Execute("INSERT INTO threshold_table (id, name) VALUES (1, 'large_data_to_trigger_threshold')"); err != nil {
		t.Fatalf("failed to insert: %v", err)
	}

	// 此时应该已经落盘且清空了内存
	if _, err := os.Stat(tblPath); err != nil {
		t.Errorf("expected .tbl file to exist after threshold trigger, but got error: %v", err)
	}

	// 验证内存未命中时可以回读磁盘
	res, _ := eng.Execute("SELECT * FROM threshold_table WHERE id=1")
	if len(res.Rows) != 1 || res.Rows[0]["name"] != "large_data_to_trigger_threshold" {
		t.Errorf("expected to read row from disk, got %v", res.Rows)
	}
}
