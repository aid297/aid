package excelV2

import "github.com/xuri/excelize/v2"

// ColumnNumberToText 列索引转文字
func ColumnNumberToText(columnNumber int) (string, error) {
	return excelize.ColumnNumberToName(columnNumber)
}

// ColumnTextToNumber 列文字转索引
func ColumnTextToNumber(columnText string) int {
	result := 0
	for i, char := range columnText {
		result += (int(char - 'A' + 1)) * pow(26, len(columnText)-i-1)
	}

	return result
}

// pow 是一个简单的幂函数计算，用于26进制转换
func pow(base, exponent int) int {
	result := 1
	for range exponent {
		result *= base
	}

	return result
}
