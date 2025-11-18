package main

import (
	"fmt"
	"strconv"
	"strings"
)

// DataType 定义数据类型
type DataType int

const (
	TypeInt DataType = iota
	TypeVarchar
	TypeText
	TypeFloat
	TypeBool
)

// Column 定义表列结构
type Column struct {
	Name     string
	Type     DataType
	Nullable bool
}

// Row 表示一行数据
type Row struct {
	Values []interface{}
}

// Table 定义表结构
type Table struct {
	Name    string
	Columns []Column
	Rows    []Row
}

// String 返回数据类型的字符串表示
func (dt DataType) String() string {
	switch dt {
	case TypeInt:
		return "INT"
	case TypeVarchar:
		return "VARCHAR"
	case TypeText:
		return "TEXT"
	case TypeFloat:
		return "FLOAT"
	case TypeBool:
		return "BOOL"
	default:
		return "UNKNOWN"
	}
}

// ParseDataType 从字符串解析数据类型
func ParseDataType(s string) (DataType, error) {
	switch strings.ToUpper(s) {
	case "INT", "INTEGER":
		return TypeInt, nil
	case "VARCHAR":
		return TypeVarchar, nil
	case "TEXT":
		return TypeText, nil
	case "FLOAT", "DOUBLE":
		return TypeFloat, nil
	case "BOOL", "BOOLEAN":
		return TypeBool, nil
	default:
		return TypeInt, fmt.Errorf("unknown data type: %s", s)
	}
}

// ParseValue 根据数据类型解析值
func ParseValue(value string, dataType DataType) (interface{}, error) {
	// 去除引号
	value = strings.TrimSpace(value)
	if (dataType == TypeVarchar || dataType == TypeText) && 
	   strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
		value = value[1 : len(value)-1]
	}

	switch dataType {
	case TypeInt:
		return strconv.Atoi(value)
	case TypeFloat:
		return strconv.ParseFloat(value, 64)
	case TypeBool:
		lower := strings.ToLower(value)
		if lower == "true" || lower == "1" {
			return true, nil
		} else if lower == "false" || lower == "0" {
			return false, nil
		}
		return false, fmt.Errorf("invalid boolean value: %s", value)
	case TypeVarchar, TypeText:
		return value, nil
	default:
		return nil, fmt.Errorf("unsupported data type: %d", dataType)
	}
}

// ValueToString 将值转换为字符串用于存储
func ValueToString(value interface{}, dataType DataType) string {
	switch dataType {
	case TypeInt, TypeFloat, TypeBool:
		return fmt.Sprintf("%v", value)
	case TypeVarchar, TypeText:
		return fmt.Sprintf("'%v'", value)
	default:
		return fmt.Sprintf("%v", value)
	}
}