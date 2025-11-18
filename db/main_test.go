package main

import (
	"os"
	"testing"
)

// 测试基本的数据库功能
func TestDatabase(t *testing.T) {
	// 创建临时目录用于测试
	testDataDir := "./test_data"
	defer os.RemoveAll(testDataDir)
	
	// 创建数据库实例
	db, err := NewDatabase(testDataDir)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	
	// 测试创建表
	err = db.Execute("CREATE TABLE users (id INT, name VARCHAR, age INT)")
	if err != nil {
		t.Errorf("Failed to create table: %v", err)
	}
	
	// 验证表是否存在
	if _, exists := db.parser.Tables["users"]; !exists {
		t.Error("Table 'users' was not created")
	}
	
	// 测试插入数据
	err = db.Execute("INSERT INTO users VALUES (1, 'Alice', 25)")
	if err != nil {
		t.Errorf("Failed to insert data: %v", err)
	}
	
	err = db.Execute("INSERT INTO users VALUES (2, 'Bob', 30)")
	if err != nil {
		t.Errorf("Failed to insert data: %v", err)
	}
	
	// 验证数据是否插入成功
	if len(db.parser.Tables["users"].Rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(db.parser.Tables["users"].Rows))
	}
	
	// 验证插入的数据
	if db.parser.Tables["users"].Rows[0].Values[0] != 1 ||
	   db.parser.Tables["users"].Rows[0].Values[1] != "Alice" ||
	   db.parser.Tables["users"].Rows[0].Values[2] != 25 {
		t.Error("Inserted data doesn't match expected values")
	}
	
	if db.parser.Tables["users"].Rows[1].Values[0] != 2 ||
	   db.parser.Tables["users"].Rows[1].Values[1] != "Bob" ||
	   db.parser.Tables["users"].Rows[1].Values[2] != 30 {
		t.Error("Inserted data doesn't match expected values")
	}
	
	// 测试重新加载数据
	newDb, err := NewDatabase(testDataDir)
	if err != nil {
		t.Fatalf("Failed to create new database instance: %v", err)
	}
	
	// 验证表是否重新加载
	table, exists := newDb.parser.Tables["users"]
	if !exists {
		t.Error("Table 'users' was not loaded from disk")
		return
	}
	
	// 验证数据是否重新加载
	if table == nil {
		t.Error("Table 'users' is nil after loading")
		return
	}
	
	if len(table.Rows) != 2 {
		t.Errorf("Expected 2 rows after reload, got %d", len(table.Rows))
	}
	
	// 验证重新加载的数据
	if len(table.Rows) > 0 && (table.Rows[0].Values[0] != 1 ||
	   table.Rows[0].Values[1] != "Alice" ||
	   table.Rows[0].Values[2] != 25) {
		t.Error("Reloaded data doesn't match expected values")
	}
}

// 测试不同的数据类型
func TestDataTypes(t *testing.T) {
	// 创建临时目录用于测试
	testDataDir := "./test_data_types"
	defer os.RemoveAll(testDataDir)
	
	// 创建数据库实例
	db, err := NewDatabase(testDataDir)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	
	// 创建包含各种数据类型的表
	err = db.Execute("CREATE TABLE test_types (id INT, name VARCHAR, score FLOAT, active BOOL)")
	if err != nil {
		t.Errorf("Failed to create table: %v", err)
	}
	
	// 插入各种类型的数据
	err = db.Execute("INSERT INTO test_types VALUES (1, 'Alice', 95.5, true)")
	if err != nil {
		t.Errorf("Failed to insert data: %v", err)
	}
	
	err = db.Execute("INSERT INTO test_types VALUES (2, 'Bob', 87.3, false)")
	if err != nil {
		t.Errorf("Failed to insert data: %v", err)
	}
	
	// 验证数据类型
	row := db.parser.Tables["test_types"].Rows[0]
	
	if _, ok := row.Values[0].(int); !ok {
		t.Error("Expected id to be int")
	}
	
	if _, ok := row.Values[1].(string); !ok {
		t.Error("Expected name to be string")
	}
	
	if _, ok := row.Values[2].(float64); !ok {
		t.Error("Expected score to be float64")
	}
	
	if _, ok := row.Values[3].(bool); !ok {
		t.Error("Expected active to be bool")
	}
}

// 测试错误处理
func TestErrorHandling(t *testing.T) {
	// 创建临时目录用于测试
	testDataDir := "./test_error_handling"
	defer os.RemoveAll(testDataDir)
	
	// 创建数据库实例
	db, err := NewDatabase(testDataDir)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	
	// 测试插入不存在的表
	err = db.Execute("INSERT INTO non_existent VALUES (1, 'test')")
	if err == nil {
		t.Error("Expected error when inserting into non-existent table")
	}
	
	// 测试列数不匹配
	err = db.Execute("CREATE TABLE test (id INT, name VARCHAR)")
	if err != nil {
		t.Errorf("Failed to create table: %v", err)
	}
	
	err = db.Execute("INSERT INTO test VALUES (1)")
	if err == nil {
		t.Error("Expected error when column count doesn't match")
	}
	
	// 测试无效的SQL
	err = db.Execute("INVALID SQL")
	if err == nil {
		t.Error("Expected error for invalid SQL")
	}
}