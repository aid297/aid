package excelV2

import (
	"fmt"
	"sync"

	"github.com/xuri/excelize/v2"
)

type (
	IRow interface {
		SetCells(cells ...ICell) IRow
		AppendCells(cells ...ICell) IRow
		GetCells() []ICell
		SetRowNum(rowNum uint) IRow
		GetRowNum() uint
	}

	Row struct {
		Error  error
		mu     sync.RWMutex
		rowNum uint
		cells  []ICell
	}
)

// NewRow 新建行数据
func NewRow(cells ...ICell) IRow { return &Row{cells: cells, mu: sync.RWMutex{}} }

// NewRowByNum 通过行号新建行数据
func NewRowByNum(rowNum uint, cells ...ICell) IRow { return NewRow(cells...).SetRowNum(rowNum) }

// SetCells 设置 cells
func (my *Row) SetCells(cells ...ICell) IRow {
	my.mu.Lock()
	defer my.mu.Unlock()

	var (
		err error
		col string
	)

	for idx := range cells {
		if col, err = excelize.ColumnNumberToName(idx + 1); err != nil {
			my.Error = fmt.Errorf("设置 cells-生成列失败：%w", err)
			return my
		}
		cells[idx].SetAttrs(Coordinate(col)).SetRowNum(my.rowNum)
	}

	my.cells = cells
	return my
}

// AppendCells 追加 cells
func (my *Row) AppendCells(cells ...ICell) IRow {
	my.mu.Lock()
	defer my.mu.Unlock()

	var (
		err error
		col string
	)

	for idx := range cells {
		if col, err = excelize.ColumnNumberToName(idx + 1 + len(my.cells)); err != nil {
			my.Error = fmt.Errorf("设置 cells-生成列失败：%w", err)
			return my
		}
		cells[idx].SetAttrs(Coordinate(col)).SetRowNum(my.rowNum)
	}

	my.cells = append(my.cells, cells...)
	return my
}

// GetCells 获取 cells
func (my *Row) GetCells() []ICell {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.cells
}

// SetRowNum 设置行号
func (my *Row) SetRowNum(rowNum uint) IRow {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.rowNum = rowNum

	var (
		err error
		col string
	)

	for idx := range my.cells {
		if col, err = excelize.ColumnNumberToName(idx + 1); err != nil {
			my.Error = fmt.Errorf("设置 cells-生成列失败：%w", err)
			return my
		}
		my.cells[idx].SetAttrs(Coordinate(col)).SetRowNum(my.rowNum)
	}

	return my
}

// GetRowNum 获取行号
func (my *Row) GetRowNum() uint {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.rowNum
}
