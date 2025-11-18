package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// StorageEngine 存储引擎
type StorageEngine struct {
	dataDir string
}

// NewStorageEngine 创建新的存储引擎
func NewStorageEngine(dataDir string) (*StorageEngine, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %v", err)
	}
	
	return &StorageEngine{
		dataDir: dataDir,
	}, nil
}

// SaveTable 保存表结构和数据到文件
func (s *StorageEngine) SaveTable(table *Table) error {
	// 保存表结构
	if err := s.saveTableSchema(table); err != nil {
		return fmt.Errorf("failed to save table schema: %v", err)
	}
	
	// 保存表数据
	if err := s.saveTableData(table); err != nil {
		return fmt.Errorf("failed to save table data: %v", err)
	}
	
	return nil
}

// saveTableSchema 保存表结构到文件
func (s *StorageEngine) saveTableSchema(table *Table) error {
	schemaPath := filepath.Join(s.dataDir, table.Name+".schema.json")
	
	schema := map[string]interface{}{
		"name":    table.Name,
		"columns": table.Columns,
	}
	
	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal schema: %v", err)
	}
	
	if err := os.WriteFile(schemaPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write schema file: %v", err)
	}
	
	return nil
}

// saveTableData 保存表数据到文件
func (s *StorageEngine) saveTableData(table *Table) error {
	dataPath := filepath.Join(s.dataDir, table.Name+".data")
	
	file, err := os.Create(dataPath)
	if err != nil {
		return fmt.Errorf("failed to create data file: %v", err)
	}
	defer file.Close()
	
	// 写入行数
	if err := binary.Write(file, binary.LittleEndian, uint32(len(table.Rows))); err != nil {
		return fmt.Errorf("failed to write row count: %v", err)
	}
	
	// 写入每一行数据
	for _, row := range table.Rows {
		if err := s.writeRow(file, row, table.Columns); err != nil {
			return fmt.Errorf("failed to write row: %v", err)
		}
	}
	
	return nil
}

// writeRow 写入一行数据
func (s *StorageEngine) writeRow(file *os.File, row Row, columns []Column) error {
	// 写入列数
	if err := binary.Write(file, binary.LittleEndian, uint32(len(row.Values))); err != nil {
		return err
	}
	
	// 写入每一列的数据
	for i, value := range row.Values {
		if err := s.writeValue(file, value, columns[i].Type); err != nil {
			return err
		}
	}
	
	return nil
}

// writeValue 写入单个值
func (s *StorageEngine) writeValue(file *os.File, value interface{}, dataType DataType) error {
	// 写入数据类型
	if err := binary.Write(file, binary.LittleEndian, uint8(dataType)); err != nil {
		return err
	}
	
	switch v := value.(type) {
	case int:
		if err := binary.Write(file, binary.LittleEndian, int64(v)); err != nil {
			return err
		}
	case int64:
		if err := binary.Write(file, binary.LittleEndian, v); err != nil {
			return err
		}
	case float64:
		if err := binary.Write(file, binary.LittleEndian, v); err != nil {
			return err
		}
	case bool:
		var boolValue uint8
		if v {
			boolValue = 1
		}
		if err := binary.Write(file, binary.LittleEndian, boolValue); err != nil {
			return err
		}
	case string:
		// 写入字符串长度
		if err := binary.Write(file, binary.LittleEndian, uint32(len(v))); err != nil {
			return err
		}
		// 写入字符串内容
		if _, err := file.WriteString(v); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported value type: %T", value)
	}
	
	return nil
}

// LoadTable 从文件加载表结构和数据
func (s *StorageEngine) LoadTable(tableName string) (*Table, error) {
	// 加载表结构
	schemaPath := filepath.Join(s.dataDir, tableName+".schema.json")
	schemaData, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %v", err)
	}
	
	var schema struct {
		Name    string   `json:"name"`
		Columns []Column `json:"columns"`
	}
	
	if err := json.Unmarshal(schemaData, &schema); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schema: %v", err)
	}
	
	table := &Table{
		Name:    schema.Name,
		Columns: schema.Columns,
		Rows:    make([]Row, 0),
	}
	
	// 加载表数据
	dataPath := filepath.Join(s.dataDir, tableName+".data")
	file, err := os.Open(dataPath)
	if err != nil {
		// 如果数据文件不存在，返回空表
		return table, nil
	}
	defer file.Close()
	
	// 读取行数
	var rowCount uint32
	if err := binary.Read(file, binary.LittleEndian, &rowCount); err != nil {
		return nil, fmt.Errorf("failed to read row count: %v", err)
	}
	
	// 读取每一行数据
	for i := uint32(0); i < rowCount; i++ {
		row, err := s.readRow(file, table.Columns)
		if err != nil {
			return nil, fmt.Errorf("failed to read row %d: %v", i, err)
		}
		table.Rows = append(table.Rows, row)
	}
	
	return table, nil
}

// readRow 从文件读取一行数据
func (s *StorageEngine) readRow(file *os.File, columns []Column) (Row, error) {
	// 读取列数
	var colCount uint32
	if err := binary.Read(file, binary.LittleEndian, &colCount); err != nil {
		return Row{}, err
	}
	
	if int(colCount) != len(columns) {
		return Row{}, fmt.Errorf("column count mismatch: expected %d, got %d", len(columns), colCount)
	}
	
	// 读取每一列的数据
	values := make([]interface{}, colCount)
	for i := uint32(0); i < colCount; i++ {
		value, err := s.readValue(file)
		if err != nil {
			return Row{}, err
		}
		values[i] = value
	}
	
	return Row{Values: values}, nil
}

// readValue 从文件读取单个值
func (s *StorageEngine) readValue(file *os.File) (interface{}, error) {
	// 读取数据类型
	var dataType uint8
	if err := binary.Read(file, binary.LittleEndian, &dataType); err != nil {
		return nil, err
	}
	
	switch DataType(dataType) {
	case TypeInt:
		var value int64
		if err := binary.Read(file, binary.LittleEndian, &value); err != nil {
			return nil, err
		}
		return int(value), nil
	case TypeFloat:
		var value float64
		if err := binary.Read(file, binary.LittleEndian, &value); err != nil {
			return nil, err
		}
		return value, nil
	case TypeBool:
		var value uint8
		if err := binary.Read(file, binary.LittleEndian, &value); err != nil {
			return nil, err
		}
		return value != 0, nil
	case TypeVarchar, TypeText:
		// 读取字符串长度
		var length uint32
		if err := binary.Read(file, binary.LittleEndian, &length); err != nil {
			return nil, err
		}
		
		// 读取字符串内容
		str := make([]byte, length)
		if _, err := file.Read(str); err != nil {
			return nil, err
		}
		
		return string(str), nil
	default:
		return nil, fmt.Errorf("unsupported data type: %d", dataType)
	}
}