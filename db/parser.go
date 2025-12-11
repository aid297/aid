package main

import (
	"fmt"
	"regexp"
	"strings"
)

// SQLParser SQL解析器
type SQLParser struct {
	Tables map[string]*Table
}

// NewSQLParser 创建新的SQL解析器
func NewSQLParser() *SQLParser {
	return &SQLParser{
		Tables: make(map[string]*Table),
	}
}

// ParseAndExecute 解析并执行SQL语句
func (p *SQLParser) ParseAndExecute(sql string) error {
	sql = strings.TrimSpace(sql)

	// 转换为大写以便判断语句类型
	upperSQL := strings.ToUpper(sql)

	if strings.HasPrefix(upperSQL, "CREATE TABLE") {
		return p.parseCreateTable(sql)
	} else if strings.HasPrefix(upperSQL, "INSERT INTO") {
		return p.parseInsert(sql)
	} else if strings.HasPrefix(upperSQL, "SELECT") {
		return p.parseSelect(sql)
	} else {
		return fmt.Errorf("unsupported SQL statement: %s", sql)
	}
}

// parseCreateTable 解析CREATE TABLE语句
func (p *SQLParser) parseCreateTable(sql string) error {
	// 简化的CREATE TABLE解析，格式: CREATE TABLE table_name (col1 type1, col2 type2, ...)
	re := regexp.MustCompile(`CREATE\s+TABLE\s+(\w+)\s*\((.+)\)`)
	matches := re.FindStringSubmatch(sql)
	if len(matches) != 3 {
		return fmt.Errorf("invalid CREATE TABLE syntax: %s", sql)
	}

	tableName := matches[1]
	columnsDef := matches[2]

	// 解析列定义
	columnsStr := strings.Split(columnsDef, ",")
	var columns []Column

	for _, colStr := range columnsStr {
		colStr = strings.TrimSpace(colStr)
		parts := strings.Fields(colStr)
		if len(parts) < 2 {
			return fmt.Errorf("invalid column definition: %s", colStr)
		}

		colName := parts[0]
		colType, err := ParseDataType(parts[1])
		if err != nil {
			return fmt.Errorf("invalid column type: %v", err)
		}

		// 检查是否可为空（简化处理）
		nullable := true
		for _, part := range parts[2:] {
			if strings.ToUpper(part) == "NOT" && len(parts) > 3 && strings.ToUpper(parts[3]) == "NULL" {
				nullable = false
				break
			}
		}

		columns = append(columns, Column{
			Name:     colName,
			Type:     colType,
			Nullable: nullable,
		})
	}

	// 创建表
	p.Tables[tableName] = &Table{
		Name:    tableName,
		Columns: columns,
		Rows:    make([]Row, 0),
	}

	fmt.Printf("Table '%s' created successfully\n", tableName)
	return nil
}

// parseInsert 解析INSERT语句
func (p *SQLParser) parseInsert(sql string) error {
	// 简化的INSERT解析，格式: INSERT INTO table_name VALUES (val1, val2, ...)
	re := regexp.MustCompile(`INSERT\s+INTO\s+(\w+)\s+VALUES\s*\((.+)\)`)
	matches := re.FindStringSubmatch(sql)
	if len(matches) != 3 {
		return fmt.Errorf("invalid INSERT syntax: %s", sql)
	}

	tableName := matches[1]
	valuesStr := matches[2]

	// 检查表是否存在
	table, exists := p.Tables[tableName]
	if !exists {
		return fmt.Errorf("table '%s' does not exist", tableName)
	}

	// 解析值
	valuesStr = strings.TrimSpace(valuesStr)
	valueStrings := strings.Split(valuesStr, ",")

	if len(valueStrings) != len(table.Columns) {
		return fmt.Errorf("column count mismatch: expected %d, got %d", len(table.Columns), len(valueStrings))
	}

	var values []interface{}
	for i, valueStr := range valueStrings {
		valueStr = strings.TrimSpace(valueStr)
		value, err := ParseValue(valueStr, table.Columns[i].Type)
		if err != nil {
			return fmt.Errorf("error parsing value '%s': %v", valueStr, err)
		}
		values = append(values, value)
	}

	// 添加行到表
	table.Rows = append(table.Rows, Row{Values: values})

	fmt.Printf("Rows inserted into table '%s' successfully\n", tableName)
	return nil
}

// parseSelect 解析SELECT语句
func (p *SQLParser) parseSelect(sql string) error {
	// 简化的SELECT解析，格式: SELECT * FROM table_name
	re := regexp.MustCompile(`SELECT\s+\*\s+FROM\s+(\w+)`)
	matches := re.FindStringSubmatch(sql)
	if len(matches) != 2 {
		return fmt.Errorf("invalid SELECT syntax: %s", sql)
	}

	tableName := matches[1]
	table, exists := p.Tables[tableName]
	if !exists {
		return fmt.Errorf("table '%s' does not exist", tableName)
	}

	// 打印表头
	header := make([]string, len(table.Columns))
	for i, col := range table.Columns {
		header[i] = col.Name
	}
	fmt.Println(strings.Join(header, " | "))

	// 打印分隔线
	separator := make([]string, len(table.Columns))
	for i := range table.Columns {
		separator[i] = strings.Repeat("-", len(header[i]))
	}
	fmt.Println(strings.Join(separator, "-+-"))

	// 打印数据行
	for _, row := range table.Rows {
		rowStr := make([]string, len(row.Values))
		for i, value := range row.Values {
			rowStr[i] = ValueToString(value, table.Columns[i].Type)
		}
		fmt.Println(strings.Join(rowStr, " | "))
	}

	return nil
}
