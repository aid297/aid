package excel

import (
	"time"

	"github.com/aid297/aid/dict"

	"github.com/xuri/excelize/v2"

	"github.com/aid297/aid/str"
)

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

func WriteDemo(filename string) {
	err := NewWriter(filename).
		CreateSheet("Sheet1").
		SetTitleRow([]string{"username", "nickname", "score"}, 1).
		SetRows([]*Row{
			NewRow().
				SetRowNumber(2).
				SetCells([]*Cell{
					NewCellAny("zhangsan"),
					NewCellAny("张三"),
					NewCellAny(100).
						SetFontRgbFunc(func() (string, bool) {
							if 100 > 80 {
								return "FF0000", true
							} else {
								return "", false
							}
						}).
						SetFontSize(28, true).
						SetFontBold(true, true).
						SetFontItalic(true, true),
				}),
			NewRow().
				SetRowNumber(3).
				SetCells([]*Cell{
					NewCellAny("lisi"),
					NewCellAny("李四"),
					NewCellAny(90).
						SetFontRgbFunc(func() (string, bool) {
							if 90 > 80 {
								return "FF0000", true
							} else {
								return "", false
							}
						}).
						SetFontSize(28, true).
						SetFontBold(true, true).
						SetFontItalic(true, true),
				}),
			NewRow().
				SetRowNumber(4).
				SetCells([]*Cell{
					NewCellAny("wangwu"),
					NewCellAny("王五"),
					NewCellAny(80).
						SetFontRgbFunc(func() (string, bool) {
							if 80 > 90 {
								return "FF0000", true
							} else {
								return "", false
							}
						}).
						SetFontSize(28, false).
						SetFontBold(true, false).
						SetFontItalic(true, false),
				}),
			NewRow().
				SetRowNumber(5).
				SetCells([]*Cell{
					NewCellAny("zhaoliu"),
					NewCellAny("赵六").
						SetPatternRgb("#00FF00", true),
					NewCellTime(time.Now()).
						SetFontRgbFunc(func() (string, bool) {
							if 70 > 80 {
								return "FF0000", true
							} else {
								return "", false
							}
						}).
						SetFontSize(28, false).
						SetFontBold(true, false).
						SetFontItalic(true, false).
						SetBorderSurrounding("0000FF", 1, true),
				}),
		}).
		Save()
	if err != nil {
		str.NewTerminalLog("保存excel失败 %v").Error(err)
	}
}

func ReadDemo(filename string) {
	excelData, err := NewReader().
		OpenFile(filename).
		SetOriginalRow(2).
		SetTitleRow(1).
		SetSheetName("Sheet1").
		ReadTitle().
		Read().
		DataWithTitle()
	if err != nil {
		str.NewTerminalLog("err: %v").Error(err)
	}

	excelData.Each(func(key uint64, value *dict.AnyDict[string, string]) {
		username := value.GetValueByKey("username")
		nickname := value.GetValueByKey("nickname")
		score := value.GetValueByKey("score")

		str.NewTerminalLog("%d行: 姓名[%s]，昵称[%s]，分数[%s]").Success(key, username, nickname, score)
	})
}
