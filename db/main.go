package main

import (
	"bufio"
	"fmt"
	"os"
)

// Database 数据库结构
type Database struct {
	parser        *SQLParser
	storageEngine *StorageEngine
	dataDir       string
}

// NewDatabase 创建新的数据库实例
func NewDatabase(dataDir string) (*Database, error) {
	storageEngine, err := NewStorageEngine(dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage engine: %v", err)
	}

	db := &Database{
		parser:        NewSQLParser(),
		storageEngine: storageEngine,
		dataDir:       dataDir,
	}

	// 加载已存在的表
	if err := db.loadExistingTables(); err != nil {
		return nil, fmt.Errorf("failed to load existing tables: %v", err)
	}

	return db, nil
}

// loadExistingTables 加载已存在的表
func (db *Database) loadExistingTables() error {
	// 列出所有schema文件
	entries, err := os.ReadDir(db.dataDir)
	if err != nil {
		if os.IsNotExist(err) {
			// 目录不存在，这不是错误
			return nil
		}
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() && len(entry.Name()) > 12 && entry.Name()[len(entry.Name())-12:] == ".schema.json" {
			tableName := entry.Name()[:len(entry.Name())-12]
			table, err := db.storageEngine.LoadTable(tableName)
			if err != nil {
				return fmt.Errorf("failed to load table '%s': %v", tableName, err)
			}
			db.parser.Tables[tableName] = table
		}
	}

	return nil
}

// Execute 执行SQL语句
func (db *Database) Execute(sql string) error {
	if err := db.parser.ParseAndExecute(sql); err != nil {
		return err
	}

	// 自动保存表数据
	for tableName, table := range db.parser.Tables {
		if err := db.storageEngine.SaveTable(table); err != nil {
			return fmt.Errorf("failed to save table '%s': %v", tableName, err)
		}
	}

	return nil
}

// ListTables 列出所有表
func (db *Database) ListTables() {
	fmt.Println("Available tables:")
	for tableName := range db.parser.Tables {
		fmt.Printf("- %s\n", tableName)
	}
}

// ShowTableSchema 显示表结构
func (db *Database) ShowTableSchema(tableName string) error {
	table, exists := db.parser.Tables[tableName]
	if !exists {
		return fmt.Errorf("table '%s' does not exist", tableName)
	}

	fmt.Printf("Table: %s\n", table.Name)
	fmt.Println("Columns:")
	for _, col := range table.Columns {
		nullable := "NULL"
		if !col.Nullable {
			nullable = "NOT NULL"
		}
		fmt.Printf("- %s %s %s\n", col.Name, col.Type.String(), nullable)
	}

	fmt.Printf("Rows: %d\n", len(table.Rows))
	return nil
}

// RunInteractive 运行交互式命令行界面
func (db *Database) RunInteractive() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Simple MySQL-like Database")
	fmt.Println("Type 'help' for available commands or 'exit' to quit")
	fmt.Println()

	for {
		fmt.Print("db> ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		if input == "" {
			continue
		}

		if input == "exit" || input == "quit" {
			break
		} else if input == "help" {
			fmt.Println("Available commands:")
			fmt.Println("- CREATE TABLE table_name (col1 type1, col2 type2, ...)")
			fmt.Println("- INSERT INTO table_name VALUES (val1, val2, ...)")
			fmt.Println("- SELECT * FROM table_name")
			fmt.Println("- LIST TABLES")
			fmt.Println("- DESC table_name")
			fmt.Println("- EXIT")
			continue
		} else if input == "LIST TABLES" {
			db.ListTables()
			continue
		} else if len(input) > 5 && input[:5] == "DESC " {
			tableName := input[5:]
			if err := db.ShowTableSchema(tableName); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			continue
		}

		if err := db.Execute(input); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Input error: %v\n", err)
	}
}

// 主函数
func main() {
	// 创建数据库实例
	db, err := NewDatabase("./data")
	if err != nil {
		fmt.Printf("Failed to create database: %v\n", err)
		return
	}

	// 运行交互式界面
	db.RunInteractive()
}
