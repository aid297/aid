package excelV2

import (
	"testing"
	"time"
)

func TestRow1(t *testing.T) {
	rows := []IRow{
		NewRowByNum(
			1,
			NewCellInt(1),
			NewCell("张三"),
			NewCellBool(true),
			NewCellTime(time.Date(2025, 1, 2, 0, 0, 0, 0, time.Local)),
		),
		NewRowByNum(
			2,
			NewCellInt(2),
			NewCell("李四"),
			NewCellBool(false),
			NewCellTime(time.Date(2024, 3, 4, 0, 0, 0, 0, time.Local)),
		),
		NewRowByNum(
			3,
			NewCellInt(3),
			NewCell("王五"),
			NewCellBool(true),
			NewCellTime(time.Date(2023, 5, 6, 0, 0, 0, 0, time.Local)),
		),
		NewRowByNum(
			4,
			NewCellInt(4),
			NewCell("赵六"),
			NewCellBool(false),
			NewCellTime(time.Date(2022, 7, 8, 0, 0, 0, 0, time.Local)),
		),
		NewRowByNum(
			5,
			NewCellInt(5),
			NewCell("孙七"),
			NewCellBool(true),
			NewCellTime(time.Date(2021, 9, 10, 0, 0, 0, 0, time.Local)),
		),
	}

	for _, row := range rows {
		t.Logf("行：%d", row.GetRowNum())
		for _, cell := range row.GetCells() {
			t.Logf("%v\t%v", cell.GetCoordinate(), cell.GetContent())
		}
	}
}
