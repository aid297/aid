package excel

import (
	"fmt"

	"github.com/aid297/aid/array"

	"github.com/xuri/excelize/v2"
)

// Row Excel行
type Row struct {
	Err       error
	cells     *array.AnyArray[*Cell]
	rowNumber uint64
}

var RowApp Row

func (*Row) New() *Row { return NewRow() }

// NewRow 构造函数
//
//go:fix 推荐使用New方法
func NewRow() *Row { return &Row{} }

// GetCells 获取单元格组
func (my *Row) GetCells() *array.AnyArray[*Cell] { return my.cells }

// SetCells 设置单元格组
func (my *Row) SetCells(cells []*Cell) *Row {
	if my.GetRowNumber() == 0 {
		my.Err = SetCellErr.New("行标必须大于0")
		return my
	}

	for colNumber, cell := range cells {
		if colText, err := excelize.ColumnNumberToName(colNumber + 1); err != nil {
			my.Err = SetCellErr.Wrap(fmt.Errorf("列索引转列文字失败：%d，%d", my.GetRowNumber(), colNumber+1))
			return my
		} else {
			cell.SetCoordinate(fmt.Sprintf("%s%d", colText, my.GetRowNumber()))
		}
	}
	my.cells = array.New(cells)

	return my
}

// GetRowNumber 获取行标
func (my *Row) GetRowNumber() uint64 { return my.rowNumber }

// SetRowNumber 设置行标
func (my *Row) SetRowNumber(rowNumber uint64) *Row {
	my.rowNumber = rowNumber

	return my
}
