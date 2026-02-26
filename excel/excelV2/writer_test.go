package excelV2

import (
	"testing"
	"time"

	"github.com/aid297/aid/filesystem/filesystemV4"
)

func TestWriter1(t *testing.T) {
	excelByFile := NewWriter().
		SetFilename(File(filesystemV4.NewFile(filesystemV4.Rel("./test-by-file.xlsx")))).
		SetSheet(SheetName("Sheet 1")) // 通过名称选择一个工作表

	if err := excelByFile.Save(); err != nil {
		t.Errorf("保存失败(by file)：%v\n", err)
	}

	excelByFilename := NewWriter().
		SetFilename(Filename("./test-by-filename.xlsx")).
		SetSheet(CreateSheet("Sheet B")) // 创建一个新的工作表

	if err := excelByFilename.Save(); err != nil {
		t.Errorf("保存失败(by filename)：%v\n", err)
	}

	t.Logf("OK")
}

func TestWriter2(t *testing.T) {
	var err error

	// 设置三行数据
	rows := NewRows(
		3, // 从第三行开始
		NewRow(
			NewCellInt(1),
			NewCell("张三"),
			NewCellBool(true),
			NewCellTime(time.Date(2025, 1, 2, 0, 0, 0, 0, time.Local)),
		),
		NewRow(
			NewCellInt(2),
			NewCell("李四"),
			NewCellBool(false),
			NewCellTime(time.Date(2024, 3, 4, 0, 0, 0, 0, time.Local)),
		),
		NewRow(
			NewCellInt(3),
			NewCell("王五"),
			NewCellBool(true),
			NewCellTime(time.Date(2023, 5, 6, 0, 0, 0, 0, time.Local)),
		),
	)

	// 追加两行数据
	rows.AppendRows(
		NewRow(
			NewCellInt(4),
			NewCell("赵六"),
			NewCellBool(false),
			NewCellTime(time.Date(2022, 7, 8, 0, 0, 0, 0, time.Local)),
		),
		NewRow(
			NewCellInt(5),
			NewCell("孙七"),
			NewCellBool(true),
			NewCellTime(time.Date(2021, 9, 10, 0, 0, 0, 0, time.Local)),
		),
	)

	excelWriter := NewWriter().SetFilename(Filename("./test.xlsx")).
		SetSheet(SheetIndex(0)). // 通过索引选择一个工作表
		Write(rows)

	if err = excelWriter.Save(); err != nil {
		t.Errorf("保存失败：%v\n", err)
	}

	t.Logf("OK")
}

func TestWriter3(t *testing.T) {
	rows := NewRows(
		5, // 从第5行开始
		NewRow(
			NewCellInt(
				1, // 设置内容
				Font(CellFontOpt{
					Family:     "宋体",
					Bold:       true,
					Italic:     false,
					RGB:        "red",
					PatternRGB: "pink",
					Size:       15,
				}), // 设置字体样式
				Border(
					CellBorderRGBOpt{
						Top:          "green",
						Bottom:       "black",
						Left:         "white",
						Right:        "red",
						DiagonalUp:   "purple",
						DiagonalDown: "blue",
					},
					CellBorderStyleOpt{
						Top:          1,
						Bottom:       2,
						Left:         3,
						Right:        4,
						DiagonalUp:   5,
						DiagonalDown: 6,
					},
				), // 设置边框样式
				Alignment(CellAlignmentOpt{
					Horizontal: "right",
					Vertical:   "bottom",
					WrapText:   true,
				}), // 设置字体对齐
			),
			NewCell("张三"),
			NewCellBool(true),
			NewCellTime(time.Date(2025, 1, 2, 0, 0, 0, 0, time.Local)),
		),
	)

	excelWriter := NewWriter().SetFilename(Filename("./test.xlsx")).
		SetSheet(SheetIndex(0)). // 通过索引选择一个工作表
		Write(rows)

	if err := excelWriter.Save(); err != nil {
		t.Errorf("保存失败：%v\n", err)
	}

	t.Logf("OK")
}
